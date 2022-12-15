package openapi3edit

import (
	"regexp"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/spectrum/openapi3"
)

func (se *SpecEdit) SchemasSetDeprecated(newDeprecated bool) {
	if se.SpecMore.Spec == nil {
		return
	}
	for _, schemaRef := range se.SpecMore.Spec.Components.Schemas {
		if len(schemaRef.Ref) == 0 && schemaRef.Value != nil {
			schemaRef.Value.Deprecated = newDeprecated
		}
	}
}

func (se *SpecEdit) OperationsSetDeprecated(newDeprecated bool) {
	if se.SpecMore.Spec == nil {
		return
	}
	openapi3.VisitOperations(
		se.SpecMore.Spec,
		func(path, method string, op *oas3.Operation) {
			if op != nil {
				op.Deprecated = newDeprecated
			}
		},
	)
}

var rxDeprecated = regexp.MustCompile(`(?i)\bdeprecated\b`)

func (se *SpecEdit) SetDeprecatedImplicit() {
	if se.SpecMore.Spec == nil {
		return
	}
	openapi3.VisitOperations(
		se.SpecMore.Spec,
		func(path, method string, op *oas3.Operation) {
			if op != nil && rxDeprecated.MatchString(op.Description) {
				op.Deprecated = true
			}
		},
	)
	for _, schemaRef := range se.SpecMore.Spec.Components.Schemas {
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
