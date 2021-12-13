package openapi3

import (
	"strings"

	"github.com/grokify/mogo/encoding/jsonutil"
)

const PointerComponentsSchemas = "#/components/schemas"

func SchemaPointerExpand(prefix, schemaName string) string {
	// https://swagger.io/docs/specification/components/
	prefix = strings.TrimSpace(prefix)
	schemaName = strings.TrimSpace(schemaName)
	pointer := schemaName
	if strings.Index(schemaName, PointerComponentsSchemas) < 0 {
		pointer = PointerComponentsSchemas + "/" + jsonutil.PropertyNameEscape(schemaName)
	}
	if len(prefix) > 0 && strings.Index(pointer, "#") == 0 {
		pointer = prefix + pointer
	}
	return pointer
}
