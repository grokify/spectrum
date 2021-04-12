package main

import (
	"fmt"
	"log"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/swaggman/openapi3"
	"github.com/grokify/swaggman/openapi3/styleguide"
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
		styleguide.RuleOpIdStyleCamelCase,
		styleguide.RuleOpSummaryNotEmpty,
		styleguide.RuleOpSummaryCaseFirstCapitalized,
		styleguide.RulePathParamStyleCamelCase,
		styleguide.RuleSchemaObjectPropsExist,
		styleguide.RuleSchemaPropEnumStylePascalCase,
		styleguide.RuleTagCaseFirstCapitalized,
	}
	ruleset := styleguide.NewPolicySimple(rules)

	vios, err := styleguide.SpecCheckViolations(spec, ruleset)
	if err != nil {
		log.Fatal(err)
	}

	fmtutil.PrintJSON(vios.LocationsByRule())

	fmtutil.PrintJSON(sm.Stats())

	fmt.Println("DONE")
}
