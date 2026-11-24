import Foundation
import Combine
import UIKit
import SwiftUI

enum DeviceOperationError: LocalizedError {
    case message(String)
    var errorDescription: String? { if case let .message(text) = self { return text }; return nil }
}

@MainActor
final class DeviceWorkflow: NSObject, ObservableObject {
    static let shared = DeviceWorkflow()
    private var groupWorker: Task<URL, Error>?
    @Published private(set) var phase = "idle"
    @Published private(set) var message = ""
    @Published private(set) var busy = false
    @Published private(set) var hasCertificate = false
    @Published private(set) var hasProfile = false
    @Published var setupVisible = false
    @Published var errorMessage = ""
    @Published private(set) var certificateName = ""
    @Published private(set) var profileName = ""
    @Published var signingPassword = ""
    @Published private(set) var progress: Double?
    @Published private(set) var activeRecord: DownloadRecord?
    private var signingRecord: DownloadRecord?
    private var signingContext: [String: Any]?
    private var forceSigningReview = false
    private var signingApprovalPending = false
    @Published var installationSetupNeeded = false
    @Published private(set) var installationReady = false
    @Published private(set) var installationProfiles: [AdditionalSigningProfile] = []
    @Published var signingReview: SigningReview?
    @Published var signedInstallation: SignedInstallation?
    private var queuedInstall: DownloadRecord?
    // Preserve legacy identifiers as protected signing targets after migration.
    private var helperBundle = "com.certlium.crumb.decrypt"
    private var job: UnsafeMutableRawPointer?
    private var generation = UUID()
    private var polling: Task<Void, Never>?
    private var background: UIBackgroundTaskIdentifier = .invalid
    private var staged: URL?
    private var managedHelper: String?
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
            hasCertificate = try DeviceKeychain.read("p12") != nil
            hasProfile = try DeviceKeychain.read("profile") != nil
            signingPassword = String(data: try DeviceKeychain.read("password") ?? Data(), encoding: .utf8) ?? ""
            certificateName = String(data: try DeviceKeychain.read("p12-name") ?? Data(), encoding: .utf8) ?? ""
            profileName = String(data: try DeviceKeychain.read("profile-name") ?? Data(), encoding: .utf8) ?? ""
            if let data = try DeviceKeychain.read("signing-info"),
               let info = try JSONSerialization.jsonObject(with: data) as? [String: String], hasCertificate && hasProfile {
                helperBundle = info["helper_bundle"] ?? helperBundle
            }
            // Remove stale state from versions that used an external helper.
            try? DeviceKeychain.remove("pending")
        } catch { message = "Could not read protected device settings." }
        loadSharedSigning()
    }
    func importFile(_ url: URL, kind: String) {
        guard !busy else { return }
        let access = url.startAccessingSecurityScopedResource(); defer { if access { url.stopAccessingSecurityScopedResource() } }
        do {
            let size = try url.resourceValues(forKeys: [.fileSizeKey]).fileSize ?? 0
            guard size > 0 && size <= 2 * 1024 * 1024 else { throw DeviceOperationError.message("Signing file must be between 1 byte and 2 MiB.") }
            let data = try Data(contentsOf: url)
            if kind == "install-extra-profile" {
                guard installationProfiles.count < 32 else { throw DeviceOperationError.message("Maximum 32 extension profiles.") }
                installationProfiles.append(AdditionalSigningProfile(id:UUID().uuidString,name:url.lastPathComponent,data:data))
                try DeviceKeychain.write(JSONEncoder().encode(installationProfiles),account:"install-profiles")
                return
            }
            try DeviceKeychain.write(data, account: kind)
            try DeviceKeychain.write(Data(url.lastPathComponent.utf8), account: kind + "-name")
            invalidateSigning()
            if kind == "p12" { certificateName = url.lastPathComponent } else { profileName = url.lastPathComponent }
            if kind == "p12" { hasCertificate = true } else { hasProfile = true }
            message = "Signing file stored on this device."
        } catch { errorMessage = error.localizedDescription }
    }
    var activityTitle: String {
        switch phase {
        case "group-decrypt": return "Decrypting on this device"
        case "signing": return "Signing app"
        default: return "Preparing"
        }
    }
    private func signingRequest(operation: String) throws -> [String: Any] {
        let fm = FileManager.default
        let base = try fm.url(for: .applicationSupportDirectory, in: .userDomainMask, appropriateFor: nil, create: true).appendingPathComponent("DeviceJobs", isDirectory: true)
        try fm.createDirectory(at: base, withIntermediateDirectories: true, attributes: [.protectionKey:FileProtectionType.complete])
        let root = base.appendingPathComponent(UUID().uuidString, isDirectory:true)
        try fm.createDirectory(at: root, withIntermediateDirectories:false, attributes:[.protectionKey:FileProtectionType.complete])
        var excluded = URLResourceValues(); excluded.isExcludedFromBackup = true
        var protectedRoot = root; try protectedRoot.setResourceValues(excluded)
        staged = root
        guard let cert = try DeviceKeychain.read("p12"), let profile = try DeviceKeychain.read("profile") else { throw DeviceOperationError.message("Import your signing files.") }
        let p12 = root.appendingPathComponent("identity.p12"), mobile = root.appendingPathComponent("profile.mobileprovision")
        try cert.write(to:p12,options:[.atomic,.completeFileProtection]); try profile.write(to:mobile,options:[.atomic,.completeFileProtection])
        var request: [String:Any] = ["operation":operation,"p12":p12.path,"profile":mobile.path,
            "password":signingPassword,
            "helper_bundle":helperBundle,
            "protected_bundle_ids":[helperBundle,managedHelper].compactMap { $0 },
            "own_bundle":Bundle.main.bundleIdentifier ?? "com.certlium.crumb", "work":root.appendingPathComponent("work").path]
        var profiles: [String] = []
        for profile in installationProfiles {
            let url = root.appendingPathComponent(profile.id + ".mobileprovision")
            try profile.data.write(to: url, options: [.atomic, .completeFileProtection])
            profiles.append(url.path)
        }
        request["profiles"] = profiles
        return request
    }
    func loadSharedSigning() {
        do {
            // Collapse the former installation-only store into the existing identity.
            // A complete existing helper identity always wins; never mix key pairs.
            if try DeviceKeychain.read("p12") == nil, try DeviceKeychain.read("profile") == nil,
               let cert = try DeviceKeychain.read("install-p12"), let profile = try DeviceKeychain.read("install-profile") {
                try DeviceKeychain.write(cert,account:"p12")
                try DeviceKeychain.write(profile,account:"profile")
                for key in ["password","p12-name","profile-name"] {
                    if let data = try DeviceKeychain.read("install-" + key) { try DeviceKeychain.write(data,account:key) }
                }
                signingPassword = String(data:try DeviceKeychain.read("password") ?? Data(),encoding:.utf8) ?? ""
            }
            hasCertificate = try DeviceKeychain.read("p12") != nil
            hasProfile = try DeviceKeychain.read("profile") != nil
            certificateName = String(data:try DeviceKeychain.read("p12-name") ?? Data(),encoding:.utf8) ?? ""
            profileName = String(data:try DeviceKeychain.read("profile-name") ?? Data(),encoding:.utf8) ?? ""
            let validated = try DeviceKeychain.read("identity-validated") != nil || DeviceKeychain.read("signing-info") != nil
            installationReady = hasCertificate && hasProfile && validated
            if let data = try DeviceKeychain.read("install-profiles") { installationProfiles = try JSONDecoder().decode([AdditionalSigningProfile].self,from:data) }
        } catch { errorMessage = error.localizedDescription }
    }
    func removeInstallationProfile(_ id:String) {
        guard !busy else { return }
        installationProfiles.removeAll { $0.id == id }
        do { try DeviceKeychain.write(JSONEncoder().encode(installationProfiles),account:"install-profiles") }
        catch { errorMessage = error.localizedDescription }
    }
    func validateSigningForInstall() {
        guard !busy else { return }
        do { start(try signingRequest(operation:"signing-info")) }
        catch { cleanup(); errorMessage = error.localizedDescription }
    }
    func invalidateSigning() {
        installationReady = false
        try? DeviceKeychain.remove("signing-info")
        try? DeviceKeychain.remove("identity-validated")
    }
    func saveSigning() {
        validateSigningForInstall()
    }
    func decrypt(_ record: DownloadRecord, install: Bool = false) {
        guard !busy, groupWorker == nil, signingContext == nil, signedInstallation == nil else {
            errorMessage = "Finish the current app first."; return
        }
        guard let root = AppGroupDecryption.container() else {
            errorMessage = "Crumb could not access its shared storage. Re-sign Crumb in KSign with an authorized App Group, then try again."
            return
        }
        guard let source = record.fileURL else { errorMessage = "The downloaded IPA is missing."; return }
        crumb_group_reset()
        busy = true; activeRecord = record; phase = "group-decrypt"; progress = 0
        message = "Decrypting on this device…"; errorMessage = ""
        let token = UUID(); generation = token
        background = UIApplication.shared.beginBackgroundTask(withName: "Crumb App Group decryption") { [weak self] in
            Task { @MainActor in self?.cancel() }
        }
        let worker = Task.detached(priority: .userInitiated) {
            try AppGroupDecryption.run(source: source, root: root, cancelled: { Task.isCancelled }) { value in
                Task { @MainActor [weak self] in
                    guard let self, self.generation == token, self.groupWorker != nil else { return }
                    self.progress = value
                }
            }
        }
        groupWorker = worker
        Task {
            do {
                let destination = try await worker.value
                guard !worker.isCancelled else { try? FileManager.default.removeItem(at: destination); throw CancellationError() }
                var result = DownloadRecord(filename: destination.lastPathComponent, appID: record.appID,
                    appName: record.appName, bundleID: record.bundleID, version: record.version,
                    build: record.build, externalVersionID: record.externalVersionID, date: Date())
                result.decrypted = true
                do { try result.save() } catch { try? FileManager.default.removeItem(at: destination); throw error }
                groupWorker = nil; busy = false; activeRecord = nil; progress = nil; phase = "idle"; endBackground()
                message = "Decrypted and verified. Ready to export or sign and install."
                AppData.shared.restoreDownloadedIPA()
                if install { self.install(result) }
            } catch {
                groupWorker = nil; busy = false; activeRecord = nil; progress = nil; phase = "idle"; endBackground()
                if worker.isCancelled || error is CancellationError { message = "Decryption stopped. Original IPA retained." }
                else { errorMessage = error.localizedDescription; message = "Decryption failed. Original IPA retained." }
            }
        }
    }
    private func continueQueuedInstallation() {
        guard installationReady, !busy, let record = queuedInstall else { return }
        queuedInstall = nil
        installationSetupNeeded = false
        install(record, review: forceSigningReview)
    }
    func install(_ record: DownloadRecord, review: Bool = false) {
        guard !busy, signingContext == nil, signedInstallation == nil else { errorMessage = "Finish the current app first."; return }
        guard record.isDecrypted else { AppData.shared.installationRequest = record; return }
        forceSigningReview = review
        loadSharedSigning()
        guard installationReady else {
            queuedInstall = record
            if hasCertificate && hasProfile { activeRecord = record; validateSigningForInstall() }
            else { installationSetupNeeded = true }
            return
        }
        do {
            guard let url = record.fileURL else { throw CocoaError(.fileNoSuchFile) }
            var request = try signingRequest(operation:"signing-plan")
            request["ipa"] = url.path; request["output"] = staged!.appendingPathComponent("signed.ipa").path
            signingRecord = record; activeRecord = record; signingContext = request
            start(request)
        } catch { cleanup(); activeRecord = nil; errorMessage = error.localizedDescription }
    }
    func installationSetupDismissed() {
        if installationReady { continueQueuedInstallation() }
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

    func cancel() {
        if let groupWorker { groupWorker.cancel(); crumb_group_cancel(); message = "Cancelling…"; return }
        guard let job else { return }
        crumb_device_cancel(job)
        message = "Cancelling…"
    }

    // A cancelled native worker owns its workspace until it has stopped.
    func stopDecryption() {
        if groupWorker != nil { cancel(); return }
        let retiredJob = job
        let retiredWorkspace = staged
        generation = UUID()
        polling?.cancel(); polling = nil; job = nil; staged = nil
        if let retiredJob { crumb_device_cancel(retiredJob) }
        queuedInstall = nil; signingRecord = nil
        signingContext = nil; signingApprovalPending = false; signingReview = nil
        activeRecord = nil
        busy = false; progress = nil; phase = "idle"; message = ""; errorMessage = ""
        lastRequest = [:]
        endBackground()
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
    }
    private func start(_ request: [String: Any]) {
        guard !busy else { return }
        do {
            let data = try JSONSerialization.data(withJSONObject: request)
            guard let json = String(data: data, encoding: .utf8), let handle = crumb_device_start(json) else { throw DeviceOperationError.message("Could not start device operation.") }
            job = handle; progress = nil; busy = true; phase = "starting"; message = ""
            lastRequest = request.filter { $0.key != "password" }
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
            if ["sign", "signing-info", "signing-plan"].contains(request["operation"] as? String ?? "") { cleanup() }
        }
    }
    private func receive(_ value: [String: Any]) {
        guard let phase = value["phase"] as? String else { return }
        self.phase = phase
        if let bytes = value["bytes"] as? Double, let total = value["total"] as? Double, total > 0 { progress = min(1,bytes / total) }
        guard DeviceJobState.terminalPhases.contains(phase) else { return }
        let request = lastRequest
        if let job { crumb_device_free(job); self.job = nil }
        busy = false; endBackground(); lastRequest = [:]
        do {
            switch phase {
            case "installation-signing-ready":
                try DeviceKeychain.write(Data(signingPassword.utf8),account:"password")
                try DeviceKeychain.write(Data("validated".utf8),account:"identity-validated")
                installationReady = true; cleanup()
                if installationSetupNeeded { installationSetupNeeded = false } else { continueQueuedInstallation() }
            case "signing-plan":
                guard let record = signingRecord, let original = value["original_bundle_id"] as? String,
                      let bundle = value["bundle_id"] as? String, let staged else { throw DeviceOperationError.message("Invalid signing plan.") }
                let changes = value["adjusted_capabilities"] as? [String] ?? []
                if forceSigningReview || original != bundle || !changes.isEmpty {
                    signingReview = SigningReview(id:UUID(uuidString:staged.lastPathComponent) ?? UUID(),appName:record.appName,
                        originalBundleID:original,bundleID:bundle,bundleCount:value["bundle_count"] as? Int ?? 1,
                        changedIDs:value["changed_ids"] as? Int ?? 0,capabilities:changes,profileName:value["profile_name"] as? String ?? "")
                } else { approveSigning() }

            case "signed":
                guard var original = signingRecord, let bundle = value["bundle_id"] as? String,
                      let path = value["path"] as? String, let workspace = staged else { throw DeviceOperationError.message("Missing signed app.") }
                original.lastSigning = SigningSummary(bundleID:bundle,originalBundleID:value["original_bundle_id"] as? String ?? original.bundleID,
                    profileName:value["profile_name"] as? String ?? "",capabilities:value["adjusted_capabilities"] as? [String] ?? [],date:Date(),
                    teamID:value["signing_team"] as? String,applicationIdentifier:value["application_identifier"] as? String,
                    originalKeychainGroups:value["original_keychain_groups"] as? [String],keychainGroups:value["keychain_groups"] as? [String],
                    originalAppGroups:value["original_app_groups"] as? [String],appGroups:value["app_groups"] as? [String])
                if let diagnostics = value["cloudkit_diagnostics"], JSONSerialization.isValidJSONObject(diagnostics) {
                    let data = try JSONSerialization.data(withJSONObject:diagnostics,options:[.prettyPrinted,.sortedKeys])
                    original.lastSigning?.cloudKitDiagnostics = String(data:data,encoding:.utf8)
                }
                try original.save()
                let record = DownloadRecord(filename:original.filename,appID:original.appID,appName:original.appName,bundleID:bundle,
                    version:original.version,build:original.build,externalVersionID:original.externalVersionID,date:original.date,decrypted:true,signed:true)
                let source = URL(fileURLWithPath:path)
                for file in try FileManager.default.contentsOfDirectory(at:workspace,includingPropertiesForKeys:nil) where file.lastPathComponent != source.lastPathComponent {
                    try? FileManager.default.removeItem(at:file)
                }
                staged = nil; signingRecord = nil; activeRecord = nil
                AppData.shared.restoreDownloadedIPA()
                signedInstallation = SignedInstallation(record:record,source:source,workspace:workspace)

            case "failed":
                let operation = request["operation"] as? String ?? ""
                message = (value["message"] as? String) ?? "Could not finish."
                errorMessage = message
                if operation == "signing-plan", signingReview != nil { return }
                if operation == "signing-info" {
                    installationReady = false
                    if queuedInstall != nil && !setupVisible { installationSetupNeeded = true }
                }
                activeRecord = nil
                signingContext = nil; cleanup()
            default:
                cleanup(); message = "Operation cancelled."
                activeRecord = nil
                queuedInstall = nil
            }
        } catch { signingContext = nil; activeRecord = nil; errorMessage = error.localizedDescription; cleanup() }
    }
    private func endBackground() { if background != .invalid { UIApplication.shared.endBackgroundTask(background); background = .invalid } }
    private func cleanup() { if let staged { try? FileManager.default.removeItem(at: staged) }; staged = nil }
}
