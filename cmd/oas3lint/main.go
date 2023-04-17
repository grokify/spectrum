package main

import (
	"fmt"
	"regexp"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/spectrum/openapi3lint"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	PolicyFile    string `short:"p" long:"policyfile" description:"Policy File" required:"true"`
	InputFileOAS3 string `short:"i" long:"inputspec" description:"Input OAS Spec File or Dir" required:"false"`
	Severity      string `short:"s" long:"severity" description:"Severity level"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)
	fmtutil.MustPrintJSON(opts)

	vsets, err := ValidateSpecFiles(opts.InputFileOAS3, opts.PolicyFile, opts.Severity)
	logutil.FatalErr(err)

	fmtutil.MustPrintJSON(vsets.LocationsByRule())
	fmtutil.MustPrintJSON(vsets.CountsByRule())

	fmt.Println("DONE")
}

func ValidateSpecFiles(specFileOrDir string, policyfile, sev string) (*lintutil.PolicyViolationsSets, error) {
	files, err := filesFromFileOrDir(specFileOrDir)
	if err != nil {
		return nil, err
	}

	polCfg, err := openapi3lint.NewPolicyConfigFile(policyfile)
	if err != nil {
		return nil, err
	}
	//polCfg.AddRuleCollection(extensions.NewRuleCollectionExtensions())
	//logutil.FatalErr(fmtutil.PrintJSON(polCfg))
	//logutil.FatalErr(fmtutil.PrintJSON(polCfg.RuleNames()))

	pol, err := polCfg.Policy()
	if err != nil {
		return nil, err
	}
	fmtutil.MustPrintJSON(pol)
	fmtutil.MustPrintJSON(pol.RuleNames())

	return pol.ValidateSpecFiles(sev, files)
}

func filesFromFileOrDir(filename string) ([]string, error) {
	return osutil.Filenames(filename, regexp.MustCompile(`(?i)\.(json|yaml|yml)$`), false, false)
}

/*
func filesFromFileOrDirOld(filename string) ([]string, error) {
	var files []string
	if len(filename) > 0 {
		isDir, err := osutil.IsDir(filename)
		if err != nil {
			return files, err
		}

		if isDir {
			entries, err := osutil.ReadDirMore(filename,
				regexp.MustCompile(`(?i)\.(json|yaml|yml)$`), false, true, false)
			logutil.FatalErr(err)

			files = osutil.DirEntries(entries).Names(filename, true)
		} else {
			files = []string{filename}
		}
	} else {
		files = []string{filename}
	}
	return files, nil
}
*/

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
