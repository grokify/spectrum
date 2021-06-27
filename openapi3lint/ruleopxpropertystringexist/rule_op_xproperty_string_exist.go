package ruleopxpropertystringexist

import (
	"errors"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
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
	return lintutil.ScopeSpecification
}

func (rule RuleOperationXPropertyStringExist) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	/*
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
	*/
	return []lintutil.PolicyViolation{}
}

func (rule RuleOperationXPropertyStringExist) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	vios := []lintutil.PolicyViolation{}
	if spec == nil {
		return []lintutil.PolicyViolation{}
	}
	propVal := strings.TrimSpace(openapi3.GetExtensionPropStringOrEmpty(
		spec.ExtensionProps, rule.xPropertyName))
	if len(propVal) > 0 {
		return []lintutil.PolicyViolation{}
	}
	for pathURL, pathItem := range spec.Paths {
		if pathItem == nil {
			continue
		}
		propVal := strings.TrimSpace(openapi3.GetExtensionPropStringOrEmpty(
			pathItem.ExtensionProps, rule.xPropertyName))
		if len(propVal) > 0 {
			continue
		}

		openapi3.VisitOperationsPathItem(pathURL, pathItem,
			func(path, method string, op *oas3.Operation) {
				if op == nil {
					return
				}
				propVal := strings.TrimSpace(openapi3.GetExtensionPropStringOrEmpty(
					op.ExtensionProps, rule.xPropertyName))
				if len(propVal) == 0 {
					vios = append(vios, lintutil.PolicyViolation{
						RuleName: rule.Name(),
						Location: jsonutil.PointerSubEscapeAll(
							"%s#/paths/%s/%s/%s",
							pointerBase, pathURL, method, rule.xPropertyName,
						),
					})
				}
			},
		)
	}
	return vios
}
