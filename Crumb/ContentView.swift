import SwiftUI
import PartyUI
import MapleSAP
import UIKit

struct ContentView: View {
    @EnvironmentObject var appData: AppData
    @StateObject private var device = DeviceWorkflow.shared
    @Environment(\.scenePhase) private var scenePhase
    @State private var confirmLogout = false
    @State private var readyInstallation: DownloadRecord?
    @State private var scrollRequest = CrumbScrollRequest()
    @Environment(\.horizontalSizeClass) private var width
    @Environment(\.accessibilityReduceMotion) private var reduceMotion
    @FocusState private var codeFocused: Bool

    private let destinations: [(tab: AppTab, title: LocalizedStringKey, icon: String, selectedIcon: String)] = [
        (.search, "Search", "magnifyingglass", "magnifyingglass"),
        (.downloads, "Downloads", "arrow.down.circle", "arrow.down.circle.fill"),
        (.favourites, "Favourites", "star", "star.fill"),
        (.account, "Account", "person.crop.circle", "person.crop.circle.fill")
    ]

    var body: some View {
        HStack(spacing: 0) {
            if width == .regular {
                VStack(alignment: .leading, spacing: 6) {
                    HStack(spacing: 10) {
                        Image("CrumbMark").resizable().scaledToFit().frame(width: 36, height: 36)
                            .clipShape(RoundedRectangle(cornerRadius: 9)).accessibilityHidden(true)
                        Text("Crumb").font(.title2.weight(.bold)).accessibilityAddTraits(.isHeader)
                    }.padding(.bottom, 20)
                    ForEach(destinations.indices, id: \.self) { index in tabButton(index, compact: false) }
                    Spacer()
                }
                .padding(18)
                .frame(width: 216)
                .background(Color(uiColor: .secondarySystemGroupedBackground))
                Divider()
            }
            ZStack {
                page(.search) { AppSearchView() }
                page(.downloads) { DownloadedAppsView(embedded: true) }
                page(.favourites) { FavouritesView() }
                page(.account) { account }
            }
        }
        .safeAreaInset(edge: .bottom, spacing: 0) {
            if width != .regular {
                HStack(spacing: 2) {
                    ForEach(destinations.indices, id: \.self) { index in tabButton(index, compact: true) }
                }
                .padding(.horizontal, 8)
                .padding(.vertical, 4)
                .background(.regularMaterial)
                .overlay(alignment: .top) { Divider() }
            }
        }
        .environment(\.crumbScrollRequest, scrollRequest)
        .sheet(isPresented:$device.setupNeeded,onDismiss:{ if device.ready { device.continueQueued() } }) {
            NavigationStack { DeviceSetupView().toolbar { ToolbarItem(placement:.confirmationAction) { Button("Done") { device.setupNeeded = false } } } }
        }
        .alert("Could not finish", isPresented:Binding(get:{ !device.setupVisible && !device.installationSetupNeeded && device.signingReview == nil && !device.errorMessage.isEmpty },set:{ if !$0 { device.errorMessage = "" } })) {
            Button("OK") { device.errorMessage = "" }
        } message: { Text(device.errorMessage) }
        .onChange(of:scenePhase) { phase in
            if phase == .background { device.becameInactive() }
            if phase == .active { device.autoConnect() }
        }
        .sheet(isPresented:$device.installationSetupNeeded,onDismiss:{ device.installationSetupDismissed() }) {
            NavigationStack { SigningSettingsView().toolbar { ToolbarItem(placement:.confirmationAction) { Button("Done") { device.installationSetupNeeded = false } } } }
        }
        .sheet(item:$device.signingReview,onDismiss:{ device.signingReviewDismissed() }) { SigningReviewView(review:$0) }
        .sheet(item:$device.signedInstallation) { artifact in
            OTAInstallationView(record:artifact.record,startImmediately:true,source:artifact.source,cleanup:artifact.workspace)
        }
        .sheet(isPresented: $appData.showStoreVersions) { StoreVersionsView() }
        .sheet(item: $appData.downloadReady, onDismiss: finishDownloadReview) { record in
            DownloadReadyView(record: record) { readyInstallation = record }
        }
        .sheet(item: $appData.installationRequest) { OTAInstallationView(record: $0, startImmediately: true) }
        .onAppear { appData.restoreStoreAccount(); appData.restoreDownloadedIPA() }
        .onChange(of: appData.downloadPhaseLabel) { phase in
            if appData.isDowngrading && appData.selectedTab == .downloads {
                UIAccessibility.post(notification: .announcement, argument: phase)
            }
        }
        .onChange(of: appData.isAuthenticated) { authenticated in
            if authenticated, appData.openVersionsAfterLogin {
                appData.openVersionsAfterLogin = false
                appData.openAppSelection(appData.appLink, country: appData.selectedCatalogCountry, resuming: true)
            }
        }
    }

    private func finishDownloadReview() {
        guard let record = readyInstallation else { return }
        readyInstallation = nil
        device.install(record)
    }

    private func page<Content: View>(_ tab: AppTab, @ViewBuilder content: () -> Content) -> some View {
        let selected = appData.selectedTab == tab
        return content()
            .frame(maxWidth: 920)
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .background(Color(uiColor: .systemGroupedBackground))
            .opacity(selected ? 1 : 0)
            .allowsHitTesting(selected)
            .accessibilityHidden(!selected)
    }

