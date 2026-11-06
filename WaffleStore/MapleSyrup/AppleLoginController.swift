import Foundation
import Security
import SwiftUI
import PartyUI
import MapleSAP

// MainActor owns UI state; SAPSession runs the blocking native guest away from it.
extension AppData {
    func startAppleLogin(_ action: AppleSignInAction = .submit) {
        guard !isAuthenticating else { return }
        let input: AppleSignInSubmission
        do { input = try signInForm.begin(action) }
        catch {
            signInForm.failed(error)
            authenticationError = authenticationMessage(error)
            return
        }
        if action == .requestNewCode { clearLoginPreparation() }
        isAuthenticating = true
        authenticationError = ""
        authenticationTask = Task {
            defer { isAuthenticating = false; authenticationTask = nil }
            var preparation: PreparedAppleLogin?
            var keepPrepared = false
            do {
                try Task.checkCancellation()
                let prepared: PreparedAppleLogin
                if let existing = preparedAppleLogin, existing.canReuse(for: input.email) {
                    prepared = existing
                } else {
                    if let existing = preparedAppleLogin { await existing.close() }
                    preparedAppleLogin = nil
                    prepared = try PreparedAppleLogin(email: input.email, cookies: input.cookies)
                    preparedAppleLogin = prepared
                }
                preparation = prepared
                let signer = try await prepared.prepare { self.setAuthenticationStage($0) }
                try Task.checkCancellation()
                let authentication = AppleAuthentication(transport: prepared.loginTransport, signer: signer,
                    persistence: KeychainStoreAccount(), automaticRecovery: true)
                guard let endpoint = prepared.endpoint else { throw CancellationError() }
                let outcome = try await authentication.login(email: input.email, password: input.password, code: input.code,
                    identity: prepared.identity, endpoint: endpoint, resolvedEndpoint: { endpoint in
                        await MainActor.run { prepared.endpoint = endpoint }
                    }) { stage in
                        await MainActor.run { self.setAuthenticationStage(stage) }
                    }
                try Task.checkCancellation()
                switch outcome {
                case .twoFactorRequired(let cookies):
                    signInForm.requireVerification(for: input, cookies: cookies)
                    keepPrepared = true
                    applicationStatus = AuthenticationStage.twoFactor.rawValue
                case .authenticated(let account):
                    applyStoreAccount(account, restored: false)
                }
            } catch let error where error is CancellationError || Task.isCancelled {
                applicationStatus = "Sign-in cancelled."
            } catch {
                if let preparation {
                    if hasSent2FACode {
                        signInForm.updateChallengeCookies(await preparation.loginTransport.cookies())
                    }
                }
                keepPrepared = signInForm.failed(error).preservesPreparedSession
                // Apple/SAP errors have sanitized, bounded descriptions. Arbitrary
                // URL errors can include routing secrets: expose only numeric codes.
                authenticationError = authenticationMessage(error)
                applicationStatus = "Sign-in failed."
                #if DEBUG
                print("Apple sign-in failed (code \((error as NSError).code)).")
                #endif
            }
            if let preparation {
                if keepPrepared, !Task.isCancelled, preparation.canReuse(for: input.email), preparedAppleLogin === preparation {
                    expireLoginPreparation(preparation)
                } else {
                    if preparedAppleLogin === preparation { preparedAppleLogin = nil }
                    await preparation.close()
                }
            }
        }
    }

    func requestNewVerificationCode() {
        guard signInForm.canRequestNewCode, !isAuthenticating, !isCodeResendCoolingDown else { return }
        // Password remains in RAM for this sign-in only. A fresh signed password-only
        // request asks Apple to start another challenge; Apple controls code delivery.
        startAppleLogin(.requestNewCode)
        guard isAuthenticating else { return }
        isCodeResendCoolingDown = true
        codeResendCooldownTask?.cancel()
        codeResendCooldownTask = Task {
            do { try await Task.sleep(nanoseconds: 30_000_000_000) }
            catch { return }
            isCodeResendCoolingDown = false
            codeResendCooldownTask = nil
        }
    }

    func changeAppleAccount() {
        guard !isAuthenticating else { return }
        cancelAppleLogin()
        appleId = ""
    }

    private func clearCodeResendCooldown() {
        codeResendCooldownTask?.cancel(); codeResendCooldownTask = nil
        isCodeResendCoolingDown = false
    }

    func cancelAppleLogin() {
        authenticationTask?.cancel()
        clearCodeResendCooldown()
        clearLoginPreparation()
        signInForm.cancel()
        authenticationError = ""
        applicationStatus = "Not logged in!".localized
    }

    func restoreStoreAccount() {
        guard !didRestoreStoreAccount else { return }
        didRestoreStoreAccount = true
        do {
            // Legacy sessions include a password and may use a key file. Discard
            // them without decrypting/importing and require one fresh SAP login.
            try LegacyCredentials.remove()
            guard let account = try KeychainStoreAccount().load() else { return }
            try account.validate(identity: KeychainMachineIdentity.loadOrCreate())
            applyStoreAccount(account, restored: true)
        } catch {
            authenticationError = "Saved session could not be restored (code \((error as NSError).code)). Sign in again."
            // Do not overwrite or delete an inaccessible Keychain item during device lock.
        }
    }

