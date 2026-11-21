import Foundation
import ZIPFoundation

/// Privacy gate for user-visible IPA exports. Apple DRM SINF material can be
/// bound to the purchasing account, so encrypted originals must stay private.
/// A decrypted record alone is not proof: old helper outputs and imported files
/// can still carry the original metadata. Inspect every ZIP path fail-closed.
enum IPAExportPrivacy {
    static func isSafeToShare(_ url: URL) -> Bool {
        guard url.isFileURL, url.pathExtension.lowercased() == "ipa",
              let archive = try? Archive(url: url, accessMode: .read) else { return false }
        var hasApp = false
        for entry in archive {
            guard entry.type != .symlink else { return false }
            let components = entry.path.lowercased().split(separator: "/").map(String.init)
            guard !components.isEmpty else { return false }
            let leaf = components.last ?? ""
            if components.first == "payload", components.contains(where: { $0.hasSuffix(".app") }) {
                hasApp = true
            }
            if components.contains("sc_info") ||
               ["itunesmetadata.plist", "itunesartwork", "itunesartwork@2x"].contains(leaf) ||
               leaf.hasSuffix(".sinf") || leaf.hasSuffix(".supp") || leaf.hasSuffix(".supf") {
                return false
            }
        }
        return hasApp
    }
}


extension DownloadRecord {
    // A persisted decrypted flag is not proof of a safe export.
    // Verify the actual ZIP archive for legacy helper outputs as well.
    var exportableFileURL: URL? {
        guard isDecrypted, let url = fileURL,
              IPAExportPrivacy.isSafeToShare(url) else { return nil }
        return url
    }
}
