use super::*;
use std::sync::{Mutex, atomic::AtomicBool};
fn entitlements(id: &str) -> plist::Dictionary {
    let mut d = plist::Dictionary::new();
    d.insert(
        "application-identifier".into(),
        format!("ORIGINAL.{id}").into(),
    );
    d.insert(
        "com.apple.developer.team-identifier".into(),
        "ORIGINAL".into(),
    );
    d.insert(
        "com.apple.developer.icloud-services".into(),
        vec![plist::Value::from("CloudKit")].into(),
    );
    d.insert(
        "com.apple.developer.icloud-container-identifiers".into(),
        vec![plist::Value::from("iCloud.example.app")].into(),
    );
    d.insert(
        "com.apple.developer.icloud-container-environment".into(),
        "Production".into(),
    );
    d.insert(
        "com.apple.developer.ubiquity-kvstore-identifier".into(),
        "ORIGINAL.example.app".into(),
    );
    d.insert(
        "keychain-access-groups".into(),
        vec![plist::Value::from("ORIGINAL.example.app")].into(),
    );
    d.insert(
        "com.apple.security.application-groups".into(),
        vec![plist::Value::from("group.example")].into(),
    );
    d
}
fn fixture(id: &str) -> Vec<u8> {
    let mut bytes = vec![0u8; 512];
    let mut u32at = |p: usize, n: u32| bytes[p..p + 4].copy_from_slice(&n.to_le_bytes());
    for (p, n) in [
        (0, 0xfeedfacf),
        (4, 0x0100000c),
        (12, 2),
        (16, 3),
        (20, 160),
        (32, 0x19),
        (36, 72),
        (104, 0x19),
        (108, 72),
        (176, 0x1d),
        (180, 16),
        (184, 512),
    ] {
        u32at(p, n)
    }
    bytes[40..46].copy_from_slice(b"__TEXT");
    bytes[112..122].copy_from_slice(b"__LINKEDIT");
    bytes[64..72].copy_from_slice(&4096u64.to_le_bytes());
    bytes[80..88].copy_from_slice(&512u64.to_le_bytes());
    bytes[136..144].copy_from_slice(&4096u64.to_le_bytes());
    bytes[144..152].copy_from_slice(&512u64.to_le_bytes());

    bytes[256..512].fill(0x25);
    sign_macho_data(
        Path::new("fixture"),
        &bytes,
        &MachOSigningConfig::new(id, "ORIGINAL", &entitlements(id), None),
    )
    .unwrap()
}
fn be(bytes: &[u8], at: usize) -> usize {
    u32::from_be_bytes(bytes[at..at + 4].try_into().unwrap()) as usize
}
fn signature(bytes: &[u8]) -> usize {
    let mut p = 32;
    for _ in 0..u32::from_le_bytes(bytes[16..20].try_into().unwrap()) {
        let cmd = u32::from_le_bytes(bytes[p..p + 4].try_into().unwrap());
        if cmd == 0x1d {
            return u32::from_le_bytes(bytes[p + 8..p + 12].try_into().unwrap()) as usize;
        }
        p += u32::from_le_bytes(bytes[p + 4..p + 8].try_into().unwrap()) as usize;
    }
    panic!("signature absent")
}
fn xml(bytes: &[u8]) -> Option<plist::Dictionary> {
    let p = signature(bytes);
    for i in 0..be(bytes, p + 8) {
        if be(bytes, p + 12 + i * 8) == 5 {
            let q = p + be(bytes, p + 16 + i * 8);
            let n = be(bytes, q + 4);
            return Some(
                plist::Value::from_reader(std::io::Cursor::new(&bytes[q + 8..q + n]))
                    .unwrap()
                    .into_dictionary()
                    .unwrap(),
            );
        }
    }
    None
}
fn der_only(mut bytes: Vec<u8>) -> Vec<u8> {
    let p = signature(&bytes);
    for i in 0..be(&bytes, p + 8) {
        if be(&bytes, p + 12 + i * 8) == 5 {
            bytes[p + 12 + i * 8..p + 16 + i * 8].copy_from_slice(&0x100u32.to_be_bytes());
        }
    }
    bytes
}
fn job(cancelled: bool) -> crate::Job {
    crate::Job {
        state: Mutex::new(serde_json::json!({})),
        cancelled: AtomicBool::new(cancelled),
        notify: tokio::sync::Notify::new(),
    }
}
#[test]
fn der_only_exports_gain_xml_without_losing_any_capability_or_code() {
    for id in ["example.app", "example.app.widget"] {
        let original = der_only(fixture(id));
        assert!(xml(&original).is_none());
        assert_eq!(
            crate::signing::entitlements(&original).unwrap(),
            entitlements(id)
        );
        let output = image(Path::new(id), &original).unwrap();
        assert_eq!(xml(&output).unwrap(), entitlements(id));
        assert_eq!(
            crate::signing::entitlements(&der_only(output.clone())).unwrap(),
            entitlements(id)
        );
        assert_eq!(&original[256..512], &output[256..512]);
    }
}
#[test]
fn zip_normalization_retains_resources_profiles_empty_folders_and_modes() {
    let root = std::env::temp_dir().join(format!("crumb-compat-{}", uuid::Uuid::new_v4()));
    fs::create_dir(&root).unwrap();
    let path = root.join("raw.ipa");
    let mut zip = zip::ZipWriter::new(fs::File::create(&path).unwrap());
    let resource = zip::write::SimpleFileOptions::default().unix_permissions(0o644);
    for (name, bytes) in [
        ("Payload/App.app/Info.plist", b"original info".as_slice()),
        (
            "Payload/App.app/embedded.mobileprovision",
            b"original profile",
        ),
        (
            "Payload/App.app/Settings.bundle/Root.plist",
            b"original settings",
        ),
        (
            "Payload/App.app/_CodeSignature/CodeResources",
            b"original seal",
        ),
    ] {
        zip.start_file(name, resource).unwrap();
        zip.write_all(bytes).unwrap();
    }
    zip.add_directory("Payload/App.app/Settings.bundle/Empty/", resource)
        .unwrap();
    zip.start_file("Payload/App.app/App", resource.unix_permissions(0o755))
        .unwrap();
    zip.write_all(&der_only(fixture("example.app"))).unwrap();
    zip.finish().unwrap();
    let before = fs::read(&path).unwrap();
    assert!(normalize(&path, &job(true)).is_err());
    assert_eq!(before, fs::read(&path).unwrap());
    normalize(&path, &job(false)).unwrap();
    let mut zip = zip::ZipArchive::new(fs::File::open(&path).unwrap()).unwrap();
    assert_eq!(zip.len(), 6);
    let mut b = Vec::new();
    zip.by_name("Payload/App.app/App")
        .unwrap()
        .read_to_end(&mut b)
        .unwrap();
    assert_eq!(xml(&b).unwrap(), entitlements("example.app"));
    assert_eq!(
        zip.by_name("Payload/App.app/App")
            .unwrap()
            .unix_mode()
            .unwrap()
            & 0o777,
        0o755
    );
    for (name, expected) in [
        ("Payload/App.app/Info.plist", b"original info".as_slice()),
        (
            "Payload/App.app/embedded.mobileprovision",
            b"original profile",
        ),
        (
            "Payload/App.app/Settings.bundle/Root.plist",
            b"original settings",
        ),
        (
            "Payload/App.app/_CodeSignature/CodeResources",
            b"original seal",
        ),
    ] {
        let mut b = Vec::new();
        zip.by_name(name).unwrap().read_to_end(&mut b).unwrap();
        assert_eq!(b, expected);
    }
    assert!(
        zip.by_name("Payload/App.app/Settings.bundle/Empty/")
            .unwrap()
            .is_dir()
    );
    drop(zip);
    fs::remove_dir_all(root).unwrap();
}

#[test]
fn executable_verification_rejects_code_and_swift_metadata_changes() {
    let original = fixture("example.app");
    let mut output = image(Path::new("App"), &original).unwrap();
    output[300] ^= 1;
    assert!(unchanged_program(&original, &output).is_err());
}
