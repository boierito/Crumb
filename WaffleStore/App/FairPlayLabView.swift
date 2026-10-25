import SwiftUI
import UniformTypeIdentifiers
import MapleSAP
import CryptoKit

nonisolated enum FairPlayNative {
    static func run(_ operation: String, source: URL? = nil, destination: URL? = nil,
                    report: URL? = nil, model: URL? = nil) throws -> String {
        let request = ["Operation": operation, "Source": source?.path ?? "",
                       "Destination": destination?.path ?? "", "Report": report?.path ?? "", "Model": model?.path ?? ""]
        let data = try JSONEncoder().encode(request)
        let text = String(decoding: data, as: UTF8.self)
        var output: UnsafeMutablePointer<UInt8>?
        var count = 0
        let status = text.withCString { WaffleFairPlayLab(UnsafeMutablePointer(mutating: $0), &output, &count) }
        defer { WaffleSAPFree(output, count) }
        guard status == 0, let output, count > 0, count < 1 << 20 else {
            throw NSError(domain: "Crumb.FairPlayLab", code: Int(status))
        }
        let result = Data(bytes: output, count: count)
        let object = try JSONSerialization.jsonObject(with: result)
        let pretty = try JSONSerialization.data(withJSONObject: object, options: [.prettyPrinted, .sortedKeys])
        return String(decoding: pretty, as: UTF8.self)
    }
    static func privateDirectory(_ url: URL) throws {
        try FileManager.default.createDirectory(at: url, withIntermediateDirectories: true,
            attributes: [.posixPermissions: 0o700, .protectionKey: FileProtectionType.complete])
        var resource = url
        var values = URLResourceValues(); values.isExcludedFromBackup = true
        try resource.setResourceValues(values)
    }
    static func privateWrite(_ data: Data, to url: URL) throws {
        try data.write(to: url, options: [.atomic, .completeFileProtection])
        try FileManager.default.setAttributes([.posixPermissions: 0o600], ofItemAtPath: url.path)
    }
}

