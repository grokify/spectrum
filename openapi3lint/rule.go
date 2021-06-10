package openapi3lint

import (
	"fmt"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/log/severity"
	"github.com/grokify/simplego/text/stringcase"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
	"github.com/grokify/spectrum/openapi3lint/ruleintstdformat"
	"github.com/grokify/spectrum/openapi3lint/ruleopidstyle"
	"github.com/grokify/spectrum/openapi3lint/ruleopsummaryexist"
	"github.com/grokify/spectrum/openapi3lint/ruleopsummarystylefirstuppercase"
	"github.com/grokify/spectrum/openapi3lint/rulepathparamstyle"
	"github.com/grokify/spectrum/openapi3lint/ruleschemaobjectpropsexist"
	"github.com/grokify/spectrum/openapi3lint/ruleschemapropenumstyle"
	"github.com/grokify/spectrum/openapi3lint/ruleschemareferences"
	"github.com/grokify/spectrum/openapi3lint/ruletagstylefirstuppercase"
)

type Rule interface {
	Name() string
	Scope() string
	Severity() string
	ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation
	ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation
}

type RulesMap map[string]Rule

func ValidateRules(rules map[string]Rule) error {
	unknownSeverities := []string{}
	for ruleName, rule := range rules {
		_, err := severity.Parse(rule.Severity())
		if err != nil {
			unknownSeverities = append(unknownSeverities, ruleName)
		}
	}
	if len(unknownSeverities) > 0 {
		unknownSeverities = stringsutil.Dedupe(unknownSeverities)
		sort.Strings(unknownSeverities)
		return fmt.Errorf(
			"rules with unknown severities [%s] valid [%s]",
			strings.Join(unknownSeverities, ","),
			strings.Join(severity.Severities(), ","))
	}
	return nil
}

type EmptyRule struct{}

func (rule EmptyRule) Name() string     { return "" }
func (rule EmptyRule) Scope() string    { return "" }
func (rule EmptyRule) Severity() string { return "" }
func (rule EmptyRule) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}
func (rule EmptyRule) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}

type StandardRuleNames struct {
	ruleNames map[string]int
}

func NewStandardRuleNames() *StandardRuleNames {
	srn := &StandardRuleNames{
		ruleNames: map[string]int{}}
	names := standardRuleNamesList()
	for _, name := range names {
		srn.ruleNames[name] = 1
	}
	return srn
}

func (srn *StandardRuleNames) Slice() []string {
	return standardRuleNamesList()
}

func (srn *StandardRuleNames) Exists(ruleName string) bool {
	if _, ok := srn.ruleNames[ruleName]; ok {
		return true
	}
	return false
}

func standardRuleNamesList() []string {
	rulenames := []string{
		lintutil.RulenameDatatypeIntFormatStandardExist,
		lintutil.RulenameOpIdStyleCamelCase,
		lintutil.RulenameOpIdStyleKebabCase,
		lintutil.RulenameOpIdStylePascalCase,
		lintutil.RulenameOpIdStyleSnakeCase,
		lintutil.RulenameOpSummaryExist,
		lintutil.RulenameOpSummaryStyleFirstUpperCase,
		lintutil.RulenamePathParamStyleCamelCase,
		lintutil.RulenamePathParamStyleKebabCase,
		lintutil.RulenamePathParamStylePascalCase,
		lintutil.RulenamePathParamStyleSnakeCase,
		lintutil.RulenameSchemaHasReference,
		lintutil.RulenameSchemaReferenceHasSchema,
		lintutil.RulenameSchemaObjectPropsExist,
		lintutil.RulenameSchemaPropEnumStyleCamelCase,
		lintutil.RulenameSchemaPropEnumStyleKebabCase,
		lintutil.RulenameSchemaPropEnumStylePascalCase,
		lintutil.RulenameSchemaPropEnumStyleSnakeCase,
		lintutil.RulenameTagStyleFirstUpperCase,
	}
	sort.Strings(rulenames)
	return rulenames
}

func NewStandardRule(name, sev string) (Rule, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case lintutil.RulenameDatatypeIntFormatStandardExist:
		return ruleintstdformat.NewRule(sev), nil

	case lintutil.RulenameOpIdStyleCamelCase:
		return ruleopidstyle.NewRule(sev, stringcase.CamelCase)
	case lintutil.RulenameOpIdStyleKebabCase:
		return ruleopidstyle.NewRule(sev, stringcase.KebabCase)
	case lintutil.RulenameOpIdStylePascalCase:
		return ruleopidstyle.NewRule(sev, stringcase.PascalCase)
	case lintutil.RulenameOpIdStyleSnakeCase:
		return ruleopidstyle.NewRule(sev, stringcase.SnakeCase)

	case lintutil.RulenameOpSummaryExist:
		return ruleopsummaryexist.NewRule(sev), nil
	case lintutil.RulenameOpSummaryStyleFirstUpperCase:
		return ruleopsummarystylefirstuppercase.NewRule(sev), nil

	case lintutil.RulenamePathParamStyleCamelCase:
		return rulepathparamstyle.NewRule(sev, stringcase.CamelCase)
	case lintutil.RulenamePathParamStyleKebabCase:
		return rulepathparamstyle.NewRule(sev, stringcase.KebabCase)
	case lintutil.RulenamePathParamStylePascalCase:
		return rulepathparamstyle.NewRule(sev, stringcase.PascalCase)
	case lintutil.RulenamePathParamStyleSnakeCase:
		return rulepathparamstyle.NewRule(sev, stringcase.SnakeCase)

	case lintutil.RulenameSchemaHasReference:
		return ruleschemareferences.NewRule(sev, lintutil.RulenameSchemaHasReference)
	case lintutil.RulenameSchemaReferenceHasSchema:
		return ruleschemareferences.NewRule(sev, lintutil.RulenameSchemaReferenceHasSchema)

	case lintutil.RulenameSchemaObjectPropsExist:
		return ruleschemaobjectpropsexist.NewRule(sev), nil

	case lintutil.RulenameSchemaPropEnumStyleCamelCase:
		return ruleschemapropenumstyle.NewRule(sev, stringcase.CamelCase)
	case lintutil.RulenameSchemaPropEnumStyleKebabCase:
		return ruleschemapropenumstyle.NewRule(sev, stringcase.KebabCase)
	case lintutil.RulenameSchemaPropEnumStylePascalCase:
		return ruleschemapropenumstyle.NewRule(sev, stringcase.PascalCase)
	case lintutil.RulenameSchemaPropEnumStyleSnakeCase:
		return ruleschemapropenumstyle.NewRule(sev, stringcase.SnakeCase)

	case lintutil.RulenameTagStyleFirstUpperCase:
		return ruletagstylefirstuppercase.NewRule(sev), nil
	}
	return EmptyRule{}, fmt.Errorf("NewStandardRule: rule [%s] not found", name)
}
