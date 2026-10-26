#!/usr/bin/env python3
"""Generate the whole-section semantic profile on a development host."""
import argparse,importlib.util,pathlib,subprocess,tempfile
ROOT=pathlib.Path(__file__).resolve().parent.parent
spec=importlib.util.spec_from_file_location('exact_daemon',ROOT/'scripts/generate-fairplay-next-block.py')
source=importlib.util.module_from_spec(spec);spec.loader.exec_module(source)
parser=argparse.ArgumentParser(description=__doc__)
parser.add_argument('--compiler',type=pathlib.Path,required=True)
parser.add_argument('--output-directory',type=pathlib.Path,required=True)
args=parser.parse_args();args.output_directory.mkdir(parents=True,exist_ok=True)
with tempfile.TemporaryDirectory(prefix='crumb-exact-generic-') as temporary:
 binary=source.exact_system_binary(temporary)
 subprocess.run([str(args.compiler.resolve()),'--binary',str(binary),'--binary-sha256',source.BINARY_HASH,'--output',str(args.output_directory/'fairplay-generic.ir.gz'),'--report',str(args.output_directory/'fairplay-generic-census.json')],check=True)
