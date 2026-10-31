import SwiftUI
import Combine
import UniformTypeIdentifiers
import MapleSAP

nonisolated private struct SystemEnvelope<T: Decodable>: Decodable {
    let ok: Bool
    let result: T?
}
nonisolated private struct SystemLocalPrepared: Decodable { let schema: String; let requiresSigning: Bool }
nonisolated private struct SystemSigningFolderResponse: Decodable { let bundlePath, bundleId, version, build: String }
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

nonisolated private enum SystemNativeProcessor {
    static func process(source: URL, directory: URL, destination: URL, identifiers: [String], route: String) throws -> String {
        let nativeDirectory = directory.appendingPathComponent(UUID().uuidString, isDirectory: true)
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
                    throw NSError(domain: "Crumb.SystemFairPlay", code: Int(result.remap_errno != 0 ? result.remap_errno : result.read_error), userInfo: [
                        "native_remap_return": result.remap_result, "native_remap_errno": result.remap_errno,
                        "native_mapping_attempted": result.mapping_attempted != 0,
                        "native_mapping_errno": result.mapping_errno,
                        "native_remap_attempted": result.remap_attempted != 0,
                        "native_read_error": result.read_error, "pages_total": result.pages_total,
                        "pages_read": result.pages_read, "pages_changed": result.pages_changed])
                }
            }
        }
        let export = try FairPlayNative.run("export-system-result", source: source, destination: destination,
            model: nativeDirectory, privateIdentifiers: identifiers)
        guard var outputEnvelope = try JSONSerialization.jsonObject(with: Data(export.utf8)) as? [String: Any],
              outputEnvelope["ok"] as? Bool == true else { throw CocoaError(.fileReadCorruptFile) }
        outputEnvelope["native_route"] = route
        return String(decoding: try JSONSerialization.data(withJSONObject: outputEnvelope, options: [.sortedKeys, .prettyPrinted]), as: UTF8.self)
    }
    static func safeFailure(_ error: Error, route: String) -> [String: Any] {
        let native = error as NSError
        var result: [String: Any] = ["route": route, "code": native.code,
            "stage": native.domain == "Crumb.SystemFairPlay" ? SystemNativeFailureStage.classify(native.userInfo) : "preparation_or_validation"]
        for key in ["native_remap_return", "native_remap_errno", "native_read_error", "native_mapping_attempted", "native_mapping_errno", "native_remap_attempted", "pages_total", "pages_read", "pages_changed"] {
            if let count = native.userInfo[key] as? NSNumber { result[key] = count }
        }
        return result
    }
}

