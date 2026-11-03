import Foundation
import SwiftUI
import Zip

@MainActor
final class DeviceWorkflow: NSObject, ObservableObject, NetServiceDelegate {
    static let shared = DeviceWorkflow()
    // Pairing introduced in 27. Decryption reference was tested only on 27.0.1.
    // Later versions require a reviewed change; do not use an unbounded >= check.
    static var pairingSupported: Bool {
        let v = ProcessInfo.processInfo.operatingSystemVersion
        return v.majorVersion == 27 && v.minorVersion == 0 && (0...1).contains(v.patchVersion)
    }
    static var decryptSupported: Bool {
        let v = ProcessInfo.processInfo.operatingSystemVersion
        return v.majorVersion == 27 && v.minorVersion == 0 && v.patchVersion == 1
    }
    @Published private(set) var phase = "idle"
    @Published private(set) var message = ""
    @Published private(set) var pin: String?
    @Published private(set) var connected = false
    @Published private(set) var paired = false
    @Published private(set) var busy = false
    @Published private(set) var output: URL?
    @Published private(set) var helperInstalled = false
    @Published private(set) var hasCertificate = false
    @Published private(set) var hasProfile = false
    @Published var helperBundle = "com.certlium.crumb.decrypt"
    @Published var route = "10.7.0.1"
    private var service: NetService?
    private var job: UnsafeMutableRawPointer?
    private var generation = UUID()
    private var polling: Task<Void, Never>?
    private var background: UIBackgroundTaskIdentifier = .invalid
    private var staged: URL?
    private var pending: [String: String]?
    private var lastRequest: [String: Any] = [:]

