package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/grokify/swaggman/openapi3"
	"github.com/grokify/swaggman/openapi3/openapi3postman2"
	"github.com/jessevdk/go-flags"
)

// Convert yaml2json: https://github.com/bronze1man/yaml2json ... yaml2json_darwin_amd64

type Options struct {
	Config      string `short:"c" long:"config" description:"Swaggman Config File"`
	PostmanBase string `short:"b" long:"base" description:"Basic Postman File"`
	Postman     string `short:"p" long:"postman" description:"Output Postman File" required:"true"`
	OpenAPIFile string `short:"" long:"openapi" description:"Input Swagger File" required:"true"`
}

func (opts *Options) TrimSpace() {
	opts.Config = strings.TrimSpace(opts.Config)
	opts.PostmanBase = strings.TrimSpace(opts.PostmanBase)
	opts.Postman = strings.TrimSpace(opts.Postman)
	opts.OpenAPIFile = strings.TrimSpace(opts.OpenAPIFile)
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	cfg3 := openapi3postman2.Configuration{}

	if len(opts.Config) > 0 {
		cfg3, err = openapi3postman2.ConfigurationReadFile(opts.Config)
		if err != nil {
			log.Fatal(err)
		}
	}

	spec, err := openapi3.ReadFile(opts.OpenAPIFile, true)
	if err != nil {
		log.Fatal(err)
	}
	conv := openapi3postman2.Converter{
		Configuration: cfg3,
		OpenAPISpec:   spec}

	err = conv.MergeConvert(
		opts.OpenAPIFile,
		opts.PostmanBase,
		opts.Postman)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("wrote Postman collection [%s]", opts.Postman)
	fmt.Println("DONE")
}
