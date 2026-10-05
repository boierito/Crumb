//
//  NavigationButtons.swift
//  WaffleStore
//
//  Created by lunginspector on 2/24/26.
//

import SwiftUI
import PartyUI

struct NavigationButtons: View {
    @EnvironmentObject var appData: AppData
    
    var body: some View {
        VStack {
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
                        ButtonLabel(text: "Send 2FA Code".localized, icon: "key")
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
                        ButtonLabel(text: "Choose version / download IPA", icon: "square.and.arrow.down")
                    }
                    .buttonStyle(FancyButtonStyle())
                    .disabled(appData.appLink.isEmpty)
                    if let url = appData.downloadedIPAURL {
                        ShareLink(item: url) { Label("Export IPA", systemImage: "square.and.arrow.up") }
                    }
                    let currentAppId = extractAppId(from: appData.appLink)
                    let existingFav = appData.favourites.first { extractAppId(from: $0.appLink) == currentAppId }
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
    }
}

func extractAppId(from link: String) -> String {
    var parsed = link.components(separatedBy: "id").last ?? ""
    for char in parsed {
        if !char.isNumber {
            if let index = parsed.firstIndex(of: char) {
                parsed = String(parsed.prefix(upTo: index))
            }
            break
        }
    }
    return parsed
}

func fetchAppNameAndBundleId(forLink: String, completion: @escaping (String, String) -> Void) {
    let parsedAppId = extractAppId(from: forLink)
    guard !parsedAppId.isEmpty, let url = URL(string: "https://itunes.apple.com/lookup?id=\(parsedAppId)") else {
        completion("", "")
        return
    }
    
    URLSession.shared.dataTask(with: url) { data, response, error in
        guard let data = data, error == nil else {
            completion("", "")
            return
        }
        do {
            if let json = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any],
               let results = json["results"] as? [[String: Any]],
               let firstResult = results.first {
                let trackName = firstResult["trackName"] as? String ?? ""
                let bundleId = firstResult["bundleId"] as? String ?? ""
                completion(trackName, bundleId)
            } else {
                completion("", "")
            }
        } catch {
            completion("", "")
        }
    }.resume()
}

