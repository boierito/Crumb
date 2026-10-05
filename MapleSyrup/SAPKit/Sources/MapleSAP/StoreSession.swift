// Native Swift Store transport adapted from ipatool 3411d57 (MIT).
import Foundation
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

public actor StoreSession {
    private let account: StoreAccount
    private let identity: MachineIdentity
    private let transport: AuthenticationTransport
    private let generator: KBSyncGenerator
    private let persistence: KBSyncPersistence
    private var bag: [String: String]?
    private var busy = false
    private let sleep: (TimeInterval) async throws -> Void
    private let progress: (StoreStage) async -> Void
    private let diagnostic: (String) async -> Void
    public init(account: StoreAccount, identity: MachineIdentity, transport: AuthenticationTransport,
                generator: KBSyncGenerator, persistence: KBSyncPersistence = NoKBSyncPersistence(),
                progress: @escaping (StoreStage) async -> Void = { _ in },
                diagnostic: @escaping (String) async -> Void = { _ in },
                sleep: @escaping (TimeInterval) async throws -> Void = { try await Task.sleep(nanoseconds: UInt64($0 * 1_000_000_000)) }) throws {
        try account.validate(identity: identity)
        self.account = account; self.identity = identity; self.transport = transport
        self.generator = generator; self.persistence = persistence; self.sleep = sleep; self.progress = progress; self.diagnostic = diagnostic
    }
    public func lookup(_ input: String) async throws -> StoreApp {
        let trimmed = input.trimmingCharacters(in: .whitespacesAndNewlines)
        var id = StoreParsing.identifier(trimmed)
        if let url = URL(string: trimmed), let host = url.host {
            guard ["apps.apple.com", "itunes.apple.com"].contains(host.lowercased()), url.scheme == "https" else { throw StoreError.invalidApp }
            id = url.pathComponents.compactMap { $0.hasPrefix("id") ? StoreParsing.identifier(String($0.dropFirst(2))) : nil }.last
            guard id != nil else { throw StoreError.invalidApp }
        }
        let bundle = id == nil ? trimmed : nil
        guard id != nil || (bundle?.contains(".") == true && !trimmed.contains("/") && !trimmed.contains(" ")) else { throw StoreError.invalidApp }
        let country = try Storefront.country(account.storefront)
        let data = try await get("https://itunes.apple.com/lookup", query: [id == nil ? "bundleId" : "id": id ?? trimmed,
            "country": country, "entity": "software,iPadSoftware", "limit": "1"])
        guard let root = try JSONSerialization.jsonObject(with: data) as? [String: Any],
              let result = (root["results"] as? [[String: Any]])?.first,
              let appID = StoreParsing.identifier(result["trackId"]), let bundleID = result["bundleId"] as? String,
              let name = result["trackName"] as? String, !bundleID.isEmpty,
              id == nil || appID == id, bundle == nil || bundleID == bundle else { throw StoreError.unavailable }
        return StoreApp(id: appID, bundleID: bundleID, name: name, price: (result["price"] as? NSNumber)?.doubleValue)
    }
    public func descriptor(app: StoreApp, externalVersionID: String = "", acquireFreeLicense: Bool = true) async throws -> StoreDownload {
        guard !busy else { throw SAPError.invalidState }
        busy = true; defer { busy = false }
        let version = externalVersionID.isEmpty ? try await latestVersion(app) : externalVersionID
        guard StoreParsing.identifier(version) != nil else { throw StoreError.invalidApp }
        try await resolveBag()
        do { return try await requestDescriptor(app, version: version) }
        catch StoreError.licenseRequired where acquireFreeLicense {
            try await purchaseFree(app)
            return try await requestDescriptor(app, version: version)
        }
    }
    private func latestVersion(_ app: StoreApp) async throws -> String {
        await progress(.latest)
        let country = try Storefront.country(account.storefront)
        for platform in ["enterprisestore", "iphone", "ipad"] {
            let data = try await get("https://uclient-api.itunes.apple.com/WebObjects/MZStorePlatform.woa/wa/lookup", query:
                ["version": "2", "id": app.id, "p": "mdm-lockup", "caller": "MDM", "platform": platform, "cc": country, "l": "en"])
            guard let root = try JSONSerialization.jsonObject(with: data) as? [String: Any],
                  let result = (root["results"] as? [String: Any])?[app.id] as? [String: Any],
                  result["bundleId"] as? String == app.bundleID,
                  let offer = (result["offers"] as? [[String: Any]])?.first else { continue }
            if let version = offer["version"] as? [String: Any], let id = StoreParsing.identifier(version["externalId"]) { return id }
            if let text = offer["buyParams"] as? String,
               let id = URLComponents(string: "https://unused.invalid/?" + text)?.queryItems?.first(where: { $0.name == "appExtVrsId" })?.value,
               StoreParsing.identifier(id) != nil { return id }
        }
        throw StoreError.unavailable
    }
    private func resolveBag() async throws {
        if bag != nil { return }
        await progress(.bag)
        let data = try await get("https://init.itunes.apple.com/bag.xml", query: ["guid": identity.guid])
        guard let root = try ApplePlist.dictionary(data), let urls = root["urlBag"] as? [String: Any] else { throw SAPError.invalidBag }
        bag = urls.compactMapValues { $0 as? String }
    }
    private func requestDescriptor(_ app: StoreApp, version: String) async throws -> StoreDownload {
        var preferredError: Error?
        if let endpoint = bag?["volumeStoreDownloadProduct"], URL(string: endpoint)?.path == "/WebObjects/DownloadDispatch.woa/wa/ent/download" {
            do {
                let url = try dispatchURL(endpoint, path: "/WebObjects/DownloadDispatch.woa/wa/ent/download")
                if let cached = try persistence.load(dsid: account.dsid, guid: account.guid) {
                    do { return try await ent(url, app: app, version: version, blob: cached) }
                    catch {
                        try Task.checkCancellation()
                        await diagnostic("recovery=rejected-cached-kbsync; category=\(ResponseDiagnostic.category(error))")
                        try? persistence.clear()
                    }
                }
                guard let dsid = UInt64(account.dsid), dsid > 0 else { throw AuthenticationError.invalidSession }
                await progress(.kbsync)
                let blob = try generator.generate(identity: identity, dsid: dsid)
                try Task.checkCancellation()
                let result = try await ent(url, app: app, version: version, blob: blob)
                // Cache failures cannot invalidate an already validated Apple reply.
                try? persistence.save(blob, dsid: account.dsid, guid: account.guid)
                return result
            } catch {
                try Task.checkCancellation()
                await diagnostic("recovery=ent-to-pod; category=\(ResponseDiagnostic.category(error)); saved-session=retained")
                preferredError = error
            }
        }
        // ipatool's legacy fallback derives the host from the authenticated pod.
        // Bag endpoints remain authoritative for ent, redownload and update.
        let prefix = account.pod.map { "p\($0)-" } ?? ""
        let legacy = URL(string: "https://\(prefix)buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/volumeStoreDownloadProduct")!
        do {
            await progress(.legacy)
            let root = try await post(legacy, body: payload(app, version: version, key: "externalVersionId"), token: false)
            return try StoreParsing.download(root, app: app, version: version, email: account.email)
        } catch StoreError.unavailable { }
        catch { throw error }
        guard let redownload = bag?["redownloadProduct"] else { throw preferredError ?? StoreError.unavailable }
        do {
            await progress(.redownload)
            let root = try await post(dispatchURL(redownload, path: "/r/redownload"), body: payload(app, version: version, key: "appExtVrsId"), token: false)
            return try StoreParsing.download(root, app: app, version: version, email: account.email)
        } catch let error as StoreError where error == .unavailable || error == .emptyRedownload {
            guard let update = bag?["updateProduct"] else { throw error }
            await progress(.update)
            let root = try await post(dispatchURL(update, path: "/up/updateProduct"), body: payload(app, version: version, key: "appExtVrsId"), token: false)
            return try StoreParsing.download(root, app: app, version: version, email: account.email)
        }
    }
    private func ent(_ url: URL, app: StoreApp, version: String, blob: Data) async throws -> StoreDownload {
        guard !blob.isEmpty else { throw StoreError.native(8) }
        var body = payload(app, version: version, key: "externalVersionId")
        body["salableAdamId"] = app.id
        body["kbsync"] = blob.base64EncodedString()
        body["serialNumber"] = (Data([0x54, 0xc8, 0xb0, 0xa9, 0x88]) + identity.hardwareID.dropFirst(2)).base64EncodedString()
        await progress(.ent)
        let root = try await post(url, body: body, token: true, ent: true)
        return try StoreParsing.download(root, app: app, version: version, email: account.email)
    }
    private func payload(_ app: StoreApp, version: String, key: String) -> [String: Any] {
        ["creditDisplay": "", "guid": identity.guid, "salableAdamId": NSNumber(value: UInt64(app.id) ?? 0), "serialNumber": "0", key: version]
    }
    private func purchaseFree(_ app: StoreApp) async throws {
        guard app.price == 0 else { throw StoreError.paidPurchase }
        guard let endpoint = bag?["buyProduct"], var components = URLComponents(string: endpoint),
              let url = components.url else { throw SAPError.invalidBag }
        _ = try SAPConfiguration.trustedAppleURL(endpoint)
        let host = url.host?.lowercased() ?? ""
        guard host == "buy.itunes.apple.com" || host.hasSuffix("-buy.itunes.apple.com"),
              ["/WebObjects/MZBuy.woa/wa/buyProduct", "/WebObjects/MZFinance.woa/wa/buyProduct"].contains(url.path),
              components.query == nil else { throw SAPError.invalidEndpoint }
        if let pod = account.pod { components.host = "p\(pod)-buy.itunes.apple.com" }
        guard let purchaseURL = components.url else { throw SAPError.invalidEndpoint }
        let body: [String: Any] = ["appExtVrsId": "0", "hasAskedToFulfillPreorder": "true", "buyWithoutAuthorization": "true",
            "hasDoneAgeCheck": "true", "guid": identity.guid, "needDiv": "0", "origPage": "Software-\(app.id)",
            "origPageLocation": "Buy", "price": "0", "pricingParameters": "STDQ", "productType": "C", "salableAdamId": NSNumber(value: UInt64(app.id) ?? 0)]
        await progress(.purchase)
        let result = try await post(purchaseURL, body: body, token: true, retry: false)
        if StoreParsing.identifier(result["failureType"]) == "5002" { return }
        try StoreParsing.failure(result)
        guard result["jingleDocType"] as? String == "purchaseSuccess", StoreParsing.identifier(result["status"]) == "0" else { throw StoreError.licenseRequired }
    }
    private func dispatchURL(_ text: String, path: String) throws -> URL {
        guard let c = URLComponents(string: text), let url = c.url, c.scheme == "https", c.host == "downloaddispatch.itunes.apple.com",
              c.port == nil, c.user == nil, c.password == nil, c.query == nil, c.fragment == nil,
              c.percentEncodedPath == path else { throw SAPError.invalidEndpoint }
        return url
    }
    private func post(_ url: URL, body: [String: Any], token: Bool, ent: Bool = false, retry: Bool = true) async throws -> [String: Any] {
        var c = URLComponents(url: url, resolvingAgainstBaseURL: false)!
        c.queryItems = [URLQueryItem(name: "guid", value: identity.guid)]
        var request = URLRequest(url: c.url!)
        request.httpMethod = "POST"
        request.httpBody = try PropertyListSerialization.data(fromPropertyList: body, format: .xml, options: 0)
        request.setValue(ent ? "application/x-www-form-urlencoded; charset=utf-8" : "application/x-apple-plist", forHTTPHeaderField: "Content-Type")
        request.setValue(account.dsid, forHTTPHeaderField: "iCloud-DSID")
        request.setValue(account.dsid, forHTTPHeaderField: "X-Dsid")
        if token {
            request.setValue(account.storefront, forHTTPHeaderField: "X-Apple-Store-Front")
            request.setValue(account.passwordToken, forHTTPHeaderField: "X-Token")
        }
        let data = try await send(request, retry: retry, secrets: [body["kbsync"] as? String ?? ""])
        guard var root = try ApplePlist.dictionary(data) else { throw StoreError.invalidResponse }
        if var message = root["customerMessage"] as? String {
            let secrets = [account.passwordToken, account.dsid, account.guid, account.email, body["kbsync"] as? String ?? ""]
            for secret in secrets.filter({ !$0.isEmpty }).sorted(by: { $0.count > $1.count }) {
                message = message.replacingOccurrences(of: secret, with: "[redacted]")
            }
            root["customerMessage"] = String(String.UnicodeScalarView(message.unicodeScalars.filter { !CharacterSet.controlCharacters.contains($0) }.prefix(400)))
        }
        return root
    }
    private func get(_ base: String, query: [String: String]) async throws -> Data {
        var c = URLComponents(string: base)!
        c.queryItems = query.sorted { $0.key < $1.key }.map { URLQueryItem(name: $0.key, value: $0.value) }
        return try await send(URLRequest(url: c.url!))
    }
    private func send(_ request: URLRequest, retry: Bool = true, secrets: [String] = []) async throws -> Data {
        var request = request
        let ent = request.url?.path == "/WebObjects/DownloadDispatch.woa/wa/ent/download"
        request.setValue(ent ? "Configurator/2.18 (Macintosh; OS X 15.3.2; 24D81) AppleWebKit/0620.2.4.11.6" : SAPProtocol.userAgent, forHTTPHeaderField: "User-Agent")
        for attempt in 0..<(retry ? 3 : 1) {
            try Task.checkCancellation()
            do {
                if let url = request.url { await diagnostic(await transport.cookieDiagnostic(for: url)) }
                let (data, response) = try await transport.send(request)
                let scope: String
                switch request.url?.path {
                case "/WebObjects/DownloadDispatch.woa/wa/ent/download": scope = "ent"
                case "/WebObjects/MZFinance.woa/wa/volumeStoreDownloadProduct": scope = "pod"
                case "/r/redownload": scope = "redownload"
                case "/up/updateProduct": scope = "update"
                case "/WebObjects/MZBuy.woa/wa/buyProduct", "/WebObjects/MZFinance.woa/wa/buyProduct": scope = "purchase"
                case "/bag.xml": scope = "bag"
                default: scope = "catalog"
                }
                await diagnostic(ResponseDiagnostic.response(data, status: response.statusCode, scope: scope, attempt: attempt + 1,
                    secrets: secrets + [account.passwordToken, account.dsid, account.guid, account.email]))
                // A bare HTTP rejection is not proof that the account token expired.
                // Decode populated Apple errors before generic HTTP classification.
                if let root = try? ApplePlist.dictionary(data),
                   !(StoreParsing.identifier(root["failureType"]) ?? "").isEmpty || !(root["customerMessage"] as? String ?? "").isEmpty {
                    return data
                }
                if response.statusCode == 200 { return data }
                if response.statusCode == 500 && data.isEmpty { throw StoreError.emptyRedownload }
                guard retry, attempt < 2, ([204, 404, 429].contains(response.statusCode) || response.statusCode / 100 == 5) else { throw StoreError.http(response.statusCode) }
                let delay = try StoreRetry.delay(response.value(forHTTPHeaderField: "Retry-After"), attempt: attempt)
                try await sleep(delay)
            } catch let error as URLError where retry && attempt < 2 && [.timedOut, .networkConnectionLost, .cannotConnectToHost].contains(error.code) {
                try await sleep(Double((attempt + 1) * 5))
            }
        }
        throw StoreError.invalidResponse
    }
}
