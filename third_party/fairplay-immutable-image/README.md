# Exact immutable image data

These parts contain only the `__TEXT` segment extracted from the public
`fairplayd` image with SHA-256
`1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784`.
They contain no account enrollment, device snapshot, SINF, executable from an
IPA, or derived account key. Binary-derived Apple data retains its provenance;
this directory does not relicense it as the MIT-licensed recovery sources.

Generation: `scripts/audit-fairplay-model-memory.py`, commit
`8e780304feed2ddc9b713005dfc72684cabf98e3`, Actions run
https://github.com/boierito/Crumb/actions/runs/40493389169.
The extractor verifies the complete thin executable hash, segment permissions
and geometry, and absence of chained fixups in `__TEXT` before exporting bytes.
`IMAGE.json` records both image provenance and separate resource hashes.

The 16,711,680-byte segment includes 4,080 model-sized pages. The pinned public
snapshot already contains 227 of them, including six that differ from the file.
Runtime installation preserves all captured pages and copies only the missing
3,853 pages. Writable image segments, heap and stack are never supplied by this
resource. These are Go byte slices used as interpreter data; they are not native
executable mappings. Generation uses a development host; use on iPhone requires
only the bundled resource.
