import Foundation
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

public enum FairPlayAuthorizationError: String, Error, Sendable {
    case accountMismatch, endpointUnavailable, endpointRejected, invalidSync
    case serverRejected, adiChallenge, incompleteLicense, invalidSignature, busy
    case httpRejected, invalidReply, appleRejected, redirectRejected, tooManyRedirects
}

public struct FairPlayAuthorizationReply: Sendable {
    public let keybag: Data?
    public let diversitybag: Data?
}

// Account/device provisioning, independent of app IDs, versions and Mach-O
// paths. One logical Bag-selected authorization, with bounded trusted routing
// redirects; no automatic retries of HTTP errors, challenges or rejection.
public actor FairPlayAuthorizationClient {
    private let account: StoreAccount
    private let identity: MachineIdentity
    private let transport: AuthenticationTransport
    private var busy = false
    private let diagnostic: @Sendable (String) async -> Void

    public init(account: StoreAccount, identity: MachineIdentity, transport: AuthenticationTransport,
                diagnostic: @escaping @Sendable (String) async -> Void = { _ in }) throws {
        try account.validate(identity: identity)
        self.account = account; self.identity = identity; self.transport = transport; self.diagnostic = diagnostic
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
        await observe(bagData, response: bagResponse, scope: "authorization-bag", attempt: 1)
        guard bagResponse.statusCode == 200, bagData.count <= SAPProtocol.maximumBodySize,
              let bag = try? ApplePlist.dictionary(bagData),
              let urls = bag["urlBag"] as? [String: Any], let endpoint = urls["authorizeMachine"] as? String else {
            throw FairPlayAuthorizationError.endpointUnavailable
        }
        let url = try StoreGUIDPolicy.authorizationURL(Self.validateEndpoint(endpoint), bag: bag, identity: identity)
        await diagnostic("authorization-guid-query=\((URLComponents(url: url, resolvingAgainstBaseURL: false)?.queryItems ?? []).contains(where: { $0.name == "guid" }) ? "present" : "not-required-by-bag")")
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
        // Fresh virtual ADI has no OTP/MD; the native unprovisioned path
        // omits these headers. A challenge is never silently bypassed/replayed.
        var redirects = 0
        var visited: Set<String> = [url.absoluteString]
        while true {
            try Task.checkCancellation()
            let (data, response) = try await transport.send(request)
            await observe(data, response: response, scope: "authorization", attempt: redirects + 1)
            if let action = response.value(forHTTPHeaderField: "x-apple-md-action"), !action.isEmpty {
                throw FairPlayAuthorizationError.adiChallenge
            }
            guard data.count <= SAPProtocol.maximumBodySize else { throw FairPlayAuthorizationError.invalidReply }
            let decoded = try? ApplePlist.dictionary(data)
            if let headers = decoded?["headers"] as? [String: Any], headers.keys.contains(where: { $0.lowercased() == "x-apple-md-action" }) {
                throw FairPlayAuthorizationError.adiChallenge
            }
            if [301, 302, 307, 308].contains(response.statusCode) {
                guard let location = response.value(forHTTPHeaderField: "Location"),
                      !location.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty,
                      let target = URL(string: location, relativeTo: request.url)?.absoluteURL else {
                    throw FairPlayAuthorizationError.redirectRejected
                }
                guard redirects < 3 else { throw FairPlayAuthorizationError.tooManyRedirects }
                let routed = try StoreGUIDPolicy.authorizationURL(Self.validateEndpoint(target.absoluteString), bag: bag, identity: identity)
                guard visited.insert(routed.absoluteString).inserted else { throw FairPlayAuthorizationError.tooManyRedirects }
                // Keep the original POST, exact signed body and signature.
                // Credentials only reach a validated Apple authorization URL.
                request.url = routed
                redirects += 1
                continue
            }
            guard response.statusCode == 200 else { throw FairPlayAuthorizationError.httpRejected }
            guard let root = decoded else { throw FairPlayAuthorizationError.invalidReply }
            if let failure = root["failureType"], String(describing: failure) != "0" {
                throw FairPlayAuthorizationError.appleRejected
            }
            // Native text+0x59e73..0x59ec7 only reads customerMessage in the
            // failureType branch. Message presence alone is not an error.
            func blob(_ key: String) -> Data? {
                if let bytes = root[key] as? Data, !bytes.isEmpty { return bytes }
                if let text = root[key] as? String, let bytes = Data(base64Encoded: text), !bytes.isEmpty { return bytes }
                return nil
            }
            let keybag = blob("keybag")
            let diversity = blob("diversitybag")
            // text+0x59bb2..0x59bb5 skips an absent keybag and proceeds to
            // diversitybag; text+0x59d1b..0x59d1e also permits no diversitybag.
            // Native post-import provisioning and SINF checks decide readiness.
            guard root["keybag"] == nil || keybag != nil,
                  root["diversitybag"] == nil || diversity != nil,
                  keybag != nil || diversity != nil,
                  (keybag?.count ?? 0) <= 1 << 20,
                  (diversity?.count ?? 0) <= 1 << 20 else {
                throw FairPlayAuthorizationError.incompleteLicense
            }
            return FairPlayAuthorizationReply(keybag: keybag, diversitybag: diversity)
        }
    }

    private func observe(_ data: Data, response: HTTPURLResponse, scope: String, attempt: Int) async {
        await diagnostic(ResponseDiagnostic.response(data, status: response.statusCode, scope: scope, attempt: attempt,
            secrets: [account.passwordToken, account.dsid, account.guid, account.email]))
    }

    public static func validateEndpoint(_ text: String) throws -> URL {
        guard text.utf8.count <= 8192, let url = try? SAPConfiguration.trustedAppleURL(text),
              let components = URLComponents(url: url, resolvingAgainstBaseURL: false),
              let host = components.host?.lowercased(),
              host == "buy.itunes.apple.com" || host.hasSuffix("-buy.itunes.apple.com"),
              ["/commerce/machine/authorize", "/commerce/machine/authorize/",
               "/WebObjects/MZFinance.woa/wa/authorizeMachine", "/WebObjects/MZFinance.woa/wa/authorizeMachine/"].contains(components.percentEncodedPath) else {
            throw FairPlayAuthorizationError.endpointRejected
        }
        return url
    }
}
