import Foundation

// Local download history, separate from the original installation history.
// No account ID, signed URL, cookies or kbsync are serialized.
struct DownloadRecord: Codable, Identifiable {
    var id: String { filename }
    let filename: String
    let appID: String
    let appName: String
    let bundleID: String
    let version: String
    let build: String?
    let externalVersionID: String
    let date: Date
    var decrypted: Bool? = nil
    var signed: Bool? = nil
    var lastSigning: SigningSummary? = nil
    var isDecrypted: Bool { decrypted == true }

    // Encrypted Apple downloads contain account-bound FairPlay material and
    // must never be placed in the Files-visible Documents container.
    static func storageDirectory(decrypted: Bool, create: Bool = false) throws -> URL {
        let fm = FileManager.default
        let base = try fm.url(for: decrypted ? .documentDirectory : .applicationSupportDirectory,
                              in: .userDomainMask, appropriateFor: nil, create: create)
        let url = base.appendingPathComponent(decrypted ? "Downloads" : "PrivateDownloads", isDirectory: true)
        if create { try fm.createDirectory(at: url, withIntermediateDirectories: true) }
        return url
    }

    var fileURL: URL? {
        guard filename == URL(fileURLWithPath: filename).lastPathComponent,
              filename.hasSuffix(".ipa"), !filename.contains(".."),
              let directory = try? Self.storageDirectory(decrypted: isDecrypted) else { return nil }
        let url = directory.appendingPathComponent(filename)
        return FileManager.default.fileExists(atPath: url.path) ? url : nil
    }
    func save() throws {
        guard let url = fileURL else { throw CocoaError(.fileNoSuchFile) }
        try JSONEncoder().encode(self).write(to: url.deletingPathExtension().appendingPathExtension("json"), options: .atomic)
    }
    // Delete only the IPA and matching nonsecret sidecar in Downloads. A stale
    // record is safe to remove even if Files already deleted its IPA.
    func deleteFiles(in directory: URL? = nil) throws {
        guard filename == URL(fileURLWithPath: filename).lastPathComponent,
              filename.hasSuffix(".ipa"), !filename.contains("..") else {
            throw CocoaError(.fileReadInvalidFileName)
        }
        let fm = FileManager.default
        let root: URL
        if let directory { root = directory }
        else { root = try Self.storageDirectory(decrypted: isDecrypted) }
        let ipa = root.appendingPathComponent(filename)
        for file in [ipa, ipa.deletingPathExtension().appendingPathExtension("json")] {
            if fm.fileExists(atPath: file.path) { try fm.removeItem(at: file) }
        }
    }
    // Clean downloaded IPAs and their sidecars, including IPAs without records.
    // Preserve unrelated documents and never follow a substituted Downloads directory.
    static func deleteAllFiles(in directory: URL? = nil) throws {
        if directory == nil {
            for decrypted in [false, true] {
                if let root = try? storageDirectory(decrypted: decrypted) {
                    try deleteAllFiles(in: root)
                }
            }
            return
        }
        let fm = FileManager.default
        let root = directory!
        guard fm.fileExists(atPath: root.path) else { return }
        let values = try root.resourceValues(forKeys: [.isDirectoryKey, .isSymbolicLinkKey])
        guard values.isDirectory == true, values.isSymbolicLink != true else {
            throw CocoaError(.fileReadInvalidFileName)
        }
        let files = try fm.contentsOfDirectory(at: root, includingPropertiesForKeys: [.isDirectoryKey])
        for file in files where file.pathExtension == "ipa" {
            guard try file.resourceValues(forKeys: [.isDirectoryKey]).isDirectory != true else { continue }
            let record = DownloadRecord(filename: file.lastPathComponent, appID: "", appName: "",
                bundleID: "", version: "", build: nil, externalVersionID: "", date: .distantPast)
            try record.deleteFiles(in: root)
        }
        // A sidecar may remain after the IPA was deleted in Files.
        for file in files where file.pathExtension == "json" && fm.fileExists(atPath: file.path) {
            guard let data = try? Data(contentsOf: file), data.count < 65536,
                  let record = try? JSONDecoder().decode(DownloadRecord.self, from: data),
                  file.lastPathComponent == URL(fileURLWithPath: record.filename)
                    .deletingPathExtension().appendingPathExtension("json").lastPathComponent else { continue }
            try record.deleteFiles(in: root)
        }
    }
    private static func migrateLegacyEncryptedDownloads() {
        let fm = FileManager.default
        guard let publicRoot = try? storageDirectory(decrypted: true),
              let privateRoot = try? storageDirectory(decrypted: false, create: true),
              let files = try? fm.contentsOfDirectory(at: publicRoot, includingPropertiesForKeys: nil) else { return }
        for sidecar in files where sidecar.pathExtension == "json" {
            guard let data = try? Data(contentsOf: sidecar), data.count < 65536,
                  let record = try? JSONDecoder().decode(DownloadRecord.self, from: data),
                  !record.isDecrypted,
                  record.filename == sidecar.deletingPathExtension().appendingPathExtension("ipa").lastPathComponent else { continue }
            let oldIPA = publicRoot.appendingPathComponent(record.filename)
            let privateIPA = privateRoot.appendingPathComponent(record.filename)
            let privateJSON = privateRoot.appendingPathComponent(sidecar.lastPathComponent)
            guard fm.fileExists(atPath: oldIPA.path),
                  !fm.fileExists(atPath: privateIPA.path),
                  !fm.fileExists(atPath: privateJSON.path) else { continue }
            // Move the license-bound binary first, then the non-secret sidecar.
            // If the second step fails, the source record is not shared.
            do {
                try fm.moveItem(at: oldIPA, to: privateIPA)
                try fm.moveItem(at: sidecar, to: privateJSON)
            } catch { try? fm.removeItem(at: sidecar) }
        }
    }

