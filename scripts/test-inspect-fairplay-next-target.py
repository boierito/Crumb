#!/usr/bin/env python3
"""Public synthetic-only tests for sanitized next-target instruction census."""
import importlib.util
from pathlib import Path
import unittest

SOURCE = Path(__file__).with_name("inspect-fairplay-next-target.py")
spec = importlib.util.spec_from_file_location("next_target_inspector", SOURCE)
assert spec and spec.loader
inspector = importlib.util.module_from_spec(spec)
spec.loader.exec_module(inspector)

class NextTargetTests(unittest.TestCase):
    def test_synthetic_indirect_branch(self):
        output = "\n".join([
            "100a4be38: 52800000    mov w0, #0",
            "100a4be3c: f9400000    ldr x0, [x0]",
            "100a4be40: d71f0800    braa x0, x0",
        ])
        got = inspector.summarize_disassembly(output, 0x100a4be38, 0x100a4be44)
        self.assertEqual(got, {
            "target_is_disassemblable": True,
            "instructions_observed": 3,
            "branch_instructions_observed": 1,
            "indirect_branches_observed": 1,
            "first_branch_kind": "indirect",
            "sequential_word_alignment": True
        })

    def test_missing_first_instruction_is_not_covered(self):
        got = inspector.summarize_disassembly("100a4be3c: 52800000  mov w0, #0",
                                             0x100a4be38, 0x100a4be44)
        self.assertFalse(got["target_is_disassemblable"])

    def test_never_saves_instruction_bytes(self):
        got = inspector.summarize_disassembly("100a4be38: 11223344  b 0x100a4be50",
                                             0x100a4be38, 0x100a4be48)
        for value in got.values():
            self.assertNotIn("11223344", str(value))
            self.assertNotIn("100a4be50", str(value))

    def test_branch_kinds(self):
        cases = {"braa": "indirect", "retab": "indirect", "b": "direct",
                 "b.eq": "conditional", "cbnz": "conditional", "bl": "call",
                 "mov": ""}
        for mnemonic, expected in cases.items():
            with self.subTest(mnemonic=mnemonic):
                self.assertEqual(inspector.branch_kind(mnemonic), expected)

if __name__ == "__main__":
    unittest.main()
