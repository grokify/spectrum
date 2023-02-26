package ruleopxpropertystringexist

import (
	"errors"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/encoding/jsonpointer"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleOperationXPropertyStringExist struct {
	name          string
	xPropertyName string
	inclSpec      bool
	inclPathItem  bool
}

func NewRule(ruleName, xPropertyName string, inclSpec, inclPathItem bool) (RuleOperationXPropertyStringExist, error) {
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
		xPropertyName: xPropertyName,
		inclSpec:      inclSpec,
		inclPathItem:  inclPathItem}
	return rule, nil
}

func (rule RuleOperationXPropertyStringExist) Name() string {
	return rule.name
}

func (rule RuleOperationXPropertyStringExist) Scope() string {
	return lintutil.ScopeSpecification
}

func (rule RuleOperationXPropertyStringExist) ProcessOperation(spec *openapi3.Spec, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}

func (rule RuleOperationXPropertyStringExist) ProcessSpec(spec *openapi3.Spec, pointerBase string) []lintutil.PolicyViolation {
	vios := []lintutil.PolicyViolation{}
	if spec == nil {
		return vios
	}
	if rule.inclSpec {
		propVal := strings.TrimSpace(openapi3.GetExtensionPropStringOrEmpty(
			spec.Extensions, rule.xPropertyName))
		// spec.ExtensionProps, rule.xPropertyName))
		if len(propVal) > 0 {
			return vios
		}
	}
	for pathURL, pathItem := range spec.Paths {
		if pathItem == nil {
			continue
		}
		if rule.inclPathItem {
			propVal := strings.TrimSpace(openapi3.GetExtensionPropStringOrEmpty(
				pathItem.Extensions, rule.xPropertyName))
			// pathItem.ExtensionProps, rule.xPropertyName))
			if len(propVal) > 0 {
				continue
			}
		}
		openapi3.VisitOperationsPathItem(pathURL, pathItem,
			func(path, method string, op *oas3.Operation) {
				if op == nil {
					return
				}
				propVal := strings.TrimSpace(openapi3.GetExtensionPropStringOrEmpty(
					op.Extensions, rule.xPropertyName))
				// op.ExtensionProps, rule.xPropertyName))
				if len(propVal) == 0 {
					vios = append(vios, lintutil.PolicyViolation{
						RuleName: rule.Name(),
						Location: jsonpointer.PointerSubEscapeAll(
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
