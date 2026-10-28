import Foundation
import MapleSAP

nonisolated final class NativeKBSyncGenerator: KBSyncGenerator {
    func generate(identity: MachineIdentity, dsid: UInt64) throws -> Data {
        let cache = try FileManager.default.url(for: .cachesDirectory, in: .userDomainMask,
            appropriateFor: nil, create: true).appendingPathComponent("MapleSAP", isDirectory: true)
        try FileManager.default.createDirectory(at: cache, withIntermediateDirectories: true)
        var output: UnsafeMutablePointer<UInt8>?
        var length = 0
        let status = cache.path.withCString { path in
            identity.hardwareID.withUnsafeBytes { bytes in
                WaffleSAPKBSync(UnsafeMutablePointer(mutating: path),
                    UnsafeMutablePointer(mutating: bytes.bindMemory(to: UInt8.self).baseAddress),
                    identity.hardwareID.count, dsid, &output, &length)
            }
        }
        defer { WaffleSAPFree(output, length) }
        guard status == 0 else { throw StoreError.native(Int32(status)) }
        guard let output = output, length > 0, length <= 16 << 20 else { throw StoreError.native(8) }
        return Data(bytes: output, count: length)
    }
}
nonisolated enum NativePackage {
    struct Parameters: Encodable {
        let appID: String
        let bundleID: String
        let externalVersionID: String
        let md5: String
        let metadata: Data
        let sinfs: [Data]
    }
    struct Info: Codable, Sendable { let bundleID: String; let version: String; let build: String? }
    static func inspect(url: URL, bundle: String) throws -> Info {
        var output: UnsafeMutablePointer<UInt8>?
        var length = 0
        let status = url.absoluteString.withCString { url in
            bundle.withCString { bundle in
                WaffleInspectIPA(UnsafeMutablePointer(mutating: url), UnsafeMutablePointer(mutating: bundle), &output, &length)
            }
        }
        defer { WaffleSAPFree(output, length) }
        guard status == 0 else { throw StoreError.native(Int32(status)) }
        guard let output = output, length > 0, length < 4096 else { throw StoreError.packageInvalid }
        return try JSONDecoder().decode(Info.self, from: Data(bytes: output, count: length))
    }
    static func prepare(source: URL, destination: URL, app: StoreApp, descriptor: StoreDownload) throws -> Info {
        let input = try JSONEncoder().encode(Parameters(appID: app.id, bundleID: app.bundleID,
            externalVersionID: descriptor.externalVersionID, md5: descriptor.md5,
            metadata: descriptor.metadata, sinfs: descriptor.sinfs))
        var output: UnsafeMutablePointer<UInt8>?
        var length = 0
        let status = source.path.withCString { source in
            destination.path.withCString { destination in
                input.withUnsafeBytes { bytes in
                    WafflePrepareIPA(UnsafeMutablePointer(mutating: source), UnsafeMutablePointer(mutating: destination),
                        UnsafeMutablePointer(mutating: bytes.bindMemory(to: UInt8.self).baseAddress), input.count, &output, &length)
                }
            }
        }
        defer { WaffleSAPFree(output, length) }
        guard status == 0 else { throw StoreError.native(Int32(status)) }
        guard let output = output, length > 0, length < 4096 else { throw StoreError.packageInvalid }
        return try JSONDecoder().decode(Info.self, from: Data(bytes: output, count: length))
    }
}

