import SwiftUI

struct DownloadDecryptionAction: Identifiable {
    let id = UUID()
    let record: DownloadRecord
    let install: Bool
}

// The same four actions appear after a download and in the saved library.
struct DownloadActions: View {
    let record: DownloadRecord
    let install: () -> Void
    let decrypt: (Bool) -> Void
    var body: some View {
        VStack(spacing: 0) {
            if let url = record.fileURL {
                ShareLink(item: url) {
                    DownloadActionLabel(title: "Save encrypted IPA", subtitle: "Original App Store download", icon: "lock.doc", prominent: false)
                }
                Divider().padding(.leading, 52)
                Button(action: install) {
                    DownloadActionLabel(title: "Save & install encrypted", subtitle: "Request the original installation", icon: "arrow.down.app", prominent: false)
                }
                Divider().padding(.leading, 52)
                Button { decrypt(false) } label: {
                    DownloadActionLabel(title: "Decrypt IPA", subtitle: "Verify every protected executable", icon: "lock.open", prominent: true)
                }
                Divider().padding(.leading, 52)
                Button { decrypt(true) } label: {
                    DownloadActionLabel(title: "Decrypt & install", subtitle: "Use your local signing identity", icon: "checkmark.shield", prominent: true)
                }
            } else {
                Label("This IPA is no longer in Files. Download the version again to use these actions.", systemImage: "exclamationmark.triangle")
                    .font(.footnote).foregroundStyle(.secondary).padding()
            }
        }
        .buttonStyle(.plain)
    }
}

private struct DownloadActionLabel: View {
    let title, subtitle, icon: String
    let prominent: Bool
    var body: some View {
        HStack(spacing: 14) {
            Image(systemName: icon).font(.title3.weight(.medium))
                .foregroundStyle(prominent ? Color.accentColor : Color.secondary)
                .frame(width: 30)
            VStack(alignment: .leading, spacing: 3) {
                Text(title).font(.body.weight(.semibold)).foregroundStyle(.primary)
                Text(subtitle).font(.caption).foregroundStyle(.secondary)
            }
            Spacer(minLength: 8)
            Image(systemName: "chevron.right").font(.caption.weight(.semibold)).foregroundStyle(.tertiary)
        }
        .frame(maxWidth: .infinity, minHeight: 56, alignment: .leading)
        .padding(.vertical, 8).contentShape(Rectangle())
        .accessibilityElement(children: .combine)
    }
}
