import importlib.util
import plistlib
from pathlib import Path
import tempfile
import unittest
import zipfile

spec = importlib.util.spec_from_file_location('prepare', Path(__file__).with_name('prepare-fairplayed-private-ipa.py'))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class PackagingTests(unittest.TestCase):
    def inputs(self, root):
        helper, original = root/'helper.ipa', root/'original.ipa'
        with zipfile.ZipFile(helper, 'w') as z:
            z.writestr('Payload/FairPlayed.app/Info.plist', plistlib.dumps({'CFBundleIdentifier': 'com.boierito.crumb.systemfairplay'}))
            z.writestr('Payload/FairPlayed.app/FairPlayed', b'helper')
        with zipfile.ZipFile(original, 'w') as z:
            z.writestr('Payload/Test.app/Info.plist', plistlib.dumps({'CFBundleExecutable': 'Test'}))
            z.writestr('Payload/Test.app/Test', b'encrypted target remains unchanged')
            z.writestr('Payload/Test.app/SC_Info/Manifest.plist', plistlib.dumps({
                'SinfPaths': ['SC_Info/Test.sinf'], 'SinfReplicationPaths': ['Frameworks/F.framework/SC_Info/F.sinf']}))
            z.writestr('Payload/Test.app/SC_Info/Test.sinf', b'original SINF')
        return helper, original

    def test_preserves_exact_binary_and_sinf_and_replicates_only_manifest_paths(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            helper, original = self.inputs(root)
            output = root/'out.ipa'
            report = module.prepare(helper, original, output, '202610301030')
            self.assertFalse(report['license_validated'])
            with zipfile.ZipFile(output) as z:
                base = 'Payload/FairPlayed.app/EmbeddedPayload/Target/'
                self.assertEqual(z.read(base+'Test'), b'encrypted target remains unchanged')
                self.assertEqual(z.read(base+'SC_Info/Test.sinf'), b'original SINF')
                self.assertEqual(z.read(base+'Frameworks/F.framework/SC_Info/F.sinf'), b'original SINF')
            with self.assertRaises(ValueError):
                module.prepare(helper, original, output, '202610301030')
            self.assertTrue(output.exists())

    def test_rejects_traversal_and_missing_license_without_leaving_output(self):
        for invalid in ('../escape', 'Payload/Test.app/SC_Info/Test.sinf'):
            with tempfile.TemporaryDirectory() as d:
                root = Path(d)
                helper, original = self.inputs(root)
                if invalid == '../escape':
                    with zipfile.ZipFile(original, 'a') as z:
                        z.writestr(invalid, b'bad')
                else:
                    with zipfile.ZipFile(original) as z:
                        entries = {n: z.read(n) for n in z.namelist() if not n.endswith('/Test.sinf')}
                    with zipfile.ZipFile(original, 'w') as z:
                        for name, data in entries.items():
                            z.writestr(name, data)
                output = root/'out.ipa'
                with self.assertRaises(ValueError):
                    module.prepare(helper, original, output, '202610301030')
                self.assertFalse(output.exists())


if __name__ == '__main__':
    unittest.main()
