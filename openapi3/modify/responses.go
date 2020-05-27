package modify

import (
	"fmt"
	"reflect"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/swaggman/openapi3"
)

// ValidateResponseTypes looks for `application/json` responses
// with response schema types that are not `array` or `object`.
func ValidateResponseTypes(spec *oas3.Swagger, fix bool) ([]*openapi3.OperationMeta, error) {
	errorOperations := []*openapi3.OperationMeta{}
	openapi3.VisitOperations(
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
				for mediaTypeOrig, mtRef := range response.Value.Content {
					mediaType := strings.ToLower(strings.TrimSpace(mediaTypeOrig))
					if mediaType == httputilmore.ContentTypeAppJson {
						schemaRef := mtRef.Schema
						if len(schemaRef.Ref) == 0 {
							schema := schemaRef.Value
							schemaType := schema.Type
							if fix && (schemaType == "string" || schemaType == "integer") {
								delete(response.Value.Content, mediaTypeOrig)
								if mtRefTry, ok := response.Value.Content[httputilmore.ContentTypeTextPlain]; ok {
									if !reflect.DeepEqual(mtRef, mtRefTry) {
										om := openapi3.OperationToMeta(path, method, op)
										om.MetaNotes = append(om.MetaNotes,
											fmt.Sprintf("E_BAD_MIME_TYPE_AND_SCHEMA_COLLISION MT[%s] type[%s]", mediaType, schemaType))
										errorOperations = append(errorOperations, &om)
									}
								} else {
									response.Value.Content[httputilmore.ContentTypeTextPlain] = mtRef
								}
							} else if schemaType != "object" && schemaType != "array" {
								om := openapi3.OperationToMeta(path, method, op)
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
