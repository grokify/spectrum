package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grokify/gocharts/data/histogram"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/simplego/type/maputil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	SpecFileOAS3 string `short:"s" long:"specfile" description:"Input OAS Spec File" required:"true"`
	XlsxWrite    string `short:"x" long:"xlsxwrite" description:"Output File" required:"false"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	spec := &openapi3.Spec{}

	if urlutil.IsHttp(opts.SpecFileOAS3, true, true) {
		spec, err = openapi3.ReadURL(opts.SpecFileOAS3)
	} else {
		spec, err = openapi3.ReadAndValidateFile(opts.SpecFileOAS3)
	}

	if err != nil {
		log.Fatal(err)
	}

	sm := openapi3.SpecMore{Spec: spec}

	log.Printf(
		"S_SPEC_VALID File [%s] Title [%s] Op Count [%d]",
		opts.SpecFileOAS3, spec.Info.Title, sm.OperationsCount())

	sortBy := histogram.SortValueDesc
	ops := sm.OperationCountsByTag()
	ops.WriteTableASCII(os.Stdout,
		[]string{"Tag", "Operation Count"}, sortBy, true)

	ops2 := ops.ItemCounts(sortBy)
	fmtutil.PrintJSON(ops2)

	ops2a := maputil.RecordSet(ops2)

	md := ops2a.Markdown("1. Count: ", ", Category: ", true, true)
	fmt.Println(md)
	opts.XlsxWrite = strings.TrimSpace(opts.XlsxWrite)
	if len(opts.XlsxWrite) > 0 {
		err := sm.WriteFileXLSX(opts.XlsxWrite, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("DONE")
}
