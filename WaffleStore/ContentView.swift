import SwiftUI
import PartyUI

struct ContentView: View {
    @EnvironmentObject var appData: AppData
    @State private var confirmLogout = false
    @State private var showLogs = false

    var body: some View {
        TabView(selection: $appData.selectedTab) {
            AppSearchView(embedded: true)
                .tabItem { Label("Search", systemImage: "magnifyingglass") }.tag(AppTab.search)
            DownloadedAppsView(embedded: true)
                .tabItem { Label("Downloads", systemImage: "square.and.arrow.down") }.tag(AppTab.downloads)
            FavouritesView(embedded: true)
                .tabItem { Label("Favourites".localized, systemImage: "star") }.tag(AppTab.favourites)
            account
                .tabItem { Label("Account", systemImage: "person.crop.circle") }.tag(AppTab.account)
            SettingsView(embedded: true)
                .tabItem { Label("Settings".localized, systemImage: "gearshape") }.tag(AppTab.settings)
        }
        .sheet(isPresented: $appData.showStoreVersions) { StoreVersionsView() }
        .sheet(item: $appData.downloadReady) { DownloadReadyView(record: $0) }
        .sheet(item: $appData.installationRequest) { OTAInstallationView(record: $0, startImmediately: true) }
        .onAppear { appData.restoreStoreAccount(); appData.restoreDownloadedIPA() }
        .onChange(of: appData.isAuthenticated) { authenticated in
            if authenticated, appData.openVersionsAfterLogin {
                appData.openVersionsAfterLogin = false
                appData.openAppSelection(appData.appLink, country: appData.selectedCatalogCountry)
            }
        }
    }

    private var account: some View {
        NavigationStack {
            List {
                if appData.isAuthenticated {
                    Section("Signed in") {
                        Label(appData.appleId, systemImage: "person.crop.circle")
                        LabeledContent("Account region", value: appData.accountCountry.uppercased())
                    }
                    Section {
                        Button("Sign out", role: .destructive) { confirmLogout = true }
                            .disabled(appData.isAuthenticating || appData.isDowngrading || appData.storeRequestCount > 0 || appData.showStoreVersions)
                    } footer: {
                        Text("Signing out keeps your favourites and downloaded IPAs.")
                    }
                } else {
                    LoginSection
                    Section { NavigationButtons() }
                }
                Section {
                    DisclosureGroup("Activity log", isExpanded: $showLogs) { LogView() }
                }
            }
            .navigationTitle("Apple account")
            .navigationBarTitleDisplayMode(.inline)
            .alert("Sign out of Apple account?", isPresented: $confirmLogout) {
                Button("Sign out", role: .destructive) {
                    appData.openVersionsAfterLogin = false
                    appData.logoutStoreAccount()
                }
                Button("Cancel", role: .cancel) {}
            } message: {
                Text("Your favourites and downloaded IPAs will stay on this device.")
            }
        }
    }

