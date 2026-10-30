import Foundation

// A negative default return or an aggregate error is not a measured syscall.
// Legacy reports omit the attempted flags and must stay explicitly ambiguous.
nonisolated enum SystemNativeFailureStage {
    static func classify(_ details: [String: Any]) -> String {
        let mapping = (details["native_mapping_attempted"] as? NSNumber)?.boolValue
        let remap = (details["native_remap_attempted"] as? NSNumber)?.boolValue
        let mappingError = (details["native_mapping_errno"] as? NSNumber)?.intValue ?? 0
        let returned = (details["native_remap_return"] as? NSNumber)?.intValue ?? -1
        let readError = (details["native_read_error"] as? NSNumber)?.intValue ?? 0
        if mapping == true && mappingError != 0 { return "native_mapping" }
        if remap == true && returned != 0 { return "native_remap" }
        if remap == true && returned == 0 && readError != 0 { return "native_pages" }
        if mapping == nil || remap == nil { return "native_mapping_or_remap" }
        if mapping == false { return "native_source_preparation" }
        return "native_incomplete"
    }
}
