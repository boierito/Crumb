import SwiftUI
import MapleSAP

struct StoreVersionsView: View {
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @State private var app: StoreApp?
    @State private var versions: [String] = []
    @State private var latest: String = ""
    @State private var error: String = ""
    @State private var loading = true
    @State private var manualID = ""
    var body: some View {
        NavigationStack {
            List {
                if loading { ProgressView("Resolving app, kbsync and available versions…") }
                if !error.isEmpty { Text(error).foregroundStyle(.red).textSelection(.enabled) }
                if let app = app {
                    Section(app.name) {
                        Text(app.bundleID).font(.caption)
                        Text("Select an externalVersionId. The displayed app version is verified from the downloaded IPA's Info.plist.")
                            .font(.caption).foregroundStyle(.secondary)
                        ForEach(versions, id: \.self) { id in
                            Button(id == latest ? "Latest iOS build — ID \(id)" : "Version ID \(id)") { start(app, version: id) }
                        }
                    }
                    Section("Specific externalVersionId") {
                        TextField("Numeric externalVersionId", text: $manualID).keyboardType(.numberPad)
                        Button("Download selected ID") { start(app, version: manualID) }
                            .disabled(StoreParsing.identifier(manualID) == nil)
                    }
                    Text("The IPA retains App Store protection. Exporting it does not prove that another sideloader can install or downgrade it.")
                        .font(.caption).foregroundStyle(.secondary)
                }
            }
            .navigationTitle("Download a version")
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
            .task {
                guard let tool = appData.ipaTool else { loading = false; return }
                do {
                    let resolved = try await tool.lookup(appData.appLink)
                    app = resolved
                    appData.appBundleID = resolved.bundleID
                    let descriptor = try await tool.descriptor(app: resolved)
                    try Task.checkCancellation()
                    latest = descriptor.externalVersionID; versions = descriptor.availableVersionIDs
                } catch {
                    guard !Task.isCancelled else { return }
                    self.error = (error as? StoreError)?.localizedDescription ?? (error as? SAPError)?.localizedDescription ?? "Store lookup failed (code \((error as NSError).code))."
                }
                loading = false
            }
        }
    }
    private func start(_ app: StoreApp, version: String) {
        guard let tool = appData.ipaTool else { return }
        appData.download(app: app, version: version, tool: tool)
        dismiss()
    }
}
