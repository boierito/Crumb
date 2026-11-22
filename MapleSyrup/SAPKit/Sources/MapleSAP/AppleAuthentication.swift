// Protocol, retry and redirect behavior adapted from majd/ipatool (MIT),
// pkg/appstore/appstore_login.go; Bag normalization follows f9aa653 (Oct 5).
// iPadn request semantics: plist content type, explicit 2FA attempt, phone Store routing.
// No CLI or password persistence. Crumb's production SAP signing now uses the iPadN WebKit/WASM bridge.
import Foundation
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

public enum AuthenticationError: Error, LocalizedError, Equatable {
    case invalidCode, invalidCredentials, twoFactorRequired, verificationRejected, accountDisabled
    case apple(failure: String, message: String)
    case invalidResponse(Int), http(Int), rateLimited, retryLater, network(Int)
    case invalidRedirect, redirectUnavailable(Int), tooManyRedirects, invalidSession

    public var errorDescription: String? {
        switch self {
        case .invalidCode: return "The verification code must contain exactly six digits."
        case .invalidCredentials: return "Apple rejected the Apple ID or password."
        case .twoFactorRequired: return "Apple requires a verification code from your trusted device."
        case .verificationRejected: return "Apple did not complete verification. Enter a fresh 2FA code."
        case .accountDisabled: return "Apple reports that this account is disabled or locked."
        case .apple(let failure, let message): return message.isEmpty ? "Apple rejected sign-in (code \(failure))." : message
        case .invalidResponse(let status): return "Apple returned no usable login response (HTTP \(status)). Try again later or from another network."
        case .http(let status): return "Apple authentication failed after bounded retries (HTTP \(status))."
        case .rateLimited: return "Apple rate limited sign-in. Wait before trying again."
        case .retryLater: return "Apple requested a wait longer than 30 seconds. Try again later."
        case .network(let code): return "Authentication network request failed (code \(code)). Check your connection."
        case .redirectUnavailable(let status): return "Apple returned an incomplete sign-in redirect (HTTP \(status)). Sign-in was not confirmed. Try again with the same verification code if it is still valid."
        case .invalidRedirect: return "Apple authentication endpoint or redirect was rejected. No credentials were forwarded."
        case .tooManyRedirects: return "Apple returned too many authentication redirects."
        case .invalidSession: return "The saved Store session is invalid. Sign in again."
        }
    }
}

public enum AuthenticationOutcome { case authenticated(StoreAccount), twoFactorRequired([StoreCookie]) }
public enum AuthenticationStage: String, Sendable {
    case bag = "Resolving Store Bag", sap = "Initializing SAP", signing = "Signing authentication request"
    case authenticating = "Contacting Apple", redirect = "Following Store pod", retrying = "Waiting to retry"
    case prepared = "Using prepared SAP session"
    case twoFactor = "Enter the code from your trusted device", saving = "Saving session in Keychain"
}

public protocol ActionSigning {
    func actionSignature(body: Data) async throws -> String
}
extension SAPSession: ActionSigning {}

// iPadN native Store login: device proof is calculated on-device, persisted in
// the Keychain and never packaged into a downloaded IPA.
public protocol StoreProofProvider {
    func currentHeaders(guid: String) async throws -> [String: String]?
    func provision(challenge: String, guid: String) async throws -> [String: String]
    func reset(guid: String) async
}

public struct TwoFactorAuthentication {
    public static func normalize(_ input: String) throws -> String {
        if input.isEmpty { return "" }
        let code = input.filter { !$0.isWhitespace }
        guard code.utf8.count == 6, code.allSatisfy({ $0.isASCII && $0 >= "0" && $0 <= "9" }) else {
            throw AuthenticationError.invalidCode
        }
        return code
    }
}

