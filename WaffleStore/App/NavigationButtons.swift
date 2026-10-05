//
//  NavigationButtons.swift
//  WaffleStore
//
//  Created by lunginspector on 2/24/26.
//

import SwiftUI
import PartyUI
import MapleSAP

struct NavigationButtons: View {
    @EnvironmentObject var appData: AppData
    
    @State private var installation: DownloadRecord?

    var body: some View {
        VStack(spacing: 10) {
            // i hate this.
            if !appData.isAuthenticated {
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
            } else {
                if appData.isDowngrading {
                    Button("Cancel download") { appData.storeTask?.cancel() }
                        .buttonStyle(FancyButtonStyle())
                } else {
                    Button(action: { appData.showStoreVersions = true }) {
                        ButtonLabel(text: "Choose version", icon: "square.and.arrow.down")
                    }
                    .buttonStyle(FancyButtonStyle())
                    .disabled(appData.appLink.isEmpty)
                    if let record = appData.completedDownloads.first, let url = record.fileURL, url == appData.downloadedIPAURL {
                        VStack(spacing: 6) {
                            Text("\(record.appName) · \(record.version)")
                                .font(.caption).foregroundStyle(.secondary).lineLimit(1)
                            HStack(spacing: 10) {
                                Button { installation = record } label: {
                                    ButtonLabel(text: "Install", icon: "arrow.down.app")
                                }
                                .accessibilityLabel("Install \(record.appName), version \(record.version)")
                                ShareLink(item: url) { ButtonLabel(text: "Export", icon: "square.and.arrow.up") }
                            }
                            .buttonStyle(TranslucentButtonStyle())
                        }
                    } else if let url = appData.downloadedIPAURL {
                        ShareLink(item: url) { ButtonLabel(text: "Export IPA", icon: "square.and.arrow.up") }
                            .buttonStyle(TranslucentButtonStyle())
                    }
                    let currentAppId = extractAppId(from: appData.appLink)
                    let existingFav = appData.favourites.first { currentAppId.isEmpty ? $0.appLink == appData.appLink : extractAppId(from: $0.appLink) == currentAppId }
                    let isFavourited = existingFav != nil

                    Button(action: {
                        Haptic.shared.play(.soft)
                        DispatchQueue.main.asyncAfter(deadline: .now() + 0.15) {
                            if appData.appLink.isEmpty {
                                Alertinator.shared.alert(title: "Cannot Add to Favourites".localized, body: "Please enter an app link first.".localized)
                                return
                            }
                            
                            if isFavourited {
                                if let fav = existingFav {
                                    appData.favourites = FavouritesStore.remove(fav)
                                    Haptic.shared.play(.soft)
                                    Alertinator.shared.alert(title: "Removed from Favourites".localized, body: "This app has been removed from your favourites.".localized)
                                }
                            } else {
                                let currentLink = appData.appLink
                                fetchAppNameAndBundleId(forLink: currentLink) { trackName, bundleId in
                                    DispatchQueue.main.async {
                                        let resolvedName = trackName.isEmpty ? (bundleId.isEmpty ? "App Store Link" : bundleId) : trackName
                                        let resolvedBundle = bundleId.isEmpty ? appData.appBundleID : bundleId
                                        let favourite = FavouriteApp(
                                            appLink: currentLink,
                                            bundleId: resolvedBundle,
                                            appName: resolvedName
                                        )
                                        appData.favourites = FavouritesStore.add(favourite)
                                        Haptic.shared.play(.soft)
                                        Alertinator.shared.alert(title: "Added to Favourites".localized, body: "This app has been added to your favourites for quick access.".localized)
                                    }
                                }
                            }
                        }
                    }) {
                        ButtonLabel(text: isFavourited ? "Remove from Favourites".localized : "Add to Favourites".localized, icon: isFavourited ? "star.fill" : "star")
                    }
                    .buttonStyle(FancyButtonStyle(color: .mint))
                    .disabled(appData.appLink.isEmpty)
                    

                }
            }
        }
        .sheet(item: $installation) { OTAInstallationView(record: $0, startImmediately: true) }
    }
}

func extractAppId(from link: String) -> String {
    let trimmed = link.trimmingCharacters(in: .whitespacesAndNewlines)
    if let id = StoreParsing.identifier(trimmed) { return id }
    guard let url = URL(string: trimmed), url.scheme == "https",
          ["apps.apple.com", "itunes.apple.com"].contains(url.host?.lowercased() ?? "") else { return "" }
    return url.pathComponents.compactMap { $0.hasPrefix("id") ? StoreParsing.identifier(String($0.dropFirst(2))) : nil }.last ?? ""
}

func fetchAppNameAndBundleId(forLink: String, completion: @escaping (String, String) -> Void) {
    guard let tool = AppData.shared.ipaTool else { completion("", ""); return }
    Task {
        do {
            let app = try await tool.lookup(forLink)
            completion(app.name, app.bundleID)
        } catch { completion("", "") }
    }
}
