mod package;
mod transport;
use base64::{Engine, engine::general_purpose::STANDARD as B64};
use idevice::remote_pairing::{PairableHost, PairableHostInfo, RpPairingFile, RpPairingSocket};
use serde_json::{Value, json};
use std::{
    ffi::{CStr, CString, c_char},
    sync::{
        Arc, Mutex,
        atomic::{AtomicBool, Ordering},
    },
    time::Duration,
};
use tokio::{net::TcpListener, sync::Notify};

pub struct Job {
    state: Mutex<Value>,
    cancelled: AtomicBool,
    notify: Notify,
}
impl Job {
    fn update(&self, value: Value) {
        *self.state.lock().unwrap_or_else(|e| e.into_inner()) = value;
    }
    fn cancel(&self) {
        self.cancelled.store(true, Ordering::SeqCst);
        self.notify.notify_one();
    }
}
fn decode_record(v: &Value) -> anyhow::Result<RpPairingFile> {
    let bytes = B64.decode(
        v["record"]
            .as_str()
            .ok_or_else(|| anyhow::anyhow!("Pair this device first"))?,
    )?;
    anyhow::ensure!(bytes.len() <= 65536, "Pairing record too large");
    Ok(RpPairingFile::from_bytes(&bytes)?)
}
async fn pair(job: &Arc<Job>, v: &Value) -> anyhow::Result<Value> {
    let listener = TcpListener::bind((std::net::Ipv4Addr::UNSPECIFIED, 0)).await?;
    let mut file = if v["record"].is_string() {
        decode_record(v)?
    } else {
        let mut f = RpPairingFile::generate("Crumb");
        f.identifier = uuid::Uuid::new_v4().to_string();
        f
    };
    let mut info = PairableHostInfo::generate("Crumb", "Mac17,7");
    if let Some(s) = v["host_irk"].as_str() {
        let b = B64.decode(s)?;
        info.alt_irk = b
            .try_into()
            .map_err(|_| anyhow::anyhow!("Invalid host identity"))?;
    }
    let irk = B64.encode(info.alt_irk);
    let ready = json!({"phase":"waiting", "port":listener.local_addr()?.port(), "service":file.identifier,
        "txt":info.mdns_txt_records(&file.identifier).into_iter().collect::<std::collections::BTreeMap<_,_>>()});
    job.update(ready.clone());
    // Ignore unrelated probes only before issuing a PIN; after that errors are actionable.
    loop {
        let (stream, _) = listener.accept().await?;
        let issued = Arc::new(AtomicBool::new(false));
        let seen = issued.clone();
        let j = job.clone();
        let status = ready.clone();
        let mut candidate = file.clone();
        let mut host = PairableHost::new(RpPairingSocket::new_device(stream), info.clone());
        match tokio::time::timeout(
            Duration::from_secs(120),
            host.accept(&mut candidate, move |pin| {
                seen.store(true, Ordering::SeqCst);
                let mut s = status.clone();
                s["phase"] = json!("pin");
                s["pin"] = json!(pin);
                j.update(s);
                async {}
            }),
        )
        .await
        {
            Ok(Ok(_peer)) => {
                file = candidate;
                break;
            }
            error if issued.load(Ordering::SeqCst) => anyhow::bail!(
                "Pairing confirmation failed: {}",
                if error.is_err() {
                    "timeout"
                } else {
                    "protocol rejected"
                }
            ),
            _ => continue,
        }
    }
    Ok(json!({"phase":"paired", "record": B64.encode(file.to_bytes()), "host_irk":irk}))
}
async fn run(job: &Arc<Job>, v: Value) -> anyhow::Result<Value> {
    match v["operation"].as_str().unwrap_or("") {
        "pair" => pair(job, &v).await,
        "verify" | "install" | "fetch" => transport::run(job, &v).await,
        "prepare" => {
            job.update(json!({"phase":"preparing"}));
            let v = v.clone();
            let path = tokio::task::spawn_blocking(move || package::prepare(&v)).await??;
            Ok(json!({"phase":"prepared", "path":path}))
        }
        _ => anyhow::bail!("Unknown operation"),
    }
}
// Safety: JSON must be a valid NUL-terminated UTF-8 string, bounded by the caller.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn crumb_device_start(ptr: *const c_char) -> *mut Arc<Job> {
    if ptr.is_null() {
        return std::ptr::null_mut();
    }
    let Ok(s) = (unsafe { CStr::from_ptr(ptr) }).to_str() else {
        return std::ptr::null_mut();
    };
    if s.len() > 2 * 1024 * 1024 {
        return std::ptr::null_mut();
    }
    let Ok(v) = serde_json::from_str::<Value>(s) else {
        return std::ptr::null_mut();
    };
    let job = Arc::new(Job {
        state: Mutex::new(json!({"phase":"starting"})),
        cancelled: AtomicBool::new(false),
        notify: Notify::new(),
    });
    let worker = job.clone();
    std::thread::spawn(move || {
        let outcome = std::panic::catch_unwind(std::panic::AssertUnwindSafe(|| {
            let rt = tokio::runtime::Builder::new_multi_thread()
                .worker_threads(2)
                .enable_all()
                .build()?;
            rt.block_on(async {
                if worker.cancelled.load(Ordering::SeqCst) { anyhow::bail!("Cancelled"); }
                tokio::select! {
                    _ = worker.notify.notified() => anyhow::bail!("Cancelled"),
                    result = tokio::time::timeout(Duration::from_secs(600), run(&worker, v)) => result?,
                }
            })
        }));
        let state = match outcome {
            Ok(Ok(s)) if !worker.cancelled.load(Ordering::SeqCst) => s,
            _ if worker.cancelled.load(Ordering::SeqCst) => json!({"phase":"cancelled"}),
            Ok(Err(e)) => json!({"phase":"failed", "message":e.to_string()}),
            _ => json!({"phase":"failed", "message":"Native operation failed"}),
        };
        worker.update(state);
    });
    Box::into_raw(Box::new(job))
}
#[unsafe(no_mangle)]
pub unsafe extern "C" fn crumb_device_poll(ptr: *mut Arc<Job>) -> *mut c_char {
    if ptr.is_null() {
        return std::ptr::null_mut();
    }
    let state = unsafe { &*ptr }
        .state
        .lock()
        .unwrap_or_else(|e| e.into_inner())
        .to_string();
    CString::new(state).unwrap().into_raw()
}
#[unsafe(no_mangle)]
pub unsafe extern "C" fn crumb_device_cancel(ptr: *mut Arc<Job>) {
    if !ptr.is_null() {
        unsafe { &*ptr }.cancel();
    }
}
#[unsafe(no_mangle)]
pub unsafe extern "C" fn crumb_device_free(ptr: *mut Arc<Job>) {
    if !ptr.is_null() {
        let job = unsafe { Box::from_raw(ptr) };
        job.cancel();
    }
}
#[unsafe(no_mangle)]
pub unsafe extern "C" fn crumb_device_string_free(ptr: *mut c_char) {
    if !ptr.is_null() {
        drop(unsafe { CString::from_raw(ptr) });
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn cancellation_closes_listener_and_advertising_has_swift_shape() {
        let input = CString::new(r#"{"operation":"pair"}"#).unwrap();
        let ptr = unsafe { crumb_device_start(input.as_ptr()) };
        let mut port = None;
        for _ in 0..200 {
            let p = unsafe { crumb_device_poll(ptr) };
            let value: Value =
                serde_json::from_slice(unsafe { CStr::from_ptr(p) }.to_bytes()).unwrap();
            unsafe { crumb_device_string_free(p) };
            if value["phase"] == "waiting" {
                assert_eq!(value["txt"]["name"], "Crumb");
                assert!(
                    value["txt"].as_object().is_some(),
                    "Swift expects a TXT dictionary"
                );
                port = Some(value["port"].as_u64().unwrap() as u16);
                break;
            }
            std::thread::sleep(Duration::from_millis(10));
        }
        let port = port.expect("Pairing listener did not start");
        unsafe { crumb_device_cancel(ptr) };
        let mut cancelled = false;
        for _ in 0..100 {
            let p = unsafe { crumb_device_poll(ptr) };
            cancelled = unsafe { CStr::from_ptr(p) }
                .to_string_lossy()
                .contains("cancelled");
            unsafe { crumb_device_string_free(p) };
            if cancelled {
                break;
            }
            std::thread::sleep(Duration::from_millis(10));
        }
        assert!(cancelled);
        assert!(std::net::TcpStream::connect((std::net::Ipv4Addr::LOCALHOST, port)).is_err());
        unsafe { crumb_device_free(ptr) };
    }
    #[test]
    fn malformed_requests_are_rejected() {
        let s = CString::new("{").unwrap();
        assert!(unsafe { crumb_device_start(s.as_ptr()) }.is_null());
    }
}