@MainActor
final class SystemFairPlaySession: ObservableObject {
    static let shared = SystemFairPlaySession()
    @Published var showIdentity = false
    @Published var failureMessage = ""
    @Published var requestedRecordID = ""
    @Published var signingPassword = ""
    @Published var identityStatus = "Not checked"
    @Published var identityValidated = false
    @Published var installAfterDecryption = false
    @Published var signedOutput: URL?
    @Published var readyToSignOutput = false
    @Published var outputBundleID = ""
    @Published var outputVersion = ""
    @Published var hasOutputProfile = false
    @Published var pendingJobID = ""
    @Published var phase = "idle"
    @Published var importName = ""
    private var automaticSigning = false
    private var reportFile: URL?
    var canResume: Bool { !pendingJobID.isEmpty && pendingJobID == (installedIdentity?["job_id"] as? String) && phase == "awaiting_installation" }
    var currentBundleID: String { Bundle.main.bundleIdentifier ?? "" }
    private func currentTeam() throws -> String {
        var buffer = [CChar](repeating: 0, count: 256)
        guard CrumbCurrentSigningTeam(&buffer, 256) == 0 else { throw NSError(domain: "Crumb.LocalSigning", code: 5) }
        return String(cString: buffer)
    }
    static func signingError(_ code: Int) -> String {
        switch code {
        case 5: return "Could not read Crumb’s signing team. Reinstall this update using the same signing identity as your current Crumb."
        case 10: return "The P12 could not be read. Import it again."
        case 11: return "The P12 could not be opened. Check its password and that this is a .p12 file, not a certificate alone."
        case 12: return "The P12 must contain a matching certificate and private key. Export both from your signer."
        case 13: return "The certificate has expired or is not valid yet. Check the date and renew the certificate."
        case 14: return "The provisioning profile is unreadable or its signature is invalid. Import the original .mobileprovision file."
        case 15: return "The provisioning profile has expired. Import a renewed profile."
        case 16: return "This certificate is not listed in the provisioning profile. Import a matching pair."
        case 17: return "This profile belongs to a different signing team. Updating Crumb requires its current team."
        case 18: return "The profile does not allow this bundle ID. Use a matching profile or an appropriate wildcard profile."
        case 19: return "The profile does not cover every app extension. Use a profile with a compatible wildcard; Crumb will not silently remove extensions."
        default: return "Signing could not finish (code \(code)). The original IPA is retained."
        }
    }
    // The view exposes one action; the coordinator chooses the next real step.
    func continueDecryption(record: DownloadRecord?, identifiers: [String]) -> Bool {
        guard !busy else { return false }
        failureMessage = ""
        installAfterDecryption = false
        if let record, requestedRecordID != record.id {
            requestedRecordID = record.id
            helper = nil; signedHelper = nil; output = nil; signedOutput = nil
            pendingJobID = ""; report = ""; reportFile = nil; phase = "idle"
            try? save()
        }
        if record != nil && output != nil && phase == "exported" { return false }
        let boundInstalledJob = !pendingJobID.isEmpty && pendingJobID == (installedIdentity?["job_id"] as? String)
        if boundInstalledJob && matchesInstalledTarget(record) && installedTarget != nil {
            if output == nil { decrypt(identifiers: identifiers) }
            return false
        }
        guard let original = record?.fileURL else {
            failureMessage = "Select an app in Downloads to decrypt it."
            return false
        }
        // No signing, reinstallation or automatic helper fallback is allowed on
        // this route. A system refusal remains an explicit failed attempt.
        helper = nil; signedHelper = nil
        decrypt(identifiers: identifiers, localIPA: original, description: record?.appName)
        return false
    }
    func validateIdentity(completion: (() -> Void)? = nil) {
        guard !busy, hasCertificate, hasProfile else { return }
        busy = true
        let password = signingPassword
        Task {
            var validated = false
            defer { busy = false; if validated { completion?() } }
            do {
                let directory = try signingStorage(), team = try currentTeam(), bundle = currentBundleID
                let code = await Task.detached {
                    CrumbValidateSigningIdentity(directory.appendingPathComponent("Identity.p12").path,
                        directory.appendingPathComponent("Profile.mobileprovision").path, password, bundle, team)
                }.value
                identityValidated = code == 0; validated = code == 0
                identityStatus = code == 0 ? "Certificate and profile match Crumb. iOS still checks trust, revocation and device eligibility during installation." : Self.signingError(Int(code))
            } catch { identityValidated = false; identityStatus = Self.signingError((error as NSError).code) }
        }
    }
    func resume(identifiers: [String]) {
        guard !busy, canResume else { return }
        decrypt(identifiers: identifiers)
    }
    func importInbox() {
        guard !busy else { return }
        do {
            let docs = try FileManager.default.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: true)
            let inbox = docs.appendingPathComponent("Signing Inbox", isDirectory: true)
            try FileManager.default.createDirectory(at: inbox, withIntermediateDirectories: true)
            let files = try FileManager.default.contentsOfDirectory(at: inbox, includingPropertiesForKeys: [.isRegularFileKey])
            let wanted = importName == "Identity.p12" ? "p12" : "mobileprovision"
            let matching = files.filter { $0.pathExtension.lowercased() == wanted }
            guard matching.count == 1, let file = matching.first else {
                message = "In Files → On My iPhone → Crumb → Signing Inbox, place exactly one .\(wanted) file, then tap Import from Signing Inbox again."
                return
            }
            importSigningFile(file, certificate: wanted == "p12", outputProfile: importName == "App.mobileprovision", removeInboxCopy: true)
        } catch { message = "Could not read Signing Inbox. Use the Files picker instead." }
    }
    func removeIdentity() {
        guard !busy else { return }
        do {
            try FileManager.default.removeItem(at: signingStorage())
            hasCertificate = false; hasProfile = false; hasOutputProfile = false
            signingPassword = ""; identityValidated = false; identityStatus = "Not configured"
        } catch { message = "Could not remove the identity. Try again while the iPhone is unlocked." }
    }

    @Published var busy = false
    @Published var helper: URL?
    @Published var output: URL?
    @Published var signedHelper: URL?
    @Published var hasCertificate = false
    @Published var hasProfile = false
    @Published var preparedDescription = ""
    @Published var exportDescription = ""
    @Published var report = ""
    @Published var message = ""
    private var installedIdentity: [String: Any]? {
        guard let data = try? Data(contentsOf: Bundle.main.bundleURL.appendingPathComponent("CrumbSystemPayload/job.json")),
              data.count <= 1 << 20 else { return nil }
        return try? JSONSerialization.jsonObject(with: data) as? [String: Any]
    }
    func matchesInstalledTarget(_ record: DownloadRecord?) -> Bool {
        guard let record else { return true }
        guard let job = installedIdentity, job["bundle_id"] as? String == record.bundleID,
              job["version"] as? String == record.version else { return false }
        return record.build == nil || job["build"] as? String == record.build
    }
    var installedDescription: String {
        guard let job = installedIdentity else { return "" }
        return [job["bundle_id"],job["version"],job["build"]].compactMap { $0 as? String }.joined(separator: " · ")
    }
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
    private func signingStorage() throws -> URL {
        let directory = try storage().appendingPathComponent("Signing", isDirectory: true)
        try FairPlayNative.privateDirectory(directory)
        return directory
    }
    func importSigningFile(_ selected: URL, certificate: Bool, outputProfile: Bool = false, removeInboxCopy: Bool = false) {
        guard !busy else { return }; busy = true
        Task {
            defer { busy = false }
            do {
                let directory = try signingStorage()
                let destination = directory.appendingPathComponent(certificate ? "Identity.p12" : (outputProfile ? "App.mobileprovision" : "Profile.mobileprovision"))
                try await Task.detached {
                    let accessing = selected.startAccessingSecurityScopedResource()
                    defer { if accessing { selected.stopAccessingSecurityScopedResource() } }
                    let values = try selected.resourceValues(forKeys: [.fileSizeKey, .isRegularFileKey, .isSymbolicLinkKey])
                    guard values.isRegularFile == true, values.isSymbolicLink != true,
                          let size = values.fileSize, size > 0, size <= 8 << 20 else { throw CocoaError(.fileReadTooLarge) }
                    try FairPlayNative.privateWrite(Data(contentsOf: selected), to: destination)
                }.value
                let temporary = FileManager.default.temporaryDirectory.standardizedFileURL.path + "/"
                let choco = NSHomeDirectory() + "/Documents/ChocoFilePickerFixFolder/"
                if removeInboxCopy || selected.standardizedFileURL.path.hasPrefix(temporary) || selected.standardizedFileURL.path.hasPrefix(choco) {
                    try? FileManager.default.removeItem(at: selected)
                }
                identityValidated = false; identityStatus = "Check this identity before continuing"
                hasOutputProfile = FileManager.default.fileExists(atPath: directory.appendingPathComponent("App.mobileprovision").path)
                hasCertificate = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Identity.p12").path)
                hasProfile = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Profile.mobileprovision").path)
                message = "Signing file stored in protected, backup-excluded local storage. The private key is not packaged. Signed Crumb updates include the public certificate, profile and target license; keep those updates private."
            } catch {
                identityValidated = false
                identityStatus = "Could not import this file. Choose it again, or use Share → Crumb from Files."
                message = identityStatus
            }
        }
    }
    func sign(password: String) {
        guard !busy, let helper, hasCertificate, hasProfile else { return }
        busy = true; phase = "signing_helper"; failureMessage = ""; message = "Signing the prepared Crumb package locally…"
        Task {
            defer { busy = false }
            do {
                let directory = try workspace(), identityDirectory = try signingStorage()
                let destination = directory.appendingPathComponent("Crumb-Signed.ipa")
                let bundleID = Bundle.main.bundleIdentifier ?? ""
                let teamID = try currentTeam()
                let text = try await Task.detached(priority: .userInitiated) {
                    let unpacked = directory.appendingPathComponent("Unpacked", isDirectory: true)
                    defer { try? FileManager.default.removeItem(at: unpacked) }
                    let unpackText = try FairPlayNative.run("unpack-system-helper", source: helper, destination: unpacked)
                    let decoder = JSONDecoder(); decoder.keyDecodingStrategy = .convertFromSnakeCase
                    let envelope = try decoder.decode(SystemEnvelope<SystemSigningFolderResponse>.self, from: Data(unpackText.utf8))
                    guard envelope.ok, let folder = envelope.result else { throw CocoaError(.fileReadCorruptFile) }
                    let app = unpacked.appendingPathComponent(folder.bundlePath)
                    let key = identityDirectory.appendingPathComponent("Identity.p12")
                    let profile = identityDirectory.appendingPathComponent("Profile.mobileprovision")
                    let code = app.path.withCString { appPath in key.path.withCString { keyPath in
                        profile.path.withCString { profilePath in password.withCString { secret in
                            bundleID.withCString { expected in teamID.withCString { team in
                                CrumbSignPreparedBundle(appPath, keyPath, profilePath, secret, expected, team)
                            } }
                        } }
                    } }
                    guard code == 0 else { throw NSError(domain: "Crumb.LocalSigning", code: Int(code)) }
                    return try FairPlayNative.run("archive-system-helper", source: unpacked, destination: destination, ipa: helper)
                }.value
                let object = try JSONSerialization.jsonObject(with: Data(text.utf8)) as? [String: Any]
                guard object?["ok"] as? Bool == true else { throw CocoaError(.fileReadCorruptFile) }
                signedHelper = destination
                phase = "awaiting_installation"
                identityValidated = true
                try save()
                message = "Local signature created and purchased payload verified unchanged. Request installation below; iOS decides whether the identity and provisioning profile are accepted."
            } catch {
                identityValidated = false; phase = "signing_failed"
                failureMessage = Self.signingError((error as NSError).code)
                message = failureMessage
            }
        }
    }
    func restore(restoringPackages: Bool = true) {
        guard !busy, let root = try? storage() else { return }
        retireSharedStaging()
        if let directory = try? signingStorage() {
            hasCertificate = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Identity.p12").path)
            hasProfile = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Profile.mobileprovision").path)
            hasOutputProfile = FileManager.default.fileExists(atPath: directory.appendingPathComponent("App.mobileprovision").path)
        }
        if let docs = try? FileManager.default.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: true) {
            try? FairPlayNative.privateDirectory(docs.appendingPathComponent("Signing Inbox", isDirectory: true))
        }
        guard restoringPackages,
              let data = try? Data(contentsOf: root.appendingPathComponent("latest.json")),
              let saved = try? JSONSerialization.jsonObject(with: data) as? [String: String] else { return }
        func item(_ key: String) -> URL? {
            guard let path = saved[key], path.split(separator: "/").count == 2,
                  path.split(separator: "/").allSatisfy({ $0 != "." && $0 != ".." }) else { return nil }
            let url = root.appendingPathComponent(path)
            return FileManager.default.fileExists(atPath: url.path) ? url : nil
        }
        requestedRecordID = saved["record_id"] ?? ""
        helper = item("helper"); output = item("output"); signedHelper = item("signed"); signedOutput = item("signed_output")
        pendingJobID = saved["job_id"] ?? ""; phase = saved["phase"] ?? "idle"
        installAfterDecryption = saved["install_after"] == "true"
        readyToSignOutput = output != nil && signedOutput == nil && installAfterDecryption
        outputBundleID = saved["output_bundle"] ?? ""; outputVersion = saved["output_version"] ?? ""
        preparedDescription = saved["prepared_description"] ?? "Previously prepared target"
        exportDescription = saved["export_description"] ?? "Previously processed target"
        reportFile = item("report")
        if let path = reportFile, let text = try? String(contentsOf: path, encoding: .utf8) { report = text }
    }
    private func save(reportURL: URL? = nil) throws {
        func relative(_ url: URL) -> String { url.deletingLastPathComponent().lastPathComponent + "/" + url.lastPathComponent }
        var value = [String: String]()
        if let helper { value["helper"] = relative(helper) }
        if let output { value["output"] = relative(output) }
        if let signedHelper { value["signed"] = relative(signedHelper) }
        if let signedOutput { value["signed_output"] = relative(signedOutput) }
        value["record_id"] = requestedRecordID
        value["job_id"] = pendingJobID; value["phase"] = phase
        value["install_after"] = installAfterDecryption ? "true" : "false"
        value["output_bundle"] = outputBundleID; value["output_version"] = outputVersion
        value["prepared_description"] = preparedDescription
        value["export_description"] = exportDescription
        if let reportURL { reportFile = reportURL }
        if let reportFile { value["report"] = relative(reportFile) }
        try FairPlayNative.privateWrite(JSONEncoder().encode(value), to: storage().appendingPathComponent("latest.json"))
    }
    func prepare(_ original: URL) {
        guard !busy else { return }; busy = true; phase = "preparing"; failureMessage = ""; message = "Preparing a private Crumb package on this iPhone…"
        Task {
            var continueSigning = false
            defer { busy = false; if continueSigning { sign(password: signingPassword) } }
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
                signedHelper = nil; output = nil; signedOutput = nil; report = ""; reportFile = nil; phase = "prepared"
                if let job = object?["result"] as? [String: Any] {
                    pendingJobID = job["job_id"] as? String ?? ""
                    outputBundleID = job["bundle_id"] as? String ?? ""
                    outputVersion = job["build"] as? String ?? ""
                    preparedDescription = [job["bundle_id"],job["version"],job["build"]].compactMap { $0 as? String }.joined(separator: " · ")
                }
                try save()
                continueSigning = automaticSigning && hasCertificate && hasProfile
                automaticSigning = false
                message = "Package ready. Configure your local signing identity, sign it, and request installation. Then open Settings → FairPlay on iPhone. The prepared package includes this app's private license; keep it private."
            } catch {
                automaticSigning = false; phase = "preparation_failed"
                failureMessage = "Could not prepare this app for decryption. Your original IPA is saved in Downloads."
                message = failureMessage
            }
        }
    }
    // Keep cleanup ownership in private app storage BEFORE copying a license
    // into a dedicated shared container. A kill/relaunch can then retire it.
    private func retireSharedStaging() {
        guard let ledger = try? storage().appendingPathComponent("shared-staging.json"),
              let data = try? Data(contentsOf: ledger), data.count < 32768,
              let names = try? JSONDecoder().decode([String].self, from: data), names.count <= 16,
              let root = FileManager.default.containerURL(forSecurityApplicationGroupIdentifier: "group." + currentBundleID + ".fairplay") else { return }
        var retained = [String]()
        for name in names {
            guard name.hasPrefix("Crumb-FairPlay-"), UUID(uuidString: String(name.dropFirst(15))) != nil,
                  name == URL(fileURLWithPath: name).lastPathComponent else { continue }
            let directory = root.appendingPathComponent(name, isDirectory: true)
            do {
                if FileManager.default.fileExists(atPath: directory.path) { try FileManager.default.removeItem(at: directory) }
            } catch { retained.append(name) }
        }
        if let data = try? JSONEncoder().encode(retained) { try? FairPlayNative.privateWrite(data, to: ledger) }
    }
    private func makeSharedWorkspaces() -> [URL] {
        var result = [URL](), seen = Set<String>()
        for index in 0..<32 {
            var buffer = [CChar](repeating: 0, count: 1024)
            let status = CrumbCurrentAppGroup(index, &buffer, 1024)
            if status == 1 { break }; if status != 0 { continue }
            let group = String(cString: buffer)
            // Do not stage private licenses in a signer's general-purpose group.
            // Only a dedicated Crumb namespace is eligible, and iOS must have
            // actually granted it; no entitlement is added or forged here.
            guard group == "group." + currentBundleID + ".fairplay" else { continue }
            guard let root = FileManager.default.containerURL(forSecurityApplicationGroupIdentifier: group),
                  seen.insert(root.standardizedFileURL.path).inserted else { continue }
            let directory = root.appendingPathComponent("Crumb-FairPlay-" + UUID().uuidString, isDirectory: true)
            var created = false
            do {
                try FileManager.default.createDirectory(at: directory, withIntermediateDirectories: false,
                    attributes: [.posixPermissions: 0o700, .protectionKey: FileProtectionType.complete])
                created = true
                var values = URLResourceValues(); values.isExcludedFromBackup = true
                var protected = directory; try protected.setResourceValues(values)
                let ledger = try storage().appendingPathComponent("shared-staging.json")
                let prior = (try? Data(contentsOf: ledger)).flatMap { try? JSONDecoder().decode([String].self, from: $0) } ?? []
                guard prior.count < 16 else { throw CocoaError(.fileWriteOutOfSpace) }
                try FairPlayNative.privateWrite(JSONEncoder().encode(prior + [directory.lastPathComponent]), to: ledger)
                result.append(directory)
                if result.count == 2 { break }
            } catch { if created { try? FileManager.default.removeItem(at: directory) } }
        }
        return result
    }
    func decrypt(identifiers: [String], localIPA: URL? = nil, description: String? = nil) {
        guard !busy, localIPA != nil || installedTarget != nil else { return }
        let bundledTarget = installedTarget
        let label = description ?? installedDescription
        let sharedWorkspaces = localIPA == nil ? [] : makeSharedWorkspaces()
        busy = true; readyToSignOutput = false; phase = "processing"; try? save()
        message = "Reading all protected pages through this iPhone's FairPlay…"
        Task {
            var signNext = false
            defer {
                retireSharedStaging()
                busy = false; if signNext { signExport() }
            }
            do {
                let directory = try workspace()
                let destination = directory.appendingPathComponent("App-Unsigned.ipa")
                let text = try await Task.detached(priority: .userInitiated) {
                    let extracted = directory.appendingPathComponent("Target.app", isDirectory: true)
                    defer { if localIPA != nil { try? FileManager.default.removeItem(at: extracted) } }
                    let source: URL
                    if let localIPA {
                        let accessing = localIPA.startAccessingSecurityScopedResource()
                        defer { if accessing { localIPA.stopAccessingSecurityScopedResource() } }
                        let preparation = try FairPlayNative.run("prepare-system-local", destination: extracted, ipa: localIPA)
                        let decoder = JSONDecoder(); decoder.keyDecodingStrategy = .convertFromSnakeCase
                        let envelope = try decoder.decode(SystemEnvelope<SystemLocalPrepared>.self, from: Data(preparation.utf8))
                        guard envelope.ok, let target = envelope.result,
                              target.schema == "crumb-system-fairplay-local-v1", !target.requiresSigning else {
                            throw CocoaError(.fileReadCorruptFile)
                        }
                        source = extracted
                    } else if let bundledTarget { source = bundledTarget }
                    else { throw CocoaError(.fileNoSuchFile) }
                    var attempts = [[String: Any]]()
                    do { return try SystemNativeProcessor.process(source: source, directory: directory, destination: destination, identifiers: identifiers,
                        route: localIPA == nil ? "installed_resource" : "data_container") }
                    catch {
                        guard (error as NSError).domain == "Crumb.SystemFairPlay" else { throw error }
                        attempts.append(SystemNativeProcessor.safeFailure(error, route: localIPA == nil ? "installed_resource" : "data_container")) }
                    // A legitimate shared container is a distinct storage class
                    // from the normal data container. Its accessibility to
                    // fairplayd is UNVERIFIED, not assumed from our own access.
                    for group in sharedWorkspaces {
                        do {
                            let sharedTarget = group.appendingPathComponent("Target.app", isDirectory: true)
                            try FileManager.default.copyItem(at: source, to: sharedTarget)
                            return try SystemNativeProcessor.process(source: sharedTarget, directory: directory, destination: destination, identifiers: identifiers, route: "authorized_app_group")
                        } catch { attempts.append(SystemNativeProcessor.safeFailure(error, route: "authorized_app_group")) }
                    }
                    throw NSError(domain: "Crumb.SystemLocations", code: 1, userInfo: ["attempts": attempts])
                }.value
                let object = try JSONSerialization.jsonObject(with: Data(text.utf8)) as? [String: Any]
                guard object?["ok"] as? Bool == true else { throw CocoaError(.fileReadCorruptFile) }
                report = text
                let reportURL = directory.appendingPathComponent("Export-report.json")
                try FairPlayNative.privateWrite(Data(text.utf8), to: reportURL)
                output = destination; phase = "exported"
                readyToSignOutput = installAfterDecryption
                signNext = installAfterDecryption && !signingPassword.isEmpty
                exportDescription = label
                try save(reportURL: reportURL)
                message = "Native page processing and export checks passed. Known licensing metadata and embedded signatures were removed. Sign the output before installation. Unknown app-defined account data and plaintext authentication are not certified."
            } catch {
                phase = "processing_failed"
                let native = error as NSError
                var failure: [String: Any] = ["schema": "crumb-system-attempt-v2", "ok": false,
                    "route": localIPA == nil ? "installed_resource" : "data_container", "code": native.code,
                    "stage": native.domain == "Crumb.SystemLocations" ? "native_attempts" :
                        (native.domain == "Crumb.SystemFairPlay" ? SystemNativeFailureStage.classify(native.userInfo) : "preparation_or_validation"),
                    "signing_requested": false, "installation_requested": false,
                    "authorized_shared_locations": sharedWorkspaces.count]
                if let attempts = native.userInfo["attempts"] as? [[String: Any]] { failure["attempts"] = attempts }
                for key in ["native_remap_return", "native_remap_errno", "native_read_error", "native_mapping_attempted", "native_mapping_errno", "native_remap_attempted", "pages_total", "pages_read", "pages_changed"] {
                    if let count = native.userInfo[key] as? NSNumber { failure[key] = count }
                }
                if let data = try? JSONSerialization.data(withJSONObject: failure, options: [.sortedKeys, .prettyPrinted]) {
                    report = String(decoding: data, as: UTF8.self)
                    if let directory = try? workspace() {
                        let file = directory.appendingPathComponent("Attempt.json")
                        try? FairPlayNative.privateWrite(data, to: file)
                        try? save(reportURL: file)
                    }
                }
                failureMessage = "Could not decrypt this app on this iPhone. Your original IPA is saved in Downloads."
                message = failureMessage
            }
        }
    }
    func signExport() {
        guard !busy, let output, hasCertificate, hasProfile else { return }
        busy = true; message = "Signing a private installation copy…"
        let password = signingPassword
        Task {
            defer { busy = false }
            do {
                let directory = try workspace(), identity = try signingStorage()
                let profile = identity.appendingPathComponent(hasOutputProfile ? "App.mobileprovision" : "Profile.mobileprovision")
                let destination = directory.appendingPathComponent("App-Decrypted-Signed.ipa")
                let text = try await Task.detached(priority: .userInitiated) {
                    let unpacked = directory.appendingPathComponent("Unpacked", isDirectory: true)
                    defer { try? FileManager.default.removeItem(at: unpacked) }
                    let text = try FairPlayNative.run("unpack-system-export", source: output, destination: unpacked)
                    let decoder = JSONDecoder(); decoder.keyDecodingStrategy = .convertFromSnakeCase
                    let envelope = try decoder.decode(SystemEnvelope<SystemSigningFolderResponse>.self, from: Data(text.utf8))
                    guard envelope.ok, let folder = envelope.result else { throw CocoaError(.fileReadCorruptFile) }
                    let app = unpacked.appendingPathComponent(folder.bundlePath)
                    let code = CrumbSignPreparedBundle(app.path, identity.appendingPathComponent("Identity.p12").path,
                        profile.path, password, folder.bundleId, "")
                    guard code == 0 else { throw NSError(domain: "Crumb.LocalSigning", code: Int(code)) }
                    return try FairPlayNative.run("archive-system-export", source: unpacked, destination: destination, ipa: output)
                }.value
                let object = try JSONSerialization.jsonObject(with: Data(text.utf8)) as? [String: Any]
                guard object?["ok"] as? Bool == true else { throw CocoaError(.fileReadCorruptFile) }
                if let info = object?["result"] as? [String: Any] {
                    outputBundleID = info["bundle_id"] as? String ?? ""
                    outputVersion = info["build"] as? String ?? ""
                }
                signedOutput = destination; readyToSignOutput = false; phase = "signed_export"; try save()
                message = "Installation copy signed. iOS must still accept it. Your unsigned export stays separate from this copy, which contains your public certificate and profile."
            } catch { message = Self.signingError((error as NSError).code) }
        }
    }
}

