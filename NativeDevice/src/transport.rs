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

async fn tunnel(v: &Value) -> anyhow::Result<(AdapterHandle, RsdHandshake)> {
    // Only local device routes. Never accept an arbitrary network host from an IPA.
    let ip: IpAddr = match v["route"].as_str().unwrap_or("10.7.0.1") {
        "10.7.0.1" => "10.7.0.1".parse()?,
        "127.0.0.1" => "127.0.0.1".parse()?,
        _ => anyhow::bail!("Unsupported local route"),
    };
    let addr = SocketAddr::new(ip, 49152);
    let mut record = decode_record(v)?;
    let stream = tokio::time::timeout(
        Duration::from_secs(10),
        tokio::net::TcpStream::connect(addr),
    )
    .await??;
    let mut client = RemotePairingClient::new(RpPairingSocket::new(stream), "Crumb");
    client.attempt_pair_verify().await?;
    // Never silently fall back to pairing with a dummy PIN.
    client.validate_pairing(&mut record).await?;
    let port = client.create_tcp_listener().await?;
    let stream = tokio::net::TcpStream::connect(SocketAddr::new(ip, port)).await?;
    let tunnel = connect_tls_psk_tunnel_native(stream, client.encryption_key()).await?;
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
    let handshake = RsdHandshake::new(handle.connect(rsd_port).await?).await?;
    Ok((handle, handshake))
}
pub(super) async fn run(job: &Arc<Job>, v: &Value) -> anyhow::Result<Value> {
    job.update(json!({"phase":"connecting"}));
    let (mut adapter, mut handshake) = tunnel(v).await?;
    let mut proxy = InstallationProxyClient::connect_rsd(&mut adapter, &mut handshake).await?;
    let apps = proxy.browse(None).await?;
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
    let own = HouseArrestClient::connect_rsd(&mut adapter, &mut handshake).await?;
    let mut own = own.vend_documents(own_id).await?;
    let mut marker = own
        .open(".crumb-device-proof", AfcFopenMode::RdOnly)
        .await?;
    let bytes = marker.read_n(64).await?;
    anyhow::ensure!(
        bytes == proof.as_bytes(),
        "The connection does not lead to this Crumb installation"
    );
    drop(marker);
    drop(own);
    match v["operation"].as_str().unwrap_or("") {
        "verify" => Ok(json!({"phase":"connected"})),
        "install" => {
            let path = v["path"]
                .as_str()
                .ok_or_else(|| anyhow::anyhow!("Missing helper"))?;
            anyhow::ensure!(
                std::path::Path::new(path).is_dir(),
                "Install requires a prepared helper bundle"
            );
            job.update(json!({"phase":"installing"}));
            let mut afc = AfcClient::connect_rsd(&mut adapter, &mut handshake).await?;
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
            let apps = proxy.browse(None).await?;
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
            let client = HouseArrestClient::connect_rsd(&mut adapter, &mut handshake).await?;
            let mut afc = client.vend_documents(id).await?;
            let manifest_size = afc.get_file_info("result.json").await?.size;
            anyhow::ensure!(
                manifest_size > 0 && manifest_size <= 16 * 1024 * 1024,
                "Invalid result manifest size"
            );
            let manifest = {
                let mut file = afc.open("result.json", AfcFopenMode::RdOnly).await?;
                file.read_n(manifest_size).await?
            };
            anyhow::ensure!(manifest.len() == manifest_size, "Oversized result manifest");
            let manifest: Value = serde_json::from_slice(&manifest)?;
            anyhow::ensure!(
                manifest["complete"] == true,
                "Helper did not decrypt every encrypted image"
            );
            let expected = v["job_id"]
                .as_str()
                .ok_or_else(|| anyhow::anyhow!("Missing job identity"))?;
            anyhow::ensure!(
                manifest["job_id"] == expected,
                "Result belongs to a different operation"
            );
            let size = afc.get_file_info("result.ipa").await?.size;
            anyhow::ensure!(size > 0 && size <= 4_000_000_000, "Invalid result size");
            let partial = out.with_extension("partial");
            let mut cleanup = PartialFile(partial.clone());
            let mut remote = afc.open("result.ipa", AfcFopenMode::RdOnly).await?;
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
            local.flush().await?;
            drop(local);
            super::package::verify_output(&partial, &manifest)?;
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
