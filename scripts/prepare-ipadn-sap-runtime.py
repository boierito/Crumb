#!/usr/bin/env python3
"""Assemble iPadN's complete, pinned SAP runtime as a private iOS bundle resource.

Do not use floating branches: every fetched byte is checked against iPadN's
published vendored manifest. Network access is required during the CI build,
not during Apple login. This archive contains third-party runtime files with
incomplete redistribution provenance: it is for an experimental test build.
"""
import hashlib
import io
import json
import pathlib
import urllib.request
import zipfile

ROOT = pathlib.Path(__file__).resolve().parents[1]
REVISION = "70f2ac36e7f87cacf82a45b58e6fa27ed23ecb80"
BASE = f"https://raw.githubusercontent.com/mengangs/ipadn/{REVISION}/iPadn-iOS/vendor/SAPRuntime/files/"
MANIFEST = ROOT / "scripts/ipadn-sap-manifest.json"
DESTINATION = ROOT / "Crumb/SAPRuntime.zip"


def prepare() -> None:
    expected = json.loads(MANIFEST.read_text(encoding="utf8"))["files"]
    if len(expected) != 9:
        raise RuntimeError("iPadN runtime manifest does not describe exactly nine files")
    DESTINATION.parent.mkdir(parents=True, exist_ok=True)
    buffer = io.BytesIO()
    with zipfile.ZipFile(buffer, "w", compression=zipfile.ZIP_DEFLATED, compresslevel=6) as archive:
        for filename, sha in sorted(expected.items()):
            if pathlib.PurePosixPath(filename).name != filename:
                raise RuntimeError("Unsafe iPadN runtime file name")
            request = urllib.request.Request(BASE + filename, headers={"User-Agent": "Crumb-CI/1"})
            with urllib.request.urlopen(request, timeout=120) as response:
                data = response.read(65 * 1024 * 1024)
            if hashlib.sha256(data).hexdigest() != sha:
                normalized = data.replace(b"\r\n", b"\n")
                if hashlib.sha256(normalized).hexdigest() != sha:
                    raise RuntimeError("iPadN SAP runtime hash mismatch: " + filename)
                data = normalized
            archive.writestr(filename, data)
            print(f"iPadN SAP verified {filename}: {len(data)} bytes")
    temporary = DESTINATION.with_suffix(".tmp")
    temporary.write_bytes(buffer.getvalue())
    temporary.replace(DESTINATION)
    print(f"iPadN SAPRuntime.zip assembled: {DESTINATION.stat().st_size} bytes")


if __name__ == "__main__":
    prepare()
