#!/usr/bin/env python3
"""Validate the produced artifact, not only xcodebuild's exit status."""
import hashlib
import json
import pathlib
import plistlib
import re
import struct
import sys
import zipfile

ipa = pathlib.Path(sys.argv[1])
project = pathlib.Path("Crumb.xcodeproj/project.pbxproj").read_text()
expected_build = re.search(r"CURRENT_PROJECT_VERSION = ([^;]+);", project)[1].strip()
expected_version = re.search(r"MARKETING_VERSION = ([^;]+);", project)[1].strip()
with zipfile.ZipFile(ipa) as archive:
    assert archive.testzip() is None, "IPA contains a bad ZIP checksum"
    root = "Payload/Crumb.app/"
    info = plistlib.loads(archive.read(root + "Info.plist"))
    assert info["CFBundleVersion"] == expected_build, "Wrong build number"
    assert info["CFBundleShortVersionString"] == expected_version, "Wrong version"
    assert info.get("UIFileSharingEnabled") is False, "Files sharing must be disabled to protect account-bound encrypted IPAs"
    assert info.get("LSSupportsOpeningDocumentsInPlace") is True, "Files in-place access missing"
    assert any("crumb" in item.get("CFBundleURLSchemes", [])
               for item in info.get("CFBundleURLTypes", [])), "Crumb URL scheme missing"
    assert info["CFBundleIdentifier"] == "com.certlium.crumb", "Wrong Crumb bundle ID"
    assert info["CFBundleDisplayName"] == "Crumb", "Wrong display name"
    assert info.get("CFBundleIcons", {}).get("CFBundlePrimaryIcon", {}).get("CFBundleIconFiles"), "App icon missing"
    executable = archive.read(root + info["CFBundleExecutable"])
    assert struct.unpack("<II", executable[:8]) == (0xFEEDFACF, 0x0100000C), "Expected arm64 Mach-O"
    picker = archive.read(root + "Frameworks/FilePickerFix.dylib")
    assert hashlib.sha256(picker).digest() == hashlib.sha256(pathlib.Path("EmbeddedLibraries/FilePickerFix.dylib").read_bytes()).digest(), "Picker fix changed"
    assert b"@executable_path/Frameworks/FilePickerFix.dylib\0" in executable, "Picker fix load command missing"
    assert not any(n.endswith("embedded.mobileprovision") for n in archive.namelist()), "IPA must be resignable"
    assert root + "CrumbDecryptTemplate.zip" not in archive.namelist(), "Legacy helper must not be bundled"
    notices = [n for n in archive.namelist() if n.endswith(("ThirdPartyNotices.txt", "GoThirdPartyNotices.txt"))]
    assert len(notices) == 2, "Linked code license notices missing"
    assert b"MIT License" in archive.read(root + "ThirdPartyNotices.txt"), "ipatool MIT notice missing"
    print(json.dumps({"ipa": str(ipa), "build": expected_build, "version": expected_version,
                      "bundleID": info["CFBundleIdentifier"], "displayName": info["CFBundleDisplayName"],
                      "minimumOS": info["MinimumOSVersion"], "cpu": "arm64",
                      "FilesEnabled": False, "CrumbURLScheme": True, "DeviceHelperTemplate": False,
                      "sha256": hashlib.file_digest(ipa.open("rb"), "sha256").hexdigest()}, indent=2))
