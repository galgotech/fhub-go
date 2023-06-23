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

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

type Function struct {
	inputValue  cue.Value `fhub:"input" fhub-unmarshal:"true"`
	outputValue cue.Value `fhub:"output" fhub-unmarshal:"true"`

	InputsLabel  []string       `validate:"min=1"`
	InputsType   []reflect.Kind `validate:"min=1"`
	OutputsLabel []string       `validate:"min=1"`
	OutputsType  []reflect.Kind `validate:"min=1"`
}

func (f *Function) Unmarshal(field string, value cue.Value) (err error) {
	if field == "input" {
		f.inputValue = value
		f.InputsLabel, f.InputsType, err = type_struct(f.inputValue)
		if err != nil {
			return err
		}
	} else if field == "output" {
		f.outputValue = value
		f.OutputsLabel, f.OutputsType, err = type_struct(f.outputValue)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid field")
	}

	return nil
}

func (f Function) ValidateInput(data []byte) bool {
	return f.validate(data, f.inputValue)
}

func (f Function) ValidateOutput(data []byte) bool {
	return f.validate(data, f.outputValue)
}

func (f *Function) validate(data []byte, value cue.Value) bool {
	valid := json.Valid(data)
	if !valid {
		return false
	}

	dataValue := cuecontext.New().CompileBytes(data)
	valueUnify := value.UnifyAccept(dataValue, value)
	if valueUnify.Err() != nil {
		return false
	}

	err := valueUnify.Validate(
		cue.Attributes(true),
		cue.Optional(true),
		cue.Hidden(true),
		cue.Concrete(true),
	)

	return err == nil
}

func type_struct(value cue.Value) ([]string, []reflect.Kind, error) {
	fields, err := value.Fields()
	if err != nil {
		return nil, nil, err
	}

	labels := []string{}
	values := []reflect.Kind{}
	for fields.Next() {
		fieldValue := fields.Value()
		label := fields.Label()
		labels = append(labels, label)

		switch fieldValue.IncompleteKind() {
		case cue.StringKind:
			values = append(values, reflect.String)
		case cue.IntKind:
			values = append(values, reflect.Int)
		case cue.BoolKind:
			values = append(values, reflect.Bool)
		default:
			return nil, nil, fmt.Errorf("invalid type input %s", fieldValue.IncompleteKind())
		}
	}

	return labels, values, nil
}
