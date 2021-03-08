package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grokify/gocharts/data/frequency"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/type/maputil"
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

	spec, err := openapi3.ReadAndValidateFile(opts.SpecFileOAS3)

	if err != nil {
		log.Fatal(err)
	}

	sm := openapi3.SpecMore{Spec: spec}

	log.Printf(
		"S_SPEC_VALID File [%s] Title [%s] Op Count [%d]",
		opts.SpecFileOAS3, spec.Info.Title, sm.OperationsCount())

	sortBy := frequency.SortValueDesc
	ops := sm.OperationCountsByTag()
	ops.WriteTableASCII(os.Stdout,
		[]string{"Tag", "Operation Count"}, sortBy, true)

	ops2 := ops.ItemCounts(sortBy)
	fmtutil.PrintJSON(ops2)

	ops2a := maputil.RecordSet(ops2)

	md := ops2a.Markdown("1. Count: ", ", Category: ", true, true)
	fmt.Println(md)

	fmt.Println("DONE")
}
