import SwiftUI

struct DownloadReadyView: View {
    let record: DownloadRecord
    @Environment(\.dismiss) private var dismiss
    @State private var install = false
    var body: some View {
        NavigationStack {
            Form {
                Section("Download ready") {
                    Text(record.appName).font(.headline)
                    Text("Version \(record.version)")
                    if let build = record.build { Text("Build \(build)") }
                    Text(record.bundleID).font(.caption).foregroundStyle(.secondary)
                }
                Section {
                    Button { install = true } label: { Label("Install now", systemImage: "arrow.down.app") }
                    if let url = record.fileURL {
                        ShareLink(item: url) { Label("Export IPA", systemImage: "square.and.arrow.up") }
                    }
                } footer: {
                    Text("The IPA is saved in Downloaded apps. iOS will ask you to confirm installation.")
                }
            }
            .navigationTitle("Download complete")
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Done") { dismiss() } } }
            .sheet(isPresented: $install) { OTAInstallationView(record: record, startImmediately: true) }
        }
    }
}
