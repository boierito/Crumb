import SwiftUI

struct DownloadReadyView: View {
    let record: DownloadRecord
    @Environment(\.dismiss) private var dismiss
    @State private var install = false
    @State private var decryption: DownloadDecryptionAction?
    var body: some View {
        NavigationStack {
            List {
                Section {
                    VStack(alignment: .leading, spacing: 10) {
                        Label("Saved on your iPhone", systemImage: "checkmark.circle.fill")
                            .font(.caption.weight(.semibold)).foregroundStyle(.green)
                        Text(record.appName).font(.system(.title2, design: .rounded, weight: .bold))
                        Text("Version \(record.version)\(record.build.map { " · Build " + $0 } ?? "")")
                            .font(.subheadline).foregroundStyle(.secondary)
                    }.padding(.vertical, 12)
                }
                Section {
                    DownloadActions(record: record, install: { install = true }, decrypt: {
                        decryption = DownloadDecryptionAction(record: record, install: $0)
                    })
                } footer: {
                    Text("Your original IPA stays in Downloads.")
                }
            }
            .navigationTitle("Choose what’s next").navigationBarTitleDisplayMode(.inline)
            .toolbar { ToolbarItem(placement: .confirmationAction) { Button("Done") { dismiss() } } }
            .sheet(isPresented: $install) { OTAInstallationView(record: record, startImmediately: true) }
            .sheet(item: $decryption) { SystemFairPlayView(record: $0.record, startAction: $0) }
        }
    }
}
