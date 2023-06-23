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
	var err error
	fhubPath, err = os.Getwd()
	if err != nil {
		panic(err)
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
