//! Bounded, opt-in self-device diagnostics. Never stages an IPA or a license.
//! Uses the authenticated LocalDevVPN/RemotePairing tunnel established by transport.
use super::*;
use idevice::RsdService;
use idevice::services::{
    afc::{AfcClient, opcode::AfcFopenMode},
    core_device::{Domain, FileServiceClient},
    os_trace_relay::OsTraceRelayClient,
};
use idevice::tcp::adapter::AdapterHandle;
use idevice::rsd::RsdHandshake;
use serde_json::{json, Value};
use std::{fs::File, os::fd::AsRawFd, path::Path, time::Duration};

// The fixed staging prefix is separate from Crumb's installer and contains only a
// generated canary. Never accept a user-supplied AFC path or stage an IPA/SINF.
const STAGE_ROOT: &str = "PublicStaging/Crumb-VPN-FairPlay-Probe";
const SELF_PROOF_NAME: &str = ".crumb-device-proof";

fn map_and_compare(path: &Path, expected: &[u8]) -> Value {
    if expected.is_empty() || expected.len() > 1024 {
        return json!({"open":false,"map":false,"reason":"invalid_canary"});
    }
    let file = match File::open(path) {
        Ok(f) => f,
        Err(e) => return json!({"open":false,"map":false,"errno":e.raw_os_error()}),
    };
    if file.metadata().map(|m| m.len()).ok() != Some(expected.len() as u64) {
        return json!({"open":true,"map":false,"reason":"size_mismatch"});
    }
    // MAP_PRIVATE / PROT_READ: never executes code or alters a bundle.
    let ptr = unsafe {
        libc::mmap(std::ptr::null_mut(), expected.len(), libc::PROT_READ,
                   libc::MAP_PRIVATE, file.as_raw_fd(), 0)
    };
    if ptr == libc::MAP_FAILED {
        return json!({"open":true,"map":false,
                      "errno":std::io::Error::last_os_error().raw_os_error()});
    }
    let equal = unsafe {
        std::slice::from_raw_parts(ptr as *const u8, expected.len()) == expected
    };
    unsafe { libc::munmap(ptr, expected.len()); }
    json!({"open":true,"map":true,"matches":equal})
}

fn valid_proof_path(path: &Path) -> bool {
    path.file_name().is_some_and(|v| v == SELF_PROOF_NAME)
        && path.parent().and_then(|p| p.file_name()).is_some_and(|v| v == "Documents")
}

async fn afc_canary(adapter: &mut AdapterHandle, handshake: &mut RsdHandshake) -> Value {
    let mut afc = match AfcClient::connect_rsd(adapter, handshake).await {
        Ok(afc) => afc,
        Err(e) => return json!({"status":"afc_unavailable","code":e.code()}),
    };
    // Cleanup an interrupted previous diagnostic, never the installation stage.
    let _ = afc.mk_dir("PublicStaging").await;
    if afc.get_file_info(STAGE_ROOT).await.is_ok() {
        if afc.remove_all(STAGE_ROOT).await.is_err() {
            return json!({"status":"previous_probe_cleanup_denied"});
        }
    }
    if afc.mk_dir(STAGE_ROOT).await.is_err() {
        return json!({"status":"afc_staging_denied"});
    }
    let result = async {
        let filename = format!("{STAGE_ROOT}/probe.bin");
        let data = format!("CRUMB_VPN_PROBE_{}", uuid::Uuid::new_v4());
        let mut out = afc.open(&filename, AfcFopenMode::WrOnly).await?;
        out.write_entire(data.as_bytes()).await?;
        out.close().await?;
        let mut input = afc.open(&filename, AfcFopenMode::RdOnly).await?;
        let downloaded = input.read_n(data.len()).await?;
        input.close().await?;
        if downloaded != data.as_bytes() {
            return Ok::<_,idevice::IdeviceError>(json!({"status":"afc_roundtrip_mismatch"}));
        }
        // Media is an AFC domain, not necessarily accessible to a sandboxed app.
        // This is a fixed, read-only local path probe, not a sandbox bypass.
        let mapped = map_and_compare(
            Path::new("/var/mobile/Media/PublicStaging/Crumb-VPN-FairPlay-Probe/probe.bin"),
            data.as_bytes());
        Ok(json!({"status":"afc_roundtrip_ok","app_local_mapping":mapped}))
    }.await;
    let cleaned = afc.remove_all(STAGE_ROOT).await.is_ok();
    match result {
        Ok(mut report) => {report["cleanup_ok"]=json!(cleaned); report},
        Err(e) => json!({"status":"afc_probe_failed","code":e.code(),"cleanup_ok":cleaned}),
    }
}

async fn core_device_check(adapter: &mut AdapterHandle, handshake: &mut RsdHandshake, own: &str) -> Value {
    let action = async {
        let mut client: FileServiceClient<Box<dyn idevice::ReadWrite>> =
            FileServiceClient::connect_rsd(adapter, handshake).await?;
        client.create_session(Domain::AppDataContainer, own).await?;
        let list = client.retrieve_directory_list("Documents").await?;
        Ok::<usize,idevice::IdeviceError>(list.len())
    };
    match tokio::time::timeout(Duration::from_secs(8), action).await {
        Ok(Ok(count)) => json!({"status":"own_container_list_ok","entry_count":count}),
        Ok(Err(e)) => json!({"status":"unavailable_or_denied","code":e.code()}),
        Err(_) => json!({"status":"timeout"}),
    }
}

