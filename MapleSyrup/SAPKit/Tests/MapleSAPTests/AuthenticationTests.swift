import XCTest
@testable import MapleSAP
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

final class AuthenticationTests: XCTestCase {
    private let endpoint = URL(string: "https://buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate")!
    private let identity = try! MachineIdentity(hardwareID: Data([2, 1, 2, 3, 4, 5]))
    private let responseHeaders = ["X-Set-Apple-Store-Front": "143441-1,29", "pod": "42"]
    private func plist(_ object: [String: Any]) throws -> Data {
        try PropertyListSerialization.data(fromPropertyList: object, format: .xml, options: 0)
    }
    private func success() throws -> Data {
        try plist(["dsPersonId": "123456789", "passwordToken": "fixture-token",
                   "accountInfo": ["appleId": "fixture@example.test", "address": ["firstName": "Test", "lastName": "Account"]]])
    }
    private func login(_ transport: FixtureAuthenticationTransport, code: String = "",
                       store: FixtureAccountStore = FixtureAccountStore(), signer: FixtureSigner = FixtureSigner(),
                       sleeps: FixtureSleeps = FixtureSleeps()) async throws -> AuthenticationOutcome {
        try await AppleAuthentication(transport: transport, signer: signer, persistence: store,
            sleep: { await sleeps.record($0) }).login(email: "fixture@example.test", password: "p<&>secret",
                code: code, identity: identity, endpoint: endpoint)
    }

    func testSignedPlistHasModernFieldsAndStoresNoPasswordOrCode() async throws {
        let transport = FixtureAuthenticationTransport([.http(200, try success(), responseHeaders)])
        let signer = FixtureSigner(); let store = FixtureAccountStore()
        guard case .authenticated(let account) = try await login(transport, code: "123 456", store: store, signer: signer) else {
            return XCTFail("Not authenticated")
        }
        let requests = await transport.requests
        let body = try XCTUnwrap(requests.first?.httpBody)
        XCTAssertEqual(signer.bodies, [body])
        XCTAssertEqual(requests[0].value(forHTTPHeaderField: "X-Apple-ActionSignature"), "fixture-signature-not-valid")
        XCTAssertEqual(requests[0].value(forHTTPHeaderField: "Content-Type"), "application/x-www-form-urlencoded")
        let parsed = try PropertyListSerialization.propertyList(from: body, format: nil) as! [String: String]
        XCTAssertEqual(parsed, ["appleId": "fixture@example.test", "password": "p<&>secret123456", "guid": identity.guid,
                                "attempt": "1", "rmp": "0", "why": "signIn"])
        XCTAssertEqual(account.dsid, "123456789"); XCTAssertEqual(account.pod, "42")
        XCTAssertEqual(account.authenticationURL, endpoint)
        XCTAssertEqual(try store.load(), account)
        let serialized = try XCTUnwrap(store.data)
        let dictionary = try JSONSerialization.jsonObject(with: serialized) as! [String: Any]
        XCTAssertNil(dictionary["password"]); XCTAssertNil(dictionary["code"])
        XCTAssertFalse(String(data: serialized, encoding: .utf8)!.contains("p<&>secret"))
    }

