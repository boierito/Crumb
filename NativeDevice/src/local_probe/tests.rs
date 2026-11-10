use super::*;
fn control() -> Control { Control { context:std::ptr::null_mut(),cancel:None } }
struct Fixture(PathBuf);
impl Fixture {
    fn new() -> Self { let p=std::env::temp_dir().join(format!("crumb-probe-test-{}",uuid::Uuid::new_v4()));fs::create_dir(&p).unwrap();Self(p) }
    fn work(&self)->PathBuf {self.0.join("work")}
    fn ipa(&self)->PathBuf {self.0.join("input.ipa")}
}
impl Drop for Fixture {fn drop(&mut self){let _=fs::remove_dir_all(&self.0);}}
fn encrypted()->Vec<u8>{
    let mut b=vec![0;8192];
    for (o,v) in [(0,0xfeedfacfu32),(4,0x0100000c),(16,1),(20,24),(32,0x2c),(36,24),(40,4096),(44,4096),(48,1)]{b[o..o+4].copy_from_slice(&v.to_le_bytes());}
    b[4096..].fill(0xa5);b
}
fn zip_fixture(f:&Fixture,extra:&[(&str,&[u8])]) {
    let mut z=zip::ZipWriter::new(fs::File::create(f.ipa()).unwrap());
    let metadata=b"<?xml version=\"1.0\"?><plist version=\"1.0\"><dict><key>CFBundleIdentifier</key><string>test.app</string><key>CFBundleExecutable</key><string>Main</string><key>CFBundleShortVersionString</key><string>1</string><key>CFBundleVersion</key><string>1</string></dict></plist>";
    for (name,data) in [("Payload/Test.app/Info.plist",metadata.as_slice()),("Payload/Test.app/Main",encrypted().as_slice()),("Payload/Test.app/SC_Info/Main.sinf",b"original-license".as_slice()),("Payload/Test.app/PlugIns/Widget.appex/Widget",encrypted().as_slice()),("Payload/Test.app/settings",b"resource".as_slice()),("iTunesMetadata.plist",b"outside-payload".as_slice())].into_iter().chain(extra.iter().copied()){
        z.start_file(name,SimpleFileOptions::default().unix_permissions(0o755)).unwrap();z.write_all(data).unwrap();
    }
    z.add_directory("Payload/Test.app/Empty/",SimpleFileOptions::default()).unwrap();z.finish().unwrap();
}
fn prepare_fixture(f:&Fixture)->Value {prepare(&json!({"work":f.work(),"ipa":f.ipa(),"bundle_id":"test.app"}),&control()).unwrap()}
fn synthesize(f:&Fixture,m:&Value)->Value {
    let mut receipts=Vec::new();
    for image in m["manifest"]["images"].as_array().unwrap(){
        let p=f.work().join("output").join(image["relative"].as_str().unwrap());let mut b=fs::read(&p).unwrap();
        let off=image["offset"].as_u64().unwrap() as usize;let len=image["length"].as_u64().unwrap() as usize;let cmd=image["command"].as_u64().unwrap() as usize;
        b[off..off+len].fill(0x42);b[cmd..cmd+4].fill(0);fs::write(p,b).unwrap();
        receipts.push(json!({"relative":image["relative"],"image_sha256":image["sha256"],"verified":true,"pages_expected":1,"pages_read":1}));
    }json!({"operation":"finalize","work":f.work(),"receipts":receipts})
}
#[test] fn preparation_preserves_complete_archive_and_missing_license_is_explicit(){
    let f=Fixture::new();zip_fixture(&f,&[]);let result=prepare_fixture(&f);
    assert_eq!(result["manifest"]["images"].as_array().unwrap().len(),2);
    assert_eq!(result["manifest"]["inventory"]["Payload/Test.app/Empty"]["directory"],true);
    assert_eq!(result["manifest"]["images"][0]["matching_sinf"],true);
    assert_eq!(result["manifest"]["images"][1]["matching_sinf"],false);
    assert_eq!(fs::read(f.work().join("output/iTunesMetadata.plist")).unwrap(),b"outside-payload");
}
#[test] fn replay_packages_all_images_and_resources_but_is_not_native_decryption(){
    let f=Fixture::new();zip_fixture(&f,&[]);let m=prepare_fixture(&f);let request=synthesize(&f,&m);
    let v=finalize(&request,&control()).unwrap();assert_eq!(v["image_count"],2);
    let mut z=zip::ZipArchive::new(fs::File::open(f.work().join("result.ipa")).unwrap()).unwrap();
    assert!(z.by_name("Payload/Test.app/Empty/").unwrap().is_dir());
    assert!(z.by_name("Payload/Test.app/SC_Info/Main.sinf").is_ok());
}
#[test] fn incomplete_receipts_and_tampered_resources_never_publish(){
    let f=Fixture::new();zip_fixture(&f,&[]);let m=prepare_fixture(&f);let mut request=synthesize(&f,&m);
    let receipts=request["receipts"].clone();request["receipts"]=json!([]);
    assert!(finalize(&request,&control()).is_err());request["receipts"]=receipts;
    fs::write(f.work().join("output/Payload/Test.app/settings"),b"changed").unwrap();
    assert!(finalize(&request,&control()).is_err());assert!(!f.work().join("result.ipa").exists());
}
#[test] fn changed_headers_originals_and_missing_extensions_are_rejected(){
    for mode in 0..3 {
        let f=Fixture::new();zip_fixture(&f,&[]);let m=prepare_fixture(&f);let request=synthesize(&f,&m);
        let p=f.work().join(if mode==1{"input/Payload/Test.app/Main"}else{"output/Payload/Test.app/Main"});
        if mode==2 {fs::remove_file(f.work().join("output/Payload/Test.app/PlugIns/Widget.appex/Widget")).unwrap();}
        else {let mut b=fs::read(&p).unwrap();b[28]=9;fs::write(p,b).unwrap();}
        assert!(finalize(&request,&control()).is_err());assert!(!f.work().join("result.ipa").exists());
    }
}
#[test] fn duplicate_traversal_and_corrupt_zip_are_rejected(){
    for name in ["Payload/Test.app/Main","../escape","Payload/Test.app/../../escape"] {
        let f=Fixture::new();zip_fixture(&f,&[(name,b"bad")]);
        assert!(prepare(&json!({"work":f.work(),"ipa":f.ipa(),"bundle_id":"test.app"}),&control()).is_err());
    }
    let f=Fixture::new();fs::write(f.ipa(),b"PK").unwrap();
    assert!(prepare(&json!({"work":f.work(),"ipa":f.ipa(),"bundle_id":"test.app"}),&control()).is_err());
}
#[test] fn manifest_replicates_original_license_and_rejects_unsafe_mapping(){
    for destination in ["PlugIns/Widget.appex/SC_Info/Widget.sinf","../escape"] {
        let f=Fixture::new();let manifest=format!("<plist version=\"1.0\"><dict><key>SinfPaths</key><array><string>SC_Info/Main.sinf</string></array><key>SinfReplicationPaths</key><array><string>{destination}</string></array></dict></plist>");
        zip_fixture(&f,&[("Payload/Test.app/SC_Info/Manifest.plist",manifest.as_bytes())]);
        let r=prepare(&json!({"work":f.work(),"ipa":f.ipa(),"bundle_id":"test.app"}),&control());
        if destination.starts_with("..") {assert!(r.is_err())}else {
            assert!(r.unwrap()["manifest"]["images"][1]["matching_sinf"].as_bool().unwrap());
            assert_eq!(fs::read(f.work().join("input/Payload/Test.app").join(destination)).unwrap(),b"original-license");
            assert!(!f.work().join("output/Payload/Test.app").join(destination).exists());
        }
    }
}
#[test] fn cancellation_halts_before_archive_work(){
    unsafe extern "C" fn stop(_: *mut std::ffi::c_void)->bool{true}
    let c=Control{context:std::ptr::null_mut(),cancel:Some(stop)};let f=Fixture::new();zip_fixture(&f,&[]);
    assert!(prepare(&json!({"work":f.work(),"ipa":f.ipa(),"bundle_id":"test.app"}),&c).is_err());
}