    static func load() -> [DownloadRecord] {
        migrateLegacyEncryptedDownloads()
        let fm = FileManager.default
        let directories = [false, true].compactMap { try? storageDirectory(decrypted: $0) }
        var records: [DownloadRecord] = []
        for root in directories {
            guard let files = try? fm.contentsOfDirectory(at: root, includingPropertiesForKeys: nil) else { continue }
            for file in files where file.pathExtension == "json" {
                guard let data = try? Data(contentsOf: file), data.count < 65536,
                      let record = try? JSONDecoder().decode(DownloadRecord.self, from: data),
                      record.fileURL != nil else { continue }
                records.append(record)
            }
        }
        return records.sorted { $0.date > $1.date }
    }
}

struct SigningSummary: Codable {
    let bundleID: String
    let originalBundleID: String
    let profileName: String
    let capabilities: [String]
    let date: Date
    var teamID: String? = nil
    var applicationIdentifier: String? = nil
    var originalKeychainGroups: [String]? = nil
    var keychainGroups: [String]? = nil
    var originalAppGroups: [String]? = nil
    var appGroups: [String]? = nil
    var cloudKitDiagnostics: String? = nil
    static func capabilityName(_ key:String)->String {
        if key == "com.apple.security.application-groups" { return "App Groups" }
        if key == "keychain-access-groups" { return "Shared Keychain" }
        if key == "aps-environment" { return "Push notifications" }
        if key.contains("icloud") || key.contains("ubiquity") { return "iCloud (" + key + ")" }
        if key == "com.apple.developer.associated-domains" { return "Associated domains" }
        return key
    }
    var report: String {
        "Crumb signing report\nOriginal: \(originalBundleID)\nSigned: \(bundleID)\nProfile: \(profileName)\nDate: \(date)\nAdjusted capabilities:\n" + capabilities.joined(separator:"\n") + "\nTeam: \(teamID ?? "")\nApplication identifier: \(applicationIdentifier ?? "")\nOriginal Keychain: \((originalKeychainGroups ?? []).joined(separator:", "))\nSigned Keychain: \((keychainGroups ?? []).joined(separator:", "))\nOriginal App Groups: \((originalAppGroups ?? []).joined(separator:", "))\nSigned App Groups: \((appGroups ?? []).joined(separator:", "))\nCloudKit entitlements (original / profile / signed):\n\(cloudKitDiagnostics ?? "Not recorded")"
    }
}
