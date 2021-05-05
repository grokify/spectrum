package openapi3edit

import (
	"regexp"

	"github.com/grokify/simplego/type/stringsutil"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

var rxOptionalDefault = regexp.MustCompile(`(?i)\boptional\b`)

func SpecSetSchemaPropertiesOptional(spec *oas3.Swagger, rxOptional *regexp.Regexp) {
	if rxOptional == nil {
		return
	}
	for _, schemaRef := range spec.Components.Schemas {
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
}
