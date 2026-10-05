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
    let externalVersionID: String
    let date: Date
    var fileURL: URL? {
        guard filename == URL(fileURLWithPath: filename).lastPathComponent,
              filename.hasSuffix(".ipa"), !filename.contains(".."),
              let docs = try? FileManager.default.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: false) else { return nil }
        let url = docs.appendingPathComponent("Downloads").appendingPathComponent(filename)
        return FileManager.default.fileExists(atPath: url.path) ? url : nil
    }
    static func load() -> [DownloadRecord] {
        let fm = FileManager.default
        guard let docs = try? fm.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: false),
              let files = try? fm.contentsOfDirectory(at: docs.appendingPathComponent("Downloads"), includingPropertiesForKeys: nil) else { return [] }
        return files.filter { $0.pathExtension == "json" }.compactMap { file in
            guard let data = try? Data(contentsOf: file), data.count < 65536,
                  let record = try? JSONDecoder().decode(DownloadRecord.self, from: data), record.fileURL != nil else { return nil }
            return record
        }.sorted { $0.date > $1.date }
    }
}
