import Foundation

// RAM-only report for this branch. Never retain a request/response body, URL,
// header value, account identifier, machine ID, code or error description.
public struct AuthenticationTestReport {
    private var lines: [String] = []
    private var startedAt: Date?
    private var trial = 0
    public var text: String {
        (["Crumb clean-login test v1", "secret-values=withheld",
          "keychain-scope=test-only", "Apple-controls-2FA=true"] + lines).joined(separator: "\n")
    }
    public init() {}

    public mutating func beginReset() {
        lines = []; trial = 0; startedAt = nil
        append("reset=started; apple-network=not-requested-by-reset")
    }

    public mutating func reset(_ evidence: AuthenticationTestResetEvidence) {
        lines = []; trial = 0; startedAt = nil
        append("reset=verified; identity-changed=\(evidence.identityChanged); identity-persisted=\(evidence.identityPersisted)")
        append("session-removed=\(evidence.sessionRemoved); kbsync-removed=\(evidence.kbsyncRemoved)")
        append("cookies=discarded; prepared-SAP=closed; SAP-assets=retained")
        append("apple-network=not-requested-by-reset")
    }
    public mutating func begin(_ action: AppleSignInAction, verification: Bool, now: Date = Date()) {
        trial += 1; startedAt = now
        let kind = action == .requestNewCode ? "request-new-code" : (verification ? "verify" : "password-login")
        append("trial=\(trial); action=\(kind)")
    }
    public mutating func stage(_ stage: AuthenticationStage, now: Date = Date()) {
        append("elapsed-ms=\(elapsed(now)); stage=\(stage.rawValue)")
    }
    public mutating func challenge(now: Date = Date()) {
        append("elapsed-ms=\(elapsed(now)); outcome=Apple-requested-2FA")
    }
    public mutating func authenticated(verification: Bool, now: Date = Date()) {
        let kind = verification ? "verified-with-submitted-code" : "Apple-did-not-request-2FA-in-this-submission"
        append("elapsed-ms=\(elapsed(now)); outcome=authenticated; \(kind)")
    }
    public mutating func restored() { append("session=restored; not-a-fresh-login-test") }
    public mutating func cancelled() { append("outcome=cancelled") }
    public mutating func failed(_ error: Error) {
        let category: String
        if let error = error as? AuthenticationError {
            switch error {
            case .http(let status), .invalidResponse(let status): category = "HTTP-\(status)"
            case .network(let code): category = "network-\(code)"
            case .verificationRejected: category = "Apple-verification-rejected"
            case .invalidCode: category = "local-code-format"
            case .invalidCredentials: category = "Apple-credentials-rejected"
            case .accountDisabled: category = "Apple-account-disabled"
            case .rateLimited: category = "rate-limited"
            case .retryLater: category = "retry-budget-or-wait"
            case .redirectUnavailable(let status): category = "incomplete-redirect-\(status)"
            case .invalidRedirect, .tooManyRedirects: category = "redirect-rejected"
            case .apple: category = "Apple-error-message-withheld"
            case .twoFactorRequired: category = "2FA-required"
            case .invalidSession: category = "invalid-session"
            }
        } else if let error = error as? SAPError {
            switch error {
            case .keychain(let status): category = "keychain-\(status)"
            case .nativeRuntime(let stage): category = "SAP-native-\(stage)"
            case .http(let status): category = "SAP-HTTP-\(status)"
            default: category = "SAP"
            }
        } else if error is AuthenticationTestResetError { category = "reset-verification-failed" }
        else { category = "operation-\((error as NSError).code)" }
        append("outcome=failed; category=\(category)")
    }

    // Accept complete fixed-field messages only, never sanitise by printing a
    // substring of an arbitrary message. Unknown/new diagnostic formats are dropped.
    public mutating func diagnostic(_ message: String) {
        guard message.utf8.count <= 220,
              Self.diagnosticPatterns.contains(where: {
                  message.range(of: $0, options: .regularExpression) != nil
              }) else { return }
        append(message)
    }
    private static let diagnosticPatterns = [
        #"\Aauthentication-recovery-attempt=[0-9]{1,2}/[0-9]{1,2}\z"#,
        #"\Acookie-jar-count=[0-9]{1,6}; request-cookie-count=[0-9]{1,6}\z"#,
        #"\Aauthentication-redirect=received; HTTP=[0-9]{3}; location-present=(true|false)\z"#,
        #"\Ascope=authentication; attempt=[0-9]{1,2}; HTTP=[0-9]{3}; body=(plist|html|empty|other); apple-failure=(absent|present-withheld|-?[0-9]{1,8})\z"#
    ]
    private func elapsed(_ now: Date) -> Int {
        guard let startedAt else { return 0 }
        return Int(max(0, min(86_400, now.timeIntervalSince(startedAt))) * 1000)
    }
    private mutating func append(_ line: String) {
        lines.append(line)
        if lines.count > 160 { lines.removeFirst(lines.count - 160) }
    }
}
