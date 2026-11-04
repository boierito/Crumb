use super::*;
struct Fixture(PathBuf);
impl Fixture {
    fn new() -> Self {
        let path = std::env::temp_dir().join(format!("crumb-root-{}", uuid::Uuid::new_v4()));
        fs::create_dir_all(path.join("input/Payload/Test.app")).unwrap();
        fs::create_dir(path.join("helper")).unwrap();
        Self(path)
    }
    fn put(&self, name: &str, bytes: &[u8]) {
        let path = self.0.join("input").join(name);
        fs::create_dir_all(path.parent().unwrap()).unwrap();
        fs::write(path, bytes).unwrap();
    }
    fn preserve(&self) -> (Value, Value) {
        preserve_archive_resources(&self.0.join("input"), &self.0.join("helper")).unwrap()
    }
}
impl Drop for Fixture {
    fn drop(&mut self) {
        let _ = fs::remove_dir_all(&self.0);
    }
}
#[test]
fn original_root_files_and_empty_directories_are_preserved_outside_payload() {
    let fixture = Fixture::new();
    let files: [(&str, &[u8]); 3] = [
        (
            "META-INF/com.apple.ZipMetadata.plist",
            b"original metadata signature",
        ),
        ("iTunesMetadata.plist", b"bplist00original-store-fields"),
        ("iTunesArtwork", b"original artwork"),
    ];
    for (name, bytes) in files {
        fixture.put(name, bytes);
    }
    fixture.put("Payload/Test.app/Test", b"payload must not be duplicated");
    fs::create_dir(fixture.0.join("input/META-INF/empty")).unwrap();
    let (resources, dirs) = fixture.preserve();
    assert_eq!(resources.as_object().unwrap().len(), 3);
    assert!(resources.get("Payload/Test.app/Test").is_none());
    assert!(dirs.as_array().unwrap().contains(&json!("META-INF/empty")));
    assert!(
        !dirs
            .as_array()
            .unwrap()
            .iter()
            .any(|v| v.as_str().unwrap().starts_with("Payload"))
    );
    for (name, bytes) in files {
        let item = &resources[name];
        let saved = fs::read(
            fixture
                .0
                .join("helper/CrumbArchiveResources")
                .join(item["file"].as_str().unwrap()),
        )
        .unwrap();
        assert_eq!(B64.decode(saved).unwrap(), bytes);
        assert_eq!(item["size"].as_u64(), Some(bytes.len() as u64));
        assert_eq!(item["sha256"], format!("{:x}", Sha256::digest(bytes)));
        assert_eq!(fs::read(fixture.0.join("input").join(name)).unwrap(), bytes);
    }
}
#[test]
fn root_resource_streaming_is_not_limited_to_the_license_metadata_size() {
    let fixture = Fixture::new();
    let bytes = vec![b'x'; 5 * 1024 * 1024 + 13];
    fixture.put("iTunesArtwork", &bytes);
    let (resources, _) = fixture.preserve();
    let item = &resources["iTunesArtwork"];
    let saved = fs::read(
        fixture
            .0
            .join("helper/CrumbArchiveResources")
            .join(item["file"].as_str().unwrap()),
    )
    .unwrap();
    assert_eq!(B64.decode(saved).unwrap(), bytes);
    assert_eq!(item["sha256"], format!("{:x}", Sha256::digest(&bytes)));
}
fn result_fixture(
    missing_metadata: bool,
    changed_metadata: bool,
    missing_dir: bool,
) -> (Fixture, PathBuf, Value) {
    let fixture = Fixture::new();
    let path = fixture.0.join("result.ipa");
    let mut zip = zip::ZipWriter::new(fs::File::create(&path).unwrap());
    let mut hashes = Vec::new();
    let original = b"bplist00original-store-fields";
    for (name, bytes) in [
        ("Payload/Test.app/Test", b"plain image".as_slice()),
        (
            "META-INF/com.apple.ZipMetadata.plist",
            b"original signature".as_slice(),
        ),
        (
            "iTunesMetadata.plist",
            if changed_metadata {
                b"changed fields".as_slice()
            } else {
                original.as_slice()
            },
        ),
    ] {
        if missing_metadata && name == "iTunesMetadata.plist" {
            continue;
        }
        zip.start_file(name, zip::write::SimpleFileOptions::default())
            .unwrap();
        zip.write_all(bytes).unwrap();
        hashes.push(json!({"path":name,"sha256":format!("{:x}",Sha256::digest(bytes))}));
    }
    for name in ["META-INF/", "META-INF/empty/"] {
        if missing_dir && name == "META-INF/empty/" {
            continue;
        }
        zip.add_directory(name, zip::write::SimpleFileOptions::default())
            .unwrap();
        hashes.push(
            json!({"path":name,"directory":true,"sha256":format!("{:x}",Sha256::digest([]))}),
        );
    }
    zip.finish().unwrap();
    let manifest = json!({"schema":3,"files":hashes,"archive_resources":{
        "iTunesMetadata.plist":{"size":original.len(),"sha256":format!("{:x}",Sha256::digest(original))}},
        "archive_directories":["META-INF","META-INF/empty"]});
    (fixture, path, manifest)
}
#[test]
fn verified_output_rejects_missing_changed_metadata_and_missing_empty_directories() {
    for flags in [
        (true, false, false),
        (false, true, false),
        (false, false, true),
    ] {
        let (_fixture, path, manifest) = result_fixture(flags.0, flags.1, flags.2);
        assert!(verify_output(&path, &manifest).is_err());
    }
}
#[test]
fn verified_output_accepts_safe_directories_and_rejects_mislabeled_entries() {
    let (_fixture, path, mut manifest) = result_fixture(false, false, false);
    verify_output(&path, &manifest).unwrap();
    let directory = manifest["files"]
        .as_array_mut()
        .unwrap()
        .iter_mut()
        .find(|f| f["directory"] == true)
        .unwrap();
    directory["directory"] = json!(false);
    assert!(verify_output(&path, &manifest).is_err());
}
