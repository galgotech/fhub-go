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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

func UnmarshalFile(path string) (Fhub, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Fhub{}, err
	}
	return UnmarshalBytes(data)
}

func UnmarshalHttp(url string) (Fhub, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Fhub{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Fhub{}, err
	}
	return UnmarshalBytes(body)
}

func UnmarshalBytes(data []byte) (Fhub, error) {
	ctx := cuecontext.New()
	value := ctx.CompileBytes(data)
	fhub, err := unmarshalStart(value)
	if err != nil {
		return Fhub{}, err
	}
	return fhub, nil
}

func UnmarshalString(data string) (Fhub, error) {
	ctx := cuecontext.New()
	value := ctx.CompileString(data)
	fhub, err := unmarshalStart(value)
	if err != nil {
		return Fhub{}, err
	}
	return fhub, nil
}

func unmarshalStart(value cue.Value) (Fhub, error) {
	fhub := Fhub{}

	outValueOf := reflect.ValueOf(&fhub).Elem()
	err := unmarshalIt(outValueOf, "fhub", value)
	if err != nil {
		return fhub, err
	}

	return fhub, nil
}

type Unmarshaler interface {
	Unmarshal(field string, value cue.Value) (err error)
}

func unmarshalIt(outValueOf reflect.Value, base string, value cue.Value) (err error) {
	fmt.Println(outValueOf.Type().Name(), base)
	outTypeOf := outValueOf.Type()

	for i := 0; i < outValueOf.NumField(); i++ {
		fieldValueOf := outValueOf.Field(i)
		fieldTypeOf := outTypeOf.Field(i)

		kind := fieldTypeOf.Type.Kind()
		tag := fieldTypeOf.Tag.Get("fhub")
		unmarshal := fieldTypeOf.Tag.Get("fhub-unmarshal") == "true"

		var name string
		if tag == "" {
			name = fieldTypeOf.Name
			name = strings.ToLower(name[:1]) + name[1:]
		} else {
			name = tag
		}

		currentValue := value.LookupPath(cue.ParsePath(name))
		if currentValue.Exists() {
			if unmarshal {
				f, ok := outValueOf.Addr().Interface().(Unmarshaler)
				if !ok {
					return errors.New("type does not implement Unmarshaler")
				}
				err = f.Unmarshal(name, currentValue)
			} else {
				switch kind {
				case reflect.String:
					err = unmarshalValueString(fieldValueOf, currentValue)
				case reflect.Struct:
					err = unmarshalIt(fieldValueOf, name, currentValue)
				case reflect.Map, reflect.Slice:
					err = unmarshalValueIt(fieldValueOf, currentValue)
				default:
					panic(fmt.Sprintf("type %q not implemented from key %q", kind, name))
				}
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func unmarshalValueIt(outValueOf reflect.Value, value cue.Value) (err error) {
	fields, err := value.Fields()
	if err != nil {
		return err
	}

	outKindOf := outValueOf.Kind()
	var valueValueOf reflect.Value
	switch outKindOf {
	case reflect.Map:
		valueValueOf = reflect.MakeMap(outValueOf.Type())
	case reflect.Slice:
		valueValueOf = reflect.MakeSlice(outValueOf.Type(), 0, 0)
	default:
		panic("invalid kind")
	}

	valueValueTypeOf := valueValueOf.Type().Elem()
	for fields.Next() {
		key := fields.Label()
		val := fields.Value()

		newValueOf := reflect.New(valueValueTypeOf).Elem()
		err := unmarshalIt(newValueOf, key, val)
		if err != nil {
			return err
		}

		switch outKindOf {
		case reflect.Map:
			valueValueOf.SetMapIndex(reflect.ValueOf(key), newValueOf)
		case reflect.Slice:
			valueValueOf = reflect.Append(valueValueOf, newValueOf)
		}
	}

	outValueOf.Set(valueValueOf)
	return nil
}

func unmarshalValueString(outValueOf reflect.Value, value cue.Value) error {
	val, err := value.String()
	if err != nil {
		return err
	}
	outValueOf.SetString(val)
	return nil
}
