package plugin

import (
	"encoding/gob"
	"errors"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

func init() {
	gob.Register([]any{})
	gob.Register(map[string]any{})
}

type FHub interface {
	Exec(function string, input map[string]any) (map[string]any, error)
}

type FHubRPC struct {
	client *rpc.Client
}

func (g *FHubRPC) Exec(function string, input map[string]any) (map[string]any, error) {
	var args any = []any{
		function,
		input,
	}
	output := map[string]any{}
	err := g.client.Call("Plugin.Exec", &args, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

type FHubRPCServer struct {
	Impl FHub
}

func (s *FHubRPCServer) Exec(args any, output *map[string]any) (err error) {
	argsList := args.([]any)
	function, ok := argsList[0].(string)
	if !ok {
		return errors.New("invalid function interface")
	}
	input, ok := argsList[1].(map[string]any)
	if !ok {
		return errors.New("invalid args interface")
	}

	*output, err = s.Impl.Exec(function, input)
	return err
}

type FHubPlugin struct {
	Impl FHub
}

func (p *FHubPlugin) Server(*plugin.MuxBroker) (any, error) {
	return &FHubRPCServer{Impl: p.Impl}, nil
}

func (FHubPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &FHubRPC{client: c}, nil
}
