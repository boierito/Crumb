use anyhow::{Context, Result, ensure};
use std::{
    fs,
    io::Read,
    path::Path,
};
use zip::ZipArchive;
const MAX_BYTES: u64 = 4_000_000_000;
pub(super) fn extract(ipa: &Path, dest: &Path) -> Result<()> {
    let mut z = ZipArchive::new(fs::File::open(ipa)?)?;
    ensure!(z.len() <= 100000, "Too many archive entries");
    let mut total = 0u64;
    let mut names = std::collections::HashSet::new();
    for i in 0..z.len() {
        let mut f = z.by_index(i)?;
        ensure!(!f.is_symlink(), "Archive links are unsupported");
        let relative = f.enclosed_name().context("Unsafe archive path")?;
        ensure!(
            !relative.as_os_str().is_empty() && names.insert(relative.clone()),
            "Duplicate archive entry"
        );
        ensure!(
            !f.name().contains('\\') && !f.name().contains('\0'),
            "Invalid archive path"
        );
        total = total
            .checked_add(f.size())
            .context("Archive size overflow")?;
        ensure!(total <= MAX_BYTES, "Expanded IPA exceeds 4 GB");
        let path = dest.join(relative);
        if f.is_dir() {
            fs::create_dir_all(&path)?;
            continue;
        }
        fs::create_dir_all(path.parent().context("Missing parent")?)?;
        let mut out = fs::OpenOptions::new()
            .write(true)
            .create_new(true)
            .open(&path)?;
        let size = std::io::copy(&mut (&mut f).take(MAX_BYTES + 1), &mut out)?;
        ensure!(size == f.size(), "Archive entry size mismatch");
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            fs::set_permissions(
                &path,
                fs::Permissions::from_mode(f.unix_mode().unwrap_or(0o644) & 0o777),
            )?;
        }
    }
    Ok(())
}
pub fn encrypted_ranges(b: &[u8]) -> Result<Vec<(usize, usize, usize)>> {
    encrypted_ranges_header(b, b.len())
}
pub(super) fn encrypted_ranges_header(
    b: &[u8],
    file_size: usize,
) -> Result<Vec<(usize, usize, usize)>> {
    fn u32at(b: &[u8], p: usize) -> Result<u32> {
        Ok(u32::from_le_bytes(
            b.get(p..p + 4).context("Truncated Mach-O")?.try_into()?,
        ))
    }
    if b.len() < 4 {
        return Ok(vec![]);
    }
    let magic = u32at(b, 0)?;
    if magic == 0xbebafeca || magic == 0xbfbafeca {
        anyhow::bail!(
            "FAT Mach-O is not supported in the first device pipeline; use an arm64 thin IPA"
        );
    }
    if magic != 0xfeedfacf {
        ensure!(magic != 0xfeedface, "32-bit Mach-O unsupported");
        return Ok(vec![]);
    }
    ensure!(
        b.len() >= 32 && u32at(b, 4)? == 0x0100000c,
        "Only arm64 Mach-O supported"
    );
    let count = u32at(b, 16)? as usize;
    let size = u32at(b, 20)? as usize;
    ensure!(
        size <= 16 * 1024 * 1024 && size <= b.len() - 32 && count <= size / 8,
        "Invalid Mach-O command table"
    );
    let end = 32 + size;
    let mut p = 32;
    let mut ranges = vec![];
    for _ in 0..count {
        let cmd = u32at(b, p)?;
        let n = u32at(b, p + 4)? as usize;
        ensure!(n >= 8 && n % 4 == 0 && n <= end - p, "Invalid load command");
        if cmd == 0x2c || cmd == 0x21 {
            ensure!(
                n >= if cmd == 0x2c { 24 } else { 20 },
                "Truncated encryption command"
            );
            let off = u32at(b, p + 8)? as usize;
            let len = u32at(b, p + 12)? as usize;
            ensure!(
                off <= file_size && len <= file_size - off,
                "Encryption range outside file"
            );
            if u32at(b, p + 16)? != 0 {
                ensure!(len > 0, "Empty encrypted range");
                ranges.push((off, len, p + 16));
            }
        }
        p += n;
    }
    ensure!(p == end, "Unparsed Mach-O load commands");
    Ok(ranges)
}
