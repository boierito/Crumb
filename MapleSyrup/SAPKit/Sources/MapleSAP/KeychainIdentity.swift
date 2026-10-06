#if canImport(Security)
import Foundation
import Security

public enum KeychainMachineIdentity {
    private static let lock = NSLock()
    private static let service = AuthenticationTestScope.keychainService
    private static let account = "machine-identity-v1"

    // No Apple ID, network interface MAC, private entitlement, file fallback,
    // or UserDefaults. Keeping this item on logout preserves the device identity.
    public static func loadOrCreate() throws -> MachineIdentity {
        lock.lock()
        defer { lock.unlock() }
        return try loadOrCreateUnlocked(service: service)
    }

    // Branch-only, local reset. UI must close live SAP/Store work first. It
    // affects only this test namespace and never asks Apple to remove trust.
    public static func resetForAuthenticationTest() throws -> AuthenticationTestResetEvidence {
        try resetForAuthenticationTest(service: service)
    }

    static func resetForAuthenticationTest(service: String) throws -> AuthenticationTestResetEvidence {
        guard service == Self.service || service.hasPrefix(Self.service + ".tests.") else {
            throw AuthenticationTestResetError.namespaceRejected
        }
        lock.lock()
        defer { lock.unlock() }
        let previous = try loadOrCreateUnlocked(service: service)
        try KeychainStoreAccount(service: service).clear()
        try KeychainKBSync(service: service).clear()
        let status = SecItemDelete([kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service, kSecAttrAccount as String: account] as CFDictionary)
        guard status == errSecSuccess || status == errSecItemNotFound else { throw SAPError.keychain(status) }
        let created = try loadOrCreateUnlocked(service: service)
        let restored = try loadOrCreateUnlocked(service: service)
        let evidence = AuthenticationTestResetEvidence(identityChanged: created != previous,
            identityPersisted: restored == created,
            sessionRemoved: try !containsItem(service: service, account: "store-account-v2"),
            kbsyncRemoved: try !containsItem(service: service, account: "kbsync-v1"))
        guard evidence.identityChanged, evidence.identityPersisted,
              evidence.sessionRemoved, evidence.kbsyncRemoved else { throw AuthenticationTestResetError.evidenceFailed }
        return evidence
    }

    private static func containsItem(service: String, account: String) throws -> Bool {
        let status = SecItemCopyMatching([kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service, kSecAttrAccount as String: account,
            kSecMatchLimit as String: kSecMatchLimitOne] as CFDictionary, nil)
        guard status == errSecSuccess || status == errSecItemNotFound else { throw SAPError.keychain(status) }
        return status == errSecSuccess
    }

    private static func loadOrCreateUnlocked(service: String) throws -> MachineIdentity {
        let query: [String: Any] = [kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service, kSecAttrAccount as String: account]
        var read = query
        read[kSecReturnData as String] = true
        read[kSecMatchLimit as String] = kSecMatchLimitOne
        var result: CFTypeRef?
        let status = SecItemCopyMatching(read as CFDictionary, &result)
        if status == errSecSuccess {
            guard let data = result as? Data else { throw SAPError.invalidIdentity }
            return try MachineIdentity(hardwareID: data)
        }
        guard status == errSecItemNotFound else { throw SAPError.keychain(status) }
        var bytes = [UInt8](repeating: 0, count: 6)
        let randomStatus = SecRandomCopyBytes(kSecRandomDefault, bytes.count, &bytes)
        guard randomStatus == errSecSuccess else { throw SAPError.keychain(randomStatus) }
        // Locally administered unicast address, stable for this installation.
        // Apple's acceptance of this identity still requires device validation.
        bytes[0] = (bytes[0] & 0xFC) | 0x02
        let identity = try MachineIdentity(hardwareID: Data(bytes))
        var insert = query
        insert[kSecValueData as String] = identity.hardwareID
        insert[kSecAttrAccessible as String] = kSecAttrAccessibleAfterFirstUnlockThisDeviceOnly
        let saved = SecItemAdd(insert as CFDictionary, nil)
        guard saved == errSecSuccess else { throw SAPError.keychain(saved) }
        return identity
    }
}
#endif
