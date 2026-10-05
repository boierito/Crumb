import Foundation
import Security
import SwiftUI
import PartyUI
import MapleSAP

// MainActor owns UI state; SAPSession runs the blocking native guest away from it.
extension AppData {
    func startAppleLogin() {
        guard !isAuthenticating, !appleId.isEmpty, !password.isEmpty else { return }
        let email = appleId.trimmingCharacters(in: .whitespacesAndNewlines)
        let secret = password
        let verification = hasSent2FACode ? code : ""
        let challengeCookies = hasSent2FACode ? pendingAuthenticationCookies : []
        isAuthenticating = true
        authenticationError = ""
        authenticationDiagnostic = ["WaffleStore authentication probe v4",
            "iOS=\(UIDevice.current.systemVersion)",
            "app-build=\(Bundle.main.object(forInfoDictionaryKey: "CFBundleVersion") ?? "unknown")",
            "password-persistence=false", "signer=tci-no-jit"].joined(separator: "\n")
        authenticationTask = Task {
            defer { isAuthenticating = false; authenticationTask = nil }
            let sapTransport = AppleSAPTransport()
            let loginTransport = AppleAuthenticationTransport(cookies: challengeCookies, isolatedConnections: true)
            defer { sapTransport.close(); loginTransport.close() }
            var sap: SAPSession?
            do {
                // Invalid codes are rejected before downloading assets or sending credentials.
                _ = try TwoFactorAuthentication.normalize(verification)
                let identity = try KeychainMachineIdentity.loadOrCreate()
                setAuthenticationStage(.bag)
                let configuration = try await SAPProtocol(transport: sapTransport).bag(identity: identity)
                _ = try AuthenticationEndpoint.validate(configuration.authenticationURL)
                setAuthenticationStage(.sap)
                let signer = try SAPSession(guest: NativeSAPGuest(), transport: sapTransport)
                sap = signer
                try await signer.initialize(configuration: configuration, identity: identity)
                try Task.checkCancellation()
                let authentication = AppleAuthentication(transport: loginTransport, signer: signer,
                    persistence: KeychainStoreAccount(), diagnostic: { event in
                        await MainActor.run {
                            self.authenticationDiagnostic += "\n\(event)"
                            print("Apple authentication diagnostic: \(event)")
                        }
                    })
                let outcome = try await authentication.login(email: email, password: secret, code: verification,
                    identity: identity, endpoint: configuration.authenticationURL) { stage in
                        await MainActor.run { self.setAuthenticationStage(stage) }
                    }
                try Task.checkCancellation()
                switch outcome {
                case .twoFactorRequired(let cookies):
                    pendingAuthenticationCookies = cookies
                    authenticationDiagnostic += "\noutcome=2FA-required"
                    hasSent2FACode = true
                    code = ""
                    applicationStatus = AuthenticationStage.twoFactor.rawValue
                case .authenticated(let account):
                    authenticationDiagnostic += "\ntwo-factor=\(verification.isEmpty ? "not-requested-in-this-login" : "submitted")"
                    applyStoreAccount(account, restored: false)
                }
            } catch let error where error is CancellationError || Task.isCancelled {
                applicationStatus = "Sign-in cancelled."
                authenticationDiagnostic += "\noutcome=cancelled"
            } catch {
                if hasSent2FACode { pendingAuthenticationCookies = await loginTransport.cookies() }
                // Apple/SAP errors have sanitized, bounded descriptions. Arbitrary
                // URL errors can include routing secrets: expose only numeric codes.
                if let error = error as? AuthenticationError { authenticationError = error.localizedDescription }
                else if let error = error as? SAPError { authenticationError = error.localizedDescription }
                else { authenticationError = "Sign-in failed (code \((error as NSError).code))." }
                code = ""
                applicationStatus = "Sign-in failed."
                print("Apple authentication failed: \(diagnosticCategory(error))")
                // Customer messages belong in the UI; exported logs contain fixed
                // stage/category information, not arbitrary Apple response text.
                authenticationDiagnostic += "\noutcome=failed; error-category=\(diagnosticCategory(error))"
            }
            if let sap = sap { await sap.close() }
        }
    }

    func cancelAppleLogin() {
        authenticationTask?.cancel()
        hasSent2FACode = false
        code = ""
        password = ""
        pendingAuthenticationCookies = []
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
            try KeychainKBSync().clear()
            try KeychainStoreAccount().clear()
            try LegacyCredentials.remove()
            ipaTool?.close()
            ipaTool = nil
            isAuthenticated = false
            hasSent2FACode = false
            appleId = ""; password = ""; code = ""
            pendingAuthenticationCookies = []
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

    private func setAuthenticationStage(_ stage: AuthenticationStage) {
        applicationStatus = stage.rawValue
        authenticationDiagnostic += "\nstage=\(stage.rawValue)"
        print("Apple authentication stage: \(stage.rawValue)")
    }
    private func applyStoreAccount(_ account: StoreAccount, restored: Bool) {
        appleId = account.email
        password = ""; code = ""; hasSent2FACode = false
        pendingAuthenticationCookies = []
        ipaTool = IPATool(account: account)
        isAuthenticated = true
        applicationStatus = restored ? "Saved session loaded; Apple validity not checked." : "Signed in. Choose an app/version to download."
        applicationIcon = "checkmark.circle.fill"
        applicationIconColor = .primary
        print("Apple authentication: \(restored ? "saved session loaded" : "DSID/token/storefront received and saved in Keychain") [values withheld]")
        authenticationDiagnostic += "\noutcome=\(restored ? "saved-session-loaded-not-validated" : "authenticated")\nDSID=present\npasswordToken=present\nstorefront=present\npod=\(account.pod == nil ? "absent" : "present")\nsecret-values=withheld"
    }
    private func diagnosticCategory(_ error: Error) -> String {
        if let error = error as? AuthenticationError {
            switch error {
            case .apple: return "Apple-account-response"
            case .http(let status), .invalidResponse(let status): return "HTTP-\(status)"
            case .network(let code): return "network-\(code)"
            case .invalidCode: return "invalid-2FA-format"
            case .invalidCredentials: return "invalid-credentials"
            case .twoFactorRequired: return "2FA-required"
            case .verificationRejected: return "2FA-rejected-or-expired"
            case .accountDisabled: return "account-disabled"
            case .rateLimited: return "rate-limit"
            case .retryLater: return "retry-after-over-budget"
            case .invalidRedirect: return "endpoint-or-redirect-rejected"
            case .tooManyRedirects: return "redirect-limit"
            case .invalidSession: return "invalid-session"
            }
        }
        if let error = error as? SAPError {
            switch error {
            case .nativeRuntime(let stage): return "SAP-native-stage-\(stage)"
            case .keychain(let status): return "keychain-\(status)"
            case .http(let status): return "SAP-HTTP-\(status)"
            default: return "SAP-protocol"
            }
        }
        return "operation-\((error as NSError).code)"
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
            kSecAttrApplicationTag as String: "com.nxtcoreee3.WaffleStore.key"] as CFDictionary)
        guard status == errSecSuccess || status == errSecItemNotFound else { throw SAPError.keychain(status) }
    }
}
