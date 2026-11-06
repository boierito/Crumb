//! Profile-aware application resigning. Helper provisioning is deliberately separate.
use super::signing;
use anyhow::{Context, Result, ensure};
use apple_codesign_quick::{BundleSigningSettings, ProvisioningProfile};
use serde_json::{Value, json};
use std::{
    collections::BTreeSet,
    fs,
    io::Read,
    path::{Path, PathBuf},
};

pub struct Profile {
    pub bytes: Vec<u8>,
    pub parsed: ProvisioningProfile,
}
struct ContextData {
    profiles: Vec<Profile>,
    signer: apple_codesign_quick::RustCryptoCmsSigner,
}
fn context(v: &Value) -> Result<ContextData> {
    let (bytes, parsed) = signing::profile(v)?;
    let (signer, leaf) = signing::identity_material(v, &parsed)?;
    let team = parsed.team_id().to_owned();
    let mut profiles = vec![Profile { bytes, parsed }];
    if let Some(paths) = v.get("profiles") {
        let paths = paths.as_array().context("Invalid extension profiles")?;
        ensure!(paths.len() <= 32, "Too many extension profiles");
        for path in paths {
            let (bytes, parsed) = signing::profile(
                &json!({"profile":path.as_str().context("Invalid extension profile path")?}),
            )?;
            ensure!(
                parsed
                    .plist()
                    .get("DeveloperCertificates")
                    .and_then(|v| v.as_array())
                    .is_some_and(|xs| xs.iter().any(|x| x.as_data() == Some(leaf.as_slice()))),
                "An extension profile belongs to a different certificate"
            );
            ensure!(
                parsed.team_id() == team,
                "Extension profiles must belong to the same signing team"
            );
            profiles.push(Profile { bytes, parsed });
        }
    }
    Ok(ContextData { profiles, signer })
}
fn pattern(grants: &plist::Dictionary) -> Result<(&str, &str)> {
    let (prefix, pattern) = grants
        .get("application-identifier")
        .and_then(|v| v.as_string())
        .context("Profile has no application identifier")?
        .split_once('.')
        .context("Invalid profile identifier")?;
    ensure!(
        !prefix.is_empty() && !pattern.is_empty(),
        "Invalid profile identifier"
    );
    let fixed = pattern.strip_suffix('*').unwrap_or(pattern);
    ensure!(
        !fixed.contains('*')
            && fixed
                .bytes()
                .all(|b| b.is_ascii_alphanumeric() || b == b'.' || b == b'-'),
        "Invalid profile identifier pattern"
    );
    Ok((prefix, pattern))
}
fn valid_id(id: &str) -> bool {
    id.len() <= 255
        && id.contains('.')
        && id
            .split('.')
            .all(|s| !s.is_empty() && s.bytes().all(|b| b.is_ascii_alphanumeric() || b == b'-'))
}
fn permits_id(grants: &plist::Dictionary, id: &str) -> bool {
    pattern(grants)
        .is_ok_and(|(_, p)| p == id || p.strip_suffix('*').is_some_and(|p| id.starts_with(p)))
}
fn select_main(
    grants: &plist::Dictionary,
    original: &str,
    custom: Option<&str>,
    protected: &[&str],
) -> Result<String> {
    ensure!(valid_id(original), "Invalid original app identifier");
    pattern(grants)?;
    let result = custom
        .filter(|s| !s.is_empty())
        .unwrap_or(original)
        .to_owned();
    ensure!(valid_id(&result), "Enter a valid bundle identifier");
    ensure!(
        !protected.iter().any(|id| result.eq_ignore_ascii_case(id)),
        "Cannot replace Crumb or its helper"
    );
    Ok(result)
}
fn remap_id(value: &str, old: &str, new: &str) -> String {
    if value == old {
        new.to_owned()
    } else if let Some(suffix) = value.strip_prefix(&format!("{old}.")) {
        format!("{new}.{suffix}")
    } else {
        value.to_owned()
    }
}
fn compatible_entitlements(
    original: plist::Dictionary,
    grants: &plist::Dictionary,
    team: &str,
    old_main: &str,
    new_main: &str,
    bundle: &str,
) -> Result<(plist::Dictionary, Vec<String>)> {
    compatible_with_identity(original, grants, team, old_main, new_main, bundle, None)
}
fn compatible_with_identity(
    original: plist::Dictionary,
    grants: &plist::Dictionary,
    team: &str,
    old_main: &str,
    new_main: &str,
    bundle: &str,
    profile_identity: Option<&str>,
) -> Result<(plist::Dictionary, Vec<String>)> {
    let (prefix, _) = pattern(grants)?;
    ensure!(
        profile_identity.is_some() || permits_id(grants, bundle),
        "No profile authorizes extension {bundle}. Import a matching extension profile."
    );
    let appid = profile_identity
        .map(str::to_owned)
        .unwrap_or_else(|| format!("{prefix}.{bundle}"));
    let mut result = plist::Dictionary::new();
    let mut changed = BTreeSet::new();
    for (key, mut value) in original {
        let original_value = value.clone();
        if [
            "application-identifier",
            "com.apple.developer.team-identifier",
            "get-task-allow",
            "beta-reports-active",
        ]
        .contains(&key.as_str())
        {
            continue;
        }
        if key == "keychain-access-groups"
            || key == "com.apple.developer.parent-application-identifiers"
        {
            if let Some(groups) = value.as_array_mut() {
                for group in groups.iter_mut() {
                    if let Some(s) = group.as_string() {
                        if let Some((_, id)) = s.split_once('.') {
                            *group = plist::Value::String(format!(
                                "{prefix}.{}",
                                remap_id(id, old_main, new_main)
                            ));
                        }
                    }
                }
            }
        }
        if key == "com.apple.developer.ubiquity-kvstore-identifier" {
            if let (Some(grant), Some(old)) = (grants.get(&key), value.as_string()) {
                if !signing::permits(grant, &value) {
                    if let Some((_, suffix)) = old.split_once('.') {
                        let candidate = plist::Value::String(format!(
                            "{team}.{}",
                            remap_id(suffix, old_main, new_main)
                        ));
                        if signing::permits(grant, &candidate) {
                            value = candidate;
                        }
                    }
                }
            }
        }
        if value != original_value {
            changed.insert(key.clone());
        }
        if let Some(grant) = grants.get(&key) {
            if signing::permits_entitlement(&key, grant, &value) {
                result.insert(key, value);
                continue;
            }
            if let (Some(gs), Some(vs)) = (grant.as_array(), value.as_array()) {
                let filtered: Vec<_> = vs
                    .iter()
                    .filter(|v| gs.iter().any(|g| signing::permits(g, v)))
                    .cloned()
                    .collect();
                if !filtered.is_empty() {
                    result.insert(key.clone(), filtered.into());
                }
            } else if key == "aps-environment" {
                result.insert(key.clone(), grant.clone());
            }
        }
        changed.insert(key);
    }
    // Keychain identity follows the selected application ID, never the helper's ID.
    if !result.contains_key("keychain-access-groups") {
        let group = plist::Value::String(appid.clone());
        if grants
            .get("keychain-access-groups")
            .and_then(|v| v.as_array())
            .is_some_and(|gs| gs.iter().any(|g| signing::permits(g, &group)))
        {
            result.insert("keychain-access-groups".into(), vec![group].into());
        }
    }
    result.insert("application-identifier".into(), appid.into());
    result.insert("com.apple.developer.team-identifier".into(), team.into());
    if let Some(v) = grants.get("get-task-allow") {
        result.insert("get-task-allow".into(), v.clone());
    }
    Ok((result, changed.into_iter().collect()))
}
struct BundlePlan {
    path: PathBuf,
    original: String,
    id: String,
    profile: usize,
    entitlements: plist::Dictionary,
    old_keychain: Vec<String>,
    old_groups: Vec<String>,
    original_cloud: plist::Dictionary,
}
struct Plan {
    app: PathBuf,
    original: String,
    id: String,
    bundles: Vec<BundlePlan>,
    changes: Vec<String>,
}
fn app_path(work: &Path) -> Result<PathBuf> {
    let apps = fs::read_dir(work.join("Payload"))?.collect::<std::io::Result<Vec<_>>>()?;
    ensure!(
        apps.len() == 1
            && apps[0].file_type()?.is_dir()
            && apps[0].path().extension().is_some_and(|e| e == "app"),
        "IPA must contain one app"
    );
    Ok(apps[0].path())
}
fn info(path: &Path) -> Result<plist::Dictionary> {
    plist::Value::from_file(path.join("Info.plist"))?
        .into_dictionary()
        .context("Invalid app metadata")
}
fn text<'a>(info: &'a plist::Dictionary, key: &str) -> Result<&'a str> {
    info.get(key)
        .and_then(|v| v.as_string())
        .with_context(|| format!("Missing {key}"))
}
fn build_plan(work: &Path, v: &Value, ctx: &ContextData) -> Result<Plan> {
    // Refuse incomplete decryption in any executable/framework, not just the main binary.
    use std::io::Read;
    for entry in walkdir::WalkDir::new(work) {
        let entry = entry?;
        if entry.file_type().is_file() {
            let size = entry.metadata()?.len();
            ensure!(size <= 512 * 1024 * 1024, "Image exceeds signing budget");
            let mut header = Vec::new();
            fs::File::open(entry.path())?
                .take(16 * 1024 * 1024 + 32)
                .read_to_end(&mut header)?;
            ensure!(
                super::package::encrypted_ranges_header(&header, size as usize)?.is_empty(),
                "Incomplete decryption: {}",
                entry.path().strip_prefix(work)?.display()
            );
        }
    }
    let app = app_path(work)?;
    let main = info(&app)?;
    let original = text(&main, "CFBundleIdentifier")?.to_owned();
    let mut protected = vec![
        v["own_bundle"].as_str().unwrap_or("com.certlium.crumb"),
        v["helper_bundle"]
            .as_str()
            .unwrap_or("com.certlium.crumb.decrypt"),
    ];
    if let Some(values) = v["protected_bundle_ids"].as_array() {
        ensure!(values.len() <= 16, "Invalid protected app identities");
        for value in values {
            protected.push(value.as_str().context("Invalid protected app identity")?);
        }
    }
    // Like zsign, keep CFBundleIdentifier separate from profile signing identity.
    let id = select_main(
        ctx.profiles[0].parsed.entitlements(),
        &original,
        v["bundle_id"].as_str(),
        &protected,
    )?;
    let mut bundles = Vec::new();
    let mut ids = BTreeSet::new();
    let mut changes = BTreeSet::new();
    for entry in walkdir::WalkDir::new(&app) {
        let entry = entry?;
        if !entry.file_type().is_dir()
            || !entry
                .path()
                .extension()
                .is_some_and(|x| x == "app" || x == "appex")
        {
            continue;
        }
        let path = entry.path();
        let metadata = info(path)?;
        let old = text(&metadata, "CFBundleIdentifier")?.to_owned();
        let target = if path == app {
            id.clone()
        } else {
            ensure!(
                old.starts_with(&format!("{original}.")),
                "Extension {old} does not belong to {original}"
            );
            remap_id(&old, &original, &id)
        };
        ensure!(
            valid_id(&target) && ids.insert(target.clone()),
            "Duplicate or invalid signed app identity"
        );
        ensure!(
            !protected.iter().any(|s| target.eq_ignore_ascii_case(s)),
            "Cannot replace Crumb or its helper"
        );
        let profile = if path == app {
            0
        } else {
            ctx.profiles
                .iter()
                .enumerate()
                .filter(|(_, p)| {
                    permits_id(p.parsed.entitlements(), &target)
                        || p.parsed
                            .plist()
                            .get("ProvisionsAllDevices")
                            .and_then(|v| v.as_boolean())
                            == Some(true)
                })
                .max_by_key(|(_, p)| {
                    if !permits_id(p.parsed.entitlements(), &target) {
                        0
                    } else {
                        pattern(p.parsed.entitlements())
                            .map(|(_, p)| {
                                if p.ends_with('*') {
                                    p.len() + 1
                                } else {
                                    usize::MAX
                                }
                            })
                            .unwrap_or(0)
                    }
                })
                .map(|(i, _)| i)
                .unwrap_or(0)
        };
        let executable = text(&metadata, "CFBundleExecutable")?;
        ensure!(
            Path::new(executable).components().count() == 1
                && executable != "."
                && executable != "..",
            "Invalid executable path"
        );
        let bytes = fs::read(path.join(executable))?;
        ensure!(
            super::package::encrypted_ranges(&bytes)?.is_empty(),
            "This app still contains encrypted code"
        );
        let ent = signing::entitlements(&bytes)?;
        let old_keychain = string_values(&ent, "keychain-access-groups");
        let old_groups = string_values(&ent, "com.apple.security.application-groups");
        let original_cloud = cloud_entitlements(&ent);
        let chosen = &ctx.profiles[profile].parsed;
        let profile_identity = if !permits_id(chosen.entitlements(), &target) {
            chosen
                .entitlements()
                .get("application-identifier")
                .and_then(|v| v.as_string())
        } else {
            None
        };
        let (entitlements, removed) = compatible_with_identity(
            ent,
            ctx.profiles[profile].parsed.entitlements(),
            ctx.profiles[profile].parsed.team_id(),
            &original,
            &id,
            &target,
            profile_identity,
        )?;
        changes.extend(removed);
        bundles.push(BundlePlan {
            path: path.to_owned(),
            original: old,
            id: target,
            profile,
            entitlements,
            old_keychain,
            old_groups,
            original_cloud,
        });
    }
    ensure!(!bundles.is_empty(), "Missing app executable");
    Ok(Plan {
        app,
        original,
        id,
        bundles,
        changes: changes.into_iter().collect(),
    })
}
fn string_values(dict: &plist::Dictionary, key: &str) -> Vec<String> {
    dict.get(key)
        .and_then(|v| v.as_array())
        .map(|vs| {
            vs.iter()
                .filter_map(|v| v.as_string().map(str::to_owned))
                .collect()
        })
        .unwrap_or_default()
}
fn cloud_entitlements(dict: &plist::Dictionary) -> plist::Dictionary {
    dict.iter()
        .filter(|(k, _)| {
            k.starts_with("com.apple.developer.icloud-")
                || k.starts_with("com.apple.developer.ubiquity-")
        })
        .map(|(k, v)| (k.clone(), v.clone()))
        .collect()
}
fn report(plan: &Plan, ctx: &ContextData) -> Value {
    let main = plan
        .bundles
        .iter()
        .find(|b| b.path == plan.app)
        .expect("Main bundle is always planned");
    json!({"phase":"signing-plan","bundle_id":plan.id,"original_bundle_id":plan.original,
      "bundle_count":plan.bundles.len(),"changed_ids":plan.bundles.iter().filter(|b|b.id!=b.original).count(),
      "adjusted_capabilities":plan.changes,"profile_name":ctx.profiles[0].parsed.plist().get("Name").and_then(|v|v.as_string()).unwrap_or("Signing profile"),
      "signing_team":ctx.profiles[0].parsed.team_id(),"application_identifier":main.entitlements.get("application-identifier").and_then(|v|v.as_string()),
      "original_keychain_groups":main.old_keychain,"original_app_groups":main.old_groups,
      "keychain_groups":string_values(&main.entitlements,"keychain-access-groups"),
      "app_groups":string_values(&main.entitlements,"com.apple.security.application-groups"),
      "cloudkit_diagnostics":plan.bundles.iter().filter(|b| !b.original_cloud.is_empty()).map(|b| json!({
          "bundle_id":b.id,"original":b.original_cloud,
          "profile":cloud_entitlements(ctx.profiles[b.profile].parsed.entitlements()),
          "signed":cloud_entitlements(&b.entitlements)
      })).collect::<Vec<_>>()})
}

