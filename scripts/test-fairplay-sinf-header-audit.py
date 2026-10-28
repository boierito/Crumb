"""Compare header predicate formulas with the original pinned ARM64 bytes."""
import gzip
import json
import hashlib
import pathlib
import random
import unittest
from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM, arm64_const as arm
from fairplay_image_resource import read_image_text
from fairplay_sinf_header_audit import HeaderAudit, expected_word, encoded_word, DECODE_PC, WORD_PC, COMPARE_PC, ERROR_PC, MobileHeaderComparison, MATERIAL_PC, KEY_MASK

ROOT = pathlib.Path(__file__).resolve().parent.parent


class HeaderContracts(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        compressed, pin = read_image_text(ROOT)
        cls.native = gzip.decompress(compressed)
        assert hashlib.sha256(cls.native).hexdigest() == pin['immutable_text']['sha256']

    def engine(self, pc, end):
        u = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
        u.mem_map(pc & ~4095, 4096)
        u.mem_write(pc, self.native[pc - 0x100000000:end - 0x100000000])
        return u

    def test_word_transform_against_native(self):
        start, end = 0x1004CCE3C, 0x1004CCE50
        u = self.engine(start, end)
        rng = random.Random(23053)
        for word in [0, 0xffffffff, int.from_bytes(b'itun', 'big')] + [rng.getrandbits(32) for _ in range(2048)]:
            for register, value in [(arm.UC_ARM64_REG_X10, word), (arm.UC_ARM64_REG_X8, word & 255), (arm.UC_ARM64_REG_X6, 0x4e757bd5), (arm.UC_ARM64_REG_X7, 0xce757bc5)]:
                u.reg_write(register, value)
            u.emu_start(start, end, count=5)
            self.assertEqual(u.reg_read(arm.UC_ARM64_REG_X2), encoded_word(word))

    def test_expected_word_against_native(self):
        start, end = 0x1004CB590, 0x1004CB5C0
        u = self.engine(start, end)
        u.mem_map(0x200000, 4096)
        u.reg_write(arm.UC_ARM64_REG_SP, 0x200000)
        rng = random.Random(42057)
        for scalar in [0, 0xffffffff] + [rng.getrandbits(32) for _ in range(2048)]:
            u.mem_write(0x200dd8, scalar.to_bytes(4, 'little'))
            u.emu_start(start, end, count=12)
            self.assertEqual(u.reg_read(arm.UC_ARM64_REG_X9), expected_word(scalar))

    def test_audit_never_infers_decoding_from_error_alone(self):
        audit = HeaderAudit()
        audit.observe(ERROR_PC, [0] * 31, 0, lambda a, n: bytes(n))
        self.assertEqual(audit.report()['decisions'], [])
        audit.observe(COMPARE_PC, [0] * 31, 0, lambda a, n: bytes(n))
        row = audit.report()['decisions'][0]
        self.assertFalse(row['decoded_word_observed'])
        self.assertFalse(row['word_transform_matches_native'])
        self.assertFalse(audit.report()['license_validated_by_audit'])

    def test_header_read_and_transform_are_correlated(self):
        # Public synthetic four-byte buffer; the observer has no write callback.
        r = [0] * 31
        r[3], r[4], r[5], r[20], r[26] = 1, 2, 3, 0, 27
        buffer = bytes(x ^ 27 for x in b'itun')
        audit = HeaderAudit()
        audit.observe(DECODE_PC, r, 0, lambda a, n: buffer[a:a+n])
        r[2] = encoded_word(int.from_bytes(b'itun', 'big'))
        audit.observe(WORD_PC, r, 0, None)
        audit.observe(COMPARE_PC, r, 0, lambda a, n: bytes(n))
        row = audit.report()['decisions'][0]
        self.assertTrue(row['decoded_header_is_itun'])
        self.assertTrue(row['word_transform_matches_native'])
        self.assertFalse(audit.report()['license_validated_by_audit'])

    def test_mobile_first_block_against_independent_public_vectors(self):
        vectors = json.loads((ROOT / 'MapleSyrup/NativeSAP/mobilelicense/testdata/reference.json').read_text())
        tested = 0
        for vector in vectors:
            if vector.get('reject'):
                continue
            sinf, credential = bytes.fromhex(vector['sinf']), bytes.fromhex(vector['credential'])
            try:
                comparison = MobileHeaderComparison(sinf)
            except ValueError:
                continue
            material = bytes(a ^ b ^ 0x30 for a, b in zip(credential, KEY_MASK))
            registers = [0] * 31
            registers[29] = 0x1000
            def read_material(address, size):
                self.assertEqual((address, size), (0x1000 - 0x70, 16))
                return material
            comparison.observe(MATERIAL_PC, registers, 0, read_material)
            self.assertEqual(comparison.candidates[0], b'itun' + bytes(12))
            registers[3], registers[4], registers[5], registers[26] = 1, 2, 3, 27
            encoded = bytes(x ^ 27 for x in comparison.candidates[0])
            comparison.observe(DECODE_PC, registers, 0, lambda address, size: encoded[address:address+size])
            report = comparison.report()
            self.assertEqual(report['comparisons'][0]['native_first_block_matches_candidates'], [True])
            self.assertFalse(report['license_validated_by_comparison'])
            tested += 1
        self.assertGreaterEqual(tested, 4)



if __name__ == '__main__':
    unittest.main()
