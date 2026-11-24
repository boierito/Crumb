mod entitlement_der;
mod package;
mod resign;
mod signing;
use serde_json::{Value, json};
use std::{
    ffi::{CStr, CString, c_char},
    sync::{
        Arc, Mutex,
        atomic::{AtomicBool, Ordering},
    },
    time::Duration,
};
use tokio::sync::Notify;

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
async fn run(job: &Arc<Job>, v: Value) -> anyhow::Result<Value> {
    match v["operation"].as_str().unwrap_or("") {
        "signing-info" | "signing-plan" | "sign" => {
            job.update(json!({"phase":"signing"}));
            let job = job.clone();
            tokio::task::spawn_blocking(move || {
                if v["operation"] == "sign" {
                    resign::sign(&v, &job)
                } else if v["operation"] == "signing-plan" {
                    resign::preview(&v, &job)
                } else {
                    signing::inspect(&v)
                }
            })
            .await?
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
                if worker.cancelled.load(Ordering::SeqCst) {
                    anyhow::bail!("Cancelled");
                }
                let blocking = matches!(
                    v["operation"].as_str(),
                    Some("sign" | "signing-info" | "signing-plan")
                );
                let operation = run(&worker, v);
                tokio::pin!(operation);
                tokio::select! {
                    _ = worker.notify.notified() => {
                        // Blocking signers must release their workspace before Swift removes it.
                        if blocking { let _ = operation.await; }
                        anyhow::bail!("Cancelled")
                    },
                    _ = tokio::time::sleep(Duration::from_secs(600)) => {
                        worker.cancelled.store(true, Ordering::SeqCst);
                        if blocking { let _ = operation.await; }
                        anyhow::bail!("Operation timed out")
                    },
                    result = &mut operation => result,
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
