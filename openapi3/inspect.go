package openapi3

import (
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

func SpecHasComponentSchema(spec *oas3.Swagger, name string, lowerCaseMatch bool) bool {
	name = strings.TrimSpace(name)
	if lowerCaseMatch {
		name = strings.ToLower(name)
	}
	if len(spec.Components.Schemas) == 0 {
		return false
	}
	if _, ok := spec.Components.Schemas[name]; ok {
		return true
	}
	if lowerCaseMatch {
		for nameTry := range spec.Components.Schemas {
			nameTry = strings.ToLower(strings.TrimSpace(nameTry))
			if nameTry == name {
				return true
			}
		}
	}
	return false
}