public struct AppleAuthentication {
    private let transport: AuthenticationTransport
    private let signer: ActionSigning
    private let persistence: StoreAccountPersistence
    private let proof: StoreProofProvider?
    private let sleep: (TimeInterval) async throws -> Void
    private let diagnostic: (String) async -> Void
    private let recoveryAttempts: Int
    private let recoveryWindow: TimeInterval?
    private let now: () -> Date
    public init(transport: AuthenticationTransport, signer: ActionSigning,
                persistence: StoreAccountPersistence,
                proof: StoreProofProvider? = nil,
                diagnostic: @escaping (String) async -> Void = { _ in },
                automaticRecovery: Bool = false,
                now: @escaping () -> Date = Date.init,
                sleep: @escaping (TimeInterval) async throws -> Void = {
                    try await Task.sleep(nanoseconds: UInt64($0 * 1_000_000_000))
                }) {
        self.transport = transport; self.signer = signer; self.persistence = persistence; self.proof = proof; self.sleep = sleep; self.diagnostic = diagnostic; self.now = now
        recoveryAttempts = automaticRecovery ? 12 : 3
        recoveryWindow = automaticRecovery ? 120 : nil
    }

    public func login(email: String, password: String, code: String = "", identity: MachineIdentity,
                      endpoint: URL, resolvedEndpoint: (URL) async -> Void = { _ in }, progress: (AuthenticationStage) async -> Void = { _ in }) async throws -> AuthenticationOutcome {
        let code = try TwoFactorAuthentication.normalize(code)
        var endpoint = try AuthenticationEndpoint.validate(endpoint)
        var recovery = AuthenticationRecoveryBudget(attempts: recoveryAttempts, window: recoveryWindow, now: now())
        let logicalAttempt = code.isEmpty ? 1 : 2
        var redirects = 0
        let body = try payload(email: email, password: password, code: code, guid: identity.guid, attempt: logicalAttempt)
        while true {
            try Task.checkCancellation()
            // Only publish endpoints after the same strict validation used for
            // credential replay. A warm 2FA/retry can skip an already resolved pod.
            await resolvedEndpoint(endpoint)
            let (data, response) = try await send(body: body, endpoint: endpoint, guid: identity.guid, progress: progress, secrets: [password, code, email], storefront: phoneNumberStorefront(email), recovery: &recovery)
            if (300..<400).contains(response.statusCode) {
                await diagnostic("authentication-redirect=received; HTTP=\(response.statusCode); location-present=\(response.value(forHTTPHeaderField: "Location") != nil)")
                guard [301, 302, 307, 308].contains(response.statusCode),
                      let location = response.value(forHTTPHeaderField: "Location") else {
                    throw AuthenticationError.invalidRedirect
                }
                guard redirects < 4 else { throw AuthenticationError.tooManyRedirects }
                endpoint = try AuthenticationEndpoint.redirect(from: endpoint, location: location)
                redirects += 1
                await progress(.redirect)
                // Same attempt and exact signed payload; never convert to GET.
                continue
            }
            guard let result = try? ApplePlist.dictionary(data) else { throw AuthenticationError.invalidResponse(response.statusCode) }
            let failure = string(result["failureType"])
            let message = string(result["customerMessage"])
            if failure.isEmpty, message == "MZFinance.BadLogin.Configurator_message" {
                if code.isEmpty { await progress(.twoFactor); return .twoFactorRequired(await transport.cookies()) }
                throw AuthenticationError.verificationRejected
            }
            if message == "Your account is disabled." { throw AuthenticationError.accountDisabled }
            if !failure.isEmpty {
                if failure == "-5000", message.isEmpty { throw AuthenticationError.invalidCredentials }
                throw AuthenticationError.apple(failure: safe(failure, secrets: [password, code, email]),
                    message: safe(message, secrets: [password, code, email, string(result["passwordToken"])]))
            }
            let dsid = sessionDSID(result, response: response)
            let token = string(result["passwordToken"])
            let storefront = response.value(forHTTPHeaderField: "X-Set-Apple-Store-Front") ?? ""
            if !message.isEmpty, token.isEmpty || dsid.isEmpty {
                throw AuthenticationError.apple(failure: "", message: safe(message, secrets: [password, code, email]))
            }
            guard response.statusCode == 200, !dsid.isEmpty, !token.isEmpty, !storefront.isEmpty else {
                throw AuthenticationError.invalidResponse(response.statusCode)
            }
            let info = result["accountInfo"] as? [String: Any] ?? [:]
            let address = info["address"] as? [String: Any] ?? [:]
            let name = [string(address["firstName"]), string(address["lastName"])].filter { !$0.isEmpty }.joined(separator: " ")
            let account = StoreAccount(email: string(info["appleId"]).isEmpty ? email : string(info["appleId"]),
                name: name, dsid: dsid, passwordToken: token, storefront: storefront,
                pod: response.value(forHTTPHeaderField: "pod"), guid: identity.guid,
                authenticationURL: endpoint, cookies: await transport.cookies())
            try account.validate(identity: identity)
            try Task.checkCancellation()
            await progress(.saving)
            try Task.checkCancellation()
            try persistence.save(account)
            return .authenticated(account)
        }
    }

