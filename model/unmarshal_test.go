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

func TestUnmarshal(t *testing.T) {
	t.Run("devenv/test.cue", func(t *testing.T) {
		fhub, err := UnmarshalFile("../devenv/test.cue")
		assert.NoError(t, err)

		assert.Equal(t, "test", fhub.Name)
		assert.Equal(t, "1.0", fhub.SpecVersion)
		assert.Equal(t, "v1", fhub.Version)
		assert.Equal(t, map[string]string{"const_name": "test"}, fhub.Constants)
		assert.Equal(t, []string{"NAME"}, fhub.Env)
		assert.Equal(t, []string{"fhub/internaltest.cue"}, fhub.Import)

		assert.Equal(t, "./", fhub.Build.Local.Source)
		assert.Equal(t, "https://fhub.dev/test", fhub.Serving.Http.Url)

		assert.Equal(t, "arg0", fhub.Functions["test"].InputsLabel[0])
		assert.Equal(t, "arg1", fhub.Functions["test"].InputsLabel[1])
		assert.Equal(t, reflect.String, fhub.Functions["test"].InputsType[0])
		assert.Equal(t, reflect.String, fhub.Functions["test"].InputsType[1])
		assert.Equal(t, "ok", fhub.Functions["test"].OutputsLabel[0])
		assert.Equal(t, reflect.Bool, fhub.Functions["test"].OutputsType[0])

		ok := fhub.Functions["test"].ValidateInput([]byte(`{"arg0": "test", "arg1": "test2"}`))
		assert.True(t, ok)
		ok = fhub.Functions["test"].ValidateOutput([]byte(`{"ok": true}`))
		assert.True(t, ok)

		ok = fhub.Functions["test"].ValidateInput([]byte(`{"arg0": "test", "arg2": "invalid"}`))
		assert.False(t, ok)
		ok = fhub.Functions["test"].ValidateOutput([]byte(`{"ok": "invalid"}`))
		assert.False(t, ok)
	})

	t.Run("devenv/test_containerfile.cue", func(t *testing.T) {
		fhub, err := UnmarshalFile("../devenv/test_containerfile.cue")
		assert.NoError(t, err)

		assert.Equal(t, "test", fhub.Name)
		assert.Equal(t, "1.0", fhub.SpecVersion)
		assert.Equal(t, "v1", fhub.Version)
		assert.Equal(t, []string{"fhub/internaltest.cue"}, fhub.Import)

		assert.Equal(t, "Containerfile", fhub.Build.Container.ContainerFile)
		assert.Equal(t, "/app", fhub.Build.Container.Source)
		assert.Equal(t, "https://fhub.dev/test", fhub.Serving.Http.Url)

		assert.Equal(t, "arg0", fhub.Functions["test"].InputsLabel[0])
		assert.Equal(t, "arg1", fhub.Functions["test"].InputsLabel[1])
		assert.Equal(t, reflect.String, fhub.Functions["test"].InputsType[0])
		assert.Equal(t, reflect.String, fhub.Functions["test"].InputsType[1])
		assert.Equal(t, "ok", fhub.Functions["test"].OutputsLabel[0])
		assert.Equal(t, reflect.Bool, fhub.Functions["test"].OutputsType[0])

		ok := fhub.Functions["test"].ValidateInput([]byte(`{"arg0": "test", "arg1": "test2"}`))
		assert.True(t, ok)
		ok = fhub.Functions["test"].ValidateOutput([]byte(`{"ok": true}`))
		assert.True(t, ok)

		ok = fhub.Functions["test"].ValidateInput([]byte(`{"arg0": "test", "arg2": "invalid"}`))
		assert.False(t, ok)
		ok = fhub.Functions["test"].ValidateOutput([]byte(`{"ok": "invalid"}`))
		assert.False(t, ok)
	})
}
