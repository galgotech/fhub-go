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

package gencode

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"

	"github.com/galgotech/fhub-runtime-go/model"
)

func Exec(root string, rootOutput string) error {
	fhubFile := filepath.Join(root, "fhub.cue")

	_, err := os.Stat(fhubFile)
	if err != nil {
		return err
	}

	fhub, err := model.UnmarshalFile(fhubFile)
	if err != nil {
		return err
	}

	code, err := gen(fhub)
	if err != nil {
		return err
	}

	pathGen := filepath.Join(rootOutput, "main.go")
	err = os.WriteFile(pathGen, code, 0760)
	if err != nil {
		return err
	}

	return nil
}

func gen(fhub model.FHub) ([]byte, error) {
	f := jen.NewFile("main")
	f.PackagePrefix = "pkg"

	for label, pkg := range fhub.Packages {
		if pkg.HasLaunch() {
			interfaceName := fmt.Sprintf("interface%s", label)
			methods := make([]jen.Code, len(fhub.Functions))
			i := 0
			for _, function := range fhub.Functions {
				inputs := make([]jen.Code, len(function.InputsType))
				for j, inputType := range function.InputsType {
					inputs[j] = jen.Id(inputType)
				}

				outputs := make([]jen.Code, len(function.InputsLabel))
				for j, outputType := range function.OutputsType {
					outputs[j] = jen.Id(outputType)
				}

				methods[i] = jen.Id(function.Launch).Params(
					inputs...,
				).Params(outputs...)
				i++
			}

			f.Type().Id(interfaceName).Interface(methods...)

			f.ImportAlias(pkg.Import, fmt.Sprintf("launch%s", label))
			f.Var().Id(label).Id(interfaceName)
		}
	}

	f.Var().Id("f").Op("=").Id("functions").Values()

	f.Func().Id("Initialize").Params(
		jen.Id("env").Map(jen.String()).String(),
		jen.Id("constants").Map(jen.String()).String(),
	).Id("error").BlockFunc(
		func(g *jen.Group) {
			g.Var().Add(jen.Id("err"), jen.Error())
			for label, pkg := range fhub.Packages {
				if pkg.HasLaunch() {
					g.Add(
						jen.List(jen.Id(label), jen.Id("err")).Op("=").Qual(pkg.Import, pkg.Launch).Call(
							jen.Id("env"),
							jen.Id("constants"),
						),
					)
					g.Add(
						jen.If(jen.Id("err").Op("!=").Nil()).Block(
							jen.Return(jen.Id("err")),
						),
					)
				}
			}
			g.Add(jen.Return(jen.Nil()))
		},
	)

	f.Func().Id("Exec").Params(
		jen.Id("function").String(), jen.Id("input").Map(jen.String()).Any(),
	).Map(jen.String()).Any().Block(
		jen.Switch(jen.Id("function").BlockFunc(func(g *jen.Group) {
			for label, function := range fhub.Functions {
				g.Add(jen.Case(jen.Lit(label)).Block(
					jen.Return(jen.Id("f").Dot(function.Launch).Call(jen.Id("input"))),
				))
			}
		})),
		jen.Return(jen.Nil()),
	)

	f.Type().Id("functions").Struct()

	for _, function := range fhub.Functions {
		functionArgs := make([]jen.Code, len(function.InputsLabel))
		for i := range function.InputsLabel {
			var t jen.Code
			switch function.InputsType[i] {
			case "string":
				t = jen.String()
			case "bool":
				t = jen.Bool()
			case "int":
				t = jen.Int()
			}
			functionArgs[i] = jen.Id("input").Index(jen.Lit(function.InputsLabel[i])).Assert(t)
		}

		f.Func().Params(
			jen.Id("f").Id("*functions"),
		).Id(function.Launch).Params(
			jen.Id("input").Map(jen.String()).Any(),
		).Map(jen.String()).Any().Block(
			func() jen.Code {
				v := jen.ListFunc(func(g *jen.Group) {
					for _, id := range function.OutputsLabel {
						g.Add(jen.Id(id))
					}
				}).Op(":=")

				if pkg, ok := fhub.Packages[function.Package]; ok {
					if pkg.HasLaunch() {
						v = v.Id(function.Package).Dot(function.Launch)
					} else {
						v = v.Qual(pkg.Import, function.Launch)
					}
				}

				v = v.Call(functionArgs...)
				return v
			}(),
			jen.Id("output").Op(":=").Map(jen.String()).Any().Values(jen.DictFunc(
				func(dict jen.Dict) {
					for _, id := range function.OutputsLabel {
						dict[jen.Lit(id)] = jen.Id(id)
					}
				},
			)),
			jen.Return(jen.Id("output")),
		)
	}

	buf := &bytes.Buffer{}
	err := f.Render(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
