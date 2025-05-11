package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/spectrum/openapi3"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	SpecFileOAS3          string `short:"s" long:"specfile" description:"Input OAS Spec File" required:"true"`
	WriteOpStatusCodeXlsx string `long:"writeopstatus" description:"Output File" required:"false"`
	XlsxWrite             string `short:"x" long:"xlsxwrite" description:"Output File" required:"false"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var spec *openapi3.Spec

	if urlutil.IsHTTP(opts.SpecFileOAS3, true, true) {
		spec, err = openapi3.ReadURL(opts.SpecFileOAS3)
	} else {
		spec, err = openapi3.ReadFile(opts.SpecFileOAS3, true)
	}
	if err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}

	sm := openapi3.SpecMore{Spec: spec}

	slog.Info(
		"validSpecInfo",
		"filename", opts.SpecFileOAS3,
		"title", spec.Info.Title,
		"opsCount", sm.OperationsCount())

	sortBy := histogram.SortValueDesc
	ops := sm.OperationCountsByTag()
	err = ops.WriteTableASCII(os.Stdout,
		[]string{"Tag", "Operation Count"}, sortBy, true)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(3)
	}

	ops2 := ops.ItemCounts(sortBy)
	err = fmtutil.PrintJSON(ops2)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(3)
	}

	md := ops2.Markdown("1. Count: ", ", Category: ", true, true)
	slog.Info(md)
	opts.XlsxWrite = strings.TrimSpace(opts.XlsxWrite)
	if len(opts.XlsxWrite) > 0 {
		err := sm.WriteFileXLSX(opts.XlsxWrite, nil, nil, nil)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(4)
		}
	}

	opts.WriteOpStatusCodeXlsx = strings.TrimSpace(opts.WriteOpStatusCodeXlsx)
	if len(opts.WriteOpStatusCodeXlsx) > 0 {
		err := sm.WriteFileXLSXOperationStatusCodes(
			opts.WriteOpStatusCodeXlsx)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(5)
		}
	}

	opIDs := sm.OperationIDs()
	slog.Info("operationIDs", "count", len(opIDs))

	endpoints := sm.PathMethods(true)
	slog.Info("endpoints", "count", len(endpoints))

	slog.Info("DONE")
	os.Exit(0)
}
