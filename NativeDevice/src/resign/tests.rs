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
