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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/galgotech/fhub-go/internal/plugin"
	"github.com/galgotech/fhub-go/model"
)

func Exec(fhubModel model.FHub, fhubExec *plugin.FHubExec) error {
	r := gin.Default()

	err := load(r, fhubModel, fhubExec)
	if err != nil {
		return err
	}

	err = r.Run()
	if err != nil {
		return err
	}

	return nil
}

func load(r *gin.Engine, fhub model.FHub, fhubExec *plugin.FHubExec) error {
	for label, function := range fhub.Functions {
		func(label string, function model.Function) {
			path := fmt.Sprintf("%s/%s/%s", fhub.Version, fhub.Name, strings.ToLower(label))
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

				output, err := fhubExec.Exec(label, input)
				if err != nil {
					fmt.Printf("fail pluginExec: %s\n", err)
					c.JSON(http.StatusInternalServerError, nil)
					return
				}
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
