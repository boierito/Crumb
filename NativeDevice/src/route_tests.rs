//! These peers exercise real local TCP and the real RemotePairing client.
//! They reject or stall its handshake; they do not emulate iOS acceptance.
use super::*;

#[test]
fn direct_routes_cover_both_families_without_leaving_this_host() {
    let addresses = local_pairing_addresses("127.0.0.1", 49152).unwrap();
    assert_eq!(addresses.len(), 2);
    assert!(addresses[0].is_ipv4());
    assert!(addresses[1].is_ipv6());
    assert!(
        addresses
            .iter()
            .all(|a| a.ip().is_loopback() && a.port() == 49152)
    );
    assert_eq!(local_pairing_addresses("10.7.0.1", 49152).unwrap().len(), 1);
    for route in ["192.168.1.1", "example.com", "::", "127.0.0.2"] {
        assert!(local_pairing_addresses(route, 49152).is_err());
    }
}
#[test]
fn advertised_port_precedes_default_without_leaving_loopback() {
    let candidates =
        pairing_candidates(&json!({"route":"127.0.0.1", "pairing_ports":[49170,49152]})).unwrap();
    assert_eq!(candidates.len(), 4);
    assert!(candidates.iter().all(|a| a.ip().is_loopback()));
    assert_eq!(
        candidates.iter().map(|a| a.port()).collect::<Vec<_>>(),
        vec![49170, 49170, 49152, 49152]
    );
    let default = pairing_candidates(&json!({"route":"127.0.0.1"})).unwrap();
    assert_eq!(default.len(), 2);
    assert!(default.iter().all(|a| a.port() == 49152));
    assert_eq!(
        pairing_candidates(&json!({"route":"127.0.0.1", "pairing_ports":[49152,49152]}))
            .unwrap()
            .len(),
        2
    );
}
#[test]
fn invalid_and_excessive_discovered_ports_are_rejected() {
    for ports in [
        json!([]),
        json!([0]),
        json!([65536]),
        json!([-1]),
        json!(["49152"]),
        json!([49152, 49153, 49154]),
        json!("49152"),
    ] {
        assert!(pairing_candidates(&json!({"route":"127.0.0.1", "pairing_ports":ports})).is_err());
    }
    assert!(pairing_candidates(&json!({"route":"192.168.1.1", "pairing_ports":[49170]})).is_err());
}
fn request() -> Value {
    json!({"record": B64.encode(RpPairingFile::generate("Loopback test").to_bytes())})
}
async fn reject_after_handshake(listener: TcpListener) {
    let (mut socket, _) = listener.accept().await.unwrap();
    let mut byte = [0];
    // A bare successful connect must not skip the other address family.
    socket.read_exact(&mut byte).await.unwrap();
}
#[tokio::test]
async fn accepted_tcp_then_failed_pairing_still_tries_ipv6() {
    let ipv4 = TcpListener::bind(("127.0.0.1", 0)).await.unwrap();
    let ipv6 = TcpListener::bind(("::1", 0)).await.unwrap();
    let addresses = [ipv4.local_addr().unwrap(), ipv6.local_addr().unwrap()];
    let first = tokio::spawn(reject_after_handshake(ipv4));
    let second = tokio::spawn(reject_after_handshake(ipv6));
    let result = tokio::time::timeout(
        Duration::from_secs(3),
        tunnel_on_candidates(&request(), &addresses, Duration::from_secs(1)),
    )
    .await
    .unwrap();
    let message = match result {
        Err(error) => error.to_string(),
        Ok(_) => panic!("TCP acceptance must not be reported as an authenticated tunnel"),
    };
    assert_eq!(message.matches("Pairing handshake:").count(), 2);
    assert!(message.contains("127.0.0.1:"));
    assert!(message.contains("[::1]:"));
    assert!(message.contains("Use LocalDevVPN"));
    first.await.unwrap();
    second.await.unwrap();
}
#[tokio::test]
async fn stalled_pairing_times_out_then_tries_next_family() {
    let ipv4 = TcpListener::bind(("127.0.0.1", 0)).await.unwrap();
    let ipv6 = TcpListener::bind(("::1", 0)).await.unwrap();
    let addresses = [ipv4.local_addr().unwrap(), ipv6.local_addr().unwrap()];
    let first = tokio::spawn(async move {
        let (mut socket, _) = ipv4.accept().await.unwrap();
        // Drain the request without replying; this finishes when the bounded
        // attempt is cancelled and its socket is closed.
        let mut buffer = [0; 4096];
        while socket.read(&mut buffer).await.unwrap_or(0) != 0 {}
    });
    let second = tokio::spawn(reject_after_handshake(ipv6));
    let result = tokio::time::timeout(
        Duration::from_secs(3),
        tunnel_on_candidates(&request(), &addresses, Duration::from_millis(250)),
    )
    .await
    .unwrap();
    let message = match result {
        Err(error) => error.to_string(),
        Ok(_) => panic!("stalled pairing must not establish a tunnel"),
    };
    assert!(message.contains("Authenticated tunnel attempt timed out"));
    assert!(message.contains("[::1]:"));
    assert!(message.contains("Pairing handshake:"));
    first.await.unwrap();
    second.await.unwrap();
}

#[tokio::test]
async fn download_metadata_cannot_be_used_as_pairing_credentials_or_reach_tcp() {
    let listener = TcpListener::bind(("127.0.0.1", 0)).await.unwrap();
    let bad = json!({"record":B64.encode(br#"{"filename":"download.ipa"}"#)});
    let result = tunnel_on_candidates(
        &bad,
        &[listener.local_addr().unwrap()],
        Duration::from_secs(1),
    )
    .await;
    let error = match result {
        Err(e) => e.to_string(),
        Ok(_) => panic!("Metadata is not a pairing record"),
    };
    assert!(error.contains("Pairing credentials are unreadable"));
    assert!(!error.contains("download.ipa"));
    assert!(
        tokio::time::timeout(Duration::from_millis(50), listener.accept())
            .await
            .is_err()
    );
}
