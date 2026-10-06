import SwiftUI
import MapleSAP

// Branch-only controls; the production main checkout has neither this view nor
// the reset API. Confirmation is centered, not a source-anchored popover.
struct AuthenticationTestSection: View {
    @EnvironmentObject private var appData: AppData
    @State private var confirmReset = false
    @State private var showReport = false

    var body: some View {
        Section {
            Label("Crumb Test · 23023", systemImage: "testtube.2")
            Button {
                confirmReset = true
            } label: {
                HStack {
                    Label("Reset test sign-in", systemImage: "arrow.counterclockwise")
                    Spacer()
                    if appData.isResettingAuthenticationTest { ProgressView() }
                }
            }
            .disabled(!appData.canResetAuthenticationTest)
            if appData.authenticationTestResetVerified {
                Label("New identity verified in Keychain", systemImage: "checkmark.shield")
                    .font(.subheadline).foregroundStyle(.secondary)
            }
            Button { showReport = true } label: {
                Label("Test report", systemImage: "doc.text.magnifyingglass")
            }
        } header: {
            Text("Fresh-login test")
        } footer: {
            Text("Tap Reset test sign-in before logging in. It resets only this test app’s session, cookies and machine identity. Apple decides whether to request 2FA. SAP assets stay cached.")
        }
        .alert("Reset this test app’s sign-in?", isPresented: $confirmReset) {
            Button("Reset test sign-in", role: .destructive) {
                Task { await appData.resetAuthenticationTest() }
            }
            Button("Cancel", role: .cancel) {}
        } message: {
            Text("You’ll enter your Apple ID and password again. The original Crumb account and files are not changed. This cannot remove trust stored by Apple.")
        }
        .sheet(isPresented: $showReport) {
            NavigationStack {
                ScrollView {
                    Text(appData.authenticationTestReport.text)
                        .font(.system(.footnote, design: .monospaced))
                        .textSelection(.enabled)
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .padding()
                }
                .navigationTitle("Test report")
                .navigationBarTitleDisplayMode(.inline)
                .toolbar {
                    ToolbarItem(placement: .cancellationAction) {
                        Button("Close") { showReport = false }
                    }
                    ToolbarItem(placement: .primaryAction) {
                        ShareLink(item: appData.authenticationTestReport.text) {
                            Label("Share report", systemImage: "square.and.arrow.up")
                        }
                    }
                }
            }
        }
    }
}
