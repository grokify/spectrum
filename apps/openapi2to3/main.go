package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/grokify/gotilla/io/ioutilmore"
	"github.com/grokify/swaggman/swagger2"
	"github.com/jessevdk/go-flags"
)

// install: go get github.com/grokify/swaggman/apps/openapi2to3

type Options struct {
	OAS2File string `short:"i" long:"input" description:"Input filepath" required:"true"`
	OAS3File string `short:"o" long:"output" description:"Output filepath" required:"true"`
	Pretty   []bool `short:"p" long:"pretty" description:"Pretty print output"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	opts.OAS2File = strings.TrimSpace(opts.OAS2File)
	opts.OAS3File = strings.TrimSpace(opts.OAS3File)
	isFile, err := ioutilmore.IsFile(opts.OAS2File)
	if err != nil {
		log.Fatal(err)
	} else if !isFile {
		log.Fatalf("E_INPUT_FILE_IS_NOT_FILE [%v]", opts.OAS2File)
	}

	wantPretty := false
	if len(opts.Pretty) > 0 {
		wantPretty = true
	}

	err = swagger2.ConvertOAS2FileToOAS3File(opts.OAS2File, opts.OAS3File, 0644, wantPretty)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", opts.OAS3File)

	fmt.Println("DONE")
}
