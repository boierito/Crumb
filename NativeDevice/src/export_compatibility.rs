//! Prepare decrypted code for re-signers, without applying a user's profile.
use anyhow::{Context, Result, ensure};
use apple_codesign_quick::{MachOSigningConfig, sign_macho_data};
use std::{
    fs,
    io::{Read, Write},
    path::Path,
    sync::atomic::Ordering,
};

fn unchanged_program(original: &[u8], output: &[u8]) -> Result<()> {
    fn word(b: &[u8], p: usize) -> Result<usize> {
        Ok(u32::from_le_bytes(
            b.get(p..p + 4)
                .context("Truncated signature header")?
                .try_into()?,
        ) as usize)
    }
    let mut p = 32;
    let mut end = None;
    let mut permitted = Vec::new();
    for _ in 0..word(original, 16)? {
        let cmd = word(original, p)?;
        let length = word(original, p + 4)?;
        ensure!(
            length >= 8 && p + length <= original.len(),
            "Invalid original load command"
        );
        if cmd == 0x1d {
            end = Some(word(original, p + 8)?);
            permitted.push((p + 8, p + 16));
        }
        if cmd == 0x19 && original.get(p + 8..p + 18) == Some(b"__LINKEDIT") {
            permitted.push((p + 32, p + 40));
            permitted.push((p + 48, p + 56));
        }
        p += length;
    }
    let end = end.context("Original code signature is missing")?;
    ensure!(
        end <= original.len() && end <= output.len(),
        "Invalid code signature boundary"
    );
    for i in 0..end {
        if !permitted.iter().any(|&(a, b)| i >= a && i < b) {
            ensure!(
                original[i] == output[i],
                "Export normalization changed executable code or metadata"
            );
        }
    }
    Ok(())
}

pub(super) fn image(path: &Path, bytes: &[u8]) -> Result<Vec<u8>> {
    let entitlements = crate::signing::entitlements(bytes)?;
    if entitlements.is_empty() {
        return Ok(bytes.to_vec());
    }
    ensure!(
        crate::package::encrypted_ranges(bytes)?.is_empty(),
        "Code must be decrypted before export normalization"
    );
    let application = entitlements
        .get("application-identifier")
        .and_then(|v| v.as_string());
    let team = entitlements
        .get("com.apple.developer.team-identifier")
        .and_then(|v| v.as_string())
        .or_else(|| application.and_then(|s| s.split_once('.').map(|p| p.0)))
        .unwrap_or("");
    let identifier = application
        .and_then(|s| s.split_once('.').map(|p| p.1))
        .or_else(|| path.file_name().and_then(|p| p.to_str()))
        .context("Missing export code identifier")?;
    let config = MachOSigningConfig::new(identifier, team, &entitlements, None);
    let normalized = sign_macho_data(path, bytes, &config)?;
    unchanged_program(bytes, &normalized)?;
    ensure!(
        crate::signing::entitlements(&normalized)? == entitlements,
        "Export normalization changed original entitlements"
    );
    Ok(normalized)
}

pub fn normalize(path: &Path, job: &crate::Job) -> Result<()> {
    let candidate = path.with_extension("compatible");
    let result = (|| {
        let mut input = zip::ZipArchive::new(fs::File::open(path)?)?;
        let mut output = zip::ZipWriter::new(fs::File::create(&candidate)?);
        for index in 0..input.len() {
            ensure!(!job.cancelled.load(Ordering::SeqCst), "Cancelled");
            let mut entry = input.by_index(index)?;
            let options = zip::write::SimpleFileOptions::default()
                .compression_method(entry.compression())
                .unix_permissions(
                    entry
                        .unix_mode()
                        .unwrap_or(if entry.is_dir() { 0o755 } else { 0o644 })
                        & 0o777,
                );
            let options = if let Some(time) = entry.last_modified() {
                options.last_modified_time(time)
            } else {
                options
            };
            if entry.is_dir() {
                output.add_directory(entry.name(), options)?;
                continue;
            }
            output.start_file(entry.name(), options)?;
            let mut prefix = [0u8; 4];
            let count = entry.size().min(4) as usize;
            entry.read_exact(&mut prefix[..count])?;
            if count == 4 && u32::from_le_bytes(prefix) == 0xfeedfacf {
                ensure!(
                    entry.size() <= 512 * 1024 * 1024,
                    "Export image exceeds memory budget"
                );
                let mut bytes = Vec::with_capacity(entry.size() as usize);
                bytes.extend_from_slice(&prefix);
                entry.read_to_end(&mut bytes)?;
                output.write_all(&image(Path::new(entry.name()), &bytes)?)?;
            } else {
                output.write_all(&prefix[..count])?;
                let mut buffer = [0u8; 65536];
                loop {
                    ensure!(!job.cancelled.load(Ordering::SeqCst), "Cancelled");
                    let n = entry.read(&mut buffer)?;
                    if n == 0 {
                        break;
                    }
                    output.write_all(&buffer[..n])?;
                }
            }
        }
        output.finish()?.sync_all()?;
        ensure!(!job.cancelled.load(Ordering::SeqCst), "Cancelled");
        drop(input);
        fs::rename(&candidate, path)?;
        Ok(())
    })();
    if result.is_err() {
        let _ = fs::remove_file(&candidate);
    }
    result
}

#[cfg(test)]
#[path = "export_compatibility_tests.rs"]
mod tests;
