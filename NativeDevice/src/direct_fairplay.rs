//! Fail-closed, no-helper FairPlay diagnostic on one licensed IPA from Crumb's downloads.
//! Extracts at most one arm64 encrypted Mach-O into Crumb's private temp directory,
//! attempts one read-only mapping and native remap, and destroys the temp file.
//! No output IPA, no signing/repackaging, no system path writes, no raw trace logs.
use super::*;
use idevice::RsdService;
use idevice::rsd::RsdHandshake;
use idevice::services::os_trace_relay::OsTraceRelayClient;
use idevice::tcp::handle::AdapterHandle;
use sha2::{Digest, Sha256};
use std::{
    fs::{self, File, OpenOptions},
    io::{Read, Write},
    os::fd::AsRawFd,
    path::{Path, PathBuf},
    time::Duration,
};
use tokio::time::Instant;
use zip::ZipArchive;

const MAX_IPA: u64 = 4 * 1024 * 1024 * 1024;
const MAX_IMAGE: u64 = 512 * 1024 * 1024;
const MAX_COMMANDS: usize = 8 * 1024 * 1024;
const MAX_ARCHIVE_ENTRIES: usize = 100_000;
const MAX_CANDIDATES: usize = 128;

struct TempImage {
    path: PathBuf,
    root: PathBuf,
    sinf_staged: bool,
    sinf_bytes: u64,
}
impl Drop for TempImage {
    fn drop(&mut self) {
        // Includes sensitive SC_Info resources. Never leave them in Crumb's tmp.
        let _ = fs::remove_dir_all(&self.root);
    }
}
// Only the original matching SINF is staged beside this test Mach-O.
// The original IPA is never rewritten or installed as an app/helper.
// Sidecar absence is explicit in the report and must not be taken as a
// FairPlay-daemon policy decision.
fn stage_matching_sinf(
    zip: &mut ZipArchive<File>, original_entry: &str, executable: &Path,
) -> anyhow::Result<(bool,u64)> {
    let parts: Vec<&str> = original_entry.split('/').collect();
    anyhow::ensure!(parts.len() >= 3 && is_app_candidate(original_entry),
        "Invalid executable archive path");
    let leaf = parts.last().copied().unwrap_or("");
    anyhow::ensure!(!leaf.is_empty() && leaf.len() < 256, "Invalid executable basename");
    let original_dir = original_entry.rsplit_once('/').map(|x| x.0).unwrap_or("");
    let source_candidates = [
        format!("{original_dir}/SC_Info/{leaf}.sinf"),
        format!("Payload/{}/SC_Info/{leaf}.sinf", parts[1]),
    ];
    for name in source_candidates {
        let Ok(mut entry) = zip.by_name(&name) else { continue; };
        anyhow::ensure!(!entry.is_dir() && !entry.is_symlink()
            && entry.enclosed_name().is_some() && entry.size() > 0
            && entry.size() <= 4 * 1024 * 1024,
            "Unsupported original license resource");
        let target = executable.parent()
            .ok_or_else(||anyhow::anyhow!("Missing executable parent"))?
            .join("SC_Info").join(format!("{leaf}.sinf"));
        fs::create_dir_all(target.parent().unwrap())?;
        let mut f = OpenOptions::new().write(true).create_new(true).open(&target)?;
        let copied = std::io::copy(&mut entry, &mut f)?;
        anyhow::ensure!(copied == entry.size(), "Truncated original SINF");
        f.sync_all()?;
        return Ok((true,copied));
    }
    Ok((false,0))
}
#[derive(Clone, Copy, Debug)]
struct Range {
    cryptoff: usize, cryptsize: usize, cputype: u32, cpusubtype: u32,
}
fn le(b: &[u8], at: usize) -> anyhow::Result<u32> {
    let end = at.checked_add(4).ok_or_else(||anyhow::anyhow!("Mach-O integer overflow"))?;
    let s = b.get(at..end).ok_or_else(||anyhow::anyhow!("Truncated Mach-O"))?;
    Ok(u32::from_le_bytes(s.try_into()?))
}
fn parse_crypt_range(b: &[u8], length: u64) -> anyhow::Result<Option<Range>> {
    if b.len() < 32 || le(b, 0)? != 0xfeedfacf || le(b,4)? != 0x0100000c {
        return Ok(None);
    }
    let ncmds = le(b,16)? as usize;
    let cmdsize = le(b,20)? as usize;
    anyhow::ensure!(ncmds <= 4096 && cmdsize <= MAX_COMMANDS
        && cmdsize.checked_add(32).is_some_and(|end|end<=b.len()),
        "Invalid load command geometry");
    let end = 32 + cmdsize;
    let mut offset = 32usize;
    let mut found = None;
    for _ in 0..ncmds {
        anyhow::ensure!(offset + 8 <= end, "Mach-O load command overrun");
        let cmd = le(b, offset)?;
        let size = le(b, offset+4)? as usize;
        anyhow::ensure!(size >= 8 && size <= end-offset, "Invalid load command length");
        if cmd == 0x2c || cmd == 0x21 {
            anyhow::ensure!(size >= 20, "Truncated encryption command");
            let start = le(b, offset+8)? as usize;
            let count = le(b, offset+12)? as usize;
            let cryptid = le(b,offset+16)?;
            if cryptid != 0 {
                anyhow::ensure!(found.is_none() && count>0, "Ambiguous encrypted ranges");
                anyhow::ensure!((start as u64) < length
                    && count as u64 <= length - start as u64,
                    "Encrypted range exceeds file");
                found = Some(Range {
                    cryptoff:start, cryptsize:count,
                    cputype:le(b,4)?, cpusubtype:le(b,8)?,
                });
            }
        }
        offset += size;
    }
    anyhow::ensure!(offset == end, "Mismatched Mach-O load commands");
    Ok(found)
}
fn is_app_candidate(name: &str) -> bool {
    let parts: Vec<_> = name.split('/').collect();
    parts.len() >= 3 && parts[0] == "Payload" && parts[1].ends_with(".app")
        && parts.iter().all(|p|!p.is_empty() && *p != "." && *p != "..")
        && !name.ends_with('/') && name.len() <= 4096
}
fn extract_one(ipa: &Path) -> anyhow::Result<(TempImage, Range, u64, String)> {
    let stat = fs::symlink_metadata(ipa)?;
    anyhow::ensure!(stat.file_type().is_file() && stat.len() >= 4 && stat.len() <= MAX_IPA,
        "IPA must be a regular file of supported size");
    let mut archive = ZipArchive::new(File::open(ipa)?)?;
    anyhow::ensure!(archive.len() <= MAX_ARCHIVE_ENTRIES, "Too many IPA entries");
    let mut candidates=0usize;
    for i in 0..archive.len() {
        let mut entry=archive.by_index(i)?;
        if entry.is_dir() || entry.is_symlink()
            || entry.enclosed_name().is_none()
            || !is_app_candidate(entry.name()) || entry.size() < 32
            || entry.size() > MAX_IMAGE { continue; }
        candidates += 1;
        anyhow::ensure!(candidates <= MAX_CANDIDATES, "Too many candidate executable files");
        let mut header=[0u8;32];
        entry.read_exact(&mut header)?;
        if le(&header,0)? != 0xfeedfacf || le(&header,4)? != 0x0100000c {continue;}
        let commands=le(&header,20)? as usize;
        anyhow::ensure!(commands <= MAX_COMMANDS, "Oversized load commands");
        let mut bytes=Vec::with_capacity(32+commands);
        bytes.extend_from_slice(&header);
        let mut rest=vec![0u8;commands];
        entry.read_exact(&mut rest)?;
        bytes.extend_from_slice(&rest);
        let Some(range)=parse_crypt_range(&bytes, entry.size())? else {continue};
        drop(entry);
        // Preserve the original relative bundle/executable basename. FairPlay
        // derives SC_Info/<name>.sinf from that name; a random *.macho filename
        // is NOT a valid licensed-file control.
        let original_entry = entry.name().to_string();
        let relative = original_entry.strip_prefix("Payload/")
            .ok_or_else(||anyhow::anyhow!("Invalid archive root"))?;
        let root=std::env::temp_dir().join(format!("crumb-fairplay-{}",uuid::Uuid::new_v4()));
        let path=root.join(relative);
        fs::create_dir_all(path.parent().unwrap())?;
        let mut cleanup=TempImage {path,root,sinf_staged:false,sinf_bytes:0};
        let mut out=OpenOptions::new().write(true).create_new(true).open(&cleanup.path)?;
        let mut fresh=archive.by_index(i)?;
        let mut hash=Sha256::new();
        let mut buffer=[0u8; 64*1024];
        let mut total=0u64;
        loop {
            let n=fresh.read(&mut buffer)?;
            if n==0 {break;}
            total=total.checked_add(n as u64).ok_or_else(||anyhow::anyhow!("Overflow"))?;
            anyhow::ensure!(total<=MAX_IMAGE,"Unexpected image size");
            out.write_all(&buffer[..n])?;
            hash.update(&buffer[..n]);
        }
        anyhow::ensure!(total==fresh.size(), "Truncated encrypted image");
        out.sync_all()?;
        drop(out);
        drop(fresh);
        let (sinf_staged,sinf_bytes)=stage_matching_sinf(&mut archive,&original_entry,&cleanup.path)?;
        cleanup.sinf_staged=sinf_staged;
        cleanup.sinf_bytes=sinf_bytes;
        return Ok((cleanup,range,total,format!("{:x}",hash.finalize())));
    }
    anyhow::bail!("IPA has no supported encrypted arm64 Mach-O")
}

