import Foundation

// A successful Go bridge call is not proof that portable recovery completed.
// Only permit an analysis IPA from this iPhone's retained capture when every
// native recovery gate is present and the output file is nonempty.
// This does NOT authenticate the plaintext; the caller must label it as such.
enum FairPlayLocalRecoveryGate {
    static func accepts(_ envelope: [String: Any], outputBytes: Int64) -> Bool {
        guard envelope["operation"] as? String == "recover-capture",
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
}
