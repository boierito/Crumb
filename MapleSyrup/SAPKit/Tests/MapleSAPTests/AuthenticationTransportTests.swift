import XCTest
@testable import MapleSAP
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

// Tests the real URLSession transport with local URLProtocol responses, not Apple.
final class AuthenticationTransportTests: XCTestCase {
    func testChallengeCookiesSurviveIsolatedRequestSessions() async throws {
        let cookie = HTTPCookie(properties: [.name: "initial", .value: "fixture-initial", .domain: ".itunes.apple.com", .path: "/", .secure: "TRUE"])!
        let transport = AppleAuthenticationTransport(cookies: [StoreCookie(cookie)], isolatedConnections: true,
            protocolClasses: [ChallengeProtocol.self])
        defer { transport.close() }
        let request = URLRequest(url: URL(string: "https://buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate")!)
        _ = try await transport.send(request)
        _ = try await transport.send(request)
        let cookies = await transport.cookies()
        XCTAssertTrue(cookies.contains { $0.name == "initial" && $0.value == "fixture-initial" })
        XCTAssertTrue(cookies.contains { $0.name == "challenge" && $0.value == "fixture-challenge" })
    }
    func testClosedTransportCannotStartAnotherLoginRequest() async throws {
        let transport = AppleAuthenticationTransport(isolatedConnections: true)
        transport.close()
        do { _ = try await transport.send(URLRequest(url: URL(string: "https://buy.itunes.apple.com/")!)); XCTFail("request started after close") }
        catch { XCTAssertTrue(error is CancellationError) }
    }
}
private final class ChallengeProtocol: URLProtocol {
    override class func canInit(with request: URLRequest) -> Bool { true }
    override class func canonicalRequest(for request: URLRequest) -> URLRequest { request }
    override func startLoading() {
        let headers = ["Content-Type": "application/x-apple-plist",
                       "Set-Cookie": "challenge=fixture-challenge; Domain=.itunes.apple.com; Path=/; Secure"]
        let response = HTTPURLResponse(url: request.url!, statusCode: 200, httpVersion: "HTTP/1.1", headerFields: headers)!
        client?.urlProtocol(self, didReceive: response, cacheStoragePolicy: .notAllowed)
        client?.urlProtocol(self, didLoad: Data("<plist version=\"1.0\"><dict/></plist>".utf8))
        client?.urlProtocolDidFinishLoading(self)
    }
    override func stopLoading() {}
}
