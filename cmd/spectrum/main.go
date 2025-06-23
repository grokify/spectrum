package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3/openapi3postman2"
	flags "github.com/jessevdk/go-flags"
)

// Convert yaml2json: https://github.com/bronze1man/yaml2json ... yaml2json_darwin_amd64

type Options struct {
	OpenAPIFile string `short:"O" long:"openapiFile" description:"Input Swagger File" required:"true"`
	Config      string `short:"C" long:"config" description:"Spectrum Config File"`
	PostmanBase string `short:"B" long:"basePostmanFile" description:"Basic Postman File"`
	Postman     string `short:"P" long:"postmanFile" description:"Output Postman File"`
	XLSXFile    string `short:"X" long:"xlsxFile" description:"Output XLSX File"`
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

	spec, err := openapi3.ReadFile(opts.OpenAPIFile, true)
	if err != nil {
		log.Fatal(err)
	}

	if len(opts.Postman) > 0 {
		cfg3 := openapi3postman2.Configuration{}

		if len(opts.Config) > 0 {
			cfg3, err = openapi3postman2.ConfigurationReadFile(opts.Config)
			if err != nil {
				log.Fatal(errorsutil.Wrap(err, "openapi3postman2.ConfigurationReadFile"))
			}
		}

		conv := openapi3postman2.Converter{
			Configuration: cfg3,
			OpenAPISpec:   spec}

		err = conv.MergeConvert(
			opts.OpenAPIFile,
			opts.PostmanBase,
			opts.Postman)

		if err != nil {
			log.Fatal(errorsutil.Wrap(err, "spectrum.main << conv.MergeConvert"))
		}

		fmt.Printf("wrote Postman collection [%s]\n", opts.Postman)
	}
	if len(opts.XLSXFile) > 0 {
		sm := openapi3.SpecMore{Spec: spec}
		err := sm.WriteFileXLSX(opts.XLSXFile, nil, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("DONE")
}
