package main

import (
	"fmt"
	"log"

	"github.com/grokify/spectrum/openapi2/openapi2postman2"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/postman2"
	"github.com/jessevdk/go-flags"
)

// Convert yaml2json: https://github.com/bronze1man/yaml2json ... yaml2json_darwin_amd64

type Options struct {
	PostmanBase string `short:"b" long:"base" description:"Basic Postman File"`
	Postman     string `short:"p" long:"postman" description:"Output Postman File" required:"true"`
	Swagger     string `short:"s" long:"swagger" description:"Input Swagger File" required:"true"`
}

func main() {

	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	if 1 == 1 {
		spec, err := openapi3.ReadFile(opts.Swagger, true)
		if err != nil {
			log.Fatal(err)
		}
		sm := openapi3.SpecMore{Spec: spec}
		if 1 == 1 {
			err := sm.SchemaPropertiesWithoutDescriptionsWriteFile("rc-platform.yml.schema-properties_missing-descriptions.txt")
			if err != nil {
				log.Fatal(err)
			}
		}
		if 1 == 1 {
			err := sm.OperationParametersWithoutDescriptionsWriteFile("rc-platform.yml.op-params_missing-descriptions.txt")
			if err != nil {
				log.Fatal(err)
			}

		}
		panic("ZZZ")
	}

	cfg := openapi2postman2.Configuration{
		PostmanURLBase: "{{RINGCENTRAL_SERVER_URL}}",
		PostmanHeaders: []postman2.Header{{
			Key:   "Authorization",
			Value: "Bearer {{my_access_token}}"}}}

	conv := openapi2postman2.NewConverter(cfg)
	err = conv.MergeConvert(opts.Swagger, opts.PostmanBase, opts.Postman)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Wrote %v\n", opts.Postman)
	}

	fmt.Println("DONE")
}
