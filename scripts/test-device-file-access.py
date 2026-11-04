#!/usr/bin/env python3
"""Actual OS file-access errors; no iOS policy or FairPlay is simulated."""
import ctypes
import errno
import fcntl
import os
from pathlib import Path
import stat
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[1]


class Access(ctypes.Structure):
    _fields_ = [("error_code", ctypes.c_int), ("metadata_available", ctypes.c_int),
                ("mode", ctypes.c_uint32), ("owner", ctypes.c_uint32), ("group", ctypes.c_uint32)]


class FileAccessTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.library_temp = tempfile.TemporaryDirectory()
        library = Path(cls.library_temp.name) / "file-access.so"
        subprocess.run(["cc", "-shared", "-fPIC", "-Wall", "-Wextra", "-Werror",
                        str(ROOT / "DecryptHelper/CrumbFileAccess.c"), "-o", str(library)], check=True)
        cls.open_image = ctypes.CDLL(str(library), use_errno=True).crumb_open_installed_image
        cls.open_image.argtypes = [ctypes.c_char_p, ctypes.POINTER(Access)]
        cls.open_image.restype = ctypes.c_int

    @classmethod
    def tearDownClass(cls):
        cls.library_temp.cleanup()

    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.path = Path(self.temp.name) / "image"

    def open(self, path):
        status = Access()
        fd = self.open_image(os.fsencode(path), ctypes.byref(status))
        saved_errno = ctypes.get_errno()
        if fd >= 0:
            self.addCleanup(os.close, fd)
        return fd, status, saved_errno

    def test_read_only_descriptor_and_metadata(self):
        self.path.write_bytes(b"image bytes")
        fd, result, _ = self.open(self.path)
        self.assertGreaterEqual(fd, 0)
        self.assertEqual(os.read(fd, 11), b"image bytes")
        self.assertEqual(fcntl.fcntl(fd, fcntl.F_GETFL) & os.O_ACCMODE, os.O_RDONLY)
        self.assertTrue(fcntl.fcntl(fd, fcntl.F_GETFD) & fcntl.FD_CLOEXEC)
        self.assertEqual(result.error_code, 0)
        self.assertEqual(result.metadata_available, 1)
        self.assertEqual(result.owner, os.geteuid())
        self.assertTrue(stat.S_ISREG(result.mode))

    def test_missing_file_preserves_open_errno(self):
        fd, result, saved = self.open(self.path)
        self.assertEqual((fd, result.error_code, saved), (-1, errno.ENOENT, errno.ENOENT))
        self.assertEqual(result.metadata_available, 0)

    def test_final_symlink_is_rejected_with_link_metadata(self):
        target = self.path.with_name("target")
        target.write_bytes(b"not followed")
        self.path.symlink_to(target)
        fd, result, saved = self.open(self.path)
        self.assertEqual((fd, result.error_code, saved), (-1, errno.ELOOP, errno.ELOOP))
        self.assertTrue(stat.S_ISLNK(result.mode))
        self.assertEqual(result.metadata_available, 1)

    @unittest.skipIf(os.geteuid() == 0, "root bypasses ordinary read permissions")
    def test_permission_denied_preserves_mode_and_original_errno(self):
        self.path.write_bytes(b"protected")
        self.path.chmod(0)
        self.addCleanup(self.path.chmod, 0o600)
        fd, result, saved = self.open(self.path)
        self.assertEqual((fd, result.error_code, saved), (-1, errno.EACCES, errno.EACCES))
        self.assertEqual(result.metadata_available, 1)
        self.assertEqual(stat.S_IMODE(result.mode), 0)


if __name__ == "__main__":
    unittest.main()