async fn trace_check(adapter: &mut AdapterHandle, handshake: &mut RsdHandshake) -> Value {
    let action = async {
        let client = OsTraceRelayClient::connect_rsd(adapter, handshake).await?;
        let mut receiver = client.start_trace(None).await?;
        let end = tokio::time::Instant::now() + Duration::from_secs(3);
        let mut total = 0u32;
        let mut fairplay = 0u32;
        let mut policy = 0u32;
        while tokio::time::Instant::now() < end && total < 128 {
            match tokio::time::timeout(Duration::from_millis(450), receiver.next()).await {
                Ok(Ok(log)) => {
                    total += 1;
                    // Never return raw log strings, paths, identifiers or licenses.
                    let name = log.filename.to_ascii_lowercase();
                    let msg = log.message.to_ascii_lowercase();
                    if name.contains("fairplay") || msg.contains("fairplay") {fairplay += 1;}
                    if msg.contains("system policy") || msg.contains("sandbox") {policy += 1;}
                }
                Ok(Err(_)) => break,
                Err(_) => continue,
            }
        }
        Ok::<_,idevice::IdeviceError>((total, fairplay, policy))
    };
    match tokio::time::timeout(Duration::from_secs(8), action).await {
        Ok(Ok((total, fairplay, policy))) =>
            json!({"status":"passive_sample_only","events":total,
                   "fairplay_mentions":fairplay,"policy_mentions":policy}),
        Ok(Err(e)) => json!({"status":"unavailable_or_denied","code":e.code()}),
        Err(_) => json!({"status":"timeout"}),
    }
}

pub(super) async fn run(
    job: &std::sync::Arc<Job>,
    adapter: &mut AdapterHandle,
    handshake: &mut RsdHandshake,
    request: &Value,
    own: &str,
) -> anyhow::Result<Value> {
    let path = request["self_probe_path"].as_str()
        .ok_or_else(|| anyhow::anyhow!("Missing local proof path"))?;
    anyhow::ensure!(valid_proof_path(Path::new(path)), "Invalid local proof path");
    let proof = request["proof"].as_str()
        .ok_or_else(|| anyhow::anyhow!("Missing local proof"))?;
    anyhow::ensure!(uuid::Uuid::parse_str(proof).is_ok(), "Invalid proof format");
    job.update(json!({"phase":"probing"}));
    let local = map_and_compare(Path::new(path), proof.as_bytes());
    if job.cancelled.load(std::sync::atomic::Ordering::SeqCst) {
        anyhow::bail!("Cancelled");
    }
    let available = json!({
        "afc":handshake.services.contains_key("com.apple.afc.shim.remote"),
        "file_service":handshake.services.contains_key("com.apple.coredevice.fileservice.control"),
        "os_trace":handshake.services.contains_key("com.apple.os_trace_relay.shim.remote"),
    });
    let afc = afc_canary(adapter, handshake).await;
    if job.cancelled.load(std::sync::atomic::Ordering::SeqCst) {
        anyhow::bail!("Cancelled");
    }
    let core = core_device_check(adapter, handshake, own).await;
    if job.cancelled.load(std::sync::atomic::Ordering::SeqCst) {
        anyhow::bail!("Cancelled");
    }
    let trace = trace_check(adapter, handshake).await;
    Ok(json!({"phase":"probe_complete","report":{
        "schema":1, "mode":"no_helper_no_ipa_no_fairplay_remap",
        "self_documents_mmap":local, "rsd_advertised":available,
        "afc_public_staging":afc,"coredevice_own_documents":core,
        "ostrace_passive":trace,
        "conclusion":"AFC/CoreDevice/trace access does not establish fairplayd authorization."
    }}))
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn proof_path_is_restricted_to_own_documents_marker() {
        assert!(valid_proof_path(Path::new("/a/Documents/.crumb-device-proof")));
        assert!(!valid_proof_path(Path::new("/a/Documents/example.ipa")));
        assert!(!valid_proof_path(Path::new("/a/tmp/.crumb-device-proof")));
    }
    #[test]
    fn readonly_mapping_recognizes_own_canary() {
        let path = std::env::temp_dir().join(format!("crumb-probe-{}", uuid::Uuid::new_v4()));
        std::fs::write(&path, b"crumb-canary").unwrap();
        let good = map_and_compare(&path, b"crumb-canary");
        assert_eq!(good["map"], true);
        assert_eq!(good["matches"], true);
        let wrong = map_and_compare(&path, b"other-canary");
        assert_eq!(wrong["matches"], false);
        std::fs::remove_file(&path).unwrap();
        assert_eq!(map_and_compare(&path,b"crumb-canary")["open"],false);
    }
}