@MainActor
final class FairPlayLab: ObservableObject {
    @Published var busy = false
    @Published var logs = ["Components": "Pending", "Transaction": "Pending", "Inspection": "Pending", "Decryption": "Blocked: import a compatible model after selecting an IPA."]
    @Published var output: URL?
    @Published var diagnostic: URL?
    @Published var captureArchive: URL?
    @Published var input: URL?
    let root = FileManager.default.temporaryDirectory.appendingPathComponent("FairPlayLab-" + UUID().uuidString)
    private var task: Task<Void, Never>?
    private var transactions = 0
    func log(_ section: String, _ text: String) {
        let line = "[\(Date().formatted(date: .omitted, time: .standard))] \(text)"
        if logs[section] == "Pending" { logs[section] = line }
        else { logs[section] = String(((logs[section] ?? "") + "\n" + line).suffix(24000)) }
    }
    func perform(_ operation: @escaping @MainActor () async throws -> Void) {
        guard !busy else { return }; busy = true
        task = Task {
            defer { busy = false; task = nil }
            do { try FairPlayNative.privateDirectory(root); try await operation() }
            catch { log("Transaction", "Operation failed (\((error as NSError).domain == NSURLErrorDomain ? "network" : "validation/storage")); no private error text logged. Retry or export diagnostic.") }
        }
    }
    func selftest() {
        perform {
            self.log("Components", "Running synthetic state/frame/SIDB vectors. No account data used.")
            let result = try await Task.detached { try FairPlayNative.run("selftest") }.value
            self.log("Components", result)
        }
    }
    func capture(_ transaction: FairPlayTransaction) throws {
        let directory = root.appendingPathComponent("captures").appendingPathComponent(UUID().uuidString)
        try FairPlayNative.privateDirectory(directory)
        var files: [[String: String]] = []
        var blobs: [(String, Data)] = [("kbsync.bin", transaction.kbsync)]
        if let keybag = transaction.keybag { blobs.append(("keybag.bin", keybag)) }
        blobs += transaction.sinfs.enumerated().map { ("sinf-\($0.offset).bin", $0.element) }
        for (name, data) in blobs {
            guard data.count <= 16 << 20 else { throw CocoaError(.fileWriteOutOfSpace) }
            try FairPlayNative.privateWrite(data, to: directory.appendingPathComponent(name))
            files.append(["name": name, "bytes": String(data.count), "sha256": SHA256.hash(data: data).map { String(format: "%02x", $0) }.joined()])
        }
        let manifest: [String: Any] = ["schema": "crumb-fairplay-capture-v1", "route": "validated-ent-download",
                                      "keybag_present": transaction.keybag != nil, "files": files]
        try FairPlayNative.privateWrite(JSONSerialization.data(withJSONObject: manifest, options: [.prettyPrinted]), to: directory.appendingPathComponent("manifest.json"))
        transactions += 1
        log("Transaction", "Validated exchange \(transactions): kbsync \(transaction.kbsync.count) bytes; keybag \(transaction.keybag?.count ?? 0) bytes; SINF count \(transaction.sinfs.count). Capture saved privately.")
    }
    func download(_ text: String, account: StoreAccount) {
        perform {
            guard self.input == nil else { return }
            self.transactions = 0
            let tool = IPATool(account: account, capture: { value in
                await MainActor.run {
                    do { try self.capture(value) }
                    catch { self.log("Transaction", "Private capture could not be saved; download may continue.") }
                }
            })
            defer { tool.close() }
            self.log("Transaction", "Looking up app and requesting latest accessible version.")
            let app = try await tool.lookup(text)
            let descriptor = try await tool.descriptor(app: app)
            if self.transactions == 0 { self.log("Transaction", "Descriptor used an alternative route: no matched ent keybag captured.") }
            let source = self.root.appendingPathComponent("download.ipa")
            _ = try await CDNDownload(destination: source, progress: { _,_,_ in }).run(descriptor.url)
            let destination = self.root.appendingPathComponent("Original.ipa")
            _ = try await Task.detached {
                try NativePackage.prepare(source: source, destination: destination, app: app, descriptor: descriptor)
            }.value
            try? FileManager.default.removeItem(at: source)
            self.input = destination
            self.log("Transaction", "IPA downloaded and prepared. Original retained until this test is erased.")
        }
    }
    func inspect() {
        guard let input else { return }
        perform {
            let result = try await Task.detached { try FairPlayNative.run("inspect", source: input) }.value
            self.log("Inspection", result)
            self.log("Inspection", "Local enrollment generation is not implemented. Keybag capture alone does not supply the missing private state input.")
        }
    }
    func decrypt(_ selected: URL) {
        guard let input else { return }
        perform {
            self.output = nil
            let model = self.root.appendingPathComponent("model.bin")
            let access = selected.startAccessingSecurityScopedResource()
            defer { if access { selected.stopAccessingSecurityScopedResource() }; try? FileManager.default.removeItem(at: model) }
            let size = try selected.resourceValues(forKeys: [.fileSizeKey]).fileSize ?? 0
            guard size > 0, size <= 1 << 20 else { throw CocoaError(.fileReadTooLarge) }
            try FairPlayNative.privateWrite(Data(contentsOf: selected), to: model)
            let output = self.root.appendingPathComponent(UUID().uuidString + ".decrypted.ipa")
            let report = self.root.appendingPathComponent("private-report.json")
            defer { try? FileManager.default.removeItem(at: report) }
            self.log("Decryption", "Checking model/SINF and processing IPA. Native failures identify the operation; private exception text is suppressed.")
            let result = try await Task.detached { try FairPlayNative.run("decrypt", source: input, destination: output, report: report, model: model) }.value
            self.log("Decryption", result)
            if let json = result.data(using: .utf8), let value = try JSONSerialization.jsonObject(with: json) as? [String: Any], value["ok"] as? Bool == true {
                self.output = output
            } else { try? FileManager.default.removeItem(at: output) }
        }
    }
    func export() {
        perform {
            let summary: [String: Any] = ["schema": "crumb-fairplay-diagnostic-v1", "ios": UIDevice.current.systemVersion,
                "app": Bundle.main.object(forInfoDictionaryKey: "CFBundleVersion") as? String ?? "unknown", "logs": self.logs]
            let file = self.root.appendingPathComponent("diagnostic.json")
            try FairPlayNative.privateWrite(JSONSerialization.data(withJSONObject: summary, options: [.prettyPrinted, .sortedKeys]), to: file)
            self.diagnostic = file
            let captures = self.root.appendingPathComponent("captures")
            if FileManager.default.fileExists(atPath: captures.path) {
                let archive = self.root.appendingPathComponent("Private-FairPlay-capture.zip")
                _ = try await Task.detached { try FairPlayNative.run("export", source: captures, destination: archive) }.value
                self.captureArchive = archive
            }
        }
    }
    func erase() {
        guard !busy else { return }
        try? FileManager.default.removeItem(at: root)
        input = nil; output = nil; diagnostic = nil; captureArchive = nil; transactions = 0
        logs = ["Components": "Pending", "Transaction": "Pending", "Inspection": "Pending", "Decryption": "Blocked: compatible model required."]
    }
}

