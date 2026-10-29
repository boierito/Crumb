import XCTest
@testable import MapleSAP

final class FairPlayCaptureTests: XCTestCase {
    func testKeybagLocationsAndEncodings() {
        let data = Data([1,2,3,4])
        for key in ["keybag", "subKeyBag", "sbKeyBag"] {
            XCTAssertEqual(FairPlayTransaction.keybag([key: data]), data)
            XCTAssertEqual(FairPlayTransaction.keybag(["songList": [[key: data.base64EncodedString()]]]), data)
            XCTAssertEqual(FairPlayTransaction.keybag(["songList": [["metadata": [key: data]]]]), data)
        }
        XCTAssertNil(FairPlayTransaction.keybag(["keybag": "not base64"]))
        XCTAssertNil(FairPlayTransaction.keybag(["keybag": Data()]))
        XCTAssertNil(FairPlayTransaction.keybag([:]))
    }
}

extension FairPlayCaptureTests {
    func testLicenseBagsKeepKindsSeparateAndAcceptExactAliases() throws {
        let key = Data([1, 2]), diversity = Data([3, 4])
        let bags = try FairPlayTransaction.licenseBags(["keybag": key, "songList": [["sbKeyBag": key.base64EncodedString(), "metadata": ["diversitybag": diversity]]]])
        XCTAssertEqual(bags.keybag, key)
        XCTAssertEqual(bags.diversitybag, diversity)
        XCTAssertNil(try FairPlayTransaction.licenseBags([:]).diversitybag)
    }
    func testLicenseBagsRejectAmbiguityMalformedAndOversizedValues() {
        let roots: [[String: Any]] = [
            ["keybag": Data([1]), "songList": [["keybag": Data([2])]]],
            ["keybag": Data([1]), "diversitybag": "invalid base64"],
            ["diversitybag": Data()],
            ["keybag": Data(repeating: 1, count: (1 << 20) + 1)]
        ]
        for root in roots { XCTAssertThrowsError(try FairPlayTransaction.licenseBags(root)) }
    }
}
