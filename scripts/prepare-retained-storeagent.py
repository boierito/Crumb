#!/usr/bin/env python3
"""Add an opt-in pre-global-initialization hook to the pinned StoreAgent source.
The ordinary SAP/download constructor keeps exactly its original behavior.
"""
from pathlib import Path
import sys

root = Path(sys.argv[1])
source = root / 'storeagent.go'
text = source.read_text()
signature = 'func openStoreAgentGlobal(ctx context.Context, bundle assets.Bundle, image, hardwareID []byte) (*StoreAgent, uint32, error) {'
init = '\tglobalContext, err := agent.initializeGlobal(hardware)'
if text.count(signature) != 1 or text.count(init) != 1:
    raise SystemExit('pinned StoreAgent initialization contract changed')
replacement = signature + '''
 return openStoreAgentGlobalPrepared(ctx, bundle, image, hardwareID, nil)
}

func openStoreAgentGlobalPrepared(ctx context.Context, bundle assets.Bundle, image, hardwareID []byte, prepare func(*StoreAgent) error) (*StoreAgent, uint32, error) {'''
text = text.replace(signature, replacement, 1)
text = text.replace(init, '''	if prepare != nil {
        if err := prepare(agent); err != nil {
            _ = agent.Close()
            return nil, 0, err
        }
    }
''' + init, 1)
shims = root / 'shims.go'
shimtext = shims.read_text()
if shimtext.count('type shims struct {') != 1 or 'guestStateFS *guestStateFS' in shimtext:
    raise SystemExit('pinned shims contract changed or already prepared')
# Validate both contracts before changing either file.
source.write_text(text)
shims.write_text(shimtext.replace('type shims struct {', 'type shims struct {\n guestStateFS *guestStateFS', 1))
