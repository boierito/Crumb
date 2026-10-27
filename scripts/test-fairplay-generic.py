#!/usr/bin/env python3
"""End-to-end generic IR differential tests; all programs/state are synthetic."""
import base64
import gzip
import hashlib
import json
import os
from pathlib import Path
import random
import platform
import sys
import struct
import subprocess
import tarfile
import tempfile

from keystone import Ks, KS_ARCH_ARM64, KS_MODE_LITTLE_ENDIAN
from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM
from unicorn import arm64_const as arm

ROOT=Path(__file__).resolve().parent.parent
BASE=0x100001000
REV='f7242cce58077a4af8db738ed295ee23ad8f067d'
ARCHIVE_SHA='e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39'
HARNESS=r'''package main
import("crypto/sha256";"encoding/json";"encoding/hex";"os";"fmt")
type rawSourceRuntime struct {stopPC,lastKnownBlock uint64;genericOnly bool;observe func(uint64);boundary func()(bool,error);additionalBlocks int;genericBlocks int;genericInstructions uint64;genericLastPC,currentInstruction uint64;executionContext string;untranslatedTarget,targetClass string}
var knownSourceBlocks=map[uint64]sourceBlock{}
func (r *rawSourceRuntime) handle(s *state,pc uint64)(uint64,error){return pc,fmt.Errorf("external operation unavailable in isolated fixture")}
type fixture struct {PC,SP uint64;Chunks int;Registers [31]uint64;Flags [4]bool;Vectors [32][16]byte;Pages map[uint64][]byte}
type result struct {PC,SP uint64;Registers [31]uint64;Flags [4]bool;Vectors [32][16]byte;Pages map[uint64][]byte}
func main(){
 data,err:=os.ReadFile(os.Args[1]);if err!=nil{panic(err)}
 hash,err:=hex.DecodeString(os.Args[2]);if err!=nil{panic(err)};var source [32]byte;copy(source[:],hash)
 genericImage,err=loadGenericInstructionProfile(data,source,genericKernelABI);if err!=nil{panic(err)}
 text,err:=os.ReadFile(os.Args[3]);if err!=nil{panic(err)};if sha256.Sum256(text)!=source{panic("native source mismatch")};if err=genericImage.attachNativeText(source,0x100000000,text);err!=nil{panic(err)}
 var inputs []fixture;if err=json.NewDecoder(os.Stdin).Decode(&inputs);err!=nil{panic(err)}
 results:=[]result{}
 for _,f:=range inputs {
  m:=&model{registers:f.Registers,pages:f.Pages,sp:f.SP,n:f.Flags[0],z:f.Flags[1],c:f.Flags[2],v:f.Flags[3]}
  s:=newState(m);s.reset();s.vectors=f.Vectors
  runtime:=&rawSourceRuntime{genericOnly:true};pc,err:=runNativeRegion(s,runtime,f.PC,nativeRegion{"differential",f.PC,0x200000,1,128});if err!=nil{panic("native region execution failed: "+err.Error())}
  pages:=map[uint64][]byte{};for address:=range f.Pages {value,err:=s.outputBytes(address,4096);if err!=nil{panic(err)};pages[address]=value}
  results=append(results,result{pc,s.sp,s.x,[4]bool{s.n,s.z,s.c,s.v},s.vectors,pages})
 }
 if err=json.NewEncoder(os.Stdout).Encode(results);err!=nil{panic(err)}
}
'''

