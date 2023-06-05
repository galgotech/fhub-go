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
	"path"
	"path/filepath"
	"plugin"

	"github.com/gin-gonic/gin"

	"github.com/galgotech/fhub-runtime-go/model"
)

func Exec(root string) error {
	r := gin.Default()

	fhubPath := path.Join(root, "fhub.cue")
	pluginPath := filepath.Join(root, "plugin.so")
	err := load(r, fhubPath, pluginPath)
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
		return fmt.Errorf("load plugin: %w", err)
	}

	pluginIntializeInterface, err := p.Lookup("Initialize")
	if err != nil {
		return fmt.Errorf("plugin lookup Initialize: %w", err)
	}

	pluginExecInterface, err := p.Lookup("Exec")
	if err != nil {
		return fmt.Errorf("plugin lookup Exec: %w", err)
	}
	pluginIntialize, ok := pluginIntializeInterface.(func(map[string]string, map[string]string) error)
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
	env := map[string]string{}
	for _, name := range fhub.Env {
		if value, ok := os.LookupEnv(name); ok {
			env[name] = value
		}
	}

	err = pluginIntialize(env, fhub.Constants)
	if err != nil {
		return err
	}
	for label, function := range fhub.Functions {
		func(label string, function model.Function) {
			path := fmt.Sprintf("%s/%s/%s", fhub.Version, fhub.Name, label)
			r.POST(path, func(c *gin.Context) {

				inputJson, err := ioutil.ReadAll(c.Request.Body)
				if err != nil {
					fmt.Printf("fail read json input: %s\n", err)
					c.JSON(http.StatusInternalServerError, nil)
					return
				}

				if ok := function.ValidateInput(inputJson); !ok {
					fmt.Printf("fail validate input\n")
					c.JSON(http.StatusBadRequest, nil)
					return
				}

				input := map[string]any{}
				err = json.Unmarshal(inputJson, &input)
				if err != nil {
					fmt.Printf("fail unmarshal input: %s\n", err)
					c.JSON(http.StatusInternalServerError, nil)
					return
				}

				output := pluginExec(label, input)
				if output == nil {
					fmt.Printf("fail pluginExec\n")
					c.JSON(http.StatusInternalServerError, nil)
					return
				}

				outputJson, err := json.Marshal(output)
				if err != nil {
					fmt.Printf("fail marshal input: %s\n", err)
					c.JSON(http.StatusBadRequest, nil)
					return
				}

				if ok := function.ValidateOutput(outputJson); !ok {
					fmt.Printf("fail validate output\n")
					c.JSON(http.StatusBadRequest, nil)
					return
				}

				c.JSON(http.StatusOK, output)
			})
		}(label, function)
	}

	return nil
}
