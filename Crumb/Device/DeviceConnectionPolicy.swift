import Foundation

// A connection check is an event-driven probe, never a keep-alive timer.
struct DeviceConnectionPolicy {
    private(set) var attempted = false
    private var networkSignature: String?
    mutating func networkChanged(_ signature: String) -> Bool {
        guard signature != networkSignature else { return false }
        networkSignature = signature
        attempted = false
        return true
    }
    mutating func becameInactive() { attempted = false }
    mutating func shouldCheck(paired: Bool, busy: Bool, connected: Bool, cellular: Bool, helperPending: Bool) -> Bool {
        guard paired, !busy, !connected, !cellular, !helperPending, !attempted else { return false }
        attempted = true
        return true
    }
}
enum DeviceRequest {
    static let reserved: Set<String> = ["record", "host_irk", "pairing_kind", "operation", "proof", "route", "own_bundle", "pairing_ports"]
    static func adding(_ extra: [String: Any], to request: [String: Any]) throws -> [String: Any] {
        guard reserved.isDisjoint(with: extra.keys) else {
            throw NSError(domain: "CrumbDeviceRequest", code: 1,
                userInfo: [NSLocalizedDescriptionKey: "Invalid device operation parameters."])
        }
        return request.merging(extra) { _, new in new }
    }
    // Pending-job metadata is local state, never authentication input.
    static func retrievalParameters(_ pending: [String: String]) -> [String: Any] {
        pending.filter { ["job_id", "helper_bundle", "output"].contains($0.key) }
    }
}

struct DeviceResultRequest {
    private(set) var requested = false
    mutating func enqueue() { requested = true }
    mutating func claim(busy: Bool, hasPendingJob: Bool) -> Bool {
        guard hasPendingJob else { requested = false; return false }
        guard requested, !busy else { return false }
        requested = false
        return true
    }
}

// Shared by foreground publication and retired-worker cleanup. A retired worker
// can release only its captured workspace, never a newer session's state.
enum DeviceJobState {
    static let terminalPhases: Set<String> = ["signing-plan", "installation-signing-ready", "legacy-result", "pairing-imported", "signing-ready", "signed", "paired", "connected", "prepared", "installed", "complete", "failed", "cancelled"]
}
