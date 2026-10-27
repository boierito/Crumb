#!/usr/bin/env python3
import pathlib
import subprocess
import tempfile
import unittest

SCRIPT = pathlib.Path(__file__).with_name('prepare-retained-storeagent.py')
SIGNATURE = 'func openStoreAgentGlobal(ctx context.Context, bundle assets.Bundle, image, hardwareID []byte) (*StoreAgent, uint32, error) {'
SOURCE = SIGNATURE + '\n\tglobalContext, err := agent.initializeGlobal(hardware)\n return agent, globalContext, err\n}\n'
SHIMS = 'package machine\ntype shims struct {\n}\n'

class PrepareRetainedStoreAgentTests(unittest.TestCase):
 def run_patch(self, root):
  return subprocess.run(['python3',str(SCRIPT),str(root)],capture_output=True)
 def test_hook_precedes_initializer_and_default_has_no_hook(self):
  with tempfile.TemporaryDirectory() as tmp:
   root=pathlib.Path(tmp)
   (root/'storeagent.go').write_text(SOURCE); (root/'shims.go').write_text(SHIMS)
   self.assertEqual(self.run_patch(root).returncode,0)
   result=(root/'storeagent.go').read_text()
   self.assertIn('return openStoreAgentGlobalPrepared(ctx, bundle, image, hardwareID, nil)',result)
   self.assertLess(result.index('prepare(agent)'),result.index('agent.initializeGlobal(hardware)'))
   self.assertIn('_ = agent.Close()',result)
   self.assertIn('guestStateFS *guestStateFS',(root/'shims.go').read_text())
   # Applying the overlay twice must fail without changing prepared files.
   self.assertNotEqual(self.run_patch(root).returncode,0)
   self.assertEqual(result,(root/'storeagent.go').read_text())
 def test_contract_drift_does_not_partially_rewrite_sources(self):
  for broken_source,broken_shims in [(SOURCE.replace('initializeGlobal','unknownGlobal'),SHIMS),(SOURCE,'package machine\n')]:
   with tempfile.TemporaryDirectory() as tmp:
    root=pathlib.Path(tmp)
    (root/'storeagent.go').write_text(broken_source);(root/'shims.go').write_text(broken_shims)
    self.assertNotEqual(self.run_patch(root).returncode,0)
    self.assertEqual((root/'storeagent.go').read_text(),broken_source)
    self.assertEqual((root/'shims.go').read_text(),broken_shims)

if __name__ == '__main__': unittest.main()
