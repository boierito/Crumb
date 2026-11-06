import Foundation

struct AdditionalSigningProfile: Codable, Identifiable {
    let id: String
    let name: String
    let data: Data
}
struct SigningReview: Identifiable {
    let id: UUID
    let appName: String
    let originalBundleID: String
    let bundleID: String
    let bundleCount: Int
    let changedIDs: Int
    let capabilities: [String]
    let profileName: String
}
struct SignedInstallation: Identifiable {
    let id = UUID()
    let record: DownloadRecord
    let source: URL
    let workspace: URL
}
