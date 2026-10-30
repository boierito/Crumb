import SwiftUI
import Combine
import UniformTypeIdentifiers
import MapleSAP

nonisolated private struct SystemEnvelope<T: Decodable>: Decodable {
    let ok: Bool
    let result: T?
}
nonisolated private struct SystemNativePlan: Decodable {
    struct File: Decodable {
        struct Range: Decodable {
            let sliceOffset, sliceSize, cryptOffset, cryptSize, cryptidField: UInt64
            let cpuType, cpuSubtype: UInt32
        }
        let relativePath: String
        let ranges: [Range]
    }
    let schema: String
    let files: [File]
    let encryptedSlices: Int
}

@MainActor
final class SystemFairPlaySession: ObservableObject {
    @Published var busy = false
    @Published var helper: URL?
    @Published var output: URL?
    @Published var report = ""
    @Published var message = ""
    var installedTarget: URL? {
        let bundle = Bundle.main.bundleURL.appendingPathComponent("CrumbSystemPayload", isDirectory: true)
        guard let data = try? Data(contentsOf: bundle.appendingPathComponent("job.json")), data.count <= 1 << 20,
              let job = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
              job["schema"] as? String == "crumb-system-fairplay-job-v1",
              let id = job["job_id"] as? String, id.count == 32,
              id.allSatisfy({ $0.isHexDigit }) else { return nil }
        return bundle.appendingPathComponent("Target", isDirectory: true)
    }
    private func storage() throws -> URL {
        let root = try FileManager.default.url(for: .applicationSupportDirectory, in: .userDomainMask,
                                              appropriateFor: nil, create: true)
            .appendingPathComponent("CrumbSystemFairPlay", isDirectory: true)
        try FairPlayNative.privateDirectory(root)
        return root
    }
    private func workspace() throws -> URL {
        let directory = try storage().appendingPathComponent(UUID().uuidString, isDirectory: true)
        try FairPlayNative.privateDirectory(directory)
        return directory
    }
    func restore() {
        guard !busy, let root = try? storage(),
              let data = try? Data(contentsOf: root.appendingPathComponent("latest.json")),
              let saved = try? JSONSerialization.jsonObject(with: data) as? [String: String] else { return }
        func item(_ key: String) -> URL? {
            guard let path = saved[key], path.split(separator: "/").count == 2,
                  path.split(separator: "/").allSatisfy({ $0 != "." && $0 != ".." }) else { return nil }
            let url = root.appendingPathComponent(path)
            return FileManager.default.fileExists(atPath: url.path) ? url : nil
        }
        helper = item("helper"); output = item("output")
        if let path = item("report"), let text = try? String(contentsOf: path, encoding: .utf8) { report = text }
    }
    private func save(reportURL: URL? = nil) throws {
        func relative(_ url: URL) -> String { url.deletingLastPathComponent().lastPathComponent + "/" + url.lastPathComponent }
        var value = [String: String]()
        if let helper { value["helper"] = relative(helper) }
        if let output { value["output"] = relative(output) }
        if let reportURL { value["report"] = relative(reportURL) }
        try FairPlayNative.privateWrite(JSONEncoder().encode(value), to: storage().appendingPathComponent("latest.json"))
    }
    func prepare(_ original: URL) {
        guard !busy else { return }; busy = true; message = "Preparing a private Crumb package on this iPhone…"
        Task {
            defer { busy = false }
            do {
                let directory = try workspace()
                let destination = directory.appendingPathComponent("Crumb-Prepared.ipa")
                let bundle = Bundle.main.bundleURL
                let text = try await Task.detached(priority: .userInitiated) {
                    let accessing = original.startAccessingSecurityScopedResource()
                    defer { if accessing { original.stopAccessingSecurityScopedResource() } }
                    return try FairPlayNative.run("prepare-system-helper", source: bundle, destination: destination, ipa: original)
                }.value
                let object = try JSONSerialization.jsonObject(with: Data(text.utf8)) as? [String: Any]
                guard object?["ok"] as? Bool == true else { throw CocoaError(.fileReadCorruptFile) }
                helper = destination
                try save()
                message = "Package ready. Sign and install it with Crumb's existing bundle identifier. Then open Settings → FairPlay on iPhone. The package includes this app's private license; keep it private."
            } catch {
                message = "Preparation failed. No successful package was accepted; the original download is retained."
            }
        }
    }
    func decrypt(identifiers: [String]) {
        guard !busy, let source = installedTarget else { return }
        busy = true; message = "Reading all protected pages through this iPhone's FairPlay…"
        Task {
            defer { busy = false }
            do {
                let directory = try workspace()
                let destination = directory.appendingPathComponent("App-Unsigned.ipa")
                let text = try await Task.detached(priority: .userInitiated) {
                    let nativeDirectory = directory.appendingPathComponent("native", isDirectory: true)
                    try FairPlayNative.privateDirectory(nativeDirectory)
                    defer { try? FileManager.default.removeItem(at: nativeDirectory) }
                    let planText = try FairPlayNative.run("plan-system-target", source: source)
                    let decoder = JSONDecoder(); decoder.keyDecodingStrategy = .convertFromSnakeCase
                    let envelope = try decoder.decode(SystemEnvelope<SystemNativePlan>.self, from: Data(planText.utf8))
                    guard envelope.ok, let plan = envelope.result,
                          plan.schema == "crumb-system-fairplay-native-plan-v1", plan.encryptedSlices > 0 else {
                        throw CocoaError(.fileReadCorruptFile)
                    }
                    for file in plan.files {
                        let original = source.appendingPathComponent(file.relativePath)
                        let candidate = nativeDirectory.appendingPathComponent(file.relativePath)
                        try FairPlayNative.privateDirectory(candidate.deletingLastPathComponent())
                        try FileManager.default.copyItem(at: original, to: candidate)
                        try FileManager.default.setAttributes([.posixPermissions: 0o600,
                            .protectionKey: FileProtectionType.complete], ofItemAtPath: candidate.path)
                        for range in file.ranges {
                            let request = CrumbSystemRange(slice_offset: range.sliceOffset, slice_size: range.sliceSize,
                                crypt_offset: range.cryptOffset, crypt_size: range.cryptSize, cryptid_field: range.cryptidField,
                                cpu_type: range.cpuType, cpu_subtype: range.cpuSubtype)
                            var result = CrumbSystemRangeResult()
                            let code = original.path.withCString { input in
                                candidate.path.withCString { output in CrumbSystemDecryptRange(input, output, request, &result) }
                            }
                            guard code == 0, result.remap_result == 0, result.pages_total > 0,
                                  result.pages_read == result.pages_total, result.pages_changed == result.pages_total else {
                                throw NSError(domain: "Crumb.SystemFairPlay", code: Int(result.remap_errno != 0 ? result.remap_errno : result.read_error))
                            }
                        }
                    }
                    return try FairPlayNative.run("export-system-result", source: source, destination: destination,
                        model: nativeDirectory, privateIdentifiers: identifiers)
                }.value
                let object = try JSONSerialization.jsonObject(with: Data(text.utf8)) as? [String: Any]
                guard object?["ok"] as? Bool == true else { throw CocoaError(.fileReadCorruptFile) }
                report = text
                let reportURL = directory.appendingPathComponent("Export-report.json")
                try FairPlayNative.privateWrite(Data(text.utf8), to: reportURL)
                output = destination
                try save(reportURL: reportURL)
                message = "Native page processing and export checks passed. Known licensing metadata and embedded signatures were removed. Sign the output before installation. Unknown app-defined account data and plaintext authentication are not certified."
            } catch {
                message = "Native processing or export validation failed. No partial IPA was accepted. The installed target remains available for a retry."
            }
        }
    }
}