    // iPadn-compatible account routing. Keep the Apple ID itself unchanged.
    func phoneNumberStorefront(_ appleID: String) -> String? {
        var number = appleID.trimmingCharacters(in: .whitespacesAndNewlines)
            .replacingOccurrences(of: " ", with: "").replacingOccurrences(of: "-", with: "")
        var international = false
        if number.hasPrefix("+") { international = true; number.removeFirst() }
        else if number.hasPrefix("00") { international = true; number.removeFirst(2) }
        guard !number.isEmpty, number.allSatisfy({ $0.isASCII && $0.isNumber }) else { return nil }
        if number.count == 13 && number.hasPrefix("861") { return "143465" }
        if number.count == 12 && number.hasPrefix("91") { number.removeFirst(2) }
        else if international { return nil }
        if number.count == 11 && number.hasPrefix("1") { return "143465" }
        if !international && number.count == 11 && number.hasPrefix("0") { number.removeFirst() }
        if number.count == 10, let first = number.first, ("6"..."9").contains(first) { return "143467" }
        return nil
    }

    private func sessionDSID(_ result: [String: Any], response: HTTPURLResponse) -> String {
        let bodyID = string(result["dsPersonId"])
        if !bodyID.isEmpty { return bodyID }
        let headerID = response.value(forHTTPHeaderField: "X-Dsid")?
            .trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        if !headerID.isEmpty && headerID.allSatisfy({ $0.isASCII && $0.isNumber }) { return headerID }
        guard let url = response.url else { return "" }
        let headers = response.allHeaderFields.reduce(into: [String: String]()) {
            $0[String(describing: $1.key)] = String(describing: $1.value)
        }
        for cookie in HTTPCookie.cookies(withResponseHeaderFields: headers, for: url) {
            for prefix in ["mz_at0-", "mz_at_ssl-"] where cookie.name.hasPrefix(prefix) {
                let candidate = String(cookie.name.dropFirst(prefix.count))
                if !candidate.isEmpty && candidate.allSatisfy({ $0.isASCII && $0.isNumber }) { return candidate }
            }
        }
        return ""
    }

    private func payload(email: String, password: String, code: String, guid: String, attempt: Int) throws -> Data {
        try PropertyListSerialization.data(fromPropertyList: ["appleId": email, "attempt": String(attempt),
            "guid": guid, "password": password + code, "rmp": "0", "why": "signIn"], format: .xml, options: 0)
    }

