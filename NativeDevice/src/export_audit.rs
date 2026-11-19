//! Compare decompressed IPA contents; only FairPlay ranges and cryptid may change.
use anyhow::{Context, Result, ensure};
use std::{
    collections::{BTreeMap, BTreeSet},
    fs::File,
    io::Read,
    path::Path,
    sync::atomic::Ordering,
};
use zip::ZipArchive;

fn header(archive: &mut ZipArchive<File>, name: &str, size: usize) -> Result<Vec<u8>> {
    let mut entry = archive.by_name(name)?;
    let mut bytes = vec![0; size.min(32)];
    entry.read_exact(&mut bytes)?;
    if bytes.get(..4) == Some(&0xfeedfacfu32.to_le_bytes()) {
        ensure!(bytes.len() == 32, "Truncated original Mach-O");
        let commands = u32::from_le_bytes(bytes[20..24].try_into()?) as usize;
        ensure!(
            commands <= 16 * 1024 * 1024 && commands <= size - 32,
            "Invalid original Mach-O header"
        );
        bytes.resize(32 + commands, 0);
        entry.read_exact(&mut bytes[32..])?;
    }
    Ok(bytes)
}
fn names(archive: &mut ZipArchive<File>) -> Result<BTreeSet<String>> {
    ensure!(archive.len() <= 100_000, "Too many IPA entries");
    let mut names = BTreeSet::new();
    let mut size = 0u64;
    for i in 0..archive.len() {
        let entry = archive.by_index(i)?;
        ensure!(
            !entry.is_symlink() && entry.enclosed_name().is_some(),
            "Invalid IPA entry"
        );
        ensure!(names.insert(entry.name().to_owned()), "Duplicate IPA entry");
        size = size
            .checked_add(entry.size())
            .context("IPA size overflow")?;
        ensure!(size <= 4_000_000_000, "IPA exceeds verification budget");
    }
    Ok(names)
}

