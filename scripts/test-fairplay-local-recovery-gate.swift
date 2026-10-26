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
        print("FairPlay iPhone local recovery publish gate: passed")
    }
}
