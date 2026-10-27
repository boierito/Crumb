import Foundation
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

public enum FairPlayAuthorizationError: String, Error, Sendable {
    case accountMismatch, endpointUnavailable, endpointRejected, invalidSync
    case serverRejected, adiChallenge, incompleteLicense, invalidSignature, busy
}

public struct FairPlayAuthorizationReply: Sendable {
    public let keybag: Data?
    public let diversitybag: Data?
}

// Account/device provisioning, independent of app IDs, versions and Mach-O
// paths. Uses one Bag-selected request; no retries of an authorization action.
public actor FairPlayAuthorizationClient {
    private let account: StoreAccount
    private let identity: MachineIdentity
    private let transport: AuthenticationTransport
    private var busy = false

    public init(account: StoreAccount, identity: MachineIdentity, transport: AuthenticationTransport) throws {
        try account.validate(identity: identity)
        self.account = account; self.identity = identity; self.transport = transport
    }

    public func authorize(kbsync: Data, needsDiversity: Bool, machineName: String,
                          sign: @Sendable (Data) async throws -> String) async throws -> FairPlayAuthorizationReply {
        guard !busy else { throw FairPlayAuthorizationError.busy }
        busy = true; defer { busy = false }
        guard !kbsync.isEmpty, kbsync.count <= 1 << 20 else { throw FairPlayAuthorizationError.invalidSync }
        try Task.checkCancellation()
        var bagRequest = URLRequest(url: URL(string: "https://init.itunes.apple.com/bag.xml?guid=\(identity.guid)")!)
        bagRequest.setValue(SAPProtocol.userAgent, forHTTPHeaderField: "User-Agent")
        let (bagData, bagResponse) = try await transport.send(bagRequest)
        guard bagResponse.statusCode == 200, bagData.count <= SAPProtocol.maximumBodySize,
              let bag = try? ApplePlist.dictionary(bagData),
              let urls = bag["urlBag"] as? [String: Any], let endpoint = urls["authorizeMachine"] as? String else {
            throw FairPlayAuthorizationError.endpointUnavailable
        }
        let url = try Self.validateEndpoint(endpoint)
        var fields: [String: Any] = ["kbsync": kbsync, "guid": identity.guid,
                                     "machineName": String(machineName.prefix(128))]
        if needsDiversity { fields["needDiv"] = "1" }
        // Original StoreAgent uses plist Data, not a base64 String field.
        let body = try PropertyListSerialization.data(fromPropertyList: fields, format: .xml, options: 0)
        let signature = try await sign(body)
        guard !signature.isEmpty, !signature.contains("\r"), !signature.contains("\n") else {
            throw FairPlayAuthorizationError.invalidSignature
        }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"; request.httpBody = body
        request.setValue(SAPProtocol.userAgent, forHTTPHeaderField: "User-Agent")
        request.setValue("application/x-apple-plist", forHTTPHeaderField: "Content-Type")
        request.setValue(account.passwordToken, forHTTPHeaderField: "X-Token")
        request.setValue(account.dsid, forHTTPHeaderField: "X-Dsid")
        request.setValue(account.dsid, forHTTPHeaderField: "iCloud-DSID")
        request.setValue(account.storefront, forHTTPHeaderField: "X-Apple-Store-Front")
        request.setValue(signature, forHTTPHeaderField: "X-Apple-ActionSignature")
        // Fresh virtual ADI has no OTP/MD. Original StoreAgent explicitly omits
        // these headers in that case. Never synthesize them or replay a challenge.
        try Task.checkCancellation()
        let (data, response) = try await transport.send(request)
        if let action = response.value(forHTTPHeaderField: "x-apple-md-action"), !action.isEmpty {
            throw FairPlayAuthorizationError.adiChallenge
        }
        guard response.statusCode == 200, data.count <= SAPProtocol.maximumBodySize,
              let root = try? ApplePlist.dictionary(data) else { throw FairPlayAuthorizationError.serverRejected }
        if let headers = root["headers"] as? [String: Any], headers.keys.contains(where: { $0.lowercased() == "x-apple-md-action" }) {
            throw FairPlayAuthorizationError.adiChallenge
        }
        if let failure = root["failureType"], String(describing: failure) != "0" {
            throw FairPlayAuthorizationError.serverRejected
        }
        guard root["customerMessage"] == nil else { throw FairPlayAuthorizationError.serverRejected }
        func blob(_ key: String) -> Data? {
            if let bytes = root[key] as? Data, !bytes.isEmpty { return bytes }
            if let text = root[key] as? String, let bytes = Data(base64Encoded: text), !bytes.isEmpty { return bytes }
            return nil
        }
        let keybag = blob("keybag")
        let diversity = blob("diversitybag")
        // text+0x59bb2..0x59bb5 skips an absent keybag and proceeds to
        // diversitybag. The already-restored download state is not replaced.
        guard root["keybag"] == nil || keybag != nil,
              root["diversitybag"] == nil || diversity != nil,
              keybag != nil || diversity != nil,
              !needsDiversity || diversity != nil,
              (keybag?.count ?? 0) <= 1 << 20,
              (diversity?.count ?? 0) <= 1 << 20 else {
            throw FairPlayAuthorizationError.incompleteLicense
        }
        return FairPlayAuthorizationReply(keybag: keybag, diversitybag: diversity)
    }

    public static func validateEndpoint(_ text: String) throws -> URL {
        guard let url = try? SAPConfiguration.trustedAppleURL(text),
              let components = URLComponents(url: url, resolvingAgainstBaseURL: false),
              let host = components.host?.lowercased(),
              host == "buy.itunes.apple.com" || host.hasSuffix("-buy.itunes.apple.com"),
              ["/commerce/machine/authorize", "/WebObjects/MZFinance.woa/wa/authorizeMachine"].contains(components.percentEncodedPath) else {
            throw FairPlayAuthorizationError.endpointRejected
        }
        return url
    }
}
