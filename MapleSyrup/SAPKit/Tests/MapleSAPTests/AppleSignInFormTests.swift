import XCTest
@testable import MapleSAP
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

final class AppleSignInFormTests: XCTestCase {
    private let email = "fixture@example.test"
    private let password = " fixture secret "
    private func cookie(_ name: String) -> StoreCookie {
        StoreCookie(HTTPCookie(properties: [.name: name, .value: "fixture-cookie",
            .domain: ".itunes.apple.com", .path: "/", .secure: "TRUE"])!)
    }
    private func challenged() throws -> AppleSignInForm {
        var form = AppleSignInForm(email: email, password: password)
        let input = try form.begin()
        form.requireVerification(for: input, cookies: [cookie("challenge")])
        return form
    }

    func testInitialInputTrimsAccountButNeverPasswordOrAppendsUnrequestedCode() throws {
        var form = AppleSignInForm(email: " \n" + email + " ", password: password)
        form.code = "999999"
        let input = try form.begin()
        XCTAssertEqual(input.email, email)
        XCTAssertEqual(input.password, password)
        XCTAssertEqual(input.code, "")
        XCTAssertTrue(input.cookies.isEmpty)
        XCTAssertThrowsError(try form.submission(.requestNewCode)) {
            XCTAssertEqual($0 as? SignInInputError, .noChallenge)
        }
    }

    func testMissingCredentialsCannotStartPreparation() throws {
        for form in [AppleSignInForm(email: " \n", password: password), AppleSignInForm(email: email)] {
            XCTAssertFalse(form.canSubmit)
            XCTAssertThrowsError(try form.submission()) {
                XCTAssertEqual($0 as? SignInInputError, .missingCredentials)
            }
        }
    }

    func testChallengeKeepsAccountBoundCredentialsIndependentOfEditablePassword() throws {
        var form = try challenged()
        form.password = "" // SecureField/autofill must not destroy the challenge.
        form.code = "123 456"
        let input = try form.begin()
        XCTAssertEqual(input.password, password)
        XCTAssertEqual(input.code, "123456")
        XCTAssertEqual(input.cookies.map(\.name), ["challenge"])
        XCTAssertTrue(form.awaitsVerification)
    }

    func testBlankOrMalformedVerificationCannotFallBackToPasswordOnlyLogin() throws {
        var form = try challenged()
        XCTAssertFalse(form.canSubmit)
        XCTAssertThrowsError(try form.begin()) {
            XCTAssertEqual($0 as? SignInInputError, .verificationCodeRequired)
        }
        for code in ["12345", "12a456", "１２３４５６", "   "] {
            form.code = code
            XCTAssertFalse(form.canSubmit)
            XCTAssertThrowsError(try form.begin())
        }
        form.code = "12\n34 56"
        XCTAssertTrue(form.canSubmit)
    }

    func testChallengeCannotBeSubmittedOrResentToAnotherAccount() throws {
        var form = try challenged()
        form.code = "123456"
        form.email = "other@example.test"
        XCTAssertFalse(form.canSubmit)
        XCTAssertFalse(form.canRequestNewCode)
        for action in [AppleSignInAction.submit, .requestNewCode] {
            XCTAssertThrowsError(try form.begin(action)) {
                XCTAssertEqual($0 as? SignInInputError, .accountChanged)
            }
        }
    }

    func testTemporaryFailuresKeepCodeChallengeAndUpdatedCookies() throws {
        let errors: [AuthenticationError] = [.http(404), .http(503), .network(-1001), .rateLimited,
            .retryLater, .invalidResponse(204), .redirectUnavailable(301)]
        for error in errors {
            var form = try challenged()
            form.code = "123456"
            form.updateChallengeCookies([cookie("pod")])
            XCTAssertEqual(form.failed(error), .retrySameSubmission)
            XCTAssertTrue(form.lastRecovery!.preservesPreparedSession)
            XCTAssertTrue(form.isRetryingVerification)
            let retry = try form.begin()
            XCTAssertEqual(retry.code, "123456")
            XCTAssertEqual(retry.password, password)
            XCTAssertEqual(retry.cookies.map(\.name), ["pod"])
            XCTAssertNil(form.lastRecovery)
        }
    }

