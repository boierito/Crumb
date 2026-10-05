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
    @State private var showAdvanced = false
    @State private var loadAttempt = 0
    var body: some View {
        NavigationStack {
            List {
                if let app {
                    Section {
                        HStack(spacing: 12) {
                            VStack(alignment: .leading, spacing: 4) {
                                Text(app.name).font(.headline)
                                Text(app.bundleID).font(.caption).foregroundStyle(.secondary)
                            }
                        }
                    }
                }
                if loading {
                    Section {
                        HStack(spacing: 12) {
                            ProgressView()
                            Text(appData.applicationStatus == "Obtaining free app…" ? "Obtaining free app…" : "Finding available versions…")
                                .font(.subheadline).foregroundStyle(.secondary)
                        }
                    }
                } else if !error.isEmpty {
                    Section {
                        Label(error, systemImage: "exclamationmark.circle")
                            .font(.subheadline).foregroundStyle(.secondary).textSelection(.enabled)
                        Button("Try again") { loadAttempt += 1 }
                    }
                } else if let app {
                    Section {
                        ForEach(versions, id: \.self) { id in
                            Button { select(id, app: app) } label: {
                                HStack(spacing: 12) {
                                    VStack(alignment: .leading, spacing: 4) {
                                        Text(labels.info[id].map { "Version \($0.version)" } ?? (labels.failures[id] != nil ? "Version ID \(id)" : (id == latest ? "Latest version" : "Loading version…")))
                                            .foregroundStyle(.primary)
                                        if id == latest {
                                            Text("Latest").font(.caption).foregroundStyle(.secondary)
                                        } else if labels.info[id] == nil {
                                            Text(labels.failures[id] == nil ? "ID \(id)" : "Number unavailable")
                                                .font(.caption).foregroundStyle(.secondary)
                                        }
                                    }
                                    Spacer(minLength: 8)
                                    if labels.active == id { ProgressView() }
                                    Image(systemName: "chevron.right").font(.caption.weight(.semibold)).foregroundStyle(.tertiary)
                                }
                                .padding(.vertical, 4)
                                .frame(minHeight: 44)
                                .contentShape(Rectangle())
                            }
                            .buttonStyle(.plain)
                            .onAppear { if let tool = appData.ipaTool { labels.request(id, app: app, tool: tool) } }
                            .onDisappear { if selected?.id != id { labels.hide(id) } }
                        }
                    } header: {
                        Text("\(versions.count) available versions")
                    }
                    Section {
                        DisclosureGroup("Advanced", isExpanded: $showAdvanced) {
                            TextField("External version ID", text: $manualID).keyboardType(.numberPad)
                            Button("Choose this version") { select(manualID, app: app) }
                                .disabled(StoreParsing.identifier(manualID) == nil)
                        }
                    }
                }
            }
            .navigationTitle("Choose version")
            .navigationBarTitleDisplayMode(.inline)
            .listStyle(.insetGrouped)
            .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { labels.stop(); dismiss() } } }
            .sheet(item: $selected) { choice in
                NavigationStack {
                    Form {
                        Section("Selected version") {
                            Text(app?.name ?? "App")
                            if let info = labels.info[choice.id] {
                                Text("Version: \(info.version)")
                                if let build = info.build { Text("Build: \(build)") }
                            } else if let failure = labels.failures[choice.id] { Text(failure).font(.caption) }
                            else { HStack(spacing: 12) { ProgressView(); Text("Checking version…").foregroundStyle(.secondary) } }
                            DisclosureGroup("Version details") {
                                LabeledContent("External version ID", value: choice.id).font(.caption).textSelection(.enabled)
                            }
                        }
                        Section {
                            Button { download(choice, install: true) } label: {
                                Label("Download and install", systemImage: "arrow.down.app")
                            }
                            Button { download(choice, install: false) } label: {
                                Label("Download IPA", systemImage: "square.and.arrow.down")
                            }
                        } footer: {
                            Text("iOS will ask you to confirm installation after the download.")
                        }
                    }
                    .navigationTitle("Selected version")
                    .navigationBarTitleDisplayMode(.inline)
                    .toolbar { ToolbarItem(placement: .cancellationAction) { Button("Close") { selected = nil } } }
                }
                .presentationDetents([.medium, .large])
            }
            .onDisappear { if selected == nil { labels.stop() } }
            .task(id: loadAttempt) {
                loading = true; error = ""
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
                    appData.applicationStatus = "Could not load versions."
                    appData.applicationIcon = "exclamationmark.circle"
                }
                loading = false
            }
        }
    }
    private func download(_ choice: VersionSelection, install: Bool) {
        guard let app, let tool = appData.ipaTool else { return }
        let expected = labels.info[choice.id]?.version
        labels.stop()
        selected = nil
        Task {
            await labels.finish()
            appData.download(app: app, version: choice.id, tool: tool,
                expectedVersion: expected, installWhenReady: install)
            dismiss()
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
