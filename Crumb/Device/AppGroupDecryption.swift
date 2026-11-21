import Foundation
import ZIPFoundation

// The system, not a folder name or a plist entry, grants access to this container.
// Never stage in Crumb.app: modifying its sealed resources invalidates signing.
enum AppGroupDecryption {
    static let configuredGroupKey = "FairPlayAppGroupOverride"
    static var signedGroups: [String]? { crumb_group_signed_groups() }
    static var profileGroups: [String] {
        guard let url = Bundle.main.url(forResource: "embedded", withExtension: "mobileprovision"),
              let data = try? Data(contentsOf: url),
              let start = data.range(of: Data("<?xml".utf8)),
              let end = data.range(of: Data("</plist>".utf8), options: .backwards),
              start.lowerBound < end.upperBound,
              let profile = try? PropertyListSerialization.propertyList(from: data.subdata(in: start.lowerBound..<end.upperBound), format: nil) as? [String: Any],
              let entitlements = profile["Entitlements"] as? [String: Any],
              let groups = entitlements["com.apple.security.application-groups"] as? [String] else { return [] }
        return groups
    }
    static func candidates(signed: [String]?, profile: [String], configured: [String]) -> [String] {
        var seen = Set<String>()
        return ((signed ?? []) + configured + profile).map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty && !$0.contains("*") && seen.insert($0).inserted }
    }
    static var groupCandidates: [String] {
        candidates(signed: signedGroups, profile: profileGroups, configured: [
            UserDefaults.standard.string(forKey: configuredGroupKey) ?? "",
            Bundle.main.object(forInfoDictionaryKey: "FairPlayAppGroup") as? String ?? ""
        ])
    }
    static func container() -> URL? {
        let fm = FileManager.default
        for group in groupCandidates {
            guard let url = fm.containerURL(forSecurityApplicationGroupIdentifier: group) else { continue }
            let probe = url.appendingPathComponent(".crumb-group-probe-" + UUID().uuidString)
            defer { try? fm.removeItem(at: probe) }
            if (try? Data().write(to: probe, options: .withoutOverwriting)) != nil { return url }
        }
        return nil
    }

    static func fail(_ message: String) -> DeviceOperationError { .message(message) }

    // Validate the relative manifest path before joining it to the app URL.
    // standardizing/resolving an existing base and a not-yet-created destination
    // can produce different /var vs /private/var aliases on Darwin. Containment
    // must not depend on comparing those physical path spellings.
    static func sinfURL(_ relative: String, in app: URL) throws -> URL {
        let parts = relative.split(separator: "/", omittingEmptySubsequences: false)
        guard !relative.isEmpty, !relative.contains("\\"), !relative.contains("\0"),
              parts.allSatisfy({ !$0.isEmpty && $0 != "." && $0 != ".." }) else {
            throw fail("Invalid SINF relative path: absolute, empty or traversal component.")
        }
        var destination = app
        for part in parts {
            destination.appendPathComponent(String(part))
            // The archive already rejects links. Also refuse a substituted link
            // before using a manifest destination, including any parent component.
            if let attributes = try? FileManager.default.attributesOfItem(atPath: destination.path),
               attributes[.type] as? FileAttributeType == .typeSymbolicLink {
                throw fail("Invalid SINF relative path: symbolic link.")
            }
        }
        return destination
    }


    // Blocking operation. Its task owns the workspace until completion, including
    // cancellation; no later job may reset the C cancellation flag before then.
    static func run(source: URL, root: URL, cancelled: () -> Bool,
                    progress: (Double) -> Void) throws -> URL {
        let fm = FileManager.default
        let work = root.resolvingSymlinksInPath().appendingPathComponent("crumb-fairplay-" + UUID().uuidString, isDirectory: true)
        try fm.createDirectory(at: work, withIntermediateDirectories: true)
        defer { try? fm.removeItem(at: work) }
        let unpacked = work.appendingPathComponent("original", isDirectory: true)
        let packed = work.appendingPathComponent("plain", isDirectory: true)
        try fm.createDirectory(at: unpacked, withIntermediateDirectories: true)
        try fm.createDirectory(at: packed, withIntermediateDirectories: true)
        let archive = try Archive(url: source, accessMode: .read)
        var seen = Set<String>()
        var total: UInt64 = 0
        // Validate the complete directory before extracting anything. Preserve
        // original root resources, SC_Info and nested bundles in the output.
        for entry in archive {
            let parts = entry.path.split(separator: "/", omittingEmptySubsequences: false)
            let normalized = entry.path.hasSuffix("/") ? String(entry.path.dropLast()) : entry.path
            guard !entry.path.hasPrefix("/"), !entry.path.contains("\\"),
                  !parts.contains(".."), !parts.contains("."), !normalized.isEmpty,
                  !normalized.split(separator: "/", omittingEmptySubsequences: false).contains(""),
                  entry.type != .symlink, seen.insert(normalized).inserted,
                  seen.count <= 100_000, entry.uncompressedSize <= 512 * 1024 * 1024,
                  entry.uncompressedSize <= 4_000_000_000 - total else {
                throw fail("Unsafe or oversized IPA archive.")
            }
            total += entry.uncompressedSize
        }
        if let free = try? root.resourceValues(forKeys: [.volumeAvailableCapacityForImportantUsageKey]).volumeAvailableCapacityForImportantUsage,
           free < Int64(total) * 4 { throw fail("Not enough space to unpack, decrypt and verify this IPA.") }
        for entry in archive {
            try Task.checkCancellation()
            if cancelled() { throw CancellationError() }
            guard try archive.extract(entry, to: unpacked.appendingPathComponent(entry.path)) == entry.checksum else { throw fail("IPA checksum verification failed.") }
        }
        let payload = unpacked.appendingPathComponent("Payload", isDirectory: true)
        let apps = try fm.contentsOfDirectory(at: payload, includingPropertiesForKeys: nil).filter { $0.pathExtension == "app" }
        guard apps.count == 1 else { throw fail("Expected exactly one Payload app.") }
        let app = apps[0]
        // Use the store's SINF replication manifest. Do not replace an existing
        // per-binary license with another binary's license.
        let manifestURL = app.appendingPathComponent("SC_Info/Manifest.plist")
        var primarySinf: URL?
        if let data = try? Data(contentsOf: manifestURL),
           let manifest = try? PropertyListSerialization.propertyList(from: data, format: nil) as? [String: Any],
           let sources = manifest["SinfPaths"] as? [String], let first = sources.first {
            let sourceSinf = try sinfURL(first, in: app)
            primarySinf = sourceSinf
            for path in manifest["SinfReplicationPaths"] as? [String] ?? [] {
                // Replication hints may be stale or refer to removed extensions.
                // Ignore unusable hints without ever writing an unsafe path;
                // actual Mach-O files below still get their own checked SINF.
                guard path.lowercased().hasSuffix(".sinf"),
                      let destination = try? sinfURL(path, in: app) else { continue }
                if !fm.fileExists(atPath: destination.path) {
                    try fm.createDirectory(at: destination.deletingLastPathComponent(), withIntermediateDirectories: true)
                    try fm.copyItem(at: sourceSinf, to: destination)
                }
            }
        }
        guard let enumerator = fm.enumerator(atPath: unpacked.path) else { throw fail("Cannot enumerate IPA.") }
        // The relative enumerator avoids /var vs /private/var spelling changes
        // in URL-based enumeration of App Group containers on Darwin.
        var files: [String] = []
        for case let relative as String in enumerator { files.append(relative) }
        let appPrefix = "Payload/" + app.lastPathComponent + "/"
        var decrypted = 0
        for (index, relative) in files.enumerated() {
            if cancelled() { throw CancellationError() }
            let parts = relative.split(separator: "/", omittingEmptySubsequences: false)
            guard !relative.contains("\\"), parts.allSatisfy({ !$0.isEmpty && $0 != "." && $0 != ".." }) else {
                throw fail("Unexpected archive enumeration path.")
            }
            let url = unpacked.appendingPathComponent(relative)
            let destination = packed.appendingPathComponent(relative)
            let values = try url.resourceValues(forKeys: [.isDirectoryKey, .isSymbolicLinkKey])
            guard values.isSymbolicLink != true else { throw fail("IPA contains a symbolic link.") }
            if values.isDirectory == true {
                try fm.createDirectory(at: destination, withIntermediateDirectories: true)
                continue
            }
            try fm.createDirectory(at: destination.deletingLastPathComponent(), withIntermediateDirectories: true)
            if relative.hasPrefix(appPrefix), let primarySinf {
                let handle = try FileHandle(forReadingFrom: url)
                let head = try handle.read(upToCount: 4)
                try handle.close()
                if head == Data([0xcf, 0xfa, 0xed, 0xfe]) {
                    let sinf = url.deletingLastPathComponent().appendingPathComponent("SC_Info/" + url.lastPathComponent + ".sinf")
                    if !fm.fileExists(atPath: sinf.path) {
                        try fm.createDirectory(at: sinf.deletingLastPathComponent(), withIntermediateDirectories: true)
                        try fm.copyItem(at: primarySinf, to: sinf)
                        let sinfRelative = (relative as NSString).deletingLastPathComponent + "/SC_Info/" + url.lastPathComponent + ".sinf"
                        let outputSinf = packed.appendingPathComponent(sinfRelative)
                        try fm.createDirectory(at: outputSinf.deletingLastPathComponent(), withIntermediateDirectories: true)
                        try fm.copyItem(at: sinf, to: outputSinf)
                    }
                }
            }
            var encrypted = ObjCBool(false)
            // This shares the helper's bounded Mach-O parser and complete page
            // checks. Unsupported FAT/32-bit images fail instead of partial export.
            var nativeError: NSError?
            guard crumb_group_decrypt(url, destination, &encrypted, &nativeError) else {
                if let nativeError { throw fail("\(relative): \(nativeError.localizedDescription)") }
                throw fail("FairPlay rejected this image.")
            }
            if encrypted.boolValue { decrypted += 1 }
            progress(0.1 + 0.75 * Double(index + 1) / Double(max(files.count, 1)))
        }
        guard decrypted > 0 else { throw fail("No encrypted images found; no decrypted result was published.") }
        // The original encrypted archive is private and contains account-bound
        // FairPlay license data. Export only from the verified plaintext tree:
        // remove all Apple purchase metadata, SINF/SUPP and per-bundle SC_Info.
        // Signed installation uses the decrypted binary, not the purchase token.
        guard let privateEnumerator = fm.enumerator(atPath: packed.path) else {
            throw fail("Cannot inspect export privacy.")
        }
        let privateEntries = privateEnumerator.allObjects.compactMap { $0 as? String }
        for relative in privateEntries.sorted(by: { $0.count > $1.count }) {
            let components = relative.split(separator: "/").map(String.init)
            let name = components.last ?? ""
            let accountMetadata = ["iTunesMetadata.plist", "iTunesArtwork", "iTunesArtwork@2x"].contains(name)
            let licenseResource = components.contains("SC_Info")
            guard accountMetadata || licenseResource else { continue }
            let candidate = packed.appendingPathComponent(relative)
            if fm.fileExists(atPath: candidate.path) {
                try fm.removeItem(at: candidate)
            }
        }
        // Fail closed: never publish an export containing original license files.
        if let residual = fm.enumerator(atPath: packed.path) {
            for case let relative as String in residual {
                let parts = relative.split(separator: "/")
                guard !parts.contains("SC_Info"), parts.last != "iTunesMetadata.plist" else {
                    throw fail("Personal Apple license metadata remains in output.")
                }
            }
        }
        if cancelled() { throw CancellationError() }
        let downloads = fm.urls(for: .documentDirectory, in: .userDomainMask)[0].appendingPathComponent("Downloads", isDirectory: true)
        try fm.createDirectory(at: downloads, withIntermediateDirectories: true)
        let output = downloads.appendingPathComponent(UUID().uuidString + "-decrypted.ipa")
        var success = false
        defer { if !success { try? fm.removeItem(at: output) } }
        try fm.zipItem(at: packed, to: output, shouldKeepParent: false)
        // Verify CRCs of every packaged entry and exact byte-for-byte content
        // against the verified plaintext tree, before recording success.
        let result = try Archive(url: output, accessMode: .read)
        for entry in result where entry.type == .file {
            if cancelled() { throw CancellationError() }
            let expected = try FileHandle(forReadingFrom: packed.appendingPathComponent(entry.path))
            defer { try? expected.close() }
            let checksum = try result.extract(entry) { bytes in
                guard try expected.read(upToCount: bytes.count) == bytes else { throw fail("Packaged IPA verification failed.") }
            }
            guard checksum == entry.checksum else { throw fail("Output IPA checksum verification failed.") }
            guard (try expected.read(upToCount: 1))?.isEmpty != false else { throw fail("Packaged IPA is truncated.") }
        }
        if cancelled() { throw CancellationError() }
        progress(1)
        success = true
        return output
    }
}