    func testExplicitAppleVerificationRejectionClearsOnlyCode() throws {
        var form = try challenged()
        form.code = "123456"
        XCTAssertEqual(form.failed(AuthenticationError.verificationRejected), .replaceCode)
        XCTAssertEqual(form.code, "")
        XCTAssertTrue(form.awaitsVerification)
        XCTAssertTrue(form.canRequestNewCode)
        XCTAssertFalse(form.canSubmit)
        form.code = "654321"
        XCTAssertEqual(try form.begin().password, password)
    }

    func testLocalFormattingErrorKeepsEnteredCodeForCorrection() throws {
        var form = try challenged()
        form.code = "12a456"
        XCTAssertEqual(form.failed(AuthenticationError.invalidCode), .editCode)
        XCTAssertEqual(form.code, "12a456")
        XCTAssertTrue(form.awaitsVerification)
        XCTAssertTrue(form.lastRecovery!.preservesPreparedSession)
        XCTAssertFalse(form.canSubmit)
    }

    func testExplicitNewCodeRequestAloneDropsOldCodeAndCookies() throws {
        var form = try challenged()
        form.code = "123456"
        let input = try form.begin(.requestNewCode)
        XCTAssertEqual(input.password, password)
        XCTAssertEqual(input.code, "")
        XCTAssertTrue(input.cookies.isEmpty)
        XCTAssertEqual(form.code, "")
        XCTAssertTrue(form.awaitsVerification)
        XCTAssertTrue(form.canRequestNewCode)
        XCTAssertFalse(form.canSubmit)
    }

    func testRejectedCredentialsUnlockFieldsAndEraseChallengeSecrets() throws {
        for error in [AuthenticationError.invalidCredentials, .accountDisabled] {
            var form = try challenged()
            form.code = "123456"
            XCTAssertEqual(form.failed(error), .credentialsRejected)
            XCTAssertFalse(form.awaitsVerification)
            XCTAssertFalse(form.canSubmit)
            XCTAssertFalse(form.canRequestNewCode)
            XCTAssertEqual(form.email, email)
            XCTAssertEqual(form.password, "")
            XCTAssertEqual(form.code, "")
        }
    }

    func testUnsafeRedirectOrUnknownErrorNeverCreatesAnAutomaticChallenge() throws {
        for error: Error in [AuthenticationError.invalidRedirect, SAPError.invalidExchange] {
            var form = try challenged()
            form.code = "123456"
            let disposition = form.failed(error)
            XCTAssertEqual(disposition, .stop)
            XCTAssertFalse(disposition.preservesPreparedSession)
            XCTAssertFalse(form.isRetryingVerification)
            XCTAssertEqual(form.code, "123456") // No evidence Apple rejected it.
        }
    }

    func testCancelSuccessAndAccountChangeEraseEphemeralCredentials() throws {
        var form = try challenged()
        form.code = "123456"
        form.cancel()
        XCTAssertEqual(form.email, email)
        XCTAssertEqual(form.password, "")
        XCTAssertEqual(form.code, "")
        XCTAssertFalse(form.awaitsVerification)
        XCTAssertNil(form.lastRecovery)
        form = try challenged()
        form.authenticated(email: "canonical@example.test")
        XCTAssertEqual(form.email, "canonical@example.test")
        XCTAssertEqual(form.password, "")
        XCTAssertEqual(form.code, "")
        XCTAssertFalse(form.canSubmit)
        form = try challenged()
        form.cancel(clearEmail: true)
        XCTAssertEqual(form.email, "")
        XCTAssertFalse(form.canRequestNewCode)
    }
}
