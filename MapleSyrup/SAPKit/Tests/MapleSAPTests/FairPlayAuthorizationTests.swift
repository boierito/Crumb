import XCTest
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif
@testable import MapleSAP

private actor LicenseTransport: AuthenticationTransport {
    var requests: [URLRequest] = []
    let body: Data
    let endpoint: String
    let status: Int
    let headers: [String: String]
    let policy: [String: Any]?
    init(reply: [String: Any], endpoint: String = "https://buy.itunes.apple.com/commerce/machine/authorize", status: Int = 200, headers: [String: String] = [:], policy: [String: Any]? = ["schemes": ["https"], "regex": ["/commerce/"]]) throws {
        self.body = try PropertyListSerialization.data(fromPropertyList: reply, format: .xml, options: 0)
        self.endpoint = endpoint; self.status = status; self.headers = headers; self.policy = policy
    }
    func cookies() async -> [StoreCookie] { [] }
    func send(_ request: URLRequest) async throws -> (Data, HTTPURLResponse) {
        requests.append(request)
        if requests.count == 1 {
            var fields: [String: Any] = ["urlBag": ["authorizeMachine": endpoint]]
            if let policy { fields["guid-urls"] = policy }
            let bag = try PropertyListSerialization.data(fromPropertyList: fields, format: .xml, options: 0)
            return (bag, HTTPURLResponse(url: request.url!, statusCode: 200, httpVersion: nil, headerFields: nil)!)
        }
        return (body, HTTPURLResponse(url: request.url!, statusCode: status, httpVersion: nil, headerFields: headers)!)
    }
}

