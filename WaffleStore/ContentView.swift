import SwiftUI
import PartyUI

struct ContentView: View {
    @State private var hasShownWelcome: Bool = false
    @State private var showLogs: Bool = false
    @State private var showSettingsView: Bool = false
    @State private var showSearchView: Bool = false
    @State private var showHistoryView: Bool = false
    @State private var showFavouritesView: Bool = false
    @State private var showDownloadedView = false
    
    @EnvironmentObject var appData: AppData
    @StateObject private var localizationManager = LocalizationManager.shared
    
    var body: some View {
        Group {
            if UIDevice.current.userInterfaceIdiom == .pad {
                NavigationSplitView(sidebar: {
                    List {
                        LogsSection
                        NavigationButtons()
                    }
                    .navigationTitle("WaffleStore")
                }) {
                    List {
                        if !appData.isAuthenticated {
                            LoginSection
                        } else {
                            if appData.isDowngrading {
                                AppInfoSection
                            } else {
                                InputAppSection
                            }
                        }
                    }
                    .toolbar {
                        ToolbarItem(placement: .topBarTrailing) {
                            AppMenu(showHistoryView: $showHistoryView, showFavouritesView: $showFavouritesView, showDownloadedView: $showDownloadedView)
                        }
                        ToolbarItem(placement: .topBarTrailing) {
                            Button(action: {
                                showSettingsView.toggle()
                            }) {
                                Image(systemName: "gear")
                            }
                        }
                    }
                }
            } else {
                NavigationStack {
                    List {
                        if appData.isAuthenticated { LogsSection }
                        if !appData.isAuthenticated {
                            LoginSection
                            LogsSection
                        } else {
                            if appData.isDowngrading {
                                AppInfoSection
                            } else {
                                InputAppSection
                            }
                        }
                    }
                    .navigationTitle("WaffleStore")
                    .toolbar {
                        ToolbarItem(placement: .topBarTrailing) {
                            AppMenu(showHistoryView: $showHistoryView, showFavouritesView: $showFavouritesView, showDownloadedView: $showDownloadedView)
                        }
                        ToolbarItem(placement: .topBarTrailing) {
                            Button(action: {
                                showSettingsView.toggle()
                            }) {
                                Image(systemName: "gear")
                            }
                        }
                    }
                    .safeAreaInset(edge: .bottom) {
                        NavigationButtons()
                            .modifier(OverlayBackground())
                    }
                }
            }
        }
        .sheet(isPresented: $showSettingsView) {
            SettingsView()
        }
        .sheet(isPresented: $showSearchView) {
            AppSearchView()
        }
        .sheet(isPresented: $showHistoryView) {
            DowngradeHistoryView()
        }
        .sheet(isPresented: $showFavouritesView) {
            FavouritesView()
        }
        .sheet(isPresented: $showDownloadedView) { DownloadedAppsView() }
        .sheet(isPresented: $appData.showStoreVersions) { StoreVersionsView() }
        .sheet(item: $appData.downloadReady) { DownloadReadyView(record: $0) }
        .sheet(item: $appData.installationRequest) { OTAInstallationView(record: $0, startImmediately: true) }
        .onAppear { appData.restoreStoreAccount(); appData.restoreDownloadedIPA() }
    }
    
    private var LogsSection: some View {
        Section {
            if appData.isAuthenticated {
                Label(appData.applicationStatus, systemImage: appData.applicationIcon)
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
                if appData.showsDowngradeProgress {
                    VStack(alignment: .leading, spacing: 6) {
                        ProgressView(value: appData.downgradeProgress)
                        HStack {
                            Text(appData.downgradeProgressDetail)
                            Spacer()
                            Text("\(Int(appData.downgradeProgress * 100))%")
                        }
                        .font(.caption)
                        .foregroundStyle(.secondary)
                    }
                    .padding(.horizontal, 2)
                }
            }
            DisclosureGroup("Activity log", isExpanded: $showLogs) {
                LogView()
                    .frame(height: 140)
                    .modifier(TerminalPlatter())
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
    
    private var InputAppSection: some View {
        Section(header: HeaderLabel(text: "Downgrade App".localized, icon: "arrow.down.app"), footer: Text("Choose a version to download or install. Your IPA remains available for export.")) {
            VStack(spacing: 12) {
                TextField("App Store link, ID or bundle ID", text: $appData.appLink)
                    .modifier(TextFieldBackground())
                    .autocorrectionDisabled()
                    .textInputAutocapitalization(.never)
                
                Button(action: {
                    Haptic.shared.play(.soft)
                    showSearchView.toggle()
                }) {
                    ButtonLabel(text: "Search App Store".localized, icon: "magnifyingglass")
                }
                .buttonStyle(TranslucentButtonStyle())
            }
        }
    }
    
    private var AppInfoSection: some View {
        Section(header: HeaderLabel(text: "App Info".localized, icon: "info.circle")) {
            ItemInfoCell(label: "App Link".localized, icon: "link", text: appData.appLink)
            ItemInfoCell(label: "App Bundle ID".localized, icon: "shippingbox", text: appData.appBundleID)

            ItemInfoCell(label: "Target App Version".localized, icon: "arrow.down.app", text: appData.appVersion)
        }
    }
    
}

struct AppMenu: View {
    @EnvironmentObject var appData: AppData
    @Binding var showHistoryView: Bool
    @Binding var showFavouritesView: Bool
    @Binding var showDownloadedView: Bool
    
    var body: some View {
        Menu {
            Button("Downloaded apps", systemImage: "square.and.arrow.down") { showDownloadedView = true }
            Button(action: {
                showFavouritesView.toggle()
            }) {
                Label("Favourites".localized, systemImage: "star.fill")
            }
            .disabled(!appData.isAuthenticated || appData.isAuthenticating)
            
            Button(action: {
                showHistoryView.toggle()
            }) {
                Label("Downgrade History".localized, systemImage: "clock.arrow.circlepath")
            }
            .disabled(!appData.isAuthenticated || appData.isAuthenticating)
            
            Button(action: {
                if let url = appData.downloadedIPAURL { presentShareSheet(with: url) }
            }) {
                Label("Export IPA".localized, systemImage: "arrow.up.doc")
            }
            .disabled(!appData.hasAppBeenServed)
            
            Button(action: {
                Haptic.shared.play(.heavy)
                DispatchQueue.main.asyncAfter(deadline: .now() + 0.25) {
                    appData.logoutStoreAccount()
                }
            }) {
                ButtonLabel(text: "Log Out".localized, icon: "arrow.right")
            }
            .disabled(!appData.isAuthenticated || appData.isAuthenticating || appData.isDowngrading || appData.showStoreVersions)
        } label: {
            Image(systemName: "line.horizontal.3")
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

struct SidebarToggleModifier: ViewModifier {
    func body(content: Content) -> some View {
        if #available(iOS 17.0, *) {
            content
                .toolbar(removing: .sidebarToggle)
        } else {
            content
        }
    }
}

#Preview {
    ContentView()
        .environmentObject(AppData.shared)
}
