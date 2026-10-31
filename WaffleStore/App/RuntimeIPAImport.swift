import Foundation

/// Copies an IPA picked by the user into private Crumb storage without
/// rebuilding or modifying the installed Crumb.app bundle.
/// This prepares input for the existing direct native FairPlay route; it does
/// not grant fairplayd additional access to the resulting data-container path.
nonisolated enum RuntimeIPAImport {
    enum ImportFailure: Error {
        case invalidInput
        case invalidArchive
    }

    nonisolated static func preflight(_ url: URL) throws {
        guard url.isFileURL, url.pathExtension.lowercased() == "ipa" else {
            throw ImportFailure.invalidInput
        }
        let values = try url.resourceValues(forKeys: [.isRegularFileKey, .isSymbolicLinkKey, .fileSizeKey])
        guard values.isRegularFile == true, values.isSymbolicLink != true,
              let bytes = values.fileSize, bytes >= 4 && Int64(bytes) <= 4_294_967_296 else {
            throw ImportFailure.invalidInput
        }
        let handle = try FileHandle(forReadingFrom: url)
        defer { try? handle.close() }
        guard try handle.read(upToCount: 4) == Data([0x50, 0x4B, 0x03, 0x04]) else {
            throw ImportFailure.invalidArchive
        }
    }

    /// The picker delivers an authorized copy. Keep access scoped only while
    /// copying; never rewrite or remove the user's selected original.
    nonisolated static func copyPickedIPA(_ selected: URL, into privateRoot: URL) throws -> URL {
        let scope = selected.startAccessingSecurityScopedResource()
        defer { if scope { selected.stopAccessingSecurityScopedResource() } }
        try preflight(selected)
        let fm = FileManager.default
        try fm.createDirectory(at: privateRoot, withIntermediateDirectories: true,
                               attributes: [.posixPermissions: 0o700])
        let destination = privateRoot.appendingPathComponent("Imported-" + UUID().uuidString + ".ipa")
        do {
            try fm.copyItem(at: selected, to: destination)
            try preflight(destination)
            try fm.setAttributes([.posixPermissions: 0o600], ofItemAtPath: destination.path)
            #if os(iOS)
            try fm.setAttributes([.protectionKey: FileProtectionType.complete], ofItemAtPath: destination.path)
            #endif
            var mutable = destination
            var values = URLResourceValues()
            values.isExcludedFromBackup = true
            try? mutable.setResourceValues(values)
            return destination
        } catch {
            try? fm.removeItem(at: destination)
            throw error
        }
    }
}
