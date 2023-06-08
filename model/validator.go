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
	"context"

	"github.com/go-playground/validator/v10"
)

type ValidatorCtxValueKey string

const ValidatorCtxValue ValidatorCtxValueKey = "value"

func Validator(model FHub) error {
	validate := validator.New()

	validate.RegisterStructValidationCtx(fhubStructLevel, FHub{})
	validate.RegisterStructValidationCtx(functionStructLevel, Function{})

	ctx := context.Background()
	ctx = context.WithValue(ctx, ValidatorCtxValue, model)
	err := validate.StructCtx(ctx, model)

	return err
}

func fhubStructLevel(ctx context.Context, sl validator.StructLevel) {
	// fhub, ok := sl.Current().Interface().(FHub)
	// if !ok {
	// 	return
	// }

}

func functionStructLevel(ctx context.Context, sl validator.StructLevel) {
	function, ok := sl.Current().Interface().(Function)
	if !ok {
		return
	}

	fhub, ok := ctx.Value(ValidatorCtxValue).(FHub)
	if !ok {
		return
	}

	if _, ok := fhub.Packages[function.Package]; !ok {
		sl.ReportError(sl.Current(), "package", "Package", "exists", "")
	}
}
