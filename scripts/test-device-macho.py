#!/usr/bin/env python3
"""Adversarial parser fixtures; no FairPlay or Apple device simulated."""
import ctypes, pathlib, struct, subprocess, tempfile, unittest
ROOT = pathlib.Path(__file__).resolve().parents[1]
class Range(ctypes.Structure):
    _fields_ = [('cpu', ctypes.c_uint32), ('subtype', ctypes.c_uint32), ('offset', ctypes.c_uint32), ('length', ctypes.c_uint32), ('command', ctypes.c_size_t)]
class ParserTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.temp = tempfile.TemporaryDirectory()
        output = pathlib.Path(cls.temp.name)/'parser.so'
        subprocess.run(['cc','-shared','-fPIC','-Wall','-Wextra','-Werror',str(ROOT/'DecryptHelper/CrumbMachO.c'),'-o',str(output)],check=True)
        cls.parse = ctypes.CDLL(str(output)).crumb_macho
        cls.parse.argtypes = [ctypes.c_void_p,ctypes.c_size_t,ctypes.POINTER(Range)]
    @classmethod
    def tearDownClass(cls): cls.temp.cleanup()
    def run_parser(self, b):
        out=Range(); result=self.parse(b,len(b),ctypes.byref(out));return result,out
    def image(self, cid=1, off=56, size=8):
        return struct.pack('<8I',0xfeedfacf,0x100000c,0,2,1,24,0,0)+struct.pack('<6I',0x2c,24,off,size,cid,0)+bytes(8)
    def test_encrypted_range(self):
        status,r=self.run_parser(self.image());self.assertEqual((status,r.offset,r.length,r.command),(1,56,8,48))
    def test_plain_image(self): self.assertEqual(self.run_parser(self.image(cid=0))[0],0)
    def test_truncated_input(self):
        b=self.image()
        for n in range(4,len(b)): self.assertEqual(self.run_parser(b[:n])[0],-1)
    def test_range_overflow(self): self.assertEqual(self.run_parser(self.image(off=0xfffffff0,size=100))[0],-1)
    def test_invalid_load_size(self):
        b=bytearray(self.image());struct.pack_into('<I',b,36,4);self.assertEqual(self.run_parser(bytes(b))[0],-1)
    def test_fat_is_explicitly_unsupported(self): self.assertEqual(self.run_parser(bytes.fromhex('cafebabe00000001'))[0],-1)
    def test_resource(self): self.assertEqual(self.run_parser(b'resource')[0],0)
if __name__=='__main__': unittest.main()
