import SwiftUI

struct DownloadReadyView: View {
    let record: DownloadRecord
    let onInstall: () -> Void
    @Environment(\.dismiss) private var dismiss
    @State private var complete = false

    var body: some View {
        NavigationStack {
            Form {
                Section {
                    HStack(spacing: 14) {
                        CrumbProgressMark(progress: 1, complete: complete)
                        VStack(alignment: .leading, spacing: 4) {
                            Text(record.appName).font(.headline)
                            Text("Version \(record.version)").font(.subheadline).foregroundStyle(.secondary)
                        }
                    }.padding(.vertical, 6)
                }
                Section {
                    Button {
                        onInstall()
                        dismiss()
                    } label: {
                        Label("Install now", systemImage: "arrow.down.app")
                    }.buttonStyle(CrumbActionStyle())
                    if let url = record.fileURL {
                        ShareLink(item: url) { Label("Export IPA", systemImage: "square.and.arrow.up") }
                    }
                }
            }
            .navigationTitle("Download ready")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Done") { dismiss() } } }
            .task { complete = true }
        }
        .presentationDetents([.medium, .large])
    }
}
