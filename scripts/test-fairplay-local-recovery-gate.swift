import Foundation

@main
struct FairPlayLocalRecoveryGateTests {
    static func main() {
        let successfulResult: [String: Any] = [
            "completed": true,
            "portable_source_derived": true,
            "enrollment_frame_derived": true,
            "source_frame_cross_checked": true,
            "model_sinf_validated": true,
            "plaintext_authenticated": false,
            "images": 2,
            "pages": 712
        ]
        func envelope(_ result: [String: Any]) -> [String: Any] {
            ["operation": "recover-capture", "ok": true, "result": result]
        }
        precondition(FairPlayLocalRecoveryGate.accepts(envelope(successfulResult), outputBytes: 4096),
                     "complete and nonempty output was rejected")

        for flag in ["completed", "portable_source_derived", "enrollment_frame_derived",
                     "source_frame_cross_checked", "model_sinf_validated"] {
            var result = successfulResult
            result[flag] = false
            precondition(!FairPlayLocalRecoveryGate.accepts(envelope(result), outputBytes: 4096),
                         "accepted false gate: \(flag)")
            result.removeValue(forKey: flag)
            precondition(!FairPlayLocalRecoveryGate.accepts(envelope(result), outputBytes: 4096),
                         "accepted absent gate: \(flag)")
        }
        for key in ["images", "pages"] {
            var result = successfulResult
            result[key] = 0
            precondition(!FairPlayLocalRecoveryGate.accepts(envelope(result), outputBytes: 4096),
                         "accepted zero \(key)")
        }
        var result = successfulResult
        result["plaintext_authenticated"] = true
        precondition(!FairPlayLocalRecoveryGate.accepts(envelope(result), outputBytes: 4096),
                     "accepted unsupported plaintext authenticity claim")
        precondition(!FairPlayLocalRecoveryGate.accepts(envelope(successfulResult), outputBytes: 0),
                     "accepted absent processed IPA")

        var failed = envelope(successfulResult)
        failed["ok"] = false
        precondition(!FairPlayLocalRecoveryGate.accepts(failed, outputBytes: 4096),
                     "accepted native failure")
        failed = envelope(successfulResult)
        failed["operation"] = "recover-enrollment"
        precondition(!FairPlayLocalRecoveryGate.accepts(failed, outputBytes: 4096),
                     "accepted unrelated recovery route")
        failed.removeValue(forKey: "result")
        precondition(!FairPlayLocalRecoveryGate.accepts(failed, outputBytes: 4096),
                     "accepted missing recovery evidence")
        // The modern native enrollment route previously treated any ok=true
        // as sufficient, even when no processed IPA existed. Assert the same
        // strict model/frame/SINF gates that protect on-device recovery.
        var native = envelope(successfulResult)
        native["operation"] = "recover-enrollment"
        precondition(FairPlayLocalRecoveryGate.acceptsNativeEnrollment(native, outputBytes: 4096),
                     "rejected valid synthetic modern enrollment")
        precondition(!FairPlayLocalRecoveryGate.acceptsNativeEnrollment(native, outputBytes: 0),
                     "accepted enrollment without any processed IPA")
        for flag in ["completed", "portable_source_derived", "enrollment_frame_derived",
                     "source_frame_cross_checked", "model_sinf_validated"] {
            var broken = successfulResult
            broken[flag] = false
            var response = native
            response["result"] = broken
            precondition(!FairPlayLocalRecoveryGate.acceptsNativeEnrollment(response, outputBytes: 4096),
                         "accepted incomplete enrollment: \(flag)")
        }
        precondition(!FairPlayLocalRecoveryGate.acceptsNativeEnrollment(envelope(successfulResult), outputBytes: 4096),
                     "accepted the wrong operation type for native enrollment")

        let decrypted: [String: Any] = [
            "completed": true, "plaintext_authenticated": false, "images": 2, "pages": 712
        ]
        let direct: [String: Any] = ["operation": "decrypt", "ok": true, "result": decrypted]
        precondition(FairPlayLocalRecoveryGate.acceptsModelDecrypt(direct, outputBytes: 4096),
                     "valid synthetic model processing rejected")
        precondition(!FairPlayLocalRecoveryGate.acceptsModelDecrypt(direct, outputBytes: 0),
                     "published nonexistent direct IPA")
        precondition(!FairPlayLocalRecoveryGate.acceptsModelDecrypt(native, outputBytes: 4096),
                     "confused native and direct operations")
        for key in ["completed", "images", "pages", "plaintext_authenticated"] {
            var broken = decrypted
            broken.removeValue(forKey: key)
            var response = direct
            response["result"] = broken
            precondition(!FairPlayLocalRecoveryGate.acceptsModelDecrypt(response, outputBytes: 4096),
                         "accepted malformed decryption: \(key)")
        }
        print("FairPlay native, local and direct IPA publish gates: passed")
    }
}
