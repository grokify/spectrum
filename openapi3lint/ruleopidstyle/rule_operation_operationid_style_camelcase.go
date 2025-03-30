package ruleopidstyle

import (
	"fmt"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/text/stringcase"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleOperationOperationIDStyle struct {
	name       string
	stringCase string
}

func NewRule(requiredStringCase string) (RuleOperationOperationIDStyle, error) {
	canonicalCase, err := stringcase.Parse(requiredStringCase)
	if err != nil {
		return RuleOperationOperationIDStyle{},
			fmt.Errorf("invalid string case [%s]", requiredStringCase)
	}
	rule := RuleOperationOperationIDStyle{
		stringCase: canonicalCase}
	switch canonicalCase {
	case stringcase.CamelCase:
		rule.name = lintutil.RulenameOpIDStyleCamelCase
	case stringcase.KebabCase:
		rule.name = lintutil.RulenameOpIDStyleKebabCase
	case stringcase.PascalCase:
		rule.name = lintutil.RulenameOpIDStylePascalCase
	case stringcase.SnakeCase:
		rule.name = lintutil.RulenameOpIDStyleSnakeCase
	default:
		return rule, fmt.Errorf("invalid string case [%s]", canonicalCase)
	}
	return rule, nil
}

func (rule RuleOperationOperationIDStyle) Name() string {
	return rule.name
}

func (rule RuleOperationOperationIDStyle) Scope() string {
	return lintutil.ScopeOperation
}

func (rule RuleOperationOperationIDStyle) ProcessOperation(spec *openapi3.Spec, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	if spec == nil || op == nil || len(op.OperationID) == 0 {
		return nil
	}
	isWantCase, err := stringcase.IsCase(rule.stringCase, op.OperationID)
	if err == nil && isWantCase {
		return nil
	}
	vio := lintutil.PolicyViolation{
		RuleName: rule.Name(),
		Location: urlutil.JoinAbsolute(opPointer, openapi3.PropertyOperationID),
		Value:    op.OperationID}
	if err != nil {
		vio.Data = map[string]string{
			"error": err.Error()}
	}
	return []lintutil.PolicyViolation{vio}
}

func (rule RuleOperationOperationIDStyle) ProcessSpec(spec *openapi3.Spec, pointerBase string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}
