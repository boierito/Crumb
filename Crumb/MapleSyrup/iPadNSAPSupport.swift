// iPadN SAP runtime integration support.
// Upstream code: mengangs/ipadn @ 70f2ac36e7f87cacf82a45b58e6fa27ed23ecb80 (GPL-2.0-or-later).
import Foundation
import Security
import ZIPFoundation

enum CrumbSAPText {
    static func text(_ value: String) -> String { value }
}
enum CrumbSAPToolError: Error, LocalizedError {
    case message(String)
    var errorDescription: String? {
        switch self { case .message(let message): return message }
    }
}
struct CrumbSAPProofRecord: Codable, Equatable, Sendable {
    var identifier: String
    var files: [String: String]
    var updatedAt: Date = Date()
}
enum CrumbSAPProofStore {
    private static let maxEntries = 8
    private static var service: String { (Bundle.main.bundleIdentifier ?? "com.certlium.crumb") + ".ipadn.storeproof" }
    private static func query() -> [String: Any] {
        [kSecClass as String: kSecClassGenericPassword,
         kSecAttrService as String: service,
         kSecAttrAccount as String: "store-proof-v1"]
    }
    static func load() -> [String: CrumbSAPProofRecord] {
        var request = query()
        request[kSecReturnData as String] = true
        request[kSecMatchLimit as String] = kSecMatchLimitOne
        var result: CFTypeRef?
        guard SecItemCopyMatching(request as CFDictionary, &result) == errSecSuccess,
              let data = result as? Data,
              let records = try? JSONDecoder().decode([String: CrumbSAPProofRecord].self, from: data) else {
            return [:]
        }
        return records
    }
    static func save(_ records: [String: CrumbSAPProofRecord]) {
        let kept = records.sorted { $0.value.updatedAt > $1.value.updatedAt }.prefix(maxEntries)
        let compact = Dictionary(uniqueKeysWithValues: kept.map { ($0.key, $0.value) })
        guard let data = try? JSONEncoder().encode(compact) else { return }
        let attributes: [String: Any] = [
            kSecValueData as String: data,
            kSecAttrAccessible as String: kSecAttrAccessibleAfterFirstUnlockThisDeviceOnly
        ]
        if SecItemUpdate(query() as CFDictionary, attributes as CFDictionary) == errSecItemNotFound {
            var item = query()
            attributes.forEach { item[$0.key] = $0.value }
            _ = SecItemAdd(item as CFDictionary, nil)
        }
    }
}
enum CrumbSAPPlist {
    static func dictionary(_ data: Data) throws -> [String: Any] {
        guard let value = try PropertyListSerialization.propertyList(from: data, options: [], format: nil) as? [String: Any] else {
            throw CrumbSAPToolError.message("Invalid Apple property list")
        }
        return value
    }
}
struct CrumbSAPHTTPStatusError: Error {
    let code: Int
}
private final class CrumbSAPRedirectDelegate: NSObject, URLSessionTaskDelegate, @unchecked Sendable {
    let allow: Bool
    init(allow: Bool) { self.allow = allow }
    func urlSession(_ session: URLSession, task: URLSessionTask,
                    willPerformHTTPRedirection response: HTTPURLResponse,
                    newRequest request: URLRequest,
                    completionHandler: @escaping (URLRequest?) -> Void) {
        guard allow, request.url?.scheme == "https",
              let host = request.url?.host?.lowercased(),
              (host == "apple.com" || host.hasSuffix(".apple.com") ||
               host == "mzstatic.com" || host.hasSuffix(".mzstatic.com")) else {
            completionHandler(nil)
            return
        }
        completionHandler(request)
    }
}
final class CrumbSAPHTTPClient: @unchecked Sendable {
    static let shared = CrumbSAPHTTPClient()
    static func encodePlist(_ value: [String: Any]) throws -> Data {
        try PropertyListSerialization.data(fromPropertyList: value, format: .xml, options: 0)
    }
    func request(_ url: URL, method: String = "GET", body: Data? = nil,
                 headers: [String: String] = [:], allowRedirects: Bool = true)
        async throws -> (Data, HTTPURLResponse) {
        guard url.scheme == "https" else { throw CrumbSAPToolError.message("HTTPS is required") }
        var request = URLRequest(url: url, cachePolicy: .reloadIgnoringLocalCacheData, timeoutInterval: 30)
        request.httpMethod = method
        request.httpBody = body
        for (key, value) in headers { request.setValue(value, forHTTPHeaderField: key) }
        let configuration = URLSessionConfiguration.ephemeral
        configuration.timeoutIntervalForRequest = 30
        configuration.timeoutIntervalForResource = 60
        configuration.httpShouldSetCookies = false
        configuration.urlCache = nil
        let session = URLSession(configuration: configuration,
                                 delegate: CrumbSAPRedirectDelegate(allow: allowRedirects), delegateQueue: nil)
        defer { session.finishTasksAndInvalidate() }
        let (data, result) = try await session.data(for: request)
        guard let response = result as? HTTPURLResponse else {
            throw CrumbSAPToolError.message("Invalid Apple response")
        }
        guard data.count <= 4 * 1024 * 1024 else {
            throw CrumbSAPToolError.message("Apple response exceeded 4 MiB")
        }
        return (data, response)
    }
}
enum CrumbSAPAssets {
    // Only the verified build-time archive is extracted. No remote JS or library is executed.
    static let root: URL? = {
        let fileManager = FileManager.default
        guard let archive = Bundle.main.url(forResource: "SAPRuntime", withExtension: "zip") else { return nil }
        let root = fileManager.urls(for: .cachesDirectory, in: .userDomainMask)[0]
            .appendingPathComponent("Crumb-iPadN-SAP-70f2ac36", isDirectory: true)
        let index = root.appendingPathComponent("index.html")
        let wasm = root.appendingPathComponent("sap.wasm")
        if fileManager.fileExists(atPath: index.path) && fileManager.fileExists(atPath: wasm.path) {
            return root
        }
        do {
            try? fileManager.removeItem(at: root)
            try fileManager.createDirectory(at: root, withIntermediateDirectories: true)
            try fileManager.unzipItem(at: archive, to: root)
            guard fileManager.fileExists(atPath: index.path),
                  fileManager.fileExists(atPath: wasm.path),
                  fileManager.fileExists(atPath: root.appendingPathComponent("sap-signer.js").path) else {
                try? fileManager.removeItem(at: root)
                return nil
            }
            return root
        } catch {
            try? fileManager.removeItem(at: root)
            return nil
        }
    }()
    static func rootURL() -> URL? { root }
}
