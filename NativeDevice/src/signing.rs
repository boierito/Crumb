use anyhow::{Context, Result, ensure};
use apple_codesign_quick::{ProvisioningProfile, RustCryptoCmsSigner};
use serde_json::{Value, json};
use std::fs;

pub(super) fn profile(v: &Value) -> Result<(Vec<u8>, ProvisioningProfile)> {
    let bytes = fs::read(
        v["profile"]
            .as_str()
            .context("Import a provisioning profile")?,
    )?;
    let profile = ProvisioningProfile::parse(&bytes)?;
    ensure!(
        profile
            .plist()
            .get("ExpirationDate")
            .and_then(|x| x.as_date())
            .is_some_and(|d| std::time::SystemTime::from(d) > std::time::SystemTime::now()),
        "Provisioning profile expired"
    );
    Ok((bytes, profile))
}
fn identity(v: &Value, profile: &ProvisioningProfile) -> Result<RustCryptoCmsSigner> {
    Ok(identity_material(v, profile)?.0)
}
pub(super) fn identity_material(
    v: &Value,
    profile: &ProvisioningProfile,
) -> Result<(RustCryptoCmsSigner, Vec<u8>)> {
    let store = p12_keystore::KeyStore::from_pkcs12(
        &fs::read(v["p12"].as_str().context("Import a certificate")?)?,
        v["password"].as_str().unwrap_or(""),
        p12_keystore::Pkcs12ImportPolicy::Strict,
    )
    .map_err(|_| anyhow::anyhow!("Could not unlock certificate. Check its password."))?;
    let (_, id) = store
        .private_key_chain()
        .context("Certificate has no private key")?;
    let leaf = id
        .certs()
        .first()
        .context("Missing signing certificate")?
        .as_der();
    ensure!(
        profile
            .plist()
            .get("DeveloperCertificates")
            .and_then(|x| x.as_array())
            .is_some_and(|xs| xs.iter().any(|x| x.as_data() == Some(leaf))),
        "Certificate does not match provisioning profile"
    );
    Ok((
        RustCryptoCmsSigner::from_der(
            id.key().as_der(),
            leaf,
            id.certs().iter().skip(1).map(|c| c.as_der()),
        )?,
        leaf.to_vec(),
    ))
}
pub fn choose_helper(profile: &plist::Dictionary, own: &str) -> Result<String> {
    let pattern = profile
        .get("Entitlements")
        .and_then(|x| x.as_dictionary())
        .and_then(|x| x.get("application-identifier"))
        .and_then(|x| x.as_string())
        .context("Profile has no app identifier")?
        .split_once('.')
        .context("Invalid app identifier")?
        .1;
    let candidate = if let Some(prefix) = pattern.strip_suffix('*') {
        format!(
            "{prefix}{}",
            if prefix.is_empty() {
                "com.certlium.crumb.decrypt"
            } else {
                "crumbdecrypt"
            }
        )
    } else {
        pattern.to_string()
    };
    super::package::helper_application_id(profile, &candidate, own)?;
    Ok(candidate)
}
pub fn inspect(v: &Value) -> Result<Value> {
    let (_, profile) = profile(v)?;
    let _ = identity(v, &profile)?;
    if v["purpose"] == "installation" {
        return Ok(
            json!({"phase":"installation-signing-ready", "name":profile.plist().get("Name").and_then(|v|v.as_string()).unwrap_or("Signing profile")}),
        );
    }
    let bundle = choose_helper(
        profile.plist(),
        v["own_bundle"].as_str().unwrap_or("com.certlium.crumb"),
    )?;
    Ok(json!({"phase":"signing-ready", "helper_bundle":bundle,
        "name":profile.plist().get("Name").and_then(|x| x.as_string()).unwrap_or("Signing identity")}))
}
// Read original XML entitlements from the existing thin arm64 code signature.
// Invalid signatures fail closed instead of silently losing app capabilities.
pub(super) fn entitlements(bytes: &[u8]) -> Result<plist::Dictionary> {
    fn le(b: &[u8], p: usize) -> Result<usize> {
        Ok(u32::from_le_bytes(b.get(p..p + 4).context("Truncated image")?.try_into()?) as usize)
    }
    fn be(b: &[u8], p: usize) -> Result<usize> {
        Ok(
            u32::from_be_bytes(b.get(p..p + 4).context("Truncated signature")?.try_into()?)
                as usize,
        )
    }
    ensure!(le(bytes, 0)? == 0xfeedfacf, "Unsupported signed image");
    let mut p = 32;
    for _ in 0..le(bytes, 16)? {
        let cmd = le(bytes, p)?;
        let size = le(bytes, p + 4)?;
        ensure!(size >= 8 && p + size <= bytes.len(), "Invalid load command");
        if cmd == 0x1d {
            let off = le(bytes, p + 8)?;
            let len = le(bytes, p + 12)?;
            let sig = bytes
                .get(off..off.checked_add(len).context("Signature overflow")?)
                .context("Missing signature")?;
            ensure!(
                be(sig, 0)? == 0xfade0cc0 && be(sig, 4)? <= sig.len(),
                "Invalid code signature"
            );
            let count = be(sig, 8)?;
            ensure!(count <= 128, "Too many signature blobs");
            let mut der = None;
            for i in 0..count {
                if be(sig, 12 + i * 8)? == 7 {
                    let at = be(sig, 16 + i * 8)?;
                    ensure!(be(sig, at)? == 0xfade7172, "Invalid DER entitlement blob");
                    let size = be(sig, at + 4)?;
                    ensure!(
                        (8..=1024 * 1024).contains(&size),
                        "Invalid DER entitlement size"
                    );
                    der = Some(
                        sig.get(at + 8..at + size)
                            .context("Truncated DER entitlements")?,
                    );
                }
                if be(sig, 12 + i * 8)? == 5 {
                    let at = be(sig, 16 + i * 8)?;
                    ensure!(be(sig, at)? == 0xfade7171, "Invalid entitlement blob");
                    let size = be(sig, at + 4)?;
                    ensure!(
                        (8..=1024 * 1024).contains(&size),
                        "Invalid entitlement size"
                    );
                    let xml = sig
                        .get(at + 8..at + size)
                        .context("Truncated entitlements")?;
                    return Ok(plist::Value::from_reader(std::io::Cursor::new(xml))?
                        .into_dictionary()
                        .context("Invalid entitlements")?);
                }
            }
            if let Some(der) = der {
                return super::entitlement_der::decode(der);
            }
            return Ok(plist::Dictionary::new());
        }
        p += size;
    }
    Ok(plist::Dictionary::new())
}
pub(super) fn permits(grant: &plist::Value, request: &plist::Value) -> bool {
    match (grant, request) {
        (plist::Value::String(g), plist::Value::String(r)) => {
            g == r || g.strip_suffix('*').is_some_and(|p| r.starts_with(p))
        }
        (plist::Value::Array(gs), plist::Value::Array(rs)) => {
            rs.iter().all(|r| gs.iter().any(|g| permits(g, r)))
        }
        _ => grant == request,
    }
}
#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn wildcard_grants_are_scoped() {
        assert!(permits(
            &plist::Value::String("TEAM.com.test.*".into()),
            &plist::Value::String("TEAM.com.test.app".into())
        ));
        assert!(!permits(
            &plist::Value::String("TEAM.com.test.*".into()),
            &plist::Value::String("TEAM.com.other.app".into())
        ));
    }
    #[test]
    fn helper_identifier_is_automatic_and_respects_profile_scope() {
        let own = "com.certlium.crumb";
        for (pattern, expected) in [
            ("*", "com.certlium.crumb.decrypt"),
            ("org.example.*", "org.example.crumbdecrypt"),
            ("org.example.helper", "org.example.helper"),
        ] {
            let mut ent = plist::Dictionary::new();
            ent.insert(
                "application-identifier".into(),
                format!("TEAM.{pattern}").into(),
            );
            let mut profile = plist::Dictionary::new();
            profile.insert("Entitlements".into(), ent.into());
            profile.insert(
                "ApplicationIdentifierPrefix".into(),
                vec![plist::Value::String("TEAM".into())].into(),
            );
            assert_eq!(choose_helper(&profile, own).unwrap(), expected);
        }
    }
    #[test]
    fn malformed_signature_fails() {
        assert!(entitlements(b"truncated").is_err());
    }
}
