use super::*;
use idevice::services::afc::AfcClient;
use idevice::{
    RsdService,
    remote_pairing::{RemotePairingClient, connect_tls_psk_tunnel_native},
    rsd::RsdHandshake,
    services::afc::opcode::AfcFopenMode,
    services::{house_arrest::HouseArrestClient, installation_proxy::InstallationProxyClient},
    tcp::{adapter::Adapter, handle::AdapterHandle},
};
use std::net::{IpAddr, SocketAddr};
use tokio::io::{AsyncReadExt, AsyncWriteExt};

// Fixed stage labels and numeric protocol codes keep credentials and device
// response payloads out of user-facing errors while preserving socket causes.
fn device_error(stage: &'static str, error: idevice::IdeviceError) -> anyhow::Error {
    match error {
        idevice::IdeviceError::Socket(error) => anyhow::anyhow!(
            "{stage}: socket {:?} (OS {})",
            error.kind(),
            error
                .raw_os_error()
                .map(|code| code.to_string())
                .unwrap_or_else(|| "unavailable".into())
        ),
        error => anyhow::anyhow!(
            "{stage}: device protocol error {}.{}",
            error.code(),
            error.sub_code()
        ),
    }
}

// VendDocuments keeps the AFC root at the app container. Documents is an
// accessible child directory, not the root of the returned file service.
#[derive(Clone, Copy)]
enum AppDocument {
    Proof,
    Manifest,
    Output,
}
impl AppDocument {
    fn path(self) -> &'static str {
        match self {
            Self::Proof => "/Documents/.crumb-device-proof",
            Self::Manifest => "/Documents/result.json",
            Self::Output => "/Documents/result.ipa",
        }
    }
    fn stage(self) -> &'static str {
        match self {
            Self::Proof => "Device proof file",
            Self::Manifest => "Helper result manifest",
            Self::Output => "Helper result IPA",
        }
    }
}
async fn read_app_document(
    client: &mut AfcClient,
    document: AppDocument,
    limit: usize,
) -> anyhow::Result<Vec<u8>> {
    let mut file = client
        .open(document.path(), AfcFopenMode::RdOnly)
        .await
        .map_err(|e| device_error(document.stage(), e))?;
    let read = file.read_n(limit).await;
    let close = file.close().await;
    let bytes = read.map_err(|e| device_error(document.stage(), e))?;
    close.map_err(|e| device_error("Close app document", e))?;
    Ok(bytes)
}
fn check_helper_manifest(manifest: &Value, expected: &str) -> anyhow::Result<()> {
    anyhow::ensure!(
        manifest["job_id"] == expected,
        "Result belongs to a different operation"
    );
    if manifest["complete"] != true {
        let message: String = manifest["message"]
            .as_str()
            .unwrap_or("")
            .chars()
            .filter(|c| !c.is_control())
            .take(1024)
            .collect();
        let message = if message.trim().is_empty() {
            "Helper did not decrypt every encrypted image"
        } else {
            message.as_str()
        };
        anyhow::bail!(
            "Helper failed after {} decrypted images: {message}",
            manifest["decrypted"].as_u64().unwrap_or(0)
        );
    }
    Ok(())
}
fn local_pairing_addresses(route: &str, port: u16) -> anyhow::Result<Vec<SocketAddr>> {
    // Keep direct mode on this host. Never turn a route supplied by an IPA or
    // callback into arbitrary LAN access, or silently connect through a VPN.
    let hosts = match route {
        "10.7.0.1" => vec![IpAddr::V4(std::net::Ipv4Addr::new(10, 7, 0, 1))],
        "127.0.0.1" => vec![
            IpAddr::V4(std::net::Ipv4Addr::LOCALHOST),
            IpAddr::V6(std::net::Ipv6Addr::LOCALHOST),
        ],
        _ => anyhow::bail!("Unsupported local route"),
    };
    Ok(hosts
        .into_iter()
        .map(|ip| SocketAddr::new(ip, port))
        .collect())
}
async fn tunnel(v: &Value) -> anyhow::Result<(AdapterHandle, RsdHandshake)> {
    let addresses = pairing_candidates(v)?;
    let seconds = if addresses.len() > 2 { 8 } else { 15 };
    tunnel_on_candidates(v, &addresses, Duration::from_secs(seconds)).await
}
fn pairing_candidates(v: &Value) -> anyhow::Result<Vec<SocketAddr>> {
    let route = v["route"].as_str().unwrap_or("10.7.0.1");
    let ports = match v.get("pairing_ports") {
        None => vec![49152],
        Some(value) => {
            let values = value
                .as_array()
                .ok_or_else(|| anyhow::anyhow!("Invalid pairing ports"))?;
            anyhow::ensure!(
                !values.is_empty() && values.len() <= 2,
                "Invalid pairing port count"
            );
            values
                .iter()
                .map(|value| {
                    value
                        .as_u64()
                        .filter(|p| (1..=65535).contains(p))
                        .map(|p| p as u16)
                        .ok_or_else(|| anyhow::anyhow!("Invalid pairing port"))
                })
                .collect::<anyhow::Result<Vec<_>>>()?
        }
    };
    let mut addresses = Vec::new();
    for port in ports {
        for address in local_pairing_addresses(route, port)? {
            if !addresses.contains(&address) {
                addresses.push(address);
            }
        }
    }
    Ok(addresses)
}
async fn tunnel_on_candidates(
    v: &Value,
    addresses: &[SocketAddr],
    attempt_timeout: Duration,
) -> anyhow::Result<(AdapterHandle, RsdHandshake)> {
    let mut failures = Vec::new();
    for &address in addresses {
        // A successful TCP connect alone does not establish a route. iOS can
        // accept loopback and reset the protocol, so try the next family only
        // after attempting authentication, TLS and RSD on this one.
        let error = match tokio::time::timeout(attempt_timeout, tunnel_at(v, address)).await {
            Ok(Ok(connection)) => return Ok(connection),
            Ok(Err(error)) => error.to_string(),
            Err(_) => "Authenticated tunnel attempt timed out".to_owned(),
        };
        failures.push(format!("{address}: {error}"));
    }
    let advice = if !addresses.is_empty() && addresses.iter().all(|a| a.ip().is_loopback()) {
        " Direct loopback could not establish a device connection. Use LocalDevVPN and check again."
    } else {
        ""
    };
    anyhow::bail!("{}{advice}", failures.join("; "))
}
async fn tunnel_at(v: &Value, addr: SocketAddr) -> anyhow::Result<(AdapterHandle, RsdHandshake)> {
    let ip = addr.ip();
    let mut record = decode_record(v)?;
    let stream = tokio::time::timeout(
        Duration::from_secs(10),
        tokio::net::TcpStream::connect(addr),
    )
    .await
    .map_err(|_| anyhow::anyhow!("Pairing TCP: connection timed out after 10 seconds"))?
    .map_err(|e| device_error("Pairing TCP", e.into()))?;
    let mut client = RemotePairingClient::new(RpPairingSocket::new(stream), "Crumb");
    client
        .attempt_pair_verify()
        .await
        .map_err(|e| device_error("Pairing handshake", e))?;
    // Never silently fall back to pairing with a dummy PIN.
    client
        .validate_pairing(&mut record)
        .await
        .map_err(|e| device_error("Pairing verification", e))?;
    let port = client
        .create_tcp_listener()
        .await
        .map_err(|e| device_error("Tunnel listener", e))?;
    let stream = tokio::net::TcpStream::connect(SocketAddr::new(ip, port))
        .await
        .map_err(|e| device_error("Tunnel TCP", e.into()))?;
    let tunnel = connect_tls_psk_tunnel_native(stream, client.encryption_key())
        .await
        .map_err(|e| device_error("Tunnel TLS", e))?;
    let our_ip = tunnel.info.client_address.parse()?;
    let their_ip = tunnel.info.server_address.parse()?;
    let rsd_port = tunnel.info.server_rsd_port;
    anyhow::ensure!(
        tunnel.info.mtu > 60 && rsd_port != 0,
        "Invalid device tunnel parameters"
    );
    let mtu = tunnel.info.mtu as usize;
    let mut adapter = Adapter::new(Box::new(tunnel.into_inner()), our_ip, their_ip);
    adapter.set_mss(mtu.min(8192) - 60);
    let mut handle = adapter.to_async_handle();
    let stream = handle
        .connect(rsd_port)
        .await
        .map_err(|e| device_error("RSD TCP", e.into()))?;
    let handshake = RsdHandshake::new(stream)
        .await
        .map_err(|e| device_error("RSD handshake", e))?;
    Ok((handle, handshake))
}
enum Connection {
    Remote(AdapterHandle, RsdHandshake),
    Legacy(idevice::provider::TcpProvider),
}
impl Connection {
    async fn new(v: &Value) -> anyhow::Result<Self> {
        if v["pairing_kind"] == "lockdown" {
            let bytes = B64.decode(
                v["record"]
                    .as_str()
                    .ok_or_else(|| anyhow::anyhow!("Import pairing file"))?,
            )?;
            anyhow::ensure!(bytes.len() <= 65536, "Pairing file too large");
            let pairing_file = idevice::pairing_file::PairingFile::from_bytes(&bytes)?;
            Ok(Self::Legacy(idevice::provider::TcpProvider {
                addr: "10.7.0.1".parse()?,
                scope_id: None,
                pairing_file,
                label: "Crumb".into(),
            }))
        } else {
            let (a, h) = tunnel(v).await?;
            Ok(Self::Remote(a, h))
        }
    }
    async fn installation(&mut self) -> Result<InstallationProxyClient, idevice::IdeviceError> {
        use idevice::IdeviceService;
        match self {
            Self::Remote(a, h) => InstallationProxyClient::connect_rsd(a, h).await,
            Self::Legacy(p) => InstallationProxyClient::connect(p).await,
        }
    }
    async fn documents(&mut self) -> Result<HouseArrestClient, idevice::IdeviceError> {
        use idevice::IdeviceService;
        match self {
            Self::Remote(a, h) => HouseArrestClient::connect_rsd(a, h).await,
            Self::Legacy(p) => HouseArrestClient::connect(p).await,
        }
    }
    async fn afc(&mut self) -> Result<AfcClient, idevice::IdeviceError> {
        use idevice::IdeviceService;
        match self {
            Self::Remote(a, h) => AfcClient::connect_rsd(a, h).await,
            Self::Legacy(p) => AfcClient::connect(p).await,
        }
    }
}
pub(super) async fn run(job: &Arc<Job>, v: &Value) -> anyhow::Result<Value> {
    job.update(json!({"phase":"connecting"}));
    let mut connection =
        tokio::time::timeout(Duration::from_secs(30), Connection::new(v)).await??;
    let mut proxy = connection
        .installation()
        .await
        .map_err(|e| device_error("Installation service", e))?;
    let apps = proxy
        .browse(None)
        .await
        .map_err(|e| device_error("Installed-app lookup", e))?;
    let own_id = v["own_bundle"]
        .as_str()
        .ok_or_else(|| anyhow::anyhow!("Missing app identity"))?;
    anyhow::ensure!(
        apps.iter().any(|a| a
            .as_dictionary()
            .and_then(|d| d.get("CFBundleIdentifier"))
            .and_then(|v| v.as_string())
            == Some(own_id)),
        "Crumb is not installed on the connected device"
    );
    let proof = v["proof"]
        .as_str()
        .ok_or_else(|| anyhow::anyhow!("Missing local installation proof"))?;
    let own = connection
        .documents()
        .await
        .map_err(|e| device_error("Container service", e))?;
    let mut own = own
        .vend_documents(own_id)
        .await
        .map_err(|e| device_error("Crumb documents", e))?;
    let bytes = read_app_document(&mut own, AppDocument::Proof, 64).await?;
    anyhow::ensure!(
        bytes == proof.as_bytes(),
        "The connection does not lead to this Crumb installation"
    );
    drop(own);
    match v["operation"].as_str().unwrap_or("") {
        "verify" => Ok(json!({"phase":"connected"})),
        "install" => {
            let helper_id = v["helper_bundle"]
                .as_str()
                .ok_or_else(|| anyhow::anyhow!("Missing helper identity"))?;
            super::package::validate_helper_id(helper_id, own_id)?;
            let exists = apps.iter().any(|a| {
                a.as_dictionary()
                    .and_then(|d| d.get("CFBundleIdentifier"))
                    .and_then(|v| v.as_string())
                    .is_some_and(|id| id.eq_ignore_ascii_case(helper_id))
            });
            anyhow::ensure!(
                !exists || v["replace_helper"] == true,
                "This bundle identifier is already installed and is not a helper managed by Crumb"
            );
            let path = v["path"]
                .as_str()
                .ok_or_else(|| anyhow::anyhow!("Missing helper"))?;
            anyhow::ensure!(
                std::path::Path::new(path).is_dir(),
                "Install requires a prepared helper bundle"
            );
            job.update(json!({"phase":"installing"}));
            let mut afc = connection.afc().await?;
            // Only Crumb's own staging directory is removed. Interrupted uploads
            // are reclaimed on the next attempt after same-installation proof.
            let staging = "PublicStaging/Crumb-DevicePipeline";
            let _ = afc.mk_dir("PublicStaging").await;
            if afc.get_file_info(staging).await.is_ok() {
                afc.remove_all(staging).await?;
            }
            afc.mk_dir(staging).await?;
            let remote = format!("{staging}/CrumbDecrypt.app");
            let installation = async {
                upload_helper(job, &mut afc, std::path::Path::new(path), &remote).await?;
                let mut options = plist::Dictionary::new();
                options.insert("PackageType".into(), "Developer".into());
                proxy
                    .upgrade(&remote, Some(plist::Value::Dictionary(options)))
                    .await?;
                Ok::<_, anyhow::Error>(())
            }
            .await;
            let _ = afc.remove_all(staging).await;
            installation?;
            let id = v["helper_bundle"]
                .as_str()
                .ok_or_else(|| anyhow::anyhow!("Missing helper identity"))?;
            let apps = proxy
                .browse(None)
                .await
                .map_err(|e| device_error("Installed-app lookup", e))?;
            anyhow::ensure!(
                apps.iter().any(|a| a
                    .as_dictionary()
                    .and_then(|d| d.get("CFBundleIdentifier"))
                    .and_then(|v| v.as_string())
                    == Some(id)),
                "Helper not found after installation"
            );
            Ok(json!({"phase":"installed"}))
        }
        "fetch" => {
            let id = v["helper_bundle"]
                .as_str()
                .ok_or_else(|| anyhow::anyhow!("Missing helper identity"))?;
            let out = std::path::PathBuf::from(
                v["output"]
                    .as_str()
                    .ok_or_else(|| anyhow::anyhow!("Missing output"))?,
            );
            let client = connection.documents().await?;
            let mut afc = client.vend_documents(id).await?;
            let manifest_size = afc
                .get_file_info(AppDocument::Manifest.path())
                .await
                .map_err(|e| device_error("Helper manifest metadata", e))?
                .size;
            anyhow::ensure!(
                manifest_size > 0 && manifest_size <= 16 * 1024 * 1024,
                "Invalid result manifest size"
            );
            let manifest =
                read_app_document(&mut afc, AppDocument::Manifest, manifest_size).await?;
            anyhow::ensure!(manifest.len() == manifest_size, "Oversized result manifest");
            let manifest: Value = serde_json::from_slice(&manifest)?;
            let expected = v["job_id"]
                .as_str()
                .ok_or_else(|| anyhow::anyhow!("Missing job identity"))?;
            check_helper_manifest(&manifest, expected)?;
            let size = afc
                .get_file_info(AppDocument::Output.path())
                .await
                .map_err(|e| device_error("Helper IPA metadata", e))?
                .size;
            anyhow::ensure!(size > 0 && size <= 4_000_000_000, "Invalid result size");
            let partial = out.with_extension("partial");
            let mut cleanup = PartialFile(partial.clone());
            let mut remote = afc
                .open(AppDocument::Output.path(), AfcFopenMode::RdOnly)
                .await
                .map_err(|e| device_error("Helper IPA open", e))?;
            let mut local = tokio::fs::File::create(&partial).await?;
            let mut total = 0usize;
            while total < size {
                let chunk = remote.read_n((size - total).min(1024 * 1024)).await?;
                anyhow::ensure!(!chunk.is_empty(), "Incomplete IPA transfer");
                total += chunk.len();
                anyhow::ensure!(total <= size, "Result grew during transfer");
                local.write_all(&chunk).await?;
                job.update(json!({"phase":"fetching", "bytes":total, "total":size}));
            }
            remote
                .close()
                .await
                .map_err(|e| device_error("Helper IPA close", e))?;
            local.flush().await?;
            drop(local);
            super::package::verify_output(&partial, &manifest)?;
            let original =
                std::path::PathBuf::from(v["original_ipa"].as_str().ok_or_else(|| {
                    anyhow::anyhow!(
                        "Keep the original encrypted IPA to verify the decrypted export."
                    )
                })?);
            let candidate = partial.clone();
            let audit_job = job.clone();
            job.update(json!({"phase":"verifying-export"}));
            tokio::task::spawn_blocking(move || {
                super::export_audit::verify(&original, &candidate, &audit_job)?;
                audit_job.update(json!({"phase":"preparing-export"}));
                super::export_compatibility::normalize(&candidate, &audit_job)
            })
            .await??;
            anyhow::ensure!(!job.cancelled.load(Ordering::SeqCst), "Cancelled");
            tokio::fs::rename(&partial, &out).await?;
            cleanup.0 = std::path::PathBuf::new();
            Ok(json!({"phase":"complete", "path":out}))
        }
        _ => anyhow::bail!("Unknown device operation"),
    }
}
async fn upload_helper(
    job: &Arc<Job>,
    afc: &mut AfcClient,
    local: &std::path::Path,
    remote: &str,
) -> anyhow::Result<()> {
    let mut buffer = vec![0u8; 1024 * 1024];
    let mut bytes = 0u64;
    for entry in walkdir::WalkDir::new(local) {
        anyhow::ensure!(!job.cancelled.load(Ordering::SeqCst), "Cancelled");
        let entry = entry?;
        anyhow::ensure!(
            !entry.file_type().is_symlink(),
            "Prepared helper contains a link"
        );
        let relative = entry.path().strip_prefix(local)?;
        let target = if relative.as_os_str().is_empty() {
            remote.to_string()
        } else {
            format!(
                "{remote}/{}",
                relative
                    .to_str()
                    .ok_or_else(|| anyhow::anyhow!("Invalid helper path"))?
            )
        };
        if entry.file_type().is_dir() {
            afc.mk_dir(&target).await?;
            continue;
        }
        anyhow::ensure!(entry.file_type().is_file(), "Unsupported helper entry");
        let mut input = tokio::fs::File::open(entry.path()).await?;
        let mut output = afc.open(&target, AfcFopenMode::WrOnly).await?;
        loop {
            anyhow::ensure!(!job.cancelled.load(Ordering::SeqCst), "Cancelled");
            let n = input.read(&mut buffer).await?;
            if n == 0 {
                break;
            }
            output.write_entire(&buffer[..n]).await?;
            bytes += n as u64;
            job.update(json!({"phase":"installing","bytes":bytes}));
        }
        output.close().await?;
    }
    Ok(())
}
struct PartialFile(std::path::PathBuf);
impl Drop for PartialFile {
    fn drop(&mut self) {
        if !self.0.as_os_str().is_empty() {
            let _ = std::fs::remove_file(&self.0);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn failed_helper_result_preserves_its_file_error_without_accepting_output() {
        let manifest = json!({"job_id":"test-job", "complete":false, "decrypted":2,
            "message":"Frameworks/Test.framework/Test: Cannot open installed image (errno 13)"});
        let message = check_helper_manifest(&manifest, "test-job")
            .unwrap_err()
            .to_string();
        assert!(message.contains("after 2 decrypted images"));
        assert!(message.contains("Frameworks/Test.framework/Test"));
        assert!(message.contains("errno 13"));
    }
    #[test]
    fn wrong_job_result_does_not_expose_its_diagnostic() {
        let manifest =
            json!({"job_id":"stale-job", "complete":false, "message":"unrelated private payload"});
        let message = check_helper_manifest(&manifest, "test-job")
            .unwrap_err()
            .to_string();
        assert!(message.contains("different operation"));
        assert!(!message.contains("private payload"));
    }
    #[test]
    fn helper_status_checks_completion_and_bounds_failure_text() {
        assert!(
            check_helper_manifest(&json!({"job_id":"test-job", "complete":true}), "test-job")
                .is_ok()
        );
        let message = check_helper_manifest(
            &json!({"job_id":"test-job", "complete":false,
            "message":format!("\u{0000}{}", "x".repeat(4096))}),
            "test-job",
        )
        .unwrap_err()
        .to_string();
        assert!(!message.contains('\0'));
        assert!(message.len() < 1200);
        let message =
            check_helper_manifest(&json!({"job_id":"test-job", "complete":false}), "test-job")
                .unwrap_err()
                .to_string();
        assert!(message.contains("did not decrypt every encrypted image"));
    }
    #[test]
    fn socket_diagnostic_preserves_stage_and_os_code() {
        let error = device_error(
            "Pairing handshake",
            std::io::Error::from_raw_os_error(54).into(),
        )
        .to_string();
        assert!(error.starts_with("Pairing handshake: socket "));
        assert!(error.contains("OS 54"));
    }
    #[test]
    fn protocol_payloads_are_excluded_from_diagnostics() {
        let message = device_error(
            "RSD handshake",
            idevice::IdeviceError::UnexpectedResponse(
                "private payload that must not be shown".into(),
            ),
        )
        .to_string();
        assert!(message.starts_with("RSD handshake: device protocol error "));
        assert!(!message.contains("private payload"));
    }
    #[test]
    fn socket_text_payloads_are_excluded_from_diagnostics() {
        let message = device_error(
            "Device proof read",
            std::io::Error::new(std::io::ErrorKind::UnexpectedEof, "private socket payload").into(),
        )
        .to_string();
        assert!(message.contains("UnexpectedEof"));
        assert!(!message.contains("private socket payload"));
    }
}

#[cfg(test)]
#[path = "document_tests.rs"]
mod document_tests;

#[cfg(test)]
#[path = "route_tests.rs"]
mod route_tests;
