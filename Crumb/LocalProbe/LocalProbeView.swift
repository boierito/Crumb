import SwiftUI

@MainActor
final class LocalProbeCoordinator: ObservableObject {
    static let shared = LocalProbeCoordinator()
    @Published private(set) var busy = false
    @Published private(set) var stage = ""
    @Published private(set) var outcome = ""
    @Published private(set) var reportURL: URL?
    @Published private(set) var outputURL: URL?
    private let worker = DispatchQueue(label: "Crumb.LocalProbe", qos: .userInitiated)
    private var stopRequested = false
    private let runner: CrumbLocalProbe
    private let directory: URL

    private init() {
        directory = FileManager.default.urls(for: .applicationSupportDirectory, in: .userDomainMask)[0]
            .appendingPathComponent("LocalProbe", isDirectory: true)
        runner = CrumbLocalProbe(directory: directory)
        // The native runner classifies any interrupted job before loading its report.
        let files = (try? FileManager.default.contentsOfDirectory(at: directory,
            includingPropertiesForKeys: [.contentModificationDateKey])) ?? []
        reportURL = files.filter { $0.lastPathComponent.hasPrefix("report-") && $0.pathExtension == "json" }
            .sorted { ((try? $0.resourceValues(forKeys: [.contentModificationDateKey]).contentModificationDate) ?? .distantPast) >
                      ((try? $1.resourceValues(forKeys: [.contentModificationDateKey]).contentModificationDate) ?? .distantPast) }.first
        if let reportURL, let data = try? Data(contentsOf: reportURL),
           let report = try? JSONSerialization.jsonObject(with: data) as? [String: Any] {
            outcome = Self.message(report["outcome"] as? String ?? "")
        }
    }
    func run(_ record: DownloadRecord) {
        guard !busy, let ipa = record.fileURL else { return }
        stopRequested = false
        busy = true; stage = "Preparing test".localized; outcome = ""; reportURL = nil; outputURL = nil
        let runner = self.runner
        worker.async {
            let result = runner.runIPA(ipa, bundleID: record.bundleID) { stage in
                Task { @MainActor in self.stage = Self.progress(stage) }
            }
            let generated = runner.outputURL
            let report = runner.reportURL
            Task { @MainActor in
                self.reportURL = report
                self.outcome = Self.message(result["outcome"] as? String ?? "")
                if self.stopRequested {
                    if let generated { try? FileManager.default.removeItem(at: generated) }
                    self.outcome = Self.message("cancelled_or_local_failure")
                    if let report, let data = try? Data(contentsOf: report),
                       var value = try? JSONSerialization.jsonObject(with: data) as? [String: Any] {
                        value["outcome"] = "cancelled_or_local_failure"
                        if let updated = try? JSONSerialization.data(withJSONObject: value, options: [.prettyPrinted, .sortedKeys]) {
                            try? updated.write(to: report, options: [.atomic, .completeFileProtection])
                        }
                    }
                    self.busy = false; self.stage = ""; return
                }
                if let generated {
                    do {
                        let downloads = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask)[0]
                            .appendingPathComponent("Downloads", isDirectory: true)
                        try FileManager.default.createDirectory(at: downloads, withIntermediateDirectories: true)
                        let destination = downloads.appendingPathComponent(UUID().uuidString + "-decrypted.ipa")
                        try FileManager.default.moveItem(at: generated, to: destination)
                        var saved = DownloadRecord(filename: destination.lastPathComponent, appID: record.appID,
                            appName: record.appName, bundleID: record.bundleID, version: record.version,
                            build: record.build, externalVersionID: record.externalVersionID, date: Date())
                        saved.decrypted = true
                        do { try saved.save() }
                        catch { try? FileManager.default.moveItem(at: destination, to: generated); throw error }
                        self.outputURL = destination
                        AppData.shared.restoreDownloadedIPA()
                    } catch {
                        self.outputURL = generated
                        self.outcome += " The verified IPA is available to export here.".localized
                    }
                }
                self.busy = false; self.stage = ""
            }
        }
    }
    func stop() { stopRequested = true; runner.cancel(); stage = "Stopping after the current operation".localized }
    private static func progress(_ stage: String) -> String {
        switch stage {
        case "archive_preparation": return "Checking the downloaded app".localized
        case "installed_discovery": return "Looking for the installed version".localized
        case "sandbox_extension": return "Checking file access".localized
        case "source_hash", "inspect", "parse": return "Checking the original executable".localized
        case "copy", "shared_container_copy": return "Preparing a local copy".localized
        case "remap", "signature", "page_read": return "Testing decryption".localized
        case "archive_verification", "verify_pages", "patch_output": return "Verifying the result".localized
        default: return "Running local checks".localized
        }
    }
    private static func message(_ code: String) -> String {
        switch code {
        case "downloaded_ipa_success": return "Downloaded IPA decrypted and verified without a helper.".localized
        case "installed_app_success": return "Decryption verified using the installed app. This route requires that installed version.".localized
        case "no_verified_route": return "No working decryption route was established. Share the report to identify which checks failed.".localized
        case "preparation_failed": return "The IPA could not be prepared for this test. The report includes the failed stage.".localized
        case "interrupted": return "The previous test was interrupted. Its report preserves the last recorded operation.".localized
        case "cancelled_or_local_failure": return "The test stopped before a complete result. The report preserves completed checks.".localized
        default: return code.isEmpty ? "" : "The test could not finish. Share its report for details.".localized
        }
    }
}

struct LocalProbeView: View {
    let record: DownloadRecord
    @ObservedObject private var probe = LocalProbeCoordinator.shared
    @Environment(\.dismiss) private var dismiss
    var body: some View {
        NavigationStack {
            List {
                Section {
                    Text(record.appName).font(.headline)
                    Text("Tests decryption inside Crumb, without another app or a VPN. If the same version is installed, it is also checked automatically.".localized)
                    Text("Keep Crumb open while the test runs. Signing into Crumb does not establish the iPhone’s App Store license.".localized)
                        .font(.footnote).foregroundStyle(.secondary)
                    Button("Try without helper".localized) { probe.run(record) }.disabled(probe.busy)
                }
                if probe.busy {
                    Section {
                        HStack { ProgressView(); Text(probe.stage) }
                        Button("Stop".localized, role: .destructive) { probe.stop() }
                    }
                }
                if !probe.outcome.isEmpty { Section { Text(probe.outcome) } }
                if !probe.busy, let report = probe.reportURL {
                    Section { ShareLink(item: report) { Label("Share test report".localized, systemImage: "square.and.arrow.up") } }
                }
                if !probe.busy, let ipa = probe.outputURL {
                    Section { ShareLink(item: ipa) { Label("Export verified IPA".localized, systemImage: "square.and.arrow.up") } }
                }
            }
            .navigationTitle("Without helper".localized)
            .toolbar { ToolbarItem(placement: .confirmationAction) { Button("Done".localized) { dismiss() } } }
        }
        .interactiveDismissDisabled(probe.busy)
    }
}
