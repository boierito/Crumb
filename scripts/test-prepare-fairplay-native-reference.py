"""Verify the host adapter refuses upstream drift without partial edits."""

import pathlib
import subprocess
import sys
import tempfile
import unittest

SCRIPT = pathlib.Path(__file__).with_name("prepare-fairplay-native-reference.py")
FLOW = """package rawsource
func runNativeRegion(s *state, runtime *rawSourceRuntime, pc uint64, region nativeRegion) (uint64, error) {
\thits := 0
\t\t\tnext, err := runtime.handle(s, pc)
}
"""
STORES = """package rawsource
func sourceOutputStore(s *state, address uint64, data []byte) error {
}
"""


class AdapterContracts(unittest.TestCase):
    def attempt(self, directory):
        return subprocess.run(
            [sys.executable, str(SCRIPT), str(directory)], capture_output=True
        )

    def fixture(self, directory, flow=FLOW, stores=STORES):
        (directory / "native_flow.go").write_text(flow)
        (directory / "source_outputs.go").write_text(stores)

    def test_drift_does_not_partially_instrument(self):
        for flow, stores in (
            (FLOW.replace("hits := 0", "hits := 1"), STORES),
            (FLOW, STORES.replace("data []byte", "data []uint8")),
        ):
            with self.subTest(flow=flow), tempfile.TemporaryDirectory() as temporary:
                directory = pathlib.Path(temporary)
                self.fixture(directory, flow, stores)
                self.assertNotEqual(self.attempt(directory).returncode, 0)
                self.assertEqual((directory / "native_flow.go").read_text(), flow)
                self.assertEqual((directory / "source_outputs.go").read_text(), stores)
                self.assertFalse((directory / "native_reference.go").exists())

    def test_single_install_then_refuse_duplicate(self):
        with tempfile.TemporaryDirectory() as temporary:
            directory = pathlib.Path(temporary)
            self.fixture(directory)
            self.assertEqual(self.attempt(directory).returncode, 0)
            installed = {p.name: p.read_bytes() for p in directory.iterdir()}
            self.assertNotEqual(self.attempt(directory).returncode, 0)
            self.assertEqual(
                installed, {p.name: p.read_bytes() for p in directory.iterdir()}
            )


if __name__ == "__main__":
    unittest.main()
