package ruleopxpropertystringexist

import (
	"errors"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

const (
	RuleName = "x-operation-x-api-group-exist"
)

type RuleOperationXPropertyStringExist struct {
	name          string
	xPropertyName string
}

func NewRule(ruleName, xPropertyName string) (RuleOperationXPropertyStringExist, error) {
	ruleName = strings.ToLower(strings.TrimSpace(ruleName))
	xPropertyName = strings.TrimSpace(xPropertyName)

	if len(ruleName) == 0 {
		return RuleOperationXPropertyStringExist{},
			errors.New("rule name not provided")
	}
	if len(xPropertyName) == 0 {
		return RuleOperationXPropertyStringExist{},
			errors.New("x-property name not provided")
	}

	rule := RuleOperationXPropertyStringExist{
		name:          ruleName,
		xPropertyName: xPropertyName}
	return rule, nil
}

func (rule RuleOperationXPropertyStringExist) Name() string {
	return rule.name
}

func (rule RuleOperationXPropertyStringExist) Scope() string {
	return lintutil.ScopeOperation
}

func (rule RuleOperationXPropertyStringExist) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	if spec == nil || op == nil || len(op.OperationID) == 0 {
		return nil
	}
	prop := strings.TrimSpace(
		openapi3.GetOperationExtensionPropStringOrEmpty(*op, rule.xPropertyName))
	if len(prop) > 0 {
		return nil
	}
	vio := lintutil.PolicyViolation{
		RuleName: rule.Name(),
		Location: urlutil.JoinAbsolute(opPointer, "operationId"),
		Value:    op.OperationID}
	return []lintutil.PolicyViolation{vio}
}

func (rule RuleOperationXPropertyStringExist) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}
