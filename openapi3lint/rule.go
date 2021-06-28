package openapi3lint

import (
	"fmt"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/log/severity"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type Rule interface {
	Name() string
	Scope() string
	ProcessSpec(spec *openapi3.Spec, pointerBase string) []lintutil.PolicyViolation
	ProcessOperation(spec *openapi3.Spec, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation
}

type PolicyRule struct {
	Rule     Rule
	Severity string
}

type RulesMap map[string]Rule

func ValidateRules(policyRules map[string]PolicyRule) error {
	unknownSeverities := []string{}
	for ruleName, policyRule := range policyRules {
		_, err := severity.Parse(policyRule.Severity)
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

/*
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

func NewStandardRule(name string) (Rule, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case lintutil.RulenameDatatypeIntFormatStandardExist:
		return ruleintstdformat.NewRule(), nil

	case lintutil.RulenameOpIdStyleCamelCase:
		return ruleopidstyle.NewRule(stringcase.CamelCase)
	case lintutil.RulenameOpIdStyleKebabCase:
		return ruleopidstyle.NewRule(stringcase.KebabCase)
	case lintutil.RulenameOpIdStylePascalCase:
		return ruleopidstyle.NewRule(stringcase.PascalCase)
	case lintutil.RulenameOpIdStyleSnakeCase:
		return ruleopidstyle.NewRule(stringcase.SnakeCase)

	case lintutil.RulenameOpSummaryExist:
		return ruleopsummaryexist.NewRule(), nil
	case lintutil.RulenameOpSummaryStyleFirstUpperCase:
		return ruleopsummarystylefirstuppercase.NewRule(), nil

	case lintutil.RulenamePathParamStyleCamelCase:
		return rulepathparamstyle.NewRule(stringcase.CamelCase)
	case lintutil.RulenamePathParamStyleKebabCase:
		return rulepathparamstyle.NewRule(stringcase.KebabCase)
	case lintutil.RulenamePathParamStylePascalCase:
		return rulepathparamstyle.NewRule(stringcase.PascalCase)
	case lintutil.RulenamePathParamStyleSnakeCase:
		return rulepathparamstyle.NewRule(stringcase.SnakeCase)

	case lintutil.RulenameSchemaHasReference:
		return ruleschemareferences.NewRule(lintutil.RulenameSchemaHasReference)
	case lintutil.RulenameSchemaReferenceHasSchema:
		return ruleschemareferences.NewRule(lintutil.RulenameSchemaReferenceHasSchema)

	case lintutil.RulenameSchemaObjectPropsExist:
		return ruleschemaobjectpropsexist.NewRule(), nil

	case lintutil.RulenameSchemaPropEnumStyleCamelCase:
		return ruleschemapropenumstyle.NewRule(stringcase.CamelCase)
	case lintutil.RulenameSchemaPropEnumStyleKebabCase:
		return ruleschemapropenumstyle.NewRule(stringcase.KebabCase)
	case lintutil.RulenameSchemaPropEnumStylePascalCase:
		return ruleschemapropenumstyle.NewRule(stringcase.PascalCase)
	case lintutil.RulenameSchemaPropEnumStyleSnakeCase:
		return ruleschemapropenumstyle.NewRule(stringcase.SnakeCase)

	case lintutil.RulenameTagStyleFirstUpperCase:
		return ruletagstylefirstuppercase.NewRule(), nil
	}
	return EmptyRule{}, fmt.Errorf("NewStandardRule: rule [%s] not found", name)
}
*/
