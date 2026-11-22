// Modern Swift Store facade. The local Go library supplies only streaming IPA
// inspection and packaging; SAP and KBSync use iPadN WebKit/WASM.
import Foundation
import MapleSAP

@MainActor
final class IPATool {
    private(set) var account: StoreAccount
    private var store: StoreSession?
    private var transport: AppleAuthenticationTransport?
    init(account: StoreAccount) { self.account = account }
    var appleId: String { account.email }
    func close() { transport?.close(); transport = nil; store = nil }
    private func session() throws -> StoreSession {
        if let store = store { return store }
        let transport = AppleAuthenticationTransport(cookies: account.cookies)
        let store = try StoreSession(account: account, identity: KeychainMachineIdentity.loadOrCreate(),
            transport: transport, generator: CrumbSAPKBSyncGenerator(), persistence: KeychainKBSync(),
            accountPersistence: KeychainStoreAccount(), progress: { stage in
                await MainActor.run {
                    switch stage {
                    case .purchase: AppData.shared.applicationStatus = "Obtaining free app…"
                    case .bag, .latest: AppData.shared.applicationStatus = "Finding available versions…"
                    default: AppData.shared.applicationStatus = "Preparing download…"
                    }
                    #if DEBUG
                    print("Apple Store stage: \(stage.rawValue)")
                    #endif
                }
            })
        self.transport = transport; self.store = store
        return store
    }
    func lookup(_ input: String, country: String? = nil) async throws -> StoreApp { try await session().lookup(input, country: country) }
    func descriptor(app: StoreApp, version: String = "") async throws -> StoreDownload {
        AppData.shared.storeRequestCount += 1
        defer { AppData.shared.storeRequestCount -= 1 }
        do {
            let current = try session()
            let result = try await current.descriptor(app: app, externalVersionID: version)
            // The validated -128 Apple response may correct the storefront.
            account = await current.currentAccount()
            return result
        }
        catch {
            if let storeError = error as? StoreError, storeError == .sessionExpired {
                // User interaction is required: passwords and 2FA are not saved.
                // Keep Keychain credentials intact instead of inventing a token.
                AppData.shared.authenticationError = storeError.localizedDescription
                AppData.shared.selectedTab = .account
            }
            #if DEBUG
            print("Apple Store failure: category=\(ResponseDiagnostic.category(error))")
            #endif
            throw error
        }
    }
}
