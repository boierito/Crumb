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

    public func send(_ request: URLRequest) async throws -> (Data, HTTPURLResponse) {
        let (file, response) = try await session.download(for: request)
        defer { try? FileManager.default.removeItem(at: file) }
        guard let response = response as? HTTPURLResponse else { throw AuthenticationError.invalidResponse(0) }
        let length = try file.resourceValues(forKeys: [.fileSizeKey]).fileSize ?? 0
        guard length <= SAPProtocol.maximumBodySize else { throw SAPError.oversizedResponse }
        return (try Data(contentsOf: file), response)
    }
    public func cookies() async -> [StoreCookie] {
        (configuration.httpCookieStorage?.cookies ?? []).map(StoreCookie.init).filter { $0.cookie() != nil }
    }
    public func urlSession(_ session: URLSession, task: URLSessionTask,
                           willPerformHTTPRedirection response: HTTPURLResponse, newRequest request: URLRequest,
                           completionHandler: @escaping (URLRequest?) -> Void) { completionHandler(nil) }
    public func close() { session.invalidateAndCancel() }
}
