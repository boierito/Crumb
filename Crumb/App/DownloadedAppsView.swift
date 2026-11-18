import SwiftUI
import PartyUI

struct DownloadedAppsView: View {
    var embedded = false
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    @StateObject private var device = DeviceWorkflow.shared
    @State private var deletion: DownloadRecord?
    @State private var signingDetails: DownloadRecord?
    @State private var confirmDeleteAll = false
    @State private var confirmDeletion = false
    @State private var deletionError = ""
    @State private var query = ""
    @State private var filter: PackageFilter = .all

    private enum PackageFilter: String, CaseIterable {
        case all, encrypted, decrypted

        var title: LocalizedStringKey {
            switch self {
            case .all: "All files"
            case .encrypted: "Encrypted"
            case .decrypted: "Decrypted"
            }
        }
    }

    private var matchingRecords: [DownloadRecord] {
        let term = query.trimmingCharacters(in: .whitespacesAndNewlines)
        return appData.completedDownloads.filter { record in
            let includesType = filter == .all || (filter == .decrypted) == record.isDecrypted
            return includesType && (term.isEmpty || [record.appName, record.version, record.bundleID]
                .contains { $0.localizedCaseInsensitiveContains(term) })
        }
    }
    var body: some View {
        NavigationStack {
            List {
                Color.clear.frame(height: 0).listRowInsets(EdgeInsets())
                    .listRowBackground(Color.clear).accessibilityHidden(true).id(AppTab.downloads)
                if appData.isDowngrading {
                    Section {
                        VStack(alignment: .leading, spacing: 12) {
                            HStack(spacing: 12) {
                                let transferring = appData.downloadPhaseLabel == "Downloading"
                                let fraction = min(max(appData.downgradeProgress / 0.85, 0), 1)
                                CrumbProgressMark(progress: transferring ? fraction : nil)
                                VStack(alignment: .leading, spacing: 4) {
                                    Text(appData.activeDownloadName.isEmpty ? "Downloading app" : appData.activeDownloadName)
                                        .font(.headline)
                                    Text(appData.downloadPhaseLabel).font(.subheadline).foregroundStyle(.secondary)
                                    if transferring {
                                        Text("\(appData.downloadProgressText) · \(Int(fraction * 100))%")
                                            .font(.caption.monospacedDigit()).foregroundStyle(.secondary)
                                    }
                                }
                                .accessibilityElement(children: .combine)
                                Spacer(minLength: 0)
                                Button { appData.storeTask?.cancel() } label: {
                                    Image(systemName: "xmark").frame(width: 44, height: 44)
                                        .background(.thinMaterial, in: Circle())
                                }
                                .buttonStyle(.borderless).accessibilityLabel("Cancel download")
                            }
                        }.padding(.vertical, 6)
                    } header: { Text("In progress") }
                }
                if !appData.storeError.isEmpty {
                    Section {
                        Label(appData.activeDownloadName, systemImage: "exclamationmark.circle").font(.headline)
                        Text(appData.storeError).font(.subheadline).foregroundStyle(.secondary)
                        if let retry = appData.retryDownload {
                            Button("Try again", action: retry).buttonStyle(CrumbActionStyle(primary: false))
                        }
                    } header: { Text("Download failed") }
                }

                if device.activeRecord != nil || device.waitingForHelper {
                    Section {
                        HStack(spacing:12) {
                            CrumbProgressMark(progress: device.progress)
                            VStack(alignment:.leading,spacing:4) {
                                Text(device.activeRecord?.appName ?? "Decrypter".localized).font(.headline)
                                Text(device.activityTitle.localized).font(.subheadline).foregroundStyle(.secondary)
                            }
                            Spacer()
                            Button("Stop", role: .destructive) { device.stopDecryption() }
                                .buttonStyle(.borderless)
                                .frame(minWidth: 44, minHeight: 44)
                        }
                        if device.waitingForHelper && !device.busy {
                            HStack {
                                Button("Open decrypter") { device.openHelper() }
                                Spacer()
                                Button("Save result") { device.retrieve() }
                            }.font(.subheadline).buttonStyle(.borderless)
                        }
                    }
                }
                if !appData.completedDownloads.isEmpty {
                    Section {
                        TextField("Search name, version or bundle ID", text: $query)
                            .autocorrectionDisabled().textInputAutocapitalization(.never)
                        Picker("File type", selection: $filter) {
                            ForEach(PackageFilter.allCases, id: \.self) { item in
                                Text(item.title).tag(item)
                            }
                        }
                    }
                    Section("Ready to install") {
                        if matchingRecords.isEmpty {
                            CrumbEmptyState(title: "No matching downloads", message: "Try another search or file type.", icon: "line.3.horizontal.decrease.circle")
                        }
                        ForEach(matchingRecords) { record in downloadRow(record) }
                    }
                } else if !appData.isDowngrading && device.activeRecord == nil && !device.waitingForHelper && appData.storeError.isEmpty {
                    Section {
                        CrumbEmptyState(title: "No downloads yet", message: "Downloaded apps appear here when they’re ready to install or export.", icon: "arrow.down.circle")
                        Button("Search App Store") { appData.selectedTab = .search }
                            .buttonStyle(CrumbActionStyle(primary: false))
                    }
                }
                if !appData.completedDownloads.isEmpty || appData.downloadedIPAURL != nil {
                    Section {
                        Button("Delete all downloads", role: .destructive) { confirmDeleteAll = true }
                            .disabled(appData.isDowngrading)
                    }
                }
            }
            .crumbPrimaryPage(.downloads)
            .navigationTitle("Downloads")
            .scrollDismissesKeyboard(.interactively)
            .listStyle(.insetGrouped)
             .toolbar {
                if !embedded { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
            }

            .sheet(item:$signingDetails) { record in
                NavigationStack {
                    if let summary = record.lastSigning { SigningDetailsView(summary:summary).toolbar { ToolbarItem(placement:.confirmationAction) { Button("Done") { signingDetails = nil } } } }
                }
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
                if let deletion {
                    Text("Delete the saved IPA for \(deletion.appName), version \(deletion.version)? The installed app and its data stay on this device.")
                }
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
                CrumbArtwork(url: appData.favourites.first { $0.bundleId == record.bundleID }?.artworkURL, size: 44)
                VStack(alignment: .leading, spacing: 4) {
                    Text(record.appName).font(.headline)
                    Text(record.version).font(.subheadline.monospacedDigit()).foregroundStyle(.secondary)
                    Text(record.isDecrypted ? "Decrypted IPA" : "Encrypted IPA")
                        .font(.caption).foregroundStyle(.secondary)
                }
                Spacer(minLength:8)
                Menu {
                    if let url = record.fileURL {
                        ShareLink(item: url) { Label("Export IPA", systemImage: "square.and.arrow.up") }
                    }
                    if record.isDecrypted {
                        Button("Signing options") { device.install(record,review:true) }
                        if record.lastSigning != nil { Button("Signing details") { signingDetails = record } }
                    }
                    if !record.isDecrypted && DeviceWorkflow.decryptSupported {
                        Button("Decrypt") { device.decrypt(record) }
                        Button("Decrypt and install") { device.decrypt(record,install:true) }
                    }
                    Button("Delete",role:.destructive) { requestDeletion(record) }
                } label: { Image(systemName:"ellipsis").frame(width:44,height:44) }
                    .accessibilityLabel("Options for \(record.appName)")
            }
            Button { device.install(record) } label: {
                Label("Install", systemImage: "arrow.down.app")
            }
            .buttonStyle(CrumbActionStyle())
            .disabled(device.busy || device.waitingForHelper)

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
