import SwiftUI
import PartyUI

struct DownloadedAppsView: View {
    var embedded = false
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    @State private var showLogs = false
    @State private var installation: DownloadRecord?
    @State private var deletion: DownloadRecord?
    @State private var confirmDeletion = false
    @State private var deletionError = ""
    var body: some View {
        NavigationStack {
            List {
                if appData.isDowngrading {
                    Section {
                        VStack(alignment: .leading, spacing: 14) {
                            HStack(alignment: .top, spacing: 12) {
                                Image(systemName: "arrow.down.app").font(.title2).foregroundStyle(.tint)
                                VStack(alignment: .leading, spacing: 4) {
                                    Text(appData.activeDownloadName.isEmpty ? "Downloading app" : appData.activeDownloadName)
                                        .font(.headline)
                                    Text(appData.downloadPhaseLabel).font(.subheadline).foregroundStyle(.secondary)
                                }
                                Spacer(minLength: 0)
                                Button { appData.storeTask?.cancel() } label: {
                                    Image(systemName: "xmark.circle.fill").font(.title2).foregroundStyle(.secondary)
                                        .frame(minWidth: 44, minHeight: 44)
                                }
                                .buttonStyle(.borderless).accessibilityLabel("Cancel download")
                            }
                            if appData.showsDowngradeProgress {
                                ProgressView(value: appData.downgradeProgress)
                                HStack {
                                    Text(appData.downloadProgressText)
                                    Spacer()
                                    Text("\(Int(appData.downgradeProgress * 100))%")
                                        .monospacedDigit()
                                }
                                .font(.caption).foregroundStyle(.secondary)
                            }
                        }
                        .padding(.vertical, 8)
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
                    Section {
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
                        // A destructive swipe optimistically removes the row before confirmation.
                        Button { requestDeletion(record) } label: { Label("Delete", systemImage: "trash") }
                            .tint(.red)
                    }
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
            }
            .sheet(item: $installation) { OTAInstallationView(record: $0, startImmediately: true) }
            .alert("Delete downloaded IPA?", isPresented: $confirmDeletion) {
                if let record = deletion {
                    Button("Delete IPA", role: .destructive) {
                        do {
                            try withAnimation(reduceMotion ? nil : .easeInOut(duration: 0.25)) {
                                try appData.deleteDownload(record)
                            }
                        }
                        catch {
                            appData.restoreDownloadedIPA()
                            deletionError = "Could not delete the download (code \((error as NSError).code))."
                        }
                        deletion = nil
                    }
                }
                Button("Cancel", role: .cancel) { deletion = nil }
            } message: {
                Text("The saved IPA will be removed. The installed app and its data will stay on this device.")
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
