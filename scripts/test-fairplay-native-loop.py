import pathlib,sys,importlib.util,gzip,hashlib,tarfile,tempfile,subprocess,json,struct,random,base64,re
root=pathlib.Path(__file__).resolve().parent.parent;sys.path.insert(0,str(root/'scripts'))
from fairplay_generic_resource import read_profile
from fairplay_image_resource import read_image_text
from unicorn import Uc,UC_ARCH_ARM64,UC_MODE_ARM,UC_HOOK_MEM_UNMAPPED,arm64_const as arm
spec=importlib.util.spec_from_file_location('differential',root/'scripts/test-fairplay-generic.py');mod=importlib.util.module_from_spec(spec);spec.loader.exec_module(mod)
enc,pin=read_profile(root);compressed,imgpin=read_image_text(root);text=gzip.decompress(compressed)
assert hashlib.sha256(text).hexdigest()==imgpin['immutable_text']['sha256']
start=0x1001005ac;end=0x1001006b0
reference=json.loads((root/'scripts/fixtures/fairplay-byte-loop.json').read_text())
assert reference['binary_sha256']==pin['binary_sha256'] and reference['start']==start
comments=re.findall(r'// 0x([0-9a-f]+): ([^\n]+)',(root/'MapleSyrup/NativeSAP/fairplayrawsupplement_byteloop.go').read_text())
assert len(reference['instructions'])==len(comments)==(end-start)//4
for item,(addr,assembly) in zip(reference['instructions'],comments):
 assert int(addr,16)==item['address'] and assembly==item['assembly']
 assert text[item['address']-0x100000000:item['address']-0x100000000+4].hex()==item['hex']