struct SystemFairPlayView: View {
    var record: DownloadRecord? = nil
    var startAction: DownloadDecryptionAction? = nil
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @StateObject private var session = SystemFairPlaySession.shared
    @State private var started = false
    @State private var coreMediaRunning = false
    @State private var coreMediaSummary = ""
    private var identifiers: [String] {
        guard let account = appData.ipaTool?.account else { return [] }
        return [account.email, account.name, account.dsid, account.guid, account.passwordToken]
            .filter { $0.utf8.count >= 6 && $0.utf8.count <= 4096 }
    }
    private var visibleOutput: URL? {
        guard started, record == nil || session.requestedRecordID == record?.id else { return nil }
        return session.output
    }
    private var progressTitle: String {
        switch session.phase {
        case "preparing": return "Preparing your app…"
        case "signing_helper": return "Preparing authorization…"
        case "processing": return "Decrypting your app…"
        default: return "Working…"
        }
    }
    private func proceed() {
        _ = session.continueDecryption(record: record, identifiers: identifiers)
    }
    var body: some View {
        NavigationStack {
            VStack(spacing: 24) {
                Spacer()
                Image(systemName: visibleOutput != nil ? "checkmark.circle.fill" : "lock.open")
                    .font(.system(size: 56, weight: .light))
                    .foregroundStyle(visibleOutput != nil ? Color.green : Color.accentColor)
                    .accessibilityHidden(true)
                Text(record?.appName ?? "Decrypt app")
                    .font(.system(.title2, design: .rounded, weight: .bold))
                    .multilineTextAlignment(.center)
                if session.busy {
                    ProgressView()
                    Text(progressTitle).font(.headline)
                    Text("Keep Crumb open.").foregroundStyle(.secondary)
                } else if let output = visibleOutput {
                    Text("Your decrypted IPA is ready.").foregroundStyle(.secondary)
                    ShareLink(item: output) { Label("Save IPA", systemImage: "square.and.arrow.up").frame(maxWidth: .infinity, minHeight: 36) }
                        .buttonStyle(.borderedProminent)
                } else if !session.failureMessage.isEmpty {
                    Text(session.failureMessage).multilineTextAlignment(.center).foregroundStyle(.secondary)
                    Button("Try again", action: proceed).buttonStyle(.borderedProminent)
                    if !session.report.isEmpty {
                        ShareLink(item: session.report) { Label("Share error details", systemImage: "square.and.arrow.up") }
                    }
                } else {
                    Text("Crumb will process the app and check the result before saving it.")
                        .multilineTextAlignment(.center).foregroundStyle(.secondary)
                    Button("Decrypt", action: proceed).buttonStyle(.borderedProminent)
                }
                if let localIPA = record?.fileURL {
                    VStack(spacing: 8) {
                        Button(coreMediaRunning ? "Checking CoreMedia…" : "Test CoreMedia file access") {
                            guard !coreMediaRunning else { return }
                            coreMediaRunning = true
                            coreMediaSummary = ""
                            Task {
                                let result = await Task.detached(priority: .utility) {
                                    CoreMediaLocalProbe.run(input: localIPA)
                                }.value
                                coreMediaSummary = result.summary
                                coreMediaRunning = false
                            }
                        }
                        .buttonStyle(.bordered)
                        .disabled(coreMediaRunning || session.busy)
                        if !coreMediaSummary.isEmpty {
                            Text(coreMediaSummary)
                                .font(.footnote)
                                .foregroundStyle(.secondary)
                                .multilineTextAlignment(.center)
                                .textSelection(.enabled)
                        }
                    }
                }
                Spacer()
            }
            .padding(28).frame(maxWidth: .infinity)
            .navigationTitle("Decrypt").navigationBarTitleDisplayMode(.inline)
            .toolbar { ToolbarItem(placement: .confirmationAction) { Button("Done") { dismiss() }.disabled(session.busy) } }
            .interactiveDismissDisabled(session.busy)
            .task {
                guard !started else { return }; started = true
                session.restore()
                if session.canResume && (record == nil || session.requestedRecordID == record?.id) {
                    session.resume(identifiers: identifiers)
                } else if startAction != nil { proceed() }
            }

        }
    }
}

