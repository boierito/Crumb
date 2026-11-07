use super::*;
use std::{
    collections::BTreeMap,
    io::Write,
    sync::{Mutex, atomic::AtomicBool},
};

fn job(cancelled: bool) -> crate::Job {
    crate::Job {
        state: Mutex::new(serde_json::json!({})),
        cancelled: AtomicBool::new(cancelled),
        notify: tokio::sync::Notify::new(),
    }
}
fn files() -> BTreeMap<String, Vec<u8>> {
    let mut image = vec![0; 8320];
    for (offset, value) in [
        (0, 0xfeedfacfu32),
        (4, 0x0100000c),
        (16, 2),
        (20, 40),
        (32, 0x2c),
        (36, 24),
        (40, 4096),
        (44, 4096),
        (48, 1),
        (56, 0x1d),
        (60, 16),
        (64, 8192),
        (68, 128),
    ] {
        image[offset..offset + 4].copy_from_slice(&value.to_le_bytes());
    }
    image[4096..8192].fill(0x99);
    let entitlements =
        b"original signature and entitlement bytes: CloudKit / original team / original keychain";
    image[8192..8192 + entitlements.len()].copy_from_slice(entitlements);
    BTreeMap::from([
        ("Payload/Test.app/Test".into(), image),
        (
            "Payload/Test.app/Info.plist".into(),
            b"original Info.plist".to_vec(),
        ),
        (
            "Payload/Test.app/Settings.bundle/Root.plist".into(),
            b"original settings".to_vec(),
        ),
        ("Payload/Test.app/Settings.bundle/Empty/".into(), vec![]),
        (
            "Payload/Test.app/embedded.mobileprovision".into(),
            b"original profile".to_vec(),
        ),
        (
            "Payload/Test.app/_CodeSignature/CodeResources".into(),
            b"original resource seal".to_vec(),
        ),
        (
            "iTunesMetadata.plist".into(),
            b"original store metadata".to_vec(),
        ),
    ])
}
fn decrypted(source: &BTreeMap<String, Vec<u8>>) -> BTreeMap<String, Vec<u8>> {
    let mut result = source.clone();
    let image = result.get_mut("Payload/Test.app/Test").unwrap();
    image[48..52].fill(0);
    image[4096..8192].fill(0x23);
    result
}
fn write(path: &Path, entries: &BTreeMap<String, Vec<u8>>) {
    let mut zip = zip::ZipWriter::new(File::create(path).unwrap());
    for (name, bytes) in entries {
        let options = zip::write::SimpleFileOptions::default();
        if name.ends_with('/') {
            zip.add_directory(name, options).unwrap();
        } else {
            zip.start_file(name, options).unwrap();
            zip.write_all(bytes).unwrap();
        }
    }
    zip.finish().unwrap();
}
fn check(
    source: &BTreeMap<String, Vec<u8>>,
    result: &BTreeMap<String, Vec<u8>>,
    cancelled: bool,
) -> Result<()> {
    let root = std::env::temp_dir().join(format!("crumb-export-{}", uuid::Uuid::new_v4()));
    std::fs::create_dir(&root)?;
    let original = root.join("original.ipa");
    let output = root.join("decrypted.ipa");
    write(&original, source);
    write(&output, result);
    let verdict = verify(&original, &output, &job(cancelled));
    std::fs::remove_dir_all(root)?;
    verdict
}
#[test]
fn preserves_every_original_byte_except_ciphertext_and_cryptid() {
    let source = files();
    let mut result = decrypted(&source);
    result.insert("Payload/".into(), vec![]);
    assert!(check(&source, &result, false).is_ok());
}
#[test]
fn rejects_modified_signature_or_entitlements() {
    let source = files();
    let mut result = decrypted(&source);
    result.get_mut("Payload/Test.app/Test").unwrap()[8210] ^= 1;
    assert!(check(&source, &result, false).is_err());
}
#[test]
fn rejects_missing_and_changed_original_resources_and_metadata() {
    let source = files();
    for name in source.keys() {
        let mut result = decrypted(&source);
        result.remove(name);
        assert!(check(&source, &result, false).is_err(), "{name}");
        if !name.ends_with('/') {
            let mut result = decrypted(&source);
            result.get_mut(name).unwrap()[0] ^= 1;
            assert!(check(&source, &result, false).is_err(), "{name}");
        }
    }
}
#[test]
fn rejects_header_mutation_remaining_encryption_and_cancellation() {
    let source = files();
    let mut result = decrypted(&source);
    result.get_mut("Payload/Test.app/Test").unwrap()[12] ^= 1;
    assert!(check(&source, &result, false).is_err());
    assert!(check(&source, &source, false).is_err());
    assert!(check(&source, &decrypted(&source), true).is_err());
}
#[test]
fn rejects_unexpected_files_including_unlisted_license_replicas() {
    let source = files();
    for name in [
        "job.json",
        "Payload/Test.app/SC_Info/unlisted.sinf",
        "Unexpected/",
    ] {
        let mut result = decrypted(&source);
        result.insert(name.into(), vec![]);
        assert!(check(&source, &result, false).is_err());
    }
}
#[test]
fn permits_only_exact_replicas_requested_by_original_manifest() {
    let mut source = files();
    let root = "Payload/Test.app/";
    let replica = "Frameworks/Example.framework/SC_Info/Example.sinf";
    let mut manifest = plist::Dictionary::new();
    manifest.insert(
        "SinfPaths".into(),
        vec![plist::Value::from("SC_Info/main.sinf")].into(),
    );
    manifest.insert(
        "SinfReplicationPaths".into(),
        vec![plist::Value::from(replica)].into(),
    );
    let mut bytes = Vec::new();
    plist::Value::Dictionary(manifest)
        .to_writer_xml(&mut bytes)
        .unwrap();
    source.insert(format!("{root}SC_Info/Manifest.plist"), bytes);
    source.insert(
        format!("{root}SC_Info/main.sinf"),
        b"original SINF".to_vec(),
    );
    let mut result = decrypted(&source);
    result.insert(format!("{root}{replica}"), b"original SINF".to_vec());
    assert!(check(&source, &result, false).is_ok());
    result.get_mut(&format!("{root}{replica}")).unwrap()[0] ^= 1;
    assert!(check(&source, &result, false).is_err());
}

