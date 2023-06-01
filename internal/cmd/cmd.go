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

package cmd

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/galgotech/fhub-runtime-go/internal/gencode"
	"github.com/galgotech/fhub-runtime-go/internal/rest"
)

func Gencode() error {
	app := &cli.App{
		Name:  "fhub-gencode",
		Usage: "",
		Authors: []*cli.Author{{
			Name:  "André Miranda",
			Email: "contact@fhub.dev",
		}},
		Action: func(c *cli.Context) (err error) {
			if c.NArg() != 2 {
				return errors.New("without schema")
			}
			return gencode.Exec(c.Args().Get(0), c.Args().Get(1))
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}

func Rest() error {
	app := &cli.App{
		Name:  "fhub-rest",
		Usage: "",
		Authors: []*cli.Author{{
			Name:  "André Miranda",
			Email: "andre@galgo.tech",
		}},
		Action: func(c *cli.Context) (err error) {
			if c.NArg() != 1 {
				return errors.New("without schema or plugin")
			}
			return rest.Exec(c.Args().Get(0))
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}
