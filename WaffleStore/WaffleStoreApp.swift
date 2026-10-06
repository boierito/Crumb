//
//  WaffleStoreApp.swift
//  WaffleStore
//
//  Created by nxtcoreee3 on 6/20/26.
//

import SwiftUI
import UniformTypeIdentifiers


@main
struct CrumbApp: App {
    @StateObject private var appData = AppData.shared
    @StateObject private var localizationManager = LocalizationManager.shared
    
    @AppStorage("autoCleanApp") var autoCleanApp: Bool = true
    
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(appData)
                .environment(\.locale, .init(identifier: localizationManager.currentLanguage.rawValue))
                .onAppear {
                    if autoCleanApp {
                        cleanUp()
                    }
                }
                .onOpenURL { schemedURL in
                    let rawURL = schemedURL.absoluteString.replacingOccurrences(of: "crumb-authtest:", with: "")
                    if let appLink = rawURL.removingPercentEncoding {
                        appData.openAppSelection(appLink)
                        #if DEBUG
                        print("Successfully received app link! \(appLink)")
                        #endif
                    }
                }
        }
    }
}

extension String: @retroactive Error {}