struct SystemFairPlayView: View {
    var record: DownloadRecord? = nil
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @StateObject private var session = SystemFairPlaySession()
    @State private var importing = false
    var body: some View {
        NavigationStack {
            Form {
                Section("Prepare a downloaded app") {
                    Text("Crumb packages your encrypted app and its original license inside a private copy of Crumb. This happens on the iPhone; no IPA or license is uploaded.")
                    if let original = record?.fileURL {
                        Button("Prepare this download") { session.prepare(original) }
                    }
                    Button("Choose an encrypted IPA from Files") { importing = true }
                    if let helper = session.helper {
                        ShareLink("Send prepared Crumb to your signer", item: helper)
                        Text("Sign and install this package, keeping Crumb's bundle identifier. A signing tool and a valid signing identity are currently required. Do not uninstall Crumb first.")
                            .font(.footnote)
                    }
                }
                if session.installedTarget != nil {
                    Section("Process the installed target") {
                        Text("The target is installed inside this Crumb package. The iPhone's system Apple Account must own its license.")
                        Button("Decrypt all supported executables and prepare export") {
                            var identifiers = [String]()
                            if let account = appData.ipaTool?.account {
                                identifiers = [account.email, account.dsid, account.guid, account.passwordToken]
                                    .filter { $0.utf8.count >= 6 && $0.utf8.count <= 4096 }
                            }
                            session.decrypt(identifiers: identifiers)
                        }
                    }
                }
                if let output = session.output {
                    Section("Unsigned export") {
                        ShareLink("Save processed IPA", item: output)
                        Text("Known license, purchase, receipt and signing metadata are removed. The export needs signing. This does not guarantee removal of arbitrary app-defined identifiers or grant redistribution rights.")
                            .font(.footnote)
                    }
                }
                if !session.message.isEmpty { Section { Text(session.message).textSelection(.enabled) } }
                if !session.report.isEmpty {
                    Section("Export checks") { Text(session.report).font(.system(.caption, design: .monospaced)).textSelection(.enabled) }
                }
            }
            .disabled(session.busy)
            .overlay { if session.busy { ProgressView().padding().background(.regularMaterial).clipShape(RoundedRectangle(cornerRadius: 12)) } }
            .navigationTitle("FairPlay on iPhone")
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() }.disabled(session.busy) } }
            .interactiveDismissDisabled(session.busy)
            .task { session.restore() }
            .fileImporter(isPresented: $importing, allowedContentTypes: [UTType(filenameExtension: "ipa") ?? .data]) { result in
                if case .success(let file) = result { session.prepare(file) }
            }
        }
    }
}
