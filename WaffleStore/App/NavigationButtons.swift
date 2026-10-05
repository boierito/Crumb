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
                ButtonLabel(text: "Log In".localized, icon: "arrow.right")
            } else {
                ButtonLabel(text: "Sign in", icon: "key")
            }
        }
        .buttonStyle(FancyButtonStyle())
        .disabled(appData.appleId.isEmpty || appData.password.isEmpty || appData.isAuthenticating)
        .disabled(appData.hasSent2FACode ? appData.code.isEmpty : false)
    }
}

func extractAppId(from link: String) -> String {
    let trimmed = link.trimmingCharacters(in: .whitespacesAndNewlines)
    if let id = StoreParsing.identifier(trimmed) { return id }
    guard let url = URL(string: trimmed), url.scheme == "https",
          ["apps.apple.com", "itunes.apple.com"].contains(url.host?.lowercased() ?? "") else { return "" }
    return url.pathComponents.compactMap { $0.hasPrefix("id") ? StoreParsing.identifier(String($0.dropFirst(2))) : nil }.last ?? ""
}
