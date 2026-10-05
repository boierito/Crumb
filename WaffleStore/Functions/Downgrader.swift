import Foundation
import SwiftUI
import PartyUI
import MapleSAP

// Download/export comes first. No localhost installer, third-party manifest,
// or fabricated installation-success history is used by this jailed flow.
func downgradeApp(appId: String, ipaTool: IPATool) { AppData.shared.showStoreVersions = true }
func cleanUp() {
    // Original auto-clean setting must not erase completed Downloads at launch.
    let fm = FileManager.default
    for name in ["app.ipa", "app"] {
        try? fm.removeItem(at: fm.temporaryDirectory.appendingPathComponent(name))
    }
}
func resetDowngradeProgress() {
    let appData = AppData.shared
    appData.isDowngrading = false; appData.showsDowngradeProgress = false
}
extension AppData {
    func download(app: StoreApp, version: String, tool: IPATool, expectedVersion: String? = nil) {
        guard storeTask == nil else { return }
        isDowngrading = true; showsDowngradeProgress = true; downgradeProgress = 0
        storeError = ""

        storeTask = Task {
            defer { storeTask = nil; isDowngrading = false; showsDowngradeProgress = false }
            let fm = FileManager.default
            let scratch = fm.temporaryDirectory.appendingPathComponent("store-" + UUID().uuidString, isDirectory: true)
            defer { try? fm.removeItem(at: scratch) }
            do {
                try fm.createDirectory(at: scratch, withIntermediateDirectories: true)
                storeStage("Requesting selected externalVersionId")
                let descriptor = try await tool.descriptor(app: app, version: version)
                try Task.checkCancellation()
                storeStage("Downloading IPA from CDN")
                let source = scratch.appendingPathComponent("source.ipa")
                var received: URL?
                for attempt in 0..<3 {
                    do {
                        let downloader = CDNDownload(destination: source) { progress, written, expected in
                            Task { @MainActor in
                                guard self.isDowngrading, self.downgradeProgress < 0.85 else { return }
                                self.downgradeProgress = progress * 0.85
                                self.downgradeProgressDetail = expected > 0 ? "\(written / 1_048_576) / \(expected / 1_048_576) MiB" : "\(written / 1_048_576) MiB"
                            }
                        }
                        received = try await downloader.run(descriptor.url)
                        break
                    } catch {
                        try Task.checkCancellation()
                        let retryable: Bool
                        if let error = error as? URLError { retryable = [.timedOut, .networkConnectionLost, .cannotConnectToHost].contains(error.code) }
                        else if let error = error as? CDNHTTPFailure { retryable = [429, 500, 502, 503, 504].contains(error.status) }
                        else { retryable = false }
                        guard retryable, attempt < 2 else { throw error }
                        try? fm.removeItem(at: source)
                        storeStage("Retrying CDN download")
                        let delay = try StoreRetry.delay((error as? CDNHTTPFailure)?.retryAfter, attempt: attempt)
                        try await Task.sleep(nanoseconds: UInt64(delay * 1_000_000_000))
                    }
                }
                guard let source = received else { throw StoreError.packageInvalid }
                storeStage("Validating ZIP, app identity and purchase data")
                downgradeProgress = 0.85
                let staged = scratch.appendingPathComponent("prepared.ipa")
                // C guest/package work is blocking; keep it off MainActor.
                // Cancellation prevents committing the file after it returns.
                let info = try await Task.detached(priority: .userInitiated) {
                    try NativePackage.prepare(source: source, destination: staged, app: app, descriptor: descriptor)
                }.value
                try Task.checkCancellation()
                if let expectedVersion, info.version != expectedVersion { throw StoreError.versionMismatch }
                let downloads = try fm.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: true)
                    .appendingPathComponent("Downloads", isDirectory: true)
                try fm.createDirectory(at: downloads, withIntermediateDirectories: true)
                let filename = "\(app.id)-\(descriptor.externalVersionID)-\(UUID().uuidString).ipa"
                let destination = downloads.appendingPathComponent(filename)
                let record = DownloadRecord(filename: filename, appID: app.id, appName: app.name, bundleID: info.bundleID,
                    version: info.version, build: info.build, externalVersionID: descriptor.externalVersionID, date: Date())
                let sidecar = destination.deletingPathExtension().appendingPathExtension("json")
                try JSONEncoder().encode(record).write(to: sidecar, options: .atomic)
                do { try fm.moveItem(at: staged, to: destination) }
                catch { try? fm.removeItem(at: sidecar); throw error }
                completedDownloads = DownloadRecord.load()
                downloadedIPAURL = destination; hasAppBeenServed = true
                appBundleID = info.bundleID; appVersion = info.version
                storeStage("IPA verified and saved. Export to Files or another app.")

                downgradeProgress = 1; applicationIcon = "checkmark.circle.fill"
            } catch {
                if error is CancellationError || Task.isCancelled {
                    applicationStatus = "Download cancelled."
                } else {
                    storeError = (error as? CDNHTTPFailure).map { "CDN request failed (HTTP \($0.status))." } ?? (error as? StoreError)?.localizedDescription ?? (error as? SAPError)?.localizedDescription ?? "Store operation failed (code \((error as NSError).code))."
                    applicationStatus = "Download failed."; applicationIcon = "xmark.circle.fill"
                    Alertinator.shared.alert(title: "Download failed", body: storeError)
                }
            }
        }
    }
    private func storeStage(_ stage: String) {
        applicationStatus = stage; downgradeProgressDetail = stage

        print("Apple Store stage: \(stage)")
    }
    func restoreDownloadedIPA() {
        completedDownloads = DownloadRecord.load()
        let fm = FileManager.default
        guard let docs = try? fm.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: false),
              let files = try? fm.contentsOfDirectory(at: docs.appendingPathComponent("Downloads"),
                  includingPropertiesForKeys: [.contentModificationDateKey]) else { return }
        downloadedIPAURL = files.filter { $0.pathExtension == "ipa" }.sorted {
            ((try? $0.resourceValues(forKeys: [.contentModificationDateKey]).contentModificationDate) ?? .distantPast) >
            ((try? $1.resourceValues(forKeys: [.contentModificationDateKey]).contentModificationDate) ?? .distantPast)
        }.first
        hasAppBeenServed = downloadedIPAURL != nil
        if let record = completedDownloads.first { appBundleID = record.bundleID; appVersion = record.version }
    }
}
