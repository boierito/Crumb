//! Bundle identifier validation, selection and boundary-aware remapping.
use anyhow::{Context, Result, ensure};

pub(super) fn pattern(grants: &plist::Dictionary) -> Result<(&str, &str)> {
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
pub(super) fn valid_id(id: &str) -> bool {
    id.len() <= 255
        && id.contains('.')
        && id
            .split('.')
            .all(|s| !s.is_empty() && s.bytes().all(|b| b.is_ascii_alphanumeric() || b == b'-'))
}
pub(super) fn permits_id(grants: &plist::Dictionary, id: &str) -> bool {
    pattern(grants)
        .is_ok_and(|(_, p)| p == id || p.strip_suffix('*').is_some_and(|p| id.starts_with(p)))
}
pub(super) fn select_main(
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
pub(super) fn remap_id(value: &str, old: &str, new: &str) -> String {
    if value == old {
        new.to_owned()
    } else if let Some(suffix) = value.strip_prefix(&format!("{old}.")) {
        format!("{new}.{suffix}")
    } else {
        value.to_owned()
    }
}