regsids=[getattr(arm,'UC_ARM64_REG_X'+str(i)) for i in range(29)]+[arm.UC_ARM64_REG_X29,arm.UC_ARM64_REG_X30]
rng=random.Random(23053);fixtures=[];expected=[]
lengths=[1,2,15,16,17,32,440,448,888,4096,145600]
for length in lengths:
 regs=[rng.getrandbits(64) for _ in range(31)];regs[9]=length
 for r,v in {0:0,2:0x1000,3:0,6:1,8:0x640000,11:3,14:0x7b,15:0,16:0xc9,23:0x600000,24:0x6a0000,25:0x660000,27:0x650000,28:0x800000,29:0x610100,30:0x630000}.items():regs[r]=v
 pages={a:bytearray(rng.randbytes(4096)) for a in [0x600000,0x610000,0x620000,0x630000,0x640000,0x650000,0x660000,0x670000,0x671000,0x690000,0x6a0000,0x101075000]}
 for a in range(0x800000,0x800000+((length+4095)//4096)*4096,4096):pages[a]=bytearray(rng.randbytes(4096))
 for a in range(0x10105b000,0x101065000,4096):pages[a]=bytearray(4096)
 for offset in range(0x6f0,0x9000,8):
  a=0x10105b000+offset;struct.pack_into('<Q',pages[a&~4095],a&4095,0x690000)
 struct.pack_into('<Q',pages[0x610000],0x100-0xa0,0x620000)
 struct.pack_into('<Q',pages[0x610000],0x100-0xb0,0x800000)
 struct.pack_into('<Q',pages[0x6a0000],3*8,start);struct.pack_into('<Q',pages[0x6a0000],2*8,end)
 u=Uc(UC_ARCH_ARM64,UC_MODE_ARM);u.ctl_set_cpu_model(arm.UC_CPU_ARM64_MAX)
 u.mem_map(0x100100000,4096);u.mem_write(start,text[start-0x100000000:end-0x100000000]);u.mem_map(end&~4095,4096) if (end&~4095)!=(start&~4095) else None
 for a,data in pages.items():u.mem_map(a,4096);u.mem_write(a,bytes(data))
 for r,v in zip(regsids,regs):u.reg_write(r,v)
 u.reg_write(arm.UC_ARM64_REG_SP,0x71010000)
 u.hook_add(UC_HOOK_MEM_UNMAPPED,lambda u,access,a,size,value,data: print("fixture fault",length,hex(u.reg_read(arm.UC_ARM64_REG_PC)),hex(a),size) or False)
 u.emu_start(start,end,count=65*length+2)
 pc=u.reg_read(arm.UC_ARM64_REG_PC);assert pc==end,(length,hex(pc),u.reg_read(arm.UC_ARM64_REG_X9))
 nzcv=u.reg_read(arm.UC_ARM64_REG_NZCV)
 expected.append({'PC':pc,'SP':u.reg_read(arm.UC_ARM64_REG_SP),'Registers':[u.reg_read(r) for r in regsids],'Flags':[bool(nzcv&(1<<(31-j))) for j in range(4)],'Vectors':[[0]*16 for _ in range(32)],'Pages':{str(a):base64.b64encode(u.mem_read(a,4096)).decode() for a in pages},'LoopVisits':length,'LoopSequences':1,'LoopBucket':'1_to_4096' if length<=4096 else '4097_to_1m'})
 fixtures.append({'PC':start,'End':end,'SP':0x71010000,'Registers':regs,'Flags':[False]*4,'Vectors':[[0]*16 for _ in range(32)],'Pages':{str(a):base64.b64encode(v).decode() for a,v in pages.items()}})
with tempfile.TemporaryDirectory(prefix='crumb-native-loop-') as folder:
 work=pathlib.Path(folder);runner=work/'runtime'
 subprocess.run(['python3',str(root/'scripts/prepare-fairplay-generic.py'),'runtime',str(runner)],check=True)
 rev='f7242cce58077a4af8db738ed295ee23ad8f067d';archive=root/'third_party/archives'/f'ipatool-fairplay-{rev}.tar.gz'
 assert hashlib.sha256(archive.read_bytes()).hexdigest()=='e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39'
 with tarfile.open(archive) as z:z.extractall(work/'upstream',filter='data')
 upstream=work/'upstream'/f'ipatool-{rev}/tools/fairplay-recovery'
 for name in ('oracle.go','generated_26b118.go','input_size_default.go'):(runner/name).write_bytes((upstream/name).read_bytes())
 for name in ('fairplaygeneric.go','fairplaygeneric_atomic.go','fairplayrawblocks.go','fairplayrawflow.go'):(runner/name).write_bytes((root/'MapleSyrup/NativeSAP'/name).read_bytes())
 for p in runner.glob('*.go'):p.write_text(p.read_text().replace('package rawsource','package main',1))
 h=mod.HARNESS.replace('PC,SP uint64;Chunks','PC,SP,End uint64;Chunks').replace('"differential",f.PC,0x200000,1,128','"differential",f.PC,f.End,1,200000').replace('sha256.Sum256(text)!=source','hex.EncodeToString(sha256Sum(text))!="'+imgpin['immutable_text']['sha256']+'"')
 h=h.replace('Vectors [32][16]byte;Pages map[uint64][]byte}\nfunc main()', 'Vectors [32][16]byte;Pages map[uint64][]byte;LoopVisits,LoopSequences uint64;LoopBucket string}\nfunc main()').replace('s.vectors,pages})','s.vectors,pages,runtime.byteLoop.visits,runtime.byteLoop.sequences,runtime.byteLoop.maxBucket})')
 # Keep native source identity pinned, while verifying the exact __TEXT segment.
 h=h.replace('func main(){','func sha256Sum(b []byte)[]byte{s:=sha256.Sum256(b);return s[:]}\nfunc main(){')
 (runner/'main.go').write_text(h);(work/'profile.gz').write_bytes(enc);(work/'text.bin').write_bytes(text)
 subprocess.run(['go','build','-o',str(work/'runner'),*map(str,runner.glob('*.go'))],check=True)
 r=subprocess.run([str(work/'runner'),str(work/'profile.gz'),pin['binary_sha256'],str(work/'text.bin')],input=json.dumps(fixtures),text=True,capture_output=True)
 if r.returncode:print(r.stderr);raise SystemExit(r.returncode)
 actual=json.loads(r.stdout)
 for length,a,b in zip(lengths,actual,expected):
  if a!=b:
   print('MISMATCH',length,[k for k in a if a[k]!=b[k]])
   for i,(x,y) in enumerate(zip(a['Registers'],b['Registers'])):
    if x!=y:print('register',i,hex(x),hex(y))
   print('flags',a['Flags'],b['Flags']);raise SystemExit(1)
 # Execute exactly the same cases through the compiled Go block as well.
 # The native emulator executes the original bytes, including the return
 # branch using its configured ARM64 PAC model; no guest data is involved.
 candidate=(root/'MapleSyrup/NativeSAP/fairplayrawsupplement_byteloop.go').read_text().replace('package rawsource','package main',1)
 (runner/'candidate.go').write_text(candidate)
 h=h.replace('var knownSourceBlocks=map[uint64]sourceBlock{}','var knownSourceBlocks=map[uint64]sourceBlock{0x1001005ac:{run:rawSourceBootstrapBlock1001005ac}}').replace('genericOnly:true','genericOnly:false')
 (runner/'main.go').write_text(h)
 subprocess.run(['go','build','-o',str(work/'runner'),*map(str,runner.glob('*.go'))],check=True)
 r=subprocess.run([str(work/'runner'),str(work/'profile.gz'),pin['binary_sha256'],str(work/'text.bin')],input=json.dumps(fixtures),text=True,capture_output=True)
 if r.returncode:print(r.stderr);raise SystemExit(r.returncode)
 for length,a,b in zip(lengths,json.loads(r.stdout),expected):
  if a!=b:
   print('AOT MISMATCH',length,[k for k in a if a[k]!=b[k]]);raise SystemExit(1)
 print('Generic and AOT execution of the exact native loop match Unicorn for 11 lengths, including 145600 bytes, registers, flags, vectors, full pages and observed iteration counts.')
