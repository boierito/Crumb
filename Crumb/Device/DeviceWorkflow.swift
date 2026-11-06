import Foundation
import Combine
import UIKit
import SwiftUI
import Zip
import Network

enum DeviceOperationError: LocalizedError {
    case message(String)
    var errorDescription: String? { if case let .message(text) = self { return text }; return nil }
}

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
    @Published var setupVisible = false
    @Published var setupNeeded = false
    @Published var errorMessage = ""
    @Published private(set) var signingReady = false
    @Published private(set) var certificateName = ""
    @Published private(set) var profileName = ""
    @Published var signingPassword = ""
    @Published private(set) var onCellular = false
    @Published private(set) var progress: Double?
    @Published private(set) var activeRecord: DownloadRecord?
    @Published private(set) var waitingForHelper = false
    private var installAfterDecrypt = false
    private var queuedRecord: DownloadRecord?
    private var signingRecord: DownloadRecord?
    private var signingContext: [String: Any]?
    private var forceSigningReview = false
    private var signingApprovalPending = false
    @Published var installationSetupNeeded = false
    @Published private(set) var installationReady = false
    @Published private(set) var hasInstallationCertificate = false
    @Published private(set) var hasInstallationProfile = false
    @Published private(set) var installationCertificateName = ""
    @Published private(set) var installationProfileName = ""
    @Published private(set) var installationProfiles: [AdditionalSigningProfile] = []
    @Published var installationPassword = ""
    @Published var signingReview: SigningReview?
    @Published var signedInstallation: SignedInstallation?
    private var queuedInstall: DownloadRecord?
    private let network = NWPathMonitor()
    private var connectionPolicy = DeviceConnectionPolicy()
    private var resultRequest = DeviceResultRequest()
    @Published private(set) var helperBundle = "com.certlium.crumb.decrypt"
    private let route = "10.7.0.1"
    private var service: NetService?
    private var job: UnsafeMutableRawPointer?
    private var generation = UUID()
    private var polling: Task<Void, Never>?
    private var background: UIBackgroundTaskIdentifier = .invalid
    private var staged: URL?
    private var managedHelper: String?
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
            if let data = try DeviceKeychain.read("managed-helper") { managedHelper = String(data: data, encoding: .utf8) }
            paired = try DeviceKeychain.read("pairing") != nil
            hasCertificate = try DeviceKeychain.read("p12") != nil
            hasProfile = try DeviceKeychain.read("profile") != nil
            signingPassword = String(data: try DeviceKeychain.read("password") ?? Data(), encoding: .utf8) ?? ""
            certificateName = String(data: try DeviceKeychain.read("p12-name") ?? Data(), encoding: .utf8) ?? ""
            profileName = String(data: try DeviceKeychain.read("profile-name") ?? Data(), encoding: .utf8) ?? ""
            if let data = try DeviceKeychain.read("signing-info"),
               let info = try JSONSerialization.jsonObject(with: data) as? [String: String], hasCertificate && hasProfile {
                helperBundle = info["helper_bundle"] ?? helperBundle
                signingReady = true
            }
            if let data = try DeviceKeychain.read("pending"),
               let value = try JSONSerialization.jsonObject(with: data) as? [String: String] {
                pending = value; helperInstalled = true; waitingForHelper = true
                installAfterDecrypt = value["install"] == "true"
                if let text = value["download_record"] ?? value["record"], let data = Data(base64Encoded: text) {
                    activeRecord = try? JSONDecoder().decode(DownloadRecord.self, from: data)
                }
            }
        } catch { message = "Could not read protected device settings." }
        seedInstallationSigning()
        network.pathUpdateHandler = { [weak self] path in
            let cellular = path.usesInterfaceType(.cellular) && !path.usesInterfaceType(.wifi)
            let signature = "\(path.status)-\(path.usesInterfaceType(.wifi))-\(cellular)"
            Task { @MainActor in
                guard let self else { return }
                self.onCellular = cellular
                if self.connectionPolicy.networkChanged(signature) {
                    self.connected = false
                    self.autoConnect()
                }
            }
        }
        network.start(queue: DispatchQueue(label: "Crumb.network"))
        if let path = pending?["output"], activeRecord == nil {
            if FileManager.default.fileExists(atPath:path) { start(["operation":"inspect-result","ipa":path]) }
            else { pending = nil; helperInstalled = false; waitingForHelper = false; try? DeviceKeychain.remove("pending") }
        }
    }
    func importFile(_ url: URL, kind: String) {
        guard !busy else { return }
        let access = url.startAccessingSecurityScopedResource(); defer { if access { url.stopAccessingSecurityScopedResource() } }
        do {
            let size = try url.resourceValues(forKeys: [.fileSizeKey]).fileSize ?? 0
            guard size > 0 && size <= 2 * 1024 * 1024 else { throw DeviceOperationError.message("Signing file must be between 1 byte and 2 MiB.") }
            let data = try Data(contentsOf: url)
            if kind == "pairing" {
                guard data.count <= 65536 else { throw DeviceOperationError.message("Pairing file is too large.") }
                start(["operation":"import-pairing", "record":data.base64EncodedString()]); return
            }
            if kind.hasPrefix("install-") {
                if kind == "install-extra-profile" {
                    guard installationProfiles.count < 32 else { throw DeviceOperationError.message("Maximum 32 extension profiles.") }
                    installationProfiles.append(AdditionalSigningProfile(id:UUID().uuidString,name:url.lastPathComponent,data:data))
                    try DeviceKeychain.write(JSONEncoder().encode(installationProfiles),account:"install-profiles")
                } else {
                    try DeviceKeychain.write(data,account:kind)
                    try DeviceKeychain.write(Data(url.lastPathComponent.utf8),account:kind + "-name")
                    invalidateInstallationSigning()
                    seedInstallationSigning()
                }
                return
            }
            try DeviceKeychain.write(data, account: kind)
            try DeviceKeychain.write(Data(url.lastPathComponent.utf8), account: kind + "-name")
            try DeviceKeychain.remove("signing-info")
            signingReady = false
            if kind == "p12" { certificateName = url.lastPathComponent } else { profileName = url.lastPathComponent }
            if kind == "p12" { hasCertificate = true } else { hasProfile = true }
            message = "Signing file stored on this device."
        } catch { errorMessage = error.localizedDescription }
    }
    private func credentials() throws -> [String: Any] {
        guard let data = try DeviceKeychain.read("pairing"),
              let value = try JSONSerialization.jsonObject(with: data) as? [String: Any] else { throw DeviceOperationError.message("Pair this iPhone first.") }
        return value
    }
    func pair() {
        guard Self.pairingSupported else { message = "On-device pairing is enabled only for iOS 27.0–27.0.1 in this build."; return }
        connected = false
        var request = (try? credentials()) ?? [:]; request["operation"] = "pair"
        start(request)
    }
    func verify() { deviceOperation("verify") }
    var ready: Bool { paired && signingReady }
    var activityTitle: String {
        if busy && lastRequest["operation"] as? String == "fetch" { return "Saving app" }
        if waitingForHelper { return "Decrypting" }
        switch phase {
        case "preparing": return "Preparing decrypter"
        case "installing": return "Installing decrypter"
        case "fetching": return "Saving app"
        case "signing": return "Signing app"
        default: return "Connecting"
        }
    }
    func becameInactive() {
        connectionPolicy.becameInactive()
        connected = false
    }
    func autoConnect() {
        guard connectionPolicy.shouldCheck(paired: paired, busy: busy, connected: connected,
            cellular: onCellular, helperPending: waitingForHelper || resultRequest.requested || signingContext != nil || signedInstallation != nil) else { return }
        verify()
    }
    private func signingRequest(operation: String, installation: Bool = false) throws -> [String: Any] {
        let fm = FileManager.default
        let base = try fm.url(for: .applicationSupportDirectory, in: .userDomainMask, appropriateFor: nil, create: true).appendingPathComponent("DeviceJobs", isDirectory: true)
        try fm.createDirectory(at: base, withIntermediateDirectories: true, attributes: [.protectionKey:FileProtectionType.complete])
        let root = base.appendingPathComponent(UUID().uuidString, isDirectory:true)
        try fm.createDirectory(at: root, withIntermediateDirectories:false, attributes:[.protectionKey:FileProtectionType.complete])
        var excluded = URLResourceValues(); excluded.isExcludedFromBackup = true
        var protectedRoot = root; try protectedRoot.setResourceValues(excluded)
        staged = root
        guard let cert = try DeviceKeychain.read(installation ? "install-p12" : "p12"), let profile = try DeviceKeychain.read(installation ? "install-profile" : "profile") else { throw DeviceOperationError.message("Import your signing files.") }
        let p12 = root.appendingPathComponent("identity.p12"), mobile = root.appendingPathComponent("profile.mobileprovision")
        try cert.write(to:p12,options:[.atomic,.completeFileProtection]); try profile.write(to:mobile,options:[.atomic,.completeFileProtection])
        var request: [String:Any] = ["operation":operation,"p12":p12.path,"profile":mobile.path,
            "password":installation ? installationPassword : signingPassword,
            "purpose":installation ? "installation" : "helper", "helper_bundle":helperBundle,
            "protected_bundle_ids":[helperBundle,managedHelper].compactMap { $0 },
            "own_bundle":Bundle.main.bundleIdentifier ?? "com.certlium.crumb", "work":root.appendingPathComponent("work").path]
        if installation {
            var profiles: [String] = []
            for profile in installationProfiles {
                let url = root.appendingPathComponent(profile.id + ".mobileprovision")
                try profile.data.write(to:url,options:[.atomic,.completeFileProtection]); profiles.append(url.path)
            }
            request["profiles"] = profiles
        }
        return request
    }
    func seedInstallationSigning() {
        do {
            // Migrate once, then keep the installation identity independent.
            for key in ["p12","profile","password","p12-name","profile-name"] {
                if try DeviceKeychain.read("install-" + key) == nil, let data = try DeviceKeychain.read(key) {
                    try DeviceKeychain.write(data,account:"install-" + key)
                }
            }
            hasInstallationCertificate = try DeviceKeychain.read("install-p12") != nil
            hasInstallationProfile = try DeviceKeychain.read("install-profile") != nil
            installationCertificateName = String(data:try DeviceKeychain.read("install-p12-name") ?? Data(),encoding:.utf8) ?? ""
            installationProfileName = String(data:try DeviceKeychain.read("install-profile-name") ?? Data(),encoding:.utf8) ?? ""
            if installationPassword.isEmpty { installationPassword = String(data:try DeviceKeychain.read("install-password") ?? Data(),encoding:.utf8) ?? "" }
            let validated = try DeviceKeychain.read("install-signing-info") != nil
            installationReady = hasInstallationCertificate && hasInstallationProfile && validated
            if let data = try DeviceKeychain.read("install-profiles") { installationProfiles = try JSONDecoder().decode([AdditionalSigningProfile].self,from:data) }
        } catch { errorMessage = error.localizedDescription }
    }
    func invalidateInstallationSigning() {
        installationReady = false; try? DeviceKeychain.remove("install-signing-info")
    }
    func removeInstallationProfile(_ id:String) {
        guard !busy else { return }
        installationProfiles.removeAll { $0.id == id }
        do { try DeviceKeychain.write(JSONEncoder().encode(installationProfiles),account:"install-profiles") }
        catch { errorMessage = error.localizedDescription }
    }
    func saveInstallationSigning() {
        guard !busy else { return }
        do { start(try signingRequest(operation:"signing-info",installation:true)) }
        catch { cleanup(); errorMessage = error.localizedDescription }
    }
    func invalidateSigning() {
        signingReady = false
        try? DeviceKeychain.remove("signing-info")
    }
    func saveSigning() {
        guard !busy else { return }
        do { start(try signingRequest(operation:"signing-info")) }
        catch { cleanup(); errorMessage = error.localizedDescription }
    }
    func decrypt(_ record: DownloadRecord, install: Bool = false) {
        guard (!busy || lastRequest["operation"] as? String == "verify"), !waitingForHelper, signingContext == nil, signedInstallation == nil else { errorMessage = "Finish the current app first."; return }
        guard Self.decryptSupported else { errorMessage = "Decryption is supported on iOS 27.0.1. Pairing on other versions does not enable this method."; return }
        queuedRecord = record; installAfterDecrypt = install
        continueQueued()
    }
    func continueQueued() {
        if ready && connected { setupNeeded = false }
        if installationReady && !busy, let record = queuedInstall {
            queuedInstall = nil; installationSetupNeeded = false; install(record,review:forceSigningReview); return
        }
        guard let record = queuedRecord, !busy else { return }
        guard ready else { setupNeeded = true; return }
        guard connected else { autoConnect(); return }
        queuedRecord = nil; activeRecord = record
        prepare(record, password:signingPassword)
    }
    func install(_ record: DownloadRecord, review: Bool = false) {
        // Installing a local decrypted IPA does not require a VPN proof probe.
        if busy && lastRequest["operation"] as? String == "verify" && queuedRecord == nil && queuedInstall == nil && !waitingForHelper {
            stopDecryption(notifyHelper:false)
        }
        guard !busy, !waitingForHelper, signingContext == nil, signedInstallation == nil else { errorMessage = "Finish the current app first."; return }
        guard record.isDecrypted else { AppData.shared.installationRequest = record; return }
        forceSigningReview = review
        seedInstallationSigning()
        guard installationReady else {
            queuedInstall = record
            if hasInstallationCertificate && hasInstallationProfile { activeRecord = record; saveInstallationSigning() }
            else { installationSetupNeeded = true }
            return
        }
        do {
            guard let url = record.fileURL else { throw CocoaError(.fileNoSuchFile) }
            var request = try signingRequest(operation:"signing-plan",installation:true)
            request["ipa"] = url.path; request["output"] = staged!.appendingPathComponent("signed.ipa").path
            signingRecord = record; activeRecord = record; signingContext = request
            start(request)
        } catch { cleanup(); activeRecord = nil; errorMessage = error.localizedDescription }
    }
    func installationSetupDismissed() {
        if installationReady { continueQueued() }
        else { queuedInstall = nil; if !busy { activeRecord = nil } }
    }
    func replanSigning(bundleID:String) {
        guard !busy, var context = signingContext else { return }
        context["bundle_id"] = bundleID; context["prepared"] = true
        signingContext = context; forceSigningReview = true; start(context)
    }
    func approveSigning() {
        guard !busy, signingContext != nil else { return }
        if let review = signingReview {
            signingContext?["bundle_id"] = review.bundleID
            signingApprovalPending = true; signingReview = nil
        }
        else { beginSigning() }
    }
    func signingReviewDismissed() {
        if signingApprovalPending { signingApprovalPending = false; beginSigning() }
        else if signingContext != nil { stopDecryption() }
    }
    private func beginSigning() {
        guard !busy, var context = signingContext else { return }
        context["operation"] = "sign"; context["prepared"] = true
        signingContext = nil; start(context)
    }

    private func deviceOperation(_ operation: String, extra: [String: Any] = [:]) {
        guard !busy else { return }
        guard paired else { setupNeeded = true; return }
        performDeviceOperation(operation, extra: extra, ports: [49152])
    }

    private func performDeviceOperation(_ operation: String, extra: [String: Any], ports: [Int]) {
        guard !busy else { return }
        do {
            var request = try credentials(); request["operation"] = operation
            request["pairing_ports"] = ports
            let proof = UUID().uuidString
            let docs = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask)[0]
            try Data(proof.utf8).write(to: docs.appendingPathComponent(".crumb-device-proof"), options: [.atomic, .completeFileProtection])
            request["proof"] = proof
            request["route"] = route; request["own_bundle"] = Bundle.main.bundleIdentifier
            start(try DeviceRequest.adding(extra, to: request))
        } catch { errorMessage = error.localizedDescription }
    }
    func prepare(_ record: DownloadRecord, password: String) {
        guard !busy, connected, Self.decryptSupported else { activeRecord = nil; errorMessage = "Connect the decrypter on iOS 27.0.1 first."; return }
        do {
            guard let ipa = record.fileURL, let template = Bundle.main.url(forResource: "CrumbDecryptTemplate", withExtension: "zip"),
                  let cert = try DeviceKeychain.read("p12"), let profile = try DeviceKeychain.read("profile") else { throw DeviceOperationError.message("Import a signing certificate and provisioning profile first.") }
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
                "password":password, "helper_bundle":helperBundle, "job_id":id, "own_bundle":Bundle.main.bundleIdentifier ?? "com.certlium.crumb"]
            let outputDir = base.appendingPathComponent("Decrypted", isDirectory: true)
            try fm.createDirectory(at: outputDir, withIntermediateDirectories: true)
            pending = ["job_id":id, "helper_bundle":helperBundle, "output":outputDir.appendingPathComponent(id + "-decrypted.ipa").path,
                       "download_record":try JSONEncoder().encode(record).base64EncodedString(), "install":installAfterDecrypt ? "true" : "false", "cancel_protocol":"1"]
            output = nil; helperInstalled = false; start(request)
        } catch { cleanup(); activeRecord = nil; errorMessage = error.localizedDescription }
    }
    @discardableResult func receiveResultURL(_ url: URL) -> Bool {
        guard ["device-result", "device-cancelled"].contains(url.host ?? ""),
              let jobID = URLComponents(url: url, resolvingAgainstBaseURL: false)?.queryItems?.first(where: { $0.name == "job" })?.value,
              UUID(uuidString: jobID) != nil, pending?["job_id"] == jobID else { return false }
        if url.host == "device-cancelled" { stopDecryption(notifyHelper: false) } else { retrieve() }
        return true
    }
    func retrieve() {
        guard Self.decryptSupported, pending != nil else { return }
        resultRequest.enqueue()
        drainResultRequest()
    }
    private func drainResultRequest() {
        guard resultRequest.claim(busy:busy, hasPendingJob:pending != nil), let pending else { return }
        deviceOperation("fetch", extra: DeviceRequest.retrievalParameters(pending))
    }
    func openHelper() {
        guard helperInstalled else { return }
        UIApplication.shared.open(URL(string: "crumb-decrypt-helper://run")!, options: [:]) { [weak self] ok in
            if !ok { Task { @MainActor in self?.errorMessage = "Open Crumb Decrypt from the Home Screen, then return to retrieve the result." } }
        }
    }
    func forgetPairing() {
        guard !busy else { return }
        do { try DeviceKeychain.remove("pairing"); paired = false; connected = false; message = "Local pairing keys removed. The system's trusted-host entry can be removed in Settings." }
        catch { message = error.localizedDescription }
    }
    func forgetSigning() {
        guard !busy else { return }
        do { try DeviceKeychain.remove("p12"); try DeviceKeychain.remove("profile"); try DeviceKeychain.remove("password"); try DeviceKeychain.remove("signing-info")
            try DeviceKeychain.remove("p12-name"); try DeviceKeychain.remove("profile-name")
            signingReady = false; signingPassword = ""; certificateName = ""; profileName = ""; hasCertificate = false; hasProfile = false }
        catch { message = error.localizedDescription }
    }
    func cancel() {
        guard let job else { return }
        resultRequest = DeviceResultRequest()
        crumb_device_cancel(job)
        message = "Cancelling…"; connected = false
        stopAdvertising(); pin = nil
    }

    // Abandon the UI session immediately, including an external helper with no
    // native Job. A cancelled worker owns its workspace until it has stopped.
    func stopDecryption(notifyHelper: Bool = true) {
        let helperJob = waitingForHelper && pending?["cancel_protocol"] == "1" ? pending?["job_id"] : nil
        let retiredJob = job
        let retiredWorkspace = staged
        generation = UUID()
        polling?.cancel(); polling = nil; job = nil; staged = nil
        if let retiredJob { crumb_device_cancel(retiredJob) }
        resultRequest = DeviceResultRequest()
        queuedRecord = nil; queuedInstall = nil; signingRecord = nil
        signingContext = nil; signingApprovalPending = false; signingReview = nil
        pending = nil; activeRecord = nil; helperInstalled = false; waitingForHelper = false
        busy = false; progress = nil; pin = nil; phase = "idle"; message = ""; errorMessage = ""
        connected = false; connectionPolicy.becameInactive(); lastRequest = [:]
        stopAdvertising(); endBackground()
        do { try DeviceKeychain.remove("pending") }
        catch { errorMessage = "Could not remove the saved job. Please retry." }
        if let retiredJob {
            Task { @MainActor in
                while true {
                    if let pointer = crumb_device_poll(retiredJob) {
                        let data = Data(String(cString: pointer).utf8)
                        crumb_device_string_free(pointer)
                        if let value = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
                           let phase = value["phase"] as? String,
                           DeviceJobState.terminalPhases.contains(phase) {
                            crumb_device_free(retiredJob)
                            if let retiredWorkspace { try? FileManager.default.removeItem(at: retiredWorkspace) }
                            return
                        }
                    }
                    try? await Task.sleep(for: .milliseconds(250))
                }
            }
        } else if let retiredWorkspace { try? FileManager.default.removeItem(at: retiredWorkspace) }
        if notifyHelper, let helperJob {
            var url = URLComponents()
            url.scheme = "crumb-decrypt-helper"; url.host = "cancel"
            url.queryItems = [URLQueryItem(name: "job", value: helperJob)]
            if let target = url.url { UIApplication.shared.open(target, options: [:], completionHandler: nil) }
        }
    }
    private func start(_ request: [String: Any]) {
        guard !busy else { return }
        do {
            let data = try JSONSerialization.data(withJSONObject: request)
            guard let json = String(data: data, encoding: .utf8), let handle = crumb_device_start(json) else { throw DeviceOperationError.message("Could not start device operation.") }
            job = handle; progress = nil; busy = true; phase = "starting"; message = ""; pin = nil
            lastRequest = request.filter { $0.key != "password" && $0.key != "record" && $0.key != "host_irk" }
            let operationID = UUID(); generation = operationID
            background = UIApplication.shared.beginBackgroundTask(withName: "Crumb device operation") { [weak self] in
                Task { @MainActor in
                    guard let self, self.generation == operationID else { return }
                    self.cancel(); self.endBackground()
                    self.message = "iOS ended the background session. Return to Crumb and retry."
                }
            }
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
        } catch {
            errorMessage = error.localizedDescription; busy = false
            if ["prepare", "install", "sign", "signing-info", "signing-plan"].contains(request["operation"] as? String ?? "") { cleanup() }
        }
    }
    private func receive(_ value: [String: Any]) {
        guard let phase = value["phase"] as? String else { return }
        self.phase = phase
        if let bytes = value["bytes"] as? Double, let total = value["total"] as? Double, total > 0 { progress = min(1,bytes / total) }
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
        guard DeviceJobState.terminalPhases.contains(phase) else { return }
        defer { drainResultRequest() }
        let request = lastRequest
        if let job { crumb_device_free(job); self.job = nil }
        busy = false; stopAdvertising(); endBackground(); pin = nil; lastRequest = [:]
        do {
            switch phase {
            case "installation-signing-ready":
                try DeviceKeychain.write(Data(installationPassword.utf8),account:"install-password")
                try DeviceKeychain.write(Data("validated".utf8),account:"install-signing-info")
                installationReady = true; cleanup()
                if installationSetupNeeded { installationSetupNeeded = false } else { continueQueued() }
            case "signing-plan":
                guard let record = signingRecord, let original = value["original_bundle_id"] as? String,
                      let bundle = value["bundle_id"] as? String, let staged else { throw DeviceOperationError.message("Invalid signing plan.") }
                let changes = value["adjusted_capabilities"] as? [String] ?? []
                if forceSigningReview || original != bundle || !changes.isEmpty {
                    signingReview = SigningReview(id:UUID(uuidString:staged.lastPathComponent) ?? UUID(),appName:record.appName,
                        originalBundleID:original,bundleID:bundle,bundleCount:value["bundle_count"] as? Int ?? 1,
                        changedIDs:value["changed_ids"] as? Int ?? 0,capabilities:changes,profileName:value["profile_name"] as? String ?? "")
                } else { approveSigning() }

            case "pairing-imported":
                guard let record = value["record"] as? String, let kind = value["pairing_kind"] as? String else { throw DeviceOperationError.message("Invalid pairing file.") }
                try DeviceKeychain.write(JSONSerialization.data(withJSONObject:["record":record,"pairing_kind":kind]),account:"pairing")
                paired = true; connected = false; connectionPolicy.becameInactive(); autoConnect()
            case "signing-ready":
                guard let bundle = value["helper_bundle"] as? String else { throw DeviceOperationError.message("Invalid signing identity.") }
                helperBundle = bundle
                try DeviceKeychain.write(Data(signingPassword.utf8), account:"password")
                try DeviceKeychain.write(JSONSerialization.data(withJSONObject:["helper_bundle":bundle]), account:"signing-info")
                signingReady = true; cleanup(); continueQueued()
            case "signed":
                guard var original = signingRecord, let bundle = value["bundle_id"] as? String,
                      let path = value["path"] as? String, let workspace = staged else { throw DeviceOperationError.message("Missing signed app.") }
                original.lastSigning = SigningSummary(bundleID:bundle,originalBundleID:value["original_bundle_id"] as? String ?? original.bundleID,
                    profileName:value["profile_name"] as? String ?? "",capabilities:value["adjusted_capabilities"] as? [String] ?? [],date:Date())
                try original.save()
                let record = DownloadRecord(filename:original.filename,appID:original.appID,appName:original.appName,bundleID:bundle,
                    version:original.version,build:original.build,externalVersionID:original.externalVersionID,date:original.date,decrypted:true,signed:true)
                let source = URL(fileURLWithPath:path)
                for file in try FileManager.default.contentsOfDirectory(at:workspace,includingPropertiesForKeys:nil) where file != source {
                    try? FileManager.default.removeItem(at:file)
                }
                staged = nil; signingRecord = nil; activeRecord = nil
                AppData.shared.restoreDownloadedIPA()
                signedInstallation = SignedInstallation(record:record,source:source,workspace:workspace)

            case "paired":
                guard let record = value["record"] as? String, let irk = value["host_irk"] as? String else { throw DeviceOperationError.message("Missing pairing credentials.") }
                try DeviceKeychain.write(JSONSerialization.data(withJSONObject: ["record":record, "host_irk":irk]), account: "pairing")
                paired = true; connectionPolicy.becameInactive(); autoConnect()
            case "connected": connected = true; continueQueued()
            case "prepared":
                guard let path = value["path"] as? String else { throw DeviceOperationError.message("Missing prepared helper.") }
                if let staged {
                    try? FileManager.default.removeItem(at: staged.appendingPathComponent("identity.p12"))
                    try? FileManager.default.removeItem(at: staged.appendingPathComponent("helper.mobileprovision"))
                }
                deviceOperation("install", extra: ["path":path, "helper_bundle":helperBundle, "replace_helper":managedHelper == helperBundle])
            case "installed":
                try DeviceKeychain.write(Data(helperBundle.utf8), account: "managed-helper")
                managedHelper = helperBundle
                helperInstalled = true; waitingForHelper = true
                if let pending { try DeviceKeychain.write(JSONSerialization.data(withJSONObject: pending), account: "pending") }
                cleanup(); message = "Helper installed. Keep Crumb Decrypt on screen until it finishes."
                openHelper()
            case "legacy-result", "complete":
                if phase == "legacy-result", let bundle = value["bundle_id"] as? String {
                    let original = DownloadRecord.load().first { $0.bundleID == bundle && $0.version == (value["version"] as? String) }
                    activeRecord = original ?? DownloadRecord(filename:"legacy.ipa",appID:"",appName:bundle,bundleID:bundle,version:value["version"] as? String ?? "",build:value["build"] as? String,externalVersionID:"",date:Date())
                }
                guard let path = value["path"] as? String else { throw DeviceOperationError.message("Missing verified result.") }
                guard let source = activeRecord else { throw DeviceOperationError.message("Missing download record.") }
                let name = UUID().uuidString + "-decrypted.ipa"
                let directory = FileManager.default.urls(for:.documentDirectory,in:.userDomainMask)[0].appendingPathComponent("Downloads",isDirectory:true)
                try FileManager.default.createDirectory(at:directory,withIntermediateDirectories:true)
                let destination = directory.appendingPathComponent(name)
                try FileManager.default.moveItem(at:URL(fileURLWithPath:path),to:destination)
                var record = DownloadRecord(filename:name,appID:source.appID,appName:source.appName,bundleID:source.bundleID,version:source.version,build:source.build,externalVersionID:source.externalVersionID,date:Date())
                record.decrypted = true
                do { try record.save() } catch { try? FileManager.default.removeItem(at:destination); throw error }
                output = destination; helperInstalled = false; waitingForHelper = false; activeRecord = nil; pending = nil
                try DeviceKeychain.remove("pending")
                AppData.shared.restoreDownloadedIPA()
                if installAfterDecrypt { install(record) }

            case "failed":
                // Local signing failures do not invalidate the verified route.
                let operation = request["operation"] as? String ?? ""
                if ["verify","install","fetch"].contains(operation) { connected = false }
                message = (value["message"] as? String) ?? "Could not finish."
                if operation != "verify" { errorMessage = message }
                if operation == "signing-plan", signingReview != nil { return }
                if operation == "signing-info", request["purpose"] as? String == "installation" {
                    installationReady = false; installationSetupNeeded = true
                }
                if operation == "verify", queuedRecord != nil { setupNeeded = true }
                if operation == "signing-info" && request["purpose"] as? String != "installation" { signingReady = false }
                if operation != "fetch" { activeRecord = nil }
                if operation == "sign" || operation == "signing-info" || operation == "signing-plan" { signingContext = nil; cleanup() }
                if ["prepare", "install", "sign", "signing-info", "signing-plan"].contains(request["operation"] as? String ?? "") { cleanup() }
            default:
                cleanup(); message = "Operation cancelled."
                if !waitingForHelper { activeRecord = nil }
                queuedRecord = nil; queuedInstall = nil
            }
        } catch { signingContext = nil; activeRecord = nil; errorMessage = error.localizedDescription; cleanup() }
    }
    private func stopAdvertising() { service?.stop(); service = nil }
    private func endBackground() { if background != .invalid { UIApplication.shared.endBackgroundTask(background); background = .invalid } }
    private func cleanup() { if let staged { try? FileManager.default.removeItem(at: staged) }; staged = nil }
    nonisolated func netService(_ sender: NetService, didNotPublish errorDict: [String: NSNumber]) {
        Task { @MainActor in self.cancel(); self.errorMessage = "Allow Local Network access in Settings, then pair again." }
    }
}
