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
    @State private var expanded = false
    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            HStack {
                Button { expanded = true } label: { Label("Expand", systemImage: "arrow.up.left.and.arrow.down.right") }
                Spacer()
                Button { UIPasteboard.general.string = activity.text } label: { Label("Copy", systemImage: "doc.on.doc") }
            }
            .font(.caption)
            .buttonStyle(.borderless)
            LogText(followsLatest: true)
        }
        .padding(.vertical, 10)
        .frame(maxWidth: .infinity, alignment: .leading)
        .sheet(isPresented: $expanded) { ActivityLogScreen() }
    }
}

private struct LogText: View {
    @ObservedObject private var activity = ActivityLog.shared
    let followsLatest: Bool
    var body: some View {
        ScrollViewReader { proxy in
            ScrollView {
                Text(activity.text.isEmpty ? "No activity yet." : activity.text)
                    .font(.system(.caption, design: .monospaced))
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .textSelection(.enabled)
                Color.clear.frame(height: 1).id("end")
            }
            .onAppear { if followsLatest { proxy.scrollTo("end", anchor: .bottom) } }
            .onChange(of: activity.text) { _ in if followsLatest { proxy.scrollTo("end", anchor: .bottom) } }
            .onChange(of: followsLatest) { enabled in if enabled { proxy.scrollTo("end", anchor: .bottom) } }
        }
    }
}

private struct ActivityLogScreen: View {
    @ObservedObject private var activity = ActivityLog.shared
    @Environment(\.dismiss) private var dismiss
    @State private var followsLatest = false
    var body: some View {
        NavigationStack {
            VStack(spacing: 12) {
                Toggle("Follow latest activity", isOn: $followsLatest).font(.subheadline)
                LogText(followsLatest: followsLatest)
            }
            .padding()
            .navigationTitle("Activity log")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) { Button("Close") { dismiss() } }
                ToolbarItem(placement: .primaryAction) {
                    Button { UIPasteboard.general.string = activity.text } label: { Label("Copy log", systemImage: "doc.on.doc") }
                }
            }
        }
    }
}
