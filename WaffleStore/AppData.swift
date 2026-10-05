// xcode: set sdk=iOS

//
//  AppData.swift
//  WaffleStore
//
//  Created by lunginspector on 2/25/26.
//

import SwiftUI
import Combine
import MapleSAP

@MainActor
final class AppData: ObservableObject {
    static let shared = AppData()
    
    @Published var applicationIcon: String = "xmark.circle.fill"
    @Published var applicationIconColor: Color = .secondary
    @Published var applicationStatus: String = "Not logged in!".localized
    @Published var downgradeProgress: Double = 0
    @Published var downgradeProgressDetail: String = ""
    @Published var showsDowngradeProgress: Bool = false
    
    @Published var appBundleID: String = ""
    @Published var appVersion: String = ""
    
    @Published var hasAppBeenServed: Bool = false
    
    @Published var ipaTool: IPATool?
    
    @Published var appleId: String = ""
    @Published var password: String = ""
    @Published var code: String = ""
    
    @Published var isAuthenticated: Bool = false
    @Published var isAuthenticating: Bool = false
    @Published var authenticationError: String = ""
    @Published var authenticationRecovery: String = ""
    @Published var authenticationDiagnostic: String = "No authentication attempted."
    var authenticationTask: Task<Void, Never>?
    var didRestoreStoreAccount = false
    var pendingAuthenticationCookies: [StoreCookie] = []
    @Published var showStoreVersions = false
    @Published var downloadedIPAURL: URL?
    @Published var completedDownloads: [DownloadRecord] = []
    @Published var storeError = ""
    @Published var storeDiagnostic = "No Store operation attempted."
    var storeTask: Task<Void, Never>?
    @Published var isDowngrading: Bool = false
    
    @Published var appLink: String = ""
    
    @Published var hasSent2FACode: Bool = false
    
    @Published var showPassword: Bool = false
    
    @Published var showFavouritesView: Bool = false
    
    @Published var downgradeHistory: [DowngradeHistoryEntry] = DowngradeHistoryStore.load()
    
    @Published var favourites: [FavouriteApp] = FavouritesStore.load()
}
