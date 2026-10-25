#!/usr/bin/env python3
"""Isolate the pinned pure-Go modern enrollment runtime (no guest execution)."""
import hashlib
import pathlib
import re
import sys

source, target = map(pathlib.Path, sys.argv[1:])
target.mkdir(parents=True, exist_ok=True)
names = ['oracle.go', 'main.go', 'dfa_aes.go', 'dfa_attack.go', 'dfa_finish.go',
         'dfa_linear.go', 'generated_26b118.go', 'input_size_default.go',
         'enrollment_frame_default.go', 'embedded_model_raw_source.go',
         'embedded_model_raw_source_test.go', 'raw_source_step.go',
         'raw_source_runtime_test.go', 'raw_source_vector_test.go']
names += [p.name for p in sorted(source.glob('generated_raw_source_*.go'))]
for name in names:
    text = (source / name).read_text()
    text = re.sub(r'^//go:build[^\n]*\n', '', text)
    text = text.replace('package main', 'package rawsource', 1)
    if name == 'raw_source_step.go':
        # Instrument the existing lookup only; preserve longest-suffix matching
        # and all guest-visible return values. Never retain requested paths.
        field = 'type rawSourceRuntime struct {'
        lookup = 'return bestSuffix, bestData, bestSuffix != ""'
        if text.count(field) != 1 or text.count(lookup) != 1:
            raise SystemExit('raw runtime lookup contract changed')
        text = text.replace(field, field + '\n lookupCounts map[string]LookupCount')
        text = text.replace(lookup, 'r.recordLookup(path, bestSuffix != "")\n\t' + lookup)
    (target / name).write_text(text)
profile = 'macos26a428-account-entry-public.model.gz'
(target / 'profiles').mkdir(exist_ok=True)
resource = (source / 'profiles' / profile).read_bytes()
(target / 'profiles' / profile).write_bytes(resource)
expected = hashlib.sha256(resource).hexdigest()
api = pathlib.Path(__file__).resolve().parent.parent / 'MapleSyrup/NativeSAP/fairplayrawapi.go'
text = api.read_text().replace('RAW_RESOURCE_SHA256', expected)
(target / 'api.go').write_text(text)

(target / "api_test.go").write_text((api.parent / "fairplayrawapi_test.go").read_text())

for name, origin in [('LICENSE-IPATOOL', source.parent.parent / 'LICENSE'),
                     ('LICENSE-DARKPHOENIX.txt', source / 'LICENSE-DARKPHOENIX.txt'),
                     ('PROFILE-README.md', source / 'profiles/README.md')]:
    (target / name).write_bytes(origin.read_bytes())
