import SwiftUI
import PartyUI

struct DownloadedAppsView: View {
    var embedded = false
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    @StateObject private var device = DeviceWorkflow.shared
    @State private var deletion: DownloadRecord?
    @State private var confirmDeleteAll = false
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

                if let record = device.activeRecord {
                    Section {
                        HStack(spacing:12) {
                            ProgressView()
                            VStack(alignment:.leading,spacing:4) {
                                Text(record.appName).font(.headline)
                                Text(device.activityTitle.localized).font(.subheadline).foregroundStyle(.secondary)
                            }
                            Spacer()
                            Button("Stop", role: .destructive) { device.stopDecryption() }
                                .buttonStyle(.borderless)
                                .frame(minWidth: 44, minHeight: 44)
                        }
                        if let progress = device.progress { ProgressView(value:progress) }
                        if device.waitingForHelper && !device.busy {
                            HStack {
                                Button("Open decrypter") { device.openHelper() }
                                Spacer()
                                Button("Save result") { device.retrieve() }
                            }.font(.subheadline).buttonStyle(.borderless)
                        }
                    }
                }
                Section("Encrypted") {
                    let records = appData.completedDownloads.filter { !$0.isDecrypted }
                    if records.isEmpty { Text("No encrypted apps").foregroundStyle(.secondary) }
                    ForEach(records) { record in downloadRow(record) }
                }
                Section("Decrypted") {
                    let records = appData.completedDownloads.filter { $0.isDecrypted }
                    if records.isEmpty { Text("No decrypted apps").foregroundStyle(.secondary) }
                    ForEach(records) { record in downloadRow(record) }
                }
                if !appData.completedDownloads.isEmpty || appData.downloadedIPAURL != nil {
                    Section {
                        Button("Delete all downloads", role: .destructive) { confirmDeleteAll = true }
                            .disabled(appData.isDowngrading)
                    }
                }
            }
            .navigationTitle("Downloads")
            .navigationBarTitleDisplayMode(.inline)
            .listStyle(.insetGrouped)
             .toolbar {
                if !embedded { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
            }

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
            .alert("Delete all downloads?", isPresented: $confirmDeleteAll) {
                Button("Delete all", role: .destructive) {
                    do {
                        try withAnimation(reduceMotion ? nil : .easeInOut(duration: 0.25)) {
                            try appData.deleteAllDownloads()
                        }
                    } catch {
                        deletionError = "Could not delete all downloads (code \((error as NSError).code))."
                    }
                }
                Button("Cancel", role: .cancel) { }
            } message: {
                Text("All downloaded IPAs will be removed. Installed apps, favourites and your Apple account will stay unchanged.")
            }
            .alert("Could not delete download", isPresented: Binding(
                get: { !deletionError.isEmpty }, set: { if !$0 { deletionError = "" } })) {
                    Button("OK") { deletionError = "" }
                } message: { Text(deletionError) }
            .onAppear { appData.restoreDownloadedIPA() }
        }
    }
    private func downloadRow(_ record: DownloadRecord) -> some View {
        VStack(alignment:.leading,spacing:12) {
            HStack(spacing:12) {
                Image(systemName:record.isDecrypted ? "lock.open.fill" : "app.fill")
                    .font(.title3).foregroundStyle(.tint)
                    .frame(width:42,height:42).background(Color.accentColor.opacity(0.1),in:RoundedRectangle(cornerRadius:10))
                VStack(alignment:.leading,spacing:3) {
                    Text(record.appName).font(.headline)
                    Text(record.version).font(.subheadline).foregroundStyle(.secondary)
                }
                Spacer(minLength:8)
                Menu {
                    if !record.isDecrypted && DeviceWorkflow.decryptSupported {
                        Button("Decrypt") { device.decrypt(record) }
                        Button("Decrypt and install") { device.decrypt(record,install:true) }
                    }
                    Button("Delete",role:.destructive) { requestDeletion(record) }
                } label: { Image(systemName:"ellipsis").frame(width:44,height:44) }
                    .accessibilityLabel("Options for \(record.appName)")
            }
            HStack(spacing:12) {
                Button { device.install(record) } label: { Label("Install",systemImage:"arrow.down.app").frame(maxWidth:.infinity) }
                    .disabled(device.busy || device.waitingForHelper)
                if let url = record.fileURL {
                    ShareLink(item:url) { Label("Export",systemImage:"square.and.arrow.up").frame(maxWidth:.infinity) }
                }
            }.buttonStyle(.bordered).controlSize(.regular)
        }
        .padding(.vertical,8)
        .swipeActions(allowsFullSwipe:false) {
            Button { requestDeletion(record) } label: { Label("Delete",systemImage:"trash") }.tint(.red)
        }
    }
    private func requestDeletion(_ record: DownloadRecord) {
        deletion = record; confirmDeletion = true
    }
}
