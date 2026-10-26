#!/usr/bin/env python3
"""Adapt the hash-pinned semantic engine and profile compiler; no address patches."""
import hashlib,json,pathlib,re,sys
ROOT=pathlib.Path(__file__).resolve().parent.parent
VENDOR=ROOT/'third_party/unfairplay-generic'
manifest=json.loads((VENDOR/'SOURCE.json').read_text())
for name,digest in manifest['files'].items():
 if hashlib.sha256((VENDOR/pathlib.Path(name).name).read_bytes()).hexdigest()!=digest:
  raise SystemExit('generic interpreter source integrity mismatch: '+name)
mode,target=sys.argv[1],pathlib.Path(sys.argv[2]);target.mkdir(parents=True,exist_ok=True)
kernel=(VENDOR/'compact_kernel.go').read_text();cfg=(VENDOR/'compact_cfg.go').read_text()
abi=hashlib.sha256((kernel+cfg).encode()).hexdigest()
if mode=='compiler':
 source=(VENDOR/'main.go').read_text()
 source=source[:source.index('func annotatedMachoBlocks(')]+source[source.index('func compactLogicalMaskFields('):]
 for line in ['\t"flag"\n','\t"os"\n','\t"sort"\n','\t"github.com/DrWhax/unfairplay/internal/fairplay"\n']:
  source=source.replace(line,'')
 (target/'upstream.go').write_text(source)
 names=re.findall(r'^\s*(compact\w+)',kernel[kernel.index('const ('):kernel.index('\n\tcompactOpcodeCount')],re.M)
 assert names[0]=='compactAnd32' and len(names)<65536
 (target/'abi.go').write_text('package main\nvar opcodeNames=[]string{'+','.join(json.dumps(n) for n in names)+'}\nconst kernelABI='+json.dumps(abi)+'\n')
 (target/'main.go').write_bytes((ROOT/'scripts/fairplay-generic-profile/main.go').read_bytes())
 (target/'go.mod').write_text('module crumb/generic-profile\n\ngo 1.26.0\n\nrequire golang.org/x/arch v0.32.0\n')
elif mode=='runtime':
 # Profile kernels never use upstream phase-specific hooks or preamble examples.
 kernel=kernel[:kernel.index('\nvar compact26A428DigestPrologue')]
 start=kernel.index('func runCompactHook(');end=kernel.index('func runCompactKernel(',start)
 kernel=kernel[:start]+'''func runCompactHook(s *state, action compactHookAction) error { return fmt.Errorf("phase hooks forbidden in generic profile") }

'''+kernel[end:]
 cfg=cfg[:cfg.index('\nfunc compact26A428SINFEntryBlock(')]
 cfg=cfg.replace('s.x[30] = terminalPC + 4\n\t\treturn s.x[block.register], false, nil','target := s.x[block.register]\n\t\ts.x[30] = terminalPC + 4\n\t\treturn target, false, nil')
 for name,source in [('generic_compact_kernel.go',kernel),('generic_compact_cfg.go',cfg)]:
  (target/name).write_text(source.replace('package recovery','package rawsource',1))
 (target/'generic_abi.go').write_text('package rawsource\nconst genericKernelABI="'+abi+'"\n')
else: raise SystemExit('expected compiler or runtime mode')
