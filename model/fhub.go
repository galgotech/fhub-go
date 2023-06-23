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

type FHub struct {
	// Function namespace
	Name string `validate:"required"`
	// Function version
	Version string `validate:"required"`
	// FHub schema version
	SpecVersion string `validate:"required"`
	Build       Build
	Serving     Serving
	Functions   map[string]Function `validate:"min=1,dive"`
}

func (in *FHub) DeepCopy() (out *FHub) {
	out = new(FHub)
	*out = *in

	out.Functions = make(map[string]Function, len(in.Functions))
	for key, function := range in.Functions {
		out.Functions[key] = function
	}

	out.Build = *in.Build.DeepCopy()
	out.Serving = *in.Serving.DeepCopy()

	return
}

type Serving struct {
	Http *Http
	Grpc *Grpc
}

func (in *Serving) DeepCopy() (out *Serving) {
	out = new(Serving)
	*out = *in

	if out.Http != nil {
		in, out := &in.Http, &out.Http
		*out = new(Http)
		**out = **in
	}
	if out.Grpc != nil {
		in, out := &in.Grpc, &out.Grpc
		*out = new(Grpc)
		**out = **in
	}
	return
}

type Http struct {
	Url string `validate:"required"`
}

func (in *Http) DeepCopy() (out *Http) {
	out = new(Http)
	*out = *in
	return
}

type Grpc struct {
	Proto string `validate:"required"`
}

func (in *Grpc) DeepCopy() (out *Grpc) {
	out = new(Grpc)
	*out = *in
	return
}

type Build struct {
	Local     *Local     `validate:"required_without=Container"`
	Container *Container `validate:"required_without=Local"`
}

func (in *Build) DeepCopy() (out *Build) {
	out = new(Build)
	*out = *in
	if out.Local != nil {
		in, out := &in.Local, &out.Local
		*out = new(Local)
		**out = **in
	}
	if out.Container != nil {
		in, out := &in.Container, &out.Container
		*out = new(Container)
		**out = **in
	}
	return out
}

type Local struct {
	Source string `validate:"required"`
}

func (in *Local) DeepCopy() (out *Local) {
	out = new(Local)
	*out = *in
	return
}

type Container struct {
	Image         string `validate:"required_without=ContainerFile"`
	ContainerFile string `validate:"required_without=Image"`
	Source        string `validate:"required"`
}

func (in *Container) DeepCopy() (out *Container) {
	out = new(Container)
	*out = *in
	return
}
