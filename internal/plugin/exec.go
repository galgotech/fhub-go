package plugin

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/galgotech/fhub-go/model"
)

type FHubExec struct {
	Model     model.FHub
	Functions reflect.Value
}

func (f *FHubExec) Exec(function string, input map[string]any) (map[string]any, error) {
	modelFunction, ok := f.Model.Functions[function]
	if !ok {
		return nil, errors.New("function not found")
	}

	args := make([]reflect.Value, len(modelFunction.InputsLabel))
	for i, label := range modelFunction.InputsLabel {
		if val, ok := input[label]; ok {
			args[i] = reflect.ValueOf(val)
		} else {
			return nil, fmt.Errorf("arg not found %q", label)
		}
	}

	method := f.Functions.MethodByName(function)
	if !method.IsValid() {
		return nil, fmt.Errorf("function not implemented %q", function)
	}
	outs := method.Call(args)
	outputLen := len(modelFunction.OutputsLabel)
	if len(outs) != outputLen {
		return nil, errors.New("invalid output")
	}

	output := make(map[string]any, outputLen)
	for i, out := range outs {
		output[modelFunction.OutputsLabel[i]] = out.Interface()
	}

	return output, nil
}
