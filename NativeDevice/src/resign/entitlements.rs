//! Preserve profile-authorized capabilities and report changed identities.
use super::identifiers::{pattern, permits_id, remap_id};
use crate::signing;
use anyhow::{Result, ensure};
use std::collections::BTreeSet;

#[cfg(test)]
pub(super) fn compatible_entitlements(
    original: plist::Dictionary,
    grants: &plist::Dictionary,
    team: &str,
    old_main: &str,
    new_main: &str,
    bundle: &str,
) -> Result<(plist::Dictionary, Vec<String>)> {
    compatible_with_identity(original, grants, team, old_main, new_main, bundle, None)
}
pub(super) fn compatible_with_identity(
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
pub(super) fn cloud_entitlements(dict: &plist::Dictionary) -> plist::Dictionary {
    dict.iter()
        .filter(|(k, _)| {
            k.starts_with("com.apple.developer.icloud-")
                || k.starts_with("com.apple.developer.ubiquity-")
        })
        .map(|(k, v)| (k.clone(), v.clone()))
        .collect()
}
