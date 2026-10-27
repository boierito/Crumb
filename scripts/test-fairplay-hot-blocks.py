import base64, gzip, hashlib, json, pathlib, random, re, struct, subprocess, tarfile, tempfile
from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM, UC_HOOK_MEM_UNMAPPED, UcError, arm64_const as a
from fairplay_image_resource import read_image_text
root=pathlib.Path(__file__).resolve().parent.parent; tmp=root/'scripts/fixtures/fairplay-hot-blocks'
rev='f7242cce58077a4af8db738ed295ee23ad8f067d'; archive=root/'third_party/archives'/f'ipatool-fairplay-{rev}.tar.gz'
encoded,pin=read_image_text(root);text=gzip.decompress(encoded)
assert hashlib.sha256(text).hexdigest()==pin['immutable_text']['sha256']
assert hashlib.sha256(archive.read_bytes()).hexdigest()=='e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39'
source=(root/'MapleSyrup/NativeSAP/fairplayrawsupplement_hot.go').read_text()
comments={int(addr,16):asm for addr,asm in re.findall(r'// 0x([0-9a-f]+): ([^\n]+)',source)}
registers=[getattr(a,'UC_ARM64_REG_X'+str(i)) for i in range(29)]+[a.UC_ARM64_REG_X29,a.UC_ARM64_REG_X30]
harness='''package main
import("encoding/json";"os";"fmt")
type rawSourceRuntime struct { stopPC uint64 }
type Fixture struct { Entry,SP uint64; Registers [31]uint64; Flags [4]bool; Pages map[uint64][]byte }
type Result struct { Next,SP uint64; Registers [31]uint64; Flags [4]bool; Pages map[uint64][]byte }
func main(){ var inputs []Fixture; if e:=json.NewDecoder(os.Stdin).Decode(&inputs);e!=nil{panic(e)}; out:=[]Result{}; for _,f:=range inputs { m:=&model{registers:f.Registers,pages:f.Pages,sp:f.SP,n:f.Flags[0],z:f.Flags[1],c:f.Flags[2],v:f.Flags[3]}; s:=newState(m);s.reset(); var next uint64; var err error; switch f.Entry {
'''
for p in sorted(tmp.glob('*.json')):
 pc=int(p.stem,16);harness+=f'case {pc}: next,_,err=rawSourceBootstrapBlock{pc:x}(s,nil)\n'
harness+='''default:panic("bad entry") };if err!=nil{panic(fmt.Sprint(err))}; pages:=map[uint64][]byte{}; for addr:=range f.Pages { v,e:=s.outputBytes(addr,4096);if e!=nil{panic(e)};pages[addr]=v };out=append(out,Result{next,s.sp,s.x,[4]bool{s.n,s.z,s.c,s.v},pages}) };json.NewEncoder(os.Stdout).Encode(out) }
'''
rng=random.Random(23200);fixtures=[];expected=[]
with tempfile.TemporaryDirectory() as folder:
 w=pathlib.Path(folder)
 with tarfile.open(archive) as tar:tar.extractall(w,filter='data')
 upstream=w/f'ipatool-{rev}/tools/fairplay-recovery'
 for name in ('oracle.go','generated_26b118.go','input_size_default.go'):(w/name).write_bytes((upstream/name).read_bytes())
 (w/'candidate.go').write_text(source.replace('package rawsource','package main',1))
 (w/'main.go').write_text(harness)
 args=['go','build','-o',str(w/'run'),'main.go','oracle.go','generated_26b118.go','input_size_default.go','candidate.go'];subprocess.run(args,cwd=w,check=True)
 for f in sorted(tmp.glob('*.json')):
  d=json.loads(f.read_text());start=d['start'];data=b''.join(bytes.fromhex(i['hex']) for i in d['instructions']);assert text[start-0x100000000:start-0x100000000+len(data)]==data
  assert d['binary_sha256']=='1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784'
  assert all(comments.get(i['address'])==i['assembly'] for i in d['instructions'])
  # BRAA -> BR as used by existing portable PAC convention.
  ins=d['instructions'][-1]['assembly'];assert ins.startswith('braa '); rn=int(ins.split()[1].split(',')[0][1:]);code=data[:-4]+struct.pack('<I',0xd61f0000|(rn<<5))
  good=0
  for trial in range(60):
   u=Uc(UC_ARCH_ARM64,UC_MODE_ARM);u.ctl_set_cpu_model(a.UC_CPU_ARM64_MAX);page=start&~4095; codeend=(start+len(code)+4095)&~4095;u.mem_map(page,codeend-page);u.mem_write(start,code)
   regs=[rng.randrange(0,2**16) for i in range(31)];regs[29]=0x700800;flags=[bool(rng.getrandbits(1)) for _ in range(4)];pages={}
   def hook(u,acc,addr,size,value,data):
    base=addr&~4095
    if len(pages)>=128 or base>=0x800000000000: return False
    if base not in pages:
     try:u.mem_map(base,4096)
     except UcError:return False
     pages[base]=bytes(4096)
    return True
   u.hook_add(UC_HOOK_MEM_UNMAPPED,hook)
   for r,v in zip(registers,regs):u.reg_write(r,v)
   u.reg_write(a.UC_ARM64_REG_SP,0x710800);u.reg_write(a.UC_ARM64_REG_NZCV,sum(int(v)<<(31-i) for i,v in enumerate(flags)))
   try:u.emu_start(start,start+len(code),count=len(d['instructions']))
   except UcError:continue
   if len(pages)==0:continue
   nzcv=u.reg_read(a.UC_ARM64_REG_NZCV)
   fixtures.append({'Entry':start,'SP':0x710800,'Registers':regs,'Flags':flags,'Pages':{str(k):base64.b64encode(v).decode() for k,v in pages.items()}})
   expected.append({'Next':u.reg_read(a.UC_ARM64_REG_PC),'SP':u.reg_read(a.UC_ARM64_REG_SP),'Registers':[u.reg_read(r) for r in registers],'Flags':[bool(nzcv&(1<<(31-i))) for i in range(4)],'Pages':{str(k):base64.b64encode(u.mem_read(k,4096)).decode() for k in pages}})
   good+=1
   if good==12:break
  print(f.stem,'cases',good,flush=True)
  if good<4:raise RuntimeError('too few fixtures')
 r=subprocess.run([str(w/'run')],input=json.dumps(fixtures),text=True,capture_output=True,check=True);actual=json.loads(r.stdout)
 for i,(x,y) in enumerate(zip(actual,expected)):
  if x!=y:print('MISMATCH',i,hex(fixtures[i]['Entry']),[k for k in x if x[k]!=y[k]]);raise SystemExit(1)
 print('All',len(actual),'public synthetic states match Unicorn for all five recovery blocks; no PAC authentication is asserted.')
