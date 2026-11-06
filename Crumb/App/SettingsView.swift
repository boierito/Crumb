import SwiftUI
import PartyUI

struct SettingsView: View {
    var embedded = false
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @Environment(\.openURL) private var openURL
    @Environment(\.dynamicTypeSize) private var dynamicTypeSize
    
    @AppStorage("autoCleanApp") var autoCleanApp: Bool = true
    @StateObject private var localizationManager = LocalizationManager.shared
    
    private var appVersionString: String {
        let version = Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as? String ?? "?"
        let build = Bundle.main.object(forInfoDictionaryKey: "CFBundleVersion") as? String ?? "?"
        return "\(version) (\(build))"
    }
    
    var body: some View {
        NavigationStack {
            List {
                Section(header: HeaderLabel(text: "About".localized, icon: "info.circle")) {
                    VStack(alignment: .leading, spacing: 10) {
                        AppInfoCell()
                        let linksLayout = dynamicTypeSize.isAccessibilitySize
                            ? AnyLayout(VStackLayout(spacing: 10))
                            : AnyLayout(HStackLayout(spacing: 10))
                        linksLayout {
                            Button { openURL(URL(string: "https://t.me/certlium")!) } label: {
                                Label("Telegram", systemImage: "paperplane.fill")
                            }
                            .buttonStyle(CommunityButtonStyle())
                            Button { openURL(URL(string: "https://github.com/boierito/Crumb")!) } label: {
                                HStack(spacing: 8) {
                                    Image("github").renderingMode(.template).resizable().scaledToFit().frame(width: 20, height: 20)
                                    Text("GitHub")
                                }
                            }
                            .buttonStyle(CommunityButtonStyle())
                        }

                    }
                }
                
                Section(header: HeaderLabel(text: "Settings".localized, icon: "gearshape")) {
                    Toggle(isOn: $autoCleanApp) {
                        Text("Auto-Clean App".localized)
                        Text("Auto-Clean Description".localized)
                    }
                }
                
                Section { NavigationLink("Decrypter") { DeviceSetupView() } }

                Section(header: HeaderLabel(text: "Language".localized, icon: "globe"), footer:
                    Button(action: {
                        openURL(URL(string: "https://poeditor.com/join/project/Ofr2qvyudt")!)
                    }) {
                        Text("Translation Thank You".localized)
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                    .buttonStyle(.plain)
                ) {
                    Picker("Language".localized, selection: $localizationManager.currentLanguage) {
                        ForEach(Language.allCases) { language in
                            Text(language.displayName).tag(language)
                        }
                    }
                    .pickerStyle(.menu)
                }
                
                Section(header: HeaderLabel(text: "Credits".localized, icon: "star")) {
                    LinkCreditCell(image: Image("mineek"), name: "mineek", description: "Original creator of MuffinStore Jailed.".localized, url: "https://github.com/mineek")
                    LinkCreditCell(image: Image("lunginspector"), name: "lunginspector", description: "Original creator of PancakeStore.".localized, url: "https://github.com/lunginspector")
                    LinkCreditCell(image: Image("skadz"), name: "skadz", description: "Original creator of PancakeStore.".localized, url: "https://github.com/skadz108")
                    LinkCreditCell(image: Image("nxtcoreee3"), name: "nxtcoreee3", description: "UI changes and feature improvements.".localized, url: "https://github.com/nxtcoreee3")
                    LinkCreditCell(image: Image(systemName: "person.crop.circle"), name: "boierito", description: "Revival coordination and iOS 27 device testing.", url: "https://github.com/boierito")
                    LinkCreditCell(image: Image(systemName: "shippingbox"), name: "majd / ipatool", description: "MIT-licensed reference for SAP, Apple authentication and App Store downloads.", url: "https://github.com/majd/ipatool")
                }
            }
            .navigationTitle("Settings".localized)
            .toolbar {
                if !embedded { ToolbarItem(placement: .topBarTrailing) {
                    Button(action: {
                        dismiss()
                    }) {
                        Image(systemName: "xmark")
                    }
                } }
            }
        }
    }
}

private struct CommunityButtonStyle: ButtonStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .font(.headline)
            .lineLimit(1)
            .padding(.horizontal, 12)
            .frame(maxWidth: .infinity, minHeight: 50)
            .foregroundStyle(Color.accentColor)
            .background(Color.accentColor.opacity(configuration.isPressed ? 0.20 : 0.10), in: RoundedRectangle(cornerRadius: 14))
            .overlay(RoundedRectangle(cornerRadius: 14).strokeBorder(Color.accentColor.opacity(0.22), lineWidth: 1))
            .contentShape(RoundedRectangle(cornerRadius: 14))
    }
}