    private func send(body: Data, endpoint: URL, guid: String, progress: (AuthenticationStage) async -> Void, secrets: [String], storefront: String?, recovery: inout AuthenticationRecoveryBudget) async throws -> (Data, HTTPURLResponse) {
        var refreshedProof: [String: String]? = nil
        var provisioningRounds = 0
        var resetRounds = 0
        for attempt in 1...recoveryAttempts {
            try Task.checkCancellation()
            try recovery.checkDeadline(now: now())
            await diagnostic("authentication-recovery-attempt=\(attempt)/\(recoveryAttempts)")
            await progress(.signing)
            var request = URLRequest(url: endpoint)
            request.timeoutInterval = try recovery.requestTimeout(now: now())
            request.httpMethod = "POST"; request.httpBody = body
            request.setValue("application/x-apple-plist", forHTTPHeaderField: "Content-Type")
            request.setValue(storefront, forHTTPHeaderField: "X-Apple-Store-Front")
            request.setValue(SAPProtocol.userAgent, forHTTPHeaderField: "User-Agent")
            request.setValue("close", forHTTPHeaderField: "Connection")
            // Native iPadN login requires one-time X-Apple-AMD device proof.
            // Only attach it to the exact trusted Apple Store authentication paths.
            if let proof, Self.needsStoreProof(endpoint) {
                let headers: [String: String]?
                if let prepared = refreshedProof { headers = prepared; refreshedProof = nil }
                else { headers = try await proof.currentHeaders(guid: guid) }
                if let headers,
                   let amd = headers["X-Apple-AMD"], !amd.isEmpty,
                   let machine = headers["X-Apple-AMD-M"], !machine.isEmpty {
                    request.setValue(amd, forHTTPHeaderField: "X-Apple-AMD")
                    request.setValue(machine, forHTTPHeaderField: "X-Apple-AMD-M")
                }
            }
            let signature = try await signer.actionSignature(body: body)
            guard !signature.isEmpty else { throw SAPError.emptySignature }
            request.setValue(signature, forHTTPHeaderField: "X-Apple-ActionSignature")
            try Task.checkCancellation()
            // Native signing can be slow: recompute the remaining network
            // timeout afterwards, not before the blocking interpreter call.
            request.timeoutInterval = try recovery.requestTimeout(now: now())
            await progress(.authenticating)
            do {
                await diagnostic(await transport.cookieDiagnostic(for: endpoint))
                let (data, response) = try await transport.send(request)
                await diagnostic(ResponseDiagnostic.response(data, status: response.statusCode, scope: "authentication", attempt: attempt, secrets: secrets))
                guard data.count <= SAPProtocol.maximumBodySize else { throw SAPError.oversizedResponse }
                let result = try? ApplePlist.dictionary(data)
                // Apple Store SP/RP requests must be processed before interpreting
                // an ordinary login response, and cannot replay credentials to
                // an untrusted host or beyond the bounded recovery window.
                if let proof, Self.needsStoreProof(endpoint),
                   let action = response.value(forHTTPHeaderField: "X-Apple-AMD-Action")?.uppercased().trimmingCharacters(in: .whitespacesAndNewlines) {
                    if action == "SP" {
                        guard provisioningRounds < 2, attempt < recoveryAttempts,
                              let challenge = response.value(forHTTPHeaderField: "X-Apple-AMD-Data"),
                              challenge.count <= 16_384, !challenge.isEmpty,
                              challenge.count % 4 == 0,
                              let decoded = Data(base64Encoded: challenge), !decoded.isEmpty,
                              decoded.base64EncodedString() == challenge else {
                            throw AuthenticationError.invalidResponse(response.statusCode)
                        }
                        refreshedProof = try await proof.provision(challenge: challenge, guid: guid)
                        provisioningRounds += 1
                        await progress(.retrying)
                        continue
                    }
                    if action == "RP" {
                        guard resetRounds < 1, attempt < recoveryAttempts else {
                            throw AuthenticationError.invalidResponse(response.statusCode)
                        }
                        await proof.reset(guid: guid)
                        refreshedProof = nil
                        resetRounds += 1
                        await progress(.retrying)
                        continue
                    }
                }
                let populated = result.map { !string($0["failureType"]).isEmpty || !string($0["customerMessage"]).isEmpty || !string($0["passwordToken"]).isEmpty } ?? false
                let status = response.statusCode
                // An allowed redirect status without a usable Location is not an
                // unsafe destination. Retry the current validated endpoint only;
                // never invent a redirect or broaden the credential allowlist.
                let missingLocation = [301, 302, 307, 308].contains(status)
                    && (response.value(forHTTPHeaderField: "Location")?
                        .trimmingCharacters(in: .whitespacesAndNewlines).isEmpty ?? true)
                if populated || (!missingLocation && (300..<400).contains(status)) || (status == 200 && result != nil) { return (data, response) }
                let failure: AuthenticationError = missingLocation ? .redirectUnavailable(status)
                    : (status == 429 ? .rateLimited : .http(status))
                guard missingLocation || [200, 204, 403, 404, 429].contains(status) || status / 100 == 5 else { throw AuthenticationError.invalidResponse(status) }
                guard attempt < recoveryAttempts else { throw failure }
                let delay = try retryDelay(response.value(forHTTPHeaderField: "Retry-After"), attempt: attempt)
                try recovery.reserveRetry(delay: delay, now: now(), failure: failure)
                await progress(.retrying)
                try await sleep(delay)
            } catch let error as URLError {
                if error.code == .cancelled || Task.isCancelled { throw CancellationError() }
                guard [.timedOut, .networkConnectionLost, .cannotConnectToHost].contains(error.code), attempt < recoveryAttempts else {
                    throw AuthenticationError.network(error.code.rawValue)
                }
                await progress(.retrying)
                let delay = fallbackDelay(attempt: attempt)
                try recovery.reserveRetry(delay: delay, now: now(), failure: .network(error.code.rawValue))
                try await sleep(delay)
            }
        }
        throw AuthenticationError.invalidResponse(0)
    }

