import SwiftUI
import MapleSAP

// Phase 3 only. No login, passwords, purchases or downloads are performed here.
struct SAPDiagnosticView: View {
    @State private var report = "No test performed."
    @State private var running = false
    @State private var includeAppleNetworkTest = false
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        NavigationStack {
            List {
                Section("SAP / jailed compatibility") {
                    Text("Tests memory permissions and the no-JIT interpreter. The optional SAP test downloads Apple's hash-verified guest assets into this app's cache, performs SAP setup and signs a test body. It sends no credentials. First initialization can take several minutes.")
                    Toggle("Initialize SAP and sign a test body", isOn: $includeAppleNetworkTest)
                        .disabled(running)
                    Button(running ? "Testing…" : "Run diagnostic") {
                        running = true
                        Task { await run() }
                    }
                    .disabled(running)
                }
                Section("Sanitized diagnostic") {
                    Text(report).font(.system(.caption, design: .monospaced)).textSelection(.enabled)
                    ShareLink(item: report) { Label("Export diagnostic", systemImage: "square.and.arrow.up") }
                }
            }
            .navigationTitle("SAP diagnostic")
            .toolbar { Button("Done") { dismiss() }.disabled(running) }
        }
    }

    @MainActor private func run() async {
        defer { running = false }
        let capability = MemoryCapability.probe()
        let interpreter = await Task.detached { waffle_probe_tci() }.value
        var lines = ["WaffleStore SAP probe v1", "iOS=\(UIDevice.current.systemVersion)",
            "device-family=\(UIDevice.current.model)",
            "app-build=\(Bundle.main.object(forInfoDictionaryKey: "CFBundleVersion") ?? "unknown")",
            capability.sanitizedReport,
            "tci-guest-status=\(interpreter.error)",
            "tci-guest-rax=\(interpreter.guest_rax)",
            "tci-instruction-hooks=\(interpreter.instruction_hooks)",
            "tci-test-is-sap=false"]
        var signatureGenerated = false
        do {
            let identity = try KeychainMachineIdentity.loadOrCreate()
            let repeated = try KeychainMachineIdentity.loadOrCreate()
            lines.append("keychain-identity-stable=\(identity == repeated)")
            if includeAppleNetworkTest {
                let transport = AppleSAPTransport()
                defer { transport.close() }
                let apple = SAPProtocol(transport: transport)
                let bag = try await apple.bag(identity: identity)
                lines.append("store-bag=validated; sap-version=\(bag.version)")
                let guest = try NativeSAPGuest()
                let session = try SAPSession(guest: guest, transport: transport)
                do {
                    try await session.initialize(configuration: bag, identity: identity)
                    lines.append("sap-initialization=completed")
                    let body = try PropertyListSerialization.data(fromPropertyList:
                        ["probe": "WaffleStore SAP no-login test", "guid": identity.guid], format: .xml, options: 0)
                    let signature = try await session.actionSignature(body: body)
                    lines.append("X-Apple-ActionSignature=generated; base64-length=\(signature.count); contents=withheld")
                    signatureGenerated = true
                    await session.close()
                } catch {
                    await session.close()
                    throw error
                }
            } else { lines.append("apple-network=not-requested") }
        } catch let error as SAPError {
            // SAPError only exposes fixed descriptions and numeric status codes.
            lines.append("probe-error=\(error.localizedDescription)")
        } catch {
            // Never copy arbitrary URL errors, query strings, response bodies,
            // headers or secrets into an exportable diagnostic.
            let nsError = error as NSError
            lines.append("probe-error-code=\(nsError.code)")
        }
        if !signatureGenerated { lines.append("X-Apple-ActionSignature=not-generated") }
        lines += ["sap-runtime=experimental-tci-static-library",
                  "apple-login=not-attempted",
                  "Result: a nonempty test signature does not prove login or Apple acceptance. Use a Release IPA on a physical device without a debugger."]
        report = lines.joined(separator: "\n")
    }
}
