import Foundation
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

public protocol AuthenticationTransport {
    func send(_ request: URLRequest) async throws -> (Data, HTTPURLResponse)
    func cookies() async -> [StoreCookie]
}

// Each login owns an ephemeral cookie jar and connections. Never allow URLSession
// to change a redirect into GET or replay credentials to an unvalidated host.
public final class AppleAuthenticationTransport: NSObject, AuthenticationTransport, URLSessionTaskDelegate {
    private let configuration: URLSessionConfiguration = {
        let configuration = URLSessionConfiguration.ephemeral
        configuration.timeoutIntervalForRequest = 30
        configuration.timeoutIntervalForResource = 30
        configuration.urlCache = nil
        configuration.httpCookieAcceptPolicy = .always
        configuration.httpShouldSetCookies = true
        return configuration
    }()
    private lazy var session = URLSession(configuration: configuration, delegate: self, delegateQueue: nil)

    public init(cookies: [StoreCookie] = []) {
        super.init()
        // Carry challenge cookies into the verification request, as ipatool's
        // shared jar does. The jar remains ephemeral; no challenge is persisted.
        for cookie in cookies.compactMap({ $0.cookie() }) { configuration.httpCookieStorage?.setCookie(cookie) }
    }

    public func send(_ request: URLRequest) async throws -> (Data, HTTPURLResponse) {
        // Login replies contain tokens. Keep them in bounded RAM, never a
        // URLSession download file, URLCache, or a persistent cookie jar.
        #if canImport(FoundationNetworking)
        let (data, response) = try await session.data(for: request)
        guard data.count <= SAPProtocol.maximumBodySize else { throw SAPError.oversizedResponse }
        #else
        let (bytes, response) = try await session.bytes(for: request)
        var data = Data()
        for try await byte in bytes {
            guard data.count < SAPProtocol.maximumBodySize else { throw SAPError.oversizedResponse }
            data.append(byte)
        }
        #endif
        guard let response = response as? HTTPURLResponse else { throw AuthenticationError.invalidResponse(0) }
        return (data, response)
    }
    public func cookies() async -> [StoreCookie] {
        (configuration.httpCookieStorage?.cookies ?? []).map(StoreCookie.init).filter { $0.cookie() != nil }
    }
    public func urlSession(_ session: URLSession, task: URLSessionTask,
                           willPerformHTTPRedirection response: HTTPURLResponse, newRequest request: URLRequest,
                           completionHandler: @escaping (URLRequest?) -> Void) { completionHandler(nil) }
    public func close() { session.invalidateAndCancel() }
}
