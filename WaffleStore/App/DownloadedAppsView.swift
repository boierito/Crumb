import SwiftUI
import PartyUI

struct DownloadedAppsView: View {
    var embedded = false
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @State private var showHistory = false
    @State private var showLogs = false
    @State private var installation: DownloadRecord?
    @State private var deletion: DownloadRecord?
    @State private var confirmDeletion = false
    @State private var deletionError = ""
    var body: some View {
        NavigationStack {
            List {
                if appData.isDowngrading {
                    Section("Downloading") {
                        HStack(spacing: 12) { ProgressView(); Text(appData.applicationStatus) }
                        if appData.showsDowngradeProgress {
                            ProgressView(value: appData.downgradeProgress)
                            Text(appData.downgradeProgressDetail).font(.caption).foregroundStyle(.secondary)
                        }
                        Button("Cancel download", role: .destructive) { appData.storeTask?.cancel() }
                    }
                }

                if appData.completedDownloads.isEmpty {
                    VStack(alignment: .leading, spacing: 8) {
                        Label("No downloads yet", systemImage: "square.and.arrow.down")
                            .font(.headline)
                        Text("Choose an app and download a version. Its IPA will appear here for installation or export.")
                            .foregroundStyle(.secondary)
                    }
                }
                ForEach(appData.completedDownloads) { record in
                    VStack(alignment: .leading, spacing: 12) {
                        HStack(alignment: .top, spacing: 12) {
                            VStack(alignment: .leading, spacing: 4) {
                                Text(record.appName).font(.headline)
                                Text("Version \(record.version)").font(.subheadline).foregroundStyle(.secondary)
                            }
                            Spacer(minLength: 0)
                            Menu {
                                Button(role: .destructive) { requestDeletion(record) } label: {
                                    Label("Delete IPA", systemImage: "trash")
                                }
                            } label: {
                                Image(systemName: "ellipsis").frame(minWidth: 44, minHeight: 44)
                            }
                            .buttonStyle(.borderless)
                            .accessibilityLabel("Options for \(record.appName), version \(record.version)")
                        }
                        HStack(spacing: 10) {
                            Button { installation = record } label: { Label("Install", systemImage: "arrow.down.app") }
                            if let url = record.fileURL {
                                ShareLink(item: url) { Label("Export", systemImage: "square.and.arrow.up") }
                            }
                        }
                        .buttonStyle(TranslucentButtonStyle())
                        DisclosureGroup("Details") {
                            if let build = record.build { Text("Build \(build)").font(.caption) }
                            Text(record.bundleID).font(.caption).textSelection(.enabled)
                            Text("externalVersionId \(record.externalVersionID)").font(.caption).textSelection(.enabled)
                            Text(record.date, style: .date).font(.caption)
                            if let url = record.fileURL,
                               let size = try? url.resourceValues(forKeys: [.fileSizeKey]).fileSize {
                                Text(ByteCountFormatter.string(fromByteCount: Int64(size), countStyle: .file)).font(.caption)
                            }
                            Text("Version and build read from the downloaded IPA.").font(.caption).foregroundStyle(.secondary)
                        }
                    }
                    .padding(.vertical, 4)
                    .swipeActions(allowsFullSwipe: false) {
                        Button(role: .destructive) { requestDeletion(record) } label: { Label("Delete", systemImage: "trash") }
                    }
                }
                Section {
                    DisclosureGroup("Activity log", isExpanded: $showLogs) { LogView() }
                }
            }
            .navigationTitle("Downloaded apps")
            .navigationBarTitleDisplayMode(.inline)
            .listStyle(.insetGrouped)
             .toolbar {
                if !embedded { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
                ToolbarItem(placement: .topBarTrailing) {
                    Button { showHistory = true } label: { Image(systemName: "clock.arrow.circlepath") }
                        .accessibilityLabel("Downgrade history")
                }
            }
            .sheet(isPresented: $showHistory) { DowngradeHistoryView() }
            .sheet(item: $installation) { OTAInstallationView(record: $0, startImmediately: true) }
            .confirmationDialog("Delete downloaded IPA?", isPresented: $confirmDeletion, titleVisibility: .visible) {
                if let record = deletion {
                    Button("Delete IPA", role: .destructive) {
                        do { try appData.deleteDownload(record) }
                        catch {
                            appData.restoreDownloadedIPA()
                            deletionError = "Could not delete the download (code \((error as NSError).code))."
                        }
                        deletion = nil
                    }
                }
                Button("Cancel", role: .cancel) { deletion = nil }
            } message: {
                Text("This removes the saved IPA and its download record. It does not uninstall the app or delete its data.")
            }
            .alert("Could not delete download", isPresented: Binding(
                get: { !deletionError.isEmpty }, set: { if !$0 { deletionError = "" } })) {
                    Button("OK") { deletionError = "" }
                } message: { Text(deletionError) }
            .onAppear { appData.restoreDownloadedIPA() }
        }
    }
    private func requestDeletion(_ record: DownloadRecord) {
        deletion = record; confirmDeletion = true
    }
}
