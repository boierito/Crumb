import Foundation

@main
struct RuntimeIPAImportTests {
    static func main() throws {
        let fm = FileManager.default
        let scratch = fm.temporaryDirectory.appendingPathComponent("crumb-import-test-" + UUID().uuidString)
        try fm.createDirectory(at: scratch, withIntermediateDirectories: true)
        defer { try? fm.removeItem(at: scratch) }
        let input = scratch.appendingPathComponent("Purchased.ipa")
        let original = Data([0x50, 0x4b, 0x03, 0x04] + Array(repeating: UInt8(0xBC), count: 2048))
        try original.write(to: input)
        let importedRoot = scratch.appendingPathComponent("Private")
        let first = try RuntimeIPAImport.copyPickedIPA(input, into: importedRoot)
        let second = try RuntimeIPAImport.copyPickedIPA(input, into: importedRoot)
        precondition(first != second, "imports must not overwrite each other")
        precondition(first.deletingLastPathComponent() == importedRoot, "private destination")
        precondition(try Data(contentsOf: first) == original, "imported bytes must match")
        precondition(try Data(contentsOf: input) == original, "original must remain unchanged")

        let badExtension = scratch.appendingPathComponent("NotAnIPA.zip")
        try original.write(to: badExtension)
        do {
            try RuntimeIPAImport.preflight(badExtension)
            fatalError("wrong extension accepted")
        } catch RuntimeIPAImport.ImportFailure.invalidInput { }

        let broken = scratch.appendingPathComponent("Broken.ipa")
        try Data(repeating: 0x00, count: 128).write(to: broken)
        do {
            try RuntimeIPAImport.preflight(broken)
            fatalError("invalid archive header accepted")
        } catch RuntimeIPAImport.ImportFailure.invalidArchive { }

        let linked = scratch.appendingPathComponent("Linked.ipa")
        try fm.createSymbolicLink(at: linked, withDestinationURL: input)
        do {
            try RuntimeIPAImport.preflight(linked)
            fatalError("symbolic link accepted")
        } catch RuntimeIPAImport.ImportFailure.invalidInput { }

        let directory = scratch.appendingPathComponent("Folder.ipa")
        try fm.createDirectory(at: directory, withIntermediateDirectories: true)
        do {
            try RuntimeIPAImport.preflight(directory)
            fatalError("directory accepted")
        } catch RuntimeIPAImport.ImportFailure.invalidInput { }

        print("Runtime IPA import: 8 positive and rejection assertions passed.")
    }
}
