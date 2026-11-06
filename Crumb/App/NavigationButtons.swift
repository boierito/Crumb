import SwiftUI
import PartyUI
import MapleSAP

// Account owns sign-in; other actions live on their respective tabs.
struct NavigationButtons: View {
    @EnvironmentObject var appData: AppData
    var body: some View {
        Button(action: {
            Haptic.shared.play(.soft)
            appData.startAppleLogin()
        }) {
            if appData.isAuthenticating {
                ButtonLabel(text: "Signing in…", icon: "hourglass")
            } else if appData.hasSent2FACode {
                ButtonLabel(text: appData.signInForm.isRetryingVerification ? "Retry verification" : "Verify code", icon: "checkmark.shield")
            } else {
                ButtonLabel(text: "Sign in", icon: "key")
            }
        }
        .buttonStyle(FancyButtonStyle())
        .disabled(!appData.canSubmitAppleLogin || appData.isAuthenticating)
    }
}

func extractAppId(from link: String) -> String {
    let trimmed = link.trimmingCharacters(in: .whitespacesAndNewlines)
    if let id = StoreParsing.identifier(trimmed) { return id }
    guard let url = URL(string: trimmed), url.scheme == "https",
          ["apps.apple.com", "itunes.apple.com"].contains(url.host?.lowercased() ?? "") else { return "" }
    return url.pathComponents.compactMap { $0.hasPrefix("id") ? StoreParsing.identifier(String($0.dropFirst(2))) : nil }.last ?? ""
}
