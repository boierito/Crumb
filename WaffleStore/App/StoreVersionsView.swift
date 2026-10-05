import SwiftUI
import Combine
import MapleSAP

private struct VersionSelection: Identifiable { let id: String }

// One worker prevents concurrent Store requests and limits inspection to visible
// rows. Selection has priority; no complete IPA is downloaded to label a row.
@MainActor
private final class VersionLabels: ObservableObject {
    @Published var info: [String: NativePackage.Info] = [:]
    @Published var failures: [String: String] = [:]
    @Published var active: String?
    private var pending: [String] = []
    private var visible: Set<String> = []
    private var worker: Task<Void, Never>?
    func request(_ id: String, app: StoreApp, tool: IPATool, priority: Bool = false) {
        visible.insert(id)
        guard info[id] == nil, failures[id] == nil else { return }
        if active != id {
            pending.removeAll { $0 == id }
            if priority { pending.insert(id, at: 0) } else { pending.append(id) }
        }
        guard worker == nil else { return }
        worker = Task {
            defer { worker = nil; active = nil }
            while !pending.isEmpty, !Task.isCancelled {
                let next = pending.removeFirst()
                guard visible.contains(next) else { continue }
                active = next
                do {
                    let descriptor = try await tool.descriptor(app: app, version: next)
                    let result = try await Task.detached(priority: .utility) {
                        try NativePackage.inspect(url: descriptor.url, bundle: app.bundleID)
                    }.value
                    try Task.checkCancellation()
                    info[next] = result
                } catch {
                    guard !Task.isCancelled else { return }
                    failures[next] = "Version label unavailable; the downloaded IPA will be checked."
                }
                active = nil
            }
        }
    }
    func hide(_ id: String) { visible.remove(id); pending.removeAll { $0 == id } }
    func stop() { worker?.cancel(); pending.removeAll(); visible.removeAll() }
}

struct StoreVersionsView: View {
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @StateObject private var labels = VersionLabels()
    @State private var app: StoreApp?
    @State private var versions: [String] = []
    @State private var latest = ""
    @State private var error = ""
    @State private var loading = true
    @State private var manualID = ""
    @State private var selected: VersionSelection?
    var body: some View {
        NavigationStack {
            List {
                if loading { ProgressView("Resolving available versions…") }
                if !error.isEmpty { Text(error).foregroundStyle(.red).textSelection(.enabled) }
                if let app {
                    Section(app.name) {
                        Text(app.bundleID).font(.caption)
                        Text("Version numbers are read from the IPA. Labels load as rows become visible.")
                            .font(.caption).foregroundStyle(.secondary)
                        ForEach(versions, id: \.self) { id in
                            Button { select(id, app: app) } label: {
                                HStack {
                                    VStack(alignment: .leading) {
                                        Text(labels.info[id].map { "Version \($0.version)" } ?? (labels.failures[id] != nil ? "Version number unavailable" : (id == latest ? "Latest version" : "Reading version…")))
                                        Text("ID \(id)\(id == latest ? " · Latest" : "")").font(.caption).foregroundStyle(.secondary)
                                    }
                                    Spacer()
                                    if labels.active == id { ProgressView() }
                                }
                            }
                            .onAppear { if let tool = appData.ipaTool { labels.request(id, app: app, tool: tool) } }
                            .onDisappear { if selected?.id != id { labels.hide(id) } }
                        }
                    }
                    Section("Specific externalVersionId") {
                        TextField("Numeric externalVersionId", text: $manualID).keyboardType(.numberPad)
                        Button("Review selected version") { select(manualID, app: app) }
                            .disabled(StoreParsing.identifier(manualID) == nil)
                    }
                }
            }
            .navigationTitle("Download a version")
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } } }
            .sheet(item: $selected) { choice in
                NavigationStack {
                    Form {
                        Section("Selected version") {
                            Text(app?.name ?? "App")
                            Text("externalVersionId: \(choice.id)").textSelection(.enabled)
                            if let info = labels.info[choice.id] {
                                Text("Version: \(info.version)")
                                if let build = info.build { Text("Build: \(build)") }
                                Text("Read from IPA Info.plist").font(.caption).foregroundStyle(.secondary)
                            } else if let failure = labels.failures[choice.id] { Text(failure).font(.caption) }
                            else { ProgressView("Reading the selected IPA version…") }
                        }
                        Button("Download this version") {
                            guard let app, let tool = appData.ipaTool else { return }
                            let expected = labels.info[choice.id]?.version
                            labels.stop()
                            selected = nil
                            // Wait for the serialized inspector before starting Store work.
                            Task {
                                await labels.finish()
                                appData.download(app: app, version: choice.id, tool: tool, expectedVersion: expected)
                                dismiss()
                            }
                        }
                    }
                    .navigationTitle("Confirm download")
                    .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { selected = nil } } }
                }
                .presentationDetents([.medium, .large])
            }
            .onDisappear { labels.stop() }
            .task {
                guard let tool = appData.ipaTool else { loading = false; return }

                do {
                    let resolved = try await tool.lookup(appData.appLink)
                    app = resolved; appData.appBundleID = resolved.bundleID
                    let descriptor = try await tool.descriptor(app: resolved)
                    try Task.checkCancellation()
                    latest = descriptor.externalVersionID; versions = descriptor.availableVersionIDs
                } catch {
                    guard !Task.isCancelled else { return }

                    self.error = (error as? StoreError)?.localizedDescription ?? "Store lookup failed. Please try again."
                }
                loading = false
            }
        }
    }
    private func select(_ id: String, app: StoreApp) {
        selected = VersionSelection(id: id)
        if let tool = appData.ipaTool { labels.request(id, app: app, tool: tool, priority: true) }
    }
}
private extension VersionLabels {
    func finish() async { await worker?.value }
}
