#!/usr/bin/env python3
"""Compare the 23042 bootstrap continuation with independent ARM64 execution.

Synthetic memory only. BRAA is represented as BR, matching the portable
runtime convention; native pointer authentication is outside this test.
"""
import base64
import hashlib
import json
from pathlib import Path
import random
import re
import struct
import subprocess
import tarfile
import tempfile

from keystone import Ks, KS_ARCH_ARM64, KS_MODE_LITTLE_ENDIAN
from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM
from unicorn import arm64_const as arm

ROOT = Path(__file__).resolve().parent.parent
REV = "f7242cce58077a4af8db738ed295ee23ad8f067d"
ARCHIVE_HASH = "e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39"
START = 0x1008C292C
STACK = 0x71010000
HARNESS = r'''package main
import ("encoding/json"; "os"; "fmt")
type rawSourceRuntime struct { stopPC uint64 }
type fixture struct { Registers [31]uint64; Pages map[uint64][]byte; Flags [4]bool; SP uint64 }
type result struct { Registers [31]uint64; Flags [4]bool; Next, SP uint64; Stack []byte }
func main() {
 var inputs []fixture
 if err:=json.NewDecoder(os.Stdin).Decode(&inputs); err!=nil {panic(err)}
 results:=[]result{}
 for _, f:=range inputs {
  m:=&model{registers:f.Registers,pages:f.Pages,sp:f.SP,n:f.Flags[0],z:f.Flags[1],c:f.Flags[2],v:f.Flags[3]}
  s:=newState(m); s.reset()
  next,stopped,err:=rawSourceBootstrapBlock1008c292c(s,nil)
  if err!=nil||stopped {panic(fmt.Sprint(err))}
  stack,err:=s.outputBytes(0x71010000,4096); if err!=nil {panic(err)}
  results=append(results,result{s.x,[4]bool{s.n,s.z,s.c,s.v},next,s.sp,stack})
 }
 if err:=json.NewEncoder(os.Stdout).Encode(results); err!=nil {panic(err)}
}
'''


def main():
    source = (ROOT / "MapleSyrup/NativeSAP/fairplayrawsupplement_exit.go").read_text()
    assembly = re.findall(r"// 0x[0-9a-f]+: (.+)", source)
    assert assembly[-1] == "braa x9, x17"
    assembly[-1] = "br x9"
    code = bytes(Ks(KS_ARCH_ARM64, KS_MODE_LITTLE_ENDIAN).asm("; ".join(assembly))[0])
    randomizer = random.Random(23043)
    register_ids = [getattr(arm, "UC_ARM64_REG_X" + str(i)) for i in range(29)] + [arm.UC_ARM64_REG_X29, arm.UC_ARM64_REG_X30]
    inputs, expected = [], []
    for i in range(64):
        index = [-256, -64, -1, 0, 1, 63, 255, randomizer.randint(-256, 255)][i % 8]
        registers = [randomizer.getrandbits(64) for _ in range(31)]
        registers[12] = (randomizer.getrandbits(32) << 32) | ((index - 0x2aa3) & 0xffffffff)
        registers[30] = 0x600800
        sp = STACK + (0x100 if i % 2 else 0)
        target = 0x200000 + i * 4
        table = bytearray(randomizer.randbytes(4096))
        struct.pack_into("<Q", table, 0x800 + index * 8, target)
        pages = {STACK: randomizer.randbytes(4096), 0x600000: bytes(table)}
        flags = [bool(randomizer.getrandbits(1)) for _ in range(4)]
        inputs.append({"Registers": registers, "Pages": {k: base64.b64encode(v).decode() for k, v in pages.items()}, "Flags": flags, "SP": sp})
        emulator = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
        emulator.mem_map(START & ~4095, 4096)
        emulator.mem_write(START, code)
        emulator.mem_map(0x200000, 4096)
        for address, data in pages.items():
            emulator.mem_map(address, 4096)
            emulator.mem_write(address, data)
        for register, value in zip(register_ids, registers):
            emulator.reg_write(register, value)
        emulator.reg_write(arm.UC_ARM64_REG_SP, sp)
        emulator.reg_write(arm.UC_ARM64_REG_NZCV, sum(int(flag) << (31 - j) for j, flag in enumerate(flags)))
        emulator.emu_start(START, START + len(code), count=len(assembly))
        nzcv = emulator.reg_read(arm.UC_ARM64_REG_NZCV)
        expected.append({"Registers": [emulator.reg_read(reg) for reg in register_ids], "Flags": [bool(nzcv & (1 << (31-j))) for j in range(4)], "Next": emulator.reg_read(arm.UC_ARM64_REG_PC), "SP": emulator.reg_read(arm.UC_ARM64_REG_SP), "Stack": base64.b64encode(bytes(emulator.mem_read(STACK, 4096))).decode()})

    with tempfile.TemporaryDirectory(prefix="crumb-exit-reference-") as directory:
        work = Path(directory)
        archive = ROOT / "third_party/archives" / f"ipatool-fairplay-{REV}.tar.gz"
        if hashlib.sha256(archive.read_bytes()).hexdigest() != ARCHIVE_HASH:
            raise ValueError("source integrity mismatch")
        with tarfile.open(archive) as tar:
            tar.extractall(work, filter="data")
        upstream = work / f"ipatool-{REV}/tools/fairplay-recovery"
        for name in ("oracle.go", "generated_26b118.go", "input_size_default.go"):
            (work / name).write_bytes((upstream / name).read_bytes())
        (work / "candidate.go").write_text(source.replace("package rawsource", "package main", 1))
        (work / "main.go").write_text(HARNESS)
        runner = work / "reference-runner"
        subprocess.run(["go", "build", "-o", str(runner), "main.go", "oracle.go", "generated_26b118.go", "input_size_default.go", "candidate.go"], cwd=work, check=True, timeout=120)
        result = subprocess.run([str(runner)], input=json.dumps(inputs), text=True, capture_output=True, check=True, timeout=30)
    observed = json.loads(result.stdout)
    assert len(observed) == len(expected) == 64
    for i, (actual, reference) in enumerate(zip(observed, expected)):
        assert actual == reference, (i, [key for key in actual if actual[key] != reference[key]])
    print("64 bootstrap-exit states match ARM64 execution: registers, NZCV, SP, signed dispatch and entire stack page.")


if __name__ == "__main__":
    main()
