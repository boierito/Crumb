import SwiftUI
import UniformTypeIdentifiers

struct FairPlayDecryptView: View {
    let record: DownloadRecord
    @Environment(\.dismiss) private var dismiss
    @State private var importing = false
    @State private var working = false
    @State private var output: URL?
    @State private var workspace: URL?
    @State private var failure = ""

    var body: some View {
        NavigationStack {
            Form {
                Section("Offline FairPlay decryption") {
                    Text(record.appName)
                    Text("Import a private enrollment model generated for this IPA's Apple Account and SINF selector. Account age alone does not determine compatibility: legacy key-5 enrollments can be supported, but other migrated layouts require validation.")
                    Text("The exported IPA requires signing before installation.")
                        .foregroundStyle(.secondary)
                    if working {
                        ProgressView("Decrypting IPA…")
                    } else if let output {
                        ShareLink(item: output) { Label("Export decrypted IPA", systemImage: "square.and.arrow.up") }
                    } else {
                        Button("Import model and decrypt…") { importing = true }
                    }
                    if !failure.isEmpty { Text(failure).foregroundStyle(.red) }
                }
            }
            .navigationTitle("Decrypt IPA")
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Close") { dismiss() }.disabled(working)
                }
            }
            .interactiveDismissDisabled(working)
            .fileImporter(isPresented: $importing, allowedContentTypes: [.data]) { result in
                switch result {
                case .success(let url): start(model: url)
                case .failure: failure = "Could not open the selected model."
                }
            }
            .onDisappear {
                if let workspace { try? FileManager.default.removeItem(at: workspace) }
            }
        }
    }

    private func start(model: URL) {
        guard let input = record.fileURL else {
            failure = "The downloaded IPA is missing."
            return
        }
        working = true
        failure = ""
        Task {
            do {
                let directory = FileManager.default.temporaryDirectory
                    .appendingPathComponent("FairPlay-" + UUID().uuidString, isDirectory: true)
                workspace = directory
                let decrypted = try await Task.detached(priority: .userInitiated) {
                    try NativeFairPlay.decrypt(input: input, model: model, directory: directory)
                }.value
                output = decrypted
            } catch {
                failure = "Decryption failed. Check that the enrollment model matches this IPA and that the app uses a supported FairPlay profile."
                if let workspace { try? FileManager.default.removeItem(at: workspace) }
                workspace = nil
            }
            working = false
        }
    }
}

nonisolated enum NativeFairPlay {
    static func decrypt(input: URL, model: URL, directory: URL) throws -> URL {
        let manager = FileManager.default
        try manager.createDirectory(at: directory, withIntermediateDirectories: true,
                                    attributes: [.posixPermissions: 0o700])
        let privateModel = directory.appendingPathComponent("enrollment.bin")
        let accessing = model.startAccessingSecurityScopedResource()
        defer { if accessing { model.stopAccessingSecurityScopedResource() } }
        try manager.copyItem(at: model, to: privateModel)
        try manager.setAttributes([.posixPermissions: 0o600,
                                   .protectionKey: FileProtectionType.complete], ofItemAtPath: privateModel.path)
        defer { try? manager.removeItem(at: privateModel) }
        let output = directory.appendingPathComponent("App.decrypted.ipa")
        let report = directory.appendingPathComponent("report.json")
        defer { try? manager.removeItem(at: report) }
        let status = input.path.withCString { source in
            output.path.withCString { destination in
                report.path.withCString { reportPath in
                    privateModel.path.withCString { modelPath in
                        WaffleDecryptIPA(UnsafeMutablePointer(mutating: source),
                            UnsafeMutablePointer(mutating: destination),
                            UnsafeMutablePointer(mutating: reportPath),
                            UnsafeMutablePointer(mutating: modelPath))
                    }
                }
            }
        }
        guard status == 0 else { throw NSError(domain: "Crumb.FairPlay", code: Int(status)) }
        return output
    }
}
