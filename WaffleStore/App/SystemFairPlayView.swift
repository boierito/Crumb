import SwiftUI
import Combine
import UniformTypeIdentifiers
import MapleSAP

nonisolated private struct SystemEnvelope<T: Decodable>: Decodable {
    let ok: Bool
    let result: T?
}
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

@MainActor
final class SystemFairPlaySession: ObservableObject {
    static let shared = SystemFairPlaySession()
    @Published var signingPassword = ""
    @Published var identityStatus = "Not checked"
    @Published var identityValidated = false
    @Published var installAfterDecryption = false
    @Published var signedOutput: URL?
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
        default: return "Signing could not finish (code \(code)). The original IPA is retained."
        }
    }
    func begin(_ record: DownloadRecord, install: Bool) {
        guard !busy, let original = record.fileURL else { return }
        installAfterDecryption = install
        automaticSigning = true
        prepare(original)
    }
    func validateIdentity() {
        guard !busy, hasCertificate, hasProfile else { return }
        busy = true
        let password = signingPassword
        Task {
            defer { busy = false }
            do {
                let directory = try signingStorage(), team = try currentTeam(), bundle = currentBundleID
                let code = await Task.detached {
                    CrumbValidateSigningIdentity(directory.appendingPathComponent("Identity.p12").path,
                        directory.appendingPathComponent("Profile.mobileprovision").path, password, bundle, team)
                }.value
                identityValidated = code == 0
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
                if removeInboxCopy { try? FileManager.default.removeItem(at: selected) }
                identityValidated = false; identityStatus = "Check this identity before continuing"
                hasOutputProfile = FileManager.default.fileExists(atPath: directory.appendingPathComponent("App.mobileprovision").path)
                hasCertificate = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Identity.p12").path)
                hasProfile = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Profile.mobileprovision").path)
                message = "Signing file stored in protected, backup-excluded local storage. The private key is not packaged. Signed Crumb updates include the public certificate, profile and target license; keep those updates private."
            } catch { message = "Could not import the signing file. No private error text was logged." }
        }
    }
    func sign(password: String) {
        guard !busy, let helper, hasCertificate, hasProfile else { return }
        busy = true; message = "Signing the prepared Crumb package locally…"
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
                identityValidated = false
                message = Self.signingError((error as NSError).code)
            }
        }
    }
    func restore(restoringPackages: Bool = true) {
        guard !busy, let root = try? storage() else { return }
        if let directory = try? signingStorage() {
            hasCertificate = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Identity.p12").path)
            hasProfile = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Profile.mobileprovision").path)
            hasOutputProfile = FileManager.default.fileExists(atPath: directory.appendingPathComponent("App.mobileprovision").path)
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
        helper = item("helper"); output = item("output"); signedHelper = item("signed"); signedOutput = item("signed_output")
        pendingJobID = saved["job_id"] ?? ""; phase = saved["phase"] ?? "idle"
        installAfterDecryption = saved["install_after"] == "true"
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
        guard !busy else { return }; busy = true; message = "Preparing a private Crumb package on this iPhone…"
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
                message = "Preparation failed. No successful package was accepted; the original download is retained."
            }
        }
    }
    func decrypt(identifiers: [String]) {
        guard !busy, let source = installedTarget else { return }
        busy = true; phase = "processing"; try? save()
        message = "Reading all protected pages through this iPhone's FairPlay…"
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
                output = destination; phase = "exported"
                exportDescription = installedDescription
                try save(reportURL: reportURL)
                message = "Native page processing and export checks passed. Known licensing metadata and embedded signatures were removed. Sign the output before installation. Unknown app-defined account data and plaintext authentication are not certified."
            } catch {
                phase = "processing_failed"; try? save()
                message = "Native processing or export validation failed (code \((error as NSError).code)). No partial IPA was accepted. The installed target remains available for a retry."
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
                signedOutput = destination; phase = "signed_export"; try save()
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
    @State private var picker = false
    @State private var identity = false
    @State private var installing = false
    @State private var installingOutput = false
    @State private var started = false
    private var identifiers: [String] {
        guard let account = appData.ipaTool?.account else { return [] }
        return [account.email, account.name, account.dsid, account.guid, account.passwordToken]
            .filter { $0.utf8.count >= 6 && $0.utf8.count <= 4096 }
    }
    var body: some View {
        NavigationStack {
            List {
                Section {
                    VStack(alignment: .leading, spacing: 12) {
                        Label("On this iPhone", systemImage: "iphone.gen3.radiowaves.left.and.right")
                            .font(.caption.weight(.semibold)).foregroundStyle(.tint)
                        Text(record?.appName ?? "Your next decrypted IPA")
                            .font(.system(.title2, design: .rounded, weight: .bold))
                        Text(record.map { "Version \($0.version)" } ?? session.installedDescription)
                            .font(.subheadline).foregroundStyle(.secondary)
                        Text("Prepare → Install Crumb update → Decrypt")
                            .font(.footnote.weight(.medium)).foregroundStyle(.secondary)
                    }.padding(.vertical, 8)
                }
                Section("Signing") {
                    Button { identity = true } label: {
                        Label(session.hasCertificate && session.hasProfile ? "Manage signing identity" : "Set up certificate and profile", systemImage: "signature")
                    }
                    Text("The profile must allow \(session.currentBundleID) and use the same team as this installed Crumb.")
                        .font(.footnote).foregroundStyle(.secondary).textSelection(.enabled)
                }
                Section("Prepare") {
                    if let original = record?.fileURL {
                        Button { session.prepare(original) } label: { Label("Prepare this download", systemImage: "shippingbox") }
                    }
                    Button { picker = true } label: { Label("Choose encrypted IPA", systemImage: "folder") }
                    if session.helper != nil {
                        Text(session.preparedDescription).font(.caption).foregroundStyle(.secondary)
                        Button("Sign Crumb update") { session.sign(password: session.signingPassword) }
                            .disabled(!session.hasCertificate || !session.hasProfile)
                    }
                    if session.signedHelper != nil {
                        Button { installing = true } label: { Label("Install Crumb update", systemImage: "arrow.down.app.fill") }
                            .buttonStyle(.borderedProminent)
                        Text("Accept iOS’s installation prompt, then reopen Crumb. Keep this screen open during transfer. Do not uninstall Crumb.")
                            .font(.footnote).foregroundStyle(.secondary)
                    }
                }
                if session.installedTarget != nil {
                    Section("Decrypt") {
                        Text(session.installedDescription).font(.caption).foregroundStyle(.secondary)
                        Button { session.decrypt(identifiers: identifiers) } label: { Label("Decrypt and verify all executables", systemImage: "lock.open") }
                            .disabled(!session.matchesInstalledTarget(record))
                        if !session.matchesInstalledTarget(record) {
                            Text("Prepare and install this download first. A different app is currently installed for processing.").font(.footnote)
                        }
                    }
                }
                if let output = session.output {
                    Section("Decrypted IPA") {
                        Label(session.exportDescription, systemImage: "checkmark.seal.fill").font(.caption)
                        ShareLink(item: output) { Label("Save decrypted IPA", systemImage: "square.and.arrow.up") }
                        Button { session.signExport() } label: { Label("Sign for installation", systemImage: "signature") }
                            .disabled(!session.hasCertificate || !session.hasProfile)
                        if session.installAfterDecryption {
                            Text("To finish Decrypt & install, sign this app using a profile that allows its bundle ID. Re-enter the P12 password in Signing if needed.")
                                .font(.footnote).foregroundStyle(.secondary)
                        }
                        Text("Known license and purchase metadata are removed. Arbitrary app-defined identifiers and plaintext authenticity are not certified.")
                            .font(.footnote).foregroundStyle(.secondary)
                    }
                }
                if session.signedOutput != nil {
                    Section("Installation copy") {
                        Button { installingOutput = true } label: { Label("Install decrypted app", systemImage: "arrow.down.app.fill") }
                            .buttonStyle(.borderedProminent)
                        Text("This private copy contains your signing profile. Use the unsigned decrypted IPA when exporting without signing metadata.")
                            .font(.footnote).foregroundStyle(.secondary)
                    }
                }
                if !session.message.isEmpty { Section { Text(session.message).font(.subheadline).textSelection(.enabled) } }
                if !session.report.isEmpty {
                    Section { DisclosureGroup("Verification report") {
                        Text(session.report).font(.system(.caption, design: .monospaced)).textSelection(.enabled)
                    } }
                }
            }
            .disabled(session.busy)
            .overlay { if session.busy { VStack(spacing: 12) {
                ProgressView(); Text("Keep Crumb open").font(.footnote)
            }.padding(24).background(.regularMaterial, in: RoundedRectangle(cornerRadius: 20)) } }
            .navigationTitle("Decrypt app")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Done") { dismiss() }.disabled(session.busy) } }
            .interactiveDismissDisabled(session.busy)
            .task {
                guard !started else { return }; started = true
                session.restore()
                if session.canResume { session.resume(identifiers: identifiers) }
                else if let action = startAction { session.begin(action.record, install: action.install) }
            }
            .onChange(of: session.signedHelper) { package in
                if package != nil && session.phase == "awaiting_installation" { installing = true }
            }
            .onChange(of: session.signedOutput) { package in
                if package != nil && session.installAfterDecryption { installingOutput = true }
            }
            .sheet(isPresented: $picker) { SigningFilePicker(selected: { picker = false; session.prepare($0) }, cancelled: { picker = false }) }
            .sheet(isPresented: $identity) { SigningIdentityView() }
            .sheet(isPresented: $installing) { if let package = session.signedHelper { SystemPackageInstallationView(package: package) } }
            .sheet(isPresented: $installingOutput) {
                if let package = session.signedOutput {
                    SystemPackageInstallationView(package: package, bundleID: session.outputBundleID,
                        appName: "Decrypted app", version: session.outputVersion)
                }
            }
        }
    }
}

