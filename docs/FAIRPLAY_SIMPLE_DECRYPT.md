# Simplified decryption, build 23056

The purchased-IPA flow now offers Save IPA, Install, and Decrypt. Decrypt &
install is removed from this interface: the unsigned verified output is the
only decryption result offered. The separate signing/export engine is preserved
internally, without placing configuration or debug controls in the product flow.

The decryption sheet shows progress, a save action on completion, or one retry
action on failure. Preparation, signing the private Crumb update, selecting the
native reader and export checks are automatic. A new download cannot inherit a
previous target's output or queued preparation. Automatic resume requires the
retained job to match the installed job, and an explicitly selected record must
match the saved request as well.

Certificate configuration contains a P12 picker, a provisioning profile picker,
a password field and Save/Continue. Validation runs as part of that action;
there is no manual preflight operation, role selection, optional output profile,
bundle ID field, report viewer or preparation/decryption stage menu. File opening
through Share → Crumb remains available. Only one Certificate row remains in the
ordinary Settings section. Lab entry points are removed from the visible
Settings/Downloads menus in both Debug and Release; research implementation is
preserved.

This does NOT implement a new native FairPlay capability. The demonstrated
system path still needs the target inside the installed Crumb bundle, which
requires an accepted signed update. iOS's installation confirmation and reopening
Crumb remain necessary. When configuration is needed, that requirement is stated
in the certificate sheet; while waiting for installation it is stated in the
progress sheet. Merely hiding UI cannot satisfy a fully unattended one-tap flow.
The current route is therefore not represented as meeting that final goal.

Original all-page verification, privacy limits and OS installation checks remain
unchanged. No new plaintext authenticity or universal anonymization claim is made.
