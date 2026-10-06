// This source exists only on feature/2fa-clean-login-test. Never use production
// Keychain names for test reset, even when a resigning tool shares access groups.
public enum AuthenticationTestScope {
    public static let keychainService = "com.certlium.crumb.authtest.sap"
    public static let bundleIdentifier = "com.certlium.crumb.authtest"
}

public struct AuthenticationTestResetEvidence {
    public let identityChanged: Bool
    public let identityPersisted: Bool
    public let sessionRemoved: Bool
    public let kbsyncRemoved: Bool
}

public enum AuthenticationTestResetError: Error {
    case namespaceRejected, evidenceFailed
}
