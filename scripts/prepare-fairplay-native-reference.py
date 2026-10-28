#!/usr/bin/env python3
"""Instrument the prepared host runtime only; never alter the iOS build."""
from pathlib import Path
import sys

root = Path(__file__).resolve().parent.parent
target = Path(sys.argv[1])
flow = target / "native_flow.go"
text = flow.read_text()
old = "func runNativeRegion(s *state, runtime *rawSourceRuntime, pc uint64, region nativeRegion) (uint64, error) {"
assert text.count(old) == 1, "native region signature changed"
text = text.replace(
    old, old.replace("(uint64, error)", "(finalPC uint64, finalErr error)"), 1
)
old = "\thits := 0\n"
assert text.count(old) == 1, "native region entry changed"
text = text.replace(
    old,
    "\tnativeReferenceStart(s, pc)\n\tdefer func() { nativeReferenceFinish(s, finalPC, finalErr) }()\n"
    + old,
    1,
)
old = "\t\t\tnext, err := runtime.handle(s, pc)"
assert text.count(old) == 1, "native service dispatch changed"
text = text.replace(
    old,
    "\t\t\tnativeReferenceBefore(s, pc)\n"
    + old
    + "\n\t\t\tnativeReferenceAfter(s, next, err)",
    1,
)
output = target / "source_outputs.go"
stores = output.read_text()
old = "func sourceOutputStore(s *state, address uint64, data []byte) error {"
assert stores.count(old) == 1, "native store contract changed"
stores = stores.replace(old, old + "\n\tnativeReferenceStore(s, address, len(data))", 1)
# Validate every anchor before making any edits.
flow.write_text(text)
output.write_text(stores)
(target / "native_reference.go").write_bytes(
    (root / "scripts/fairplay-native-reference/runtime.go").read_bytes()
)