#[test]
fn keeps_empty_and_short_resources_and_rejects_signature_overlap() {
    let mut source = files();
    source.insert("Payload/Test.app/empty".into(), vec![]);
    source.insert("Payload/Test.app/one-byte".into(), vec![1]);
    assert!(check(&source, &decrypted(&source), false).is_ok());
    source.get_mut("Payload/Test.app/Test").unwrap()[44..48]
        .copy_from_slice(&4224u32.to_le_bytes());
    assert!(check(&source, &decrypted(&source), false).is_err());
}

#[test]
fn rejects_changed_original_executable_permissions() {
    let root = std::env::temp_dir().join(format!("crumb-export-mode-{}", uuid::Uuid::new_v4()));
    std::fs::create_dir(&root).unwrap();
    let source = root.join("source.ipa");
    let output = root.join("result.ipa");
    let original = files();
    write(&source, &original);
    let mut zip = zip::ZipWriter::new(File::create(&output).unwrap());
    for (name, bytes) in decrypted(&original) {
        let options = zip::write::SimpleFileOptions::default().unix_permissions(
            if name == "Payload/Test.app/Test" {
                0o755
            } else {
                0o644
            },
        );
        if name.ends_with('/') {
            zip.add_directory(name, options).unwrap();
        } else {
            zip.start_file(name, options).unwrap();
            zip.write_all(&bytes).unwrap();
        }
    }
    zip.finish().unwrap();
    assert!(
        verify(&source, &output, &job(false))
            .unwrap_err()
            .to_string()
            .contains("permissions")
    );
    std::fs::remove_dir_all(root).unwrap();
}
