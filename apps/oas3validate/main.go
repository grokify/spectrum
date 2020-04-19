package main

import (
	"fmt"

	"github.com/grokify/swaggman/openapi3"
	"github.com/jessevdk/go-flags"

	log "github.com/sirupsen/logrus"
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

	spec, err := openapi3.ReadAndValidateFile(opts.SpecFileOAS3)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Infof("S_SPEC_VALID File [%s] Title [%s]", opts.SpecFileOAS3, spec.Info.Title)
	}

	fmt.Println("DONE")
}
