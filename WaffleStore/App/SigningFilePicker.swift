import SwiftUI
import UniformTypeIdentifiers
import Darwin

// User-supplied compatibility binary, pinned by the artifact verifier. It is
// loaded only when using a picker, not while handling Apple authentication.
@MainActor
enum FilePickerCompatibility {
    private static var handle: UnsafeMutableRawPointer?
    static func load() {
        guard handle == nil else { return }
        let file = Bundle.main.bundleURL.appendingPathComponent("Frameworks/FilePickerFix.dylib")
        handle = dlopen(file.path, RTLD_NOW | RTLD_LOCAL)
    }
}

struct SigningFilePicker: UIViewControllerRepresentable {
    let selected: (URL) -> Void
    let cancelled: () -> Void
    func makeCoordinator() -> Coordinator { Coordinator(self) }
    func makeUIViewController(context: Context) -> UIDocumentPickerViewController {
        FilePickerCompatibility.load()
        // Copy mode avoids relying on sideloaded security-scope entitlements.
        // Any moving performed by FilePickerFix operates on the provider copy.
        let picker = UIDocumentPickerViewController(forOpeningContentTypes: [.item], asCopy: true)
        picker.allowsMultipleSelection = false
        picker.delegate = context.coordinator
        return picker
    }
    func updateUIViewController(_ controller: UIDocumentPickerViewController, context: Context) {}
    final class Coordinator: NSObject, UIDocumentPickerDelegate {
        let parent: SigningFilePicker
        init(_ parent: SigningFilePicker) { self.parent = parent }
        func documentPicker(_ controller: UIDocumentPickerViewController, didPickDocumentsAt urls: [URL]) {
            guard let file = urls.first else { parent.cancelled(); return }
            parent.selected(file)
        }
        func documentPickerWasCancelled(_ controller: UIDocumentPickerViewController) { parent.cancelled() }
    }
}
