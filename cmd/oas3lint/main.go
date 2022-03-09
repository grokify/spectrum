package main

import (
	"fmt"
	"regexp"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
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
	logutil.FatalOnError(err)

	var files []string
	if len(opts.InputFileOAS3) > 0 {
		isDir, err := osutil.IsDir(opts.InputFileOAS3)
		logutil.FatalOnError(err)

		if isDir {
			entries, err := osutil.ReadDirMore(opts.InputFileOAS3,
				regexp.MustCompile(`(?i)\.(json|yaml|yml)$`), false, true, false)
			logutil.FatalOnError(err)

			files = osutil.DirEntries(entries).Names(opts.InputFileOAS3, true)
		} else {
			files = []string{opts.InputFileOAS3}
		}
		err = fmtutil.PrintJSON(files)
		logutil.FatalOnError(err)
	}

	polCfg, err := openapi3lint.NewPolicyConfigFile(opts.PolicyFile)
	logutil.FatalOnError(err)

	polCfg.AddRuleCollection(extensions.NewRuleCollectionExtensions())
	logutil.FatalOnError(fmtutil.PrintJSON(polCfg))
	logutil.FatalOnError(fmtutil.PrintJSON(polCfg.RuleNames()))

	pol, err := polCfg.Policy()
	logutil.FatalOnError(errorsutil.Wrap(err, "polCfg.Policy()"))
	logutil.FatalOnError(fmtutil.PrintJSON(pol))
	logutil.FatalOnError(fmtutil.PrintJSON(pol.RuleNames()))

	severityLevel := severity.SeverityError
	if len(opts.Severity) > 0 {
		severityTry, err := severity.Parse(opts.Severity)
		logutil.FatalOnError(err)
		severityLevel = severityTry
	}

	vsets := lintutil.NewPolicyViolationsSets()
	for _, file := range files {
		spec, err := openapi3.ReadFile(file, false)
		logutil.FatalOnError(err)

		vsetsRule, err := pol.ValidateSpec(spec, filepathutil.FilepathLeaf(file), severityLevel)
		logutil.FatalOnError(err)
		logutil.FatalOnError(vsets.UpsertSets(vsetsRule))
	}

	logutil.FatalOnError(fmtutil.PrintJSON(vsets.LocationsByRule()))
	logutil.FatalOnError(fmtutil.PrintJSON(vsets.CountsByRule()))

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
