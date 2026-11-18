import Foundation
import ZIPFoundation

// The system, not a folder name or a plist entry, grants access to this container.
// Never stage in Crumb.app: modifying its sealed resources invalidates signing.
enum AppGroupDecryption {
    static func container() -> URL? {
        var groups: [String] = []
        if let configured = Bundle.main.object(forInfoDictionaryKey: "FairPlayAppGroup") as? String {
            groups.append(configured)
        }
        if let url = Bundle.main.url(forResource: "embedded", withExtension: "mobileprovision"),
           let data = try? Data(contentsOf: url),
           let start = data.range(of: Data("<?xml".utf8)),
           let end = data.range(of: Data("</plist>".utf8), options: .backwards),
           start.lowerBound < end.upperBound,
           let profile = try? PropertyListSerialization.propertyList(from: data.subdata(in: start.lowerBound..<end.upperBound), format: nil) as? [String: Any],
           let entitlements = profile["Entitlements"] as? [String: Any],
           let authorized = entitlements["com.apple.security.application-groups"] as? [String] {
            groups += authorized
        }
        for group in groups where !group.contains("*") {
            if let url = FileManager.default.containerURL(forSecurityApplicationGroupIdentifier: group) { return url }
        }
        return nil
    }

    static func fail(_ message: String) -> DeviceOperationError { .message(message) }

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
           free < Int64(total) * 3 { throw fail("Not enough space to unpack, decrypt and verify this IPA.") }
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
            func safeURL(_ relative: String) throws -> URL {
                let url = app.appendingPathComponent(relative).standardizedFileURL
                let base = app.resolvingSymlinksInPath().pathComponents
                let parts = url.resolvingSymlinksInPath().pathComponents
                guard parts.count > base.count, Array(parts.prefix(base.count)) == base,
                      !relative.hasPrefix("/") else { throw fail("Invalid SINF replication path.") }
                return url
            }
            let sourceSinf = try safeURL(first)
            primarySinf = sourceSinf
            for path in manifest["SinfReplicationPaths"] as? [String] ?? [] {
                let destination = try safeURL(path)
                if !fm.fileExists(atPath: destination.path) {
                    try fm.createDirectory(at: destination.deletingLastPathComponent(), withIntermediateDirectories: true)
                    try fm.copyItem(at: sourceSinf, to: destination)
                }
            }
        }
        guard let enumerator = fm.enumerator(at: unpacked, includingPropertiesForKeys: [.isDirectoryKey, .isSymbolicLinkKey]) else { throw fail("Cannot enumerate IPA.") }
        var files: [URL] = []
        for case let url as URL in enumerator { files.append(url) }
        let baseComponents = unpacked.resolvingSymlinksInPath().pathComponents
        var decrypted = 0
        for (index, url) in files.enumerated() {
            if cancelled() { throw CancellationError() }
            let components = url.resolvingSymlinksInPath().pathComponents
            guard components.count > baseComponents.count,
                  Array(components.prefix(baseComponents.count)) == baseComponents else {
                throw fail("Unexpected archive enumeration path.")
            }
            let relative = components.dropFirst(baseComponents.count).joined(separator: "/")
            let destination = packed.appendingPathComponent(relative)
            let values = try url.resourceValues(forKeys: [.isDirectoryKey, .isSymbolicLinkKey])
            guard values.isSymbolicLink != true else { throw fail("IPA contains a symbolic link.") }
            if values.isDirectory == true {
                try fm.createDirectory(at: destination, withIntermediateDirectories: true)
                continue
            }
            try fm.createDirectory(at: destination.deletingLastPathComponent(), withIntermediateDirectories: true)
            let appComponents = app.resolvingSymlinksInPath().pathComponents
            if components.count > appComponents.count,
               Array(components.prefix(appComponents.count)) == appComponents, let primarySinf {
                let handle = try FileHandle(forReadingFrom: url)
                let head = try handle.read(upToCount: 4)
                try handle.close()
                if head == Data([0xcf, 0xfa, 0xed, 0xfe]) {
                    let sinf = url.deletingLastPathComponent().appendingPathComponent("SC_Info/" + url.lastPathComponent + ".sinf")
                    if !fm.fileExists(atPath: sinf.path) {
                        try fm.createDirectory(at: sinf.deletingLastPathComponent(), withIntermediateDirectories: true)
                        try fm.copyItem(at: primarySinf, to: sinf)
                        let sinfRelative = sinf.resolvingSymlinksInPath().pathComponents.dropFirst(baseComponents.count).joined(separator: "/")
                        let outputSinf = packed.appendingPathComponent(sinfRelative)
                        try fm.createDirectory(at: outputSinf.deletingLastPathComponent(), withIntermediateDirectories: true)
                        try fm.copyItem(at: sinf, to: outputSinf)
                    }
                }
            }
            var encrypted = ObjCBool(false)
            // This shares the helper's bounded Mach-O parser and complete page
            // checks. Unsupported FAT/32-bit images fail instead of partial export.
            try crumb_group_decrypt(url, destination, &encrypted)
            if encrypted.boolValue { decrypted += 1 }
            progress(0.1 + 0.75 * Double(index + 1) / Double(max(files.count, 1)))
        }
        guard decrypted > 0 else { throw fail("No encrypted images found; no decrypted result was published.") }
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
