package openapi3edit

import (
	"errors"
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
)

func (se *SpecEdit) SchemasFlatten() {
	if se.SpecMore.Spec == nil {
		return
	}
	for schName, schRef := range se.SpecMore.Spec.Components.Schemas {
		if schRef == nil || schRef.Value == nil {
			continue
		}
		se.SchemasFlattenSchemaRef("", schName, schRef)
	}
}

func (se *SpecEdit) SchemasFlattenSchemaRef(baseName, schName string, schRef *oas3.SchemaRef) {
	if se.SpecMore.Spec == nil || schRef == nil {
		return
	}
	spec := se.SpecMore.Spec

	basePlusSchName := baseName + stringsutil.ToUpperFirst(schName, false)
	if len(baseName) == 0 {
		basePlusSchName = schName
	}

	for propName, propRef := range schRef.Value.Properties {
		if propRef == nil || propRef.Value == nil {
			continue
		}
		if openapi3.TypesRefIs(propRef.Value.Type, openapi3.TypeArray) {
			itemsRef := propRef.Value.Items
			if itemsRef == nil {
				continue
			}
			itemsRef.Ref = strings.TrimSpace(itemsRef.Ref)
			if itemsRef.Value != nil &&
				(openapi3.TypesRefIs(itemsRef.Value.Type, openapi3.TypeArray, openapi3.TypeObject)) {
				if len(itemsRef.Ref) > 0 {
					propRef.Value.Items = oas3.NewSchemaRef(itemsRef.Ref, nil)
				} else {
					newSchemaName := basePlusSchName + stringsutil.ToUpperFirst(propName, false)
					if _, ok := spec.Components.Schemas[newSchemaName]; ok {
						fmt.Printf("BASE_NAME [%s] SCH_NAME [%s] PROP_NAME [%s] ARRAY\n", baseName, schName, propName)
						panic("collision")
					}
					se.SchemasFlattenSchemaRef(basePlusSchName, propName, itemsRef)
					spec.Components.Schemas[newSchemaName] = itemsRef
					propRef.Value.Items = oas3.NewSchemaRef(openapi3.PointerComponentsSchemas+"/"+newSchemaName, nil)
				}
				//itemsRef.Value = nil
				//itemsRef.Ref = openapi3.PointerComponentsSchemas + "/" + newSchemaName
			}
			//} else if propRef.Value.Type == openapi3.TypeObject {
		} else if openapi3.TypesRefIs(propRef.Value.Type, openapi3.TypeObject) {
			if len(propRef.Value.Properties) > 0 {
				newSchemaName := basePlusSchName + stringsutil.ToUpperFirst(propName, false)
				if _, ok := spec.Components.Schemas[newSchemaName]; ok {
					fmt.Printf("BASE_NAME [%s] SCH_NAME [%s] PROP_NAME [%s] OBJECT\n", baseName, schName, propName)
					panic("collision")
				}
				//SpecSchemasFlattenSchemaRef(spec, basePlusSchName, propName, propRef)
				se.SchemasFlattenSchemaRef("", newSchemaName, propRef)
				spec.Components.Schemas[newSchemaName] = propRef
				schRef.Value.Properties[propName] = oas3.NewSchemaRef(openapi3.PointerComponentsSchemas+"/"+newSchemaName, nil)
				//SpecSchemasFlattenSchemaRef(spec, basePlusSchName, propName, propRef)
				//schRef.Value.Properties[propName] = oas3.NewSchemaRef(newSchemaName, nil)
			}
		}
	}
}

var ErrEmptySchemaName = errors.New("empty schema name encountered")

// SchemaRefsFlatten flattens Schema refs.
func (se *SpecEdit) SchemaRefsFlatten() error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	}
	spec := se.SpecMore.Spec
	// func SpecFlattenSchemaRefs(spec *openapi3.Spec, visitSchemaRefFunc func(schName string, schRef *oas3.SchemaRef) error) error {
	for schName, schRef := range spec.Components.Schemas {
		for propSchemaName, propSchemaRef := range schRef.Value.Properties {
			if len(strings.TrimSpace(propSchemaName)) == 0 {
				return ErrEmptySchemaName
			}
			// visitSchemaRefFunc(propSchemaName, propSchemaRef)
			if len(propSchemaRef.Ref) == 0 && propSchemaRef.Value != nil {
				if openapi3.TypesRefIs(propSchemaRef.Value.Type, openapi3.TypeObject, openapi3.TypeArray) {
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
