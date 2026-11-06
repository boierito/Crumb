//! Plan, sign and package local IPAs using the shared signing identity.
mod entitlements;
mod identifiers;

use super::signing;
use anyhow::{Context, Result, ensure};
use apple_codesign_quick::{BundleSigningSettings, ProvisioningProfile};
#[cfg(test)]
use entitlements::compatible_entitlements;
use entitlements::{cloud_entitlements, compatible_with_identity};
use identifiers::{pattern, permits_id, remap_id, select_main, valid_id};
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
#[path = "resign_tests.rs"]
mod integration;
#[cfg(test)]
mod tests;
