import Foundation
import MapleSAP

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
                CrumbInspectIPA(UnsafeMutablePointer(mutating: url), UnsafeMutablePointer(mutating: bundle), &output, &length)
            }
        }
        defer { CrumbIPAFree(output, length) }
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
                    CrumbPrepareIPA(UnsafeMutablePointer(mutating: source), UnsafeMutablePointer(mutating: destination),
                        UnsafeMutablePointer(mutating: bytes.bindMemory(to: UInt8.self).baseAddress), input.count, &output, &length)
                }
            }
        }
        defer { CrumbIPAFree(output, length) }
        guard status == 0 else { throw StoreError.native(Int32(status)) }
        guard let output = output, length > 0, length < 4096 else { throw StoreError.packageInvalid }
        return try JSONDecoder().decode(Info.self, from: Data(bytes: output, count: length))
    }
}
