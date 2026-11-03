use anyhow::{Context, Result, ensure};
use serde_json::{Value, json};
use sha2::{Digest, Sha256};
use std::{
    fs,
    io::Read,
    path::{Path, PathBuf},
};
use zip::ZipArchive;
const MAX_BYTES: u64 = 4_000_000_000;
fn input<'a>(v: &'a Value, key: &str) -> Result<&'a str> {
    v[key].as_str().with_context(|| format!("Missing {key}"))
}
fn extract(ipa: &Path, dest: &Path) -> Result<()> {
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
    }
    Ok(())
}
fn plist(path: &Path) -> Result<plist::Value> {
    Ok(plist::Value::from_file(path)?)
}
fn strval<'a>(p: &'a plist::Value, key: &str) -> Result<&'a str> {
    p.as_dictionary()
        .and_then(|d| d.get(key))
        .and_then(|v| v.as_string())
        .with_context(|| format!("Missing {key}"))
}
fn safe_relative(path: &str) -> Result<&Path> {
    let path = Path::new(path);
    ensure!(
        !path.as_os_str().is_empty()
            && path
                .components()
                .all(|c| matches!(c, std::path::Component::Normal(_))),
        "Unsafe SINF replication path"
    );
    Ok(path)
}
pub(super) fn validate_helper_id(bundle: &str, own: &str) -> Result<()> {
    ensure!(
        !bundle.eq_ignore_ascii_case(own),
        "The helper must not replace Crumb"
    );
    ensure!(
        bundle.len() <= 255
            && bundle.contains('.')
            && bundle.split('.').all(|part| !part.is_empty()
                && part.bytes().all(|b| b.is_ascii_alphanumeric() || b == b'-')),
        "Invalid helper bundle identifier"
    );
    Ok(())
}
pub fn prepare(v: &Value) -> Result<String> {
    let work = PathBuf::from(input(v, "work")?);
    // Caller creates an exclusive, disposable per-operation directory.
    fs::create_dir(&work)?;
    let result = (|| {
        let unpack = work.join("input");
        fs::create_dir(&unpack)?;
        extract(Path::new(input(v, "ipa")?), &unpack)?;
        let payload = unpack.join("Payload");
        let apps = fs::read_dir(&payload)?.collect::<std::io::Result<Vec<_>>>()?;
        ensure!(
            apps.len() == 1 && apps[0].file_type()?.is_dir(),
            "IPA must contain exactly one app"
        );
        let source = apps[0].path();
        ensure!(
            source.extension().is_some_and(|x| x == "app"),
            "Invalid payload"
        );
        let app_info = plist(&source.join("Info.plist"))?;
        let bundle = input(v, "helper_bundle")?;
        validate_helper_id(bundle, input(v, "own_bundle")?)?;
        let template = Path::new(input(v, "template")?);
        let helper = work.join("CrumbDecrypt.app");
        copy_tree(template, &helper)?;
        let embedded = helper.join("EmbeddedPayload");
        fs::create_dir(&embedded)?;
        let target = embedded.join("Target");
        fs::rename(&source, &target)?;
        // Apple SINF replication normally happens in installd; embedded resources
        // need the same per-image license placement before helper signing.
        let manifest_path = target.join("SC_Info/Manifest.plist");
        if manifest_path.exists() {
            let manifest = plist(&manifest_path)?;
            let dict = manifest.as_dictionary().context("Invalid SINF manifest")?;
            if let Some(paths) = dict.get("SinfReplicationPaths").and_then(|v| v.as_array()) {
                let main = dict
                    .get("SinfPaths")
                    .and_then(|v| v.as_array())
                    .and_then(|a| a.first())
                    .and_then(|v| v.as_string())
                    .context("Missing main SINF")?;
                let main = safe_relative(main)?;
                let sinf = target.join(main);
                for path in paths {
                    let relative =
                        safe_relative(path.as_string().context("Invalid SINF replication path")?)?;
                    let dest = target.join(relative);
                    if !dest.exists() {
                        fs::create_dir_all(dest.parent().context("Missing SINF parent")?)?;
                        fs::copy(&sinf, dest)?;
                    }
                }
            }
        }
        // Refuse unsupported images and unreasonable per-file sizes before signing.
        let mut encrypted = 0;
        for entry in walkdir::WalkDir::new(&target) {
            let entry = entry?;
            if entry.file_type().is_file() {
                ensure!(
                    entry.metadata()?.len() <= 512 * 1024 * 1024,
                    "Payload file exceeds 512 MiB"
                );
                let mut file = fs::File::open(entry.path())?;
                let mut header = Vec::new();
                (&mut file)
                    .take(16 * 1024 * 1024 + 32)
                    .read_to_end(&mut header)?;
                encrypted +=
                    encrypted_ranges_header(&header, entry.metadata()?.len() as usize)?.len();
            }
        }
        ensure!(encrypted > 0, "No supported encrypted images found");
        let meta = json!({"job_id":input(v,"job_id")?,"app_name":apps[0].file_name().to_string_lossy(),"bundle_id":strval(&app_info,"CFBundleIdentifier")?});
        fs::write(helper.join("job.json"), serde_json::to_vec(&meta)?)?;
        let mut info = plist(&helper.join("Info.plist"))?;
        info.as_dictionary_mut()
            .context("Invalid helper plist")?
            .insert("CFBundleIdentifier".into(), bundle.into());
        info.to_file_binary(helper.join("Info.plist"))?;
        let profile_bytes = fs::read(input(v, "profile")?)?;
        let profile = apple_codesign_quick::ProvisioningProfile::parse(&profile_bytes)?;
        let ent = profile.entitlements();
        let appid = ent
            .get("application-identifier")
            .and_then(|v| v.as_string())
            .context("Profile has no app identifier")?;
        let expected = format!("{}.{}", profile.team_id(), bundle);
        ensure!(
            appid == expected
                || appid
                    .strip_suffix('*')
                    .is_some_and(|prefix| expected.starts_with(prefix)),
            "Profile does not permit the helper bundle identifier"
        );
        ensure!(
            profile
                .plist()
                .get("ExpirationDate")
                .and_then(|v| v.as_date())
                .is_some_and(|d| std::time::SystemTime::from(d) > std::time::SystemTime::now()),
            "Provisioning profile expired"
        );
        let store = p12_keystore::KeyStore::from_pkcs12(
            &fs::read(input(v, "p12")?)?,
            input(v, "password")?,
            p12_keystore::Pkcs12ImportPolicy::Strict,
        )?;
        let (_, identity) = store
            .private_key_chain()
            .context("Certificate has no private key")?;
        let leaf = identity
            .certs()
            .first()
            .context("Missing signing certificate")?
            .as_der();
        ensure!(
            profile
                .plist()
                .get("DeveloperCertificates")
                .and_then(|v| v.as_array())
                .is_some_and(|certs| certs.iter().any(|c| c.as_data() == Some(leaf))),
            "Certificate is absent from provisioning profile"
        );
        let signer = apple_codesign_quick::RustCryptoCmsSigner::from_der(
            identity.key().as_der(),
            leaf,
            identity.certs().iter().skip(1).map(|c| c.as_der()),
        )?;
        let mut ent = ent.clone();
        ent.insert("application-identifier".into(), expected.into());
        let mut settings =
            apple_codesign_quick::BundleSigningSettings::new(profile.team_id(), ent, Some(&signer));
        settings.cms_signer = Some(&signer);
        settings.embedded_mobileprovision = Some(&profile_bytes);
        // The vendored signer treats EmbeddedPayload as sealed resources and
        // preserves the target's Apple signature, SINF and encrypted bytes.
        apple_codesign_quick::sign_bundle(&helper, &settings)?;
        Ok(helper.to_string_lossy().to_string())
    })();
    if result.is_err() {
        let _ = fs::remove_dir_all(&work);
    }
    result
}
fn copy_tree(from: &Path, to: &Path) -> Result<()> {
    for e in walkdir::WalkDir::new(from) {
        let e = e?;
        ensure!(!e.file_type().is_symlink(), "Template contains a link");
        let out = to.join(e.path().strip_prefix(from)?);
        if e.file_type().is_dir() {
            fs::create_dir_all(out)?;
        } else {
            fs::copy(e.path(), out)?;
        }
    }
    Ok(())
}
// Strict Mach-O parser shared by output verification and malformed-input tests.
pub fn encrypted_ranges(b: &[u8]) -> Result<Vec<(usize, usize, usize)>> {
    encrypted_ranges_header(b, b.len())
}
fn encrypted_ranges_header(b: &[u8], file_size: usize) -> Result<Vec<(usize, usize, usize)>> {
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
pub fn verify_output(path: &Path, manifest: &Value) -> Result<()> {
    let mut zip = ZipArchive::new(fs::File::open(path)?)?;
    ensure!(
        zip.len() > 0 && zip.len() <= 100000,
        "Invalid result archive"
    );
    let files = manifest["files"]
        .as_array()
        .context("Missing result hashes")?;
    ensure!(files.len() == zip.len(), "Result entry count mismatch");
    let mut total = 0u64;
    let mut names = std::collections::HashSet::new();
    for f in files {
        let name = f["path"].as_str().context("Missing result path")?;
        ensure!(names.insert(name), "Duplicate result hash");
        let mut entry = zip.by_name(name)?;
        ensure!(
            !entry.is_dir() && !entry.is_symlink() && entry.enclosed_name().is_some(),
            "Unsafe result entry"
        );
        total = total
            .checked_add(entry.size())
            .context("Result size overflow")?;
        ensure!(total <= MAX_BYTES, "Result too large");
        // Per-file cap also bounds verification memory on-device.
        ensure!(
            entry.size() <= 512 * 1024 * 1024,
            "Result image exceeds verification memory budget"
        );
        let size = entry.size() as usize;
        let mut digest = Sha256::new();
        let mut header = Vec::new();
        let mut buffer = vec![0u8; 1024 * 1024];
        let mut actual = 0;
        loop {
            let n = entry.read(&mut buffer)?;
            if n == 0 {
                break;
            }
            actual += n;
            ensure!(actual <= size, "Result exceeds declared size");
            digest.update(&buffer[..n]);
            let keep = (16 * 1024 * 1024 + 32usize)
                .saturating_sub(header.len())
                .min(n);
            header.extend_from_slice(&buffer[..keep]);
        }
        ensure!(actual == size, "Result truncated");
        ensure!(
            format!("{:x}", digest.finalize()) == f["sha256"].as_str().context("Missing SHA256")?,
            "Result hash mismatch"
        );
        ensure!(
            encrypted_ranges_header(&header, size)?.is_empty(),
            "Result still contains encrypted Mach-O"
        );
    }
    Ok(())
}
#[cfg(test)]
mod tests {
    use super::*;
    use std::io::Write;
    #[test]
    fn rejects_truncated_and_out_of_bounds_commands() {
        let mut b = vec![0u8; 64];
        b[0..4].copy_from_slice(&0xfeedfacfu32.to_le_bytes());
        b[4..8].copy_from_slice(&0x0100000cu32.to_le_bytes());
        b[16..20].copy_from_slice(&1u32.to_le_bytes());
        b[20..24].copy_from_slice(&24u32.to_le_bytes());
        b[32..36].copy_from_slice(&0x2cu32.to_le_bytes());
        b[36..40].copy_from_slice(&24u32.to_le_bytes());
        b[40..44].copy_from_slice(&60u32.to_le_bytes());
        b[44..48].copy_from_slice(&16u32.to_le_bytes());
        b[48..52].copy_from_slice(&1u32.to_le_bytes());
        assert!(encrypted_ranges(&b).is_err());
    }
    fn result_zip(bytes: &[u8], hash_override: Option<&str>) -> (PathBuf, Value) {
        let path = std::env::temp_dir().join(format!("{}.zip", uuid::Uuid::new_v4()));
        let mut zip = zip::ZipWriter::new(fs::File::create(&path).unwrap());
        zip.start_file(
            "Payload/Target.app/Target",
            zip::write::SimpleFileOptions::default(),
        )
        .unwrap();
        zip.write_all(bytes).unwrap();
        zip.finish().unwrap();
        let digest = format!("{:x}", Sha256::digest(bytes));
        let manifest = json!({"files":[{"path":"Payload/Target.app/Target","sha256":hash_override.unwrap_or(&digest)}]});
        (path, manifest)
    }
    #[test]
    fn rejects_hash_mismatch_and_truncated_archive() {
        let (path, manifest) = result_zip(b"plain resource", Some("wrong hash"));
        assert!(verify_output(&path, &manifest).is_err());
        fs::write(&path, b"PK").unwrap();
        assert!(verify_output(&path, &manifest).is_err());
        fs::remove_file(path).unwrap();
    }
    #[test]
    fn accepts_matching_archive_and_rejects_missing_file_manifest() {
        let (path, manifest) = result_zip(b"plain resource", None);
        assert!(verify_output(&path, &manifest).is_ok());
        assert!(verify_output(&path, &json!({"files":[]})).is_err());
        fs::remove_file(path).unwrap();
    }
    #[test]
    fn rejects_still_encrypted_image_even_with_matching_hash() {
        let mut b = vec![0u8; 64];
        for (offset, value) in [
            (0, 0xfeedfacfu32),
            (4, 0x0100000c),
            (16, 1),
            (20, 24),
            (32, 0x2c),
            (36, 24),
            (40, 56),
            (44, 8),
            (48, 1),
        ] {
            b[offset..offset + 4].copy_from_slice(&value.to_le_bytes());
        }
        let (path, manifest) = result_zip(&b, None);
        assert!(verify_output(&path, &manifest).is_err());
        fs::remove_file(path).unwrap();
    }
    #[test]
    fn rejects_sinf_escape_paths() {
        for path in ["../secret", "/absolute", "SC_Info/../../secret", ""] {
            assert!(safe_relative(path).is_err());
        }
        assert!(safe_relative("Frameworks/Foo.framework/SC_Info/Foo.sinf").is_ok());
    }
    #[test]
    fn helper_cannot_replace_crumb_or_use_invalid_identifiers() {
        assert!(validate_helper_id("COM.CERTLIUM.CRUMB", "com.certlium.crumb").is_err());
        assert!(validate_helper_id("../Crumb", "com.certlium.crumb").is_err());
        assert!(validate_helper_id("com.certlium.crumb.decrypt", "com.certlium.crumb").is_ok());
    }
    #[test]
    fn empty_data_is_not_an_encrypted_image() {
        assert!(encrypted_ranges(b"resource").unwrap().is_empty());
    }
    #[test]
    fn rejects_zip_path_traversal() {
        let temp = std::env::temp_dir().join(uuid::Uuid::new_v4().to_string());
        fs::create_dir(&temp).unwrap();
        let path = temp.join("bad.zip");
        let mut z = zip::ZipWriter::new(fs::File::create(&path).unwrap());
        let opt = zip::write::SimpleFileOptions::default();
        z.start_file("../escape", opt).unwrap();
        z.write_all(b"x").unwrap();
        z.finish().unwrap();
        assert!(extract(&path, &temp.join("out")).is_err());
        fs::remove_dir_all(temp).unwrap();
    }
}
