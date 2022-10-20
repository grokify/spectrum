package openapi3edit

import (
	"errors"
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
)

func SpecSchemasFlatten(spec *openapi3.Spec) {
	for schName, schRef := range spec.Components.Schemas {
		if schRef == nil || schRef.Value == nil {
			continue
		}
		SpecSchemasFlattenSchemaRef(spec, "", schName, schRef)
	}
}

func SpecSchemasFlattenSchemaRef(spec *openapi3.Spec, baseName, schName string, schRef *oas3.SchemaRef) {
	if schRef == nil {
		return
	}

	basePlusSchName := baseName + stringsutil.ToUpperFirst(schName, false)
	if len(baseName) == 0 {
		basePlusSchName = schName
	}

	for propName, propRef := range schRef.Value.Properties {
		if propRef == nil || propRef.Value == nil {
			continue
		}
		if propRef.Value.Type == openapi3.TypeArray {
			itemsRef := propRef.Value.Items
			if itemsRef == nil {
				continue
			}
			itemsRef.Ref = strings.TrimSpace(itemsRef.Ref)
			if itemsRef.Value != nil &&
				(itemsRef.Value.Type == openapi3.TypeArray || itemsRef.Value.Type == openapi3.TypeObject) {
				if len(itemsRef.Ref) > 0 {
					propRef.Value.Items = oas3.NewSchemaRef(itemsRef.Ref, nil)
				} else {
					newSchemaName := basePlusSchName + stringsutil.ToUpperFirst(propName, false)
					if _, ok := spec.Components.Schemas[newSchemaName]; ok {
						fmt.Printf("BASE_NAME [%s] SCH_NAME [%s] PROP_NAME [%s] ARRAY\n", baseName, schName, propName)
						panic("collision")
					}
					SpecSchemasFlattenSchemaRef(spec, basePlusSchName, propName, itemsRef)
					spec.Components.Schemas[newSchemaName] = itemsRef
					propRef.Value.Items = oas3.NewSchemaRef(openapi3.PointerComponentsSchemas+"/"+newSchemaName, nil)
				}
				//itemsRef.Value = nil
				//itemsRef.Ref = openapi3.PointerComponentsSchemas + "/" + newSchemaName
			}
			//} else if propRef.Value.Type == openapi3.TypeObject {
		} else if propRef.Value.Type == openapi3.TypeObject {
			if len(propRef.Value.Properties) > 0 {
				newSchemaName := basePlusSchName + stringsutil.ToUpperFirst(propName, false)
				if _, ok := spec.Components.Schemas[newSchemaName]; ok {
					fmt.Printf("BASE_NAME [%s] SCH_NAME [%s] PROP_NAME [%s] OBJECT\n", baseName, schName, propName)
					panic("collision")
				}
				//SpecSchemasFlattenSchemaRef(spec, basePlusSchName, propName, propRef)
				SpecSchemasFlattenSchemaRef(spec, "", newSchemaName, propRef)
				spec.Components.Schemas[newSchemaName] = propRef
				schRef.Value.Properties[propName] = oas3.NewSchemaRef(openapi3.PointerComponentsSchemas+"/"+newSchemaName, nil)
				//SpecSchemasFlattenSchemaRef(spec, basePlusSchName, propName, propRef)
				//schRef.Value.Properties[propName] = oas3.NewSchemaRef(newSchemaName, nil)
			}
		}
	}
}

var ErrEmptySchemaName = errors.New("empty schema name encountered")

// SpecSchemaRefsFlatten flattens Schema refs.
func SpecSchemaRefsFlatten(spec *openapi3.Spec) error {
	// func SpecFlattenSchemaRefs(spec *openapi3.Spec, visitSchemaRefFunc func(schName string, schRef *oas3.SchemaRef) error) error {
	for schName, schRef := range spec.Components.Schemas {
		for propSchemaName, propSchemaRef := range schRef.Value.Properties {
			if len(strings.TrimSpace(propSchemaName)) == 0 {
				return ErrEmptySchemaName
			}
			// visitSchemaRefFunc(propSchemaName, propSchemaRef)
			if len(propSchemaRef.Ref) == 0 && propSchemaRef.Value != nil {
				if propSchemaRef.Value.Type == openapi3.TypeObject || propSchemaRef.Value.Type == openapi3.TypeArray {
					newRootSchemaName := propSchemaName
					if _, ok := spec.Components.Schemas[newRootSchemaName]; ok {
						newRootSchemaName = schName + stringsutil.ToUpperFirst(newRootSchemaName, false)
						if _, ok := spec.Components.Schemas[newRootSchemaName]; ok {
							return fmt.Errorf("schema collision [%s]", newRootSchemaName)
						}
					}
					spec.Components.Schemas[newRootSchemaName] = propSchemaRef
					propSchemaRef.Value = nil
					propSchemaRef.Ref = openapi3.PointerComponentsSchemas + "/" + newRootSchemaName
				}
			}
		}
	}
	return nil
}
