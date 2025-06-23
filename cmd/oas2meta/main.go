package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/spectrum/openapi2"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("usage: oas2meta ...<filepath>")
		os.Exit(1)
	} else if set, err := openapi2.NewSpecMetaSetFilepaths(os.Args[1:]); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	} else {
		fmtutil.MustPrintJSON(set)
		fmt.Println("DONE")
		os.Exit(0)
	}
}
