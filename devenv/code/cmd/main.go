package main

import (
	"log"

	fhub "github.com/galgotech/fhub-runtime-go"
	"github.com/galgotech/fhub-runtime-go/devenv/code/pkg"
)

func main() {
	f := &pkg.Functions{}
	fhub.SetPath("devenv/code")
	err := fhub.Run(f)
	if err != nil {
		log.Fatal(err)
	}
}
