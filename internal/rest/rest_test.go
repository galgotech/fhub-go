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

package rest

import (
	"testing"

	"github.com/galgotech/fhub-go/model"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	fhub, err := model.UnmarshalFile("../../devenv/test.cue")
	if !assert.NoError(t, err) {
		return
	}

	data := []byte(`{
"arg0": "test",
"arg1":"1",
"arg2": {
  "test1": "valor1", "test2": 1
},
"arg3": ["a"],
"arg4": 1.0
}`)

	dataMap, err := fhub.Functions["test"].UnmarshalInput(data)
	if !assert.NoError(t, err) {
		return
	}

	_, ok := dataMap["arg2"].(map[string]any)["test2"].(int)
	assert.True(t, ok)

	_, ok = dataMap["arg4"].(float64)
	assert.True(t, ok)
}
