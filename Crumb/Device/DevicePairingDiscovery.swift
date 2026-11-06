import Foundation
import Darwin

struct DevicePairingDiscoveryResult: Sendable {
    let port: Int?
    let summary: String
}

// Resolve only advertisements whose addresses belong to this iPhone. This
// supplies a port, never a remote host or permission to skip authentication.
@MainActor
final class DevicePairingDiscovery: NSObject, NetServiceBrowserDelegate, NetServiceDelegate {
    private let browser = NetServiceBrowser()
    private var services: [NetService] = []
    private var names = Set<String>()
    private var timeout: Task<Void, Never>?
    private var completion: ((DevicePairingDiscoveryResult) -> Void)?
    private let ownAddresses = DevicePairingDiscovery.interfaceAddresses()
    private var finished = false

    func start(_ completion: @escaping (DevicePairingDiscoveryResult) -> Void) {
        self.completion = completion
        browser.delegate = self
        browser.remove(from: .main, forMode: .default)
        browser.schedule(in: .main, forMode: .common)
        browser.searchForServices(ofType: "_remotepairing._tcp.", inDomain: "local.")
        timeout = Task { [weak self] in
            do { try await Task.sleep(for: .seconds(3)) } catch { return }
            self?.finish(port: nil, summary: "No pairing service from this iPhone was resolved. Trying the default port.")
        }
    }
    func cancel() {
        completion = nil
        stop()
    }
    private func stop() {
        finished = true
        timeout?.cancel(); timeout = nil
        browser.delegate = nil; browser.stop()
        for service in services { service.delegate = nil; service.stop() }
        services.removeAll()
    }
    private func finish(port: Int?, summary: String) {
        guard !finished else { return }
        let callback = completion; completion = nil
        stop()
        callback?(DevicePairingDiscoveryResult(port: port, summary: summary))
    }
    private func resolve(name: String, type: String, domain: String) {
        guard !finished, services.count < 32, names.insert(name).inserted else { return }
        let service = NetService(domain: domain, type: type, name: name)
        services.append(service)
        service.delegate = self
        service.remove(from: .main, forMode: .default)
        service.schedule(in: .main, forMode: .common)
        service.resolve(withTimeout: 2)
    }
    private func resolved(port: Int, addresses: [Data]) {
        guard !finished, (1...65535).contains(port) else { return }
        let belongsHere = addresses.contains { data in
            data.withUnsafeBytes { raw in
                guard raw.count >= MemoryLayout<sockaddr>.size, let base = raw.baseAddress else { return false }
                let address = base.assumingMemoryBound(to: sockaddr.self)
                guard let host = Self.numericHost(address, length: socklen_t(raw.count)) else { return false }
                return ownAddresses.contains(host)
            }
        }
        if belongsHere { finish(port: port, summary: "This iPhone advertises pairing on port \(port).") }
    }
    private static func numericHost(_ address: UnsafePointer<sockaddr>, length: socklen_t) -> String? {
        let family = Int32(address.pointee.sa_family)
        guard (family == AF_INET && Int(length) >= MemoryLayout<sockaddr_in>.size) ||
              (family == AF_INET6 && Int(length) >= MemoryLayout<sockaddr_in6>.size) else { return nil }
        var buffer = [CChar](repeating: 0, count: Int(NI_MAXHOST))
        let capacity = socklen_t(buffer.count)
        guard getnameinfo(address, length, &buffer, capacity, nil, 0, NI_NUMERICHOST) == 0 else { return nil }
        return String(cString: buffer).split(separator: "%").first.map(String.init)
    }
    private static func interfaceAddresses() -> Set<String> {
        var result: Set<String> = ["127.0.0.1", "::1"]
        var first: UnsafeMutablePointer<ifaddrs>?
        guard getifaddrs(&first) == 0 else { return result }
        defer { freeifaddrs(first) }
        var cursor = first
        while let current = cursor {
            defer { cursor = current.pointee.ifa_next }
            if let address = current.pointee.ifa_addr,
               let host = numericHost(address, length: socklen_t(address.pointee.sa_len)) { result.insert(host) }
        }
        return result
    }
    nonisolated func netServiceBrowser(_ browser: NetServiceBrowser, didFind service: NetService, moreComing: Bool) {
        let name = service.name, type = service.type, domain = service.domain
        Task { @MainActor [weak self] in self?.resolve(name: name, type: type, domain: domain) }
    }
    nonisolated func netServiceBrowser(_ browser: NetServiceBrowser, didNotSearch errorDict: [String: NSNumber]) {
        let code = errorDict[NetService.errorCode]?.intValue ?? 0
        Task { @MainActor [weak self] in
            self?.finish(port: nil, summary: "Pairing service lookup failed (\(code)). Trying the default port.")
        }
    }
    nonisolated func netServiceDidResolveAddress(_ sender: NetService) {
        let port = sender.port, addresses = sender.addresses ?? []
        Task { @MainActor [weak self] in self?.resolved(port: port, addresses: addresses) }
    }
}
