package openapi3edit

import (
	"errors"
	"regexp"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/stringsutil/transform"
	"github.com/grokify/spectrum/openapi3"
)

func (se *SpecEdit) SchemaKeysModify(xf func(string) string) error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	} else if xf == nil {
		return nil
	}
	spec := se.SpecMore.Spec
	specSchemaKeysModifySchemaRefs(spec, xf)
	err := specSchemaKeysModifySchemaKeys(spec, xf)
	if err != nil {
		return errorsutil.Wrap(err, "SpecSchemaKeysModifySchemaKeys")
	}
	return se.SpecMore.Validate()
}

func specSchemaKeysModifySchemaRefs(spec *openapi3.Spec, xf func(string) string) {
	se := NewSpecEdit(spec)
	se.SchemaRefsModify(FuncSchemaRefModFromSchemaKeyMod(xf))
}

var rxComponentSchemasKey = regexp.MustCompile(`^(.*#/components/schemas/)([^/]+)(.*)$`)

// FuncSchemaRefModFromSchemaKeyMod takles a function for modifying schema keys and turns
// it into a function for modifying JSON schema pointers for schemas keys.
func FuncSchemaRefModFromSchemaKeyMod(xf func(string) string) func(string) string {
	return func(s string) string {
		m := rxComponentSchemasKey.FindStringSubmatch(s)
		if len(m) == 0 {
			return s
		}
		substr := m[2]
		mod := xf(substr)
		if substr == mod {
			return s
		}
		return m[1] + mod + m[3]
	}
}

// specSchemaKeysModifySchemaKeys only modifies keys in `components.schemas`. Running this
// by itself does not result in a validating spec.
func specSchemaKeysModifySchemaKeys(spec *openapi3.Spec, xf func(string) string) error {
	schKeys := maputil.StringKeys(spec.Components.Schemas, nil)
	xfMap, _, err := transform.TransformMap(xf, schKeys)
	if err != nil {
		return err
	}
	if !maputil.UniqueValues(xfMap) {
		return errors.New("collisions")
	}
	for _, schKey := range schKeys {
		schRef, ok := spec.Components.Schemas[schKey]
		if !ok {
			panic("schema key not found")
		}
		newSchKey := xf(schKey)
		if schKey == newSchKey {
			continue
		}
		spec.Components.Schemas[newSchKey] = schRef
		delete(spec.Components.Schemas, schKey)
	}
	schKeysNew := maputil.StringKeys(spec.Components.Schemas, nil)
	if len(schKeys) != len(schKeysNew) {
		return errors.New("old and new key mismatch")
	}
	return nil
}
