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
    init(reply: [String: Any], endpoint: String = "https://buy.itunes.apple.com/commerce/machine/authorize", status: Int = 200, headers: [String: String] = [:]) throws {
        self.body = try PropertyListSerialization.data(fromPropertyList: reply, format: .xml, options: 0)
        self.endpoint = endpoint; self.status = status; self.headers = headers
    }
    func cookies() async -> [StoreCookie] { [] }
    func send(_ request: URLRequest) async throws -> (Data, HTTPURLResponse) {
        requests.append(request)
        if requests.count == 1 {
            let bag = try PropertyListSerialization.data(fromPropertyList: ["urlBag": ["authorizeMachine": endpoint]], format: .xml, options: 0)
            return (bag, HTTPURLResponse(url: request.url!, statusCode: 200, httpVersion: nil, headerFields: nil)!)
        }
        return (body, HTTPURLResponse(url: request.url!, statusCode: status, httpVersion: nil, headerFields: headers)!)
    }
}

final class FairPlayAuthorizationTests: XCTestCase {
    private func client(_ transport: LicenseTransport) throws -> FairPlayAuthorizationClient {
        let id = try MachineIdentity(hardwareID: Data([2,0,0,0,0,1]))
        let account = StoreAccount(email: "fixture@example.invalid", name: "Fixture", dsid: "1", passwordToken: "PRIVATE-TOKEN", storefront: "143441-1,29", pod: nil, guid: id.guid, authenticationURL: URL(string: "https://buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate")!)
        return try FairPlayAuthorizationClient(account: account, identity: id, transport: transport)
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