struct SigningIdentityView: View {
    var continuation: (() -> Void)? = nil
    @Environment(\.dismiss) private var dismiss
    @StateObject private var session = SystemFairPlaySession.shared
    @State private var picking = false
    var body: some View {
        NavigationStack {
            Form {
                if continuation != nil {
                    Section {
                        Text("This decryption method needs iOS to authorize a private Crumb update. Add your signing files once, accept installation, then reopen Crumb.")
                            .font(.subheadline).foregroundStyle(.secondary)
                    }
                }
                Section {
                    Button { session.importName = "Identity.p12"; picking = true } label: {
                        LabeledContent("P12 certificate", value: session.hasCertificate ? "Added" : "Choose file")
                    }
                    Button { session.importName = "Profile.mobileprovision"; picking = true } label: {
                        LabeledContent("Provisioning profile", value: session.hasProfile ? "Added" : "Choose file")
                    }
                    SecureField("P12 password", text: $session.signingPassword)
                        .textInputAutocapitalization(.never).autocorrectionDisabled()
                } footer: {
                    Text("Your signing files stay on this iPhone. Use the identity that signed Crumb. The password is kept only while Crumb is open.")
                }
                if session.identityValidated {
                    Section { Label("Certificate checked", systemImage: "checkmark.circle").foregroundStyle(.green) }
                } else if session.identityStatus != "Not checked" && session.identityStatus != "Check this identity before continuing" {
                    Section { Text(session.identityStatus).font(.footnote).foregroundStyle(.secondary) }
                }
                Section {
                    Button(continuation == nil ? "Save" : "Continue") {
                        session.validateIdentity {
                            if let continuation { continuation() } else { dismiss() }
                        }
                    }
                    .frame(maxWidth: .infinity).disabled(!session.hasCertificate || !session.hasProfile)
                }
            }
            .disabled(session.busy)
            .overlay { if session.busy { ProgressView().padding().background(.regularMaterial, in: RoundedRectangle(cornerRadius: 12)) } }
            .navigationTitle("Certificate").navigationBarTitleDisplayMode(.inline)
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Cancel") { dismiss() }.disabled(session.busy) } }
            .interactiveDismissDisabled(session.busy)
            .task { session.restore(restoringPackages: false) }
            .sheet(isPresented: $picking) {
                SigningFilePicker(selected: { file in
                    picking = false
                    session.importSigningFile(file, certificate: session.importName == "Identity.p12")
                }, cancelled: { picking = false })
            }
        }
    }
}
