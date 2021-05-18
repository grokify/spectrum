package ruleopidstyle

import (
	"fmt"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/simplego/text/stringcase"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleOperationOperationIdStyle struct {
	name       string
	stringCase string
	severity   string
}

func NewRuleOperationOperationIdStyle(sev, requiredStringCase string) (RuleOperationOperationIdStyle, error) {
	canonicalCase, err := stringcase.Parse(requiredStringCase)
	if err != nil {
		return RuleOperationOperationIdStyle{},
			fmt.Errorf("invalid string case [%s]", requiredStringCase)
	}
	rule := RuleOperationOperationIdStyle{
		name:       lintutil.RulenameDatatypeIntFormatIsInt32OrInt64,
		stringCase: canonicalCase,
		severity:   sev}
	switch canonicalCase {
	case stringcase.CamelCase:
		rule.name = lintutil.RulenameOpIdStyleCamelCase
	case stringcase.KebabCase:
		rule.name = lintutil.RulenameOpIdStyleKebabCase
	case stringcase.PascalCase:
		rule.name = lintutil.RulenameOpIdStylePascalCase
	case stringcase.SnakeCase:
		rule.name = lintutil.RulenameOpIdStyleSnakeCase
	default:
		return rule, fmt.Errorf("invalid string case [%s]", canonicalCase)
	}
	return rule, nil
}

func (rule RuleOperationOperationIdStyle) Name() string {
	return rule.name
}

func (rule RuleOperationOperationIdStyle) Scope() string {
	return lintutil.ScopeOperation
}

func (rule RuleOperationOperationIdStyle) Severity() string {
	return rule.severity
}

func (rule RuleOperationOperationIdStyle) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	if spec == nil || op == nil || len(op.OperationID) == 0 {
		return nil
	}
	isWantCase, err := stringcase.IsCase(rule.stringCase, op.OperationID)
	if err == nil && isWantCase {
		return nil
	}
	vio := lintutil.PolicyViolation{
		RuleName: rule.Name(),
		Location: urlutil.JoinAbsolute(opPointer, "operationId"),
		Value:    op.OperationID}
	if err != nil {
		vio.Data = map[string]string{
			"error": err.Error()}
	}
	return []lintutil.PolicyViolation{vio}
}

func (rule RuleOperationOperationIdStyle) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}
