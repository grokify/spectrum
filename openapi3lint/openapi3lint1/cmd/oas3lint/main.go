package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/io/ioutilmore"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/openapi3lint1"
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

	files := []string{}
	isDir, err := ioutilmore.IsDir(opts.SpecFileOAS3)
	if err != nil {
		log.Fatal(err)
	}
	if isDir {
		_, files, err = ioutilmore.ReadDirMore(opts.SpecFileOAS3,
			regexp.MustCompile(`\.(yaml|yml|json)$`), true, true)
	} else {
		files = []string{opts.SpecFileOAS3}
	}

	// fmtutil.PrintJSON(files)

	policies := openapi3lint1.NewPolicySimple([]string{
		openapi3lint1.RuleDatatypeIntFormatIsInt32OrInt64})

	vsetsByFile := openapi3lint1.NewPolicyViolationsSetsByFile()

	for _, file := range files {
		spec, err := openapi3.ReadFile(file, false)
		if err != nil {
			log.Fatal(err)
		}

		vsetSpec, err := openapi3lint1.SpecCheckViolations(spec, policies)
		if err != nil {
			log.Fatal(err)
		}
		vsetsByFile.Sets[file] = vsetSpec

	}
	fmtutil.PrintJSON(vsetsByFile.LocationsByRule(true, true))
	fmt.Printf("violations [%d]\n", vsetsByFile.Count())
	fmt.Println("DONE")
}

func getRulesSimple() []string {
	return []string{
		openapi3lint1.RuleOpIdStyleCamelCase,
		openapi3lint1.RuleOpSummaryExist,
		openapi3lint1.RuleOpSummaryCaseFirstCapitalized,
		openapi3lint1.RulePathParamStyleCamelCase,
		openapi3lint1.RuleSchemaObjectPropsExist,
		openapi3lint1.RuleSchemaPropEnumStylePascalCase,
		openapi3lint1.RuleTagCaseFirstCapitalized,
	}
}
