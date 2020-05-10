package openapi3

import (
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

// SpecExtensionNames is not complete yet.
func SpecExtensionNames(spec *oas3.Swagger) map[string]int {
	extNames := map[string]int{}
	for _, schRef := range spec.Components.Schemas {
		if schRef.Value == nil {
			continue
		}
		sch := schRef.Value
		for extName := range sch.ExtensionProps.Extensions {
			count, ok := extNames[extName]
			if !ok {
				extNames[extName] = 1
			}
			extNames[extName] = count + 1
		}
	}
	return extNames
}

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