    private func tabButton(_ index: Int, compact: Bool) -> some View {
        let item = destinations[index]
        let selected = appData.selectedTab == item.tab
        let activeDownload = item.tab == .downloads && (appData.isDowngrading || device.activeRecord != nil)
        let layout = compact ? AnyLayout(VStackLayout(spacing: 3)) : AnyLayout(HStackLayout(spacing: 12))
        return Button {
            if selected { scrollRequest = CrumbScrollRequest(tab: item.tab, count: scrollRequest.count + 1) }
            else {
                withAnimation(reduceMotion ? nil : .easeOut(duration: 0.18)) { appData.selectedTab = item.tab }
            }
        } label: {
            layout {
                Image(systemName: selected ? item.selectedIcon : item.icon)
                    .font(.system(size: 19, weight: selected ? .semibold : .regular))
                    .frame(width: 28, height: 23)
                    .overlay(alignment: .topTrailing) {
                        if activeDownload { Circle().fill(Color.accentColor).frame(width: 7, height: 7) }
                    }
                    .accessibilityHidden(true)
                Text(item.title)
                    .font(compact ? .caption2.weight(selected ? .semibold : .medium) : .subheadline.weight(.medium))
                    .multilineTextAlignment(.center)
            }
            .padding(.horizontal, compact ? 1 : 10)
            .padding(.vertical, 4)
            .frame(maxWidth: .infinity, minHeight: 48, alignment: compact ? .center : .leading)
            .foregroundStyle(selected ? Color.primary : Color.secondary)
            .background(selected ? Color.accentColor.opacity(0.14) : Color.clear, in: RoundedRectangle(cornerRadius: 12))
            .contentShape(Rectangle())
        }
        .buttonStyle(.plain)
        .accessibilityAddTraits(selected ? .isSelected : [])
        .accessibilityValue(activeDownload ? Text("Activity in progress") : Text(""))
    }

    private var account: some View {
        NavigationStack {
            List {
                if appData.isAuthenticated {
                    Section {
                        HStack(spacing: 12) {
                            CrumbArtwork(symbol: "person.crop.circle", size: 54)
                            VStack(alignment: .leading, spacing: 4) {
                                Text(appData.appleId).font(.headline)
                                Text(regionName(appData.accountCountry)).font(.subheadline).foregroundStyle(.secondary)
                            }
                        }.padding(.vertical, 6).id(AppTab.account)
                    } header: { Text("Apple account") }
                    Section {
                        Button("Sign out", role: .destructive) { confirmLogout = true }
                            .disabled(appData.isAuthenticating || appData.isDowngrading || appData.storeRequestCount > 0 || appData.showStoreVersions)
                    } footer: {
                        Text("Signing out keeps your favourites and downloaded IPAs.")
                    }
                } else {
                    loginSection
                }
                Section {
                    NavigationLink { SettingsView() } label: { Label("Settings", systemImage: "gearshape") }
                }
            }
            .crumbPrimaryPage(.account)
            .navigationTitle("Account")
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

    private var loginSection: some View {
        Group {
            if appData.hasSent2FACode {
                Section {
                    TextField("Six-digit code", text: $appData.code)
                        .font(.title2.monospacedDigit())
                        .keyboardType(.numberPad)
                        .textContentType(.oneTimeCode)
                        .focused($codeFocused)
                        .disabled(appData.isAuthenticating)
                    NavigationButtons()
                    if let deadline = appData.codeResendAvailableAt, appData.isCodeResendCoolingDown {
                        HStack {
                            Text("Request a new code in")
                            Text(timerInterval: Date()...max(Date(), deadline), countsDown: true)
                                .monospacedDigit()
                        }.font(.subheadline).foregroundStyle(.secondary)
                    } else {
                        Button("Request new code") { appData.requestNewVerificationCode() }
                            .disabled(appData.isAuthenticating || !appData.signInForm.canRequestNewCode)
                    }
                    Button("Change Apple ID") { appData.changeAppleAccount() }
                        .disabled(appData.isAuthenticating)
                } header: {
                    Text("Verify your account").id(AppTab.account)
                } footer: {
                    Text("Check a trusted Apple device for a sign-in request or get a verification code in its Apple Account settings.")
                }
            } else {
                Section {
                    TextField("Apple ID".localized, text: $appData.appleId)
                        .autocorrectionDisabled().textInputAutocapitalization(.never)
                        .keyboardType(.emailAddress).textContentType(.username)
                        .disabled(appData.isAuthenticating)
                    HStack {
                        Group {
                            if appData.showPassword { TextField("Password".localized, text: $appData.password) }
                            else { SecureField("Password".localized, text: $appData.password) }
                        }
                        .textContentType(.password)
                        .autocorrectionDisabled().textInputAutocapitalization(.never)
                        .disabled(appData.isAuthenticating)
                        Button { appData.showPassword.toggle() } label: {
                            Image(systemName: appData.showPassword ? "eye.slash" : "eye").frame(width: 44, height: 44)
                        }
                        .buttonStyle(.borderless)
                        .accessibilityLabel(appData.showPassword ? "Hide password" : "Show password")
                    }
                    NavigationButtons()
                } header: {
                    Text("Sign in to Apple").id(AppTab.account)
                } footer: {
                    Text("Sign in to choose and download App Store versions.")
                }
            }
            if appData.isAuthenticating {
                Section {
                    Label(appData.applicationStatus, systemImage: "hourglass")
                        .foregroundStyle(.secondary)
                    Button("Cancel") { appData.cancelAppleLogin() }
                }
            }
            if !appData.authenticationError.isEmpty {
                Label(appData.authenticationError, systemImage: "exclamationmark.circle")
                    .font(.subheadline).foregroundStyle(.red).textSelection(.enabled)
            }
        }
        .onChange(of: appData.hasSent2FACode) { codeFocused = $0 }
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
