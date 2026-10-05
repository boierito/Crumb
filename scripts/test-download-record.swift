import Foundation

@main struct DownloadRecordTests {
    static func main() throws {
        let fm = FileManager.default
        let root = fm.temporaryDirectory.appendingPathComponent(UUID().uuidString, isDirectory: true)
        try fm.createDirectory(at: root, withIntermediateDirectories: true)
        defer { try? fm.removeItem(at: root) }
        func record(_ name: String) -> DownloadRecord {
            DownloadRecord(filename: name, appID: "1", appName: "Fixture", bundleID: "test.fixture",
                version: "1.4.0", build: "14", externalVersionID: "123", date: Date())
        }
        for name in ["selected.ipa", "selected.json", "keep.ipa", "keep.json"] {
            try Data("fixture".utf8).write(to: root.appendingPathComponent(name))
        }
        try record("selected.ipa").deleteFiles(in: root)
        precondition(!fm.fileExists(atPath: root.appendingPathComponent("selected.ipa").path))
        precondition(!fm.fileExists(atPath: root.appendingPathComponent("selected.json").path))
        precondition(fm.fileExists(atPath: root.appendingPathComponent("keep.ipa").path))
        precondition(fm.fileExists(atPath: root.appendingPathComponent("keep.json").path))
        try record("selected.ipa").deleteFiles(in: root) // already removed through Files
        try Data().write(to: root.appendingPathComponent("stale.json"))
        try record("stale.ipa").deleteFiles(in: root)
        precondition(!fm.fileExists(atPath: root.appendingPathComponent("stale.json").path))
        for name in ["../keep.ipa", "sub/keep.ipa", "keep.json", ""] {
            do {
                try record(name).deleteFiles(in: root)
                preconditionFailure("Unsafe filename accepted")
            } catch let error as CocoaError { precondition(error.code == .fileReadInvalidFileName) }
        }
        print("Download deletion: IPA/sidecar, stale records, unrelated files and path boundaries passed")
    }
}
