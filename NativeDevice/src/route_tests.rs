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