    func testPodRedirectsPreservePOSTBodyAndAttempt() async throws {
        let destination = "https://p42-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate"
        let transport = FixtureAuthenticationTransport([.http(302, Data(), ["Location": destination]),
            .http(307, Data(), ["Location": "?route=next"]), .http(200, try success(), responseHeaders)])
        let signer = FixtureSigner()
        _ = try await login(transport, signer: signer)
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 3)
        XCTAssertEqual(requests[1].url?.host, "p42-buy.itunes.apple.com")
        XCTAssertEqual(requests[2].url?.query, "route=next")
        for request in requests { XCTAssertEqual(request.httpMethod, "POST"); XCTAssertEqual(request.httpBody, requests[0].httpBody) }
        XCTAssertEqual(signer.bodies, requests.map { $0.httpBody! })
    }

    func testBagLoginUsesCanonicalInitialURLAndPreservesPodRedirectExactly() async throws {
        let bag = try plist(["urlBag": ["authenticateAccount": endpoint.absoluteString + "?route=a%2Fb+c",
            "sign-sap-setup": "https://play.itunes.apple.com/setup",
            "sign-sap-setup-cert": "https://s.mzstatic.com/sap/cert", "sign-sap-version": "200"]])
        let initial = try SAPConfiguration.parse(bag: bag).authenticationURL
        let pod = URL(string: "https://p42-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate?Pod=42&route=b%2Fc+d")!
        let transport = FixtureAuthenticationTransport([.http(302, Data(), ["Location": pod.absoluteString]),
            .http(200, try success(), responseHeaders)])
        let signer = FixtureSigner(); let sleeps = FixtureSleeps()
        _ = try await AppleAuthentication(transport: transport, signer: signer, persistence: FixtureAccountStore(),
            automaticRecovery: true, sleep: { await sleeps.record($0) })
            .login(email: "fixture@example.test", password: "secret", identity: identity, endpoint: initial)
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 2)
        XCTAssertEqual(requests[0].url?.absoluteString, endpoint.absoluteString + "/?route=a%2Fb+c")
        XCTAssertEqual(requests[1].url, pod) // Never rewrite Apple's redirect URL.
        XCTAssertEqual(requests[0].httpBody, requests[1].httpBody)
        XCTAssertTrue(requests.allSatisfy { $0.httpMethod == "POST" })
        XCTAssertEqual(signer.bodies.count, 2)
        let delays = await sleeps.values; XCTAssertTrue(delays.isEmpty)
    }

    func testResolvedPodIsReusableForTwoFactorWithoutAnotherRedirect() async throws {
        let pod = URL(string: "https://p42-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate")!
        let challenge = try plist(["customerMessage": "MZFinance.BadLogin.Configurator_message"])
        let transport = FixtureAuthenticationTransport([.http(302, Data(), ["Location": pod.absoluteString]),
            .http(200, challenge, [:]), .http(200, try success(), responseHeaders)])
        let signer = FixtureSigner()
        let auth = AppleAuthentication(transport: transport, signer: signer, persistence: FixtureAccountStore())
        let route = FixtureRoute()
        guard case .twoFactorRequired = try await auth.login(email: "fixture@example.test", password: "secret",
            identity: identity, endpoint: endpoint, resolvedEndpoint: { await route.record($0) }) else {
            return XCTFail("No challenge")
        }
        let resolved = await route.last
        XCTAssertEqual(resolved, pod)
        _ = try await auth.login(email: "fixture@example.test", password: "secret", code: "123456",
            identity: identity, endpoint: XCTUnwrap(resolved))
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 3)
        XCTAssertEqual(requests[2].url, pod)
        XCTAssertEqual(signer.bodies.count, 3) // fresh signing, including 2FA
        let payload = try PropertyListSerialization.propertyList(from: requests[2].httpBody!, format: nil) as! [String: String]
        XCTAssertEqual(payload["password"], "secret123456")
    }

    func testPreparedSAPSessionHandlesChallengeAndVerificationWithoutAnotherHandshake() async throws {
        let bag = try plist(["urlBag": ["authenticateAccount": endpoint.absoluteString,
            "sign-sap-setup": "https://play.itunes.apple.com/setup",
            "sign-sap-setup-cert": "https://play.itunes.apple.com/cert", "sign-sap-version": "200"]])
        let sapTransport = AuthenticationSAPFixture([bag,
            try plist(["sign-sap-setup-cert": Data([1])]), try plist(["sign-sap-setup-buffer": Data([2])])])
        let configuration = try await SAPProtocol(transport: sapTransport).bag(identity: identity)
        let guest = AuthenticationGuestFixture()
        let signer = try SAPSession(guest: guest, transport: sapTransport)
        try await signer.initialize(configuration: configuration, identity: identity)
        let pod = URL(string: "https://p42-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate/")!
        let challenge = try plist(["customerMessage": "MZFinance.BadLogin.Configurator_message"])
        let transport = FixtureAuthenticationTransport([.http(302, Data(), ["Location": pod.absoluteString]),
            .http(200, challenge, [:]), .http(200, try success(), responseHeaders)])
        let store = FixtureAccountStore(); let sleeps = FixtureSleeps(); let route = FixtureRoute()
        let auth = AppleAuthentication(transport: transport, signer: signer, persistence: store,
            automaticRecovery: true, sleep: { await sleeps.record($0) })
        guard case .twoFactorRequired = try await auth.login(email: "fixture@example.test", password: "secret",
            identity: identity, endpoint: configuration.authenticationURL, resolvedEndpoint: { await route.record($0) }) else {
            return XCTFail("No challenge")
        }
        XCTAssertNil(try store.load())
        let resolved = await route.last
        guard case .authenticated = try await auth.login(email: "fixture@example.test", password: "secret", code: "123456",
            identity: identity, endpoint: XCTUnwrap(resolved)) else { return XCTFail("Not verified") }
        let setupRequests = await sapTransport.requests
        XCTAssertEqual(setupRequests.count, 3) // One Bag, one certificate and one setup POST.
        XCTAssertEqual(guest.initializations, 1)
        XCTAssertEqual(guest.exchanges, 2)
        XCTAssertEqual(guest.signedBodies.count, 3)
        XCTAssertEqual(guest.signedBodies[0], guest.signedBodies[1]) // routing only
        XCTAssertNotEqual(guest.signedBodies[1], guest.signedBodies[2]) // password + code
        let requests = await transport.requests; XCTAssertEqual(requests[2].url, pod)
        let delays = await sleeps.values; XCTAssertTrue(delays.isEmpty)
        await signer.close()
    }

    func testSigningTimeAndRedirectsShareOneRecoveryDeadline() async throws {
        let clock = FixtureAuthenticationClock()
        let signer = AdvancingAuthenticationSigner(clock: clock)
        let transport = FixtureAuthenticationTransport([.http(302, Data(), ["Location": endpoint.absoluteString])])
        let auth = AppleAuthentication(transport: transport, signer: signer, persistence: FixtureAccountStore(),
            automaticRecovery: true, now: { clock.date }, sleep: { _ in })
        do {
            _ = try await auth.login(email: "fixture@example.test", password: "secret", identity: identity, endpoint: endpoint)
            XCTFail("Deadline reset on redirect")
        } catch { XCTAssertEqual(error as? AuthenticationError, .retryLater) }
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 1) // second signature exceeds the shared deadline
        XCTAssertEqual(requests[0].timeoutInterval, 30)
        XCTAssertEqual(signer.calls, 2)
    }

    func testLastRequestDoesNotRoundRemainingDeadlineUpToOneSecond() async throws {
        let clock = FixtureAuthenticationClock()
        let transport = FixtureAuthenticationTransport([.http(200, try success(), responseHeaders)])
        _ = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: FixtureAccountStore(),
            automaticRecovery: true, now: { clock.date })
            .login(email: "fixture@example.test", password: "secret", identity: identity, endpoint: endpoint) { stage in
                if stage == .signing { clock.date.addTimeInterval(119.75) }
            }
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 1)
        XCTAssertEqual(requests[0].timeoutInterval, 0.25, accuracy: 0.001)
    }

    func testBareDocumentPairsAuthenticateOnlyWithCompleteSession() async throws {
        let xml = "<Document><Protocol><key>dsPersonId</key><string>123456789</string><key>passwordToken</key><string>fixture-token</string></Protocol></Document>"
        guard case .authenticated = try await login(FixtureAuthenticationTransport([.http(200, Data(xml.utf8), responseHeaders)])) else {
            return XCTFail("Document pairs not parsed")
        }
        let incomplete = Data("<Document><key>dsPersonId</key><string>123456789</string></Document>".utf8)
        let store = FixtureAccountStore()
        do { _ = try await login(FixtureAuthenticationTransport([.http(200, incomplete, responseHeaders)]), store: store); XCTFail("Incomplete session saved") }
        catch { XCTAssertEqual(error as? AuthenticationError, .invalidResponse(200)) }
        XCTAssertNil(try store.load())
        XCTAssertThrowsError(try ApplePlist.dictionary(Data("<html><key>dsPersonId</key><string>123456789</string></html>".utf8)))
    }

    func testCredentialRedirectsRejectUntrustedDestinationsAnd303() async throws {
        for location in ["https://evil.test/login", "http://p42-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate",
                         "https://user:pass@buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate",
                         "https://buy.itunes.apple.com:444/WebObjects/MZFinance.woa/wa/authenticate",
                         "https://buy.itunes.apple.com/unrelated", "https://buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/%61uthenticate",
                         "https://buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate#fragment"] {
            let transport = FixtureAuthenticationTransport([.http(302, Data(), ["Location": location])])
            do { _ = try await login(transport); XCTFail("Followed unsafe redirect") }
            catch { XCTAssertEqual(error as? AuthenticationError, .invalidRedirect) }
            let requests = await transport.requests; XCTAssertEqual(requests.count, 1)
        }
        let transport = FixtureAuthenticationTransport([.http(303, Data(), ["Location": endpoint.absoluteString])])
        do { _ = try await login(transport); XCTFail("Followed 303") }
        catch { XCTAssertEqual(error as? AuthenticationError, .invalidRedirect) }
    }

    func testRedirectLoopHasFourHopLimit() async throws {
        let transport = FixtureAuthenticationTransport(Array(repeating: .http(302, Data(), ["Location": endpoint.absoluteString]), count: 5))
        do { _ = try await login(transport); XCTFail("Unlimited redirects") }
        catch { XCTAssertEqual(error as? AuthenticationError, .tooManyRedirects) }
        let requests = await transport.requests; XCTAssertEqual(requests.count, 5)
    }

    func testInvalidCredentialsLogicalRetryChangesOnlyAttemptAndDoesNotSaveFailure() async throws {
        let failure = try plist(["failureType": "-5000", "customerMessage": "Incorrect password"])
        let transport = FixtureAuthenticationTransport([.http(200, failure, [:]), .http(200, failure, [:])])
        let store = FixtureAccountStore()
        do { _ = try await login(transport, store: store); XCTFail("Accepted invalid credentials") }
        catch { XCTAssertEqual(error as? AuthenticationError, .apple(failure: "-5000", message: "Incorrect password")) }
        let requests = await transport.requests
        let attempts = try requests.map { try PropertyListSerialization.propertyList(from: $0.httpBody!, format: nil) as! [String: String] }
        XCTAssertEqual(attempts.map { $0["attempt"]! }, ["1", "2"])
        XCTAssertEqual(attempts[0]["password"], attempts[1]["password"])
        XCTAssertNil(try store.load())
    }

    func testTwoFactorChallengeDoesNotPersistAndFreshCodeFailureIsExplicit() async throws {
        let challenge = try plist(["customerMessage": "MZFinance.BadLogin.Configurator_message"])
        let store = FixtureAccountStore()
        guard case .twoFactorRequired = try await login(FixtureAuthenticationTransport([.http(200, challenge, [:])]), store: store) else {
            return XCTFail("No challenge")
        }
        XCTAssertNil(try store.load())
        do { _ = try await login(FixtureAuthenticationTransport([.http(200, challenge, [:])]), code: "654321", store: store); XCTFail("Accepted rejected code") }
        catch { XCTAssertEqual(error as? AuthenticationError, .verificationRejected) }
        XCTAssertNil(try store.load())
    }

    func testChallengeCookiesAreCarriedInMemoryIntoEphemeralVerificationJar() async throws {
        let rawCookie = HTTPCookie(properties: [.name: "challenge", .value: "fixture-secret-cookie", .domain: ".itunes.apple.com", .path: "/", .secure: "TRUE"])!
        let challenge = try plist(["customerMessage": "MZFinance.BadLogin.Configurator_message"])
        let transport = FixtureAuthenticationTransport([.http(200, challenge, [:])], cookies: [StoreCookie(rawCookie)])
        let store = FixtureAccountStore()
        guard case .twoFactorRequired(let cookies) = try await login(transport, store: store) else { return XCTFail("No challenge") }
        XCTAssertNil(try store.load())
        let verificationTransport = AppleAuthenticationTransport(cookies: cookies)
        defer { verificationTransport.close() }
        let carried = await verificationTransport.cookies()
        XCTAssertTrue(carried.contains { $0.name == "challenge" && $0.value == "fixture-secret-cookie" })
        let unrelated = AppleAuthenticationTransport()
        defer { unrelated.close() }
        let isolated = await unrelated.cookies()
        XCTAssertFalse(isolated.contains { $0.name == "challenge" })
    }

    func testCodeNormalizationRejectsNonASCIIAndInvalidLengthsBeforeRequests() async throws {
        XCTAssertEqual(try TwoFactorAuthentication.normalize("12\n34 56"), "123456")
        for code in ["12345", "1234567", "１２３４５６", "12a456", "123-456", "   "] {
            let transport = FixtureAuthenticationTransport([])
            do { _ = try await login(transport, code: code); XCTFail("Accepted invalid code") }
            catch { XCTAssertEqual(error as? AuthenticationError, .invalidCode) }
            let requests = await transport.requests; XCTAssertTrue(requests.isEmpty)
        }
    }

    func testEmptyHTTPResponsesRetryThreeTimesWithExactBody() async throws {
        for status in [204, 403, 404, 429, 500, 503] {
            let transport = FixtureAuthenticationTransport(Array(repeating: .http(status, Data("<html>temporary</html>".utf8), [:]), count: 3))
            let sleeps = FixtureSleeps()
            do { _ = try await login(transport, sleeps: sleeps); XCTFail("Accepted HTTP \(status)") }
            catch { XCTAssertEqual(error as? AuthenticationError, status == 429 ? .rateLimited : .http(status)) }
            let requests = await transport.requests
            XCTAssertEqual(requests.count, 3)
            XCTAssertTrue(requests.allSatisfy { $0.httpBody == requests[0].httpBody })
            let delays = await sleeps.values; XCTAssertEqual(delays, [10, 20])
        }
    }

    func testTwoFactorRecoversIncompleteRedirectAtValidatedPod() async throws {
        let pod = URL(string: "https://p42-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate")!
        let transport = FixtureAuthenticationTransport([
            .http(301, Data("<html>temporary</html>".utf8), [:]),
            .http(302, Data(), ["Location": "   "]),
            .http(200, try success(), responseHeaders)])
        let signer = FixtureSigner(); let store = FixtureAccountStore(); let sleeps = FixtureSleeps()
        let outcome = try await AppleAuthentication(transport: transport, signer: signer, persistence: store,
            automaticRecovery: true, sleep: { await sleeps.record($0) })
            .login(email: "fixture@example.test", password: "secret", code: "123456", identity: identity, endpoint: pod)
        guard case .authenticated = outcome else { return XCTFail("2FA recovery failed") }
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 3)
        XCTAssertTrue(requests.allSatisfy { $0.url == pod && $0.httpMethod == "POST" && $0.httpBody == requests[0].httpBody })
        XCTAssertEqual(signer.bodies.count, 3)
        let payload = try PropertyListSerialization.propertyList(from: XCTUnwrap(requests[0].httpBody), format: nil) as! [String: String]
        XCTAssertEqual(payload["password"], "secret123456")
        let delays = await sleeps.values; XCTAssertEqual(delays, [2, 4])
        XCTAssertNotNil(try store.load())
    }

    func testIncompleteRedirectRecoveryIsBoundedAndDoesNotSaveSession() async throws {
        let transport = FixtureAuthenticationTransport(Array(repeating: .http(301, Data(), [:]), count: 12))
        let store = FixtureAccountStore()
        do {
            _ = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: store,
                automaticRecovery: true, sleep: { _ in })
                .login(email: "fixture@example.test", password: "secret", code: "123456", identity: identity, endpoint: endpoint)
            XCTFail("Unbounded missing redirect recovery")
        } catch { XCTAssertEqual(error as? AuthenticationError, .redirectUnavailable(301)) }
        let requests = await transport.requests; XCTAssertEqual(requests.count, 12)
        XCTAssertNil(try store.load())
    }

    func testMissingRedirectStillUsesSharedRecoveryDeadline() async throws {
        let clock = FixtureAuthenticationClock()
        let transport = FixtureAuthenticationTransport([.http(302, Data(), [:])])
        do {
            _ = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: FixtureAccountStore(),
                automaticRecovery: true, now: { clock.date }, sleep: { _ in clock.date.addTimeInterval(121) })
                .login(email: "fixture@example.test", password: "secret", identity: identity, endpoint: endpoint)
            XCTFail("Missing Location reset deadline")
        } catch { XCTAssertEqual(error as? AuthenticationError, .retryLater) }
        let requests = await transport.requests; XCTAssertEqual(requests.count, 1)
    }

    func testIncompleteRedirectRecoveryDoesNotPermitUnsafeDestination() async throws {
        let transport = FixtureAuthenticationTransport([.http(301, Data(), [:]),
            .http(302, Data(), ["Location": "https://evil.test/login"])])
        do {
            _ = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: FixtureAccountStore(),
                automaticRecovery: true, sleep: { _ in })
                .login(email: "fixture@example.test", password: "secret", code: "123456", identity: identity, endpoint: endpoint)
            XCTFail("Followed unsafe destination")
        } catch { XCTAssertEqual(error as? AuthenticationError, .invalidRedirect) }
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 2); XCTAssertTrue(requests.allSatisfy { $0.url == endpoint })
    }

    func testFreshCodeChallengeDoesNotAppendRejectedCodeOrPersistPassword() async throws {
        let challenge = try plist(["customerMessage": "MZFinance.BadLogin.Configurator_message"])
        let transport = FixtureAuthenticationTransport([.http(200, challenge, [:]), .http(200, challenge, [:]),
            .http(200, challenge, [:]), .http(200, try success(), responseHeaders)])
        let store = FixtureAccountStore(); let signer = FixtureSigner()
        let auth = AppleAuthentication(transport: transport, signer: signer, persistence: store)
        guard case .twoFactorRequired = try await auth.login(email: "fixture@example.test", password: "secret",
            identity: identity, endpoint: endpoint) else { return XCTFail("Initial challenge absent") }
        do {
            _ = try await auth.login(email: "fixture@example.test", password: "secret", code: "111111", identity: identity, endpoint: endpoint)
            XCTFail("Accepted rejected code")
        } catch { XCTAssertEqual(error as? AuthenticationError, .verificationRejected) }
        XCTAssertNil(try store.load())
        guard case .twoFactorRequired = try await auth.login(email: "fixture@example.test", password: "secret", code: "",
            identity: identity, endpoint: endpoint) else { return XCTFail("New challenge absent") }
        guard case .authenticated = try await auth.login(email: "fixture@example.test", password: "secret", code: "222222",
            identity: identity, endpoint: endpoint) else { return XCTFail("Fresh verification failed") }
        let requests = await transport.requests
        let passwords = try requests.map { request in
            (try PropertyListSerialization.propertyList(from: XCTUnwrap(request.httpBody), format: nil) as! [String: String])["password"]
        }
        XCTAssertEqual(passwords, ["secret", "secret111111", "secret", "secret222222"])
        XCTAssertEqual(signer.bodies.count, 4)
        let saved = try XCTUnwrap(store.data)
        for secret in ["secret", "111111", "222222"] { XCTAssertFalse(String(decoding: saved, as: UTF8.self).contains(secret)) }
    }

    func testAutomaticRecoverySurvivesMoreThanThreeTemporaryReplies() async throws {
        let transport = FixtureAuthenticationTransport(Array(repeating: .http(404, Data("<html>temporary</html>".utf8), [:]), count: 5) + [.http(200, try success(), responseHeaders)])
        let sleeps = FixtureSleeps()
        let result = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: FixtureAccountStore(), automaticRecovery: true,
            sleep: { await sleeps.record($0) }).login(email: "fixture@example.test", password: "secret", identity: identity, endpoint: endpoint)
        guard case .authenticated = result else { return XCTFail("Recovery failed") }
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 6)
        XCTAssertTrue(requests.allSatisfy { $0.httpBody == requests.first?.httpBody })
        let delays = await sleeps.values
        XCTAssertEqual(delays, [2, 4, 8, 15, 15])
    }
    func testAutomaticRecoveryStopsAfterTwelveAttemptsAndHonorsRateLimit() async throws {
        let transport = FixtureAuthenticationTransport(Array(repeating: .http(503, Data(), [:]), count: 12))
        do {
            _ = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: FixtureAccountStore(), automaticRecovery: true, sleep: { _ in })
                .login(email: "fixture@example.test", password: "secret", identity: identity, endpoint: endpoint)
            XCTFail("Unlimited recovery")
        } catch { XCTAssertEqual(error as? AuthenticationError, .http(503)) }
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 12)
        let rate = FixtureAuthenticationTransport([.http(429, Data(), ["Retry-After": "300"])])
        do {
            _ = try await AppleAuthentication(transport: rate, signer: FixtureSigner(), persistence: FixtureAccountStore(), automaticRecovery: true, sleep: { _ in })
                .login(email: "fixture@example.test", password: "secret", identity: identity, endpoint: endpoint)
            XCTFail("Ignored Apple retry deadline")
        } catch { XCTAssertEqual(error as? AuthenticationError, .retryLater) }
        let rateRequests = await rate.requests
        XCTAssertEqual(rateRequests.count, 1)
    }

    func testAutomaticRetryAllowanceIsNotResetByPodRedirect() async throws {
        let pod = "https://p42-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate/"
        let temporary = FixtureAuthenticationTransport.Reply.http(503, Data(), [:])
        let transport = FixtureAuthenticationTransport(Array(repeating: temporary, count: 6)
            + [.http(302, Data(), ["Location": pod])] + Array(repeating: temporary, count: 12))
        let sleeps = FixtureSleeps(); let store = FixtureAccountStore()
        do {
            _ = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: store,
                automaticRecovery: true, sleep: { await sleeps.record($0) })
                .login(email: "fixture@example.test", password: "secret", identity: identity, endpoint: endpoint)
            XCTFail("Pod reset recovery allowance")
        } catch { XCTAssertEqual(error as? AuthenticationError, .http(503)) }
        let requests = await transport.requests
        XCTAssertEqual(requests.count, 13) // 12 unusable replies and one useful redirect.
        XCTAssertEqual(requests.last?.url?.absoluteString, pod)
        let delays = await sleeps.values; XCTAssertEqual(delays.count, 11)
        XCTAssertNil(try store.load())
    }

    func testNetworkAndHTTPRecoveryUseTheSameAllowance() async throws {
        let pod = "https://p42-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate/"
        let transport = FixtureAuthenticationTransport(Array(repeating: .error(.timedOut), count: 6)
            + [.http(302, Data(), ["Location": pod])] + Array(repeating: .http(404, Data(), [:]), count: 12))
        let sleeps = FixtureSleeps()
        do {
            _ = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: FixtureAccountStore(),
                automaticRecovery: true, sleep: { await sleeps.record($0) })
                .login(email: "fixture@example.test", password: "secret", identity: identity, endpoint: endpoint)
            XCTFail("Transport recovery had a separate allowance")
        } catch { XCTAssertEqual(error as? AuthenticationError, .http(404)) }
        let requests = await transport.requests; XCTAssertEqual(requests.count, 13)
        let delays = await sleeps.values; XCTAssertEqual(delays.count, 11)
    }

    func testFinalAllowedRecoveryCanStillSucceedAfterPodRedirect() async throws {
        let pod = "https://p42-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate/"
        let transport = FixtureAuthenticationTransport(Array(repeating: .http(204, Data(), [:]), count: 6)
            + [.http(302, Data(), ["Location": pod])] + Array(repeating: .http(404, Data(), [:]), count: 5)
            + [.http(200, try success(), responseHeaders)])
        let sleeps = FixtureSleeps(); let store = FixtureAccountStore()
        guard case .authenticated = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: store,
            automaticRecovery: true, sleep: { await sleeps.record($0) })
            .login(email: "fixture@example.test", password: "secret", code: "123456", identity: identity, endpoint: endpoint) else {
            return XCTFail("Last eligible response was discarded")
        }
        let requests = await transport.requests; XCTAssertEqual(requests.count, 13)
        let delays = await sleeps.values; XCTAssertEqual(delays.count, 11)
        XCTAssertNotNil(try store.load())
    }

    func testLogicalCredentialAttemptDoesNotResetAutomaticRetryAllowance() async throws {
        let failedCredential = try plist(["failureType": "-5000", "customerMessage": "Incorrect password"])
        let transport = FixtureAuthenticationTransport(Array(repeating: .http(204, Data(), [:]), count: 6)
            + [.http(200, failedCredential, [:])] + Array(repeating: .http(503, Data(), [:]), count: 12))
        do {
            _ = try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: FixtureAccountStore(),
                automaticRecovery: true, sleep: { _ in })
                .login(email: "fixture@example.test", password: "secret", code: "123456", identity: identity, endpoint: endpoint)
            XCTFail("Logical attempt reset recovery allowance")
        } catch { XCTAssertEqual(error as? AuthenticationError, .http(503)) }
        let requests = await transport.requests; XCTAssertEqual(requests.count, 13)
        let last = try ApplePlist.dictionary(XCTUnwrap(requests.last?.httpBody))
        XCTAssertEqual(last?["attempt"] as? String, "2")
        XCTAssertEqual(last?["password"] as? String, "secret123456")
    }

    func testAppleCredentialErrorOn403IsNotTreatedAsTransientHTML() async throws {
        let transport = FixtureAuthenticationTransport([.http(403, try plist(["failureType": "bad-password", "customerMessage": "Denied"]), [:])])
        do { _ = try await login(transport); XCTFail("Accepted error") }
        catch { XCTAssertEqual(error as? AuthenticationError, .apple(failure: "bad-password", message: "Denied")) }
        let requests = await transport.requests; XCTAssertEqual(requests.count, 1)
    }

    func testRetryAfterHonorsAppleDeadlineAndRejectsWaitOverBudget() async throws {
        let sleeps = FixtureSleeps()
        _ = try await login(FixtureAuthenticationTransport([.http(429, Data(), ["Retry-After": "2"]), .http(200, try success(), responseHeaders)]), sleeps: sleeps)
        let delays = await sleeps.values; XCTAssertEqual(delays, [2])
        for value in ["31", "99999999999999999999999999999999"] {
            let transport = FixtureAuthenticationTransport([.http(429, Data(), ["Retry-After": value])])
            do { _ = try await login(transport); XCTFail("Retried too early") }
            catch { XCTAssertEqual(error as? AuthenticationError, .retryLater) }
            let requests = await transport.requests; XCTAssertEqual(requests.count, 1)
        }
    }

    func testTransportTimeoutRetriesButCancellationDoesNot() async throws {
        let transport = FixtureAuthenticationTransport([.error(.timedOut), .http(200, try success(), responseHeaders)])
        let sleeps = FixtureSleeps(); _ = try await login(transport, sleeps: sleeps)
        let delays = await sleeps.values; XCTAssertEqual(delays, [10])
        let cancelled = FixtureAuthenticationTransport([.error(.cancelled)])
        do { _ = try await login(cancelled); XCTFail("Ignored cancellation") }
        catch { XCTAssertTrue(error is CancellationError) }
        let requests = await cancelled.requests; XCTAssertEqual(requests.count, 1)
    }

    func testRetryAfterPastHTTPDateWaitsAtLeastOneSecond() async throws {
        let sleeps = FixtureSleeps()
        _ = try await login(FixtureAuthenticationTransport([.http(429, Data(), ["Retry-After": "Wed, 21 Oct 2015 07:28:00 GMT"]),
            .http(200, try success(), responseHeaders)]), sleeps: sleeps)
        let delays = await sleeps.values; XCTAssertEqual(delays, [1])
    }

    func testMissingTokenDSIDOrStorefrontCannotCreateSession() async throws {
        for response in [try plist(["dsPersonId": "1"]), try plist(["passwordToken": "token"]), try plist([:])] {
            let store = FixtureAccountStore()
            do { _ = try await login(FixtureAuthenticationTransport([.http(200, response, responseHeaders)]), store: store); XCTFail("Saved incomplete response") }
            catch { XCTAssertEqual(error as? AuthenticationError, .invalidResponse(200)) }
            XCTAssertNil(try store.load())
        }
        do { _ = try await login(FixtureAuthenticationTransport([.http(200, try success(), [:])])); XCTFail("Saved without storefront") }
        catch { XCTAssertEqual(error as? AuthenticationError, .invalidResponse(200)) }
    }

    func testDisabledAccountAndUsefulMessagesAreSanitized() async throws {
        do { _ = try await login(FixtureAuthenticationTransport([.http(200, try plist(["customerMessage": "Your account is disabled."]), [:])])); XCTFail("Accepted disabled account") }
        catch { XCTAssertEqual(error as? AuthenticationError, .accountDisabled) }
        do {
            _ = try await login(FixtureAuthenticationTransport([.http(200,
                try plist(["failureType": "locked", "customerMessage": "Account p<&>secret fixture@example.test needs attention\n"]), [:])]))
            XCTFail("Accepted locked account")
        } catch {
            XCTAssertEqual(error as? AuthenticationError, .apple(failure: "locked", message: "Account [redacted] [redacted] needs attention"))
        }
    }

    func testPersistenceErrorDoesNotReturnAuthenticated() async throws {
        let store = FixtureAccountStore(); store.failSave = true
        do { _ = try await login(FixtureAuthenticationTransport([.http(200, try success(), responseHeaders)]), store: store); XCTFail("Ignored failed persistence") }
        catch { XCTAssertEqual(error as? SAPError, .keychain(-50)) }
    }

    func testCancelledBeforePersistenceDoesNotSaveReceivedCredentials() async throws {
        let store = FixtureAccountStore()
        let transport = FixtureAuthenticationTransport([.http(200, try success(), responseHeaders)])
        let task = Task {
            try await AppleAuthentication(transport: transport, signer: FixtureSigner(), persistence: store)
                .login(email: "fixture@example.test", password: "fixture-password", identity: identity, endpoint: endpoint) { stage in
                    if stage == .saving { withUnsafeCurrentTask { $0?.cancel() } }
                }
        }
        do { _ = try await task.value; XCTFail("Saved after cancellation") }
        catch { XCTAssertTrue(error is CancellationError) }
        XCTAssertNil(try store.load())
    }

    func testAccountRoundTripLogoutAndIdentityMismatch() throws {
        let store = FixtureAccountStore()
        let account = StoreAccount(email: "fixture@example.test", name: "Fixture", dsid: "123", passwordToken: "token",
            storefront: "143441-1,29", pod: nil, guid: identity.guid, authenticationURL: endpoint)
        try store.save(account); XCTAssertEqual(try store.load(), account)
        try account.validate(identity: identity)
        XCTAssertThrowsError(try account.validate(identity: MachineIdentity(hardwareID: Data(repeating: 4, count: 6))))
        try store.clear(); try store.clear(); XCTAssertNil(try store.load())
        XCTAssertEqual(identity.guid, "020102030405")
    }

    func testCookiesRoundTripExpiryAndHostRestriction() throws {
        let cookie = HTTPCookie(properties: [.name: "fixture", .value: "secret-cookie", .domain: ".itunes.apple.com",
            .path: "/", .secure: "TRUE", .expires: Date().addingTimeInterval(3600)])!
        let stored = StoreCookie(cookie)
        XCTAssertEqual(try JSONDecoder().decode(StoreCookie.self, from: JSONEncoder().encode(stored)), stored)
        XCTAssertEqual(stored.cookie()?.value, "secret-cookie")
        XCTAssertNil(stored.cookie(now: Date().addingTimeInterval(7200)))
        let malicious = HTTPCookie(properties: [.name: "fixture", .value: "secret", .domain: "evil.test", .path: "/"])!
        XCTAssertNil(StoreCookie(malicious).cookie())
    }

    #if canImport(Security)
    func testRealKeychainSaveUpdateLoadAndIdempotentLogout() throws {
        let store = KeychainStoreAccount(service: "com.nxtcoreee3.WaffleStore.tests.\(UUID().uuidString)")
        defer { try? store.clear() }
        XCTAssertNil(try store.load())
        let first = StoreAccount(email: "fixture@example.test", name: "Fixture", dsid: "123", passwordToken: "fixture-token",
            storefront: "143441-1,29", pod: nil, guid: identity.guid, authenticationURL: endpoint)
        let updated = StoreAccount(email: first.email, name: first.name, dsid: first.dsid, passwordToken: "updated-fixture-token",
            storefront: first.storefront, pod: "42", guid: first.guid, authenticationURL: endpoint)
        try store.save(first); XCTAssertEqual(try store.load(), first)
        try store.save(updated); XCTAssertEqual(try store.load(), updated)
        try store.clear(); try store.clear(); XCTAssertNil(try store.load())
    }
    #endif
}

