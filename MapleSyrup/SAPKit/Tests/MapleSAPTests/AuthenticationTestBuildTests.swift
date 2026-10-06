import XCTest
@testable import MapleSAP
#if canImport(Security)
import Security
#endif

final class AuthenticationTestBuildTests: XCTestCase {
    func testTestScopeIsDifferentFromProductionEvenWithSharedAccessGroup() {
        XCTAssertNotEqual(AuthenticationTestScope.keychainService, "com.certlium.crumb.sap")
        XCTAssertNotEqual(AuthenticationTestScope.bundleIdentifier, "com.certlium.crumb")
    }

    func testReportIncludesResetEvidenceWithoutIdentityValues() {
        var report = AuthenticationTestReport()
        report.reset(AuthenticationTestResetEvidence(identityChanged: true, identityPersisted: true,
            sessionRemoved: true, kbsyncRemoved: true))
        XCTAssertTrue(report.text.contains("identity-changed=true; identity-persisted=true"))
        XCTAssertTrue(report.text.contains("session-removed=true; kbsync-removed=true"))
        XCTAssertTrue(report.text.contains("cookies=discarded; prepared-SAP=closed; SAP-assets=retained"))
        XCTAssertTrue(report.text.contains("apple-network=not-requested-by-reset"))
        XCTAssertFalse(report.text.contains("GUID="))
    }

    func testResetFailureCannotRetainEarlierSuccessReport() {
        var report = AuthenticationTestReport()
        report.reset(AuthenticationTestResetEvidence(identityChanged: true, identityPersisted: true,
            sessionRemoved: true, kbsyncRemoved: true))
        report.beginReset()
        report.failed(SAPError.keychain(-25308))
        XCTAssertFalse(report.text.contains("reset=verified"))
        XCTAssertTrue(report.text.contains("reset=started"))
        XCTAssertTrue(report.text.contains("category=keychain--25308"))
    }

    func testReportDistinguishesActualChallengeVerificationAndNoChallenge() {
        let clock = Date(timeIntervalSince1970: 10)
        var report = AuthenticationTestReport()
        report.begin(.submit, verification: false, now: clock)
        report.stage(.sap, now: clock.addingTimeInterval(1))
        report.challenge(now: clock.addingTimeInterval(2))
        report.begin(.submit, verification: true, now: clock.addingTimeInterval(20))
        report.authenticated(verification: true, now: clock.addingTimeInterval(21))
        XCTAssertTrue(report.text.contains("elapsed-ms=1000; stage=Initializing SAP"))
        XCTAssertTrue(report.text.contains("elapsed-ms=2000; outcome=Apple-requested-2FA"))
        XCTAssertTrue(report.text.contains("trial=2; action=verify"))
        XCTAssertTrue(report.text.contains("verified-with-submitted-code"))
        var noChallenge = AuthenticationTestReport()
        noChallenge.begin(.submit, verification: false, now: clock)
        noChallenge.authenticated(verification: false, now: clock)
        XCTAssertTrue(noChallenge.text.contains("Apple-did-not-request-2FA-in-this-submission"))
        XCTAssertFalse(noChallenge.text.contains("outcome=Apple-requested-2FA"))
    }

    func testUnknownHeadersURLsCredentialsAndMessagesNeverEnterReport() {
        var report = AuthenticationTestReport()
        let secret = "fixture-sensitive-value"
        let valid = "scope=authentication; attempt=1; HTTP=404; body=html; apple-failure=absent"
        report.diagnostic(valid)
        report.diagnostic("cookie-jar-count=1; request-cookie-count=1")
        for text in ["password=" + secret, "Cookie: " + secret, "code=123456", "GUID=020000000000",
                     "https://buy.itunes.apple.com/?token=" + secret, valid + "; token=" + secret,
                     valid + "\npassword=" + secret, String(repeating: "x", count: 1000)] {
            report.diagnostic(text)
        }
        report.failed(AuthenticationError.apple(failure: secret, message: secret))
        report.failed(NSError(domain: secret, code: -1, userInfo: [NSLocalizedDescriptionKey: secret]))
        XCTAssertTrue(report.text.contains(valid))
        XCTAssertTrue(report.text.contains("cookie-jar-count=1; request-cookie-count=1"))
        for forbidden in [secret, "123456", "020000000000", "https://", "token="] {
            XCTAssertFalse(report.text.contains(forbidden))
        }
    }

