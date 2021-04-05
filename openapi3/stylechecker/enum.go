package stylechecker

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/text/stringcase"
)

func SpecCheckSchemas(spec *oas3.Swagger, rules RuleSet) PolicyViolationsSets {
	vsets := NewPolicyViolationsSets()

	enumRules := rules.RulesWithPrefix(PrefixSchemaPropertyEnum)

	for _, enumRule := range enumRules {
		vsets.UpsertSets(SpecCheckSchemaPropertyEnumCaseStyle(
			spec, enumRule))
	}
	if rules.HasRule(RuleSchemaObjectPropsExist) {
		vsets.UpsertSets(SpecCheckSchemaObjectPropsExist(
			spec))
	}

	return vsets
}

func SpecCheckSchemaObjectPropsExist(spec *oas3.Swagger) PolicyViolationsSets {
	vsets := NewPolicyViolationsSets()
	for schemaName, schemaRef := range spec.Components.Schemas {
		if schemaRef.Value == nil {
			continue
		}
		if schemaRef.Value.Type == "object" &&
			len(schemaRef.Value.Properties) == 0 {
			vsets.AddSimple(
				RuleSchemaObjectPropsExist,
				PointerSubEscapeAll("#/components/schemas/%s", schemaName),
				schemaName)
		}
	}
	return vsets
}

func SpecCheckSchemaPropertyEnumCaseStyle(spec *oas3.Swagger, rule string) PolicyViolationsSets {
	vsets := NewPolicyViolationsSets()

	if spec == nil {
		return vsets
	}
	caseStyle := RuleToCaseStyle(rule)
	if len(strings.TrimSpace(caseStyle)) == 0 {
		return vsets
	}

	for schemaName, schemaRef := range spec.Components.Schemas {
		if schemaRef.Value == nil {
			continue
		}
		for propName, propRef := range schemaRef.Value.Properties {
			if propRef.Value == nil ||
				propRef.Value.Type != "string" ||
				len(propRef.Value.Enum) == 0 {
				continue
			}
			for i, enumValue := range propRef.Value.Enum {
				location := fmt.Sprintf(
					"#/components/schemas/%s/properties/%s/%d",
					schemaName, propName, i)
				if enumValueString, ok := enumValue.(string); ok {
					ok, err := stringcase.IsCase(caseStyle, enumValueString)
					if err != nil {
						vsets.AddSimple(rule, location, err.Error())
					} else if !ok {
						vsets.AddSimple(rule, location, enumValueString)
					}
				}
			}
		}
	}
	return vsets
}
