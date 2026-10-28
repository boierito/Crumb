#!/usr/bin/env python3
"""Load the shipped real profile and exercise generic device regression tests.

The machine state and memory in the tests are synthetic; no capture is read.
"""
import gzip,hashlib,json,pathlib,re,subprocess,tarfile,tempfile
from fairplay_generic_resource import read_profile
from fairplay_image_resource import read_image_text
ROOT=pathlib.Path(__file__).resolve().parent.parent
REV='f7242cce58077a4af8db738ed295ee23ad8f067d'
ARCHIVE_SHA='e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39'
compressed,pin=read_profile(ROOT)
with tempfile.TemporaryDirectory(prefix='crumb-generic-real-') as directory:
    work=pathlib.Path(directory)
    subprocess.run(['python3',str(ROOT/'scripts/prepare-fairplay-generic.py'),'runtime',str(work)],check=True)
    archive=ROOT/'third_party/archives'/f'ipatool-fairplay-{REV}.tar.gz'
    if hashlib.sha256(archive.read_bytes()).hexdigest()!=ARCHIVE_SHA:raise ValueError('source archive integrity mismatch')
    with tarfile.open(archive) as tar:tar.extractall(work/'upstream',filter='data')
    upstream=work/'upstream'/f'ipatool-{REV}/tools/fairplay-recovery'
    for name in ('oracle.go','generated_26b118.go','input_size_default.go'):
        source=(upstream/name).read_text().replace('package main','package rawsource',1)
        source=re.sub(r'^//go:build[^\n]*\n','',source)
        (work/name).write_text(source)
    for name in ('fairplaygeneric.go','fairplaygeneric_atomic.go','fairplaygeneric_atomic_test.go','fairplaygeneric_test.go','fairplaygeneric_real_test.go','fairplayrawblocks.go','fairplayrawflow.go','fairplayrawflow_test.go'):
        (work/name).write_bytes((ROOT/'MapleSyrup/NativeSAP'/name).read_bytes())
    (work/'profiles').mkdir();(work/'profiles/generic.ir.gz').write_bytes(compressed)
    encoded_image,image_pin=read_image_text(ROOT)
    native_text=gzip.decompress(encoded_image)
    assert hashlib.sha256(native_text).hexdigest()==image_pin['immutable_text']['sha256']
    (work/'profiles/native-text.bin').write_bytes(native_text)
    source=(upstream/'raw_source_step.go').read_text()
    definition=re.search(r'var rawSourceBinaryHash = \[32\]byte\{[^}]+\}',source).group(0)
    (work/'fixture_runtime_test.go').write_text('''package rawsource
import (_ "embed";"testing";"fmt")
type rawSourceRuntime struct{stopPC,lastKnownBlock uint64;genericOnly bool;outputScratchEnabled bool;outputScratchReady uint8;byteLoop genericLoopObservation;observe func(uint64);boundary func()(bool,error);additionalBlocks,genericBlocks int;genericInstructions uint64;genericLastPC,currentInstruction uint64;executionContext string;untranslatedTarget,targetClass string}
var knownSourceBlocks=map[uint64]sourceBlock{}
func (r *rawSourceRuntime) handle(s *state,pc uint64)(uint64,error){return pc,fmt.Errorf("external operation unavailable in isolated fixture")}
'''+definition+'''
//go:embed profiles/generic.ir.gz
var encoded []byte
//go:embed profiles/native-text.bin
var nativeText []byte
func init(){var err error;genericImage,err=loadGenericInstructionProfile(encoded,rawSourceBinaryHash,genericKernelABI);if err!=nil{panic(err)};if err=genericImage.attachNativeText(rawSourceBinaryHash,0x100000000,nativeText);err!=nil{panic(err)}}
func TestRealGenericResourceLoaded(t *testing.T){if genericImage==nil{t.Fatal("profile absent")}}
''')
    (work/'go.mod').write_text('module genericreal\n\ngo 1.26.0\n')
    subprocess.run(['go','test','-race','-count=1','./...'],cwd=work,check=True,timeout=180)
print('Shipped real profile passes integrity, generic dispatch and current device landing regression; no AOT handler for that PC.')