private actor FixtureAuthenticationTransport: AuthenticationTransport {
    enum Reply { case http(Int, Data, [String: String]), error(URLError.Code) }
    var replies: [Reply]
    var requests: [URLRequest] = []
    let storedCookies: [StoreCookie]
    init(_ replies: [Reply], cookies: [StoreCookie] = []) { self.replies = replies; self.storedCookies = cookies }
    func send(_ request: URLRequest) async throws -> (Data, HTTPURLResponse) {
        requests.append(request)
        guard !replies.isEmpty else { throw SAPError.invalidState }
        switch replies.removeFirst() {
        case .http(let status, let data, let headers):
            return (data, HTTPURLResponse(url: request.url!, statusCode: status, httpVersion: nil, headerFields: headers)!)
        case .error(let code): throw URLError(code)
        }
    }
    func cookies() async -> [StoreCookie] { storedCookies }
}
private final class FixtureSigner: ActionSigning {
    var bodies: [Data] = []
    func actionSignature(body: Data) async throws -> String { bodies.append(body); return "fixture-signature-not-valid" }
}
private final class FixtureAccountStore: StoreAccountPersistence {
    var data: Data?
    var failSave = false
    func load() throws -> StoreAccount? { try data.map { try JSONDecoder().decode(StoreAccount.self, from: $0) } }
    func save(_ account: StoreAccount) throws { if failSave { throw SAPError.keychain(-50) }; data = try JSONEncoder().encode(account) }
    func clear() throws { data = nil }
}
private actor FixtureSleeps {
    var values: [TimeInterval] = []
    func record(_ seconds: TimeInterval) { values.append(seconds) }
}