final class FairPlayAuthorizationTests: XCTestCase {
    private func client(_ transport: AuthenticationTransport, diagnostic: @escaping @Sendable (String) async -> Void = { _ in }) throws -> FairPlayAuthorizationClient {
        let id = try MachineIdentity(hardwareID: Data([2,0,0,0,0,1]))
        let account = StoreAccount(email: "fixture@example.invalid", name: "Fixture", dsid: "1", passwordToken: "PRIVATE-TOKEN", storefront: "143441-1,29", pod: nil, guid: id.guid, authenticationURL: URL(string: "https://buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate")!)
        return try FairPlayAuthorizationClient(account: account, identity: id, transport: transport, diagnostic: diagnostic)
    }
    func testTypedSyncAndDistinctBagsWithoutAppSpecificFields() async throws {
        let transport = try LicenseTransport(reply: ["keybag": Data([1]), "diversitybag": Data([2])])
        let reply = try await client(transport).authorize(kbsync: Data([3,4]), needsDiversity: true, machineName: "Fixture", sign: { body in
            let fields = try PropertyListSerialization.propertyList(from: body, format: nil) as! [String: Any]
            XCTAssertEqual(fields["kbsync"] as? Data, Data([3,4])); return "SIGNED-FIXTURE"
        })
        XCTAssertEqual(reply.keybag, Data([1])); XCTAssertEqual(reply.diversitybag, Data([2]))
        let requests = await transport.requests; XCTAssertEqual(requests.count, 2)
        let request = requests[1]
        let fields = try PropertyListSerialization.propertyList(from: request.httpBody!, format: nil) as! [String: Any]
        XCTAssertEqual(Set(fields.keys), Set(["kbsync","needDiv","guid","machineName"]))
        XCTAssertEqual(fields["needDiv"] as? String, "1")
        XCTAssertEqual(URLComponents(url: request.url!, resolvingAgainstBaseURL: false)?.queryItems?.filter { $0.name == "guid" }.map(\.value), ["020000000001"])
        XCTAssertEqual(request.value(forHTTPHeaderField: "X-Token"), "PRIVATE-TOKEN")
        XCTAssertNil(request.value(forHTTPHeaderField: "X-Apple-MD"))
        XCTAssertNil(request.value(forHTTPHeaderField: "X-Apple-MD-M"))
    }
    func testDiversityOnlyReplyMatchesNativeOptionalKeybagBranch() async throws {
        let transport = try LicenseTransport(reply: ["diversitybag": Data([2])])
        let reply = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
        XCTAssertNil(reply.keybag)
        XCTAssertEqual(reply.diversitybag, Data([2]))
    }
    func testMalformedPresentBagCannotBeIgnored() async throws {
        let replies: [[String: Any]] = [["keybag": "invalid!", "diversitybag": Data([2])], ["keybag": Data(), "diversitybag": Data([2])], [:]]
        for reply in replies {
            let transport = try LicenseTransport(reply: reply)
            do { _ = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" }); XCTFail("malformed bag accepted") }
            catch { XCTAssertEqual((error as? FairPlayAuthorizationError)?.rawValue, "incompleteLicense") }
        }
    }
    func testKeybagOnlyReplyStillRequiresNativeProvisioningAndSINFGates() async throws {
        let transport = try LicenseTransport(reply: ["keybag": Data([1])])
        let reply = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
        XCTAssertEqual(reply.keybag, Data([1]))
        XCTAssertNil(reply.diversitybag)
    }
    func testNoDiversityFlagWhenNotRequested() async throws {
        let transport = try LicenseTransport(reply: ["keybag": Data([1])])
        _ = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: false, machineName: "Fixture", sign: { _ in "SIGN" })
        let requests = await transport.requests
        let fields = try PropertyListSerialization.propertyList(from: requests[1].httpBody!, format: nil) as! [String: Any]
        XCTAssertNil(fields["needDiv"])
    }
    func testFailureAndChallengeCannotYieldLicenseOrTriggerReplay() async throws {
        for (reply,status,headers) in [([:] as [String: Any],200,[:]), (["failureType":"123","customerMessage":"PRIVATE"],200,[:]), (["keybag":Data([1]),"diversitybag":Data([2])],302,[:]), (["keybag":Data([1]),"diversitybag":Data([2])],200,["x-apple-md-action":"SP"])] {
            let transport = try LicenseTransport(reply: reply, status: status, headers: headers)
            do { _ = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" }); XCTFail("invalid reply yielded a license") }
            catch { XCTAssertFalse(String(describing: error).contains("PRIVATE")) }
            let requests = await transport.requests; XCTAssertEqual(requests.count,2)
        }
    }
    func testEndpointCannotSendCredentialsToAnotherService() async throws {
        for endpoint in ["https://example.org/commerce/machine/authorize", "http://buy.itunes.apple.com/commerce/machine/authorize", "https://buy.itunes.apple.com/commerce/machine/%61uthorize", "https://buy.itunes.apple.com/other", "https://evil-buy.itunes.apple.com.evil.example/commerce/machine/authorize"] {
            let transport = try LicenseTransport(reply: [:], endpoint: endpoint)
            do { _ = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" }); XCTFail("untrusted endpoint accepted") } catch {}
            let requests = await transport.requests; XCTAssertEqual(requests.count,1)
        }
    }
}

private struct AuthorizationHTTPReply {
    let data: Data
    let status: Int
    let headers: [String: String]
    init(_ fields: [String: Any] = [:], status: Int = 200, headers: [String: String] = [:]) throws {
        data = try PropertyListSerialization.data(fromPropertyList: fields, format: .xml, options: 0)
        self.status = status; self.headers = headers
    }
    init(raw: Data, status: Int) { data = raw; self.status = status; headers = [:] }
}
private actor AuthorizationSequenceTransport: AuthenticationTransport {
    var requests: [URLRequest] = []
    let replies: [AuthorizationHTTPReply]
    init(_ replies: [AuthorizationHTTPReply]) { self.replies = replies }
    func cookies() async -> [StoreCookie] { [] }
    func send(_ request: URLRequest) async throws -> (Data, HTTPURLResponse) {
        requests.append(request)
        if requests.count == 1 {
            let bag = try PropertyListSerialization.data(fromPropertyList: [
                "urlBag": ["authorizeMachine": "https://buy.itunes.apple.com/commerce/machine/authorize"],
                "guid-urls": ["schemes": ["https"], "regex": ["/commerce/", "/WebObjects/"]]
            ], format: .xml, options: 0)
            return (bag, HTTPURLResponse(url: request.url!, statusCode: 200, httpVersion: nil, headerFields: nil)!)
        }
        guard requests.count - 2 < replies.count else { throw URLError(.badServerResponse) }
        let reply = replies[requests.count - 2]
        return (reply.data, HTTPURLResponse(url: request.url!, statusCode: reply.status, httpVersion: nil, headerFields: reply.headers)!)
    }
}
private actor AuthorizationTrace {
    var lines: [String] = []
    func append(_ line: String) { lines.append(line) }
}

extension FairPlayAuthorizationTests {
    func testBagGUIDRulesPreserveRoutingAndSignedBody() async throws {
        let base = "https://buy.itunes.apple.com/commerce/machine/authorize?route=a%2Fb%2Bc&opaque=A+B&opaque=second"
        let transport = try LicenseTransport(reply: ["keybag": Data([1])], endpoint: base,
            policy: ["schemes": ["https"], "regex": ["/COMMERCE/"]])
        _ = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
        let requests = await transport.requests
        XCTAssertEqual(requests[1].url!.absoluteString, base + "&guid=020000000001")
        let body = try PropertyListSerialization.propertyList(from: requests[1].httpBody!, format: nil) as! [String: Any]
        XCTAssertEqual(body["guid"] as? String, "020000000001")
        XCTAssertEqual(body["kbsync"] as? Data, Data([3]))
    }
    func testBagCanExcludeURLOrSchemeFromGUIDPolicy() async throws {
        let policies: [[String: Any]] = [["schemes": ["http"], "regex": ["/commerce/"]], ["schemes": ["https"], "regex": ["/unrelated/"]]]
        for policy in policies {
            let transport = try LicenseTransport(reply: ["keybag": Data([1])], policy: policy)
            _ = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
            let requests = await transport.requests
            XCTAssertNil(requests[1].url!.query)
        }
    }
    func testExistingGUIDMustMatchAndCannotBeDuplicated() async throws {
        for query in ["guid=other", "GUID=020000000001", "guid=020000000001&guid=020000000001", "%67uid=other"] {
            let transport = try LicenseTransport(reply: ["keybag": Data([1])], endpoint: "https://buy.itunes.apple.com/commerce/machine/authorize?" + query)
            do {
                _ = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
                XCTFail("Ambiguous device routing accepted")
            } catch { XCTAssertEqual(error as? FairPlayAuthorizationError, .accountMismatch) }
            let requests = await transport.requests; XCTAssertEqual(requests.count, 1)
        }
    }
    func testCustomerMessageWithoutFailureCannotHideGenuineBags() async throws {
        for message in ["", "PRIVATE informational text"] {
            let transport = try LicenseTransport(reply: ["keybag": Data([1]), "customerMessage": message])
            let reply = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
            XCTAssertEqual(reply.keybag, Data([1]))
        }
    }
    func testTrustedRedirectPreservesExactPOSTAndSignature() async throws {
        for status in [301, 302, 307, 308] {
            let transport = AuthorizationSequenceTransport([
                try AuthorizationHTTPReply(status: status, headers: ["Location": "https://p99-buy.itunes.apple.com/commerce/machine/authorize/?route=a%2Fb"]),
                try AuthorizationHTTPReply(["keybag": Data([1])])
            ])
            _ = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
            let requests = await transport.requests
            XCTAssertEqual(requests.count, 3)
            XCTAssertEqual(requests[1].httpBody, requests[2].httpBody)
            XCTAssertEqual(requests[2].httpMethod, "POST")
            XCTAssertEqual(requests[2].value(forHTTPHeaderField: "X-Apple-ActionSignature"), "SIGN")
            XCTAssertEqual(requests[2].url!.absoluteString, "https://p99-buy.itunes.apple.com/commerce/machine/authorize/?route=a%2Fb&guid=020000000001")
        }
    }
    func testRedirectCannotLeakCredentialsOrChangeService() async throws {
        for location in ["https://example.org/commerce/machine/authorize?private=PRIVATE", "http://buy.itunes.apple.com/commerce/machine/authorize", "https://buy.itunes.apple.com/other", "https://buy.itunes.apple.com/commerce/machine/authorize?guid=other"] {
            let transport = AuthorizationSequenceTransport([try AuthorizationHTTPReply(status: 307, headers: ["Location": location])])
            let trace = AuthorizationTrace()
            do {
                _ = try await client(transport, diagnostic: { await trace.append($0) }).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
                XCTFail("Unsafe redirect accepted")
            } catch {}
            let requests = await transport.requests; XCTAssertEqual(requests.count, 2)
            let lines = await trace.lines
            XCTAssertFalse(lines.joined().contains("PRIVATE"))
        }
    }
    func testRedirectLoopsAndHopBudgetStopWithoutRetry() async throws {
        let loops = AuthorizationSequenceTransport([try AuthorizationHTTPReply(status: 307, headers: ["Location": "/commerce/machine/authorize"])])
        do {
            _ = try await client(loops).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
            XCTFail("Loop accepted")
        } catch { XCTAssertEqual(error as? FairPlayAuthorizationError, .tooManyRedirects) }
        let loopRequests = await loops.requests; XCTAssertEqual(loopRequests.count, 2)
        let replies = try (1...4).map { try AuthorizationHTTPReply(status: 307, headers: ["Location": "?hop=\($0)"]) }
        let chain = AuthorizationSequenceTransport(replies)
        do {
            _ = try await client(chain).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
            XCTFail("Unbounded redirect chain")
        } catch { XCTAssertEqual(error as? FairPlayAuthorizationError, .tooManyRedirects) }
        let requests = await chain.requests; XCTAssertEqual(requests.count, 5) // Bag + four POST transmissions.
    }
    func testHTTPAppleAndInvalidReplyAreDistinctAndDoNotLogSecrets() async throws {
        let fixtures: [(AuthorizationHTTPReply, FairPlayAuthorizationError, String)] = [
            (try AuthorizationHTTPReply(["failureType": "123", "customerMessage": "PRIVATE-TOKEN fixture@example.invalid"], status: 403), .httpRejected, "HTTP=403"),
            (try AuthorizationHTTPReply(["failureType": "123", "customerMessage": "PRIVATE-TOKEN fixture@example.invalid"]), .appleRejected, "apple-failure=123"),
            (AuthorizationHTTPReply(raw: Data("<html>PRIVATE-TOKEN</html>".utf8), status: 200), .invalidReply, "body=html")
        ]
        for (reply, expected, marker) in fixtures {
            let transport = AuthorizationSequenceTransport([reply])
            let trace = AuthorizationTrace()
            do {
                _ = try await client(transport, diagnostic: { await trace.append($0) }).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
                XCTFail("Failed reply accepted")
            } catch { XCTAssertEqual(error as? FairPlayAuthorizationError, expected) }
            let lines = await trace.lines; let log = lines.joined(separator: "\n")
            XCTAssertTrue(log.contains(marker))
            for secret in ["PRIVATE-TOKEN", "fixture@example.invalid", "020000000001", "SIGNED-FIXTURE"] { XCTAssertFalse(log.contains(secret)) }
            let requests = await transport.requests; XCTAssertEqual(requests.count, 2)
        }
    }
    func testADIChallengeTakesPrecedenceOverRouting() async throws {
        let transport = AuthorizationSequenceTransport([
            try AuthorizationHTTPReply(["headers": ["x-apple-md-action": "SP"]], status: 307, headers: ["Location": "?next=1"])
        ])
        do {
            _ = try await client(transport).authorize(kbsync: Data([3]), needsDiversity: true, machineName: "Fixture", sign: { _ in "SIGN" })
            XCTFail("ADI challenge followed as redirect")
        } catch { XCTAssertEqual(error as? FairPlayAuthorizationError, .adiChallenge) }
        let requests = await transport.requests; XCTAssertEqual(requests.count, 2)
    }
}
