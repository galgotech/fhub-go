// Copyright 2023 The fhub-runtime-go Authors
// This file is part of fhub-runtime-go.
//
// This file is part of fhub-runtime-go.
// fhub-runtime-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// fhub-runtime-go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with fhub-runtime-go. If not, see <https://www.gnu.org/licenses/>.

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

	args := make([]reflect.Value, len(modelFunction.Inputs))
	for i, label := range modelFunction.Inputs {
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
	outputLen := len(modelFunction.Outputs)
	if len(outs) != outputLen {
		return nil, errors.New("invalid output")
	}

	output := make(map[string]any, outputLen)
	for i, out := range outs {
		output[modelFunction.Outputs[i]] = out.Interface()
	}

	return output, nil
}
