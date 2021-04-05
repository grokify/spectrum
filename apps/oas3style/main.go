package main

import (
	"fmt"
	"log"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/swaggman/openapi3"
	"github.com/grokify/swaggman/openapi3/stylechecker"
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
		stylechecker.RuleOpIdStyleCamelCase,
		stylechecker.RulePathParamStyleCamelCase,
		stylechecker.RuleSchemaObjectPropsExist,
		stylechecker.RuleSchemaPropEnumStylePascalCase,
		stylechecker.RuleTagCaseFirstAlphaUpper,
	}
	ruleset := stylechecker.NewRuleSet(rules)

	vios, err := stylechecker.SpecCheckViolations(spec, ruleset)
	if err != nil {
		log.Fatal(err)
	}

	fmtutil.PrintJSON(vios.LocationsByRule())

	fmt.Println("DONE")
}
