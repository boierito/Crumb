#!/usr/bin/env bash
# Development reference acquisition only. No Crumb account/license inputs.
set -euo pipefail

research_dir="${1:?provide an empty research output directory}"
if [[ -e "$research_dir" ]]; then
  echo 'Research directory must not already exist' >&2
  exit 1
fi
mkdir -p "$research_dir"/{tool,input,plain,mount,artifacts/macho}
research_dir="$(cd "$research_dir" && pwd)"
tool_archive="$research_dir/tool/ipsw.tar.gz"
curl --fail --location --retry 3 \
  https://github.com/blacktop/ipsw/releases/download/v3.1.739/ipsw_3.1.739_macOS_universal.tar.gz \
  --output "$tool_archive"
printf '%s  %s\n' '4a63c4a8a50f59c1d6255c7e958e1de3f1a3c3e9bd8a4e5b68d339a8b780cd74' "$tool_archive" | shasum -a 256 -c -
tar -xzf "$tool_archive" -C "$research_dir/tool"
ipsw_tool="$research_dir/tool/ipsw"
"$ipsw_tool" version > "$research_dir/artifacts/tool-version.txt"

firmware_url='https://updates.cdn-apple.com/2026FallFCS/59ffd938-75b6-41e5-af04-f4295ee8aeb3/iPhone17,1_27.0.1_24A446_Restore.ipsw'
"$ipsw_tool" --no-color extract --remote --flat \
  --pattern '^043-69419-713\.dmg\.aea$' \
  --output "$research_dir/input" "$firmware_url"
archive_path="$(python3 - "$research_dir/input" <<'PY'
import pathlib, sys
files = list(pathlib.Path(sys.argv[1]).rglob('043-69419-713.dmg.aea'))
if len(files) != 1 or files[0].stat().st_size != 2348810240:
    raise SystemExit('Expected exactly the pinned 24A446 Cryptex component')
print(files[0])
PY
)"
shasum -a 256 "$archive_path" > "$research_dir/artifacts/public-component-sha256.txt"
printf '%s  %s\n' '23bd07e51c154ea4c1afd5f59c24e05ddeab2b6d0cbd380d4623c13fb5a93091' "$archive_path" | shasum -a 256 -c -
"$ipsw_tool" --no-color fw aea "$archive_path" --output "$research_dir/plain"
plain_path="$(python3 - "$research_dir/plain" <<'PY'
import pathlib, sys
files = list(pathlib.Path(sys.argv[1]).rglob('043-69419-713.dmg'))
if len(files) != 1:
    raise SystemExit('Expected exactly one decoded public filesystem image')
print(files[0])
PY
)"
mounted=false
cleanup_mount() {
  if [[ "$mounted" == true ]]; then
    hdiutil detach "$research_dir/mount" || true
  fi
}
trap cleanup_mount EXIT
hdiutil attach -readonly -nobrowse -mountpoint "$research_dir/mount" "$plain_path"
mounted=true
cache_path="$(python3 - "$research_dir/mount" <<'PY'
import pathlib, sys
import json, re
root = pathlib.Path(sys.argv[1])
candidates = sorted(p for p in root.rglob('dyld_shared_cache*') if p.is_file())
print(json.dumps({'cache_candidates': [str(p.relative_to(root)) for p in candidates]}, indent=2), file=sys.stderr)
system_cache_directory = root / 'System/Library/Caches/com.apple.dyld'
files = [p for p in candidates if p.parent == system_cache_directory
         and re.fullmatch(r'dyld_shared_cache_arm64e(?:\.x[0-9]+)?', p.name)]
if len(files) != 1:
    raise SystemExit('Expected one unambiguous arm64e / arm64e.xN shared-cache root')
print(files[0])
PY
)"
"$ipsw_tool" --no-color dyld info --dylibs "$cache_path" > "$research_dir/artifacts/cache-images.txt"
for framework in CoreML MediaToolbox StoreServices; do
  case "$framework" in
    CoreML|MediaToolbox) image="/System/Library/Frameworks/$framework.framework/$framework" ;;
    StoreServices) image="/System/Library/PrivateFrameworks/$framework.framework/$framework" ;;
  esac
  if ! python3 - "$research_dir/artifacts/cache-images.txt" "$image" <<'PY'
import pathlib, sys
raise SystemExit(0 if sys.argv[2] in pathlib.Path(sys.argv[1]).read_text() else 1)
PY
  then
    echo "Reference image not listed: $framework" >> "$research_dir/artifacts/absent-images.txt"
    continue
  fi
  # Local symbols are exported by default. Global stub-island analysis builds
  # an index of the entire OS cache and is unnecessary for this acquisition.
  "$ipsw_tool" --no-color dyld extract --slide \
    --output "$research_dir/artifacts/macho" "$cache_path" "$image"
done
python3 "$(dirname "$0")/collect-fairplay-client-stubs.py" \
  --ipsw "$ipsw_tool" --cache "$cache_path" \
  --output "$research_dir/artifacts/client-stubs.json"
python3 - "$research_dir/artifacts" <<'PY'
import hashlib, json, pathlib, struct, sys
root = pathlib.Path(sys.argv[1])
images = []
for path in sorted((root / 'macho').rglob('*')):
    if not path.is_file():
        continue
    data = path.read_bytes()
    if data[:4] != struct.pack('<I', 0xfeedfacf):
        continue
    images.append({'file': str(path.relative_to(root)), 'bytes': len(data),
                   'sha256': hashlib.sha256(data).hexdigest()})
if not any(pathlib.PurePosixPath(i['file']).name == 'CoreML' for i in images):
    raise SystemExit('The required CoreML reference was not acquired')
if sum(i['bytes'] for i in images) > 512 * 1024 * 1024:
    raise SystemExit('Reference artifact size exceeds its bounded scope')
report = {'schema': 'crumb-public-client-reference-v1', 'build': '24A446',
          'component': '043-69419-713.dmg.aea', 'ipsw_version': '3.1.739',
          'extraction': 'original cache image with ipsw slide normalization and local symbols',
          'account_inputs_used': False, 'license_inputs_used': False,
          'device_service_executed': False, 'sinf_validated': False,
          'crumb_runtime_dependency': False, 'images': images}
(root / 'acquisition.json').write_text(json.dumps(report, indent=2) + '\n')
print(json.dumps(report, indent=2))
PY
