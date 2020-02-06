package main

import (
	"fmt"
	"log"

	"github.com/grokify/swaggman/openapi3"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	SpecFileOAS3 string `short:"s" long:"specfile" description:"Input OAS Spec File" required:"true"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	valid, err := openapi3.ValidateSpec(opts.SpecFileOAS3)

	if err != nil {
		log.Fatal(err)
	} else if !valid {
		log.Fatalf("E_SPEC_NOT_VALID [%v]", opts.SpecFileOAS3)
	} else {
		log.Fatalf("S_SPEC_VALID [%v]", opts.SpecFileOAS3)
	}

	fmt.Println("DONE")
}
