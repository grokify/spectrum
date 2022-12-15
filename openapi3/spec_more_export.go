package openapi3

import (
	"errors"
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/http/pathmethod"
)

// ExportByTags creates individual specs by tag.
func (sm *SpecMore) ExportByTags() (map[string]*Spec, error) {
	specs := map[string]*Spec{}
	if sm.Spec == nil {
		return specs, ErrSpecNotSet
	}
	tags := sm.Tags(false, true)

	for _, tag := range tags {
		tagSpec, err := sm.ExportByTag(tag)
		if err != nil {
			return specs, errorsutil.Wrapf(err, "error SpecMore.ExportByTag(\"%s\")", tag)
		}
		if tagSpec != nil {
			specs[tag] = tagSpec
		}
	}
	return specs, nil
}

// ExportByTag creates an individual specs for one tag.
func (sm *SpecMore) ExportByTag(tag string) (*Spec, error) {
	if (sm.Spec) == nil {
		return nil, ErrSpecNotSet
	}
	oms := sm.OperationMetas([]string{tag})
	if len(oms) == 0 {
		return nil, nil
	}
	tagSpec := &Spec{
		Components: oas3.Components{
			ExtensionProps:  sm.Spec.Components.ExtensionProps,
			SecuritySchemes: sm.Spec.Components.SecuritySchemes,
		},
		Info:    sm.Spec.Info,
		Servers: sm.Spec.Servers,
		OpenAPI: sm.Spec.OpenAPI,
	}
	for _, om := range oms {
		op, err := sm.OperationByPathMethod(om.Path, om.Method)
		if err != nil {
			return nil, errorsutil.Wrapf(err, "error `OperationByPathMethod` pathmethod: (%s)", pathmethod.PathMethod(om.Path, om.Method))
		} else if op == nil {
			continue
		}
		tagSpec.AddOperation(om.Path, om.Method, op)
		err = sm.SchemasCopyOperation(tagSpec, op)
		if err != nil {
			return nil, errorsutil.Wrapf(err, "error `SpecMore.SchemasCopyOperation()` tag: (%s)", tag)
		}
	}
	tagSm := SpecMore{Spec: tagSpec}
	return tagSm.Clone()
}

var ErrJSONPointerNotParamOrSchema = errors.New("pointer is not components/parameters or components/schemas")

func (sm *SpecMore) SchemasCopyOperation(destSpec *Spec, op *oas3.Operation) error {
	if sm.Spec == nil || destSpec == nil || op == nil {
		return errors.New("source spec, dest spec, op cannot be nil")
	}
	omr := OperationMore{Operation: op}
	refs := omr.JSONPointers()
	for refJSONPointer := range refs {
		ptr, err := ParseJSONPointer(refJSONPointer)
		if err != nil {
			return errorsutil.Wrapf(err, "error ParseJSONPointer() jsonpointer: (%s)", refJSONPointer)
		}
		paramName, isParam := ptr.IsTopParameter()
		if isParam {
			if paramRef, ok := sm.Spec.Components.Parameters[paramName]; ok {
				if destSpec.Components.Parameters == nil {
					destSpec.Components.Parameters = oas3.ParametersMap{}
				}
				destSpec.Components.Parameters[paramName] = paramRef
			}
		}
		_, isSchema := ptr.IsTopSchema()
		if isSchema {
			err := sm.SchemasCopyJSONPointer(destSpec, refJSONPointer)
			if err != nil {
				return err
			}
		}
		if !isParam && !isSchema {
			return errorsutil.Wrapf(ErrJSONPointerNotParamOrSchema, "jsonpointer: (%s)", refJSONPointer)
		}
	}
	return nil
}

func (sm *SpecMore) SchemasCopySchemaRef(destSpec *Spec, schRef *oas3.SchemaRef) error {
	if sm.Spec == nil || destSpec == nil || schRef == nil {
		return nil
	}
	if len(strings.TrimSpace(schRef.Ref)) > 0 {
		err := sm.SchemasCopyJSONPointer(destSpec, schRef.Ref)
		if err != nil {
			return err
		}
	}
	if schRef.Value == nil {
		return nil
	}
	if schRef.Value.Items != nil {
		err := sm.SchemasCopySchemaRef(destSpec, schRef.Value.Items)
		if err != nil {
			return err
		}
	}
	for _, schRefProp := range schRef.Value.Properties {
		err := sm.SchemasCopySchemaRef(destSpec, schRefProp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sm *SpecMore) SchemasCopyJSONPointer(destSpec *Spec, jsonPointer string) error {
	ptr, err := ParseJSONPointer(jsonPointer)
	if err != nil {
		return err
	}
	schName, ok := ptr.IsTopSchema()
	if !ok {
		return errors.New("json pointer is not schema pointer")
	}
	destSM := SpecMore{Spec: destSpec}
	destSchRef := destSM.SchemaRef(schName)
	if destSchRef != nil { // already present.
		return nil
	}
	srcSchRef := sm.SchemaRef(schName)
	if srcSchRef == nil {
		return fmt.Errorf("json pointer not found [%s][%s]", jsonPointer, schName)
	}
	err = destSM.SchemaRefSet(schName, srcSchRef)
	if err != nil {
		return err
	}
	// Check recursive.
	return sm.SchemasCopySchemaRef(destSpec, srcSchRef)
}