#[cfg(all(target_os="ios",target_arch="aarch64"))]
unsafe fn remap_native(address: *mut libc::c_void, size: usize, cpu: u32, subtype: u32) -> i64 {
    // The existing PR #32 uses this on-device syscall ABI (syscall 489).
    // No external system entitlement is invented. iOS may reject the operation.
    let mut ret=address as u64;
    let failed: u64;
    unsafe {
        std::arch::asm!(
            "svc #0x80",
            "cset {failed:w}, cs",
            failed=out(reg) failed,
            inout("x0") ret,
            in("x1") size as u64,
            in("x2") 2u64,
            in("x3") cpu as u64,
            in("x4") subtype as u64,
            in("x16") 489u64,
            options(nostack)
        );
    }
    if failed != 0 { -(ret as i64) } else { ret as i64 }
}

fn map_and_probe(path: &Path, size: u64, range: Range) -> Value {
    #[cfg(not(all(target_os="ios",target_arch="aarch64")))]
    { let _=(path,size,range); return json!({"stage":"host_not_supported",
        "mapping_attempted":false,"remap_attempted":false,"pages_decrypted":0}); }
    #[cfg(all(target_os="ios",target_arch="aarch64"))]
    {
        let f=match File::open(path) {
            Ok(f)=>f,Err(e)=>return json!({"stage":"open","errno":e.raw_os_error(),
                "mapping_attempted":false,"remap_attempted":false,"pages_decrypted":0})
        };
        let page=unsafe {libc::sysconf(libc::_SC_PAGESIZE)};
        if page<=0 {return json!({"stage":"page_geometry","pages_decrypted":0});}
        let page=page as usize;
        if range.cryptoff%page!=0 || range.cryptsize==0
            || range.cryptoff.checked_add(range.cryptsize).is_none_or(|end|end>size as usize) {
            return json!({"stage":"unsupported_range_alignment","pages_decrypted":0});
        }
        let len=range.cryptoff+range.cryptsize;
        let Some(mapped_len)=len.checked_add(page-1).map(|n|n/page*page) else {
            return json!({"stage":"mapping_overflow","pages_decrypted":0});
        };
        let p=unsafe{libc::mmap(std::ptr::null_mut(),mapped_len,libc::PROT_READ,
             libc::MAP_PRIVATE,f.as_raw_fd(),0)};
        if p==libc::MAP_FAILED {
            return json!({"stage":"mmap","mapping_attempted":true,"mapping_errno":
                std::io::Error::last_os_error().raw_os_error(),"remap_attempted":false,
                "pages_decrypted":0});
        }
        // Do NOT read or export decrypted bytes in this diagnostic. A successful
        // syscall still needs independent per-page verification before claiming decrypt.
        let result=unsafe{remap_native((p as *mut u8).add(range.cryptoff).cast(),
            range.cryptsize,range.cputype,range.cpusubtype)};
        unsafe{libc::munmap(p,mapped_len);}
        json!({"stage":if result==0 {"remap_accepted_unverified"} else {"remap_rejected"},
            "mapping_attempted":true,"mapping_errno":0,"remap_attempted":true,
            "remap_result":if result<0 {-1} else {result},
            "remap_errno":if result<0 {-result} else {0},
            "encrypted_pages":(range.cryptsize+page-1)/page,"pages_decrypted":0})
    }
}

