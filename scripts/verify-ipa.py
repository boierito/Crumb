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
project = pathlib.Path("WaffleStore.xcodeproj/project.pbxproj").read_text()
expected_build = re.search(r"CURRENT_PROJECT_VERSION = ([^;]+);", project)[1].strip()
expected_version = re.search(r"MARKETING_VERSION = ([^;]+);", project)[1].strip()
with zipfile.ZipFile(ipa) as archive:
    assert archive.testzip() is None, "IPA contains a bad ZIP checksum"
    root = "Payload/WaffleStore.app/"
    info = plistlib.loads(archive.read(root + "Info.plist"))
    assert info["CFBundleVersion"] == expected_build, "Wrong build number"
    assert info["CFBundleShortVersionString"] == expected_version, "Wrong version"
    assert info.get("UIFileSharingEnabled") is True, "Files sharing missing"
    assert info.get("LSSupportsOpeningDocumentsInPlace") is True, "Files in-place access missing"
    assert any("wafflestore" in item.get("CFBundleURLSchemes", [])
               for item in info.get("CFBundleURLTypes", [])), "Original URL scheme missing"
    executable = archive.read(root + info["CFBundleExecutable"])
    assert struct.unpack("<II", executable[:8]) == (0xFEEDFACF, 0x0100000C), "Expected arm64 Mach-O"
    assert not any(n.endswith("embedded.mobileprovision") for n in archive.namelist()), "IPA must be resignable"
    notices = [n for n in archive.namelist() if n.endswith(("ThirdPartyNotices.txt", "UNICORN-COPYING", "GoThirdPartyNotices.txt"))]
    assert len(notices) == 3, "Linked code license notices missing"
    assert b"MIT License" in archive.read(root + "ThirdPartyNotices.txt"), "ipatool MIT notice missing"
    print(json.dumps({"ipa": str(ipa), "build": expected_build, "version": expected_version,
                      "minimumOS": info["MinimumOSVersion"], "cpu": "arm64",
                      "FilesEnabled": True, "originalURLScheme": True,
                      "sha256": hashlib.file_digest(ipa.open("rb"), "sha256").hexdigest()}, indent=2))
