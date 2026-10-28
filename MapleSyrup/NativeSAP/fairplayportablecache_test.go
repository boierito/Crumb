package evaluator

import (
	"bytes"
	"reflect"
	"sync"
	"testing"
)

// Copied into the public state/frame evaluators only. A shared template must
// never retain one call's input, output, registers or writable pages.
func TestPublicTemplateConcurrentEvaluationsRemainIsolated(t *testing.T) {
	template, err := immutableEvaluatorModel()
	if err != nil {
		t.Fatal(err)
	}
	fresh, err := loadModelBytes(embeddedModelBytes())
	if err != nil {
		t.Fatal(err)
	}
	inputs, expected := make([][]byte, 8), make([][]byte, 8)
	for i := range inputs {
		input := make([]byte, generatedInputBytes())
		for j := range input {
			input[j] = byte(i*37 + j)
		}
		inputs[i] = input
		prepared, err := prepareEvaluatorInput(input)
		if err != nil {
			t.Fatal(err)
		}
		worker := newState(fresh)
		if err := worker.beginBytes(prepared); err != nil {
			t.Fatal(err)
		}
		if err := runGenerated(worker, fault{}); err != nil {
			t.Fatal(err)
		}
		expected[i], err = worker.outputBytes(fresh.block, len(prepared))
		if err != nil {
			t.Fatal(err)
		}
	}
	var wg sync.WaitGroup
	for i := range inputs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			before := bytes.Clone(inputs[i])
			for attempt := 0; attempt < 8; attempt++ {
				output, err := Evaluate(inputs[i], nil)
				if err != nil || !bytes.Equal(output, expected[i]) || !bytes.Equal(inputs[i], before) {
					t.Error("evaluation changed its caller input or used another call's state")
					return
				}
				// Returned output must not alias the template or a later worker.
				clear(output)
			}
		}(i)
	}
	wg.Wait()
	if !reflect.DeepEqual(template, fresh) {
		t.Fatal("evaluation modified the public template")
	}
	if _, err := Evaluate(nil, nil); err == nil {
		t.Fatal("shared template bypassed input validation")
	}
}