pub fn verify(original: &Path, decrypted: &Path, job: &crate::Job) -> Result<()> {
    let mut original = ZipArchive::new(File::open(original)?)?;
    let mut decrypted = ZipArchive::new(File::open(decrypted)?)?;
    let expected = names(&mut original)?;
    let actual = names(&mut decrypted)?;
    for name in &expected {
        ensure!(
            actual.contains(name),
            "Original IPA entry is missing: {name}"
        );
        ensure!(!job.cancelled.load(Ordering::SeqCst), "Cancelled");
        let source = original.by_name(name)?;
        let size = usize::try_from(source.size())?;
        let directory = source.is_dir();
        drop(source);
        let target = decrypted.by_name(name)?;
        ensure!(
            target.size() == size as u64 && target.is_dir() == directory,
            "Original IPA entry changed size/type: {name}"
        );
        if !directory {
            if let (Some(left), Some(right)) =
                (original.by_name(name)?.unix_mode(), target.unix_mode())
            {
                ensure!(
                    left & 0o777 == right & 0o777,
                    "Original file permissions changed: {name}"
                );
            }
        }
        drop(target);
        if directory {
            continue;
        }
        let original_header = header(&mut original, name, size)?;
        let ranges = crate::package::encrypted_ranges_header(&original_header, size)?;
        let result_header = header(&mut decrypted, name, size)?;
        ensure!(
            crate::package::encrypted_ranges_header(&result_header, size)?.is_empty(),
            "Export still contains encrypted code: {name}"
        );
        if !ranges.is_empty() {
            let mut position = 32;
            while position < original_header.len() {
                let command =
                    u32::from_le_bytes(original_header[position..position + 4].try_into()?);
                let length =
                    u32::from_le_bytes(original_header[position + 4..position + 8].try_into()?)
                        as usize;
                if command == 0x1d {
                    ensure!(length >= 16, "Invalid original code signature command");
                    let start = u32::from_le_bytes(
                        original_header[position + 8..position + 12].try_into()?,
                    ) as usize;
                    let length = u32::from_le_bytes(
                        original_header[position + 12..position + 16].try_into()?,
                    ) as usize;
                    ensure!(
                        start <= size && length <= size - start,
                        "Invalid original code signature range"
                    );
                    for &(offset, count, _) in &ranges {
                        ensure!(
                            offset + count <= start || start + length <= offset,
                            "Encryption overlaps original code signature"
                        );
                    }
                }
                position += length;
            }
        }
        for &(offset, _, command) in &ranges {
            ensure!(
                offset >= original_header.len(),
                "Encryption overlaps original load commands"
            );
            ensure!(
                result_header.get(command..command + 4) == Some(&[0, 0, 0, 0]),
                "Export cryptid is not zero: {name}"
            );
        }
        let mut source = original.by_name(name)?;
        let mut target = decrypted.by_name(name)?;
        let mut left = [0u8; 65536];
        let mut right = [0u8; 65536];
        let mut position = 0;
        loop {
            ensure!(!job.cancelled.load(Ordering::SeqCst), "Cancelled");
            let count = source.read(&mut left)?;
            if count == 0 {
                break;
            }
            target.read_exact(&mut right[..count])?;
            for &(offset, length, command) in &ranges {
                for (start, length) in [(offset, length), (command, 4)] {
                    let from = start.max(position);
                    let to = (start + length).min(position + count);
                    if from < to {
                        left[from - position..to - position].fill(0);
                        right[from - position..to - position].fill(0);
                    }
                }
            }
            ensure!(
                left[..count] == right[..count],
                "Original bytes changed outside decryption ranges: {name}"
            );
            position += count;
        }
        ensure!(
            target.read(&mut right[..1])? == 0 && position == size,
            "Export entry size mismatch: {name}"
        );
    }
    // SINF replication is explicitly listed in the original Apple manifest.
    let mut replicas = BTreeMap::new();
    for name in &expected {
        if let Some(root) = name.strip_suffix("/SC_Info/Manifest.plist") {
            let mut entry = original.by_name(name)?;
            ensure!(
                entry.size() <= 16 * 1024 * 1024,
                "Oversized original SINF manifest"
            );
            let mut bytes = Vec::new();
            entry.read_to_end(&mut bytes)?;
            let manifest = plist::Value::from_reader(std::io::Cursor::new(bytes))?;
            let Some(dict) = manifest.as_dictionary() else {
                continue;
            };
            if let Some(paths) = dict.get("SinfReplicationPaths").and_then(|v| v.as_array()) {
                let source = dict
                    .get("SinfPaths")
                    .and_then(|v| v.as_array())
                    .and_then(|v| v.first())
                    .and_then(|v| v.as_string())
                    .context("Missing original SINF")?;
                for path in paths {
                    let path = path
                        .as_string()
                        .context("Invalid original SINF replication path")?;
                    replicas.insert(format!("{root}/{path}"), format!("{root}/{source}"));
                }
            }
        }
    }
    for name in actual.difference(&expected) {
        let mut entry = decrypted.by_name(name)?;
        if entry.is_dir() {
            ensure!(
                expected
                    .iter()
                    .chain(replicas.keys())
                    .any(|file| file.starts_with(name)),
                "Unexpected export directory: {name}"
            );
            continue;
        }
        let source = replicas
            .get(name)
            .context("Unexpected file in decrypted export")?;
        let mut source = original.by_name(source)?;
        ensure!(
            source.size() <= 4 * 1024 * 1024 && source.size() == entry.size(),
            "SINF replica size differs from original"
        );
        let mut left = Vec::new();
        let mut right = Vec::new();
        source.read_to_end(&mut left)?;
        entry.read_to_end(&mut right)?;
        ensure!(left == right, "SINF replica differs from original");
    }
    Ok(())
}

#[cfg(test)]
#[path = "export_audit_tests.rs"]
mod tests;
