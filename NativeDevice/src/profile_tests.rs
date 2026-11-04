use super::*;
const OWN: &str = "com.certlium.crumb";
const HELPER: &str = "com.certlium.crumb.decrypt";
fn profile(appid: &str) -> plist::Dictionary {
    let mut ent = plist::Dictionary::new();
    ent.insert("application-identifier".into(), appid.into());
    let mut profile = plist::Dictionary::new();
    profile.insert("Entitlements".into(), ent.into());
    profile.insert(
        "TeamIdentifier".into(),
        vec![plist::Value::from("NEWTEAM123")].into(),
    );
    profile.insert(
        "ApplicationIdentifierPrefix".into(),
        vec![plist::Value::from("OLDPREFIX1")].into(),
    );
    profile
}
#[test]
fn legacy_app_prefix_is_preserved_independently_of_signing_team() {
    let profile = profile(&format!("OLDPREFIX1.{HELPER}"));
    assert_eq!(
        helper_application_id(&profile, HELPER, OWN).unwrap(),
        format!("OLDPREFIX1.{HELPER}")
    );
    assert_eq!(
        profile["TeamIdentifier"].as_array().unwrap()[0].as_string(),
        Some("NEWTEAM123")
    );
}
#[test]
fn wildcard_profiles_preserve_prefix_and_restrict_bundle_scope() {
    let scoped = profile("OLDPREFIX1.com.certlium.*");
    assert_eq!(
        helper_application_id(&scoped, HELPER, OWN).unwrap(),
        format!("OLDPREFIX1.{HELPER}")
    );
    assert!(helper_application_id(&scoped, "com.certliumevil.decrypt", OWN).is_err());
    let any = profile("OLDPREFIX1.*");
    assert_eq!(
        helper_application_id(&any, "org.example.decrypt", OWN).unwrap(),
        "OLDPREFIX1.org.example.decrypt"
    );
}
#[test]
fn mismatched_explicit_profile_identifies_allowed_and_requested_bundles() {
    let profile = profile("OLDPREFIX1.org.example.otherhelper");
    let message = helper_application_id(&profile, HELPER, OWN)
        .unwrap_err()
        .to_string();
    assert!(message.contains("Profile permits org.example.otherhelper"));
    assert!(message.contains(HELPER));
    assert!(message.contains("Import a matching profile"));
}
#[test]
fn crumb_only_profile_explains_the_need_for_a_separate_helper() {
    let profile = profile(&format!("OLDPREFIX1.{OWN}"));
    let message = helper_application_id(&profile, HELPER, OWN)
        .unwrap_err()
        .to_string();
    assert!(message.contains("This profile is for Crumb itself"));
    assert!(
        helper_application_id(&profile, OWN, OWN)
            .unwrap_err()
            .to_string()
            .contains("must not replace Crumb")
    );
}
#[test]
fn undeclared_prefixes_and_malformed_patterns_are_rejected() {
    let wrong_prefix = profile("NEWTEAM123.*");
    assert!(
        helper_application_id(&wrong_prefix, HELPER, OWN)
            .unwrap_err()
            .to_string()
            .contains("ApplicationIdentifierPrefix")
    );
    for appid in [
        "OLDPREFIX1.com.*.decrypt",
        "OLDPREFIX1.**",
        "OLDPREFIX1.",
        "OLDPREFIX1",
    ] {
        assert!(helper_application_id(&profile(appid), HELPER, OWN).is_err());
    }
}
