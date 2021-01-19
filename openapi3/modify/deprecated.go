package modify

import (
	"regexp"

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

var rxDeprecated = regexp.MustCompile(`(?i)\bdeprecated\b`)

func SpecSetDeprecatedImplicit(spec *oas3.Swagger) {
	openapi3.VisitOperations(
		spec,
		func(path, method string, op *oas3.Operation) {
			if op != nil && rxDeprecated.MatchString(op.Description) {
				op.Deprecated = true
			}
		},
	)
	for _, schemaRef := range spec.Components.Schemas {
		if len(schemaRef.Ref) == 0 && schemaRef.Value != nil {
			if rxDeprecated.MatchString(schemaRef.Value.Description) {
				schemaRef.Value.Deprecated = true
			}
			for _, propRef := range schemaRef.Value.Properties {
				if len(propRef.Ref) == 0 && propRef.Value != nil {
					if rxDeprecated.MatchString(propRef.Value.Description) {
						propRef.Value.Deprecated = true
					}
				}
			}
		}
	}
}