pub fn preview(v: &Value, job: &super::Job) -> Result<Value> {
    let work = PathBuf::from(v["work"].as_str().context("Missing signing workspace")?);
    let prepared = v["prepared"] == true;
    if prepared {
        ensure!(work.is_dir(), "Signing preview expired");
    } else {
        fs::create_dir(&work)?;
    }
    let result = (|| {
        if !prepared {
            super::package::extract(Path::new(v["ipa"].as_str().context("Missing IPA")?), &work)?;
        }
        ensure!(
            !job.cancelled.load(std::sync::atomic::Ordering::SeqCst),
            "Cancelled"
        );
        let ctx = context(v)?;
        let plan = build_plan(&work, v, &ctx)?;
        Ok(report(&plan, &ctx))
    })();
    if result.is_err() && !prepared {
        let _ = fs::remove_dir_all(work);
    }
    result
}
fn rewrite_links(value: &mut plist::Value, old: &str, new: &str) {
    match value {
        plist::Value::Dictionary(dict) => {
            for (key, value) in dict.iter_mut() {
                if ["WKCompanionAppBundleIdentifier", "WKAppBundleIdentifier"]
                    .contains(&key.as_str())
                {
                    if let Some(s) = value.as_string() {
                        *value = remap_id(s, old, new).into();
                    }
                } else {
                    rewrite_links(value, old, new);
                }
            }
        }
        plist::Value::Array(values) => {
            for v in values {
                rewrite_links(v, old, new)
            }
        }
        _ => {}
    }
}
pub fn sign(v: &Value, job: &super::Job) -> Result<Value> {
    let work = PathBuf::from(v["work"].as_str().context("Missing workspace")?);
    let output = PathBuf::from(v["output"].as_str().context("Missing signed output")?);
    ensure!(
        Path::new(v["ipa"].as_str().context("Missing source IPA")?) != output,
        "Signing must preserve the original decrypted IPA"
    );
    let result = (|| {
        if v["prepared"] != true {
            fs::create_dir(&work)?;
            super::package::extract(Path::new(v["ipa"].as_str().unwrap()), &work)?;
        }
        let ctx = context(v)?;
        let plan = build_plan(&work, v, &ctx)?;
        let mut receipt = report(&plan, &ctx);
        let mut settings = BundleSigningSettings::new(
            ctx.profiles[0].parsed.team_id(),
            plist::Dictionary::new(),
            Some(&ctx.signer),
        );
        settings.embedded_mobileprovision = Some(&ctx.profiles[0].bytes);
        for bundle in &plan.bundles {
            ensure!(
                !job.cancelled.load(std::sync::atomic::Ordering::SeqCst),
                "Cancelled"
            );
            let mut metadata = plist::Value::Dictionary(info(&bundle.path)?);
            rewrite_links(&mut metadata, &plan.original, &plan.id);
            metadata
                .as_dictionary_mut()
                .unwrap()
                .insert("CFBundleIdentifier".into(), bundle.id.clone().into());
            metadata.to_file_binary(bundle.path.join("Info.plist"))?;
            settings
                .entitlements_by_bundle_id
                .insert(bundle.id.clone(), bundle.entitlements.clone());
            settings
                .embedded_mobileprovisions_by_bundle_id
                .insert(bundle.id.clone(), &ctx.profiles[bundle.profile].bytes);
        }
        apple_codesign_quick::sign_bundle(&plan.app, &settings)?;
        let partial = output.with_extension("partial");
        let mut zip = zip::ZipWriter::new(fs::File::create(&partial)?);
        for entry in walkdir::WalkDir::new(&work).min_depth(1) {
            ensure!(
                !job.cancelled.load(std::sync::atomic::Ordering::SeqCst),
                "Cancelled"
            );
            let entry = entry?;
            let name = entry
                .path()
                .strip_prefix(&work)?
                .to_str()
                .context("Invalid archive path")?
                .replace('\\', "/");
            let mut magic = [0u8; 4];
            let executable = !entry.file_type().is_dir()
                && fs::File::open(entry.path())?.read(&mut magic)? == 4
                && matches!(
                    u32::from_le_bytes(magic),
                    0xfeedfacf
                        | 0xfeedface
                        | 0xcefaedfe
                        | 0xcffaedfe
                        | 0xbebafeca
                        | 0xcafebabe
                        | 0xbfbafeca
                        | 0xcafebabf
                );
            let permissions = if entry.file_type().is_dir() || executable {
                0o755
            } else {
                #[cfg(unix)]
                {
                    use std::os::unix::fs::PermissionsExt;
                    entry.metadata()?.permissions().mode() & 0o777
                }
                #[cfg(not(unix))]
                {
                    0o644
                }
            };
            let options = zip::write::SimpleFileOptions::default()
                .compression_method(zip::CompressionMethod::Deflated)
                .unix_permissions(permissions);
            if entry.file_type().is_dir() {
                zip.add_directory(format!("{name}/"), options)?;
            } else {
                zip.start_file(name, options)?;
                std::io::copy(&mut fs::File::open(entry.path())?, &mut zip)?;
            }
        }
        zip.finish()?.sync_all()?;
        ensure!(
            !job.cancelled.load(std::sync::atomic::Ordering::SeqCst),
            "Cancelled"
        );
        fs::rename(&partial, &output)?;
        receipt["phase"] = json!("signed");
        receipt["path"] = json!(output);
        Ok(receipt)
    })();
    let _ = fs::remove_dir_all(&work);
    if result.is_err() {
        let _ = fs::remove_file(output.with_extension("partial"));
    }
    result
}
#[cfg(test)]
mod tests {
    use super::*;
    fn grants(pattern: &str) -> plist::Dictionary {
        let mut d = plist::Dictionary::new();
        d.insert(
            "application-identifier".into(),
            format!("PREFIX.{pattern}").into(),
        );
        d
    }
    #[test]
    fn ids_preserve_original_or_custom_with_shared_profile() {
        assert_eq!(
            select_main(&grants("*"), "us.humanzoo.pinlist", None, &[]).unwrap(),
            "us.humanzoo.pinlist"
        );
        assert_eq!(
            select_main(&grants("org.example.app"), "us.humanzoo.pinlist", None, &[]).unwrap(),
            "us.humanzoo.pinlist"
        );
        let a = select_main(&grants("org.example.*"), "us.humanzoo.pinlist", None, &[]).unwrap();
        assert_eq!(a, "us.humanzoo.pinlist");
        assert_eq!(
            a,
            select_main(&grants("org.example.*"), "us.humanzoo.pinlist", None, &[]).unwrap()
        );
        assert!(
            select_main(
                &grants("com.certlium.crumb.decrypt"),
                "us.humanzoo.pinlist",
                None,
                &["com.certlium.crumb.decrypt"]
            )
            .unwrap()
                == "us.humanzoo.pinlist"
        );
        assert!(
            select_main(
                &grants("org.example.*"),
                "us.humanzoo.pinlist",
                Some("us.humanzoo.pinlist"),
                &[]
            )
            .is_ok()
        );
    }
    #[test]
    fn actual_crumb_bundle_collision_is_rejected_without_reserving_its_profile() {
        assert!(
            select_main(
                &grants("com.certlium.crumb.decrypt"),
                "us.humanzoo.pinlist",
                Some("com.certlium.crumb"),
                &["com.certlium.crumb"]
            )
            .is_err()
        );
    }
    #[test]
    fn wildcard_services_and_kvstore_do_not_grant_empty_cloudkit_containers() {
        // Synthetic equivalent of the supplied profile's permission shapes.
        let mut allowed = grants("*");
        allowed.insert("com.apple.developer.icloud-services".into(), "*".into());
        allowed.insert(
            "com.apple.developer.icloud-container-identifiers".into(),
            plist::Value::Array(vec![]),
        );
        allowed.insert(
            "com.apple.developer.ubiquity-kvstore-identifier".into(),
            "TEAM.*".into(),
        );
        allowed.insert("com.apple.developer.associated-domains".into(), "*".into());
        let mut original = plist::Dictionary::new();
        original.insert(
            "com.apple.developer.icloud-services".into(),
            vec![plist::Value::from("CloudKit")].into(),
        );
        original.insert(
            "com.apple.developer.icloud-container-identifiers".into(),
            vec![plist::Value::from("iCloud.example.app")].into(),
        );
        original.insert(
            "com.apple.developer.icloud-container-environment".into(),
            "Production".into(),
        );
        original.insert(
            "com.apple.developer.ubiquity-kvstore-identifier".into(),
            "OLD.example.app".into(),
        );
        original.insert(
            "com.apple.developer.associated-domains".into(),
            vec![plist::Value::from("applinks:example.com")].into(),
        );
        let (signed, changed) = compatible_entitlements(
            original,
            &allowed,
            "TEAM",
            "example.app",
            "custom.app",
            "custom.app",
        )
        .unwrap();
        assert_eq!(
            signed["com.apple.developer.icloud-services"],
            plist::Value::Array(vec!["CloudKit".into()])
        );
        assert_eq!(
            signed["com.apple.developer.associated-domains"],
            plist::Value::Array(vec!["applinks:example.com".into()])
        );
        assert_eq!(
            signed["com.apple.developer.ubiquity-kvstore-identifier"].as_string(),
            Some("TEAM.custom.app")
        );
        assert!(!signed.contains_key("com.apple.developer.icloud-container-identifiers"));
        assert!(!signed.contains_key("com.apple.developer.icloud-container-environment"));
        assert_eq!(
            changed,
            vec![
                "com.apple.developer.icloud-container-environment",
                "com.apple.developer.icloud-container-identifiers",
                "com.apple.developer.ubiquity-kvstore-identifier"
            ]
        );
        allowed.insert(
            "com.apple.developer.ubiquity-kvstore-identifier".into(),
            "TEAM.other".into(),
        );
        let mut original = plist::Dictionary::new();
        original.insert(
            "com.apple.developer.ubiquity-kvstore-identifier".into(),
            "OLD.example.app".into(),
        );
        let (signed, _) = compatible_entitlements(
            original,
            &allowed,
            "TEAM",
            "example.app",
            "custom.app",
            "custom.app",
        )
        .unwrap();
        assert!(!signed.contains_key("com.apple.developer.ubiquity-kvstore-identifier"));
    }
    #[test]
    fn retains_complete_authorized_icloud_capabilities_without_changing_container() {
        let mut cloud = plist::Dictionary::new();
        cloud.insert(
            "com.apple.developer.icloud-container-environment".into(),
            "Production".into(),
        );
        cloud.insert(
            "com.apple.developer.icloud-container-identifiers".into(),
            vec![plist::Value::from("iCloud.us.humanzoo.pinlist")].into(),
        );
        cloud.insert(
            "com.apple.developer.icloud-services".into(),
            vec![plist::Value::from("CloudKit")].into(),
        );
        cloud.insert(
            "com.apple.developer.ubiquity-kvstore-identifier".into(),
            "OLD.us.humanzoo.pinlist".into(),
        );
        let mut allowed = grants("*");
        allowed.extend(cloud.clone());
        allowed.insert(
            "com.apple.developer.icloud-container-environment".into(),
            vec![
                plist::Value::from("Development"),
                plist::Value::from("Production"),
            ]
            .into(),
        );
        let (signed, changed) = compatible_entitlements(
            cloud.clone(),
            &allowed,
            "TEAM",
            "us.humanzoo.pinlist",
            "us.humanzoo.pinlist",
            "us.humanzoo.pinlist",
        )
        .unwrap();
        assert_eq!(cloud_entitlements(&signed), cloud);
        assert!(changed.is_empty());
        // Environment permission alone cannot authorize a different container.
        allowed.insert(
            "com.apple.developer.icloud-container-identifiers".into(),
            vec![plist::Value::from("iCloud.other")].into(),
        );
        let (signed, changed) = compatible_entitlements(
            cloud,
            &allowed,
            "TEAM",
            "us.humanzoo.pinlist",
            "us.humanzoo.pinlist",
            "us.humanzoo.pinlist",
        )
        .unwrap();
        assert!(!signed.contains_key("com.apple.developer.icloud-container-identifiers"));
        assert_eq!(
            changed,
            vec!["com.apple.developer.icloud-container-identifiers"]
        );
    }
    #[test]
    fn compatible_signing_reports_features_instead_of_failing_or_fabricating_grants() {
        let mut g = grants("*");
        g.insert(
            "keychain-access-groups".into(),
            vec![plist::Value::String("PREFIX.*".into())].into(),
        );
        let mut original = plist::Dictionary::new();
        original.insert(
            "com.apple.security.application-groups".into(),
            vec![plist::Value::String("group.original".into())].into(),
        );
        original.insert(
            "keychain-access-groups".into(),
            vec![plist::Value::String("OLD.us.humanzoo.pinlist".into())].into(),
        );
        let (ent, changes) = compatible_entitlements(
            original,
            &g,
            "TEAM",
            "us.humanzoo.pinlist",
            "org.example.app",
            "org.example.app",
        )
        .unwrap();
        assert_eq!(
            ent["application-identifier"].as_string(),
            Some("PREFIX.org.example.app")
        );
        assert_eq!(
            ent["keychain-access-groups"].as_array().unwrap()[0].as_string(),
            Some("PREFIX.org.example.app")
        );
        assert!(!ent.contains_key("com.apple.security.application-groups"));
        assert_eq!(
            changes,
            vec![
                "com.apple.security.application-groups",
                "keychain-access-groups"
            ]
        );
    }
    #[test]
    fn provisioned_app_groups_survive_resigning_for_settings_storage() {
        let mut grants = grants("*");
        let groups = plist::Value::Array(vec![plist::Value::String(
            "group.us.humanzoo.pinlist".into(),
        )]);
        grants.insert(
            "com.apple.security.application-groups".into(),
            groups.clone(),
        );
        let mut original = plist::Dictionary::new();
        original.insert(
            "com.apple.security.application-groups".into(),
            groups.clone(),
        );
        let (signed, changes) = compatible_entitlements(
            original,
            &grants,
            "TEAM",
            "us.humanzoo.pinlist",
            "us.humanzoo.pinlist",
            "us.humanzoo.pinlist",
        )
        .unwrap();
        assert_eq!(signed["com.apple.security.application-groups"], groups);
        assert!(changes.is_empty());
    }
    #[test]
    fn id_remapping_is_boundary_aware_and_updates_nested_references() {
        assert_eq!(
            remap_id("old.application.Widget", "old.app", "new.app"),
            "old.application.Widget"
        );
        let mut value = plist::Value::Dictionary(plist::Dictionary::new());
        let d = value.as_dictionary_mut().unwrap();
        d.insert("WKCompanionAppBundleIdentifier".into(), "old.app".into());
        d.insert("CFBundleURLName".into(), "old.app".into());
        rewrite_links(&mut value, "old.app", "new.app");
        assert_eq!(
            value.as_dictionary().unwrap()["WKCompanionAppBundleIdentifier"].as_string(),
            Some("new.app")
        );
        assert_eq!(
            value.as_dictionary().unwrap()["CFBundleURLName"].as_string(),
            Some("old.app")
        );
    }
}
#[cfg(test)]
#[path = "resign_tests.rs"]
mod integration;
