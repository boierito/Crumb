import Foundation
import CryptoKit
import MapleSAP

nonisolated enum NativeFairPlayAuthorization {
    enum Failure: Error { case native(String), malformedReply }
    static func stateDirectory(account: StoreAccount) throws -> URL {
        let key = SHA256.hash(data: Data((account.dsid + ":" + account.guid).utf8)).map { String(format: "%02x", $0) }.joined()
        let documents = try FileManager.default.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: true)
        let directory = documents.appendingPathComponent("FairPlayAccounts", isDirectory: true).appendingPathComponent(key, isDirectory: true)
        try FairPlayNative.privateDirectory(directory)
        return directory
    }
    static func hasSavedState(account: StoreAccount) -> Bool {
        guard let directory = try? stateDirectory(account: account) else { return false }
        return FileManager.default.fileExists(atPath: directory.appendingPathComponent("account-state.json").path)
    }
    private static func wire(_ request: Data) throws -> Data {
        guard request.count <= 4 << 20 else { throw Failure.malformedReply }
        let text = String(decoding: request, as: UTF8.self)
        var output: UnsafeMutablePointer<UInt8>?
        var count = 0
        let status = text.withCString { WaffleFairPlayAuthorization(UnsafeMutablePointer(mutating: $0), &output, &count) }
        defer { WaffleSAPFree(output, count) }
        guard status == 0, let output, count > 0, count <= 4 << 20 else { throw Failure.malformedReply }
        return Data(bytes: output, count: count)
    }
    private static func call(_ fields: [String: Any]) async throws -> Data {
        let request = try JSONSerialization.data(withJSONObject: fields)
        let reply = try await Task.detached { try wire(request) }.value
        guard let root = try JSONSerialization.jsonObject(with: reply) as? [String: Any] else { throw Failure.malformedReply }
        guard root["ok"] as? Bool == true else {
            let code = root["code"] as? String ?? "authorization_failed"
            guard code.hasPrefix("authorization_"), code.count <= 80,
                  code.allSatisfy({ $0.isASCII && ($0.isLetter || $0 == "_") }) else { throw Failure.malformedReply }
            throw Failure.native(code)
        }
        return reply
    }
    private static func close(_ handle: UInt64) async {
        guard handle != 0 else { return }
        _ = try? await call(["Operation": "close", "Handle": handle])
    }

    // PRIVATE bootstrap material is decoded here, never sent to diagnostic logs.
    // Only finish returns a public, key-free result. Native state is per account;
    // each target IPA still has to pass every SINF and executable gate.
    static func run(account: StoreAccount, source: URL?, ipa: URL, destination: URL) async throws -> Data {
        let identity = try KeychainMachineIdentity.loadOrCreate()
        try account.validate(identity: identity)
        let directory = try stateDirectory(account: account)
        let cache = try FileManager.default.url(for: .cachesDirectory, in: .userDomainMask, appropriateFor: nil, create: true).appendingPathComponent("MapleSAP", isDirectory: true)
        try FairPlayNative.privateDirectory(cache)
        var handle: UInt64 = 0
        do {
            let raw = try await call(["Operation":"open", "Source":source?.path ?? "", "IPA":ipa.path,
                "DSID":account.dsid, "GUID":account.guid, "Cache":cache.path,
                "Directory":directory.path, "Destination":destination.path])
            guard let envelope = try JSONSerialization.jsonObject(with: raw) as? [String: Any],
                  let result = envelope["result"] as? [String: Any],
                  let number = result["handle"] as? NSNumber,
                  let needs = result["needs_diversity"] as? Bool else { throw Failure.malformedReply }
            handle = number.uint64Value
            guard handle != 0 else { throw Failure.malformedReply }
            var finish: [String: Any] = ["Operation":"finish", "Handle":handle]
            if needs {
                guard let encoded = result["kbsync"] as? String, let sync = Data(base64Encoded: encoded), !sync.isEmpty else {
                    throw Failure.malformedReply
                }
                let sapTransport = AppleSAPTransport()
                let protocolClient = SAPProtocol(transport: sapTransport)
                let sap = try SAPSession(guest: NativeSAPGuest(), transport: sapTransport)
                let transport = AppleAuthenticationTransport(cookies: account.cookies)
                do {
                    let configuration = try await protocolClient.bag(identity: identity)
                    try await sap.initialize(configuration: configuration, identity: identity)
                    let client = try FairPlayAuthorizationClient(account: account, identity: identity, transport: transport)
                    let license = try await client.authorize(kbsync: sync, needsDiversity: true, machineName: "Crumb iPhone", sign: { body in
                        try await sap.actionSignature(body: body)
                    })
                    finish["Keybag"] = license.keybag.base64EncodedString()
                    finish["Diversitybag"] = license.diversitybag?.base64EncodedString()
                    await sap.close(); sapTransport.close(); transport.close()
                } catch {
                    await sap.close(); sapTransport.close(); transport.close()
                    throw error
                }
            }
            try Task.checkCancellation()
            let reply = try await call(finish)
            handle = 0 // Native finish consumes and closes the session exactly once.
            return reply
        } catch {
            await close(handle)
            throw error
        }
    }
}