def macho(text,stub):
    base=BASE-4096
    stub_offset=((4096+len(text)+4095)//4096)*4096+4096
    total=stub_offset+len(stub)
    header=struct.pack('<8I',0xfeedfacf,0x100000c,2,2,1,232,0,0)
    segment=struct.pack('<II16sQQQQIIII',0x19,232,b'__TEXT',base,((total+4095)//4096)*4096,0,total,7,5,2,0)
    section=struct.pack('<16s16sQQIIIIIIII',b'__text',b'__TEXT',BASE,len(text),4096,2,0,0,0x80000400,0,0,0)
    helper=struct.pack('<16s16sQQIIIIIIII',b'__auth_stubs',b'__TEXT',base+stub_offset,len(stub),stub_offset,2,0,0,0x80000408,0,0,0)
    data=(header+segment+section+helper).ljust(4096,b'\0')+text
    return data.ljust(stub_offset,b'\xff')+stub,base+stub_offset


def main():
    ks=Ks(KS_ARCH_ARM64,KS_MODE_LITTLE_ENDIAN)
    rng=random.Random(23044)
    pool=bytearray();cases=[]
    branch=bytes(ks.asm('br x28')[0])
    def add(code,flags=None,pages=None,registers=None,vectors=None,sp=0x700100):
        registers=registers or [rng.getrandbits(64) for _ in range(31)]
        registers[28]=0x200000
        cases.append({'PC':BASE+len(pool),'SP':sp,'Registers':registers,'Flags':flags or [False]*4,'Pages':pages or {},'Vectors':vectors or [[0]*16 for _ in range(32)],'code':code})
        pool.extend(code)
    # Canonical CSEL/CSINC/CSINV/CSNEG: every condition and every NZCV,
    # W/X widths, zero-register operands and destination, wrapping operations.
    for width in (32,64):
        for operation in range(4):
            for condition in range(16):
                for nzcv in range(16):
                    destination=31 if nzcv%7==0 else 2
                    left=31 if nzcv%5==0 else 0
                    right=31 if nzcv%3==0 else 1
                    word=0x1a800000 | (int(width==64)<<31) | ((operation>>1)<<30) | ((operation&1)<<10) | (right<<16) | (condition<<12) | (left<<5) | destination
                    flags=[bool(nzcv&(8>>j)) for j in range(4)]
                    add(struct.pack('<I',word)+branch,flags=flags)
    for _ in range(16):
        regs=[rng.getrandbits(64) for _ in range(31)];regs[30]=0x200000
        add(bytes(ks.asm('blr x30')[0]),registers=regs)
    condition_names=['eq','ne','hs','lo','mi','pl','vs','vc','hi','ls','ge','lt','gt','le','al','nv']
    for width in (32,64):
        prefix='w' if width==32 else 'x'
        for operation in ('ccmp','ccmn'):
            for condition in condition_names:
                for nzcv in range(16):
                    operand=f'#{nzcv}' if nzcv%2 else prefix+'1'
                    code=bytes(ks.asm(f'{operation} {prefix}0,{operand},#{15-nzcv},{condition}')[0])+branch
                    add(code,flags=[bool(nzcv&(8>>j)) for j in range(4)])
    for i in range(64):
        code=bytes(ks.asm('add w0,w0,w2; str w0,[sp,#4]; ldr x5,[x4,w3,sxtw#3]; br x5')[0])
        regs=[rng.getrandbits(64) for _ in range(31)]
        index=[-256,-1,0,255][i%4]
        regs[3]=(rng.getrandbits(32)<<32)|(index&0xffffffff);regs[4]=0x600800
        table=bytearray(rng.randbytes(4096));struct.pack_into('<Q',table,0x800+index*8,0x200000)
        pages={0x600000:bytes(table),0x700000:rng.randbytes(4096)}
        add(code,flags=[bool(rng.getrandbits(1)) for _ in range(4)],pages=pages,registers=regs,sp=0x700ff8 if i%2 else 0x700100)
    for instruction in ['add v0.4s,v1.4s,v2.4s']+[f'{op} v0.{width}b,v1.{width}b,v2.{width}b' for op in ('and','orr','eor','bic','orn','bsl','bit','bif') for width in (8,16)]:
        for _ in range(8):
            vectors=[list(rng.randbytes(16)) for _ in range(32)]
            add(bytes(ks.asm(instruction)[0])+branch,vectors=vectors)
    # Full native execution of data-dependent loops and calls/returns. The
    # region runner must follow PCs rather than an assumed chunk count.
    for n in range(1,33):
        code=bytes(ks.asm(f'mov x0,#{n}; mov x1,#0; again: add x1,x1,#3; subs x0,x0,#1; b.ne again; br x28')[0])
        add(code);cases[-1]['NativeCount']=3*n+3
    for _ in range(32):
        add(bytes(ks.asm('bl inner; add x4,x4,#2; br x28; inner: add x3,x3,#7; ret')[0]))
    # Every scalar CAS width and memory-order variant, successful and failed
    # comparisons, SP, zero registers, and aliased operands. Independent Unicorn
    # execution checks ALL registers, flags, vectors and the entire memory page.
    cas_cases=0
    for size in range(4):
        width=1<<size;mask=(1<<(width*8))-1
        for order in range(4):
            for equal in (False,True):
                for rs,rt,rn in ((11,12,19),(11,11,19),(19,12,19),(11,19,19),(31,12,19),(11,31,19),(11,12,31),(31,31,31)):
                    regs=[rng.getrandbits(64) for _ in range(31)]
                    address=0x600100
                    if rn!=31:regs[rn]=address
                    expected=0 if rs==31 else regs[rs]&mask
                    old=expected if equal else expected^1
                    page=bytearray(rng.randbytes(4096));page[0x100:0x100+width]=old.to_bytes(width,'little')
                    word=0x08a07c00|(size<<30)|((order>>1)<<22)|((order&1)<<15)|(rs<<16)|(rn<<5)|rt
                    add(struct.pack('<I',word)+branch,registers=regs,pages={0x600000:bytes(page)},sp=address,flags=[bool(rng.getrandbits(1)) for _ in range(4)])
                    cas_cases+=1
    # 64-bit SIMD pair load/store, including the exact epilogue forms found
    # in the pinned daemon, negative offsets, SP and register writeback.
    pair_cases=0
    for operation in ('ldp','stp'):
        for mode in ('offset','pre','post'):
            for offset in (-128,-8,0,8,120):
                for base in ('sp','x19'):
                    address=f'[{base}, #{offset}]' if mode=='offset' else (f'[{base}, #{offset}]!' if mode=='pre' else f'[{base}], #{offset}')
                    code=bytes(ks.asm(f'{operation} d9, d8, {address}')[0])+branch
                    regs=[rng.getrandbits(64) for _ in range(31)];regs[19]=0x600800
                    vectors=[list(rng.randbytes(16)) for _ in range(32)]
                    add(code,registers=regs,vectors=vectors,pages={0x600000:rng.randbytes(4096)},sp=0x600800)
                    pair_cases+=1
    # Branch from __text to a different executable section through a register.
    # A deliberate non-executable gap separates the two sections.
    add(bytes(ks.asm('br x5')[0]))
    cross=cases[-1]
    stub_code=bytes(ks.asm('mov w0,#7; br x28')[0])
    binary,stub_address=macho(bytes(pool),stub_code)
    cross['Registers'][5]=stub_address;cross['Chunks']=2;cross['NativeCount']=3
    cases.append({'PC':stub_address,'SP':0x700100,'Registers':[0]*28+[0x200000,0,0],'Flags':[False]*4,'Pages':{},'Vectors':[[0]*16 for _ in range(32)],'code':stub_code})
    with tempfile.TemporaryDirectory(prefix='crumb-generic-differential-') as folder:
        work=Path(folder);compiler=work/'compiler';runner=work/'runtime'
        subprocess.run(['python3',str(ROOT/'scripts/prepare-fairplay-generic.py'),'compiler',str(compiler)],check=True)
        subprocess.run(['go','mod','tidy'],cwd=compiler,check=True)
        subprocess.run(['go','build','-o','compiler','.'],cwd=compiler,check=True,timeout=120)
        binary_path=work/'synthetic.macho';binary_path.write_bytes(binary)
        source_hash=hashlib.sha256(binary).hexdigest();profile=work/'synthetic.ir.gz';report=work/'census.json'
        subprocess.run([str(compiler/'compiler'),'--binary',str(binary_path),'--binary-sha256',source_hash,'--output',str(profile),'--report',str(report)],check=True)
        census=json.loads(report.read_text())
        assert census['unsupported_slots']==cas_cases+pair_cases,census
        assert census['span_slots']>census['instruction_slots'] and len(census['executable_sections'])==2,census
        subprocess.run(['python3',str(ROOT/'scripts/prepare-fairplay-generic.py'),'runtime',str(runner)],check=True)
        archive=ROOT/'third_party/archives'/f'ipatool-fairplay-{REV}.tar.gz'
        assert hashlib.sha256(archive.read_bytes()).hexdigest()==ARCHIVE_SHA
        with tarfile.open(archive) as tar:tar.extractall(work/'upstream',filter='data')
        upstream=work/'upstream'/f'ipatool-{REV}/tools/fairplay-recovery'
        for name in ('oracle.go','generated_26b118.go','input_size_default.go'):(runner/name).write_bytes((upstream/name).read_bytes())
        (runner/'generic.go').write_bytes((ROOT/'MapleSyrup/NativeSAP/fairplaygeneric.go').read_bytes())
        for name in ('fairplaygeneric_atomic.go','fairplayrawblocks.go','fairplayrawflow.go'):(runner/name).write_bytes((ROOT/'MapleSyrup/NativeSAP'/name).read_bytes())
        for name in runner.glob('*.go'):name.write_text(name.read_text().replace('package rawsource','package main',1))
        (runner/'main.go').write_text(HARNESS)
        subprocess.run(['go','build','-o',str(work/'runner'),*map(str,runner.glob('*.go'))],check=True,timeout=120)
        registers=[getattr(arm,'UC_ARM64_REG_X'+str(i)) for i in range(29)]+[arm.UC_ARM64_REG_X29,arm.UC_ARM64_REG_X30]
        expected=[];inputs=[]
        for case in cases:
            emulator=Uc(UC_ARCH_ARM64,UC_MODE_ARM)
            emulator.ctl_set_cpu_model(arm.UC_CPU_ARM64_MAX)
            page=case['PC']&~4095
            size=((case['PC']-page+len(case['code'])+4095)//4096)*4096
            emulator.mem_map(page,size);emulator.mem_write(case['PC'],case['code']);emulator.mem_map(0x200000,4096)
            if page!=(stub_address&~4095):emulator.mem_map(stub_address&~4095,4096);emulator.mem_write(stub_address,stub_code)
            for address,data in case['Pages'].items():emulator.mem_map(address,4096);emulator.mem_write(address,data)
            for register,value in zip(registers,case['Registers']):emulator.reg_write(register,value)
            for i,value in enumerate(case['Vectors']):emulator.reg_write(getattr(arm,'UC_ARM64_REG_Q'+str(i)),int.from_bytes(bytes(value),'little'))
            emulator.reg_write(arm.UC_ARM64_REG_SP,case['SP']);emulator.reg_write(arm.UC_ARM64_REG_NZCV,sum(int(flag)<<(31-j) for j,flag in enumerate(case['Flags'])))
            emulator.emu_start(case['PC'],case['PC']+len(case['code']),count=case.get('NativeCount',len(case['code'])//4))
            nzcv=emulator.reg_read(arm.UC_ARM64_REG_NZCV)
            expected.append({'PC':emulator.reg_read(arm.UC_ARM64_REG_PC),'SP':emulator.reg_read(arm.UC_ARM64_REG_SP),'Registers':[emulator.reg_read(reg) for reg in registers],'Flags':[bool(nzcv&(1<<(31-j))) for j in range(4)],'Vectors':[list(emulator.reg_read(getattr(arm,'UC_ARM64_REG_Q'+str(i))).to_bytes(16,'little')) for i in range(32)],'Pages':{str(address):base64.b64encode(bytes(emulator.mem_read(address,4096))).decode() for address in case['Pages']}})
            inputs.append({key:({str(address):base64.b64encode(data).decode() for address,data in value.items()} if key=='Pages' else value) for key,value in case.items() if key not in ('code','NativeCount')})
        command=[str(work/'runner'),str(profile),source_hash,str(binary_path)]
        if sys.platform=='linux' and platform.machine()=='x86_64':
            launcher=work/'deny-exec-runner'
            subprocess.run(['cc','-std=c11','-Wall','-Wextra','-Werror','-I',str(ROOT/'scripts'),str(ROOT/'scripts/fixtures/tci/generic-deny-exec-launcher.c'),'-o',str(launcher)],check=True)
            command=[str(launcher)]+command
        result=subprocess.run(command,input=json.dumps(inputs),text=True,capture_output=True,check=True,timeout=120)
        (runner/'generic_test.go').write_text((ROOT/'MapleSyrup/NativeSAP/fairplaygeneric_test.go').read_text().replace('package rawsource','package main',1))
        subprocess.run(['go','test','-race','-count=1',*map(str,runner.glob('*.go'))],check=True,timeout=120)
        actual=json.loads(result.stdout);assert len(actual)==len(expected)
        for i,(observed,reference) in enumerate(zip(actual,expected)):
            assert observed==reference,(i,[key for key in observed if observed[key]!=reference[key]])
    print(f'{len(cases)} native-region programs match independent ARM64 execution, including loops, calls/returns and executable-section transitions.')
    if sys.platform=='linux' and platform.machine()=='x86_64':print('Generic execution also passed with new mmap/mprotect PROT_EXEC mappings denied by seccomp.')

if __name__=='__main__':main()
