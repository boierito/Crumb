import Foundation

public enum AppleSignInAction: Equatable { case submit, requestNewCode }

public enum SignInInputError: Error, LocalizedError, Equatable {
    case missingCredentials, verificationCodeRequired, accountChanged, noChallenge
    public var errorDescription: String? {
        switch self {
        case .missingCredentials: return "Enter your Apple ID and password."
        case .verificationCodeRequired: return "Enter the six-digit verification code."
        case .accountChanged: return "Change Apple ID to start a new sign-in."
        case .noChallenge: return "Apple has not requested a verification code."
        }
    }
}

// Describes UI recovery, not another network retry or challenge request.
public enum SignInFailureRecovery: Equatable {
    case retrySameSubmission, editCode, replaceCode, credentialsRejected, stop

    public var preservesPreparedSession: Bool {
        switch self {
        case .retrySameSubmission, .editCode, .replaceCode: return true
        case .credentialsRejected, .stop: return false
        }
    }
    static func classify(_ error: Error) -> Self {
        if let error = error as? AuthenticationError {
            switch error {
            case .http, .network, .rateLimited, .retryLater, .invalidResponse, .redirectUnavailable:
                return .retrySameSubmission
            case .invalidCode: return .editCode
            case .verificationRejected: return .replaceCode
            case .invalidCredentials, .accountDisabled: return .credentialsRejected
            case .apple, .twoFactorRequired, .invalidRedirect, .tooManyRedirects, .invalidSession:
                return .stop
            }
        }
        if let error = error as? SignInInputError,
           error == .verificationCodeRequired { return .editCode }
        return .stop
    }
}

// Immutable RAM-only snapshot. Never Codable, persisted or included in logs.
public struct AppleSignInSubmission {
    public let email: String
    public let password: String
    public let code: String
    public let cookies: [StoreCookie]
}

// A single source of truth for the editable form and account-bound challenge.
// Password/code are erased on cancellation, account change and success. A
// transport error must not masquerade as rejection of an otherwise usable code.
public struct AppleSignInForm {
    public var email: String
    public var password: String
    public var code = ""
    public private(set) var lastRecovery: SignInFailureRecovery?
    private var challenge: Challenge?
    private struct Challenge {
        let email: String
        let password: String
        var cookies: [StoreCookie]
    }
    public init(email: String = "", password: String = "") {
        self.email = email; self.password = password
    }
    public var awaitsVerification: Bool { challenge != nil }
    public var canSubmit: Bool { (try? submission(.submit)) != nil }
    public var canRequestNewCode: Bool { (try? submission(.requestNewCode)) != nil }
    public var isRetryingVerification: Bool {
        awaitsVerification && lastRecovery == .retrySameSubmission && canSubmit
    }

    public func submission(_ action: AppleSignInAction = .submit) throws -> AppleSignInSubmission {
        let email = email.trimmingCharacters(in: .whitespacesAndNewlines)
        if let challenge {
            guard email == challenge.email else { throw SignInInputError.accountChanged }
            let code: String
            switch action {
            case .requestNewCode: code = ""
            case .submit:
                code = try TwoFactorAuthentication.normalize(self.code)
                guard !code.isEmpty else { throw SignInInputError.verificationCodeRequired }
            }
            return AppleSignInSubmission(email: challenge.email, password: challenge.password,
                code: code, cookies: action == .requestNewCode ? [] : challenge.cookies)
        }
        guard action == .submit else { throw SignInInputError.noChallenge }
        guard !email.isEmpty, !password.isEmpty else { throw SignInInputError.missingCredentials }
        return AppleSignInSubmission(email: email, password: password, code: "", cookies: [])
    }

    public mutating func begin(_ action: AppleSignInAction = .submit) throws -> AppleSignInSubmission {
        let input = try submission(action)
        lastRecovery = nil
        if action == .requestNewCode {
            code = ""
            challenge?.cookies = []
        }
        return input
    }

    public mutating func requireVerification(for input: AppleSignInSubmission, cookies: [StoreCookie]) {
        challenge = Challenge(email: input.email, password: input.password, cookies: cookies)
        code = ""
        lastRecovery = nil
    }

    public mutating func updateChallengeCookies(_ cookies: [StoreCookie]) {
        challenge?.cookies = cookies
    }

    @discardableResult
    public mutating func failed(_ error: Error) -> SignInFailureRecovery {
        let recovery = SignInFailureRecovery.classify(error)
        lastRecovery = recovery
        switch recovery {
        case .replaceCode: code = "" // Apple explicitly rejected verification.
        case .credentialsRejected:
            challenge = nil; password = ""; code = ""
        case .retrySameSubmission, .editCode, .stop: break
        }
        return recovery
    }

    public mutating func cancel(clearEmail: Bool = false) {
        password = ""; code = ""; challenge = nil; lastRecovery = nil
        if clearEmail { email = "" }
    }
    public mutating func authenticated(email: String) {
        cancel()
        self.email = email
    }
}
