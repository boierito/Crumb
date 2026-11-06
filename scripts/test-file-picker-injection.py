import importlib.util
import pathlib
import struct
import unittest
spec=importlib.util.spec_from_file_location('injection',pathlib.Path(__file__).with_name('embed-file-picker-fix.py'))
injection=importlib.util.module_from_spec(spec);spec.loader.exec_module(injection)
class Tests(unittest.TestCase):
    def fixture(self):
        data=bytearray(1024)
        struct.pack_into('<8I',data,0,0xfeedfacf,0x0100000c,0,2,1,152,0,0)
        struct.pack_into('<II',data,32,0x19,152)
        struct.pack_into('<I',data,32+64,1)
        struct.pack_into('<I',data,32+72+48,512)
        data[512:]=b'x'*512
        return data
    def test_load_command_preserves_sections_and_is_idempotent(self):
        source=self.fixture();result=injection.inject(source)
        self.assertEqual(result[512:],source[512:])
        self.assertEqual(struct.unpack_from('<I',result,16)[0],2)
        self.assertIn(injection.LOAD_PATH,result[:512])
        self.assertEqual(injection.inject(result),result)
    def test_refuses_to_overwrite_nonzero_header(self):
        source=self.fixture();source[184]=1
        with self.assertRaises(ValueError):injection.inject(source)
    def test_rejects_invalid_arch_and_commands(self):
        with self.assertRaises(ValueError):injection.inject(b'bad')
        source=self.fixture();struct.pack_into('<I',source,36,0)
        with self.assertRaises(ValueError):injection.inject(source)
unittest.main()
