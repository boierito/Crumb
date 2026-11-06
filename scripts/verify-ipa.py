#!/usr/bin/env python3
"""Validate the produced artifact, not only xcodebuild's exit status."""
import hashlib
import json
import io
import pathlib
import plistlib
import re
import struct
import sys
import zipfile

ipa = pathlib.Path(sys.argv[1])
project = pathlib.Path("WaffleStore.xcodeproj/project.pbxproj").read_text()
expected_build = re.search(r"CURRENT_PROJECT_VERSION = ([^;]+);", project)[1].strip()
expected_version = re.search(r"MARKETING_VERSION = ([^;]+);", project)[1].strip()
with zipfile.ZipFile(ipa) as archive:
    assert archive.testzip() is None, "IPA contains a bad ZIP checksum"
    root = "Payload/Crumb.app/"
    info = plistlib.loads(archive.read(root + "Info.plist"))
    assert info["CFBundleVersion"] == expected_build, "Wrong build number"
    assert info["CFBundleShortVersionString"] == expected_version, "Wrong version"
    assert info.get("UIFileSharingEnabled") is True, "Files sharing missing"
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
    assert "_remotepairing-pairable-host._tcp" in info.get("NSBonjourServices", []), "Pairing service declaration missing"
    assert "_remotepairing._tcp" in info.get("NSBonjourServices", []), "Pairing port discovery declaration missing"
    assert info.get("NSLocalNetworkUsageDescription"), "Local network permission text missing"
    with zipfile.ZipFile(io.BytesIO(archive.read(root + "CrumbDecryptTemplate.zip"))) as helper:
        assert set(helper.namelist()) == {"CrumbDecrypt", "Info.plist", "Icon60@2x.png", "Icon60@3x.png", "Icon76@2x.png", "Icon83.5@2x.png"}, "Unexpected helper template contents"
        helper_info = plistlib.loads(helper.read("Info.plist"))
        assert helper_info["CFBundleIcons"]["CFBundlePrimaryIcon"]["CFBundleIconFiles"] == ["Icon60"], "Helper icon missing"
        assert helper_info["CFBundleExecutable"] == "CrumbDecrypt", "Wrong helper executable"
        assert helper_info["UIFileSharingEnabled"] is True, "Helper retrieval unavailable"
        assert struct.unpack("<II", helper.read("CrumbDecrypt")[:8]) == (0xFEEDFACF, 0x0100000C), "Helper must be arm64"
    notices = [n for n in archive.namelist() if n.endswith(("ThirdPartyNotices.txt", "UNICORN-COPYING", "GoThirdPartyNotices.txt"))]
    assert len(notices) == 3, "Linked code license notices missing"
    assert b"MIT License" in archive.read(root + "ThirdPartyNotices.txt"), "ipatool MIT notice missing"
    print(json.dumps({"ipa": str(ipa), "build": expected_build, "version": expected_version,
                      "bundleID": info["CFBundleIdentifier"], "displayName": info["CFBundleDisplayName"],
                      "minimumOS": info["MinimumOSVersion"], "cpu": "arm64",
                      "FilesEnabled": True, "CrumbURLScheme": True, "DeviceHelperTemplate": True,
                      "sha256": hashlib.file_digest(ipa.open("rb"), "sha256").hexdigest()}, indent=2))
