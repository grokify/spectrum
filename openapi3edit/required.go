package openapi3edit

import (
	"regexp"

	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
)

// SchemaPropertiesSetOptional sets properites as optional if the description matches a regexp
// such as var rxOptionalDefault = regexp.MustCompile(`(?i)\boptional\b`)
func (se *SpecEdit) SchemaPropertiesSetOptional(rxOptional *regexp.Regexp) error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	}
	if rxOptional == nil {
		return nil
	}
	for _, schemaRef := range se.SpecMore.Spec.Components.Schemas {
		if len(schemaRef.Ref) == 0 && schemaRef.Value != nil {
			required := []string{}
			for propName, propRef := range schemaRef.Value.Properties {
				if len(propRef.Ref) == 0 && propRef.Value != nil {
					if len(propRef.Value.Description) > 0 &&
						!rxOptional.MatchString(propRef.Value.Description) {
						required = append(required, propName)
					}
				}
			}
			if len(required) > 1 {
				required = stringsutil.SliceCondenseSpace(required, true, true)
			}
			schemaRef.Value.Required = required
		}
	}
	return nil
}
