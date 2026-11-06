//! AFC wire fixture whose root models the app container returned by VendDocuments.
//! This tests real AfcClient open/read/close messages, not a replacement client.
use super::*;
use std::collections::HashMap;
use tokio::io::{AsyncReadExt, AsyncWriteExt, DuplexStream};

async fn response(
    stream: &mut DuplexStream,
    number: u64,
    opcode: u64,
    header: &[u8],
    payload: &[u8],
) {
    let mut packet = Vec::new();
    for value in [
        0x4141504c36414643,
        40 + (header.len() + payload.len()) as u64,
        40 + header.len() as u64,
        number,
        opcode,
    ] {
        packet.extend(value.to_le_bytes());
    }
    packet.extend(header);
    packet.extend(payload);
    stream.write_all(&packet).await.unwrap();
}
async fn container_server(mut stream: DuplexStream) -> (Vec<String>, usize) {
    let files = HashMap::from([
        (
            "/Documents/.crumb-device-proof",
            b"fresh-installation-nonce".to_vec(),
        ),
        (
            "/Documents/result.json",
            br#"{"complete":true,"job_id":"fixture"}"#.to_vec(),
        ),
        (
            "/Documents/result.ipa",
            b"PK\x03\x04fixture-IPA-bytes".to_vec(),
        ),
    ]);
    let mut handles: HashMap<u64, (Vec<u8>, usize)> = HashMap::new();
    let mut requests = vec![];
    let mut next = 1;
    let mut closed = 0;
    loop {
        let mut header = [0u8; 40];
        if stream.read_exact(&mut header).await.is_err() {
            break;
        }
        let word = |i| u64::from_le_bytes(header[i..i + 8].try_into().unwrap());
        assert_eq!(word(0), 0x4141504c36414643);
        let total = word(8) as usize;
        assert!((40..=4096).contains(&total));
        let number = word(24);
        let op = word(32);
        let mut body = vec![0u8; total - 40];
        stream.read_exact(&mut body).await.unwrap();
        match op {
            13 => {
                assert!(body.len() >= 8);
                let path = String::from_utf8(body[8..].to_vec())
                    .unwrap()
                    .trim_end_matches('\0')
                    .to_string();
                requests.push(path.clone());
                if let Some(bytes) = files.get(path.as_str()) {
                    let fd = next;
                    next += 1;
                    handles.insert(fd, (bytes.clone(), 0));
                    response(&mut stream, number, 14, &fd.to_le_bytes(), &[]).await;
                } else {
                    response(&mut stream, number, 1, &8u64.to_le_bytes(), &[]).await;
                }
            }
            15 => {
                let fd = u64::from_le_bytes(body[..8].try_into().unwrap());
                let n = u64::from_le_bytes(body[8..16].try_into().unwrap()) as usize;
                let (bytes, pos) = handles
                    .get_mut(&fd)
                    .expect("read requires a live file descriptor");
                let end = (*pos + n).min(bytes.len());
                let data = bytes[*pos..end].to_vec();
                *pos = end;
                response(&mut stream, number, 2, &[], &data).await;
            }
            20 => {
                let fd = u64::from_le_bytes(body[..8].try_into().unwrap());
                assert!(handles.remove(&fd).is_some());
                closed += 1;
                response(&mut stream, number, 1, &0u64.to_le_bytes(), &[]).await;
            }
            _ => panic!("unexpected AFC operation {op}"),
        }
    }
    assert!(
        handles.is_empty(),
        "every opened app document must be closed"
    );
    (requests, closed)
}
#[tokio::test]
async fn documents_access_uses_container_relative_paths_and_closes_descriptors() {
    let (client, server) = tokio::io::duplex(65536);
    let server = tokio::spawn(container_server(server));
    let mut client = AfcClient::new(idevice::Idevice::new(Box::new(client), "Crumb test"));
    // Reproduce the old failure: the file exists under Documents, not at root.
    let old = client
        .open(".crumb-device-proof", AfcFopenMode::RdOnly)
        .await;
    assert!(matches!(
        old,
        Err(idevice::IdeviceError::Afc(
            idevice::services::afc::errors::AfcError::ObjectNotFound
        ))
    ));
    drop(old);
    let proof = read_app_document(&mut client, AppDocument::Proof, 64)
        .await
        .unwrap();
    assert_eq!(proof, b"fresh-installation-nonce");
    let manifest = read_app_document(&mut client, AppDocument::Manifest, 256)
        .await
        .unwrap();
    assert_eq!(
        serde_json::from_slice::<Value>(&manifest).unwrap()["complete"],
        true
    );
    // The IPA retrieval uses the same document selector for streaming reads.
    let ipa = read_app_document(&mut client, AppDocument::Output, 64)
        .await
        .unwrap();
    assert!(ipa.starts_with(b"PK\x03\x04"));
    drop(client);
    let (paths, closed) = tokio::time::timeout(Duration::from_secs(2), server)
        .await
        .unwrap()
        .unwrap();
    assert_eq!(closed, 3);
    assert_eq!(
        paths,
        vec![
            ".crumb-device-proof",
            "/Documents/.crumb-device-proof",
            "/Documents/result.json",
            "/Documents/result.ipa"
        ]
    );
}
