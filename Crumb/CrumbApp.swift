//
//  CrumbApp.swift
//  Crumb
//
//  Created by nxtcoreee3 on 6/20/26.
//

import SwiftUI
import UniformTypeIdentifiers


@main
struct CrumbApp: App {
    @StateObject private var appData = AppData.shared
    @StateObject private var localizationManager = LocalizationManager.shared
    @StateObject private var sap = CrumbSAPSigner.shared
    
    @AppStorage("autoCleanApp") var autoCleanApp: Bool = true
    
    var body: some Scene {
        WindowGroup {
            ContentView()
                .overlay(alignment: .bottomTrailing) {
                    if sap.isPresented, let webView = sap.webView {
                        CrumbSAPWebView(webView: webView)
                            .frame(width: 4, height: 4)
                            .opacity(0.02)
                            .allowsHitTesting(false)
                            .accessibilityHidden(true)
                    }
                }
                .environmentObject(appData)
                .environment(\.locale, .init(identifier: localizationManager.currentLanguage.rawValue))
                .onAppear {
                    if autoCleanApp {
                        cleanUp()
                    }
                    // Load iPadN WebKit runtime ahead of the next login.
                    if !appData.isAuthenticated { sap.prewarm() }
                }
                .onOpenURL { schemedURL in
                    if ["device-result", "device-cancelled"].contains(schemedURL.host ?? "") {
                        if DeviceWorkflow.shared.receiveResultURL(schemedURL) { appData.selectedTab = .downloads }
                        return
                    }
                    let rawURL = schemedURL.absoluteString.replacingOccurrences(of: "crumb:", with: "")
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
