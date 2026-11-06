#!/usr/bin/env python3
"""Keep project paths, product identity and first-party copy consistent."""
import pathlib
import re
import subprocess
import xml.etree.ElementTree as ET
root = pathlib.Path(__file__).resolve().parents[1]
old = 'waffle' + 'store'
paths = subprocess.check_output(['git','ls-files','-z'],cwd=root).decode().split('\0')
notices = {'CREDITS.md','THIRD_PARTY_NOTICES.md','Crumb/Resources/ThirdPartyNotices.txt'}
for name in filter(None, paths):
    assert old not in name.lower(), f'Obsolete project path: {name}'
    if name in notices: continue
    try: text = (root/name).read_text()
    except UnicodeError: continue
    if name == 'README.md': text = text.split('## Origins and credits')[0]
    if name == 'index.html':
        start = text.index('<section id="credits">')
        end = text.index('<!-- COMMUNITY -->',start)
        text = text[:start] + text[end:]
    assert old not in text.lower(), f'Obsolete product copy: {name}'
project = (root/'Crumb.xcodeproj/project.pbxproj').read_text()
assert 'path = Crumb;' in project and 'name = Crumb;' in project
assert 'INFOPLIST_FILE = "Crumb-Info.plist";' in project
assert 'Crumb/MapleSyrup/Native/SAP-Bridging-Header.h' in project
scheme = ET.parse(root/'Crumb.xcodeproj/xcshareddata/xcschemes/Crumb.xcscheme')
for ref in scheme.findall('.//BuildableReference'):
    assert ref.get('BlueprintName') == 'Crumb'
    assert ref.get('ReferencedContainer') == 'container:Crumb.xcodeproj'
    assert ref.get('BlueprintIdentifier') in project
print('Crumb branding, project paths and shared scheme verified')
