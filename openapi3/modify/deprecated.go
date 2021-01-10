package modify

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/swaggman/openapi3"
)

func SpecSchemasSetDeprecated(spec *oas3.Swagger, newDeprecated bool) {
	for _, schemaRef := range spec.Components.Schemas {
		if len(schemaRef.Ref) == 0 && schemaRef.Value != nil {
			schemaRef.Value.Deprecated = newDeprecated
		}
	}
}

func SpecOperationsSetDeprecated(spec *oas3.Swagger, newDeprecated bool) {
	openapi3.VisitOperations(
		spec,
		func(path, method string, op *oas3.Operation) {
			if op != nil {
				op.Deprecated = newDeprecated
			}
		},
	)
}
