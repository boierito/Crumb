import SwiftUI
import Combine
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
    @Published var logs = ["Components": "Pending", "Transaction": "Pending", "Inspection": "Pending", "Enrollment": "Pending", "Decryption": "Blocked: import a compatible model after selecting an IPA."]
    @Published var output: URL?
    @Published var diagnostic: URL?
    @Published var captureArchive: URL?
    @Published var enrollmentStages: [String: Bool] = [:]
    @Published var input: URL?
    let root = FileManager.default.temporaryDirectory.appendingPathComponent("FairPlayLab-" + UUID().uuidString)
    private var task: Task<Void, Never>?
    private var transactions = 0
    private var auditedTransactions = -1
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
    func capture(_ transaction: FairPlayTransaction, retained: NativeRetainedKBSyncGenerator) async throws {
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
        if transaction.keybag != nil {
            log("Transaction", "Running the observed native keybag consumer in the retained interpreter. This can take up to a minute.")
            let consumed = try await Task.detached { try retained.bind(transaction) }.value
            for (name, data) in [("context.json", consumed.context), ("consumer-state.json", consumed.state)] {
                try FairPlayNative.privateWrite(data, to: directory.appendingPathComponent(name))
                files.append(["name":name, "bytes":String(data.count), "sha256":SHA256.hash(data: data).map { String(format:"%02x",$0) }.joined()])
            }
            log("Transaction", consumed.summary)
        }
        let manifest: [String: Any] = ["schema": "crumb-fairplay-capture-v3", "route": "validated-ent-download",
                                      "keybag_present": transaction.keybag != nil, "files": files]
        try FairPlayNative.privateWrite(JSONSerialization.data(withJSONObject: manifest, options: [.prettyPrinted]), to: directory.appendingPathComponent("manifest.json"))
        transactions += 1
        auditedTransactions = -1
        captureArchive = nil
        diagnostic = nil
        enrollmentStages = [:]
        log("Transaction", "Validated exchange \(transactions): kbsync \(transaction.kbsync.count) bytes; keybag \(transaction.keybag?.count ?? 0) bytes; SINF count \(transaction.sinfs.count). Capture saved privately.")
    }
    func download(_ text: String, account: StoreAccount) {
        perform {
            guard self.input == nil else { return }
            self.transactions = 0
            let retained = NativeRetainedKBSyncGenerator()
            let tool = IPATool(account: account, capture: { value in
                do { try await self.capture(value, retained: retained) }
                catch { await self.log("Transaction", "Private capture or consumer failed; download may continue. No recovery success is implied.") }
            }, retainedGenerator: retained)
            defer { tool.close() }
            self.log("Transaction", "Looking up app and requesting latest accessible version.")
            let app = try await tool.lookup(text)
            let descriptor = try await tool.descriptor(app: app)
            tool.close() // Release the retained emulator before the CDN transfer.
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
    func analyzeCapture(_ selected: URL) {
        perform {
            let access = selected.startAccessingSecurityScopedResource()
            defer { if access { selected.stopAccessingSecurityScopedResource() } }
            self.log("Inspection", "Validating capture ZIP, manifest hashes and SINF structure. No account data will be logged.")
            let result = try await Task.detached { try FairPlayNative.run("analyze-capture", source: selected) }.value
            self.log("Inspection", result)
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
    // Native verification works directly on the capture created in this session.
    // The user no longer has to import their own private ZIP to analyze it.
    private func analyzeCurrentCaptureFiles() async throws {
        let captures = root.appendingPathComponent("captures")
        guard FileManager.default.fileExists(atPath: captures.path) else {
            log("Enrollment", "No new capture in this session. Download an app before analyzing.")
            return
        }
        let archive = root.appendingPathComponent("Private-FairPlay-capture.zip")
        // Once checked, the same immutable transaction does not need to be
        // zipped, inspected and logged again by Prepare exports.
        if auditedTransactions == transactions,
           captureArchive != nil,
           FileManager.default.fileExists(atPath: archive.path) {
            return
        }
        let exported = try await Task.detached {
            try FairPlayNative.run("export", source: captures, destination: archive)
        }.value
        guard let exportData = exported.data(using: .utf8),
              let exportJSON = try JSONSerialization.jsonObject(with: exportData) as? [String: Any],
              exportJSON["ok"] as? Bool == true else {
            throw NSError(domain: "Crumb.FairPlayLab.CaptureExport", code: 1)
        }
        // Go's exporter applies mode 0600; protect the archive on iOS as well.
        try FileManager.default.setAttributes(
            [.protectionKey: FileProtectionType.complete], ofItemAtPath: archive.path)
        let analysis = try await Task.detached {
            try FairPlayNative.run("analyze-capture", source: archive)
        }.value
        guard let analysisData = analysis.data(using: .utf8),
              let value = try JSONSerialization.jsonObject(with: analysisData) as? [String: Any],
              value["ok"] as? Bool == true else {
            throw NSError(domain: "Crumb.FairPlayLab.CaptureAnalysis", code: 2)
        }
        guard let analyzed = value["result"] as? [String: Any],
              let transactions = analyzed["transactions"] as? [[String: Any]],
              !transactions.isEmpty,
              transactions.count == self.transactions,
              analyzed["ready_to_decrypt"] as? Bool == false else {
            throw NSError(domain: "Crumb.FairPlayLab.IncompleteAudit", code: 3)
        }
        // Only the native analyzer's sanitized flags are surfaced here. A
        // legacy SIDB never becomes a portable frame or decrypt-ready model.
        // Alternative download routes may have no matched keybag or consumer
        // report. Represent those stages as unavailable, not as an invalid ZIP.
        let bridges = transactions.map { ($0["enrollment_bridge"] as? [String: Any]) ?? [:] }
        enrollmentStages = [
            "Native keybag import": transactions.allSatisfy { $0["consumer_import_succeeded"] as? Bool == true },
            "Legacy SIDB captured": bridges.allSatisfy { $0["legacy_sidb_observed"] as? Bool == true },
            "Portable raw state verified": bridges.allSatisfy { $0["portable_raw_state_verified"] as? Bool == true },
            "Enrollment frame verified": bridges.allSatisfy { $0["portable_enrollment_frame_verified"] as? Bool == true },
            "Portable model validated": bridges.allSatisfy { $0["validated_portable_model"] as? Bool == true }
        ]
        captureArchive = archive
        auditedTransactions = self.transactions
        log("Enrollment", "Verified current capture locally. The following report contains counts and readiness flags only, not private state bytes.")
        log("Enrollment", analysis)
    }
    func analyzeCurrentCapture() {
        perform { try await self.analyzeCurrentCaptureFiles() }
    }
    func export() {
        perform {
            // Verify the real capture first and include this sanitized stage
            // report in diagnostic.json. The ZIP is still always private.
            try await self.analyzeCurrentCaptureFiles()
            let summary: [String: Any] = ["schema": "crumb-fairplay-diagnostic-v1", "ios": UIDevice.current.systemVersion,
                "app": Bundle.main.object(forInfoDictionaryKey: "CFBundleVersion") as? String ?? "unknown", "logs": self.logs]
            let file = self.root.appendingPathComponent("diagnostic.json")
            try FairPlayNative.privateWrite(JSONSerialization.data(withJSONObject: summary, options: [.prettyPrinted, .sortedKeys]), to: file)
            self.diagnostic = file
        }
    }
    func erase() {
        guard !busy else { return }
        try? FileManager.default.removeItem(at: root)
        input = nil; output = nil; diagnostic = nil; captureArchive = nil; transactions = 0
        auditedTransactions = -1
        enrollmentStages = [:]
        logs = ["Components": "Pending", "Transaction": "Pending", "Inspection": "Pending", "Enrollment": "Pending", "Decryption": "Blocked: compatible model required."]
    }
}

struct FairPlayLabView: View {
    var record: DownloadRecord? = nil
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @StateObject private var lab = FairPlayLab()
    @State private var appInput = ""
    @State private var importing = false
    @State private var importingCapture = false
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
                    Button("Verify current capture — no import required") { lab.analyzeCurrentCapture() }
                        .disabled(lab.input == nil)
                    Button("Import an older capture ZIP (optional)") { importingCapture = true; importing = true }
                    Text("A zero StoreAgent status means keybag import was accepted. It does not produce the portable enrollment frame or raw-state source needed for decryption.")
                }
                if !lab.enrollmentStages.isEmpty {
                    Section("Enrollment readiness") {
                        ForEach(["Native keybag import", "Legacy SIDB captured",
                                 "Portable raw state verified", "Enrollment frame verified",
                                 "Portable model validated"], id: \.self) { stage in
                            HStack {
                                Text(stage)
                                Spacer()
                                Image(systemName: lab.enrollmentStages[stage] == true
                                      ? "checkmark.circle.fill" : "minus.circle")
                                    .foregroundStyle(lab.enrollmentStages[stage] == true
                                                     ? Color.green : Color.orange)
                            }
                        }
                        Text("Native import and a nonempty legacy SIDB do not establish portable enrollment or decrypted plaintext.")
                            .font(.footnote)
                    }
                }
                Section("4 · Test with an existing model") {
                    Button("Import model and decrypt") { importingCapture = false; importing = true }.disabled(lab.input == nil)
                    if let output = lab.output { ShareLink("Export decrypted IPA", item: output) }
                    Text("Sign the exported IPA before installation. Record whether it installs and opens; compilation and structural checks alone do not prove correct plaintext.")
                }
                ForEach(["Components", "Transaction", "Inspection", "Enrollment", "Decryption"], id: \.self) { section in
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
            .fileImporter(isPresented: $importing, allowedContentTypes: importingCapture ? [.zip] : [.data]) { result in
                if case .success(let url) = result {
                    if importingCapture { lab.analyzeCapture(url) } else { lab.decrypt(url) }
                }
            }
            .onDisappear { lab.erase() }
        }
    }
}
