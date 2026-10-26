#!/usr/bin/env python3
"""Fail closed before an iPhone FairPlay IPA build with no new code coverage.

This checks a necessary release prerequisite, NOT successful decryption.
No keybags, SIDB data, Apple IDs, executable binaries or private captures.
"""
import argparse
import json
from pathlib import Path
import re
import sys

DEFAULT_TARGET = 0x100A4BE38


def coverage_report(root: Path, target: int) -> dict:
    native = root / "MapleSyrup" / "NativeSAP"
    processor = (root / "scripts" / "prepare-fairplay-raw.py").read_text()
    implementation = "\n".join(
        p.read_text() for p in sorted(native.glob("fairplayrawsupplement*.go"))
        if not p.name.endswith("_test.go")
    )
    tests = "\n".join(
        p.read_text() for p in sorted(native.glob("fairplayrawsupplement*_test.go"))
    )
    address = f"{target:x}"
    function = f"rawSourceBootstrapBlock{address}"
    exists = bool(re.search(
        r"(?m)^func\s+" + re.escape(function) + r"\s*\(", implementation
    ))
    entry = bool(re.search(
        r"(?m)^blocks\[\s*(?:0x" + address + r"|supplement_entry)\s*\]"
        r"\s*=\s*'" + re.escape(function) + r"'", processor
    )) and (f"0x{address}" in processor)
    test = bool(re.search(
        r"func\s+Test\w*" + address + r"\w*\s*\(", tests
    ))
    return {
        "schema": "crumb-fairplay-release-coverage-gate-v1",
        "required_pc": f"text+0x{target - 0x100000000:x}",
        "aot_handler_present": exists,
        "aot_handler_registered": entry,
        "target_specific_synthetic_test_present": test,
        "ready_for_device_smoke_test": exists and entry and test,
        "decryption_validated": False,
    }


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", default=str(Path(__file__).resolve().parent.parent))
    parser.add_argument("--target", type=lambda value: int(value, 0),
                        default=DEFAULT_TARGET)
    parser.add_argument("--report", default="")
    args = parser.parse_args()
    result = coverage_report(Path(args.root), args.target)
    payload = json.dumps(result, indent=2, sort_keys=True) + "\n"
    print(payload, end="")
    if args.report:
        dest = Path(args.report)
        dest.parent.mkdir(parents=True, exist_ok=True)
        dest.write_text(payload)
    if not result["ready_for_device_smoke_test"]:
        print("Release blocked: new target lacks native coverage and a targeted test.",
              file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
