import Foundation
import ZIPFoundation

enum DeviceOperationError: LocalizedError {
    case message(String)
    var errorDescription: String? { if case let .message(s) = self { return s }; return nil }
}
// Only the FairPlay syscall boundary is stubbed. Exercise the production archive,
// staging, cancellation and publication code without an iOS license/keychain.
func crumb_group_signed_groups() -> [String]? { nil }
func crumb_group_decrypt(_ source: URL, _ destination: URL, _ encrypted: UnsafeMutablePointer<ObjCBool>, _ error: NSErrorPointer) -> Bool {
    do {
        let bytes = try Data(contentsOf: source)
        encrypted.pointee = ObjCBool(bytes == Data("encrypted-fixture".utf8))
        try (encrypted.pointee.boolValue ? Data("verified-plaintext".utf8) : bytes).write(to: destination)
        return true
    } catch let caught { error?.pointee = caught as NSError; return false }
}
func require(_ condition: @autoclosure () throws -> Bool, _ message: String) throws {
    if try !condition() { throw DeviceOperationError.message(message) }
}
let fm = FileManager.default
let root = fm.temporaryDirectory.appendingPathComponent(UUID().uuidString, isDirectory: true)
try fm.createDirectory(at: root, withIntermediateDirectories: true)
defer { try? fm.removeItem(at: root) }
let tree = root.appendingPathComponent("input")
let app = tree.appendingPathComponent("Payload/Test.app")
try fm.createDirectory(at: app.appendingPathComponent("SC_Info"), withIntermediateDirectories: true)
try Data("encrypted-fixture".utf8).write(to: app.appendingPathComponent("Test"))
let license = Data("original-license".utf8)
try license.write(to: app.appendingPathComponent("SC_Info/Test.sinf"))
let sinfManifest: [String: Any] = [
    "SinfPaths": ["SC_Info/Test.sinf"],
    "SinfReplicationPaths": ["PlugIns/Widget.appex/SC_Info/Widget.sinf", "Frameworks/A Framework.framework/SC_Info/A Framework.sinf", "", "../escape.sinf", "/escape.sinf", "obsolete-extension"]
]
try PropertyListSerialization.data(fromPropertyList: sinfManifest, format: .binary, options: 0)
    .write(to: app.appendingPathComponent("SC_Info/Manifest.plist"))
try Data("root-resource".utf8).write(to: tree.appendingPathComponent("iTunesMetadata.plist"))
let ipa = root.appendingPathComponent("source.ipa")
try fm.zipItem(at: tree, to: ipa, shouldKeepParent: false)
let original = try Data(contentsOf: ipa)
let output = try AppGroupDecryption.run(source: ipa, root: root, cancelled: { false }, progress: { _ in })
defer { try? fm.removeItem(at: output) }
let unpack = root.appendingPathComponent("check")
try fm.unzipItem(at: output, to: unpack)
print("Verified archive paths:", try fm.subpathsOfDirectory(atPath: unpack.path))
try require(try Data(contentsOf: unpack.appendingPathComponent("Payload/Test.app/Test")) == Data("verified-plaintext".utf8), "plaintext missing")
try require(try Data(contentsOf: unpack.appendingPathComponent("Payload/Test.app/SC_Info/Test.sinf")) == license, "license changed")
for path in (sinfManifest["SinfReplicationPaths"] as! [String]).prefix(2) {
    try require(try Data(contentsOf: unpack.appendingPathComponent("Payload/Test.app/" + path)) == license, "nested SINF replication lost")
}
try require(fm.fileExists(atPath: unpack.appendingPathComponent("iTunesMetadata.plist").path), "root resource lost")
try require(try Data(contentsOf: ipa) == original, "original IPA changed")
func reject(_ input: URL, cancelled: Bool = false) throws {
    do {
        let unexpected = try AppGroupDecryption.run(source: input, root: root, cancelled: { cancelled }, progress: { _ in })
        try? fm.removeItem(at: unexpected)
        throw DeviceOperationError.message("unexpected successful result")
    } catch DeviceOperationError.message(let text) where text == "unexpected successful result" { throw DeviceOperationError.message(text) }
    catch { }
}
try reject(ipa, cancelled: true)
let malicious = root.appendingPathComponent("traversal.ipa")
let archive = try Archive(url: malicious, accessMode: .create)
try archive.addEntry(with: "../escaped", fileURL: app.appendingPathComponent("Test"))
try reject(malicious)
try require(!fm.fileExists(atPath: root.appendingPathComponent("escaped").path), "path traversal escaped")
try require(!fm.fileExists(atPath: root.appendingPathComponent("escape.sinf").path), "unsafe replication hint escaped staging")
let plainTree = root.appendingPathComponent("plain-input")
try fm.createDirectory(at: plainTree.appendingPathComponent("Payload/Plain.app"), withIntermediateDirectories: true)
try Data("already plain".utf8).write(to: plainTree.appendingPathComponent("Payload/Plain.app/Test"))
let plainIPA = root.appendingPathComponent("plain.ipa")
try fm.zipItem(at: plainTree, to: plainIPA, shouldKeepParent: false)
try reject(plainIPA)
let remaining = try fm.contentsOfDirectory(atPath: root.path).filter { $0.hasPrefix("crumb-fairplay-") }
try require(remaining.isEmpty, "staging workspace leaked")
print("App Group archive tests passed: plaintext round trip, original/license/root preservation, traversal rejection, cancellation, no-encryption rejection, workspace cleanup.")

try require(AppGroupDecryption.candidates(signed: ["group.signed"], profile: ["group.*", "group.signed", "group.profile"], configured: [" group.manual ", ""]) == ["group.signed", "group.manual", "group.profile"], "effective groups must take priority; wildcard must not suppress them")
try require(AppGroupDecryption.candidates(signed: nil, profile: ["group.profile"], configured: []) == ["group.profile"], "runtime-unavailable fallback failed")
try require(AppGroupDecryption.candidates(signed: [], profile: [], configured: []) == [], "invented an App Group")
print("Effective-entitlement group discovery tests passed.")

for base in [URL(fileURLWithPath: "/var/mobile/Containers/Shared/AppGroup/fixture/Test.app", isDirectory: true),
             URL(fileURLWithPath: "/private/var/mobile/Containers/Shared/AppGroup/fixture/Test.app", isDirectory: true)] {
    let target = try AppGroupDecryption.sinfURL("PlugIns/Widget.appex/SC_Info/Widget.sinf", in: base)
    try require(target.lastPathComponent == "Widget.sinf", "valid Darwin SINF path rejected")
}
for path in ["", "/SC_Info/x.sinf", "../x.sinf", "SC_Info/../../x.sinf", "SC_Info//x.sinf", "./x.sinf", "SC_Info/..", "SC_Info/x\\y.sinf", "SC_Info/x\0.sinf"] {
    do { _ = try AppGroupDecryption.sinfURL(path, in: app); fatalError("unsafe SINF path accepted") }
    catch { }
}
let link = app.appendingPathComponent("Linked")
try fm.createSymbolicLink(at: link, withDestinationURL: root)
do { _ = try AppGroupDecryption.sinfURL("Linked/escape.sinf", in: app); fatalError("SINF symlink accepted") }
catch { }
print("SINF tests passed: existing source, missing nested destinations, spaces, Darwin aliases, traversal and symlink rejection.")
