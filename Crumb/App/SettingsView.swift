import SwiftUI
import PartyUI

struct SettingsView: View {
    @AppStorage("autoCleanApp") private var autoCleanApp = true
    @StateObject private var localizationManager = LocalizationManager.shared
    @AppStorage("crumbAppearance") private var appearance = CrumbAppearance.system

    var body: some View {
        List {
            Section("Language and appearance") {
                Picker("Appearance", selection: $appearance) {
                    Text("Use system setting").tag(CrumbAppearance.system)
                    Text("Light").tag(CrumbAppearance.light)
                    Text("Dark").tag(CrumbAppearance.dark)
                }
                Picker("Language".localized, selection: $localizationManager.currentLanguage) {
                    ForEach(Language.allCases) { language in
                        Text(language.displayName).tag(language)
                    }
                }
            }
            Section {
                Toggle("Auto-Clean App".localized, isOn: $autoCleanApp)
            } header: {
                Text("Storage")
            } footer: {
                Text("Automatically removes temporary files at launch and keeps completed downloads.")
            }
            Section("Device and signing") {
                NavigationLink { DeviceSetupView() } label: { Label("Device tools", systemImage: "iphone") }
                NavigationLink { SigningSettingsView() } label: { Label("Signing identity", systemImage: "checkmark.seal") }
            }
            Section {
                NavigationLink { CrumbAboutView() } label: { Label("About Crumb", systemImage: "info.circle") }
            }
        }
        .navigationTitle("Settings".localized)
    }
}

private struct CrumbAboutView: View {
    @Environment(\.openURL) private var openURL
    @Environment(\.dynamicTypeSize) private var dynamicTypeSize

    var body: some View {
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
                
            Section("Translations") {
                Link("Help translate Crumb", destination: URL(string: "https://poeditor.com/join/project/Ofr2qvyudt")!)
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
        .navigationTitle("About Crumb")
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