    override private init() {
        super.init()
        do {
            let fm = FileManager.default
            let base = try fm.url(for: .applicationSupportDirectory, in: .userDomainMask, appropriateFor: nil, create: true).appendingPathComponent("DeviceJobs", isDirectory: true)
            // The old process has no surviving worker; remove interrupted signing workspaces.
            if let children = try? fm.contentsOfDirectory(at: base, includingPropertiesForKeys: nil) {
                for child in children where UUID(uuidString: child.lastPathComponent) != nil { try? fm.removeItem(at: child) }
            }
            paired = try DeviceKeychain.read("pairing") != nil
            hasCertificate = try DeviceKeychain.read("p12") != nil
            hasProfile = try DeviceKeychain.read("profile") != nil
            if let data = try DeviceKeychain.read("pending"),
               let value = try JSONSerialization.jsonObject(with: data) as? [String: String] {
                pending = value; helperBundle = value["helper_bundle"] ?? helperBundle; helperInstalled = true
            }
        } catch { message = "Could not read protected device settings." }
    }
    func importFile(_ url: URL, kind: String) {
        guard !busy else { return }
        let access = url.startAccessingSecurityScopedResource(); defer { if access { url.stopAccessingSecurityScopedResource() } }
        do {
            let size = try url.resourceValues(forKeys: [.fileSizeKey]).fileSize ?? 0
            guard size > 0 && size <= 2 * 1024 * 1024 else { throw "Signing file must be between 1 byte and 2 MiB." }
            try DeviceKeychain.write(Data(contentsOf: url), account: kind)
            if kind == "p12" { hasCertificate = true } else { hasProfile = true }
            message = "Signing file stored on this device."
        } catch { message = error.localizedDescription }
    }
    private func credentials() throws -> [String: Any] {
        guard let data = try DeviceKeychain.read("pairing"),
              let value = try JSONSerialization.jsonObject(with: data) as? [String: Any] else { throw "Pair this iPhone first." }
        return value
    }
    func pair() {
        guard Self.pairingSupported else { message = "On-device pairing is enabled only for iOS 27.0–27.0.1 in this build."; return }
        connected = false
        var request = (try? credentials()) ?? [:]; request["operation"] = "pair"
        start(request)
    }
    func verify() { deviceOperation("verify") }
    private func deviceOperation(_ operation: String, extra: [String: Any] = [:]) {
        guard !busy else { return }
        guard Self.pairingSupported else { message = "This system version has not been enabled for the device connection."; return }
        do {
            var request = try credentials(); request["operation"] = operation
            let proof = UUID().uuidString
            let docs = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask)[0]
            try Data(proof.utf8).write(to: docs.appendingPathComponent(".crumb-device-proof"), options: [.atomic, .completeFileProtection])
            request["proof"] = proof
            request["route"] = route; request["own_bundle"] = Bundle.main.bundleIdentifier
            extra.forEach { request[$0.key] = $0.value }; start(request)
        } catch { message = error.localizedDescription }
    }
    func prepare(_ record: DownloadRecord, password: String) {
        guard !busy, connected, Self.decryptSupported else { message = "Connect on iOS 27.0.1 before preparing a helper."; return }
        do {
            guard let ipa = record.fileURL, let template = Bundle.main.url(forResource: "CrumbDecryptTemplate", withExtension: "zip"),
                  let cert = try DeviceKeychain.read("p12"), let profile = try DeviceKeychain.read("profile") else { throw "Import a signing certificate and provisioning profile first." }
            let fm = FileManager.default
            let base = try fm.url(for: .applicationSupportDirectory, in: .userDomainMask, appropriateFor: nil, create: true)
                .appendingPathComponent("DeviceJobs", isDirectory: true)
            try fm.createDirectory(at: base, withIntermediateDirectories: true, attributes: [.protectionKey: FileProtectionType.complete])
            let id = UUID().uuidString; let root = base.appendingPathComponent(id, isDirectory: true)
            try fm.createDirectory(at: root, withIntermediateDirectories: false, attributes: [.protectionKey: FileProtectionType.complete])
            var excluded = URLResourceValues(); excluded.isExcludedFromBackup = true
            var excludedRoot = root; try excludedRoot.setResourceValues(excluded)
            staged = root
            let certURL = root.appendingPathComponent("identity.p12"), profileURL = root.appendingPathComponent("helper.mobileprovision")
            try cert.write(to: certURL, options: [.atomic, .completeFileProtection]); try profile.write(to: profileURL, options: [.atomic, .completeFileProtection])
            let unzipped = root.appendingPathComponent("template", isDirectory: true)
            try fm.createDirectory(at: unzipped, withIntermediateDirectories: false)
            // This archive is generated by our build, never supplied by an IPA.
            try Zip.unzipFile(template, destination: unzipped, overwrite: false, password: nil)
            let request: [String: Any] = ["operation":"prepare", "ipa":ipa.path, "template":unzipped.path,
                "work":root.appendingPathComponent("work").path, "p12":certURL.path, "profile":profileURL.path,
                "password":password, "helper_bundle":helperBundle, "job_id":id]
            let outputDir = base.appendingPathComponent("Decrypted", isDirectory: true)
            try fm.createDirectory(at: outputDir, withIntermediateDirectories: true)
            pending = ["job_id":id, "helper_bundle":helperBundle, "output":outputDir.appendingPathComponent(id + "-decrypted.ipa").path]
            output = nil; helperInstalled = false; start(request)
        } catch { cleanup(); message = error.localizedDescription }
    }
    func retrieve() {
        guard Self.decryptSupported, let pending else { return }
        deviceOperation("fetch", extra: pending)
    }
    func openHelper() {
        guard helperInstalled else { return }
        UIApplication.shared.open(URL(string: "crumb-decrypt-helper://run")!, options: [:]) { [weak self] ok in
            if !ok { Task { @MainActor in self?.message = "Open Crumb Decrypt from the Home Screen, then return to retrieve the result." } }
        }
    }
    func forgetPairing() {
        guard !busy else { return }
        do { try DeviceKeychain.remove("pairing"); paired = false; connected = false; message = "Local pairing keys removed. The system's trusted-host entry can be removed in Settings." }
        catch { message = error.localizedDescription }
    }
    func forgetSigning() {
        guard !busy else { return }
        do { try DeviceKeychain.remove("p12"); try DeviceKeychain.remove("profile"); hasCertificate = false; hasProfile = false }
        catch { message = error.localizedDescription }
    }
    func cancel() {
        guard let job else { return }
        crumb_device_cancel(job)
        message = "Cancelling…"
        connected = false
        stopAdvertising(); pin = nil
        // Keep polling until the worker has stopped before deleting its workspace.
    }
    private func start(_ request: [String: Any]) {
        guard !busy else { return }
        do {
            let data = try JSONSerialization.data(withJSONObject: request)
            guard let json = String(data: data, encoding: .utf8), let handle = crumb_device_start(json) else { throw "Could not start device operation." }
            job = handle; busy = true; phase = "starting"; message = ""; pin = nil
            lastRequest = request.filter { $0.key != "password" && $0.key != "record" && $0.key != "host_irk" }
            background = UIApplication.shared.beginBackgroundTask(withName: "Crumb device operation") { [weak self] in
                Task { @MainActor in self?.cancel(); self?.message = "iOS ended the background session. Return to Crumb and retry." }
            }
            let operationID = UUID(); generation = operationID
            polling = Task { [weak self] in
                while !Task.isCancelled {
                    guard let self, self.generation == operationID, let job = self.job else { return }
                    if let ptr = crumb_device_poll(job) {
                        let bytes = Data(String(cString: ptr).utf8); crumb_device_string_free(ptr)
                        if let value = try? JSONSerialization.jsonObject(with: bytes) as? [String: Any] { self.receive(value) }
                    }
                    if !self.busy || self.generation != operationID { return }
                    try? await Task.sleep(for: .milliseconds(200))
                }
            }
        } catch { message = error.localizedDescription; busy = false }
    }
    private func receive(_ value: [String: Any]) {
        guard let phase = value["phase"] as? String else { return }
        self.phase = phase
        if phase == "waiting" || phase == "pin" {
            if service == nil, let name = value["service"] as? String, let port = value["port"] as? Int,
               let txt = value["txt"] as? [String: String] {
                let service = NetService(domain: "local.", type: "_remotepairing-pairable-host._tcp.", name: name, port: Int32(port))
                service.delegate = self; service.setTXTRecord(NetService.data(fromTXTRecord: txt.mapValues { Data($0.utf8) }))
                self.service = service; service.publish()
            }
            pin = value["pin"] as? String
            return
        }
        guard ["paired", "connected", "prepared", "installed", "complete", "failed", "cancelled"].contains(phase) else { return }
        let request = lastRequest
        if let job { crumb_device_free(job); self.job = nil }
        busy = false; stopAdvertising(); endBackground(); pin = nil; lastRequest = [:]
        do {
            switch phase {
            case "paired":
                guard let record = value["record"] as? String, let irk = value["host_irk"] as? String else { throw "Missing pairing credentials." }
                try DeviceKeychain.write(JSONSerialization.data(withJSONObject: ["record":record, "host_irk":irk]), account: "pairing")
                paired = true; message = "Paired. Connect LocalDevVPN and check the connection."
            case "connected": connected = true; message = "Authenticated device connection available."
            case "prepared":
                guard let path = value["path"] as? String else { throw "Missing prepared helper." }
                if let staged {
                    try? FileManager.default.removeItem(at: staged.appendingPathComponent("identity.p12"))
                    try? FileManager.default.removeItem(at: staged.appendingPathComponent("helper.mobileprovision"))
                }
                deviceOperation("install", extra: ["path":path, "helper_bundle":helperBundle])
            case "installed":
                helperInstalled = true
                if let pending { try DeviceKeychain.write(JSONSerialization.data(withJSONObject: pending), account: "pending") }
                cleanup(); message = "Helper installed. Open Crumb Decrypt, keep it on screen until it finishes, then return here."
            case "complete":
                guard let path = value["path"] as? String else { throw "Missing verified result." }
                output = URL(fileURLWithPath: path); message = "IPA retrieved; hashes and encryption markers verified. Application execution has not been tested."
            case "failed":
                connected = false; message = (value["message"] as? String) ?? "Device operation failed."
                if ["prepare", "install"].contains(request["operation"] as? String ?? "") { cleanup() }
            default: cleanup(); message = "Operation cancelled."
            }
        } catch { message = error.localizedDescription; cleanup() }
    }
    private func stopAdvertising() { service?.stop(); service = nil }
    private func endBackground() { if background != .invalid { UIApplication.shared.endBackgroundTask(background); background = .invalid } }
    private func cleanup() { if let staged { try? FileManager.default.removeItem(at: staged) }; staged = nil }
    nonisolated func netService(_ sender: NetService, didNotPublish errorDict: [String: NSNumber]) {
        Task { @MainActor in self.cancel(); self.message = "Could not advertise the pairing host. Allow Local Network access in Settings." }
    }
}
