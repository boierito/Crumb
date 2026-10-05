import SwiftUI

struct DownloadedAppsView: View {
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @State private var installation: DownloadRecord?
    @State private var deletion: DownloadRecord?
    @State private var confirmDeletion = false
    @State private var deletionError = ""
    var body: some View {
        NavigationStack {
            List {
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
                        Text(record.appName).font(.headline)
                        Text("Version \(record.version)").font(.subheadline)
                        if let build = record.build { Text("Build \(build)").font(.caption).foregroundStyle(.secondary) }
                        HStack {
                            Button { installation = record } label: { Label("Install", systemImage: "arrow.down.app") }
                            Spacer()
                            if let url = record.fileURL {
                                ShareLink(item: url) { Label("Export", systemImage: "square.and.arrow.up") }
                            }
                        }
                        .buttonStyle(.borderless)
                        DisclosureGroup("Details") {
                            Text(record.bundleID).font(.caption).textSelection(.enabled)
                            Text("externalVersionId \(record.externalVersionID)").font(.caption).textSelection(.enabled)
                            Text(record.date, style: .date).font(.caption)
                            if let url = record.fileURL,
                               let size = try? url.resourceValues(forKeys: [.fileSizeKey]).fileSize {
                                Text(ByteCountFormatter.string(fromByteCount: Int64(size), countStyle: .file)).font(.caption)
                            }
                            Text("Version and build read from the downloaded IPA.").font(.caption).foregroundStyle(.secondary)
                        }
                        Button(role: .destructive) { requestDeletion(record) } label: {
                            Label("Delete IPA", systemImage: "trash")
                        }
                        .buttonStyle(.borderless)
                    }
                    .padding(.vertical, 4)
                    .swipeActions(allowsFullSwipe: false) {
                        Button(role: .destructive) { requestDeletion(record) } label: { Label("Delete", systemImage: "trash") }
                    }
                }
            }
            .navigationTitle("Downloaded apps")
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
            .sheet(item: $installation) { OTAInstallationView(record: $0) }
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
