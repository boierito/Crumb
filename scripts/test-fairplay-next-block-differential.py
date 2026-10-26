#!/usr/bin/env python3
"""Compare the integrated block against independently assembled ARM64 semantics.

Only public synthetic state is used. The final BRAA is replaced with BR,
matching the portable runtime's existing PAC convention; this does not validate
native pointer authentication or a real enrollment.
"""
import re, random, base64, json, struct, subprocess, hashlib, tarfile, tempfile
from pathlib import Path
from keystone import Ks, KS_ARCH_ARM64, KS_MODE_LITTLE_ENDIAN
from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM
from unicorn import arm64_const as ar
ROOT=Path(__file__).resolve().parent.parent
REV='f7242cce58077a4af8db738ed295ee23ad8f067d'
ARCHIVE_HASH='e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39'
HARNESS='package main\nimport("encoding/json";"os";"fmt")\ntype rawSourceRuntime struct{stopPC uint64}\ntype fixture struct {Registers [31]uint64; Pages map[uint64][]byte; Flags [4]bool}\ntype result struct {Registers [31]uint64; Flags [4]bool; Next uint64; Output []byte}\nfunc main(){\n var inputs []fixture; if err:=json.NewDecoder(os.Stdin).Decode(&inputs);err!=nil{panic(err)}\n results:=[]result{}\n for _,f:=range inputs{\n  m:=&model{registers:f.Registers,pages:f.Pages,n:f.Flags[0],z:f.Flags[1],c:f.Flags[2],v:f.Flags[3]}\n  s:=newState(m);s.reset()\n  next,stopped,err:=rawSourceBootstrapBlock100a4be38(s,nil);if err!=nil||stopped{panic(fmt.Sprint(err))}\n  output,err:=s.outputBytes(0x500000,65536);if err!=nil{panic(err)}\n  results=append(results,result{s.x,[4]bool{s.n,s.z,s.c,s.v},next,output})\n }\n json.NewEncoder(os.Stdout).Encode(results)\n}\n'
with tempfile.TemporaryDirectory(prefix='crumb-block-reference-') as directory:
 work=Path(directory)
 archive=ROOT/'third_party/archives'/f'ipatool-fairplay-{REV}.tar.gz'
 if hashlib.sha256(archive.read_bytes()).hexdigest()!=ARCHIVE_HASH:
  raise ValueError('source integrity mismatch')
 with tarfile.open(archive) as tar: tar.extractall(work,filter='data')
 upstream=work/f'ipatool-{REV}/tools/fairplay-recovery'
 for name in ('oracle.go','generated_26b118.go','input_size_default.go'):
  (work/name).write_bytes((upstream/name).read_bytes())
 (work/'main.go').write_text(HARNESS)
 (work/'candidate.go').write_text((ROOT/'MapleSyrup/NativeSAP/fairplayrawsupplement_next.go').read_text().replace('package rawsource','package main',1))
 runner=work/'reference-runner'
 subprocess.run(['go','build','-o',str(runner),'main.go','oracle.go','generated_26b118.go','input_size_default.go','candidate.go'],cwd=work,check=True,timeout=120)
 src=(ROOT / 'MapleSyrup/NativeSAP/fairplayrawsupplement_next.go').read_text()
 asm=[m.group(1) for m in re.finditer(r'// 0x[0-9a-f]+: (.+)',src)]
 assert asm[-1]=='braa x6, x17';asm[-1]='br x6'
 code=bytes(Ks(KS_ARCH_ARM64,KS_MODE_LITTLE_ENDIAN).asm('; '.join(asm))[0]);start=0x100a4be38
 rng=random.Random(23042);inputs=[];expected=[]
 regsids=[getattr(ar,'UC_ARM64_REG_X'+str(i)) for i in range(29)]+[ar.UC_ARM64_REG_X29,ar.UC_ARM64_REG_X30]
 for i in range(40):
  regs=[rng.getrandbits(64) for _ in range(31)]
  counter=[0,62,63,64,127,128,191,192,254,255][i%10]
  regs[0]=37;regs[1]=0x10383d;regs[2]=(rng.getrandbits(56)<<8)|counter;regs[3]=73;regs[4]=165;regs[5]=110;regs[9]=[0,0xffffffff][i%2]
  regs[24]=[0,1,15,16,255,0xffffffff][i%6];regs[10]=(i-regs[24])&0xffffffff
  regs[11]=13;regs[13]=0x300000;regs[14]=111;regs[15]=0x14307;regs[16]=176;regs[22]=0x500000;regs[23]=0x400000;regs[30]=0x600800
  pages={}
  for base in (0x300000,0x400000,0x500000):
   for offset in range(0,65536,4096):pages[base+offset]=rng.randbytes(4096)
  dispatch=bytearray(4096)
  for index,target in ((regs[9],0x200000),(110^regs[9],0x200004)):
   signed=index if index<0x80000000 else index-(1<<32)
   struct.pack_into('<Q',dispatch,0x800+signed*8,target)
  pages[0x600000]=bytes(dispatch)
  flags=[bool(rng.getrandbits(1)) for _ in range(4)]
  inputs.append({'Registers':regs,'Pages':{k:base64.b64encode(v).decode() for k,v in pages.items()},'Flags':flags})
  em=Uc(UC_ARCH_ARM64,UC_MODE_ARM);em.mem_map(start&~4095,4096);em.mem_write(start,code);em.mem_map(0x200000,4096)
  for addr,data in pages.items():em.mem_map(addr,4096);em.mem_write(addr,data)
  for reg,value in zip(regsids,regs):em.reg_write(reg,value)
  em.reg_write(ar.UC_ARM64_REG_NZCV,sum(int(b)<<(31-j) for j,b in enumerate(flags)))
  em.emu_start(start,start+len(code),count=len(asm))
  nzcv=em.reg_read(ar.UC_ARM64_REG_NZCV)
  expected.append({'Registers':[em.reg_read(r) for r in regsids],'Flags':[bool(nzcv&(1<<(31-j))) for j in range(4)],'Next':em.reg_read(ar.UC_ARM64_REG_PC),'Output':base64.b64encode(bytes(em.mem_read(0x500000,65536))).decode()})
 r=subprocess.run([str(runner)],input=json.dumps(inputs),text=True,capture_output=True,check=True)
 actual=json.loads(r.stdout)
 for i,(a,b) in enumerate(zip(actual,expected)):
  assert a==b,(i,[key for key in a if a[key]!=b[key]])
 print(f'{len(actual)} real-block synthetic states match ARM64 execution (registers, flags, target, complete output region).')
