import Foundation

@main
struct NativeFailureContracts {
    static func main() throws {
        // Actual sanitized 23057 report. Its default -1 does not distinguish
        // mmap(ENOMEM) from a measured mremap(ENOMEM).
        let legacy = try JSONSerialization.jsonObject(with: Data("""
        {"native_remap_return":-1,"native_remap_errno":12,
         "native_read_error":0,"pages_total":0,"pages_read":0,"pages_changed":0}
        """.utf8)) as! [String: Any]
        precondition(SystemNativeFailureStage.classify(legacy) == "native_mapping_or_remap")
        var mapping = legacy
        mapping["native_mapping_attempted"] = true
        mapping["native_mapping_errno"] = 12
        mapping["native_remap_attempted"] = false
        precondition(SystemNativeFailureStage.classify(mapping) == "native_mapping")
        var remap = mapping
        remap["native_mapping_errno"] = 0
        remap["native_remap_attempted"] = true
        precondition(SystemNativeFailureStage.classify(remap) == "native_remap")
        remap["native_remap_return"] = 0
        remap["native_read_error"] = 5
        precondition(SystemNativeFailureStage.classify(remap) == "native_pages")
        remap["native_read_error"] = 0
        precondition(SystemNativeFailureStage.classify(remap) == "native_incomplete")
        mapping["native_mapping_attempted"] = false
        mapping["native_mapping_errno"] = 0
        precondition(SystemNativeFailureStage.classify(mapping) == "native_source_preparation")
        print("Native stage contracts passed; no FairPlay service or decryption was simulated.")
    }
}
