"""Independent ARM64 checks for the optional host-only origin observer."""
import struct
import unittest

from unicorn import Uc, UC_ARCH_ARM64, UC_MODE_ARM, UC_HOOK_CODE, UC_HOOK_MEM_READ, UC_HOOK_MEM_WRITE
from unicorn import arm64_const as arm
from fairplay_candidate_provenance import CandidateProvenance, ByteOrigins, image_segments, canonical_register, BITS, BASE
from keystone import Ks, KS_ARCH_ARM64, KS_MODE_LITTLE_ENDIAN


def image():
    data = bytearray(4096)
    struct.pack_into('<8I', data, 0, 0xFEEDFACF, 0, 0, 2, 2, 144, 0, 0)
    for offset, name, address, protection in ((32, b'__TEXT', BASE, 5), (104, b'__DATA_CONST', BASE+4096, 3)):
        struct.pack_into('<II16sQQQQiiII', data, offset, 0x19, 72, name,
                         address, 4096, 0, 0, 7, protection, 0, 0)
    # ldrb w0,[x1]; eor w0,w0,w2; strb w0,[x3]
    struct.pack_into('<III', data, 256, 0x39400020, 0x4A020000, 0x39000060)
    return bytes(data)


class OriginTests(unittest.TestCase):
    def test_missing_writes_are_not_reported_as_native_material(self):
        audit = CandidateProvenance(image())
        u = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
        u.mem_map(0x7000, 4096)
        u.reg_write(arm.UC_ARM64_REG_X29, 0x7170)
        audit.candidate(u)
        self.assertFalse(audit.report()['runs'][0]['all_bytes_written_by_native'])

    def test_sparse_storage_bound_aborts_instead_of_losing_origins(self):
        memory = ByteOrigins()
        memory.pages = {i: None for i in range(16384)}
        with self.assertRaisesRegex(ValueError, 'memory bound'):
            memory[16384*4096] = (BITS['sidd'], BASE+256)

    def test_sparse_bulk_access_crosses_pages_without_erasing_neighbors(self):
        memory = ByteOrigins()
        memory[4093] = (BITS['sidb'], BASE+256)
        memory.store(4094, 4, BITS['sidd'], BASE+260)
        memory.store(4095, 2, BITS['sinf'], BASE+264)
        fallback = lambda address: BITS['initial_context']
        self.assertEqual(memory.origin(4094, 4, fallback), BITS['sidd'] | BITS['sinf'])
        self.assertEqual(memory[4093], (BITS['sidb'], BASE+256))
        self.assertEqual(memory[4097], (BITS['sidd'], BASE+260))
        self.assertEqual(memory.origin(4098, 1, fallback), BITS['initial_context'])

    def test_macho_permissions_not_image_address_range(self):
        audit = CandidateProvenance(image())
        self.assertEqual(audit.default_origin(BASE+256), BITS['readonly_image'])
        self.assertEqual(audit.default_origin(BASE+4096), BITS['mutable_image'])
        self.assertEqual(audit.default_origin(0x7000), BITS['initial_context'])
        self.assertEqual(canonical_register('w3'), 'x3')
        self.assertEqual(canonical_register('q3'), 'v3')

    def test_actual_arm64_load_transform_store(self):
        audit = CandidateProvenance(image())
        u = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
        u.mem_map(BASE, 4096)
        u.mem_write(BASE, image())
        u.mem_map(0x7000, 4096)
        u.mem_write(0x7000, b'\x13')
        u.reg_write(arm.UC_ARM64_REG_X1, 0x7000)
        u.reg_write(arm.UC_ARM64_REG_X2, 0x24)
        u.reg_write(arm.UC_ARM64_REG_X3, 0x7100)
        audit.registers['x2'] = BITS['native_literal']
        audit.memory[0x7000] = (BITS['sidb'], None)
        u.hook_add(UC_HOOK_CODE, audit.instruction)
        u.hook_add(UC_HOOK_MEM_READ, audit.read)
        u.hook_add(UC_HOOK_MEM_WRITE, audit.write)
        u.emu_start(BASE+256, BASE+268)
        audit.finish()
        self.assertEqual(bytes(u.mem_read(0x7100, 1)), b'\x37')
        origin, writer = audit.memory[0x7100]
        self.assertTrue(origin & BITS['sidb'])
        self.assertTrue(origin & BITS['initial_context'])
        self.assertEqual(writer, BASE+264)
        self.assertFalse(audit.report()['causal_chain_complete'])

    def test_visits_preserved_without_exporting_material(self):
        audit = CandidateProvenance(image())
        u = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
        u.mem_map(0x7000, 4096)
        u.reg_write(arm.UC_ARM64_REG_X29, 0x7170)
        for a in range(0x7100, 0x7110):
            audit.memory[a] = (BITS['sidd'], BASE+264)
        audit.candidate(u)
        audit.candidate(u)
        u.mem_write(0x7100, b'\x01')
        audit.candidate(u)
        report = audit.report()
        self.assertEqual(report['candidate_count'], 2)
        self.assertEqual(report['landmark_visits'], 3)
        self.assertEqual([r['visits'] for r in report['runs']], [2, 1])
        self.assertTrue(report['runs'][0]['all_bytes_written_by_native'])
        self.assertFalse(report['license_validated'])

    def test_whitebox_index_origin_survives_readonly_table_lookup(self):
        native = bytearray(image())
        code, _ = Ks(KS_ARCH_ARM64, KS_MODE_LITTLE_ENDIAN).asm('ldrb w0, [x1, x2]; strb w0, [x3]')
        native[256:256+len(code)] = bytes(code)
        audit = CandidateProvenance(bytes(native))
        u = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
        u.mem_map(BASE, 4096)
        u.mem_write(BASE, bytes(native))
        u.mem_map(0x7000, 4096)
        u.reg_write(arm.UC_ARM64_REG_X1, BASE+512)
        u.reg_write(arm.UC_ARM64_REG_X2, 7)
        u.reg_write(arm.UC_ARM64_REG_X3, 0x7100)
        audit.registers['x1'] = BITS['readonly_image']
        audit.registers['x2'] = BITS['sidb']
        u.hook_add(UC_HOOK_CODE, audit.instruction)
        u.hook_add(UC_HOOK_MEM_READ, audit.read)
        u.hook_add(UC_HOOK_MEM_WRITE, audit.write)
        u.emu_start(BASE+256, BASE+256+len(code))
        audit.finish()
        self.assertTrue(audit.memory_origin(0x7100, 1) & BITS['sidb'])
        self.assertTrue(audit.memory_origin(0x7100, 1) & BITS['readonly_image'])

    def test_invalid_segment_contract_rejected(self):
        data = bytearray(image())
        struct.pack_into('<I', data, 36, 8)
        with self.assertRaises(ValueError):
            image_segments(data)

    def test_read_contract_tags_same_value_bytes_and_preserves_neighbors(self):
        audit = CandidateProvenance(image())
        u = Uc(UC_ARCH_ARM64, UC_MODE_ARM)
        u.mem_map(0x7000, 4096)
        audit.fd[4] = 'sidd'
        args = [0] * 31
        args[:3] = [4, 0x7100, 16]
        returned = args.copy()
        returned[0] = 16
        audit.service(u, {'CPU': {'PC': 0xF100000000008000, 'Registers': args, 'Flags': [False]*4},
                          'After': {'Registers': returned, 'Flags': [False]*4}})
        self.assertEqual(audit.memory_origin(0x7100, 16), BITS['sidd'])
        self.assertEqual(audit.memory_origin(0x7110, 1), BITS['initial_context'])
        self.assertEqual(audit.reg_origin('w0'), BITS['service'])
        self.assertFalse(audit.report()['unknown_service_same_value_writes_verified'])


if __name__ == '__main__':
    unittest.main()
