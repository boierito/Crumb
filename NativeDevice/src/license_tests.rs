//! License resource preservation is independent of iOS/FairPlay acceptance.
use super::*;
struct Fixture(PathBuf);
impl Fixture {
    fn new() -> Self {
        let root = std::env::temp_dir().join(format!("crumb-license-{}", uuid::Uuid::new_v4()));
        fs::create_dir_all(root.join("target")).unwrap();
        fs::create_dir(root.join("helper")).unwrap();
        Self(root)
    }
    fn put(&self, path: &str, bytes: &[u8]) {
        let path = self.0.join("target").join(path);
        fs::create_dir_all(path.parent().unwrap()).unwrap();
        fs::write(path, bytes).unwrap();
    }
    fn preserve(&self) -> Result<Value> {
        preserve_license_resources(&self.0.join("target"), &self.0.join("helper"))
    }
}
impl Drop for Fixture {
    fn drop(&mut self) {
        let _ = fs::remove_dir_all(&self.0);
    }
}
#[test]
fn root_and_extension_metadata_are_preserved_without_altering_originals() {
    let fixture = Fixture::new();
    let files: [(&str, &[u8]); 3] = [
        ("SC_Info/App.sinf", b"original SINF\0bytes"),
        (
            "PlugIns/Share.appex/SC_Info/Share.v4.supp",
            b"original supplemental license",
        ),
        ("PlugIns/Share.appex/SC_Info/Share.supf", b"original SUPF"),
    ];
    for (path, bytes) in files {
        fixture.put(path, bytes);
    }
    fixture.put("Assets.car", b"ordinary resource");
    let result = fixture.preserve().unwrap();
    assert_eq!(result.as_object().unwrap().len(), files.len());
    assert!(result.get("Assets.car").is_none());
    for (path, bytes) in files {
        let item = &result[path];
        let file = item["file"].as_str().unwrap();
        assert!(!file.contains('/') && file.ends_with(".bin"));
        let saved = fs::read(fixture.0.join("helper/CrumbLicenseResources").join(file)).unwrap();
        assert_eq!(B64.decode(saved).unwrap(), bytes);
        assert_eq!(item["size"].as_u64().unwrap(), bytes.len() as u64);
        assert_eq!(item["sha256"], format!("{:x}", Sha256::digest(bytes)));
        assert_eq!(
            fs::read(fixture.0.join("target").join(path)).unwrap(),
            bytes
        );
    }
}
#[test]
fn encrypted_image_cannot_be_treated_as_license_metadata() {
    let fixture = Fixture::new();
    let mut bytes = vec![0u8; 64];
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
        bytes[offset..offset + 4].copy_from_slice(&value.to_le_bytes());
    }
    fixture.put("SC_Info/not-license.supp", &bytes);
    assert!(
        fixture
            .preserve()
            .unwrap_err()
            .to_string()
            .contains("Encrypted image")
    );
}
#[test]
fn oversized_metadata_is_rejected_before_reading_it_into_memory() {
    let fixture = Fixture::new();
    fixture.put("SC_Info/large.supp", b"");
    fs::File::options()
        .write(true)
        .open(fixture.0.join("target/SC_Info/large.supp"))
        .unwrap()
        .set_len(MAX_LICENSE_FILE + 1)
        .unwrap();
    assert!(
        fixture
            .preserve()
            .unwrap_err()
            .to_string()
            .contains("preservation limits")
    );
}
#[cfg(unix)]
#[test]
fn metadata_links_do_not_read_files_outside_the_target() {
    let fixture = Fixture::new();
    fs::create_dir(fixture.0.join("target/SC_Info")).unwrap();
    fs::write(fixture.0.join("outside"), b"not a target resource").unwrap();
    std::os::unix::fs::symlink(
        fixture.0.join("outside"),
        fixture.0.join("target/SC_Info/link.supp"),
    )
    .unwrap();
    assert!(
        fixture
            .preserve()
            .unwrap_err()
            .to_string()
            .contains("contains a link")
    );
}
