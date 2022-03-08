package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/severity"
	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/path/filepathutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint"
	"github.com/grokify/spectrum/openapi3lint/extensions"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	PolicyFile    string `short:"p" long:"policyfile" description:"Policy File" required:"true"`
	InputFileOAS3 string `short:"i" long:"inputspec" description:"Input OAS Spec File or Dir" required:"false"`
	Severity      string `short:"s" long:"severity" description:"Severity level"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	if len(opts.InputFileOAS3) > 0 {
		isDir, err := osutil.IsDir(opts.InputFileOAS3)
		if err != nil {
			log.Fatal(err)
		}
		if isDir {
			entries, err := osutil.ReadDirMore(opts.InputFileOAS3,
				regexp.MustCompile(`(?i)\.(json|yaml|yml)$`), false, true, false)
			if err != nil {
				log.Fatal(err)
			}
			files = osutil.DirEntries(entries).Names(opts.InputFileOAS3, true)
		} else {
			files = []string{opts.InputFileOAS3}
		}
		err = fmtutil.PrintJSON(files)
		if err != nil {
			log.Fatal(err)
		}
	}

	polCfg, err := openapi3lint.NewPolicyConfigFile(opts.PolicyFile)
	if err != nil {
		log.Fatal(err)
	}
	polCfg.AddRuleCollection(extensions.NewRuleCollectionExtensions())
	err = fmtutil.PrintJSON(polCfg)
	if err != nil {
		log.Fatal(err)
	}
	err = fmtutil.PrintJSON(polCfg.RuleNames())
	if err != nil {
		log.Fatal(err)
	}

	pol, err := polCfg.Policy()
	if err != nil {
		log.Fatal(errorsutil.Wrap(err, "polCfg.Policy()"))
	}
	fmtutil.PrintJSON(pol)
	fmtutil.PrintJSON(pol.RuleNames())

	severityLevel := severity.SeverityError
	if len(opts.Severity) > 0 {
		severityTry, err := severity.Parse(opts.Severity)
		if err != nil {
			log.Fatal(err)
		}
		severityLevel = severityTry
	}

	vsets := lintutil.NewPolicyViolationsSets()
	for _, file := range files {
		spec, err := openapi3.ReadFile(file, false)
		if err != nil {
			log.Fatal(err)
		}
		vsetsRule, err := pol.ValidateSpec(spec, filepathutil.FilepathLeaf(file), severityLevel)
		if err != nil {
			log.Fatal(err)
		}
		vsets.UpsertSets(vsetsRule)
	}

	fmtutil.PrintJSON(vsets.LocationsByRule())

	fmtutil.PrintJSON(vsets.CountsByRule())

	fmt.Println("DONE")
}

/*
func getPolicyConfig() openapi3lint.PolicyConfig {
	return openapi3lint.PolicyConfig{
		Rules: map[string]openapi3lint.RuleConfig{
			openapi3lint.RuleOpIdStyleCamelCase: {
				Severity: severity.SeverityError},
			openapi3lint.RuleOpSummaryExist: {
				Severity: severity.SeverityError},
			openapi3lint.RuleOpSummaryCaseFirstCapitalized: {
				Severity: severity.SeverityError},
			openapi3lint.RulePathParamStyleCamelCase: {
				Severity: severity.SeverityError},
			openapi3lint.RuleSchemaObjectPropsExist: {
				Severity: severity.SeverityError},
			openapi3lint.RuleSchemaPropEnumStylePascalCase: {
				Severity: severity.SeverityError},
			openapi3lint.RuleTagCaseFirstCapitalized: {
				Severity: severity.SeverityError},
		},
	}
}
*/
