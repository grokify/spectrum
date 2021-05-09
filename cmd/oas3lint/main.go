package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/io/ioutilmore"
	"github.com/grokify/simplego/log/severity"
	"github.com/grokify/simplego/path/filepathutil"
	"github.com/grokify/swaggman/openapi3"
	"github.com/grokify/swaggman/openapi3lint"
	"github.com/grokify/swaggman/openapi3lint/lintutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	SpecFileOAS3 string `short:"s" long:"specfile" description:"Input OAS Spec File" required:"true"`
	PolicyFile   string `short:"p" long:"policyfile" description:"Policy File" required:"true"`
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

	fmtutil.PrintJSON(files)

	polCfg, err := openapi3lint.NewPolicyConfigFile(opts.PolicyFile)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(polCfg)
	pol, err := polCfg.StandardPolicy()
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(pol)

	vsets := lintutil.NewPolicyViolationsSets()
	for _, file := range files {

		spec, err := openapi3.ReadFile(file, false)
		if err != nil {
			log.Fatal(err)
		}
		vsetsRule, err := pol.ValidateSpec(spec, filepathutil.FilepathLeaf(file), severity.SeverityWarning)
		if err != nil {
			log.Fatal(err)
		}
		vsets.UpsertSets(vsetsRule)
	}

	fmtutil.PrintJSON(vsets.LocationsByRule())
	/*
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
		fmt.Printf("violations [%d]\n", vsetsByFile.Count())
	*/
	fmt.Println("DONE")
}

/*
func FilepathLeaf(s string) string {
	_, file := filepath.Split(s)
	return file
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
