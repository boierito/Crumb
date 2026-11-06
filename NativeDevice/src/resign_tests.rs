use super::*;
use std::{
    io::{Read, Write},
    process::Command,
    sync::{Arc, Mutex, atomic::AtomicBool},
    time::{Duration, SystemTime},
};
fn command(root: &Path, args: &[&str]) {
    let status = Command::new("openssl")
        .args(args)
        .current_dir(root)
        .output()
        .expect("openssl is required for signing fixtures");
    assert!(
        status.status.success(),
        "Fixture generation failed: {}",
        args[0]
    );
}
fn profile_fixture(root: &Path, name: &str, pattern: &str) -> PathBuf {
    let mut grants = plist::Dictionary::new();
    grants.insert(
        "application-identifier".into(),
        format!("PREFIX.{pattern}").into(),
    );
    grants.insert("com.apple.developer.team-identifier".into(), "TEAM".into());
    grants.insert(
        "keychain-access-groups".into(),
        vec![plist::Value::String("PREFIX.*".into())].into(),
    );
    let mut p = plist::Dictionary::new();
    p.insert("Name".into(), name.into());
    if name == "enterprise" {
        p.insert("ProvisionsAllDevices".into(), true.into());
    }
    p.insert("Entitlements".into(), grants.into());
    p.insert(
        "TeamIdentifier".into(),
        vec![plist::Value::String("TEAM".into())].into(),
    );
    p.insert(
        "ApplicationIdentifierPrefix".into(),
        vec![plist::Value::String("PREFIX".into())].into(),
    );
    p.insert(
        "ExpirationDate".into(),
        plist::Date::from(SystemTime::now() + Duration::from_secs(86400)).into(),
    );
    p.insert(
        "DeveloperCertificates".into(),
        vec![plist::Value::Data(fs::read(root.join("cert.der")).unwrap())].into(),
    );
    let xml = format!("{name}.plist");
    plist::Value::Dictionary(p)
        .to_file_xml(root.join(&xml))
        .unwrap();
    let out = format!("{name}.mobileprovision");
    command(
        root,
        &[
            "cms",
            "-sign",
            "-binary",
            "-nodetach",
            "-in",
            &xml,
            "-signer",
            "cert.pem",
            "-inkey",
            "key.pem",
            "-outform",
            "DER",
            "-out",
            &out,
        ],
    );
    root.join(out)
}
fn image() -> Vec<u8> {
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
    let mut original = plist::Dictionary::new();
    original.insert(
        "application-identifier".into(),
        "OLD.us.humanzoo.pinlist".into(),
    );
    original.insert(
        "keychain-access-groups".into(),
        vec![plist::Value::String("OLD.us.humanzoo.pinlist".into())].into(),
    );
    original.insert(
        "com.apple.security.application-groups".into(),
        vec![plist::Value::String("group.original".into())].into(),
    );
    let config = apple_codesign_quick::MachOSigningConfig::new(
        "us.humanzoo.pinlist",
        "OLD",
        &original,
        None,
    );
    apple_codesign_quick::sign_macho_data(Path::new("fixture"), &bytes, &config).unwrap()
}
fn request(root: &Path, ipa: &Path, profile: &Path, name: &str) -> Value {
    json!({"ipa":ipa,"profile":profile,"p12":root.join("identity.p12"),"password":"fixture", "work":root.join(name),
        "output":root.join(format!("{name}-signed.ipa")),"own_bundle":"com.certlium.crumb","helper_bundle":"com.certlium.crumb.decrypt"})
}
fn job() -> Arc<super::super::Job> {
    Arc::new(super::super::Job {
        state: Mutex::new(json!({})),
        cancelled: AtomicBool::new(false),
        notify: tokio::sync::Notify::new(),
    })
}
#[test]
fn signs_main_and_extension_with_real_cms_without_mutating_source_or_losing_settings_resources() {
    let root = std::env::temp_dir().join(format!("crumb-resign-{}", uuid::Uuid::new_v4()));
    fs::create_dir(&root).unwrap();
    command(
        &root,
        &[
            "req",
            "-x509",
            "-newkey",
            "rsa:2048",
            "-nodes",
            "-keyout",
            "key.pem",
            "-out",
            "cert.pem",
            "-subj",
            "/CN=CrumbFixture",
            "-days",
            "1",
        ],
    );
    command(
        &root,
        &[
            "x509", "-in", "cert.pem", "-outform", "DER", "-out", "cert.der",
        ],
    );
    command(
        &root,
        &[
            "pkcs12",
            "-export",
            "-inkey",
            "key.pem",
            "-in",
            "cert.pem",
            "-passout",
            "pass:fixture",
            "-out",
            "identity.p12",
        ],
    );
    let main_profile = profile_fixture(&root, "main", "org.example.pinlist");
    let extension_profile = profile_fixture(&root, "widget", "org.example.pinlist.Widget");
    let source = root.join("source.ipa");
    let executable = image();
    let mut z = zip::ZipWriter::new(fs::File::create(&source).unwrap());
    let options = zip::write::SimpleFileOptions::default();
    for (path, id) in [
        ("Payload/Pinlist.app", "us.humanzoo.pinlist"),
        (
            "Payload/Pinlist.app/PlugIns/Widget.appex",
            "us.humanzoo.pinlist.Widget",
        ),
    ] {
        let mut info = plist::Dictionary::new();
        info.insert("CFBundleIdentifier".into(), id.into());
        info.insert("CFBundleExecutable".into(), "Program".into());
        info.insert(
            "CFBundlePackageType".into(),
            if path.ends_with(".appex") {
                "XPC!"
            } else {
                "APPL"
            }
            .into(),
        );
        let mut xml = Vec::new();
        plist::Value::Dictionary(info)
            .to_writer_xml(&mut xml)
            .unwrap();
        z.start_file(format!("{path}/Info.plist"), options).unwrap();
        z.write_all(&xml).unwrap();
        z.start_file(format!("{path}/Program"), options.unix_permissions(0o755))
            .unwrap();
        z.write_all(&executable).unwrap();
    }
    z.start_file("Payload/Pinlist.app/Settings.bundle/Root.plist", options)
        .unwrap();
    z.write_all(b"settings resource must survive").unwrap();
    z.add_directory("Payload/Pinlist.app/Settings.bundle/Empty/", options)
        .unwrap();
    z.finish().unwrap();
    let source_bytes = fs::read(&source).unwrap();
    let enterprise_profile = profile_fixture(&root, "enterprise", "com.certlium.crumb.decrypt");
    let mut enterprise = request(&root, &source, &enterprise_profile, "enterprise-work");
    let enterprise_plan = preview(&enterprise, &job()).unwrap();
    assert_eq!(enterprise_plan["bundle_id"], "us.humanzoo.pinlist");
    assert_eq!(enterprise_plan["bundle_count"], 2);
    enterprise["prepared"] = json!(true);
    sign(&enterprise, &job()).unwrap();
    let mut z =
        zip::ZipArchive::new(fs::File::open(enterprise["output"].as_str().unwrap()).unwrap())
            .unwrap();
    for (path, id) in [
        ("Payload/Pinlist.app", "us.humanzoo.pinlist"),
        (
            "Payload/Pinlist.app/PlugIns/Widget.appex",
            "us.humanzoo.pinlist.Widget",
        ),
    ] {
        let mut bytes = Vec::new();
        z.by_name(&format!("{path}/Info.plist"))
            .unwrap()
            .read_to_end(&mut bytes)
            .unwrap();
        assert_eq!(
            plist::Value::from_reader(std::io::Cursor::new(bytes))
                .unwrap()
                .as_dictionary()
                .unwrap()["CFBundleIdentifier"]
                .as_string(),
            Some(id)
        );
        let mut bytes = Vec::new();
        z.by_name(&format!("{path}/Program"))
            .unwrap()
            .read_to_end(&mut bytes)
            .unwrap();
        assert_eq!(
            signing::entitlements(&bytes).unwrap()["application-identifier"].as_string(),
            Some("PREFIX.com.certlium.crumb.decrypt")
        );
    }
    drop(z);
    let bad = request(&root, &source, &main_profile, "missing-extension");
    assert!(
        preview(&bad, &job())
            .unwrap_err()
            .to_string()
            .contains("extension profile")
    );
    let mut v = request(&root, &source, &main_profile, "prepared");
    v["profiles"] = json!([extension_profile]);
    let plan = preview(&v, &job()).unwrap();
    assert_eq!(plan["bundle_id"], "org.example.pinlist");
    assert_eq!(plan["bundle_count"], 2);
    assert!(
        plan["adjusted_capabilities"]
            .as_array()
            .unwrap()
            .iter()
            .any(|v| v == "com.apple.security.application-groups")
    );
    v["prepared"] = json!(true);
    v["bundle_id"] = json!("not.permitted.identifier");
    assert!(preview(&v, &job()).is_err());
    assert!(Path::new(v["work"].as_str().unwrap()).is_dir());
    v["bundle_id"] = json!("org.example.pinlist");
    assert!(preview(&v, &job()).is_ok());
    let signed = sign(&v, &job()).unwrap();
    assert_eq!(signed["bundle_id"], "org.example.pinlist");
    assert_eq!(fs::read(&source).unwrap(), source_bytes);
    let mut z =
        zip::ZipArchive::new(fs::File::open(v["output"].as_str().unwrap()).unwrap()).unwrap();
    for (path, id) in [
        ("Payload/Pinlist.app", "org.example.pinlist"),
        (
            "Payload/Pinlist.app/PlugIns/Widget.appex",
            "org.example.pinlist.Widget",
        ),
    ] {
        let mut bytes = Vec::new();
        z.by_name(&format!("{path}/Info.plist"))
            .unwrap()
            .read_to_end(&mut bytes)
            .unwrap();
        let info = plist::Value::from_reader(std::io::Cursor::new(bytes)).unwrap();
        assert_eq!(
            info.as_dictionary().unwrap()["CFBundleIdentifier"].as_string(),
            Some(id)
        );
        let mut bytes = Vec::new();
        z.by_name(&format!("{path}/Program"))
            .unwrap()
            .read_to_end(&mut bytes)
            .unwrap();
        let ent = signing::entitlements(&bytes).unwrap();
        assert_eq!(
            ent["application-identifier"].as_string(),
            Some(format!("PREFIX.{id}").as_str())
        );
        assert!(!ent.contains_key("com.apple.security.application-groups"));
        assert!(
            z.by_name(&format!("{path}/embedded.mobileprovision"))
                .is_ok()
        );
    }
    let mut settings = String::new();
    z.by_name("Payload/Pinlist.app/Settings.bundle/Root.plist")
        .unwrap()
        .read_to_string(&mut settings)
        .unwrap();
    assert_eq!(settings, "settings resource must survive");
    assert!(
        z.by_name("Payload/Pinlist.app/Settings.bundle/Empty/")
            .is_ok()
    );
    drop(z);
    let mut cancelled = request(&root, &source, &main_profile, "cancelled");
    cancelled["profiles"] = v["profiles"].clone();
    preview(&cancelled, &job()).unwrap();
    cancelled["prepared"] = json!(true);
    let stopped = job();
    stopped
        .cancelled
        .store(true, std::sync::atomic::Ordering::SeqCst);
    assert!(sign(&cancelled, &stopped).is_err());
    assert!(!Path::new(cancelled["output"].as_str().unwrap()).exists());
    assert_eq!(fs::read(&source).unwrap(), source_bytes);
    fs::remove_dir_all(root).unwrap();
}
