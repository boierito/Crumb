import SwiftUI
import Combine
import UniformTypeIdentifiers
import MapleSAP

nonisolated private struct SystemEnvelope<T: Decodable>: Decodable {
    let ok: Bool
    let result: T?
}
nonisolated private struct SystemSigningFolderResponse: Decodable { let bundlePath: String }
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
    func importSigningFile(_ selected: URL, certificate: Bool) {
        guard !busy else { return }; busy = true
        Task {
            defer { busy = false }
            do {
                let directory = try signingStorage()
                let destination = directory.appendingPathComponent(certificate ? "Identity.p12" : "Profile.mobileprovision")
                try await Task.detached {
                    let accessing = selected.startAccessingSecurityScopedResource()
                    defer { if accessing { selected.stopAccessingSecurityScopedResource() } }
                    let values = try selected.resourceValues(forKeys: [.fileSizeKey, .isRegularFileKey, .isSymbolicLinkKey])
                    guard values.isRegularFile == true, values.isSymbolicLink != true,
                          let size = values.fileSize, size > 0, size <= 8 << 20 else { throw CocoaError(.fileReadTooLarge) }
                    try FairPlayNative.privateWrite(Data(contentsOf: selected), to: destination)
                }.value
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
                var teamBuffer = [CChar](repeating: 0, count: 256)
                guard CrumbCurrentSigningTeam(&teamBuffer, 256) == 0 else {
                    throw NSError(domain: "Crumb.LocalSigning", code: 5)
                }
                let teamID = String(cString: teamBuffer)
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
                try save()
                message = "Local signature created and purchased payload verified unchanged. Request installation below; iOS decides whether the identity and provisioning profile are accepted."
            } catch {
                message = "Local signing rejected (code \((error as NSError).code)). Check the P12 password, certificate expiry and profile for Crumb's existing bundle identifier. No installed success is implied."
            }
        }
    }
    func restore(restoringPackages: Bool = true) {
        guard !busy, let root = try? storage() else { return }
        if let directory = try? signingStorage() {
            hasCertificate = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Identity.p12").path)
            hasProfile = FileManager.default.fileExists(atPath: directory.appendingPathComponent("Profile.mobileprovision").path)
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
        helper = item("helper"); output = item("output"); signedHelper = item("signed")
        preparedDescription = saved["prepared_description"] ?? "Previously prepared target"
        exportDescription = saved["export_description"] ?? "Previously processed target"
        if let path = item("report"), let text = try? String(contentsOf: path, encoding: .utf8) { report = text }
    }
    private func save(reportURL: URL? = nil) throws {
        func relative(_ url: URL) -> String { url.deletingLastPathComponent().lastPathComponent + "/" + url.lastPathComponent }
        var value = [String: String]()
        if let helper { value["helper"] = relative(helper) }
        if let output { value["output"] = relative(output) }
        if let signedHelper { value["signed"] = relative(signedHelper) }
        value["prepared_description"] = preparedDescription
        value["export_description"] = exportDescription
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
                signedHelper = nil; output = nil; report = ""
                if let job = object?["result"] as? [String: Any] {
                    preparedDescription = [job["bundle_id"],job["version"],job["build"]].compactMap { $0 as? String }.joined(separator: " · ")
                }
                try save()
                message = "Package ready. Configure your local signing identity, sign it, and request installation. Then open Settings → FairPlay on iPhone. The prepared package includes this app's private license; keep it private."
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
                exportDescription = installedDescription
                try save(reportURL: reportURL)
                message = "Native page processing and export checks passed. Known licensing metadata and embedded signatures were removed. Sign the output before installation. Unknown app-defined account data and plaintext authentication are not certified."
            } catch {
                message = "Native processing or export validation failed (code \((error as NSError).code)). No partial IPA was accepted. The installed target remains available for a retry."
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
    @State private var importingSigning = false
    @State private var importingCertificate = true
    @State private var password = ""
    @State private var installing = false
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
                        Text(session.preparedDescription).font(.caption).textSelection(.enabled)
                        ShareLink("Save private prepared package", item: helper)
                        Text("This package includes the selected app's original license. Keep it private. Do not uninstall Crumb before installing the signed update.")
                            .font(.footnote)
                    }
                }
                Section("Local signing identity") {
                    Text("Import your own development certificate (.p12) and a valid provisioning profile matching Crumb's existing bundle identifier and signing team. Signing runs inside Crumb; these files are not FairPlay keys and are never uploaded.")
                    Button(session.hasCertificate ? "Replace certificate (.p12)" : "Import certificate (.p12)") {
                        importingCertificate = true; importingSigning = true
                    }
                    Button(session.hasProfile ? "Replace provisioning profile" : "Import provisioning profile") {
                        importingCertificate = false; importingSigning = true
                    }
                    SecureField("P12 password (not saved)", text: $password)
                    Button("Sign prepared Crumb locally") {
                        session.sign(password: password); password = ""
                    }.disabled(session.helper == nil || !session.hasCertificate || !session.hasProfile)
                    if session.signedHelper != nil {
                        Button("Request installation of signed Crumb update") { installing = true }
                    }
                }
                if session.installedTarget != nil {
                    Section("Process the installed target") {
                        Text(session.installedDescription).font(.caption).textSelection(.enabled)
                        Text("The target is installed inside this Crumb package. The iPhone's system Apple Account must own its license.")
                        Button("Decrypt all supported executables and prepare export") {
                            var identifiers = [String]()
                            if let account = appData.ipaTool?.account {
                                identifiers = [account.email, account.name, account.dsid, account.guid, account.passwordToken]
                                    .filter { $0.utf8.count >= 6 && $0.utf8.count <= 4096 }
                            }
                            session.decrypt(identifiers: identifiers)
                        }
                        .disabled(!session.matchesInstalledTarget(record))
                        if !session.matchesInstalledTarget(record) {
                            Text("The installed target belongs to a different download. Prepare and install this download before processing it.").font(.footnote)
                        }
                    }
                }
                if let output = session.output {
                    Section("Unsigned export") {
                        Text(session.exportDescription).font(.caption).textSelection(.enabled)
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
            .task { session.restore(restoringPackages: record == nil) }
            .fileImporter(isPresented: $importing, allowedContentTypes: [UTType(filenameExtension: "ipa") ?? .data]) { result in
                if case .success(let file) = result { session.prepare(file) }
            }
            .fileImporter(isPresented: $importingSigning, allowedContentTypes: [.data]) { result in
                if case .success(let file) = result { session.importSigningFile(file, certificate: importingCertificate) }
            }
            .sheet(isPresented: $installing) {
                if let package = session.signedHelper { SystemPackageInstallationView(package: package) }
            }
        }
    }
}
