package openapi3lint

import (
	"github.com/grokify/simplego/log/severity"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

func RulesConfigExample1() map[string]RuleConfig {
	rules := map[string]RuleConfig{
		lintutil.RulenameOpIdStyleCamelCase:            {},
		lintutil.RulenameOpSummaryExist:                {},
		lintutil.RulenameOpSummaryStyleFirstUpperCase:  {},
		lintutil.RulenamePathParamStylePascalCase:      {},
		lintutil.RulenameSchemaObjectPropsExist:        {},
		lintutil.RulenameSchemaPropEnumStylePascalCase: {},
		lintutil.RulenameTagStyleFirstUpperCase:        {},
	}
	for name, cfg := range rules {
		cfg.Severity = severity.SeverityError
		rules[name] = cfg
	}
	return rules
}
