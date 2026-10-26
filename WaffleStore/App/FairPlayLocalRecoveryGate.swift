import Foundation

// Reject all empty, partial or unverified native outputs. A successful C
// bridge status by itself is NOT evidence that a model was recovered.
enum FairPlayLocalRecoveryGate {
    static func accepts(_ envelope: [String: Any], outputBytes: Int64) -> Bool {
        acceptsRecoveredModel(envelope, operation: "recover-capture", outputBytes: outputBytes)
    }

    static func acceptsNativeEnrollment(_ envelope: [String: Any], outputBytes: Int64) -> Bool {
        acceptsRecoveredModel(envelope, operation: "recover-enrollment", outputBytes: outputBytes)
    }

    private static func acceptsRecoveredModel(_ envelope: [String: Any],
                                              operation: String, outputBytes: Int64) -> Bool {
        guard envelope["operation"] as? String == operation,
              envelope["ok"] as? Bool == true,
              let result = envelope["result"] as? [String: Any],
              result["completed"] as? Bool == true,
              result["portable_source_derived"] as? Bool == true,
              result["enrollment_frame_derived"] as? Bool == true,
              result["source_frame_cross_checked"] as? Bool == true,
              result["model_sinf_validated"] as? Bool == true,
              result["plaintext_authenticated"] as? Bool == false,
              let images = result["images"] as? Int, images > 0,
              let pages = result["pages"] as? Int, pages > 0,
              outputBytes > 0 else {
            return false
        }
        return true
    }

    static func acceptsModelDecrypt(_ envelope: [String: Any], outputBytes: Int64) -> Bool {
        guard envelope["operation"] as? String == "decrypt",
              envelope["ok"] as? Bool == true,
              let result = envelope["result"] as? [String: Any],
              result["completed"] as? Bool == true,
              result["plaintext_authenticated"] as? Bool == false,
              let images = result["images"] as? Int, images > 0,
              let pages = result["pages"] as? Int, pages > 0,
              outputBytes > 0 else {
            return false
        }
        return true
    }
}
