import XCTest
@testable import MapleSAP

final class StoreTests: XCTestCase {
    func testStorefrontUsesAccountCountry() throws {
        XCTAssertEqual(try Storefront.country("143505-1,29"), "ar")
        XCTAssertEqual(try Storefront.country("143441-1,29"), "us")
        XCTAssertThrowsError(try Storefront.country("unknown"))
    }
    func testIdentifiersRejectBooleansAndFractions() {
        XCTAssertNil(StoreParsing.identifier(true))
        XCTAssertNil(StoreParsing.identifier(1.5))
        XCTAssertNil(StoreParsing.identifier("123&evil=1"))
        XCTAssertEqual(StoreParsing.identifier(NSNumber(value: UInt64.max)), String(UInt64.max))
    }
    func testResponseMustMatchSelectedExternalVersion() throws {
        let app = StoreApp(id: "123", bundleID: "test.app", name: "Test", price: 0)
        let root: [String: Any] = ["songList": [["URL": "https://iosapps.itunes.apple.com/test.ipa?token=withheld",
            "metadata": ["itemId": 123, "softwareVersionExternalIdentifier": "999", "softwareVersionBundleId": "test.app",
                         "softwareVersionExternalIdentifiers": ["999", 888]], "sinfs": []]]]
        let download = try StoreParsing.download(root, app: app, version: "999", email: "test@example.invalid")
        XCTAssertEqual(download.availableVersionIDs, ["999", "888"])
        XCTAssertThrowsError(try StoreParsing.download(root, app: app, version: "888", email: "test@example.invalid"))
    }
    func testCDNRejectsInsecureAndLookalikeHosts() throws {
        for url in ["http://iosapps.itunes.apple.com/a", "https://apple.com.evil.test/a", "https://user@iosapps.itunes.apple.com/a"] {
            XCTAssertThrowsError(try CDNPolicy.validate(url))
        }
        XCTAssertNoThrow(try CDNPolicy.validate("https://iosapps.itunes.apple.com/a"))
    }
    func testLicenseAndSessionErrorsRemainSpecific() {
        XCTAssertThrowsError(try StoreParsing.failure(["failureType": "9610"])) { XCTAssertEqual($0 as? StoreError, .licenseRequired) }
        XCTAssertThrowsError(try StoreParsing.failure(["failureType": "2034"])) { XCTAssertEqual($0 as? StoreError, .sessionExpired) }
    }
}
