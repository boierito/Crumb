import Foundation

public enum SAPError: Error, LocalizedError, Equatable {
    case invalidBag, unsupportedVersion, invalidEndpoint, invalidIdentity
    case http(Int), oversizedResponse, invalidCertificate, invalidExchange
    case runtimeUnavailable, executableRuntimeRejected, invalidState, emptySignature
    case keychain(Int32)

    public var errorDescription: String? {
        switch self {
        case .invalidBag: return "Store Bag does not contain a valid SAP configuration."
        case .unsupportedVersion: return "Store Bag requires an unsupported SAP version."
        case .invalidEndpoint: return "Apple endpoint rejected: HTTPS and a trusted Apple host are required."
        case .invalidIdentity: return "Machine identity is invalid."
        case .http(let status): return "Apple returned HTTP \(status)."
        case .oversizedResponse: return "Apple SAP response exceeds the 1 MiB limit."
        case .invalidCertificate: return "Apple returned no SAP certificate."
        case .invalidExchange: return "Apple returned an invalid SAP setup exchange."
        case .runtimeUnavailable: return "SAP guest interpreter is not implemented. No ActionSignature was generated. See IMPLEMENTATION.md."
        case .executableRuntimeRejected: return "A runtime that generates executable memory is incompatible with this jailed SAP module."
        case .invalidState: return "SAP operation called in an invalid session state."
        case .emptySignature: return "SAP runtime returned an empty signature."
        case .keychain(let status): return "Keychain operation failed (OSStatus \(status))."
        }
    }
}

public struct MachineIdentity: Equatable {
    public let hardwareID: Data
    public var guid: String { hardwareID.map { String(format: "%02X", $0) }.joined() }
    public init(hardwareID: Data) throws {
        guard hardwareID.count == 6 else { throw SAPError.invalidIdentity }
        self.hardwareID = hardwareID
    }
}

public struct SAPConfiguration: Equatable {
    public let authenticationURL: URL
    public let setupURL: URL
    public let certificateURL: URL
    public let version: UInt32

    public static func parse(bag data: Data) throws -> SAPConfiguration {
        guard let root = try PropertyListSerialization.propertyList(from: data, format: nil) as? [String: Any],
              let bag = root["urlBag"] as? [String: Any],
              let authentication = bag["authenticateAccount"] as? String,
              let setup = bag["sign-sap-setup"] as? String,
              let certificate = bag["sign-sap-setup-cert"] as? String,
              let rawVersion = bag["sign-sap-version"] else { throw SAPError.invalidBag }
        let versionString = (rawVersion as? String) ?? (rawVersion as? NSNumber)?.stringValue
        guard let string = versionString, let version = UInt32(string) else { throw SAPError.invalidBag }
        guard version == 200 else { throw SAPError.unsupportedVersion }
        let authURL = try trustedAppleURL(authentication)
        let host = authURL.host!.lowercased()
        guard (host == "auth.itunes.apple.com" || host.hasSuffix("-buy.itunes.apple.com")),
              ["/auth/v1/native", "/auth/v1/native/"].contains(authURL.path) else { throw SAPError.invalidEndpoint }
        return SAPConfiguration(authenticationURL: authURL, setupURL: try trustedAppleURL(setup),
                                certificateURL: try trustedAppleURL(certificate), version: version)
    }

    public static func trustedAppleURL(_ string: String) throws -> URL {
        guard let url = URL(string: string), url.scheme == "https",
              let host = url.host?.lowercased(),
              (host == "apple.com" || host.hasSuffix(".apple.com") || host.hasSuffix(".mzstatic.com")),
              url.user == nil, url.password == nil, url.fragment == nil,
              url.port == nil || url.port == 443 else { throw SAPError.invalidEndpoint }
        return url
    }
}
