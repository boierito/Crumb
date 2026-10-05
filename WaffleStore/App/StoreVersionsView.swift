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
    @State private var inspecting = false
    @State private var selectedID = ""
    @State private var confirmation = false
    @State private var versionDetail = ""
    var body: some View {
        NavigationStack {
            List {
                if inspecting { ProgressView("Reading selected IPA version…") }
                if loading { ProgressView("Resolving app, kbsync and available versions…") }
                if !error.isEmpty { Text(error).foregroundStyle(.red).textSelection(.enabled) }
                if let app = app {
                    Section(app.name) {
                        Text(app.bundleID).font(.caption)
                        Text("Select an externalVersionId. The displayed app version is verified from the downloaded IPA's Info.plist.")
                            .font(.caption).foregroundStyle(.secondary)
                        ForEach(versions, id: \.self) { id in
                            Button(id == latest ? "Latest iOS build — ID \(id)" : "Version ID \(id)") { inspect(app, version: id) }
                        }
                    }
                    Section("Specific externalVersionId") {
                        TextField("Numeric externalVersionId", text: $manualID).keyboardType(.numberPad)
                        Button("Download selected ID") { inspect(app, version: manualID) }
                            .disabled(StoreParsing.identifier(manualID) == nil)
                    }
                    Text("The IPA retains App Store protection. Exporting it does not prove that another sideloader can install or downgrade it.")
                        .font(.caption).foregroundStyle(.secondary)
                }
            }
            .disabled(inspecting)
            .confirmationDialog(versionDetail, isPresented: $confirmation, titleVisibility: .visible) {
                Button("Download IPA — ID \(selectedID)") {
                    guard let app = app, let tool = appData.ipaTool else { return }
                    appData.download(app: app, version: selectedID, tool: tool)
                    dismiss()
                }
                Button("Cancel", role: .cancel) {}
            }
            .navigationTitle("Download a version")
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
            .task {
                guard let tool = appData.ipaTool else { loading = false; return }
                appData.storeDiagnostic = "WaffleStore Store probe v3\napp-build=23003\nkbsync-runtime=tci-no-jit\nsecret-values=withheld"
                do {
                    let resolved = try await tool.lookup(appData.appLink)
                    app = resolved
                    appData.appBundleID = resolved.bundleID
                    let descriptor = try await tool.descriptor(app: resolved)
                    try Task.checkCancellation()
                    latest = descriptor.externalVersionID; versions = descriptor.availableVersionIDs
                } catch {
                    guard !Task.isCancelled else { return }
                    appData.storeDiagnostic += "\noutcome=versions-request-failed; code=\((error as NSError).code)"
                    self.error = (error as? StoreError)?.localizedDescription ?? (error as? SAPError)?.localizedDescription ?? "Store lookup failed (code \((error as NSError).code))."
                }
                loading = false
            }
        }
    }
    private func inspect(_ app: StoreApp, version: String) {
        guard let tool = appData.ipaTool, !inspecting else { return }
        inspecting = true; selectedID = version; error = ""
        Task {
            defer { inspecting = false }
            do {
                let descriptor = try await tool.descriptor(app: app, version: version)
                let info = try? await Task.detached(priority: .userInitiated) {
                    try NativePackage.inspect(url: descriptor.url, bundle: app.bundleID)
                }.value
                versionDetail = info.map { "Verified IPA version \($0.version) — externalVersionId \(version)" } ?? "externalVersionId \(version). CDN range inspection unavailable; version will be verified after download."
                confirmation = true
            } catch {
                self.error = (error as? StoreError)?.localizedDescription ?? (error as? SAPError)?.localizedDescription ?? "Version request failed (code \((error as NSError).code))."
            }
        }
    }
}
