package openapi3

import (
	"os"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/encoding/jsonpointer"
	"github.com/grokify/mogo/type/stringsutil"
)

func SchemaPointerExpand(prefix, schemaName string) string {
	// https://swagger.io/docs/specification/components/
	prefix = strings.TrimSpace(prefix)
	schemaName = strings.TrimSpace(schemaName)
	pointer := schemaName
	if !strings.Contains(schemaName, PointerComponentsSchemas) {
		pointer = PointerComponentsSchemas + "/" + jsonpointer.PropertyNameEscape(schemaName)
	}
	if len(prefix) > 0 && strings.Index(pointer, "#") == 0 {
		pointer = prefix + pointer
	}
	return pointer
}

func ReadSchemaFile(filename string) (*oas3.Schema, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	sch := oas3.NewSchema()
	err = sch.UnmarshalJSON(data)
	return sch, err
}

// AdditionalPropertiesAllowed checks for additional properties, which exists in Schema structs.
func AdditionalPropertiesAllowed(aprops oas3.AdditionalProperties) bool {
	if aprops.Has != nil {
		return *aprops.Has
	} else {
		return false
	}
}

func AdditionalPropertiesExists(props oas3.AdditionalProperties) bool {
	if props.Has == nil || !*props.Has || props.Schema == nil {
		return false
	}
	if strings.TrimSpace(props.Schema.Ref) != "" {
		return true
	}
	if props.Schema.Value == nil {
		return false
	}
	return true
}

// NewTypes returns a value that is suitable for `oas3.Schema.Value.Type`.
// In `0.123.0` and earlier, this is a `string`. In `0.124.0` and later, this is
// of type `oas3.Types`. `oas3` is `github.com/getkin/kin-openapi/openapi3`.
func NewTypesRef(t ...string) *oas3.Types {
	ot := oas3.Types{}
	for _, ti := range t {
		ot = append(ot, ti)
	}
	return &ot
}

// TypesRefIs returns if the supplied `*oas3.Types` is any of the supplied values.
// It returns false if `*oas3.Types` is false, or `type` an empty slice, or
// none of the supplied types match. `oas3` is `github.com/getkin/kin-openapi/openapi3`.
func TypesRefIs(t *oas3.Types, types ...string) bool {
	types = stringsutil.SliceCondenseSpace(types, true, false)
	if t == nil || len(*t) == 0 {
		if len(types) == 0 {
			return true
		} else {
			return false
		}
	} else if len(types) == 0 {
		return false
	}
	for _, typ := range types {
		if t.Is(typ) {
			return true
		}
	}
	return false
}

// TypesRefString returns a string if `oas3.Types` is not nil and has a single type.
func TypesRefString(t *oas3.Types) string {
	if t == nil {
		return ""
	} else if len(*t) != 1 {
		return ""
	} else {
		return (*t)[0]
	}
}