struct FairPlayLabView: View {
    var record: DownloadRecord? = nil
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @StateObject private var lab = FairPlayLab()
    @State private var appInput = ""
    @State private var importing = false
    var body: some View {
        NavigationStack {
            Form {
                Section("1 · Check portable components") {
                    Text("Synthetic tests verify state, frame and SIDB operations on this device. They do not recover your enrollment.")
                    Button("Run self-tests") { lab.selftest() }
                }
                Section("2 · Capture a real download") {
                    Text("Sign in first. Choose a small free app available to your account. Captures contain private account-derived material.")
                    TextField("App ID, bundle ID or App Store link", text: $appInput).textInputAutocapitalization(.never).autocorrectionDisabled()
                    Button("Download and capture transaction") {
                        if let account = appData.ipaTool?.account { lab.download(appInput, account: account) }
                    }.disabled(appData.ipaTool == nil || appInput.isEmpty || lab.input != nil)
                    if let record, let url = record.fileURL {
                        Button("Use selected downloaded IPA") { lab.input = url }
                    }
                }
                Section("3 · Inspect and identify missing material") {
                    Button("Inspect IPA") { lab.inspect() }.disabled(lab.input == nil)
                    Text("Enrollment generation: not implemented. The upstream still needs private state inputs; captured keybags are not enrollment models.")
                }
                Section("4 · Test with an existing model") {
                    Button("Import model and decrypt") { importing = true }.disabled(lab.input == nil)
                    if let output = lab.output { ShareLink("Export decrypted IPA", item: output) }
                    Text("Sign the exported IPA before installation. Record whether it installs and opens; compilation and structural checks alone do not prove correct plaintext.")
                }
                ForEach(["Components", "Transaction", "Inspection", "Decryption"], id: \.self) { section in
                    Section(section + " logs") {
                        ScrollView { Text(lab.logs[section] ?? "Pending").font(.system(.caption, design: .monospaced)).frame(maxWidth: .infinity, alignment: .leading).textSelection(.enabled) }.frame(maxHeight: 220)
                    }
                }
                Section("5 · Export evidence and repeat") {
                    Button("Prepare exports") { lab.export() }
                    if let diagnostic = lab.diagnostic { ShareLink("Share diagnostic (no secret blobs)", item: diagnostic) }
                    if let capture = lab.captureArchive { ShareLink("Share PRIVATE capture — contains secrets", item: capture) }
                    Button("Copy logs") { UIPasteboard.general.string = lab.logs.keys.sorted().map { $0 + "\n" + (lab.logs[$0] ?? "") }.joined(separator: "\n\n") }
                    Button("Erase this test", role: .destructive) { lab.erase() }
                    Text("Export before closing. This session’s temporary files are removed when you close the laboratory. Download history and your Apple session are preserved.")
                }
            }
            .disabled(lab.busy)
            .overlay { if lab.busy { ProgressView("Processing…").padding().background(.regularMaterial).clipShape(RoundedRectangle(cornerRadius: 12)) } }
            .navigationTitle("FairPlay Laboratory")
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() }.disabled(lab.busy) } }
            .interactiveDismissDisabled(lab.busy)
            .fileImporter(isPresented: $importing, allowedContentTypes: [.data]) { result in
                if case .success(let url) = result { lab.decrypt(url) }
            }
            .onDisappear { lab.erase() }
        }
    }
}
