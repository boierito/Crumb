import Foundation

@main
struct CoreMediaLocalProbeTests {
    static func main() throws {
        let fm = FileManager.default
        let parent = fm.temporaryDirectory.appendingPathComponent(UUID().uuidString, isDirectory: true)
        let home = parent.appendingPathComponent("Crumb", isDirectory: true)
        let own = home.appendingPathComponent("Downloads", isDirectory: true)
        try fm.createDirectory(at: own, withIntermediateDirectories: true)
        defer { try? fm.removeItem(at: parent) }

        let ipa = own.appendingPathComponent("licensed.IPA")
        try Data(repeating: 0x42, count: 64).write(to: ipa)
        precondition(CoreMediaLocalProbe.isEligible(ipa, home: home), "own IPA must be eligible")

        let otherExtension = own.appendingPathComponent("not-an-ipa.txt")
        try Data(repeating: 0x42, count: 64).write(to: otherExtension)
        precondition(!CoreMediaLocalProbe.isEligible(otherExtension, home: home), "reject non-IPA")

        let outside = parent.appendingPathComponent("other.IPA")
        try Data(repeating: 0x42, count: 64).write(to: outside)
        precondition(!CoreMediaLocalProbe.isEligible(outside, home: home), "reject external container")

        let link = own.appendingPathComponent("linked.IPA")
        try fm.createSymbolicLink(at: link, withDestinationURL: outside)
        precondition(!CoreMediaLocalProbe.isEligible(link, home: home), "reject symlink outside")

        let missing = own.appendingPathComponent("missing.IPA")
        precondition(!CoreMediaLocalProbe.isEligible(missing, home: home), "reject missing input")
        precondition(!CoreMediaLocalProbe.isEligible(URL(string: "https://example.test/app.ipa")!, home: home),
                     "reject remote URLs")

        // Host builds must not report a successful CoreMedia or FairPlay read.
        let host = CoreMediaLocalProbe.run(input: outside)
        precondition(host.stage == "blocked" && host.comparedBytes == 0, "fail closed")
        precondition(!host.summary.contains(outside.path), "redact filesystem paths")
        print("CoreMedia local diagnostics: 7 checks passed; no exploit or decryption claimed.")
    }
}
