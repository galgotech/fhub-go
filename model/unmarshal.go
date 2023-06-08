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

func UnmarshalFile(path string) (FHub, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FHub{}, err
	}
	return UnmarshalBytes(data)
}

func UnmarshalHttp(url string) (FHub, error) {
	resp, err := http.Get(url)
	if err != nil {
		return FHub{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FHub{}, err
	}
	return UnmarshalBytes(body)
}

func UnmarshalString(data string) (FHub, error) {
	return UnmarshalBytes([]byte(data))
}

func UnmarshalBytes(data []byte) (FHub, error) {
	ctx := cuecontext.New()
	value := ctx.CompileBytes(data)
	fhub, err := unmarshalStart(value)
	if err != nil {
		return FHub{}, err
	}

	err = Validator(fhub)
	if err != nil {
		return FHub{}, err
	}

	return fhub, nil
}

func unmarshalStart(value cue.Value) (FHub, error) {
	fhub := FHub{}

	outValueOf := reflect.ValueOf(&fhub).Elem()
	err := unmarshalDiscoverType("fhub", outValueOf, value)
	if err != nil {
		return fhub, err
	}

	return fhub, nil
}

type Unmarshaler interface {
	Unmarshal(field string, value cue.Value) (err error)
}

func unmarshalDiscoverType(namespace string, outValueOf reflect.Value, value cue.Value) (err error) {
	kind := outValueOf.Type().Kind()

	switch kind {
	case reflect.String:
		err = unmarshalValueString(outValueOf, value)
	case reflect.Struct:
		err = unmarshalStructIt(namespace, outValueOf, value)
	case reflect.Map:
		err = unmarshalMapIt(namespace, outValueOf, value)
	case reflect.Slice:
		err = unmarshalListIt(namespace, outValueOf, value)
	case reflect.Ptr:
		outValueOf.Set(reflect.New(outValueOf.Type().Elem()))
		err = unmarshalDiscoverType(namespace, outValueOf.Elem(), value)
	default:
		panic(fmt.Sprintf("type %q not implemented from key", kind))
	}

	return err
}

func unmarshalStructIt(namespace string, outValueOf reflect.Value, value cue.Value) (err error) {
	outTypeOf := outValueOf.Type()

	for i := 0; i < outValueOf.NumField(); i++ {
		fieldValueOf := outValueOf.Field(i)
		fieldTypeOf := outTypeOf.Field(i)

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
				namespace = fmt.Sprintf("%s.%s", namespace, name)
				err = unmarshalDiscoverType(namespace, fieldValueOf, currentValue)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func unmarshalMapIt(namespace string, outValueOf reflect.Value, value cue.Value) (err error) {
	fields, err := value.Fields()
	if err != nil {
		return err
	}

	valueValueOf := reflect.MakeMap(outValueOf.Type())
	valueValueTypeOf := valueValueOf.Type().Elem()
	for fields.Next() {
		key := fields.Label()
		val := fields.Value()

		newValueOf := reflect.New(valueValueTypeOf).Elem()
		namespace = fmt.Sprintf("%s.%s", namespace, key)
		err := unmarshalDiscoverType(namespace, newValueOf, val)
		if err != nil {
			return err
		}

		valueValueOf.SetMapIndex(reflect.ValueOf(key), newValueOf)
	}

	outValueOf.Set(valueValueOf)
	return nil
}

func unmarshalListIt(namespace string, outValueOf reflect.Value, value cue.Value) (err error) {
	fields, err := value.List()
	if err != nil {
		return err
	}

	valueValueOf := reflect.MakeSlice(outValueOf.Type(), 0, 0)
	valueValueTypeOf := valueValueOf.Type().Elem()
	for fields.Next() {
		// key := fields.Label()
		val := fields.Value()

		newValueOf := reflect.New(valueValueTypeOf).Elem()
		err := unmarshalDiscoverType(namespace, newValueOf, val)
		if err != nil {
			return err
		}

		valueValueOf = reflect.Append(valueValueOf, newValueOf)
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
