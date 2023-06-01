package gencode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fhub-runtime-go/model"
)

func Test_gen(t *testing.T) {

	t.Run("use package", func(t *testing.T) {
		fhub := model.Fhub{
			Packages: map[string]model.Package{
				"test": {
					Import: "fhub.dev/test",
				},
			},
			Functions: map[string]model.Function{
				"function_test": {
					Package: "test",
					Launch:  "function_test",
					InputsLabel: []string{
						"arg0",
					},
					InputsType: []string{
						"string",
					},
					OutputsLabel: []string{
						"output1",
					},
					OutputsType: []string{
						"string",
					},
				},
			},
		}

		code, err := gen(fhub)
		assert.NoError(t, err)

		assert.Equal(t, `package main

import pkg_test "fhub.dev/test"

var f = functions{}

func Initialize(env map[string]string, constants map[string]string) error {
	return nil
}
func Exec(function string, input map[string]any) map[string]any {
	switch function {
	case "function_test":
		return f.function_test(input)
	}
	return nil
}

type functions struct{}

func (f *functions) function_test(input map[string]any) map[string]any {
	output1 := pkg_test.function_test(input["arg0"].(string))
	output := map[string]any{"output1": output1}
	return output
}
`, string(code))
	})

	t.Run("use package", func(t *testing.T) {
		fhub := model.Fhub{
			Packages: map[string]model.Package{
				"test": {
					Import: "fhub.dev/test",
					Launch: "Start",
				},
			},
			Functions: map[string]model.Function{
				"function_test": {
					Package: "test",
					Launch:  "function_test",
					InputsLabel: []string{
						"arg0",
					},
					InputsType: []string{
						"string",
					},
					OutputsLabel: []string{
						"output1",
					},
					OutputsType: []string{
						"string",
					},
				},
			},
		}

		code, err := gen(fhub)
		assert.NoError(t, err)

		assert.Equal(t, `package main

import pkg_launchtest "fhub.dev/test"

type interfacetest interface {
	function_test(string) string
}

var test interfacetest
var f = functions{}

func Initialize(env map[string]string, constants map[string]string) error {
	test = (interfacetest)(pkg_launchtest.Start(env))
	return nil
}
func Exec(function string, input map[string]any) map[string]any {
	switch function {
	case "function_test":
		return f.function_test(input)
	}
	return nil
}

type functions struct{}

func (f *functions) function_test(input map[string]any) map[string]any {
	output1 := test.function_test(input["arg0"].(string))
	output := map[string]any{"output1": output1}
	return output
}
`, string(code))
	})

}