    private var LoginSection: some View {
        Group {
            Section(header: HeaderLabel(text: "Login".localized, icon: "icloud"), footer: Text("Sign in to choose and download App Store versions.")) {
                VStack {
                    TextField("Apple ID".localized, text: $appData.appleId)
                        .modifier(TextFieldBackground())
                        .disabled(appData.hasSent2FACode || appData.isAuthenticating)
                        .autocorrectionDisabled()
                        .textInputAutocapitalization(.never)
                        .keyboardType(.emailAddress)
                        .textContentType(.username)
                    
                    HStack {
                        if appData.showPassword {
                            TextField("Password".localized, text: $appData.password)
                                .textContentType(.password)
                                .modifier(TextFieldBackground())
                                .disabled(appData.hasSent2FACode || appData.isAuthenticating)
                                .autocorrectionDisabled()
                                .textInputAutocapitalization(.never)
                        } else {
                            SecureField("Password".localized, text: $appData.password)
                                .textContentType(.password)
                                .modifier(TextFieldBackground())
                                .disabled(appData.hasSent2FACode || appData.isAuthenticating)
                                .autocorrectionDisabled()
                                .textInputAutocapitalization(.never)
                        }
                        
                        Button(action: {
                            appData.showPassword.toggle()
                        }) {
                            Image(systemName: appData.showPassword ? "eye" : "eye.slash")
                                .frame(width: 22, height: 22, alignment: .center)
                        }
                        .buttonStyle(TranslucentButtonStyle(useFullWidth: false))
                    }
                }
                if appData.isAuthenticating {
                    HStack(spacing: 12) {
                        ProgressView()
                        VStack(alignment: .leading, spacing: 4) {
                            Text(appData.applicationStatus).font(.subheadline)
                            if !appData.authenticationRecovery.isEmpty {
                                Text(appData.authenticationRecovery).font(.caption).foregroundStyle(.secondary)
                            }
                        }
                        Spacer(minLength: 0)
                        Button("Cancel") { appData.cancelAppleLogin() }
                            .buttonStyle(.borderless)
                    }
                    .accessibilityElement(children: .contain)
                }
                if !appData.authenticationError.isEmpty {
                    Text(appData.authenticationError).font(.subheadline).foregroundStyle(.red).textSelection(.enabled)
                }
            }
            if appData.hasSent2FACode {
                Section(header: HeaderLabel(text: "Verification Code".localized, icon: "key.viewfinder")) {
                    TextField("2FA Code".localized, text: $appData.code)
                        .modifier(TextFieldBackground())
                        .keyboardType(.numberPad)
                        .textContentType(.oneTimeCode)
                        .disabled(appData.isAuthenticating)
                    Button("Use another Apple ID") { appData.cancelAppleLogin() }
                        .disabled(appData.isAuthenticating)
                }
            }
        }
    }
    
}

struct DowngradeHistoryView: View {
    @EnvironmentObject var appData: AppData
    @Environment(\.dismiss) private var dismiss
    
    var body: some View {
        NavigationStack {
            List {
                Section(header: HeaderLabel(text: "Downgrade History".localized, icon: "clock.arrow.circlepath")) {
                    if appData.downgradeHistory.isEmpty {
                        VStack(alignment: .leading, spacing: 4) {
                            Text("No Downgrades Yet".localized)
                                .font(.headline)
                            Text("No Downgrades Description".localized)
                                .font(.subheadline)
                                .foregroundStyle(.secondary)
                        }
                    } else {
                        ForEach(appData.downgradeHistory) { entry in
                            DowngradeHistoryCell(entry: entry)
                        }
                    }
                }
            }
            .navigationTitle("History".localized)
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    Button(action: {
                        dismiss()
                    }) {
                        Image(systemName: "xmark")
                    }
                }
            }
        }
    }
}

struct DowngradeHistoryCell: View {
    let entry: DowngradeHistoryEntry
    
    private static let dateFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateStyle = .medium
        formatter.timeStyle = .short
        return formatter
    }()
    
    var body: some View {
        VStack(alignment: .leading, spacing: 10) {
            HStack(alignment: .top) {
                Image(systemName: "arrow.down.app")
                    .frame(width: 22, height: 22, alignment: .center)
                VStack(alignment: .leading, spacing: 4) {
                    Text(entry.bundleId)
                        .font(.headline)
                    Text("Version \(entry.installedVersion)")
                        .font(.subheadline)
                        .foregroundStyle(.secondary)
                    Text(Self.dateFormatter.string(from: entry.date))
                        .font(.caption)
                        .foregroundStyle(.secondary)
                    Text(entry.dataNote)
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
            }
        }
        .padding(.vertical, 4)
    }
}

struct ItemInfoCell: View {
    var label: String
    var icon: String
    var text: String
    
    var body: some View {
        LabeledContent {
            if text.isEmpty {
                ProgressView()
            } else {
                Text(text)
            }
        } label: {
            HStack {
                Image(systemName: icon)
                    .frame(width: 22, height: 22, alignment: .center)
                Text(label)
            }
        }
        .contextMenu {
            Button(action: {
                UIPasteboard.general.string = text
            }) {
                Label("Copy Value".localized, systemImage: "character.cursor.ibeam")
            }
        }
    }
}

#Preview {
    ContentView()
        .environmentObject(AppData.shared)
}
