#!/usr/bin/env python3
"""Tests use only synthetic source text, never private enrollment material."""
import importlib.util
from pathlib import Path
import tempfile
import unittest

FILE = Path(__file__).with_name("verify-fairplay-release-coverage.py")
spec = importlib.util.spec_from_file_location("fairplay_release_gate", FILE)
gate = importlib.util.module_from_spec(spec)
spec.loader.exec_module(gate)


class ReleaseGateTests(unittest.TestCase):
    def root(self, handler="", registry="", tests=""):
        directory = tempfile.TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        root = Path(directory.name)
        native = root / "MapleSyrup" / "NativeSAP"
        native.mkdir(parents=True)
        scripts = root / "scripts"
        scripts.mkdir()
        (native / "fairplayrawsupplement.go").write_text(handler)
        (native / "fairplayrawsupplement_test.go").write_text(tests)
        (scripts / "prepare-fairplay-raw.py").write_text(registry)
        return root

    def test_known_missing_target_is_rejected(self):
        root = self.root("package rawsource", "supplement_entry = 0x1008fff2c")
        result = gate.coverage_report(root, 0x100a4be38)
        self.assertFalse(result["ready_for_device_smoke_test"])
        self.assertFalse(result["decryption_validated"])

    def test_comment_does_not_satisfy_function(self):
        root = self.root(
            "// func rawSourceBootstrapBlock100a4be38()\n",
            "supplement_entry = 0x100a4be38\n"
            "blocks[supplement_entry] = 'rawSourceBootstrapBlock100a4be38'",
            "func TestPinned100a4be38Synthetic(t *testing.T) {}",
        )
        self.assertFalse(gate.coverage_report(root, 0x100a4be38)
                         ["ready_for_device_smoke_test"])

    def test_handler_alone_is_not_enough(self):
        root = self.root("func rawSourceBootstrapBlock100a4be38() {}")
        self.assertFalse(gate.coverage_report(root, 0x100a4be38)
                         ["ready_for_device_smoke_test"])

    def test_requires_handler_dispatch_and_synthetic_test(self):
        root = self.root(
            "func rawSourceBootstrapBlock100a4be38() {}",
            "supplement_entry = 0x100a4be38\n"
            "blocks[supplement_entry] = 'rawSourceBootstrapBlock100a4be38'",
            "func TestPinned100a4be38Synthetic(t *testing.T) {}",
        )
        result = gate.coverage_report(root, 0x100a4be38)
        self.assertTrue(result["ready_for_device_smoke_test"])
        # Even passing this *necessary* gate never authenticates FairPlay.
        self.assertFalse(result["decryption_validated"])


if __name__ == "__main__":
    unittest.main()
