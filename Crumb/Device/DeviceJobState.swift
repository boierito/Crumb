import Foundation

// Shared by foreground publication and retired-worker cleanup.
enum DeviceJobState {
    static let terminalPhases: Set<String> = ["signing-plan", "installation-signing-ready", "signed", "failed", "cancelled"]
}
