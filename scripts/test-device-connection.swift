import Foundation
@main struct ConnectionTests {
    static func main() throws {
        let credentials: [String:Any] = ["record":"PAIRING-PLIST", "operation":"fetch", "route":"10.7.0.1"]
        // Reproduce the old saved-job collision: JSON download metadata named record.
        let pending = ["record":"DOWNLOAD-JSON", "download_record":"DOWNLOAD-JSON", "job_id":"job", "helper_bundle":"helper", "output":"/result.ipa", "install":"true", "cancel_protocol":"1"]
        let request = try DeviceRequest.adding(DeviceRequest.retrievalParameters(pending), to:credentials)
        precondition(request["record"] as? String == "PAIRING-PLIST")
        precondition(request["job_id"] as? String == "job")
        precondition(request["download_record"] == nil && request["install"] == nil && request["cancel_protocol"] == nil)
        for key in DeviceRequest.reserved {
            do { _ = try DeviceRequest.adding([key:"overridden"],to:credentials); fatalError("Credential override allowed") }
            catch { }
        }
        var policy = DeviceConnectionPolicy()
        func check(_ busy:Bool = false, _ helper:Bool = false) -> Bool {
            policy.shouldCheck(paired:true,busy:busy,connected:false,cellular:false,helperPending:helper)
        }
        precondition(policy.networkChanged("wifi"))
        precondition(check())
        for _ in 0..<100 { precondition(!check()); precondition(!policy.networkChanged("wifi")) }
        policy.becameInactive()
        precondition(!check(true)) // a callback must not start a parallel native job
        precondition(!check(false,true)) // returning from helper must not trigger a probe
        precondition(check())
        precondition(!check())
        precondition(policy.networkChanged("cellular"))
        precondition(!policy.shouldCheck(paired:true,busy:false,connected:false,cellular:true,helperPending:false))
        precondition(policy.networkChanged("wifi")); precondition(check())
        var result = DeviceResultRequest()
        result.enqueue()
        precondition(!result.claim(busy:true,hasPendingJob:true))
        precondition(result.requested)
        result.enqueue() // duplicate URL and foreground events coalesce
        precondition(result.claim(busy:false,hasPendingJob:true))
        precondition(!result.claim(busy:false,hasPendingJob:true))
        result.enqueue()
        precondition(!result.claim(busy:false,hasPendingJob:false))
        precondition(!result.requested) // completing a job must not suppress future connection checks
        precondition(!result.claim(busy:false,hasPendingJob:true))
        result.enqueue()
        precondition(result.claim(busy:false,hasPendingJob:true))
        precondition(DeviceJobState.terminalPhases.contains("signing-plan"))
        precondition(DeviceJobState.terminalPhases.contains("installation-signing-ready"))
        precondition(DeviceJobState.terminalPhases.contains("cancelled"))
        precondition(DeviceJobState.terminalPhases.contains("prepared"))
        for phase in ["starting", "preparing", "fetching", "installing", "signing"] {
            precondition(!DeviceJobState.terminalPhases.contains(phase))
        }
        print("Connection scheduling and saved-result credential regressions passed")
    }
}
