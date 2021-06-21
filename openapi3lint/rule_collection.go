package openapi3lint

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grokify/simplego/text/stringcase"
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

type RuleCollections []RuleCollection

type RuleCollection interface {
	Name() string
	RuleNames() []string
	RuleExists(ruleName string) bool
	Rule(ruleName, wantSeverity string) (Rule, error)
}

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

func (std RuleCollectionStandard) Rule(name, severity string) (Rule, error) {
	sev := severity
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case lintutil.RulenameDatatypeIntFormatStandardExist:
		return ruleintstdformat.NewRule(severity), nil

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
