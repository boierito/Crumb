#!/usr/bin/env python3
"""Public synthetic fixtures; no enrollment or daemon bytes are used."""
import hashlib
import importlib.util
import json
from pathlib import Path
import random
import subprocess
import sys
import tarfile
import tempfile
from types import SimpleNamespace
import unittest

from capstone import CS_ARCH_ARM64, CS_MODE_ARM, Cs
from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM
from unicorn.arm64_const import UC_ARM64_REG_X0, UC_ARM64_REG_X1, UC_ARM64_REG_X2, UC_ARM64_REG_PC

spec = importlib.util.spec_from_file_location("candidate", Path(__file__).with_name("generate-fairplay-next-block.py"))
candidate = importlib.util.module_from_spec(spec)
spec.loader.exec_module(candidate)


class GenerationTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.directory = tempfile.TemporaryDirectory()
        archive = candidate.ROOT / "third_party/archives" / f"ipatool-fairplay-{candidate.REV}.tar.gz"
        assert hashlib.sha256(archive.read_bytes()).hexdigest() == candidate.ARCHIVE_HASH
        with tarfile.open(archive) as tar:
            tar.extractall(cls.directory.name, filter="data")
        sys.path.insert(0, str(Path(cls.directory.name) / f"ipatool-{candidate.REV}/tools/fairplay-recovery/python"))
        import generate_macos26a428_bootstrap_cfg_go
        cls.translator = generate_macos26a428_bootstrap_cfg_go

    @classmethod
    def tearDownClass(cls):
        cls.directory.cleanup()

    def fixture(self, raw, executable=True):
        start = candidate.TARGET
        section = SimpleNamespace(contains_instructions=executable, vmaddr=start, size=len(raw), fileoff=0)
        database = SimpleNamespace(sections=[section], image=raw,
                                   bytes=SimpleNamespace(get_byte_at=lambda a, **kw: raw[a-start]))
        decoder = Cs(CS_ARCH_ARM64, CS_MODE_ARM)
        decoder.detail = True
        return database, decoder

    def test_requires_executable_section(self):
        db, decoder = self.fixture(bytes.fromhex("20001fd6"), False)
        with self.assertRaisesRegex(ValueError, "executable-section"):
            candidate.generate(db, decoder, self.translator)

    def test_does_not_continue_beyond_section(self):
        db, decoder = self.fixture(bytes.fromhex("0000028b"))
        with self.assertRaisesRegex(ValueError, "terminator"):
            candidate.generate(db, decoder, self.translator)

    def test_rejects_invalid_instruction(self):
        db, decoder = self.fixture(bytes.fromhex("ffffffff"))
        with self.assertRaisesRegex(ValueError, "decoded"):
            candidate.generate(db, decoder, self.translator)

    def test_rejects_wrong_binary_before_generating(self):
        with tempfile.TemporaryDirectory() as folder:
            binary, output = Path(folder)/"wrong", Path(folder)/"candidate.go"
            binary.write_bytes(b"synthetic incorrect daemon")
            result = subprocess.run([sys.executable, str(Path(candidate.__file__)), "--binary", str(binary), "--output", str(output)], capture_output=True)
            self.assertNotEqual(result.returncode, 0)
            self.assertFalse(output.exists())

    def test_generated_go_matches_independent_arm64_execution(self):
        # add x0,x0,x2; br x1. Include overflowing values and varying targets.
        raw = bytes.fromhex("0000028b20001fd6")
        db, decoder = self.fixture(raw)
        source = candidate.generate(db, decoder, self.translator).replace("package rawsource", "package main", 1)
        rng = random.Random(23042)
        vectors = [[0xffffffffffffffff, 0x200000, 1], [0, 0x200004, 0]]
        vectors += [[rng.getrandbits(64), 0x200000 + 4*i, rng.getrandbits(64)] for i in range(20)]
        expected = []
        for x0, x1, x2 in vectors:
            emu = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
            base = candidate.TARGET & ~0xfff
            emu.mem_map(base, 0x1000)
            emu.mem_map(0x200000, 0x1000)
            emu.mem_write(candidate.TARGET, raw)
            for reg, value in zip((UC_ARM64_REG_X0, UC_ARM64_REG_X1, UC_ARM64_REG_X2), (x0, x1, x2)):
                emu.reg_write(reg, value)
            emu.emu_start(candidate.TARGET, candidate.TARGET + 8, count=2)
            expected.append([emu.reg_read(UC_ARM64_REG_X0), emu.reg_read(UC_ARM64_REG_PC)])
        harness = '''package main
import ("encoding/json"; "os")
type state struct { x [31]uint64 }
type rawSourceRuntime struct { stopPC uint64 }
func main() {
 var vectors [][3]uint64; json.NewDecoder(os.Stdin).Decode(&vectors)
 result := [][2]uint64{}
 for _, v := range vectors {
  s := &state{}; s.x[0],s.x[1],s.x[2] = v[0],v[1],v[2]
  next, stopped, err := rawSourceBootstrapBlock100a4be38(s,nil)
  if err != nil || stopped { panic("unexpected stop") }
  result = append(result,[2]uint64{s.x[0],next})
 }
 json.NewEncoder(os.Stdout).Encode(result)
}
'''
        with tempfile.TemporaryDirectory() as folder:
            Path(folder, "candidate.go").write_text(source)
            Path(folder, "main.go").write_text(harness)
            run = subprocess.run(["go", "run", "candidate.go", "main.go"], cwd=folder, input=json.dumps(vectors), text=True, capture_output=True, timeout=90, check=True)
            self.assertEqual(json.loads(run.stdout), expected)


if __name__ == "__main__":
    unittest.main()
