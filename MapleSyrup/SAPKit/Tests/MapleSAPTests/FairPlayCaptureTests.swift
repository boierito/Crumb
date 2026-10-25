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
