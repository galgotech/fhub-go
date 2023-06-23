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

package fhub

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/galgotech/fhub-go/internal/plugin"
	"github.com/galgotech/fhub-go/internal/rest"
	"github.com/galgotech/fhub-go/model"
)

var fhubPath string

func init() {
	fhubPath = os.Getenv("FHUB_PATH")
	if fhubPath == "" {
		var err error
		fhubPath, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
}

func SetPath(path string) {
	fhubPath = path
}

func Run(functions any) error {
	functionsValueOf := reflect.ValueOf(functions)
	if functionsValueOf.Kind() != reflect.Pointer {
		return errors.New("functions need be a pointer")
	}

	if functionsValueOf.NumMethod() == 0 {
		return errors.New("any functions found")
	}

	fhubModel, err := model.UnmarshalFile(fmt.Sprintf("%s/fhub.cue", fhubPath))
	if err != nil {
		return err
	}

	fHubExec := &plugin.FHubExec{
		Model:     fhubModel,
		Functions: functionsValueOf,
	}

	return rest.Exec(fhubModel, fHubExec)

}
