package main

import (
	"fmt"
	"log"

	"github.com/grokify/simplego/fmt/fmtutil"
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

	spec, err := openapi3.ReadFile(opts.SpecFileOAS3, true)
	if err != nil {
		log.Fatal(err)
	}
	sm := openapi3.SpecMore{Spec: spec}
	fmtutil.PrintJSON(sm.Stats())

	rules := []string{
		openapi3lint.RuleOpIdStyleCamelCase,
		openapi3lint.RuleOpSummaryExist,
		openapi3lint.RuleOpSummaryCaseFirstCapitalized,
		openapi3lint.RulePathParamStyleCamelCase,
		openapi3lint.RuleSchemaObjectPropsExist,
		openapi3lint.RuleSchemaPropEnumStylePascalCase,
		openapi3lint.RuleTagCaseFirstCapitalized,
	}
	ruleset := openapi3lint.NewPolicySimple(rules)

	vios, err := openapi3lint.SpecCheckViolations(spec, ruleset)
	if err != nil {
		log.Fatal(err)
	}

	fmtutil.PrintJSON(vios.LocationsByRule())

	fmtutil.PrintJSON(sm.Stats())

	fmt.Println("DONE")
}
