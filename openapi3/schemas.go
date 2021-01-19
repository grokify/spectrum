package openapi3

import "strings"

const schemaBasePath = "#/components/schemas"

func SchemaPathExpand(schemaName string) string {
	// https://swagger.io/docs/specification/components/
	schemaName = strings.TrimSpace(schemaName)
	idx := strings.Index(schemaName, schemaBasePath)
	if idx < 0 {
		return schemaBasePath + "/" + schemaName
	}
	return schemaName
}
