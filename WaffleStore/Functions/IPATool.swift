// Modern Swift Store facade. The Go library supplies only interpreted guest
// operations and streaming ZIP preparation; no ipatool executable is bundled.
import Foundation
import MapleSAP

@MainActor
final class IPATool {
    let account: StoreAccount
    private var store: StoreSession?
    private var transport: AppleAuthenticationTransport?
    init(account: StoreAccount) { self.account = account }
    var appleId: String { account.email }
    func close() { transport?.close(); transport = nil; store = nil }
    private func session() throws -> StoreSession {
        if let store = store { return store }
        let transport = AppleAuthenticationTransport(cookies: account.cookies)
        let store = try StoreSession(account: account, identity: KeychainMachineIdentity.loadOrCreate(),
            transport: transport, generator: NativeKBSyncGenerator(), persistence: KeychainKBSync(), progress: { stage in
                await MainActor.run {
                    switch stage {
                    case .purchase: AppData.shared.applicationStatus = "Obtaining free app…"
                    case .bag, .latest: AppData.shared.applicationStatus = "Finding available versions…"
                    default: AppData.shared.applicationStatus = "Preparing download…"
                    }
                    print("Apple Store stage: \(stage.rawValue)")
                }
            })
        self.transport = transport; self.store = store
        return store
    }
    func lookup(_ input: String) async throws -> StoreApp { try await session().lookup(input) }
    func descriptor(app: StoreApp, version: String = "") async throws -> StoreDownload {
        try await session().descriptor(app: app, externalVersionID: version)
    }
}
