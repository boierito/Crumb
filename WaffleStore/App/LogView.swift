import SwiftUI
import Combine

// Drain stdout even while the terminal is collapsed. Never let an invisible
// view leave a full pipe blocking login/download. UI updates stay on MainActor.
@MainActor
final class ActivityLog: ObservableObject {
    static let shared = ActivityLog()
    @Published private(set) var text = ""
    private var capturing = false
    func beginCapture() {
        guard !capturing else { return }
        capturing = true
        pipe.fileHandleForReading.readabilityHandler = { handle in
            let data = handle.availableData
            guard !data.isEmpty else { handle.readabilityHandler = nil; return }
            Task { @MainActor in
                let log = ActivityLog.shared
                log.text = String((log.text + String(decoding: data.suffix(32_768), as: UTF8.self)).suffix(32_768))
            }
        }
    }
}

struct LogView: View {
    @ObservedObject private var activity = ActivityLog.shared
    var body: some View {
        ScrollViewReader { proxy in
            ScrollView {
                Text(activity.text.isEmpty ? "No activity yet." : activity.text)
                    .font(.system(.caption2, design: .monospaced))
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .textSelection(.enabled)
                Color.clear.frame(height: 1).id("end")
            }
            .onChange(of: activity.text) { _ in proxy.scrollTo("end", anchor: .bottom) }
            .contextMenu {
                Button { UIPasteboard.general.string = activity.text } label: {
                    Label("Copy Output", systemImage: "doc.on.doc")
                }
            }
        }
    }
}
