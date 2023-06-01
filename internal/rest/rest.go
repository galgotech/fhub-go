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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"plugin"

	"github.com/gin-gonic/gin"

	"github.com/galgotech/fhub-runtime/model"
)

func Exec(root string) error {
	r := gin.Default()
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filename := filepath.Base(path)
		if matched, err := filepath.Match("*.cue", filename); err != nil {
			return err
		} else if matched {
			name := filename[:len(filename)-len(filepath.Ext(filename))]
			pluginPath := filepath.Join(root, fmt.Sprintf("%s.so", name))
			err := load(r, path, pluginPath)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = r.Run()
	if err != nil {
		return err
	}

	return nil
}

func load(r *gin.Engine, path, pluginPath string) error {
	fhub, err := model.UnmarshalFile(path)
	if err != nil {
		return err
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return err
	}

	pluginIntializeInterface, err := p.Lookup("Initialize")
	if err != nil {
		return err
	}

	pluginExecInterface, err := p.Lookup("Exec")
	if err != nil {
		return err
	}
	pluginIntialize, ok := pluginIntializeInterface.(func(map[string]string) error)
	if !ok {
		fmt.Printf("%V\n", pluginIntializeInterface)
		return errors.New("invalid interface initialize")
	}
	pluginExec, ok := pluginExecInterface.(func(string, map[string]any) map[string]any)
	if !ok {
		fmt.Printf("%V\n", pluginExecInterface)
		return errors.New("invalid interface exec")
	}

	// TODO: Check the security, in the same runtime, start different clients
	env := map[string]string{
		"DATABASE": os.Getenv("DATABASE"),
	}
	err = pluginIntialize(env)
	if err != nil {
		return err
	}
	for _, function := range fhub.Functions {
		func(function model.Function) {
			path := fmt.Sprintf("%s/%s/%s", fhub.Version, fhub.Name, function.Label)
			r.POST(path, func(c *gin.Context) {
				inputJson, err := ioutil.ReadAll(c.Request.Body)
				if err != nil {
					c.JSON(http.StatusInternalServerError, nil)
					return
				}

				if ok := function.ValidateInput(inputJson); !ok {
					c.JSON(http.StatusBadRequest, nil)
					return
				}

				input := map[string]any{}
				err = json.Unmarshal(inputJson, &input)
				if err != nil {
					c.JSON(http.StatusInternalServerError, nil)
					return
				}

				output := pluginExec(function.Label, input)
				if output == nil {
					c.JSON(http.StatusInternalServerError, nil)
					return
				}

				outputJson, err := json.Marshal(output)
				if err != nil {
					c.JSON(http.StatusBadRequest, nil)
					return
				}

				if ok := function.ValidateOutput(outputJson); !ok {
					c.JSON(http.StatusBadRequest, nil)
					return
				}

				c.JSON(http.StatusOK, output)
			})
		}(function)
	}

	return nil
}