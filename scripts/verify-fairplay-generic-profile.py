#!/usr/bin/env python3
"""Generic-runtime release gate. No per-address handler requirement."""
from pathlib import Path
import json,sys
root=Path(__file__).resolve().parent.parent
profile=root/'third_party/unfairplay-generic/fairplay-generic.ir.gz'
manifest=root/'third_party/unfairplay-generic/PROFILE.json'
ready=profile.is_file() and manifest.is_file()
print(json.dumps({'schema':'crumb-generic-release-v1','profile_present':ready,'ready_for_device_test':False,'decryption_validated':False},indent=2))
if not ready:
 print('Generic profile/integration not validated; IPA build blocked.',file=sys.stderr)
 sys.exit(1)
sys.exit('Generic integration checks must be completed before release.')