fn probe_local(ipa: &str) -> Value {
    match extract_one(Path::new(ipa)) {
        Ok((image,range,size,sha))=>{
            let direct=map_and_probe(&image.path,size,range);
            json!({"mode":"no_helper_readonly_direct_remap","schema":3,
                "input_image_bytes":size,"encrypted_size":range.cryptsize,
                "image_sha256":sha,
                "staged_original_executable_basename":true,
                "matching_sinf_staged":image.sinf_staged,
                "matching_sinf_bytes":image.sinf_bytes,
                "native":direct,
                "note":"An SINF next to a data-container file does not grant fairplayd access. A remap return is not decrypted-page evidence; no IPA is exported."})
        }
        Err(e)=>json!({"mode":"no_helper_readonly_direct_remap","schema":2,
            "stage":"preflight_failed","reason":e.to_string()})
    }
}
pub(super) fn local(request:&Value)->anyhow::Result<Value>{
    let ipa=request["ipa"].as_str()
        .ok_or_else(||anyhow::anyhow!("Missing IPA for direct FairPlay diagnostic"))?;
    Ok(json!({"phase":"probe_complete","report":{
        "local_fairplay":probe_local(ipa),
        "trace":{"status":"not_started_without_authenticated_tunnel"}}}))
}
pub(super) async fn correlated(
    job:&std::sync::Arc<Job>,adapter:&mut AdapterHandle,handshake:&mut RsdHandshake,
    request:&Value
)->anyhow::Result<Value>{
    let ipa=request["ipa"].as_str()
        .ok_or_else(||anyhow::anyhow!("Missing IPA for direct FairPlay diagnostic"))?.to_string();
    job.update(json!({"phase":"probing"}));
    let trace=if handshake.services.contains_key("com.apple.os_trace_relay.shim.remote") {
        match OsTraceRelayClient::connect_rsd(adapter,handshake).await {
            Ok(client)=>match client.start_trace(None).await {
                Ok(receiver)=>Some(receiver),Err(_)=>None
            }, Err(_)=>None
        }
    } else {None};
    let start=Instant::now();
    let task=tokio::task::spawn_blocking(move||probe_local(&ipa));
    // Capture only counts; strings are discarded, paths and account/IPA identifiers
    // are never exposed. This is a time-window correlation, not causal proof.
    let mut counts=(0u64,0u64,0u64);
    let mut state=if trace.is_some() {"sampling"} else {"unavailable"};
    if let Some(mut receiver)=trace {
        let until=start+Duration::from_millis(2600);
        let mut finished_at: Option<Instant> = None;
        while Instant::now()<until && counts.0<8192 {
            if job.cancelled.load(Ordering::SeqCst) {state="cancelled";break;}
            if task.is_finished() && finished_at.is_none() {finished_at=Some(Instant::now());}
            if finished_at.is_some_and(|t|t.elapsed()>Duration::from_millis(450)) {break;}
            match tokio::time::timeout(Duration::from_millis(150),receiver.next()).await {
                Ok(Ok(log))=>{
                    counts.0+=1;
                    let process=log.filename.to_ascii_lowercase();
                    let image=log.image_name.to_ascii_lowercase();
                    let message=log.message.to_ascii_lowercase();
                    if process.contains("fairplay")||image.contains("fairplay")||message.contains("fairplay") {counts.1+=1;}
                    if message.contains("sandbox")||message.contains("system policy") {counts.2+=1;}
                }, Ok(Err(_))=>{state="stream_ended";break;},Err(_)=>{}
            }
        }
    }
    let truncated=counts.0>=8192;
    if truncated {state="event_limit_reached";} else if state=="sampling" {state="sampled";}
    let direct=task.await?;
    Ok(json!({"phase":"probe_complete","report":{
        "local_fairplay":direct,
        "trace":{"status":state,"event_count":counts.0,"fairplay_mentions":counts.1,
            "policy_mentions":counts.2,"window_ms":start.elapsed().as_millis().min(60000),
            "event_limit_reached":truncated,"max_events":8192},
        "trace_limitations":"Events are an aggregate in a bounded window; no raw messages or causal service return is available."
    }}))
}

