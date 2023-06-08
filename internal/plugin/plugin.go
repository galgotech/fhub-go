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

package plugin

import (
	"errors"
	"os"
	"os/exec"
	"reflect"

	"github.com/galgotech/fhub-runtime-go/model"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "plugin",
	MagicCookieValue: "fhub",
}

func Client(pluginName, pluginPath string) (*plugin.Client, FHub, error) {
	var pluginMap = map[string]plugin.Plugin{
		pluginName: &FHubPlugin{},
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(pluginPath),
		Logger:          logger,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, err
	}

	// TODO: Check the security, in the same runtime, start different clients
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		return nil, nil, err
	}

	fhub, ok := raw.(FHub)
	if !ok {
		return nil, nil, errors.New("invalid interface")
	}

	return client, fhub, nil
}

func Server(pluginName string, model model.FHub, functions reflect.Value) {
	var pluginMap = map[string]plugin.Plugin{
		pluginName: &FHubPlugin{
			Impl: &FHubExec{
				Model:     model,
				Functions: functions,
			},
		},
	}

	// logger := hclog.New(&hclog.LoggerOptions{
	// 	Level:      hclog.Trace,
	// 	Output:     os.Stderr,
	// 	JSONFormat: true,
	// })

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
