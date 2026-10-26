#!/usr/bin/env python3
"""Inspect a hosted macOS ARM runner for the EXACT pinned FairPlay executable.

No login, private Apple state, device, root, jailbreak, JIT or uploaded binary.
Produces only bounded public metadata; raw Apple executables are never uploaded.
A hash mismatch is an explicit stop, not permission to invent branch code.
"""
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile

EXPECTED = "1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784"
TARGET = 0x1008fff2c
CANDIDATES = (
    "/usr/libexec/fairplayd",
    "/System/Library/PrivateFrameworks/CoreFP.framework/Versions/A/fairplayd",
    "/System/Library/PrivateFrameworks/CoreFP.framework/fairplayd",
    "/System/Volumes/Preboot/Cryptexes/OS/System/Library/PrivateFrameworks/CoreFP.framework/Versions/A/fairplayd",
    "/System/Library/PrivateFrameworks/FairPlay.framework/Versions/A/fairplayd",
    "/System/Library/PrivateFrameworks/FairPlay.framework/fairplayd",
)

def command(args, timeout=10):
    try:
        p = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                           text=True, timeout=timeout, check=False)
        return p.returncode, p.stdout[:8192]
    except (OSError, subprocess.TimeoutExpired):
        return -1, ""

def main():
    report = {
        "schema": "crumb-fairplay-host-coverage-v1",
        "target": "text+0x8fff2c",
        "pinned_binary_sha256": EXPECTED,
        "host_macos": command(["sw_vers", "-productVersion"])[1].strip(),
        "host_macos_build": command(["sw_vers", "-buildVersion"])[1].strip(),
        "jit_used": False,
        "jailbreak_used": False,
        "private_material_accessed": False,
        "candidate_count": 0,
        "exact_build_found": False,
        "status": "no_matching_public_system_executable",
    }
    seen = set()
    candidates = [Path(p) for p in CANDIDATES]
    code, found = command(["mdfind", "kMDItemFSName == 'fairplayd'"], timeout=12)
    if code == 0:
        candidates += [Path(s) for s in found.splitlines()[:30] if s.startswith("/")]
    with tempfile.TemporaryDirectory(prefix="crumb-public-profile-") as temp:
        for candidate in candidates:
            if str(candidate) in seen:
                continue
            seen.add(str(candidate))
            try:
                if not candidate.is_file() or candidate.is_symlink() or candidate.stat().st_size > 128<<20:
                    continue
            except OSError:
                continue
            report["candidate_count"] += 1
            exact = None
            paths = [candidate]
            # The pinned model is arm64e; check architecture slices without
            # assuming a universal container has the same whole-file SHA.
            for arch in ("arm64e", "arm64"):
                thin = Path(temp)/("fairplayd-" + arch)
                rc, _ = command(["lipo", "-thin", arch, str(candidate), "-output", str(thin)])
                if rc == 0 and thin.exists():
                    paths.append(thin)
            for path in paths:
                try:
                    digest = hashlib.sha256(path.read_bytes()).hexdigest()
                except OSError:
                    continue
                if digest == EXPECTED:
                    exact = path
                    break
            if exact is None:
                continue
            report["exact_build_found"] = True
            report["status"] = "exact_binary_located"
            report["source_type"] = "system_executable"
            # Never publish raw instructions from a mismatched build.
            tool = shutil.which("llvm-objdump")
            if tool is None:
                rc, output = command(["xcrun", "-f", "llvm-objdump"])
                if rc == 0:
                    tool = output.strip()
            if tool:
                rc, output = command([tool, "--disassemble", "--start-address="+hex(TARGET-64),
                                     "--stop-address="+hex(TARGET+192), str(exact)], timeout=25)
                if rc == 0:
                    rows = [ln.strip() for ln in output.splitlines()
                            if ln.strip().startswith("1008fff") or "1008fff" in ln][:55]
                    report["target_disassembly_available"] = any("1008fff2c" in row for row in rows)
                    report["target_disassembly_lines"] = len(rows)
                    # This bounded public instruction excerpt permits a
                    # deterministic handwritten Go translation. Never attach
                    # or export the Apple executable, model or private data.
                    target_index = next((i for i, row in enumerate(rows)
                                         if "1008fff2c" in row), -1)
                    if target_index >= 0:
                        report["target_code_excerpt"] = rows[target_index:target_index+35]
                else:
                    report["status"] = "exact_binary_disassembler_unavailable"
            break
    destination = Path(sys.argv[1]) if len(sys.argv)>1 else None
    encoded = json.dumps(report, indent=2, sort_keys=True)+"\n"
    print(encoded)
    if destination:
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_text(encoded)

if __name__ == "__main__":
    main()
