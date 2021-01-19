package openapi3

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

func OperationToMeta(url, method string, op *oas3.Operation) OperationMeta {
	return OperationMeta{
		OperationID: strings.TrimSpace(op.OperationID),
		Summary:     strings.TrimSpace(op.Summary),
		Method:      strings.ToUpper(strings.TrimSpace(method)),
		Path:        strings.TrimSpace(url),
		Tags:        op.Tags,
		MetaNotes:   []string{}}
}

type OperationMeta struct {
	OperationID string
	Summary     string
	Method      string
	Path        string
	Tags        []string
	MetaNotes   []string
}

func OperationSetRequestBodySchemaRef(op *oas3.Operation, mediaType string, schemaRef *oas3.SchemaRef) {
	if op.RequestBody == nil {
		op.RequestBody = &oas3.RequestBodyRef{}
	}
	if op.RequestBody.Value == nil {
		op.RequestBody.Value = &oas3.RequestBody{
			Content: oas3.NewContent()}
	}
	op.RequestBody.Value.Content[mediaType] = oas3.NewMediaType().WithSchemaRef(schemaRef)
}

/*
	op.RequestBody = &oas3.RequestBodyRef{
		Value: &oas3.RequestBody{
			Content: map[string]*oas3.MediaType{
				"application/json": &oas3.MediaType{
					Schema: &oas3.SchemaRef{
						Ref: ref,
					},
					// Example: //
				},
			},
		},
	}

*/

func OperationSetResponseBodySchemaRef(op *oas3.Operation, status, description, mediaType string, schemaRef *oas3.SchemaRef) error {
	description = strings.TrimSpace(description)
	if len(description) == 0 {
		return fmt.Errorf("no response description for operationId [%s]", op.OperationID)
	}
	if op.Responses == nil {
		op.Responses = oas3.Responses{}
	}
	status = strings.TrimSpace(status)
	mediaType = strings.ToLower(strings.TrimSpace(mediaType))
	if _, ok := op.Responses[status]; !ok || op.Responses[status] == nil {
		op.Responses[status] = &oas3.ResponseRef{}
	}
	resRef, _ := op.Responses[status]
	if resRef.Value == nil {
		resRef.Value = &oas3.Response{}
	}
	resRef.Value.Description = &description
	if resRef.Value.Content == nil {
		resRef.Value.Content = oas3.NewContent()
	}
	resRef.Value.Content[mediaType] = oas3.NewMediaType().WithSchemaRef(schemaRef)
	return nil
}
