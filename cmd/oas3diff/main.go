// oas3diff compares two OpenAPI 3 specs (JSON or YAML) and reports the
// endpoints, operation IDs, and schema names unique to each.
//
// Usage:
//
//	oas3diff <spec1> <spec2>
//
// Endpoints are compared generically (path variables normalized), so a diff in
// operationId names does not hide the fact that the same endpoint exists in both.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/grokify/spectrum/openapi3"
)

func main() {
	if len(os.Args) < 3 {
		slog.Error("usage: oas3diff <spec1> <spec2>")
		os.Exit(1)
	}
	file1, file2 := os.Args[1], os.Args[2]

	diff, err := openapi3.ReadSpecsDiff(file1, file2, false)
	if err != nil {
		slog.Error(strconv.Quote(err.Error()))
		os.Exit(2)
	}

	fmt.Printf("spec 1: %s\nspec 2: %s\n\n", file1, file2)
	fmt.Print(diff.String())

	if diff.IsEmpty() {
		fmt.Println("\nno differences")
	}
	os.Exit(0)
}
