import SwiftUI
import PartyUI

struct SettingsView: View {
    var embedded = false
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    @Environment(\.openURL) private var openURL
    
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
                        HStack {
                            Button(action: {
                                openURL(URL(string: "https://discord.com/invite/tweakbreak-1443331342799601666")!)
                            }) {
                                ButtonLabel(text: "Discord".localized, icon: "discord", useImage: true)
                            }
                            .buttonStyle(TranslucentButtonStyle(color: .discord))
                            Button(action: {
                                openURL(URL(string: "https://github.com/nxtcoreee3/WaffleStore")!)
                            }) {
                                ButtonLabel(text: "GitHub".localized, icon: "github", useImage: true)
                            }
                            .buttonStyle(TranslucentButtonStyle(color: .github))
                        }
                        Button(action: {
                            openURL(URL(string: "https://nxtcoreee3.github.io/WaffleStore/")!)
                        }) {
                            ButtonLabel(text: "Website".localized, icon: "globe")
                        }
                        .buttonStyle(TranslucentButtonStyle())
                    }
                }
                
                Section(header: HeaderLabel(text: "Settings".localized, icon: "gearshape")) {
                    Toggle(isOn: $autoCleanApp) {
                        Text("Auto-Clean App".localized)
                        Text("Auto-Clean Description".localized)
                    }
                }
                
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
                
                if !appData.completedDownloads.isEmpty {
                    Section("Downloaded IPAs") {
                        ForEach(appData.completedDownloads) { record in
                            if let url = record.fileURL {
                                ShareLink(item: url) {
                                    VStack(alignment: .leading) {
                                        Text("\(record.appName) — \(record.version)")
                                        Text("externalVersionId \(record.externalVersionID)").font(.caption).foregroundStyle(.secondary)
                                    }
                                }
                            }
                        }
                    }
                }
                Section(header: HeaderLabel(text: "Data".localized, icon: "loupe"), footer: Text("Storage Warning".localized)) {
                    VStack {
                        Button(action: {
                            if let url = appData.downloadedIPAURL { presentShareSheet(with: url) }
                        }) {
                            ButtonLabel(text: "Export IPA".localized, icon: "arrow.up.doc")
                        }
                        .buttonStyle(TranslucentButtonStyle())
                        .disabled(!appData.hasAppBeenServed)
                        Button(action: {
                            cleanUp()
                        }) {
                            ButtonLabel(text: "Clean Documents".localized, icon: "trash")
                        }
                        .buttonStyle(TranslucentButtonStyle())
                        
                    }
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
