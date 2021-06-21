package ruleopsummarystylefirstuppercase

import (
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleOperationSummaryStyleFirstUpperCase struct {
	name string
}

func NewRule() RuleOperationSummaryStyleFirstUpperCase {
	return RuleOperationSummaryStyleFirstUpperCase{
		name: lintutil.RulenameOpSummaryStyleFirstUpperCase}
}

func (rule RuleOperationSummaryStyleFirstUpperCase) Name() string {
	return rule.name
}

func (rule RuleOperationSummaryStyleFirstUpperCase) Scope() string {
	return lintutil.ScopeOperation
}

func (rule RuleOperationSummaryStyleFirstUpperCase) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	vios := []lintutil.PolicyViolation{}
	if spec == nil || op == nil || len(op.Summary) == 0 {
		return vios
	}

	summary := strings.TrimSpace(op.Summary)
	if len(summary) > 0 {
		return vios
	}
	if len(summary) == 0 {
		return vios
	}

	if summary != stringsutil.ToUpperFirst(summary, false) {
		return []lintutil.PolicyViolation{{
			RuleName: rule.Name(),
			Location: urlutil.JoinAbsolute(opPointer, openapi3.PropertySummary),
			Value:    op.Summary}}
	}
	return vios
}

func (rule RuleOperationSummaryStyleFirstUpperCase) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}