// Laboratory retained flow. Restore account state before creating this download context.
nonisolated final class NativeRetainedKBSyncGenerator: KBSyncGenerator, @unchecked Sendable {
    private let accountStateDirectory: URL
    init(accountStateDirectory: URL) { self.accountStateDirectory = accountStateDirectory }
    private let lock = NSLock()
    private var handle: UInt64 = 0
    private var identity: MachineIdentity?
    private var dsid: UInt64 = 0
    func generate(identity: MachineIdentity, dsid: UInt64) throws -> Data {
        lock.lock(); defer { lock.unlock() }
        WaffleKBSyncClose(handle); handle = 0
        let cache = try FileManager.default.url(for: .cachesDirectory, in: .userDomainMask, appropriateFor: nil, create: true).appendingPathComponent("MapleSAP")
        try FileManager.default.createDirectory(at: cache, withIntermediateDirectories: true)
        var output: UnsafeMutablePointer<UInt8>?; var count = 0
        let status = cache.path.withCString { path in
            accountStateDirectory.path.withCString { directory in
                identity.hardwareID.withUnsafeBytes { bytes in
                    WaffleKBSyncOpenWithState(UnsafeMutablePointer(mutating: path), UnsafeMutablePointer(mutating: directory), UnsafeMutablePointer(mutating: bytes.bindMemory(to: UInt8.self).baseAddress), identity.hardwareID.count, dsid, &handle, &output, &count)
                }
            }
        }
        defer { WaffleSAPFree(output, count) }
        guard status == 0, let output, count > 0, count <= 16 << 20 else { throw StoreError.native(Int32(status)) }
        self.identity = identity; self.dsid = dsid
        return Data(bytes: output, count: count)
    }
    func bind(_ transaction: FairPlayTransaction) throws -> (context: Data, state: Data, summary: String) {
        lock.lock(); defer { lock.unlock() }
        guard let identity, let keybag = transaction.keybag else { throw StoreError.invalidResponse }
        let status = transaction.kbsync.withUnsafeBytes { blob in
            keybag.withUnsafeBytes { bag in
                WaffleKBSyncBind(handle, UnsafeMutablePointer(mutating: blob.bindMemory(to: UInt8.self).baseAddress), transaction.kbsync.count, UnsafeMutablePointer(mutating: bag.bindMemory(to: UInt8.self).baseAddress), keybag.count)
            }
        }
        guard status == 0 else { throw StoreError.native(Int32(status)) }
        var output: UnsafeMutablePointer<UInt8>?; var count = 0
        let consumeStatus = WaffleKBSyncConsume(handle, &output, &count)
        defer { WaffleSAPFree(output, count) }
        guard consumeStatus == 0, let output, count > 0, count <= 8 << 20 else { throw StoreError.native(Int32(consumeStatus)) }
        let state = Data(bytes: output, count: count)
        guard let report = try JSONSerialization.jsonObject(with: state) as? [String: Any],
              let nativeStatus = report["native_status"] as? Int,
              let completed = report["emulator_completed"] as? Bool,
              let imported = report["import_succeeded"] as? Bool else { throw StoreError.invalidResponse }
        // Go encodes nil []byte as JSON null. A mixed map (one empty lock
        // file and one real SIDB) cannot bridge to [String: String] in Swift.
        // Count entries independently, and decode only valid base64 strings.
        let files = report["state_files"] as? [String: Any] ?? [:]
        let bytes = files.values.reduce(0) { count, value in
            guard let encoded = value as? String, let decoded = Data(base64Encoded: encoded) else {
                return count
            }
            return count + decoded.count
        }
        let context = try JSONEncoder().encode(["schema":"crumb-fairplay-context-v2", "hardware_id":identity.hardwareID.base64EncodedString(), "dsid":String(dsid), "producer":"retained-storeagent", "keybag_consumer":"storeagent-10.9-import-keybag"])
        return (context, state, "Keybag consumer completed=\(completed); account state restored=\(report["account_state_restored"] as? Bool == true), saved=\(report["account_state_saved"] as? Bool == true), save failed=\(report["account_state_save_failed"] as? Bool == true); native status=\(nativeStatus); import succeeded=\(imported); virtual state files=\(files.count), bytes=\(bytes). A portable enrollment model has not been produced.")
    }
    func close() {
        lock.lock(); defer { lock.unlock() }
        WaffleKBSyncClose(handle); handle = 0; identity = nil; dsid = 0
    }
    deinit { close() }
}