struct SigningIdentityView: View {
    @Environment(\.dismiss) private var dismiss
    @StateObject private var session = SystemFairPlaySession.shared
    @State private var picking = false
    var body: some View {
        NavigationStack {
            Form {
                Section {
                    Label("Your signature stays on your iPhone", systemImage: "key.fill").font(.headline)
                    Text("Import a P12 containing your private key and its matching provisioning profile. Your Apple Account password is not used for signing.")
                        .font(.subheadline).foregroundStyle(.secondary)
                }
                Section("Certificate") {
                    importButton(session.hasCertificate ? "Replace P12" : "Import P12", name: "Identity.p12")
                    SecureField("P12 password", text: $session.signingPassword)
                        .textInputAutocapitalization(.never).autocorrectionDisabled()
                    Text("Password kept in memory only; re-enter after restarting Crumb.").font(.footnote).foregroundStyle(.secondary)
                }
                Section("Crumb profile") {
                    importButton(session.hasProfile ? "Replace Crumb profile" : "Import .mobileprovision", name: "Profile.mobileprovision")
                    Text(session.currentBundleID).font(.caption.monospaced()).textSelection(.enabled)
                    Button("Check certificate and profile") { session.validateIdentity() }
                        .disabled(!session.hasCertificate || !session.hasProfile)
                    Label(session.identityStatus, systemImage: session.identityValidated ? "checkmark.seal" : "info.circle")
                        .font(.footnote).foregroundStyle(session.identityValidated ? Color.green : Color.secondary)
                }
                Section("Profile for decrypted apps · optional") {
                    importButton(session.hasOutputProfile ? "Replace app profile" : "Import app profile", name: "App.mobileprovision")
                    Text("Used only when signing the decrypted app. Without this, Crumb tries the Crumb profile. An appropriate wildcard profile may cover both; extensions also need a compatible profile.")
                        .font(.footnote).foregroundStyle(.secondary)
                }
                Section("If the picker does not open") {
                    Picker("Import as", selection: $session.importName) {
                        Text("P12").tag("Identity.p12")
                        Text("Crumb profile").tag("Profile.mobileprovision")
                        Text("App profile").tag("App.mobileprovision")
                    }
                    Button("Import from Signing Inbox") { session.importInbox() }
                    Text("In Files → On My iPhone → Crumb → Signing Inbox, copy one P12 or profile. Crumb moves that inbox copy to protected storage after import. Your original file stays where you keep it.")
                        .font(.footnote).foregroundStyle(.secondary)
                }
                if !session.message.isEmpty { Section { Text(session.message).font(.footnote).textSelection(.enabled) } }
                Section { Button("Remove signing identity", role: .destructive) { session.removeIdentity() } }
            }
            .disabled(session.busy)
            .navigationTitle("Signing identity").navigationBarTitleDisplayMode(.inline)
            .toolbar { ToolbarItem(placement: .confirmationAction) { Button("Done") { dismiss() }.disabled(session.busy) } }
            .interactiveDismissDisabled(session.busy)
            .task { session.restore(); if session.importName.isEmpty { session.importName = "Identity.p12" } }
            .sheet(isPresented: $picking) {
                SigningFilePicker(selected: { file in
                    picking = false
                    session.importSigningFile(file, certificate: session.importName == "Identity.p12", outputProfile: session.importName == "App.mobileprovision")
                }, cancelled: { picking = false })
            }
        }
    }
    private func importButton(_ title: String, name: String) -> some View {
        Button(title) { session.importName = name; picking = true }
    }
}
