import SwiftUI

struct DownloadReadyView: View {
    let record: DownloadRecord
    @Environment(\.dismiss) private var dismiss
    var body: some View {
        NavigationStack {
            Form {
                Section("Download ready") {
                    Text(record.appName).font(.headline)
                    Text("Version \(record.version)")

                }
                Section {
                    Button { dismiss(); DeviceWorkflow.shared.install(record) } label: { Label("Install now", systemImage: "arrow.down.app") }
                    if let url = record.fileURL {
                        ShareLink(item: url) { Label("Export IPA", systemImage: "square.and.arrow.up") }
                    }
                }
            }
            .navigationTitle("Download complete")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Done") { dismiss() } } }

        }
    }
}