private actor FixtureRoute {
    private(set) var last: URL?
    func record(_ url: URL) { last = url }
}
private final class FixtureAuthenticationClock {
    var date = Date(timeIntervalSince1970: 1_000_000)
}
private final class AdvancingAuthenticationSigner: ActionSigning {
    let clock: FixtureAuthenticationClock
    var calls = 0
    init(clock: FixtureAuthenticationClock) { self.clock = clock }
    func actionSignature(body: Data) async throws -> String {
        calls += 1
        clock.date.addTimeInterval(70)
        return "fixture-signature-not-valid"
    }
}

private actor AuthenticationSAPFixture: SAPTransport {
    var replies: [Data]
    private(set) var requests: [URLRequest] = []
    init(_ replies: [Data]) { self.replies = replies }
    func send(_ request: URLRequest) async throws -> (Data, HTTPURLResponse) {
        requests.append(request)
        guard !replies.isEmpty else { throw SAPError.invalidState }
        return (replies.removeFirst(), HTTPURLResponse(url: request.url!, statusCode: 200, httpVersion: nil, headerFields: nil)!)
    }
}
private final class AuthenticationGuestFixture: AppleSAPGuest {
    let executionMode = SAPExecutionMode.interpreted
    private(set) var initializations = 0
    private(set) var exchanges = 0
    private(set) var signedBodies: [Data] = []
    func initialize(hardwareID: Data) throws { initializations += 1 }
    func exchange(version: UInt32, hardwareID: Data, input: Data) throws -> (output: Data, state: Int32) {
        exchanges += 1
        return exchanges == 1 ? (Data([3]), 1) : (Data(), 0)
    }
    func sign(body: Data) throws -> Data { signedBodies.append(body); return Data([4]) }
    func close() {}
}
