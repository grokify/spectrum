package openapi3

import (
	"strings"
)

// ExtensionNames is not complete yet.
func (sm *SpecMore) ExtensionNames() map[string]int {
	extNames := map[string]int{}
	for _, schRef := range sm.Spec.Components.Schemas {
		if schRef.Value == nil {
			continue
		}
		sch := schRef.Value
		for extName := range sch.Extensions {
			// for extName := range sch.ExtensionProps.Extensions {
			count, ok := extNames[extName]
			if !ok {
				extNames[extName] = 1
			}
			extNames[extName] = count + 1
		}
	}
	return extNames
}

func (sm *SpecMore) HasComponentSchema(componentSchemaName string, caseInsensitiveCaseMatch bool) bool {
	componentSchemaName = strings.TrimSpace(componentSchemaName)
	if caseInsensitiveCaseMatch {
		componentSchemaName = strings.ToLower(componentSchemaName)
	}
	if len(sm.Spec.Components.Schemas) == 0 {
		return false
	}
	if _, ok := sm.Spec.Components.Schemas[componentSchemaName]; ok {
		return true
	}
	if caseInsensitiveCaseMatch {
		for nameTry := range sm.Spec.Components.Schemas {
			if strings.EqualFold(strings.TrimSpace(nameTry), componentSchemaName) {
				return true
			}
		}
	}
	return false
}
