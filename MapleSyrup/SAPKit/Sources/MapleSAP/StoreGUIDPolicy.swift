import Foundation

// StoreAgent Bag.sendGUIDWithURL: and CommerceKit's standard query builder.
// These rules affect the URL as well as the independently serialized plist.
// Never replace Apple routing parameters or change the signed request body.
enum StoreGUIDPolicy {
    static func authorizationURL(_ url: URL, bag: [String: Any], identity: MachineIdentity) throws -> URL {
        guard var components = URLComponents(url: url, resolvingAgainstBaseURL: false) else {
            throw FairPlayAuthorizationError.endpointRejected
        }
        let existing = (components.queryItems ?? []).filter { $0.name.lowercased() == "guid" }
        guard existing.count <= 1,
              existing.allSatisfy({ $0.name == "guid" && $0.value?.uppercased() == identity.guid }) else {
            throw FairPlayAuthorizationError.accountMismatch
        }
        guard try requiresGUID(url, bag: bag), existing.isEmpty else { return url }
        // Append without round-tripping existing queryItems: encoded routing
        // values, plus signs and duplicate non-GUID parameters stay byte exact.
        let query = components.percentEncodedQuery ?? ""
        components.percentEncodedQuery = query + (query.isEmpty ? "" : "&") + "guid=" + identity.guid
        guard let routed = components.url else { throw FairPlayAuthorizationError.endpointRejected }
        return routed
    }

    private static func requiresGUID(_ url: URL, bag: [String: Any]) throws -> Bool {
        // The original HTTPS fallback applies when guid-urls is not a dict.
        guard let policy = bag["guid-urls"] as? [String: Any] else { return url.scheme == "https" }
        if let schemes = policy["schemes"] as? [Any] {
            guard schemes.count <= 128 else { throw FairPlayAuthorizationError.endpointUnavailable }
            guard schemes.compactMap({ $0 as? String }).contains(url.scheme ?? "") else { return false }
        }
        guard let patterns = policy["regex"] as? [Any] else { return false }
        guard patterns.count <= 128 else { throw FairPlayAuthorizationError.endpointUnavailable }
        let text = url.absoluteString
        let range = NSRange(text.startIndex..<text.endIndex, in: text)
        for pattern in patterns.compactMap({ $0 as? String }) {
            guard pattern.utf8.count <= 4096,
                  let regex = try? NSRegularExpression(pattern: pattern, options: [.caseInsensitive]) else {
                throw FairPlayAuthorizationError.endpointUnavailable
            }
            if regex.firstMatch(in: text, range: range) != nil { return true }
        }
        return false
    }
}
