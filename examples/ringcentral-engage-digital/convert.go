package main

import (
	"fmt"
	"log"

	"github.com/grokify/spectrum/openapi3/openapi3postman2"
	"github.com/grokify/spectrum/postman2"
	"github.com/jessevdk/go-flags"
)

// Convert yaml2json: https://github.com/bronze1man/yaml2json ... yaml2json_darwin_amd64
// var ApiUrlFormat string = "https://{account}.api.engagement.dimelo.com/1.0"

type Options struct {
	PostmanBase string `short:"b" long:"base" description:"Basic Postman File"`
	Postman     string `short:"p" long:"postman" description:"Output Postman File" required:"true"`
	OpenAPISpec string `short:"s" long:"openapispec" description:"Input OpenAPI Spec File" required:"true"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	cfg := openapi3postman2.Configuration{
		UseXTagGroups:            true,
		PostmanServerURLBasePath: "1.0",
		PostmanServerURL:         "{{RINGCENTRAL_ENGAGE_SERVER_URL}}",
		PostmanHeaders: []postman2.Header{{
			Key:   "Authorization",
			Value: "Bearer {{RINGCENTRAL_ENGAGE_ACCESS_TOKEN}}"}}}

	conv := openapi3postman2.NewConverter(cfg)

	merge := true

	if merge && len(opts.PostmanBase) > 0 {
		err = conv.MergeConvert(opts.OpenAPISpec, opts.PostmanBase, opts.Postman)
	} else {
		err = conv.ConvertFile(opts.OpenAPISpec, opts.Postman)
	}

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Wrote %v\n", opts.Postman)
	}

	fmt.Println("DONE")
}
