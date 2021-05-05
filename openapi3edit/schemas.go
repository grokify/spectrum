package openapi3edit

import oas3 "github.com/getkin/kin-openapi/openapi3"

func SpecSchemaNames(spec *oas3.Swagger) []string {
	schemas := []string{}
	for name := range spec.Components.Schemas {
		schemas = append(schemas, name)
	}
	return schemas
}
