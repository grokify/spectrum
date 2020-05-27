package openapi3

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

// ValidateResponseTypes looks for `application/json` responses
// with response schema types that are not `array` or `object`.
func ValidateResponseTypes(spec *oas3.Swagger) ([]*OperationMeta, error) {
	errorOperations := []*OperationMeta{}
	VisitOperations(
		spec,
		func(path, method string, op *oas3.Operation) {
			if op == nil {
				return
			}
			for _, resRef := range op.Responses {
				if resRef == nil {
					continue
				}
				if resRef.Value == nil {
					continue
				}
				response := resRef
				for mediaType, mtRef := range response.Value.Content {
					mediaType = strings.ToLower(strings.TrimSpace(mediaType))
					if mediaType == "application/json" {
						schemaRef := mtRef.Schema
						if len(schemaRef.Ref) == 0 {
							schema := schemaRef.Value
							schemaType := schema.Type
							if schemaType != "object" && schemaType != "array" {
								om := OperationToMeta(path, method, op)
								om.MetaNotes = append(om.MetaNotes,
									fmt.Sprintf("E_BAD_MIME_TYPE_AND_SCHEMA MT[%s] type[%s]", mediaType, schemaType))
								errorOperations = append(errorOperations, &om)
							}
						}
					}
				}
			}
		},
	)

	if len(errorOperations) > 0 {
		return errorOperations, fmt.Errorf("E_NUM_VALIDATION_ERRORS [%v]", len(errorOperations))
	}
	return errorOperations, nil
}
