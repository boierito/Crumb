import Foundation

/// A single successful ent/download exchange. Never persisted by StoreSession.
public struct FairPlayTransaction: Sendable {
    public let kbsync: Data
    public let keybag: Data?
    public let sinfs: [Data]
    public let diversitybag: Data?
    public init(kbsync: Data, keybag: Data?, sinfs: [Data], diversitybag: Data? = nil) {
        self.kbsync = kbsync; self.keybag = keybag; self.sinfs = sinfs; self.diversitybag = diversitybag
    }
    /// Keep distinct bag kinds separate. Multiple different values for one
    /// kind have no demonstrated single-producer contract and are rejected.
    public static func licenseBags(_ root: [String: Any]) throws -> (keybag: Data?, diversitybag: Data?) {
        let items = root["songList"] as? [[String: Any]] ?? []
        let containers = [root] + items.flatMap { [$0, $0["metadata"] as? [String: Any] ?? [:]] }
        func unique(_ keys: [String]) throws -> Data? {
            var result: Data?
            for container in containers {
                for key in keys {
                    guard let value = container[key] else { continue }
                    let bytes: Data
                    if let data = value as? Data { bytes = data }
                    else if let text = value as? String, let data = Data(base64Encoded: text) { bytes = data }
                    else { throw StoreError.invalidResponse }
                    guard !bytes.isEmpty, bytes.count <= 1 << 20,
                          result == nil || result == bytes else { throw StoreError.invalidResponse }
                    result = bytes
                }
            }
            return result
        }
        return try (unique(["keybag", "subKeyBag", "sbKeyBag"]), unique(["diversitybag"]))
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