    func logoutStoreAccount() {
        guard !isAuthenticating, storeTask == nil, !showStoreVersions else { return }
        do {
            clearCodeResendCooldown()
            clearLoginPreparation()
            try KeychainKBSync().clear()
            try KeychainStoreAccount().clear()
            try LegacyCredentials.remove()
            ipaTool?.close()
            ipaTool = nil
            isAuthenticated = false
            signInForm.cancel(clearEmail: true)
            authenticationError = ""
            hasAppBeenServed = false
            applicationStatus = "Not logged in!".localized
            applicationIcon = "xmark.circle.fill"
            // The machine-identity Keychain item deliberately survives logout.
        } catch {
            authenticationError = "Could not remove the saved session (code \((error as NSError).code))."
            Alertinator.shared.alert(title: "Log Out".localized, body: authenticationError)
        }
    }

    private func expireLoginPreparation(_ prepared: PreparedAppleLogin) {
        loginPreparationExpiry?.cancel()
        loginPreparationExpiry = Task {
            let remaining = max(0, prepared.expiresAt.timeIntervalSinceNow)
            try? await Task.sleep(nanoseconds: UInt64(remaining * 1_000_000_000))
            guard !Task.isCancelled, preparedAppleLogin === prepared, !isAuthenticating else { return }
            preparedAppleLogin = nil
            await prepared.close()
        }
    }
    private func clearLoginPreparation() {
        loginPreparationExpiry?.cancel(); loginPreparationExpiry = nil
        let prepared = preparedAppleLogin
        preparedAppleLogin = nil
        Task { await prepared?.close() }
    }

    private func setAuthenticationStage(_ stage: AuthenticationStage) {
        switch stage {
        case .bag, .sap, .signing, .prepared: applicationStatus = "Preparing sign-in…"
        case .authenticating, .redirect: applicationStatus = "Connecting to Apple…"
        case .retrying: applicationStatus = "Apple is taking longer. Retrying…"
        case .saving: applicationStatus = "Completing sign-in…"
        case .twoFactor: applicationStatus = "Enter your verification code."
        }
        #if DEBUG
        print("Apple authentication stage: \(stage.rawValue)")
        #endif
    }
    private func applyStoreAccount(_ account: StoreAccount, restored: Bool) {
        clearCodeResendCooldown()
        signInForm.authenticated(email: account.email)
        ipaTool = IPATool(account: account)
        isAuthenticated = true
        applicationStatus = "Signed in. Choose an app to get started."
        applicationIcon = "checkmark.circle.fill"
        applicationIconColor = .primary
        #if DEBUG
        print("Apple authentication: \(restored ? "saved session loaded" : "DSID/token/storefront received and saved in Keychain") [values withheld]")
        #endif
    }

    private func authenticationMessage(_ error: Error) -> String {
        if let error = error as? AuthenticationError { return error.localizedDescription }
        if let error = error as? SignInInputError { return error.localizedDescription }
        if let error = error as? SAPError { return error.localizedDescription }
        return "Sign-in failed (code \((error as NSError).code))."
    }

}

// A short-lived, same-account preparation avoids repeating the Bag/certificate/
// SAP handshake for a manual retry or 2FA. It holds no password/code, writes
// nothing to disk, signs every request freshly and expires after five minutes.
@MainActor
final class PreparedAppleLogin {
    let email: String
    let identity: MachineIdentity
    let loginTransport: AppleAuthenticationTransport
    private let sapTransport = AppleSAPTransport()
    private var signer: SAPSession?
    private var closed = false
    var endpoint: URL?
    private(set) var expiresAt = Date.distantPast
    init(email: String, cookies: [StoreCookie]) throws {
        self.email = email
        identity = try KeychainMachineIdentity.loadOrCreate()
        loginTransport = AppleAuthenticationTransport(cookies: cookies, isolatedConnections: true)
    }
    func canReuse(for email: String) -> Bool {
        !closed && self.email == email && signer != nil && expiresAt > Date()
    }
    func prepare(progress: (AuthenticationStage) -> Void) async throws -> SAPSession {
        guard !closed else { throw CancellationError() }
        if let signer { progress(.prepared); return signer }
        progress(.bag)
        let configuration = try await SAPProtocol(transport: sapTransport).bag(identity: identity)
        endpoint = try AuthenticationEndpoint.validate(configuration.authenticationURL)
        progress(.sap)
        let created = try SAPSession(guest: NativeSAPGuest(), transport: sapTransport)
        // Own the guest immediately, including cancellation during initialization.
        signer = created
        try await created.initialize(configuration: configuration, identity: identity)
        try Task.checkCancellation()
        guard !closed else { throw CancellationError() }
        expiresAt = Date().addingTimeInterval(300)
        return created
    }
    func close() async {
        guard !closed else { return }
        closed = true
        sapTransport.close(); loginTransport.close()
        if let signer { await signer.close() }
        signer = nil; endpoint = nil
    }
}

private enum LegacyCredentials {
    static func remove() throws {
        let fm = FileManager.default
        let documents = try fm.url(for: .documentDirectory, in: .userDomainMask, appropriateFor: nil, create: false)
        let library = try fm.url(for: .libraryDirectory, in: .userDomainMask, appropriateFor: nil, create: false)
        for file in [documents.appendingPathComponent("authinfo"), library.appendingPathComponent(".authkey")] {
            if fm.fileExists(atPath: file.path) { try fm.removeItem(at: file) }
        }
        let status = SecItemDelete([kSecClass as String: kSecClassKey,
            kSecAttrApplicationTag as String: "com.certlium.crumb.key"] as CFDictionary)
        guard status == errSecSuccess || status == errSecItemNotFound else { throw SAPError.keychain(status) }
    }
}
