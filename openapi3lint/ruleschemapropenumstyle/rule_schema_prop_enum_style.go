package ruleschemapropenumstyle

import (
	"fmt"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/simplego/text/stringcase"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleSchemaPropEnumStyle struct {
	name       string
	stringCase string
}

func NewRule(requiredStringCase string) (RuleSchemaPropEnumStyle, error) {
	canonicalCase, err := stringcase.Parse(requiredStringCase)
	if err != nil {
		return RuleSchemaPropEnumStyle{},
			fmt.Errorf("invalid string case [%s]", requiredStringCase)
	}
	rule := RuleSchemaPropEnumStyle{
		stringCase: canonicalCase}
	switch canonicalCase {
	case stringcase.CamelCase:
		rule.name = lintutil.RulenameSchemaPropEnumStyleCamelCase
	case stringcase.KebabCase:
		rule.name = lintutil.RulenameSchemaPropEnumStyleKebabCase
	case stringcase.PascalCase:
		rule.name = lintutil.RulenameSchemaPropEnumStylePascalCase
	case stringcase.SnakeCase:
		rule.name = lintutil.RulenameSchemaPropEnumStyleSnakeCase
	default:
		return rule, fmt.Errorf("invalid string case [%s]", canonicalCase)
	}
	return rule, nil
}

func (rule RuleSchemaPropEnumStyle) Name() string {
	return rule.name
}

func (rule RuleSchemaPropEnumStyle) Scope() string {
	return lintutil.ScopeSpecification
}

func (rule RuleSchemaPropEnumStyle) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return nil
}

func (rule RuleSchemaPropEnumStyle) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	vios := []lintutil.PolicyViolation{}

	for schName, schRef := range spec.Components.Schemas {
		if schRef == nil || schRef.Value == nil ||
			schRef.Value.Type != openapi3.TypeObject {
			continue
		}

		for propName, propRef := range schRef.Value.Properties {
			if propRef.Value == nil ||
				propRef.Value.Type != "string" ||
				len(propRef.Value.Enum) == 0 {
				continue
			}
			for i, enumValue := range propRef.Value.Enum {
				if enumValueString, ok := enumValue.(string); ok {
					jsPtr := jsonutil.PointerSubEscapeAll(
						"%s#/components/schemas/%s/properties/%s/%d",
						pointerBase, schName, propName, i)
					isWantCase, err := stringcase.IsCase(rule.stringCase, enumValueString)
					if err != nil {
						// should never happen as rule.stringCase should be validated.
						vios = append(vios, lintutil.PolicyViolation{
							RuleName: rule.Name(),
							Location: jsPtr,
							Value:    enumValueString})
					} else if !isWantCase {
						vios = append(vios, lintutil.PolicyViolation{
							RuleName: rule.Name(),
							Location: jsPtr,
							Value:    enumValueString})
					}
				}
			}
		}
	}

	return vios
}
