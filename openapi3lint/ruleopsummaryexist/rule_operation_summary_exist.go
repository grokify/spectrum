package ruleopsummaryexist

import (
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleOperationSummaryExist struct {
	name string
}

func NewRule() RuleOperationSummaryExist {
	return RuleOperationSummaryExist{
		name: lintutil.RulenameOpSummaryExist}
}

func (rule RuleOperationSummaryExist) Name() string {
	return rule.name
}

func (rule RuleOperationSummaryExist) Scope() string {
	return lintutil.ScopeOperation
}

func (rule RuleOperationSummaryExist) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	vios := []lintutil.PolicyViolation{}
	if spec == nil || op == nil {
		return vios
	}

	summary := strings.TrimSpace(op.Summary)
	if len(summary) > 0 {
		return vios
	}

	return []lintutil.PolicyViolation{{
		RuleName: rule.Name(),
		Location: urlutil.JoinAbsolute(opPointer, openapi3.PropertySummary),
		Value:    op.Summary}}
}

func (rule RuleOperationSummaryExist) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}
