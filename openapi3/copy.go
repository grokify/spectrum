package openapi3

import (
	"encoding/json"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

func CopySchemaStandard(schema oas3.Schema) (oas3.Schema, error) {
	bytes, err := json.Marshal(schema)
	if err != nil {
		return oas3.Schema{}, err
	}
	var newSchema oas3.Schema
	err = json.Unmarshal(bytes, &newSchema)
	if err != nil {
		return oas3.Schema{}, err
	}
	newSchema.ExtensionProps = oas3.ExtensionProps{}
	return newSchema, nil
}
