package openapi3

import (
	"os"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/encoding/jsonpointer"
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
