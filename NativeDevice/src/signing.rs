use anyhow::{Context, Result, ensure};
use serde_json::{Value, json};
use std::{fs,  path::{Path, PathBuf}};
use apple_codesign_quick::{ProvisioningProfile, BundleSigningSettings, RustCryptoCmsSigner};

fn profile(v: &Value) -> Result<(Vec<u8>, ProvisioningProfile)> {
    let bytes = fs::read(v["profile"].as_str().context("Import a provisioning profile")?)?;
    let profile = ProvisioningProfile::parse(&bytes)?;
    ensure!(profile.plist().get("ExpirationDate").and_then(|x| x.as_date())
        .is_some_and(|d| std::time::SystemTime::from(d) > std::time::SystemTime::now()), "Provisioning profile expired");
    Ok((bytes, profile))
}
fn identity(v: &Value, profile: &ProvisioningProfile) -> Result<RustCryptoCmsSigner> {
    let store = p12_keystore::KeyStore::from_pkcs12(
        &fs::read(v["p12"].as_str().context("Import a certificate")?)?,
        v["password"].as_str().unwrap_or(""), p12_keystore::Pkcs12ImportPolicy::Strict)
        .map_err(|_| anyhow::anyhow!("Could not unlock certificate. Check its password."))?;
    let (_, id) = store.private_key_chain().context("Certificate has no private key")?;
    let leaf = id.certs().first().context("Missing signing certificate")?.as_der();
    ensure!(profile.plist().get("DeveloperCertificates").and_then(|x| x.as_array())
        .is_some_and(|xs| xs.iter().any(|x| x.as_data() == Some(leaf))), "Certificate does not match provisioning profile");
    Ok(RustCryptoCmsSigner::from_der(id.key().as_der(), leaf, id.certs().iter().skip(1).map(|c| c.as_der()))?)
}
pub fn choose_helper(profile: &plist::Dictionary, own: &str) -> Result<String> {
    let pattern = profile.get("Entitlements").and_then(|x| x.as_dictionary())
        .and_then(|x| x.get("application-identifier")).and_then(|x| x.as_string())
        .context("Profile has no app identifier")?.split_once('.').context("Invalid app identifier")?.1;
    let candidate = if let Some(prefix) = pattern.strip_suffix('*') {
        format!("{prefix}{}", if prefix.is_empty() { "com.certlium.crumb.decrypt" } else { "crumbdecrypt" })
    } else { pattern.to_string() };
    super::package::helper_application_id(profile, &candidate, own)?;
    Ok(candidate)
}
pub fn inspect(v: &Value) -> Result<Value> {
    let (_, profile) = profile(v)?;
    let _ = identity(v, &profile)?;
    let bundle = choose_helper(profile.plist(), v["own_bundle"].as_str().unwrap_or("com.certlium.crumb"))?;
    Ok(json!({"phase":"signing-ready", "helper_bundle":bundle,
        "name":profile.plist().get("Name").and_then(|x| x.as_string()).unwrap_or("Signing identity")}))
}
// Read original XML entitlements from the existing thin arm64 code signature.
// Invalid signatures fail closed instead of silently losing app capabilities.
fn entitlements(bytes: &[u8]) -> Result<plist::Dictionary> {
    fn le(b:&[u8], p:usize)->Result<usize>{ Ok(u32::from_le_bytes(b.get(p..p+4).context("Truncated image")?.try_into()?) as usize) }
    fn be(b:&[u8], p:usize)->Result<usize>{ Ok(u32::from_be_bytes(b.get(p..p+4).context("Truncated signature")?.try_into()?) as usize) }
    ensure!(le(bytes,0)? == 0xfeedfacf, "Unsupported signed image");
    let mut p=32;
    for _ in 0..le(bytes,16)? {
        let cmd=le(bytes,p)?; let size=le(bytes,p+4)?;
        ensure!(size >= 8 && p+size <= bytes.len(), "Invalid load command");
        if cmd == 0x1d {
            let off=le(bytes,p+8)?; let len=le(bytes,p+12)?;
            let sig=bytes.get(off..off.checked_add(len).context("Signature overflow")?).context("Missing signature")?;
            ensure!(be(sig,0)? == 0xfade0cc0 && be(sig,4)? <= sig.len(), "Invalid code signature");
            let count=be(sig,8)?; ensure!(count <= 128, "Too many signature blobs");
            for i in 0..count {
                if be(sig,12+i*8)? == 5 {
                    let at=be(sig,16+i*8)?; ensure!(be(sig,at)? == 0xfade7171, "Invalid entitlement blob");
                    let size=be(sig,at+4)?; ensure!((8..=1024*1024).contains(&size), "Invalid entitlement size");
                    let xml=sig.get(at+8..at+size).context("Truncated entitlements")?;
                    return Ok(plist::Value::from_reader(std::io::Cursor::new(xml))?.into_dictionary().context("Invalid entitlements")?);
                }
            }
            return Ok(plist::Dictionary::new());
        }
        p+=size;
    }
    Ok(plist::Dictionary::new())
}
fn permits(grant: &plist::Value, request: &plist::Value) -> bool {
    match (grant, request) {
        (plist::Value::String(g), plist::Value::String(r)) => g == r || g.strip_suffix('*').is_some_and(|p| r.starts_with(p)),
        (plist::Value::Array(gs), plist::Value::Array(rs)) => rs.iter().all(|r| gs.iter().any(|g| permits(g,r))),
        _ => grant == request,
    }
}
fn target_entitlements(original: plist::Dictionary, profile:&ProvisioningProfile, bundle:&str) -> Result<plist::Dictionary> {
    target_entitlements_with_grants(original, profile.entitlements(), profile.team_id(), bundle)
}
fn target_entitlements_with_grants(original:plist::Dictionary, grants:&plist::Dictionary, team:&str, bundle:&str) -> Result<plist::Dictionary> {
    // Validate profile pattern without the helper-only collision rule.
    let appid=grants.get("application-identifier").and_then(|x| x.as_string()).context("Profile has no app identifier")?;
    let prefix=appid.split_once('.').context("Invalid profile app identifier")?.0;
    let id=format!("{prefix}.{bundle}");
    ensure!(permits(&plist::Value::String(appid.into()), &plist::Value::String(id.clone())), "Import a profile that permits {bundle}");
    let mut result=plist::Dictionary::new();
    for (key, mut value) in original {
        if ["application-identifier","com.apple.developer.team-identifier","get-task-allow","beta-reports-active"].contains(&key.as_str()) { continue; }
        if key == "keychain-access-groups" {
            if let Some(groups)=value.as_array_mut() {
                for group in groups {
                    if let Some(s)=group.as_string() {
                        if let Some((_,suffix))=s.split_once('.') { *group=plist::Value::String(format!("{prefix}.{suffix}")); }
                    }
                }
            }
        }
        ensure!(grants.get(&key).is_some_and(|g| permits(g,&value)), "Profile does not permit app capability: {key}");
        result.insert(key,value);
    }
    result.insert("application-identifier".into(), id.into());
    result.insert("com.apple.developer.team-identifier".into(),team.into());
    if let Some(v)=grants.get("get-task-allow") { result.insert("get-task-allow".into(),v.clone()); }
    Ok(result)
}
pub fn sign(v:&Value, job:&super::Job) -> Result<Value> {
    let (profile_bytes, profile)=profile(v)?;
    let signer=identity(v,&profile)?;
    let work=PathBuf::from(v["work"].as_str().context("Missing workspace")?);
    fs::create_dir(&work)?;
    let result=(|| {
        super::package::extract(Path::new(v["ipa"].as_str().context("Missing IPA")?), &work)?;
        let apps=fs::read_dir(work.join("Payload"))?.collect::<std::io::Result<Vec<_>>>()?;
        ensure!(apps.len()==1 && apps[0].path().extension().is_some_and(|e|e=="app"),"Invalid app payload");
        let app=apps[0].path();
        let mut settings=BundleSigningSettings::new(profile.team_id(),plist::Dictionary::new(),Some(&signer));
        settings.embedded_mobileprovision=Some(&profile_bytes);
        for entry in walkdir::WalkDir::new(&app) {
            let entry=entry?;
            if !entry.file_type().is_dir() || !entry.path().join("Info.plist").exists() { continue; }
            let is_app=entry.path().extension().is_some_and(|e|e=="app" || e=="appex");
            if !is_app { continue; }
            let info=plist::Value::from_file(entry.path().join("Info.plist"))?.into_dictionary().context("Invalid bundle metadata")?;
            let bundle=info.get("CFBundleIdentifier").and_then(|x|x.as_string()).context("Missing bundle ID")?;
            ensure!(bundle != v["own_bundle"].as_str().unwrap_or("com.certlium.crumb"),"Cannot replace Crumb");
            let executable=info.get("CFBundleExecutable").and_then(|x|x.as_string()).context("Missing app executable")?;
            ensure!(Path::new(executable).components().count()==1,"Invalid app executable");
            let bytes=fs::read(entry.path().join(executable))?;
            ensure!(super::package::encrypted_ranges(&bytes)?.is_empty(),"Decrypt this app before signing");
            let ent=target_entitlements(entitlements(&bytes)?,&profile,bundle)?;
            settings.entitlements_by_bundle_id.insert(bundle.into(),ent);
            settings.embedded_mobileprovisions_by_bundle_id.insert(bundle.into(),&profile_bytes);
        }
        apple_codesign_quick::sign_bundle(&app,&settings)?;
        let output=PathBuf::from(v["output"].as_str().context("Missing signed output")?);
        let partial=output.with_extension("partial");
        let mut zip=zip::ZipWriter::new(fs::File::create(&partial)?);
        let options=zip::write::SimpleFileOptions::default().compression_method(zip::CompressionMethod::Deflated).unix_permissions(0o755);
        for entry in walkdir::WalkDir::new(&work).min_depth(1) {
            let entry=entry?; let name=entry.path().strip_prefix(&work)?.to_str().context("Invalid archive path")?.replace('\\',"/");
            if entry.file_type().is_dir() { zip.add_directory(format!("{name}/"),options)?; }
            else { zip.start_file(name,options)?; std::io::copy(&mut fs::File::open(entry.path())?,&mut zip)?; }
        }
        zip.finish()?.sync_all()?;
        // A failed signing job never replaces the previously published IPA.
        ensure!(!job.cancelled.load(std::sync::atomic::Ordering::SeqCst),"Cancelled");
        fs::rename(&partial,&output)?;
        Ok(json!({"phase":"signed","path":output}))
    })();
    let _=fs::remove_dir_all(&work);
    if result.is_err() { if let Some(out)=v["output"].as_str() { let _=fs::remove_file(Path::new(out).with_extension("partial")); } }
    result
}
#[cfg(test)] mod tests {
 use super::*;
 #[test] fn wildcard_grants_are_scoped() {
  assert!(permits(&plist::Value::String("TEAM.com.test.*".into()),&plist::Value::String("TEAM.com.test.app".into())));
  assert!(!permits(&plist::Value::String("TEAM.com.test.*".into()),&plist::Value::String("TEAM.com.other.app".into())));
 }
 #[test] fn helper_identifier_is_automatic_and_respects_profile_scope() {
    let own="com.certlium.crumb";
    for (pattern,expected) in [("*","com.certlium.crumb.decrypt"),("org.example.*","org.example.crumbdecrypt"),("org.example.helper","org.example.helper")] {
      let mut ent=plist::Dictionary::new();ent.insert("application-identifier".into(),format!("TEAM.{pattern}").into());
      let mut profile=plist::Dictionary::new();profile.insert("Entitlements".into(),ent.into());
      profile.insert("ApplicationIdentifierPrefix".into(),vec![plist::Value::String("TEAM".into())].into());
      assert_eq!(choose_helper(&profile,own).unwrap(),expected);
    }
 }
 #[test] fn final_signing_preserves_allowed_capabilities_and_rejects_missing_grants() {
    let mut grants=plist::Dictionary::new();grants.insert("application-identifier".into(),"PREFIX.org.example.*".into());
    grants.insert("keychain-access-groups".into(),vec![plist::Value::String("PREFIX.*".into())].into());
    let mut ent=plist::Dictionary::new();ent.insert("application-identifier".into(),"OLD.org.example.app".into());
    ent.insert("keychain-access-groups".into(),vec![plist::Value::String("OLD.org.example.app".into())].into());
    let output=target_entitlements_with_grants(ent.clone(),&grants,"NEWTEAM","org.example.app").unwrap();
    assert_eq!(output["application-identifier"].as_string(),Some("PREFIX.org.example.app"));
    assert_eq!(output["com.apple.developer.team-identifier"].as_string(),Some("NEWTEAM"));
    assert_eq!(output["keychain-access-groups"].as_array().unwrap()[0].as_string(),Some("PREFIX.org.example.app"));
    ent.insert("com.apple.security.application-groups".into(),vec![plist::Value::String("group.original".into())].into());
    assert!(target_entitlements_with_grants(ent,&grants,"NEWTEAM","org.example.app").is_err());
    assert!(target_entitlements_with_grants(plist::Dictionary::new(),&grants,"NEWTEAM","com.other.app").is_err());
 }
 #[test] fn malformed_signature_fails() { assert!(entitlements(b"truncated").is_err()); }
}
