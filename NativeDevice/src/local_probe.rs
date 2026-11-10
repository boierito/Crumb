//! Archive-only half of the in-process experiment. No transport, account or signing API.
use crate::package;
use anyhow::{Context, Result, ensure};
use serde_json::{Value, json};
use sha2::{Digest, Sha256};
use std::{
    collections::BTreeMap,
    fs,
    io::{Read, Write},
    path::Path,
};
use zip::write::SimpleFileOptions;

type Cancel = Option<unsafe extern "C" fn(*mut std::ffi::c_void) -> bool>;
struct Control {
    context: *mut std::ffi::c_void,
    cancel: Cancel,
}
impl Control {
    fn check(&self) -> Result<()> {
        ensure!(
            !self.cancel.is_some_and(|f| unsafe { f(self.context) }),
            "cancelled"
        );
        Ok(())
    }
}
fn path<'a>(v: &'a Value, key: &str) -> Result<&'a Path> {
    Ok(Path::new(v[key].as_str().context("missing path")?))
}
fn relative(s: &str) -> Result<&Path> {
    ensure!(
        !s.is_empty()
            && !s.contains('\\')
            && !s.contains('\0')
            && s.split('/').all(|p| !p.is_empty() && p != "." && p != ".."),
        "unsafe relative path"
    );
    Ok(Path::new(s))
}
fn hash(path: &Path, c: &Control) -> Result<String> {
    let mut file = fs::File::open(path)?;
    let mut h = Sha256::new();
    let mut buf = vec![0; 65536];
    loop {
        c.check()?;
        let n = file.read(&mut buf)?;
        if n == 0 {
            break;
        }
        h.update(&buf[..n]);
    }
    Ok(format!("{:x}", h.finalize()))
}
fn inventory(root: &Path, c: &Control) -> Result<BTreeMap<String, Value>> {
    let mut out = BTreeMap::new();
    let mut total = 0u64;
    for entry in walkdir::WalkDir::new(root).min_depth(1) {
        c.check()?;
        let entry = entry?;
        ensure!(!entry.file_type().is_symlink(), "links unsupported");
        let name = entry
            .path()
            .strip_prefix(root)?
            .to_str()
            .context("non-UTF8 path")?
            .to_owned();
        relative(&name)?;
        let st = entry.metadata()?;
        if entry.file_type().is_dir() {
            out.insert(name, json!({"directory":true}));
            continue;
        }
        ensure!(
            entry.file_type().is_file() && st.len() <= 512 * 1024 * 1024,
            "unsupported file"
        );
        total = total.checked_add(st.len()).context("size overflow")?;
        ensure!(
            total <= 4_000_000_000 && out.len() < 100000,
            "archive budget"
        );
        #[cfg(unix)]
        let mode = {
            use std::os::unix::fs::PermissionsExt;
            st.permissions().mode() & 0o777
        };
        #[cfg(not(unix))]
        let mode = 0o644;
        out.insert(
            name,
            json!({"size":st.len(),"sha256":hash(entry.path(),c)?,"mode":mode}),
        );
    }
    Ok(out)
}
fn info(path: &Path) -> Result<plist::Dictionary> {
    ensure!(fs::metadata(path)?.len() <= 1024 * 1024, "metadata budget");
    plist::Value::from_file(path)?
        .into_dictionary()
        .context("invalid app metadata")
}
fn text<'a>(d: &'a plist::Dictionary, k: &str) -> Result<&'a str> {
    d.get(k)
        .and_then(|x| x.as_string())
        .context("missing app identity")
}
fn copy_tree(from: &Path, to: &Path, c: &Control) -> Result<()> {
    for e in walkdir::WalkDir::new(from) {
        c.check()?;
        let e = e?;
        ensure!(!e.file_type().is_symlink(), "links unsupported");
        let dest = to.join(e.path().strip_prefix(from)?);
        if e.file_type().is_dir() {
            fs::create_dir_all(dest)?;
        } else {
            let mut input = fs::File::open(e.path())?;
            let mut output = fs::OpenOptions::new()
                .write(true)
                .create_new(true)
                .open(&dest)?;
            fs::set_permissions(&dest, e.metadata()?.permissions())?;
            let mut buffer = [0u8; 65536];
            loop {
                c.check()?;
                let n = input.read(&mut buffer)?;
                if n == 0 {
                    break;
                }
                output.write_all(&buffer[..n])?;
            }
        }
    }
    Ok(())
}
fn replicate(app: &Path) -> Result<()> {
    let manifest = app.join("SC_Info/Manifest.plist");
    if !manifest.exists() {
        return Ok(());
    }
    let d = info(&manifest)?;
    let Some(paths) = d.get("SinfReplicationPaths").and_then(|x| x.as_array()) else {
        return Ok(());
    };
    if paths.is_empty() {
        return Ok(());
    }
    let sources = d
        .get("SinfPaths")
        .and_then(|x| x.as_array())
        .context("missing original SINF mapping")?;
    ensure!(
        sources.len() <= 4096 && paths.len() <= 4096,
        "license resource budget"
    );
    let mut license_bytes = 0u64;
    let mut sinfs = Vec::new();
    for p in sources {
        let name = p.as_string().context("invalid SINF mapping")?;
        let p = relative(name)?;
        ensure!(
            name.split('/').any(|s| s == "SC_Info"),
            "SINF outside license directory"
        );
        ensure!(
            fs::metadata(app.join(p))?.len() <= 4 * 1024 * 1024,
            "SINF budget"
        );
        license_bytes += fs::metadata(app.join(p))?.len();
        ensure!(license_bytes <= 16 * 1024 * 1024, "license resource budget");
        sinfs.push(fs::read(app.join(p))?);
    }
    ensure!(
        !sinfs.is_empty() && sinfs.iter().all(|s| s == &sinfs[0]),
        "ambiguous SINF replication"
    );
    let mut replica_bytes = 0usize;
    for p in paths {
        let name = p.as_string().context("invalid replica")?;
        let p = relative(name)?;
        ensure!(
            name.split('/').any(|s| s == "SC_Info"),
            "replica outside license directory"
        );
        let dest = app.join(p);
        if !dest.exists() {
            replica_bytes += sinfs[0].len();
            ensure!(replica_bytes <= 16 * 1024 * 1024, "license resource budget");
            fs::create_dir_all(dest.parent().context("missing parent")?)?;
            fs::write(dest, &sinfs[0])?;
        }
    }
    Ok(())
}
fn uuid(header: &[u8]) -> Option<String> {
    let count = u32::from_le_bytes(header.get(16..20)?.try_into().ok()?);
    let mut p = 32;
    for _ in 0..count {
        let cmd = u32::from_le_bytes(header.get(p..p + 4)?.try_into().ok()?);
        let size = u32::from_le_bytes(header.get(p + 4..p + 8)?.try_into().ok()?) as usize;
        if cmd == 0x1b && size >= 24 {
            return Some(
                header
                    .get(p + 8..p + 24)?
                    .iter()
                    .map(|b| format!("{b:02x}"))
                    .collect(),
            );
        }
        p += size;
    }
    None
}
fn prepare(v: &Value, c: &Control) -> Result<Value> {
    let work = path(v, "work")?;
    fs::create_dir(work)?;
    let input = work.join("input");
    fs::create_dir(&input)?;
    package::extract_checked(path(v, "ipa")?, &input, || c.check())?;
    let original = inventory(&input, c)?;
    let apps = fs::read_dir(input.join("Payload"))?.collect::<std::io::Result<Vec<_>>>()?;
    ensure!(
        apps.len() == 1
            && apps[0].file_type()?.is_dir()
            && apps[0].path().extension().is_some_and(|s| s == "app"),
        "one app required"
    );
    let app = apps[0].path();
    let d = info(&app.join("Info.plist"))?;
    ensure!(
        Some(text(&d, "CFBundleIdentifier")?) == v["bundle_id"].as_str(),
        "app identity mismatch"
    );
    let exe = text(&d, "CFBundleExecutable")?;
    ensure!(
        relative(exe)?.components().count() == 1,
        "invalid main executable"
    );
    let output = work.join("output");
    copy_tree(&input, &output, c)?;
    replicate(&app)?;
    let mut images = Vec::new();
    for e in walkdir::WalkDir::new(&app) {
        c.check()?;
        let e = e?;
        if !e.file_type().is_file() {
            continue;
        }
        let size = e.metadata()?.len() as usize;
        let mut header = Vec::new();
        fs::File::open(e.path())?
            .take(16 * 1024 * 1024 + 32)
            .read_to_end(&mut header)?;
        let ranges = package::encrypted_ranges_header(&header, size)?;
        ensure!(ranges.len() <= 1, "multiple encrypted ranges unsupported");
        if let Some((off, len, command)) = ranges.first() {
            let relative = e
                .path()
                .strip_prefix(&input)?
                .to_str()
                .context("invalid image path")?;
            let app_relative = e
                .path()
                .strip_prefix(&app)?
                .to_str()
                .context("invalid image path")?;
            let sinf = e
                .path()
                .parent()
                .context("missing image directory")?
                .join("SC_Info")
                .join(format!("{}.sinf", e.file_name().to_string_lossy()));
            let sinf_size = fs::metadata(sinf).map(|s| s.len()).unwrap_or(0);
            images.push(
                json!({"relative":relative,"app_relative":app_relative,"sha256":hash(e.path(),c)?,
                "is_main":e.path()==app.join(exe),"size":size,"uuid":uuid(&header),"offset":off,"length":len,"command":command,
                "matching_sinf":sinf_size>0 && sinf_size<=4*1024*1024,"sinf_size":sinf_size}),
            );
        }
    }
    ensure!(!images.is_empty(), "no supported encrypted images");
    images.sort_by(|a, b| {
        (
            !a["is_main"].as_bool().unwrap_or(false),
            a["relative"].as_str().unwrap_or(""),
        )
            .cmp(&(
                !b["is_main"].as_bool().unwrap_or(false),
                b["relative"].as_str().unwrap_or(""),
            ))
    });
    let manifest = json!({"schema":1,"bundle_id":text(&d,"CFBundleIdentifier")?,"version":text(&d,"CFBundleShortVersionString")?,
        "build":text(&d,"CFBundleVersion")?,"app_relative":format!("Payload/{}",apps[0].file_name().to_string_lossy()),
        "images":images,"inventory":original});
    let bytes = serde_json::to_vec(&manifest)?;
    ensure!(bytes.len() <= 32 * 1024 * 1024, "manifest budget");
    fs::write(work.join("manifest.json"), bytes)?;
    let mut summary = manifest;
    summary.as_object_mut().unwrap().remove("inventory");
    Ok(json!({"ok":true,"manifest":summary}))
}
fn compare_image(original: &Path, output: &Path, image: &Value, c: &Control) -> Result<()> {
    let off = image["offset"].as_u64().context("missing offset")?;
    let len = image["length"].as_u64().context("missing length")?;
    let cmd = image["command"].as_u64().context("missing command")?;
    ensure!(
        fs::metadata(output)?.len() == image["size"].as_u64().context("missing size")?,
        "image size changed"
    );
    let mut a = fs::File::open(original)?;
    let mut b = fs::File::open(output)?;
    let mut left = [0u8; 65536];
    let mut right = [0u8; 65536];
    let mut pos = 0u64;
    let mut changed = false;
    loop {
        c.check()?;
        let n = a.read(&mut left)?;
        if n == 0 {
            break;
        }
        b.read_exact(&mut right[..n])?;
        for i in 0..n {
            let p = pos + i as u64;
            if (off..off + len).contains(&p) {
                changed |= left[i] != right[i];
            } else if (cmd..cmd + 4).contains(&p) {
                ensure!(right[i] == 0, "cryptid not cleared");
            } else {
                ensure!(left[i] == right[i], "bytes outside allowed range changed");
            }
        }
        pos += n as u64;
    }
    ensure!(changed, "encrypted bytes unchanged");
    Ok(())
}
fn finalize(v: &Value, c: &Control) -> Result<Value> {
    let work = path(v, "work")?;
    let manifest: Value = serde_json::from_slice(&fs::read(work.join("manifest.json"))?)?;
    let expected = manifest["inventory"]
        .as_object()
        .context("missing inventory")?;
    let output = work.join("output");
    let actual = inventory(&output, c)?;
    ensure!(actual.len() == expected.len(), "output inventory changed");
    let images = manifest["images"].as_array().context("missing images")?;
    let receipts = v["receipts"].as_array().context("missing page receipts")?;
    ensure!(
        receipts.len() == images.len(),
        "incomplete decryption receipts"
    );
    for (name, entry) in expected {
        let current = actual.get(name).context("output resource missing")?;
        if let Some(image) = images.iter().find(|i| i["relative"] == *name) {
            let receipt = receipts
                .iter()
                .filter(|r| r["relative"] == *name)
                .collect::<Vec<_>>();
            ensure!(
                receipt.len() == 1
                    && receipt[0]["verified"] == true
                    && receipt[0]["image_sha256"] == image["sha256"],
                "image not verified"
            );
            ensure!(
                receipt[0]["pages_expected"]
                    .as_u64()
                    .is_some_and(|n| n > 0 && receipt[0]["pages_read"].as_u64() == Some(n)),
                "incomplete page coverage"
            );
            let page_size = receipt[0]["page_size"]
                .as_u64()
                .context("missing page size")?;
            let encrypted_len = image["length"]
                .as_u64()
                .context("missing encrypted length")?;
            ensure!(
                page_size.is_power_of_two()
                    && (4096..=65536).contains(&page_size)
                    && receipt[0]["pages_expected"].as_u64()
                        == Some(encrypted_len.div_ceil(page_size)),
                "incomplete page coverage"
            );
            let original = work.join("input").join(relative(name)?);
            ensure!(
                hash(&original, c)? == image["sha256"].as_str().context("missing original hash")?,
                "original image changed"
            );
            compare_image(&original, &output.join(relative(name)?), image, c)?;
            ensure!(current["mode"] == entry["mode"], "executable mode changed");
        } else {
            ensure!(current == entry, "original resource changed");
        }
    }
    let partial = work.join("result.partial");
    let finished = work.join("result.ipa");
    let result = (|| -> Result<()> {
        let mut zip = zip::ZipWriter::new(
            fs::OpenOptions::new()
                .write(true)
                .create_new(true)
                .open(&partial)?,
        );
        let mut files = Vec::new();
        for (name, entry) in &actual {
            c.check()?;
            if entry["directory"] == true {
                zip.add_directory(
                    format!("{name}/"),
                    SimpleFileOptions::default().unix_permissions(0o755),
                )?;
                files.push(json!({"path":format!("{name}/"),"directory":true,"sha256":format!("{:x}",Sha256::digest([]))}));
            } else {
                zip.start_file(
                    name,
                    SimpleFileOptions::default()
                        .unix_permissions(entry["mode"].as_u64().context("mode")? as u32),
                )?;
                let mut f = fs::File::open(output.join(relative(name)?))?;
                let mut buf = [0u8; 65536];
                loop {
                    c.check()?;
                    let n = f.read(&mut buf)?;
                    if n == 0 {
                        break;
                    }
                    zip.write_all(&buf[..n])?;
                }
                files.push(json!({"path":name,"directory":false,"sha256":entry["sha256"]}));
            }
        }
        zip.finish()?.sync_all()?;
        package::verify_output(
            &partial,
            &json!({"schema":3,"archive_resources":{},"archive_directories":[],"files":files}),
        )?;
        c.check()?;
        ensure!(!finished.exists(), "result exists");
        fs::rename(&partial, &finished)?;
        Ok(())
    })();
    if result.is_err() {
        let _ = fs::remove_file(&partial);
    }
    result?;
    Ok(
        json!({"ok":true,"output_path":finished,"image_count":images.len(),"file_count":actual.len()}),
    )
}
fn error_reason(error: &anyhow::Error) -> &'static str {
    for item in error.chain() {
        let text = item.to_string();
        let code = match text.as_str() {
            "cancelled" => "cancelled",
            "Unsafe archive path" | "unsafe relative path" => "unsafe_archive_path",
            "Duplicate archive entry" => "duplicate_archive_entry",
            "Archive links are unsupported" | "links unsupported" => "archive_link_unsupported",
            "one app required" => "single_app_required",
            "app identity mismatch" => "app_identity_mismatch",
            "missing app identity" | "invalid app metadata" => "app_metadata_invalid",
            "ambiguous SINF replication" => "license_replication_ambiguous",
            "missing original SINF mapping" | "invalid SINF mapping" | "invalid replica" => {
                "license_mapping_invalid"
            }
            "SINF outside license directory" | "replica outside license directory" => {
                "license_path_invalid"
            }
            "no supported encrypted images" => "encrypted_images_missing",
            "FAT Mach-O is not supported in the first device pipeline; use an arm64 thin IPA" => {
                "fat_macho_unsupported"
            }
            "Only arm64 Mach-O supported" | "32-bit Mach-O unsupported" => {
                "architecture_unsupported"
            }
            "Invalid Mach-O command table"
            | "Invalid load command"
            | "Unparsed Mach-O load commands"
            | "Encryption range outside file"
            | "Truncated encryption command"
            | "Empty encrypted range" => "macho_malformed",
            "license resource budget"
            | "SINF budget"
            | "manifest budget"
            | "archive budget"
            | "unsupported file"
            | "Expanded IPA exceeds 4 GB"
            | "Archive file exceeds 512 MiB"
            | "Too many archive entries" => "resource_budget_exceeded",
            "incomplete decryption receipts"
            | "image not verified"
            | "incomplete page coverage" => "decryption_incomplete",
            "output inventory changed"
            | "output resource missing"
            | "original resource changed" => "resource_verification_failed",
            "original image changed"
            | "bytes outside allowed range changed"
            | "image size changed"
            | "cryptid not cleared"
            | "encrypted bytes unchanged" => "image_verification_failed",
            _ => continue,
        };
        return code;
    }
    if error
        .chain()
        .any(|e| e.downcast_ref::<zip::result::ZipError>().is_some())
    {
        "zip_invalid"
    } else if error
        .chain()
        .any(|e| e.downcast_ref::<std::io::Error>().is_some())
    {
        "file_io_failed"
    } else {
        "archive_validation_failed"
    }
}
// Synchronous archive-only ABI, invoked on the probe's serial background queue.
// The cancellation context belongs to that queue's retained Objective-C runner.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn crumb_probe_archive(
    ptr: *const std::ffi::c_char,
    ctx: *mut std::ffi::c_void,
    cancel: Cancel,
) -> *mut std::ffi::c_char {
    let result = std::panic::catch_unwind(std::panic::AssertUnwindSafe(|| -> Result<Value> {
        ensure!(!ptr.is_null(), "missing request");
        let bytes = unsafe { std::ffi::CStr::from_ptr(ptr) }.to_bytes();
        ensure!(bytes.len() <= 2 * 1024 * 1024, "request budget");
        let v: Value = serde_json::from_slice(bytes)?;
        let c = Control {
            context: ctx,
            cancel,
        };
        match v["operation"].as_str() {
            Some("prepare") => prepare(&v, &c),
            Some("finalize") => finalize(&v, &c),
            Some("relocate") => {
                copy_tree(path(&v, "from")?, path(&v, "to")?, &c)?;
                Ok(json!({"ok":true}))
            }
            _ => anyhow::bail!("unknown operation"),
        }
    }));
    let value = match result {
        Ok(Ok(v)) => v,
        Ok(Err(e)) => {
            json!({"ok":false,"code":"archive_operation_failed","reason":error_reason(&e),
        "os_errno":e.chain().find_map(|x|x.downcast_ref::<std::io::Error>().and_then(|x|x.raw_os_error()))})
        }
        Err(_) => json!({"ok":false,"code":"archive_operation_panicked"}),
    };
    std::ffi::CString::new(value.to_string())
        .unwrap()
        .into_raw()
}
#[cfg(test)]
mod tests;
