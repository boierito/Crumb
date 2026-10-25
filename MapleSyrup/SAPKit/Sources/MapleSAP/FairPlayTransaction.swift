import Foundation

/// A single successful ent/download exchange. Never persisted by StoreSession.
public struct FairPlayTransaction: Sendable {
    public let kbsync: Data
    public let keybag: Data?
    public let sinfs: [Data]
    public init(kbsync: Data, keybag: Data?, sinfs: [Data]) {
        self.kbsync = kbsync; self.keybag = keybag; self.sinfs = sinfs
    }
    public static func keybag(_ root: [String: Any]) -> Data? {
        let keys = ["keybag", "subKeyBag", "sbKeyBag"]
        let items = root["songList"] as? [[String: Any]] ?? []
        let containers = [root] + items.flatMap { [$0, $0["metadata"] as? [String: Any] ?? [:]] }
        for container in containers {
            for key in keys {
                if let data = container[key] as? Data, !data.isEmpty { return data }
                if let text = container[key] as? String, let data = Data(base64Encoded: text), !data.isEmpty { return data }
            }
        }
        return nil
    }
}
