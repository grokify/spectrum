package main

import (
	"fmt"
	"log"

	"github.com/andrewcretin/swaggman/openapi3conv"
	"github.com/andrewcretin/swaggman/postman2"
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

	cfg := openapi3conv.Configuration{
		PostmanServerURLBasePath: "1.0",
		PostmanServerURL:         "{{RINGCENTRAL_ENGAGE_SERVER_URL}}",
		PostmanHeaders: []postman2.Header{{
			Key:   "Authorization",
			Value: "Bearer {{RINGCENTRAL_ENGAGE_ACCESS_TOKEN}}"}}}

	conv := openapi3conv.NewConverter(cfg)

	if len(opts.PostmanBase) > 0 {
		err = conv.MergeConvert(opts.Swagger, opts.PostmanBase, opts.Postman)
	} else {
		err = conv.Convert(opts.Swagger, opts.Postman)
	}

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Wrote %v\n", opts.Postman)
	}

	fmt.Println("DONE")
}
