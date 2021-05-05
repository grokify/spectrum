package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/io/ioutilmore"
	"github.com/grokify/swaggman/openapi3"
	"github.com/grokify/swaggman/openapi3/openapi3lint"
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

	policies := openapi3lint.NewPolicySimple([]string{
		openapi3lint.RuleDatatypeIntFormatIsInt32OrInt64})

	vsetsByFile := openapi3lint.NewPolicyViolationsSetsByFile()

	for _, file := range files {
		spec, err := openapi3.ReadFile(file, false)
		if err != nil {
			log.Fatal(err)
		}

		vsetSpec, err := openapi3lint.SpecCheckViolations(spec, policies)
		if err != nil {
			log.Fatal(err)
		}
		vsetsByFile.Sets[file] = vsetSpec

	}
	fmtutil.PrintJSON(vsetsByFile.LocationsByRule(true, true))

	fmt.Println("DONE")
}

func getRulesSimple() []string {
	return []string{
		openapi3lint.RuleOpIdStyleCamelCase,
		openapi3lint.RuleOpSummaryExist,
		openapi3lint.RuleOpSummaryCaseFirstCapitalized,
		openapi3lint.RulePathParamStyleCamelCase,
		openapi3lint.RuleSchemaObjectPropsExist,
		openapi3lint.RuleSchemaPropEnumStylePascalCase,
		openapi3lint.RuleTagCaseFirstCapitalized,
	}
}
