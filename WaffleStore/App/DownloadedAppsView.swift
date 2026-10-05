import SwiftUI

struct DownloadedAppsView: View {
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @State private var installation: DownloadRecord?
    var body: some View {
        NavigationStack {
            List {
                if appData.completedDownloads.isEmpty { Text("No verified downloads yet.") }
                ForEach(appData.completedDownloads) { record in
                    Section(record.appName) {
                        Text("Version \(record.version)")
                        if let build = record.build { Text("Build \(build)") }
                        Text("externalVersionId \(record.externalVersionID)").font(.caption).textSelection(.enabled)
                        Text(record.bundleID).font(.caption)
                        Text("Version read from IPA Info.plist; Store metadata matched the selected externalVersionId.").font(.caption).foregroundStyle(.secondary)
                        Button("Install this downloaded version") { installation = record }
                        if let url = record.fileURL { ShareLink(item: url) { Label("Export IPA", systemImage: "square.and.arrow.up") } }
                    }
                }
            }
            .navigationTitle("Downloaded apps")
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
            .sheet(item: $installation) { OTAInstallationView(record: $0) }
            .onAppear { appData.completedDownloads = DownloadRecord.load() }
        }
    }
}
