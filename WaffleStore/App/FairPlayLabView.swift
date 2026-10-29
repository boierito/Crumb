import SwiftUI
import Combine
import UniformTypeIdentifiers
import MapleSAP
import CryptoKit

nonisolated enum FairPlayNative {
    static func run(_ operation: String, source: URL? = nil, destination: URL? = nil,
                    report: URL? = nil, model: URL? = nil, ipa: URL? = nil) throws -> String {
        var cachePath = ""
        if operation == "probe-legacy-producer" {
            let cache = try FileManager.default.url(for: .cachesDirectory, in: .userDomainMask,
                appropriateFor: nil, create: true).appendingPathComponent("MapleSAP", isDirectory: true)
            try privateDirectory(cache)
            cachePath = cache.path
        }
        let request = ["Operation": operation, "Cache": cachePath, "Source": source?.path ?? "",
                       "Destination": destination?.path ?? "", "Report": report?.path ?? "", "Model": model?.path ?? "", "IPA": ipa?.path ?? ""]
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
    static func importEnrollment(_ selected: URL, to destination: URL) throws {
        let documentURL = selected.appendingPathComponent("manifest.json")
        let keys: Set<URLResourceKey> = [.isSymbolicLinkKey, .isRegularFileKey, .fileSizeKey]
        let documentInfo = try documentURL.resourceValues(forKeys: keys)
        guard documentInfo.isSymbolicLink != true, documentInfo.isRegularFile == true,
              let documentSize = documentInfo.fileSize, documentSize <= 1 << 20 else { throw CocoaError(.fileReadTooLarge) }
        let document = try Data(contentsOf: documentURL)
        guard document.count <= 1 << 20,
              let manifest = try JSONSerialization.jsonObject(with: document) as? [String: Any],
              manifest["schema"] as? String == "ipatool-fairplay-native-enrollment-capture-v1",
              let entries = manifest["files"] as? [[String: Any]], entries.count <= 256,
              let hardware = manifest["hardware_iomac"] as? [String: Any],
              let hardwarePath = hardware["path"] as? String else { throw CocoaError(.fileReadCorruptFile) }
        try privateDirectory(destination)
        try privateWrite(document, to: destination.appendingPathComponent("manifest.json"))
        var total = 0
        for entry in entries {
            guard let relative = entry["path"] as? String else { throw CocoaError(.fileReadCorruptFile) }
            guard relative == "app/executable" || relative == hardwarePath || relative.hasPrefix("app/SC_Info/") || relative.hasPrefix("state/") else { continue }
            let parts = relative.components(separatedBy: "/")
            guard !relative.contains("\\"), parts.count <= 32,
                  parts.allSatisfy({ !$0.isEmpty && $0 != "." && $0 != ".." }) else { throw CocoaError(.fileReadInvalidFileName) }
            var source = selected
            for part in parts {
                source.appendPathComponent(part)
                guard try source.resourceValues(forKeys: [.isSymbolicLinkKey]).isSymbolicLink != true else { throw CocoaError(.fileReadInvalidFileName) }
            }
            let info = try source.resourceValues(forKeys: keys)
            guard info.isRegularFile == true, let size = info.fileSize, size <= 64 << 20,
                  total + size <= 128 << 20 else { throw CocoaError(.fileReadTooLarge) }
            let data = try Data(contentsOf: source)
            guard data.count <= 64 << 20, total + data.count <= 128 << 20 else { throw CocoaError(.fileReadTooLarge) }
            total += data.count
            let target = destination.appendingPathComponent(relative)
            try privateDirectory(target.deletingLastPathComponent())
            try privateWrite(data, to: target)
        }
    }
    // Copy a bounded model registry into the protected session sandbox. The
    // native resolver still validates each pack and its selector association.
    static func importModel(_ selected: URL, to destination: URL) throws {
        let keys: Set<URLResourceKey> = [.isDirectoryKey, .isRegularFileKey, .isSymbolicLinkKey, .fileSizeKey]
        let info = try selected.resourceValues(forKeys: keys)
        guard info.isSymbolicLink != true else { throw CocoaError(.fileReadInvalidFileName) }
        func copyModel(_ source: URL, to target: URL) throws {
            let values = try source.resourceValues(forKeys: keys)
            guard values.isRegularFile == true, values.isSymbolicLink != true,
                  let size = values.fileSize, size > 0, size <= 1 << 20 else {
                throw CocoaError(.fileReadTooLarge)
            }
            let data = try Data(contentsOf: source)
            guard !data.isEmpty, data.count <= 1 << 20 else { throw CocoaError(.fileReadTooLarge) }
            try privateWrite(data, to: target)
        }
        guard info.isDirectory == true else {
            try copyModel(selected, to: destination)
            return
        }
        let entries = try FileManager.default.contentsOfDirectory(at: selected, includingPropertiesForKeys: Array(keys))
        guard entries.count <= 512 else { throw CocoaError(.fileReadTooLarge) }
        let models = entries.filter { $0.pathExtension == "bin" }
        guard !models.isEmpty, models.count <= 256 else { throw CocoaError(.fileReadTooLarge) }
        try privateDirectory(destination)
        for (index, source) in models.enumerated() {
            try copyModel(source, to: destination.appendingPathComponent("model-\(index).bin"))
        }
    }
    static func privateWrite(_ data: Data, to url: URL) throws {
        try data.write(to: url, options: [.atomic, .completeFileProtection])
        try FileManager.default.setAttributes([.posixPermissions: 0o600], ofItemAtPath: url.path)
    }
}

@MainActor
final class FairPlayLab: ObservableObject {
    @Published var busy = false
    @Published var logs = ["Components": "Pending", "Transaction": "Pending", "Inspection": "Pending", "Enrollment": "Pending", "Decryption": "Pending: run enrollment recovery or import a compatible model."]
    @Published var output: URL?
    @Published var diagnostic: URL?
    @Published var captureArchive: URL?
    @Published var importedCapture: URL?
    @Published var authorizedStateExport: URL?
    @Published var enrollmentStages: [String: Bool] = [:]
    @Published var nativeRecoveryVerified = false
    @Published var input: URL? {
        didSet {
            if oldValue != input {
                output = nil; nativeRecoveryVerified = false
                for stage in ["Portable raw state verified", "Enrollment frame verified", "Portable model validated"] { enrollmentStages[stage] = false }
            }
        }
    }
    let root = FileManager.default.temporaryDirectory.appendingPathComponent("FairPlayLab-" + UUID().uuidString)
    private var task: Task<Void, Never>?
    private var transactions = 0
    private var auditedTransactions = -1
    private var probeCompletedForCapture = false
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
        if let diversitybag = transaction.diversitybag { blobs.append(("diversitybag.bin", diversitybag)) }
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
        probeCompletedForCapture = false
        captureArchive = nil
        importedCapture = nil
        diagnostic = nil
        enrollmentStages = [:]
        log("Transaction", "Validated exchange \(transactions): kbsync \(transaction.kbsync.count) bytes; keybag \(transaction.keybag?.count ?? 0) bytes; SINF count \(transaction.sinfs.count). Capture saved privately.")
    }
    func download(_ text: String, account: StoreAccount, autoRecover: Bool = false, deepProbe: Bool = false) {
        perform {
            guard self.input == nil else { return }
            self.transactions = 0
            let retained = NativeRetainedKBSyncGenerator(accountStateDirectory: try NativeFairPlayAuthorization.stateDirectory(account: account))
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
            if autoRecover {
                // All operations happen inside the iOS application sandbox.
                // Missing ent/keybag data must not trigger speculative recovery.
                if self.transactions == 1 {
                    self.log("Enrollment", "Automatic on-device recovery requested for the retained transaction.")
                    do { try await self.attemptCurrentCaptureRecovery() }
                    catch {
                        self.log("Enrollment", "Local recovery unavailable (capture/validation). Original IPA retained; see the sanitized diagnostic.")
                    }
                } else {
                    self.log("Enrollment", "Automatic recovery skipped: exactly one retained transaction is required. Original IPA retained.")
                }
            }
            if deepProbe {
                if self.transactions == 1 {
                    do { try await self.performDeepLocalProbe() }
                    catch { self.log("Enrollment", "Multi-probe unavailable (capture/validation); original IPA retained.") }
                } else {
                    self.log("Enrollment", "Multi-probe skipped: exactly one retained transaction is required.")
                }
            }
        }
    }
    // A saved private capture and a saved original encrypted IPA are sufficient
    // to re-run all local audits. Importing either does not authenticate DRM.
    private func verifyImportedBinding(_ archive: URL, ipa: URL) async throws {
        let report = try await Task.detached {
            try FairPlayNative.run("verify-capture-ipa", source: archive, ipa: ipa)
        }.value
        guard let data = report.data(using: .utf8),
              let envelope = try JSONSerialization.jsonObject(with: data) as? [String: Any],
              envelope["ok"] as? Bool == true,
              let result = envelope["result"] as? [String: Any],
              result["schema"] as? String == "crumb-capture-ipa-binding-v1",
              result["matched"] as? Bool == true,
              result["encrypted_domain_count"] as? Int == 1,
              result["model_validated"] as? Bool == false else {
            throw CocoaError(.fileReadCorruptFile)
        }
    }

    func useExistingOriginalIPA(_ selected: URL) {
        guard input == nil else { return }
        perform {
            let access = selected.startAccessingSecurityScopedResource()
            defer { if access { selected.stopAccessingSecurityScopedResource() } }
            let info = try selected.resourceValues(forKeys: [.isSymbolicLinkKey, .isRegularFileKey, .fileSizeKey])
            guard selected.pathExtension.lowercased() == "ipa",
                  info.isSymbolicLink != true, info.isRegularFile == true,
                  let size = info.fileSize, size > 0, size <= 512 << 20 else {
                throw CocoaError(.fileReadCorruptFile)
            }
            let destination = self.root.appendingPathComponent("Original.ipa")
            try? FileManager.default.removeItem(at: destination)
            try FileManager.default.copyItem(at: selected, to: destination)
            do {
                try FileManager.default.setAttributes([.posixPermissions: 0o600,
                    .protectionKey: FileProtectionType.complete], ofItemAtPath: destination.path)
                let result = try await Task.detached {
                    try FairPlayNative.run("inspect", source: destination)
                }.value
                guard let bytes = result.data(using: .utf8),
                      let envelope = try JSONSerialization.jsonObject(with: bytes) as? [String: Any],
                      envelope["ok"] as? Bool == true else {
                    throw CocoaError(.fileReadCorruptFile)
                }
                if let archive = self.importedCapture {
                    try await self.verifyImportedBinding(archive, ipa: destination)
                }
                self.input = destination
                self.log("Transaction", "Existing encrypted IPA imported and inspected offline. No App Store download.")
            } catch {
                try? FileManager.default.removeItem(at: destination)
                throw error
            }
        }
    }

    func analyzeCapture(_ selected: URL) {
        perform {
            let access = selected.startAccessingSecurityScopedResource()
            defer { if access { selected.stopAccessingSecurityScopedResource() } }
            self.log("Inspection", "Validating imported capture ZIP and manifest hashes locally.")
            let report = try await Task.detached { try FairPlayNative.run("analyze-capture", source: selected) }.value
            guard let bytes = report.data(using: .utf8),
                  let response = try JSONSerialization.jsonObject(with: bytes) as? [String: Any],
                  response["ok"] as? Bool == true,
                  let result = response["result"] as? [String: Any],
                  result["transaction_count"] as? Int == 1,
                  let transactions = result["transactions"] as? [[String: Any]],
                  transactions.count == 1,
                  transactions[0]["integrity"] as? String == "verified" else {
                throw CocoaError(.fileReadCorruptFile)
            }
            let info = try selected.resourceValues(forKeys: [.isSymbolicLinkKey, .isRegularFileKey, .fileSizeKey])
            guard info.isSymbolicLink != true, info.isRegularFile == true,
                  let size = info.fileSize, size > 0, size <= 128 << 20 else {
                throw CocoaError(.fileReadTooLarge)
            }
            let copied = self.root.appendingPathComponent("Imported-FairPlay-capture.zip")
            try? FileManager.default.removeItem(at: copied)
            try FileManager.default.copyItem(at: selected, to: copied)
            try FileManager.default.setAttributes([.posixPermissions: 0o600,
                .protectionKey: FileProtectionType.complete], ofItemAtPath: copied.path)
            do {
                if let original = self.input {
                    try await self.verifyImportedBinding(copied, ipa: original)
                }
            } catch {
                try? FileManager.default.removeItem(at: copied)
                throw error
            }
            self.importedCapture = copied
            self.captureArchive = nil
            self.auditedTransactions = -1
            self.transactions = 1
            self.log("Inspection", report)
            self.log("Enrollment", "Verified private capture imported. Select its matching original encrypted IPA from Files; no download is required.")
        }
    }
    func inspect() {
        guard let input else { return }
        perform {
            let result = try await Task.detached { try FairPlayNative.run("inspect", source: input) }.value
            self.log("Inspection", result)
            self.log("Inspection", "Inspection reports encryption geometry only. Attempt recovery from current capture tests enrollment compatibility; inspection does not decrypt pages.")
        }
    }
    // Used by both the one-tap iPhone flow and manual troubleshooting.
    // The native bridge reports success only after source/frame and SINF/IPA
    // validation. The additional Swift gate rejects partial/malformed replies.
    private func attemptCurrentCaptureRecovery() async throws {
        guard let input else { throw CocoaError(.fileReadNoSuchFile) }
        output = nil
        nativeRecoveryVerified = false
        for stage in ["Portable raw state verified", "Enrollment frame verified", "Portable model validated"] {
            enrollmentStages[stage] = false
        }
        try await analyzeCurrentCaptureFiles()
        guard transactions == 1, let archive = captureArchive else {
            log("Enrollment", "Recovery requires exactly one verified retained iPhone transaction.")
            return
        }
        let destination = root.appendingPathComponent(UUID().uuidString + ".recovery")
        let processed = destination.appendingPathComponent("processed.ipa")
        var validated = false
        defer {
            // Never persist enrollment models or reports outside this attempt.
            try? FileManager.default.removeItem(at: destination.appendingPathComponent("models"))
            try? FileManager.default.removeItem(at: destination.appendingPathComponent("private-report.json"))
            if !validated { try? FileManager.default.removeItem(at: destination) }
        }
        log("Enrollment", "Attempting local recovery with the complete retained SIDB and its original identity. Compatibility remains experimental.")
        let response = try await Task.detached {
            try FairPlayNative.run("recover-capture", source: archive, destination: destination, ipa: input)
        }.value
        log("Enrollment", response)
        let envelope = (response.data(using: .utf8)
            .flatMap { try? JSONSerialization.jsonObject(with: $0) }) as? [String: Any] ?? [:]
        let size = ((try? FileManager.default.attributesOfItem(atPath: processed.path))?[.size] as? NSNumber)?.int64Value ?? 0
        guard FairPlayLocalRecoveryGate.accepts(envelope, outputBytes: size) else {
            log("Decryption", "Local recovery rejected or incomplete. Original IPA retained; no processed IPA published.")
            return
        }
        // These three indicators describe recovered model consistency only,
        // never independently authenticated plaintext or installation.
        enrollmentStages["Portable raw state verified"] = true
        enrollmentStages["Enrollment frame verified"] = true
        enrollmentStages["Portable model validated"] = true
        nativeRecoveryVerified = true
        output = processed
        validated = true
        log("Decryption", "Source/frame and model/SINF checks passed locally. IPA is available for analysis; plaintext correctness and installation remain unverified.")
    }
    func recoverCurrentCapture() {
        guard input != nil else { return }
        perform { try await self.attemptCurrentCaptureRecovery() }
    }
    // Compare multiple virtual file layouts from one verified, unchanged
    // private StoreAgent capture. The native result is fixed-label metadata.
    private func performDeepLocalProbe() async throws {
        guard let input else { throw CocoaError(.fileReadNoSuchFile) }
        try await analyzeCurrentCaptureFiles()
        guard transactions == 1, let archive = captureArchive else {
            log("Enrollment", "Multi-probe requires exactly one retained local transaction.")
            return
        }
        guard !probeCompletedForCapture else {
            log("Enrollment", "This unchanged capture was already tested in all four modes. Export its diagnostic; repeated probes cannot change its inputs.")
            return
        }
        log("Enrollment", "Running four bounded on-device recovery probes with unchanged SIDB; no private values or models exported.")
        let response = try await Task.detached {
            try FairPlayNative.run("probe-capture", source: archive, ipa: input)
        }.value
        let envelope = (response.data(using: .utf8)
            .flatMap { try? JSONSerialization.jsonObject(with: $0) }) as? [String: Any] ?? [:]
        // Only cache a completed four-way response, never a native failure or
        // truncated run. On a new capture the cache is explicitly invalidated.
        if envelope["ok"] as? Bool == true,
           let values = envelope["result"] as? [String: Any],
           let attempts = values["attempts"] as? Int, attempts == 4,
           values["truncated_at_time_budget"] as? Bool == false {
            probeCompletedForCapture = true
        }
        log("Enrollment", response)
        log("Enrollment", "All probe results are diagnostic only; they do not prove decryption.")
    }
    func deepProbeCurrentCapture() {
        guard input != nil else { return }
        perform { try await self.performDeepLocalProbe() }
    }
    // A single local differential experiment on a retained, verified capture.
    // Does not download another IPA, extract credentials or export a model.
    func researchCurrentCaptureCausality() {
        guard input != nil else { return }
        perform {
            try await self.analyzeCurrentCaptureFiles()
            guard self.transactions == 1, let archive = self.captureArchive, let input = self.input else {
                self.log("Enrollment", "Causal research requires one retained transaction and its original IPA; no download attempted.")
                return
            }
            self.log("Enrollment", "Studying SIDB and virtual identity causal dependencies offline. Original capture unchanged; private bytes remain on this iPhone.")
            let response = try await Task.detached {
                try FairPlayNative.run("research-causality", source: archive, ipa: input)
            }.value
            self.log("Enrollment", response)
            self.log("Enrollment", "This research measures sensitivity only. It does not recover a valid model or prove decrypted plaintext.")
        }
    }
    func authorizeAccountAndRecover(_ account: StoreAccount) {
        guard let original = input else { return }
        perform {
            self.output = nil
            self.nativeRecoveryVerified = false
            if let previous = self.authorizedStateExport { try? FileManager.default.removeItem(at: previous) }
            self.authorizedStateExport = nil
            if !NativeFairPlayAuthorization.hasSavedState(account: account) {
                try await self.analyzeCurrentCaptureFiles()
                guard self.transactions == 1, self.captureArchive != nil else {
                    self.log("Enrollment", "The first account authorization needs one saved capture and its original IPA.")
                    return
                }
            }
            let destination = self.root.appendingPathComponent("authorized-" + UUID().uuidString)
            self.log("Enrollment", "Preparing account authorization on this iPhone. A required authorization follows validated Apple routing without automatic retries; saved account state is reused across apps. Keep Crumb open.")
            do {
                let result = try await NativeFairPlayAuthorization.run(account: account, source: self.captureArchive, ipa: original, destination: destination,
                    diagnostic: { message in self.log("Enrollment", message) })
                let object = try JSONSerialization.jsonObject(with: result)
                let pretty = try JSONSerialization.data(withJSONObject: object, options: [.prettyPrinted, .sortedKeys])
                self.log("Enrollment", String(decoding: pretty, as: UTF8.self))
                let envelope = object as? [String: Any] ?? [:]
                let report = envelope["result"] as? [String: Any] ?? [:]
                if report["account_state_saved"] as? Bool == true {
                    await self.prepareAuthorizedStateExport(account)
                }
                let processed = destination.appendingPathComponent("processed.ipa")
                let size = ((try? FileManager.default.attributesOfItem(atPath: processed.path))?[.size] as? NSNumber)?.int64Value ?? 0
                if report["completed"] as? Bool == true,
                   report["sinf_model_verified"] as? Bool == true,
                   report["source_frame_verified"] as? Bool == true,
                   report["plaintext_authenticated"] as? Bool == false, size > 0 {
                    self.output = processed
                    self.nativeRecoveryVerified = true
                    self.log("Decryption", "All license domains and IPA processing checks passed. Processed output is available; plaintext authenticity remains unverified.")
                } else {
                    try? FileManager.default.removeItem(at: destination)
                    self.log("Decryption", "Account state and per-app license validation are separate. No completed IPA was produced or exported.")
                }
            } catch {
                try? FileManager.default.removeItem(at: destination)
                let code: String
                if let failure = error as? FairPlayAuthorizationError { code = failure.rawValue }
                else if let failure = error as? NativeFairPlayAuthorization.Failure,
                        case .native(let native) = failure { code = native }
                else { code = "authorization_failed" }
                self.log("Enrollment", "Account authorization stopped: " + code + ". No automatic replay; original IPA and capture retained.")
            }
        }
    }

    private func prepareAuthorizedStateExport(_ account: StoreAccount) async {
        if let previous = authorizedStateExport { try? FileManager.default.removeItem(at: previous) }
        authorizedStateExport = nil
        let formatter = DateFormatter()
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.timeZone = TimeZone(secondsFromGMT: 0)
        formatter.dateFormat = "yyyy-MM-dd-HHmmss"
        let file = root.appendingPathComponent("PRIVATE-FairPlay-authorized-state-" + formatter.string(from: Date()) + "-UTC.json")
        do {
            try await NativeFairPlayAuthorization.exportState(account: account, to: file)
            authorizedStateExport = file
            log("Enrollment", "Private authorized-state copy ready to share. It contains licensing material; no password or Apple session cookies. SINF validation is a separate check.")
        } catch {
            try? FileManager.default.removeItem(at: file)
            if let failure = error as? NativeFairPlayAuthorization.Failure,
               case .native(let code) = failure {
                switch code {
                case "authorization_state_missing":
                    log("Enrollment", "No saved authorization exists for the current account and device identity. Download the app, then use Authorize account and validate this IPA before exporting authorization. Export did not modify internal state.")
                case "authorization_cached_state_invalid":
                    log("Enrollment", "Saved authorization could not be validated for the current account and device identity. Internal files were not changed. Preserve this installation and copy these logs.")
                default:
                    log("Enrollment", "Authorized-state export stopped: " + code + ". Internal files were not changed.")
                }
            } else {
                log("Enrollment", "Authorized-state export failed during local file preparation or protection. Internal files were not changed.")
            }
        }
    }

    func exportSavedAuthorization(_ account: StoreAccount) {
        perform { await self.prepareAuthorizedStateExport(account) }
    }

    // Reopen an isolated state copy in the original legacy producer.
    // This does not perform SINF validation.
    func probeOriginalProducer() {
        guard input != nil else { return }
        perform {
            try await self.analyzeCurrentCaptureFiles()
            guard self.transactions == 1, let archive = self.captureArchive, let original = self.input else {
                self.log("Enrollment", "Original producer test requires the saved capture and its original encrypted IPA.")
                return
            }
            self.log("Enrollment", "Restoring a private copy in the original StoreAgent 10.9. This can take several minutes; no new app download or keybag request.")
            let response = try await Task.detached {
                try FairPlayNative.run("probe-legacy-producer", source: archive, ipa: original)
            }.value
            self.log("Enrollment", response)
            self.log("Enrollment", "Provisioning query only: neither outcome validates the SINF or produces a decrypted IPA.")
        }
    }

    func researchCurrentCaptureBridge() {
        guard input != nil else { return }
        perform {
            try await self.analyzeCurrentCaptureFiles()
            guard self.transactions == 1, let archive = self.captureArchive, let original = self.input else {
                self.log("Enrollment", "Bridge audit requires one retained transaction and its original IPA.")
                return
            }
            self.log("Enrollment", "Auditing StoreAgent 10.9 state against the modern fairplayd consumer. No downloads or mutations.")
            let response = try await Task.detached {
                try FairPlayNative.run("research-bridge", source: archive, ipa: original)
            }.value
            self.log("Enrollment", response)
            self.log("Enrollment", "Only strict SINF validation can establish a usable model. This report exports no keys.")
        }
    }
    func recoverEnrollment(_ selected: URL) {
        nativeRecoveryVerified = false
        guard let input else { return }
        perform {
            self.output = nil
            let source = self.root.appendingPathComponent(UUID().uuidString + ".enrollment")
            let destination = self.root.appendingPathComponent(UUID().uuidString + ".recovery")
            let access = selected.startAccessingSecurityScopedResource()
            defer {
                if access { selected.stopAccessingSecurityScopedResource() }
                try? FileManager.default.removeItem(at: source)
                // Only the processed IPA may remain for export.
                try? FileManager.default.removeItem(at: destination.appendingPathComponent("models"))
                try? FileManager.default.removeItem(at: destination.appendingPathComponent("private-report.json"))
            }
            self.log("Enrollment", "Importing modern native enrollment. Legacy StoreAgent captures cannot supply these inputs.")
            let result = try await Task.detached {
                try FairPlayNative.importEnrollment(selected, to: source)
                return try FairPlayNative.run("recover-enrollment", source: source, destination: destination, ipa: input)
            }.value
            self.log("Enrollment", result)
            let envelope = (result.data(using: .utf8)
                .flatMap { try? JSONSerialization.jsonObject(with: $0) }) as? [String: Any] ?? [:]
            let processed = destination.appendingPathComponent("processed.ipa")
            let size = ((try? FileManager.default.attributesOfItem(atPath: processed.path))?[.size] as? NSNumber)?.int64Value ?? 0
            guard FairPlayLocalRecoveryGate.acceptsNativeEnrollment(envelope, outputBytes: size) else {
                try? FileManager.default.removeItem(at: destination)
                self.log("Decryption", "Native enrollment recovery rejected or produced no validated IPA. Nothing exported.")
                return
            }
            self.nativeRecoveryVerified = true
            self.output = processed
            self.log("Decryption", "Source/frame, model/SINF and processed IPA checks passed. Plaintext authenticity remains unverified.")
        }
    }
    func decrypt(_ selected: URL) {
        guard let input else { return }
        perform {
            self.output = nil
            self.nativeRecoveryVerified = false
            let model = self.root.appendingPathComponent(UUID().uuidString + ".models")
            let access = selected.startAccessingSecurityScopedResource()
            defer { if access { selected.stopAccessingSecurityScopedResource() }; try? FileManager.default.removeItem(at: model) }
            try FairPlayNative.importModel(selected, to: model)
            let output = self.root.appendingPathComponent(UUID().uuidString + ".decrypted.ipa")
            let report = self.root.appendingPathComponent("private-report.json")
            defer { try? FileManager.default.removeItem(at: report) }
            self.log("Decryption", "Checking model/SINF and processing IPA. Native failures identify the operation; private exception text is suppressed.")
            let result = try await Task.detached { try FairPlayNative.run("decrypt", source: input, destination: output, report: report, model: model) }.value
            self.log("Decryption", result)
            let envelope = (result.data(using: .utf8)
                .flatMap { try? JSONSerialization.jsonObject(with: $0) }) as? [String: Any] ?? [:]
            let size = ((try? FileManager.default.attributesOfItem(atPath: output.path))?[.size] as? NSNumber)?.int64Value ?? 0
            guard FairPlayLocalRecoveryGate.acceptsModelDecrypt(envelope, outputBytes: size) else {
                try? FileManager.default.removeItem(at: output)
                self.log("Decryption", "Model-based processing failed validation or produced no IPA. Nothing exported.")
                return
            }
            self.output = output
            self.log("Decryption", "Processed IPA is available for inspection; plaintext authenticity remains unverified.")
        }
    }
    // Native verification works directly on the capture created in this session.
    // The user no longer has to import their own private ZIP to analyze it.
    private func analyzeCurrentCaptureFiles() async throws {
        let captures = root.appendingPathComponent("captures")
        guard importedCapture != nil || FileManager.default.fileExists(atPath: captures.path) else {
            log("Enrollment", "No new capture in this session. Download an app before analyzing.")
            return
        }
        let archive = importedCapture ?? root.appendingPathComponent("Private-FairPlay-capture.zip")
        // Once checked, the same immutable transaction does not need to be
        // zipped, inspected and logged again by Prepare exports.
        if auditedTransactions == transactions,
           captureArchive != nil,
           FileManager.default.fileExists(atPath: archive.path) {
            return
        }
        if importedCapture == nil {
            let exported = try await Task.detached {
                try FairPlayNative.run("export", source: captures, destination: archive)
            }.value
            guard let exportData = exported.data(using: .utf8),
                  let exportJSON = try JSONSerialization.jsonObject(with: exportData) as? [String: Any],
                  exportJSON["ok"] as? Bool == true else {
                throw NSError(domain: "Crumb.FairPlayLab.CaptureExport", code: 1)
            }
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
        nativeRecoveryVerified = false
        guard !busy else { return }
        try? FileManager.default.removeItem(at: root)
        input = nil; output = nil; diagnostic = nil; captureArchive = nil; importedCapture = nil; authorizedStateExport = nil; transactions = 0
        auditedTransactions = -1
        probeCompletedForCapture = false
        enrollmentStages = [:]
        logs = ["Components": "Pending", "Transaction": "Pending", "Inspection": "Pending", "Enrollment": "Pending", "Decryption": "Pending: run enrollment recovery or import a compatible model."]
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
    @State private var importingOriginalIPA = false
    @State private var importingEnrollment = false
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
                    Button("Download + attempt local recovery on iPhone") {
                        if let account = appData.ipaTool?.account {
                            lab.download(appInput, account: account, autoRecover: true)
                        }
                    }.disabled(appData.ipaTool == nil || appInput.isEmpty || lab.input != nil)
                    Text("One-tap experimental path: download, retain the matching keybag context, audit the capture, and attempt recovery entirely on this iPhone. A failed recovery leaves the original IPA available.")
                        .font(.footnote)
                    if let record, let url = record.fileURL {
                        Button("Use selected downloaded IPA") { lab.input = url }
                    }
                    Button("Import existing encrypted IPA from Files (no download)") {
                        importingOriginalIPA = true
                    }.disabled(lab.input != nil)
                }
                Section("3 · Inspect and identify missing material") {
                    Button("Inspect IPA") { lab.inspect() }.disabled(lab.input == nil)
                    Button("Verify current capture — no import required") { lab.analyzeCurrentCapture() }
                        .disabled(lab.input == nil)
                    Button("Import an older capture ZIP for offline research") { importingEnrollment = false; importingCapture = true; importing = true }
                    Text("Pair the ZIP with its original encrypted IPA from Files. Both are verified locally and available for all recovery/bridge tests during this session.")
                        .font(.footnote)
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
                Section("4 · Recover or import an enrollment") {
                    Button("Attempt recovery from current capture") { lab.recoverCurrentCapture() }.disabled(lab.input == nil)
                    Button("Authorize account and validate this IPA") {
                        if let account = appData.ipaTool?.account { lab.authorizeAccountAndRecover(account) }
                    }.disabled(lab.input == nil || appData.ipaTool == nil)
                    Text("May authorize this virtual device with Apple. Account state is saved locally for reuse with other apps; every app's licenses are validated separately. An ADI challenge stops the attempt.")
                        .font(.footnote)
                    Button("Test original StoreAgent state (no download)") { lab.probeOriginalProducer() }.disabled(lab.input == nil)
                    Text("Reopens a private copy with the original producer and checks account provisioning. It does not decrypt the SINF or export an IPA.")
                        .font(.footnote)
                    Button("Audit StoreAgent → fairplayd bridge (no download)") { lab.researchCurrentCaptureBridge() }.disabled(lab.input == nil)
                    Text("Verifies SIDB/SIDD provenance, records only static import slots and per-phase state-file lookups, and checks the original SINF. No private material leaves this iPhone.")
                        .font(.footnote)
                    Button("Research SIDB causality (no download)") { lab.researchCurrentCaptureCausality() }.disabled(lab.input == nil)
                    Text("Uses the already captured SIDB and IPA. Runs bounded one-byte changes only on private copies and compares their effect on Source, Frame and model, without exporting secrets.")
                        .font(.footnote)
                    Button("Run four-way diagnostic (reuse current capture)") { lab.deepProbeCurrentCapture() }.disabled(lab.input == nil)
                    Text("Compares dual state paths, fpsd-only, shared-only, and SIDB-only without changing the encrypted state bytes. Results appear in Enrollment logs.")
                        .font(.footnote)
                    Text("Experimental compatibility attempt using the complete retained SIDB and its original identity. Rejection leaves enrollment unverified.")
                    Button("Recover model from native enrollment folder") { importingEnrollment = true; importingCapture = false; importing = true }.disabled(lab.input == nil)
                    if lab.nativeRecoveryVerified { Label("Derived source and frame cross-checked; model matches target SINF", systemImage: "checkmark.circle.fill").foregroundStyle(.green) }
                    Text("Requires a modern native enrollment capture from your own App Store account, including its hardware identity and state files. The current StoreAgent capture is a separate format.")
                    Button("Import model or model folder and decrypt") { importingEnrollment = false; importingCapture = false; importing = true }.disabled(lab.input == nil)
                    Text("Import one model file or a folder of selector-bound .bin models. Each executable uses its associated license. Output plaintext is not authenticated.")
                    if let output = lab.output { ShareLink("Export processed IPA (unverified)", item: output) }
                    Text("Sign the exported IPA before installation. Record whether it installs and opens; compilation and structural checks alone do not prove correct plaintext.")
                }
                ForEach(["Components", "Transaction", "Inspection", "Enrollment", "Decryption"], id: \.self) { section in
                    Section(section + " logs") {
                        ScrollView { Text(lab.logs[section] ?? "Pending").font(.system(.caption, design: .monospaced)).frame(maxWidth: .infinity, alignment: .leading).textSelection(.enabled) }.frame(maxHeight: 220)
                    }
                }
                Section("5 · Export evidence and repeat") {
                    Button("Prepare PRIVATE authorized state (no new authorization)") {
                        if let account = appData.ipaTool?.account { lab.exportSavedAuthorization(account) }
                    }.disabled(appData.ipaTool == nil)
                    if let state = lab.authorizedStateExport {
                        ShareLink("Share PRIVATE authorized state — contains licensing secrets", item: state)
                    }
                    Button("Prepare exports") { lab.export() }
                    if let diagnostic = lab.diagnostic { ShareLink("Share diagnostic (no secret blobs)", item: diagnostic) }
                    if let capture = lab.captureArchive { ShareLink("Share PRIVATE capture — contains secrets", item: capture) }
                    if let original = lab.input, original.pathExtension.lowercased() == "ipa" {
                        ShareLink("Save original encrypted IPA for offline reuse", item: original)
                    }
                    Text("Keep both files together. The capture ZIP contains account-derived material. Avoid public sharing; both files are needed for offline replay.")
                        .font(.footnote)
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
            .fileImporter(isPresented: $importingOriginalIPA, allowedContentTypes: [UTType(filenameExtension: "ipa") ?? .data]) { result in
                if case .success(let url) = result { lab.useExistingOriginalIPA(url) }
            }
            .fileImporter(isPresented: $importing, allowedContentTypes: importingEnrollment ? [.folder] : (importingCapture ? [.zip] : [.data, .folder])) { result in
                if case .success(let url) = result {
                    if importingEnrollment { lab.recoverEnrollment(url) } else if importingCapture { lab.analyzeCapture(url) } else { lab.decrypt(url) }
                }
            }
            .onDisappear { lab.erase() }
        }
    }
}
