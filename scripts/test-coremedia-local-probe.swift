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

        // A second diagnostic target is allowed only if it is exactly
        // Crumb's own executable, not an arbitrary app bundle or file.
        let bundle = parent.appendingPathComponent("Crumb.app", isDirectory: true)
        try fm.createDirectory(at: bundle, withIntermediateDirectories: true)
        let ownExecutable = bundle.appendingPathComponent("Crumb")
        try Data(repeating: 0xAA, count: 64).write(to: ownExecutable)
        precondition(CoreMediaLocalProbe.isEligibleOwnExecutable(ownExecutable, bundleURL: bundle,
                                                               executableURL: ownExecutable),
                     "own installed executable accepted")
        precondition(!CoreMediaLocalProbe.isEligibleOwnExecutable(outside, bundleURL: bundle,
                                                                executableURL: ownExecutable),
                     "reject external executable")
        let linkedExecutable = bundle.appendingPathComponent("LinkedCrumb")
        try fm.createSymbolicLink(at: linkedExecutable, withDestinationURL: outside)
        precondition(!CoreMediaLocalProbe.isEligibleOwnExecutable(linkedExecutable, bundleURL: bundle,
                                                                  executableURL: linkedExecutable),
                     "reject symbolic executable")
        let nonExecutable = bundle.appendingPathComponent("Other")
        try Data(repeating: 0xAA, count: 64).write(to: nonExecutable)
        precondition(!CoreMediaLocalProbe.isEligibleOwnExecutable(nonExecutable, bundleURL: bundle,
                                                                executableURL: ownExecutable),
                     "reject non-identity matched resource")

        // Host builds must not report a successful CoreMedia or FairPlay read.
        let host = CoreMediaLocalProbe.run(input: outside)
        precondition(host.stage == "blocked" && host.comparedBytes == 0, "fail closed")
        precondition(!host.summary.contains(outside.path), "redact filesystem paths")
        print("CoreMedia local diagnostics: 11 checks passed; no exploit or decryption claimed.")
    }
}
