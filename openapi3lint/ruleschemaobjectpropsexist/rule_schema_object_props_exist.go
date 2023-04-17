// ruleschemaobjectpropsexist ensures that schema objects have properties.
package ruleschemaobjectpropsexist

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/encoding/jsonpointer"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleSchemaObjectPropsExist struct {
	name string
}

func NewRule() RuleSchemaObjectPropsExist {
	return RuleSchemaObjectPropsExist{
		name: lintutil.RulenameSchemaObjectPropsExist}
}

func (rule RuleSchemaObjectPropsExist) Name() string {
	return rule.name
}

func (rule RuleSchemaObjectPropsExist) Scope() string {
	return lintutil.ScopeSpecification
}

func (rule RuleSchemaObjectPropsExist) ProcessOperation(spec *openapi3.Spec, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return nil
}

func (rule RuleSchemaObjectPropsExist) ProcessSpec(spec *openapi3.Spec, pointerBase string) []lintutil.PolicyViolation {
	vios := []lintutil.PolicyViolation{}

	for schName, schRef := range spec.Components.Schemas {
		if schRef == nil || schRef.Value == nil || schRef.Value.Type != openapi3.TypeObject {
			continue
		}
		/*
			if len(schRef.Value.Properties) == 0 && schRef.Value.AdditionalProperties == nil &&
				(schRef.Value.AdditionalPropertiesAllowed == nil || !*schRef.Value.AdditionalPropertiesAllowed) {
				vios = append(vios, lintutil.PolicyViolation{
					RuleName: rule.Name(),
					Location: jsonpointer.PointerSubEscapeAll(
						"%s#/components/schemas/%s",
						pointerBase, schName)})
			}
		*/
		if len(schRef.Value.Properties) == 0 && !openapi3.AdditionalPropertiesExists(schRef.Value.AdditionalProperties) {
			// (!openapi3.AdditionalPropertiesAllowed(schRef.Value.AdditionalProperties) ||
			// 	len(schRef.Value.AdditionalProperties) == 0) {
			vios = append(vios, lintutil.PolicyViolation{
				RuleName: rule.Name(),
				Location: jsonpointer.PointerSubEscapeAll(
					"%s#/components/schemas/%s",
					pointerBase, schName)})
		}
		for propName, propRef := range schRef.Value.Properties {
			if propRef == nil || propRef.Value == nil || propRef.Value.Type != openapi3.TypeObject {
				continue
			}
			/*
				if len(propRef.Value.Properties) == 0 &&
					propRef.Value.AdditionalProperties == nil &&
					(propRef.Value.AdditionalPropertiesAllowed == nil || !*propRef.Value.AdditionalPropertiesAllowed) {
					vios = append(vios, lintutil.PolicyViolation{
						RuleName: rule.Name(),
						Location: jsonpointer.PointerSubEscapeAll(
							"%s#/components/schemas/%s/properties/%s",
							pointerBase, schName, propName)})
				}
			*/
			if len(propRef.Value.Properties) == 0 && !openapi3.AdditionalPropertiesExists(propRef.Value.AdditionalProperties) {
				// propRef.Value.AdditionalProperties == nil &&
				// (propRef.Value.AdditionalPropertiesAllowed == nil || !*propRef.Value.AdditionalPropertiesAllowed) {
				vios = append(vios, lintutil.PolicyViolation{
					RuleName: rule.Name(),
					Location: jsonpointer.PointerSubEscapeAll(
						"%s#/components/schemas/%s/properties/%s",
						pointerBase, schName, propName)})
			}
		}
	}
	return vios
}