    private static func needsStoreProof(_ url: URL) -> Bool {
        guard (try? AuthenticationEndpoint.validate(url)) != nil,
              let host = url.host?.lowercased() else { return false }
        let path = URLComponents(url: url, resolvingAgainstBaseURL: false)?.percentEncodedPath ?? ""
        if host == "auth.itunes.apple.com" {
            return ["/auth/v1/native", "/auth/v1/native/", "/auth/v1/native/fast", "/auth/v1/native/fast/"].contains(path)
        }
        return (host == "buy.itunes.apple.com" || host.hasSuffix("-buy.itunes.apple.com")) &&
            ["/WebObjects/MZFinance.woa/wa/authenticate", "/WebObjects/MZFinance.woa/wa/authenticate/"].contains(path)
    }

    private func retryDelay(_ header: String?, attempt: Int) throws -> TimeInterval {
        guard let header = header?.trimmingCharacters(in: .whitespacesAndNewlines) else {
            return fallbackDelay(attempt: attempt)
        }
        if !header.isEmpty, header.allSatisfy({ $0.isASCII && $0.isNumber }) {
            guard let seconds = UInt64(header), seconds <= 30 else { throw AuthenticationError.retryLater }
            return max(1, Double(seconds))
        }
        for format in ["EEE, dd MMM yyyy HH:mm:ss z", "EEEE, dd-MMM-yy HH:mm:ss z", "EEE MMM d HH:mm:ss yyyy"] {
            let formatter = DateFormatter(); formatter.locale = Locale(identifier: "en_US_POSIX")
            formatter.timeZone = TimeZone(secondsFromGMT: 0); formatter.dateFormat = format
            if let date = formatter.date(from: header) {
                let value = max(0, date.timeIntervalSince(now()))
                guard value <= 30 else { throw AuthenticationError.retryLater }
                return max(1, value)
            }
        }
        return fallbackDelay(attempt: attempt)
    }
    private func fallbackDelay(attempt: Int) -> TimeInterval {
        recoveryWindow == nil ? Double(10 << (attempt - 1)) : min(15, Double(2 << min(attempt - 1, 3)))
    }
    private func string(_ value: Any?) -> String { (value as? String) ?? (value as? NSNumber)?.stringValue ?? "" }
    private func safe(_ text: String, secrets: [String]) -> String {
        var text = text
        for secret in secrets.filter({ !$0.isEmpty }).sorted(by: { $0.count > $1.count }) {
            text = text.replacingOccurrences(of: secret, with: "[redacted]")
        }
        return String(String.UnicodeScalarView(text.unicodeScalars.filter {
            !CharacterSet.controlCharacters.contains($0)
        }.prefix(500)))
    }
}

// Each login/verification owns one budget. Successful routing or Apple's logical
// attempt=2 must not grant another eleven automatic retries. Keep the reference
// three-attempt policy per request when automatic recovery is disabled.
private struct AuthenticationRecoveryBudget {
    let deadline: Date?
    private var remainingRetries: Int?
    init(attempts: Int, window: TimeInterval?, now: Date) {
        deadline = window.map { now.addingTimeInterval($0) }
        remainingRetries = window == nil ? nil : attempts - 1
    }
    func checkDeadline(now: Date) throws {
        if let deadline, now >= deadline { throw AuthenticationError.retryLater }
    }
    func requestTimeout(now: Date) throws -> TimeInterval {
        try checkDeadline(now: now)
        return deadline.map { min(30, $0.timeIntervalSince(now)) } ?? 30
    }
    mutating func reserveRetry(delay: TimeInterval, now: Date, failure: AuthenticationError) throws {
        if let remainingRetries, remainingRetries == 0 { throw failure }
        if let deadline, now.addingTimeInterval(delay) >= deadline { throw failure }
        if let remainingRetries { self.remainingRetries = remainingRetries - 1 }
    }
}
