package openapi3lint

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grokify/mogo/text/stringcase"
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

type RuleCollectionStandard struct {
	name      string
	ruleNames map[string]int
}

func NewRuleCollectionStandard() RuleCollectionStandard {
	rules := RuleCollectionStandard{
		name:      "Spectrum OpenAPI 3 Lint Standard Rule Collection",
		ruleNames: map[string]int{}}
	names := rules.RuleNames()
	for _, name := range names {
		rules.ruleNames[name] = 1
	}
	return rules
}

func (std RuleCollectionStandard) Name() string {
	return std.name
}

func (std RuleCollectionStandard) RuleExists(ruleName string) bool {
	if _, ok := std.ruleNames[ruleName]; ok {
		return true
	}
	return false
}

func (std RuleCollectionStandard) RuleNames() []string {
	rulenames := []string{
		lintutil.RulenameDatatypeIntFormatStandardExist,
		lintutil.RulenameOpIDStyleCamelCase,
		lintutil.RulenameOpIDStyleKebabCase,
		lintutil.RulenameOpIDStylePascalCase,
		lintutil.RulenameOpIDStyleSnakeCase,
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

func (std RuleCollectionStandard) Rule(name string) (Rule, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case lintutil.RulenameDatatypeIntFormatStandardExist:
		return ruleintstdformat.NewRule(), nil

	case lintutil.RulenameOpIDStyleCamelCase:
		return ruleopidstyle.NewRule(stringcase.CamelCase)
	case lintutil.RulenameOpIDStyleKebabCase:
		return ruleopidstyle.NewRule(stringcase.KebabCase)
	case lintutil.RulenameOpIDStylePascalCase:
		return ruleopidstyle.NewRule(stringcase.PascalCase)
	case lintutil.RulenameOpIDStyleSnakeCase:
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