    func testReportAcceptsOnlyFixedAuthenticationDiagnosticsAndIsBounded() {
        var report = AuthenticationTestReport()
        report.diagnostic("authentication-recovery-attempt=1/12")
        report.diagnostic("authentication-redirect=received; HTTP=302; location-present=true")
        report.diagnostic("scope=authentication; attempt=2; HTTP=200; body=plist; apple-failure=-5000")
        XCTAssertTrue(report.text.contains("authentication-recovery-attempt=1/12"))
        XCTAssertTrue(report.text.contains("location-present=true"))
        XCTAssertTrue(report.text.contains("apple-failure=-5000"))
        for _ in 0..<500 { report.stage(.retrying) }
        XCTAssertLessThanOrEqual(report.text.split(separator: "\n").count, 164)
    }

    func testRestoredSessionIsNotLabelledAsFreshLogin() {
        var report = AuthenticationTestReport()
        report.restored()
        XCTAssertTrue(report.text.contains("not-a-fresh-login-test"))
        XCTAssertFalse(report.text.contains("outcome=authenticated"))
    }

    #if canImport(Security)
    func testRealKeychainResetRotatesIdentityAndRemovesOnlySpecifiedTestItems() throws {
        let service = AuthenticationTestScope.keychainService + ".tests." + UUID().uuidString
        let base: [String: Any] = [kSecClass as String: kSecClassGenericPassword, kSecAttrService as String: service]
        defer { SecItemDelete(base as CFDictionary) }
        let previous = try MachineIdentity(hardwareID: Data([2, 1, 2, 3, 4, 5]))
        func insert(_ account: String, value: Data) throws {
            let query = base.merging([kSecAttrAccount as String: account,
                kSecValueData as String: value]) { _, new in new }
            XCTAssertEqual(SecItemAdd(query as CFDictionary, nil), errSecSuccess)
        }
        try insert("machine-identity-v1", value: previous.hardwareID)
        try insert("unrelated-item", value: Data("fixture-unrelated".utf8))
        let store = KeychainStoreAccount(service: service)
        let account = StoreAccount(email: "fixture@example.test", name: "Fixture", dsid: "123",
            passwordToken: "fixture-token", storefront: "143441-1,29", pod: "42", guid: previous.guid,
            authenticationURL: URL(string: "https://buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate/")!, cookies: [])
        try store.save(account)
        let sync = KeychainKBSync(service: service)
        try sync.save(Data("fixture-kbsync".utf8), dsid: account.dsid, guid: previous.guid)
        let evidence = try KeychainMachineIdentity.resetForAuthenticationTest(service: service)
        XCTAssertTrue(evidence.identityChanged)
        XCTAssertTrue(evidence.identityPersisted)
        XCTAssertTrue(evidence.sessionRemoved)
        XCTAssertTrue(evidence.kbsyncRemoved)
        XCTAssertNil(try store.load())
        XCTAssertNil(try sync.load(dsid: account.dsid, guid: previous.guid))
        let untouched = base.merging([kSecAttrAccount as String: "unrelated-item"]) { _, new in new }
        XCTAssertEqual(SecItemCopyMatching(untouched as CFDictionary, nil), errSecSuccess)
    }

    func testResetRejectsProductionAndUnrelatedNamespacesBeforeAccess() throws {
        for service in ["com.certlium.crumb.sap", "", "other.test.service",
                        AuthenticationTestScope.keychainService + "-lookalike"] {
            XCTAssertThrowsError(try KeychainMachineIdentity.resetForAuthenticationTest(service: service)) {
                guard case AuthenticationTestResetError.namespaceRejected = $0 else {
                    return XCTFail("Unexpected error type")
                }
            }
        }
    }
    #endif
}
