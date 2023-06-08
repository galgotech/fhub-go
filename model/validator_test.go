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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var structFhubDefault = FHub{
	Name:        "namespace",
	Version:     "0.1",
	SpecVersion: "0.1",
	Build: Build{
		Local: &Local{
			Source: "./",
		},
	},
	Functions: map[string]Function{
		"fuctionTest": {
			// Package:      "test",
			// Launch:       "Test",
			InputsLabel:  []string{"arg1"},
			InputsType:   []reflect.Kind{reflect.String},
			OutputsLabel: []string{"out1"},
			OutputsType:  []reflect.Kind{reflect.String},
		},
	},
}

type validatorTestCase struct {
	Name  string
	Model FHub
	Err   string
}

func TestValidator(t *testing.T) {
	testCases := []validatorTestCase{{
		Name:  "success",
		Model: structFhubDefault,
		Err:   ``,
	}, {
		Name: "",
		Model: func() FHub {
			f := structFhubDefault
			f.Name = ""
			f.Version = ""
			f.SpecVersion = ""
			return f
		}(),
		Err: `Key: 'FHub.Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'FHub.Version' Error:Field validation for 'Version' failed on the 'required' tag
Key: 'FHub.SpecVersion' Error:Field validation for 'SpecVersion' failed on the 'required' tag`,
	}, {
		Name: "build required",
		Model: func() FHub {
			f := structFhubDefault.DeepCopy()
			f.Build.Local = nil
			return *f
		}(),
		Err: `Key: 'FHub.Build.Local' Error:Field validation for 'Local' failed on the 'required_without' tag
Key: 'FHub.Build.Container' Error:Field validation for 'Container' failed on the 'required_without' tag`,
	}}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := Validator(tc.Model)
			if tc.Err != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.Err, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
