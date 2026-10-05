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
                    Text("Tests the memory protection sequence required by ipatool's Unicorn runtime. No credentials are sent. An experimental interpreter is linked for a synthetic CPU test. The Apple SAP guest adapter is not yet implemented.")
                    Toggle("Fetch Apple Bag and SAP certificate", isOn: $includeAppleNetworkTest)
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
                _ = try await apple.certificate(configuration: bag)
                lines.append("sap-certificate=fetched-and-parsed")
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
        lines += ["sap-runtime=unavailable", "sap-initialization=not-performed",
                  "X-Apple-ActionSignature=not-generated",
                  "Result: permission probes do not prove SAP or Apple acceptance. A Release IPA on a physical device without a debugger is required."]
        report = lines.joined(separator: "\n")
    }
}
