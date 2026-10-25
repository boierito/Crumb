#!/usr/bin/env python3
"""Turn pinned standalone evaluators into isolated, testable in-process packages."""
import pathlib, re, sys, hashlib
source, target = map(pathlib.Path, sys.argv[1:])
for stage, tag in [('state','fairplaystate'), ('frame','fairplayframe'), ('sidb','fairplaysidb')]:
    out = target / stage
    out.mkdir(parents=True, exist_ok=True)
    names = ['oracle.go', 'main.go', 'dfa_aes.go', 'dfa_attack.go', 'dfa_finish.go', 'dfa_linear.go',
             'dfa_test.go', 'generated_'+stage+'.go', 'embedded_model_'+stage+'.go', 'embedded_model_'+stage+'_test.go']
    names += ['input_size_state.go', 'generated_state_input.go'] if stage == 'state' else ['input_size_default.go']
    names += ['enrollment_frame_sidb.go'] if stage == 'sidb' else ['enrollment_frame_default.go']
    for name in names:
        text = (source / name).read_text()
        text = re.sub(r'^//go:build[^\n]*\n', '', text)
        text = text.replace('package main', 'package '+stage, 1)
        (out / name).write_text(text)
    (out/'profiles').mkdir(exist_ok=True)
    name = 'macos26a428-'+stage+'-public.model.gz'
    (out/'profiles'/name).write_bytes((source/'profiles'/name).read_bytes())
    (out/'api.go').write_text('''package '''+stage+'''
import "fmt"
// Evaluate runs only the pinned generated operations; it executes no guest code.
func Evaluate(input, frame []byte) ([]byte, error) {
    value, err := loadModelBytes(embeddedModelBytes()); if err != nil { return nil, err }
    if embeddedModelNeedsEnrollmentFrame() {
        if len(frame) != 16 { return nil, fmt.Errorf("enrollment frame required") }
        address, _ := enrollmentFrameAddress(); if err := value.patch(address, frame); err != nil { return nil, err }
    }
    prepared, err := prepareEvaluatorInput(input); if err != nil { return nil, err }
    worker := newState(value); if err := worker.beginBytes(prepared); err != nil { return nil, err }
    if err := runGenerated(worker, fault{}); err != nil { return nil, err }
    return worker.outputBytes(value.block, len(prepared))
}
''')
    if stage == 'sidb':
        with (out/'api.go').open('a') as f:
            f.write('''func Recover(frame []byte, workers int) ([]byte, error) {
    if len(frame) != 16 { return nil, fmt.Errorf("enrollment frame required") }
    value, err := loadModelBytes(embeddedModelBytes()); if err != nil { return nil, err }
    address, _ := enrollmentFrameAddress(); if err := value.patch(address, frame); err != nil { return nil, err }
    server, err := newOracleServer(value, workers); if err != nil { return nil, err }
    recovered, err := dpRecover(&dpTarget{server:server}); if err != nil { return nil, err }
    return json.Marshal(recovered.modelJSON())
}
''')
        p=out/'api.go'; p.write_text(p.read_text().replace('import "fmt"','import ("fmt"; "encoding/json")'))

    api = out / 'api.go'
    text = api.read_text()
    text = text.replace('import "fmt"', 'import ("fmt"; "crypto/sha256")').replace('"encoding/json")', '"encoding/json"; "crypto/sha256")')
    expected = hashlib.sha256((source/'profiles'/name).read_bytes()).hexdigest()
    text += '\nfunc VerifyResource() bool { return fmt.Sprintf("%x", sha256.Sum256(embeddedModelBytes())) == "' + expected + '" }\n'
    api.write_text(text)