#[cfg(test)]
mod tests {
    use super::*;
    fn fixture(commands:&[u8],ncmds:u32)->Vec<u8>{
        let mut h=vec![0u8;32];
        h[0..4].copy_from_slice(&0xfeedfacfu32.to_le_bytes());
        h[4..8].copy_from_slice(&0x0100000cu32.to_le_bytes());
        h[8..12].copy_from_slice(&0u32.to_le_bytes());
        h[16..20].copy_from_slice(&ncmds.to_le_bytes());
        h[20..24].copy_from_slice(&(commands.len() as u32).to_le_bytes());
        h.extend_from_slice(commands);h
    }
    #[test] fn matches_encryption_command_without_claiming_decryption(){
        let mut c=vec![0u8;24];
        for (at,val) in [(0,0x2c_u32),(4,24),(8,16384),(12,16384),(16,1)] {
            c[at..at+4].copy_from_slice(&val.to_le_bytes());
        }
        let p=parse_crypt_range(&fixture(&c,1),32768).unwrap().unwrap();
        assert_eq!(p.cryptoff,16384);assert_eq!(p.cryptsize,16384);
        assert!(parse_crypt_range(&fixture(&c,1),32767).is_err());
    }
    #[test] fn plain_binary_and_unsafe_entry_names_are_rejected(){
        assert!(parse_crypt_range(&fixture(&[],0),4096).unwrap().is_none());
        assert!(!is_app_candidate("../Payload/Test.app/Test"));
        assert!(!is_app_candidate("Payload/Test.app/../../private"));
        assert!(is_app_candidate("Payload/Test.app/Test"));
        assert!(!is_app_candidate("Library/Other.app/Test"));
    }
    #[test] fn corrupt_command_geometry_fails_closed(){
        let h=fixture(&[0u8;8],1);
        assert!(parse_crypt_range(&h, 1234).is_err());
    }
}
