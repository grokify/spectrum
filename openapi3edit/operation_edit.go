package openapi3edit

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/spectrum/openapi3"
)

// OperationEdit is used for two purposes: (a) to store path and method information with the operation and
// (b) to provide a container to organize operation related functions.
type OperationEdit struct {
	openapi3.OperationMore
}

func (ope *OperationEdit) SetExternalDocs(docURL, docDescription string, preserveIfReqEmpty bool) error {
	return operationAddExternalDocs(ope.OperationMore.Operation, docURL, docDescription, preserveIfReqEmpty)
}

func (ope *OperationEdit) SetRequestBodyAttrs(description string, required bool) error {
	if ope.OperationMore.Operation == nil {
		return openapi3.ErrOperationNotSet
	}
	op := ope.OperationMore.Operation
	if op.RequestBody == nil {
		op.RequestBody = &oas3.RequestBodyRef{}
	}
	if op.RequestBody.Value == nil {
		op.RequestBody.Value = &oas3.RequestBody{
			Content: oas3.NewContent()}
	}
	op.RequestBody.Value.Description = strings.TrimSpace(description)
	op.RequestBody.Value.Required = required
	return nil
}

func (ope *OperationEdit) SetRequestBodySchemaRef(mediaType string, schemaRef *oas3.SchemaRef) error {
	if ope.OperationMore.Operation == nil {
		return openapi3.ErrOperationNotSet
	}
	op := ope.OperationMore.Operation
	mediaType = strings.TrimSpace(mediaType)
	if op.RequestBody == nil {
		op.RequestBody = &oas3.RequestBodyRef{}
	}
	if op.RequestBody.Value == nil {
		op.RequestBody.Value = &oas3.RequestBody{
			Content: oas3.NewContent()}
	}
	op.RequestBody.Value.Content[mediaType] = oas3.NewMediaType().WithSchemaRef(schemaRef)
	return nil
}

/*
func (ope *OperationEdit) SetRequestBodySchemaRefMore(description string, required bool, contentType string, schemaRef *oas3.SchemaRef) error {
	return operationAddRequestBodySchemaRef(ope.OperationMore.Operation, description, required, contentType, schemaRef)
}
*/

func (ope *OperationEdit) SetResponseBodySchemaRefMore(statusCode, description, contentType string, schemaRef *oas3.SchemaRef) error {
	return operationAddResponseBodySchemaRef(ope.OperationMore.Operation, statusCode, description, contentType, schemaRef)
}

func (ope *OperationEdit) AddToSpec(spec *openapi3.Spec, force bool) (bool, error) {
	sm := openapi3.SpecMore{Spec: spec}
	op, err := sm.OperationByPathMethod(ope.OperationMore.Path, ope.OperationMore.Method)
	if err != nil {
		return false, err
	}
	if op == nil || force {
		spec.AddOperation(ope.OperationMore.Path, ope.OperationMore.Method, ope.OperationMore.Operation)
		return true, nil
	}
	return false, nil
}

/*
func operationAddRequestBodySchemaRef(op *oas3.Operation, description string, required bool, contentType string, schemaRef *oas3.SchemaRef) error {
	if op == nil {
		return fmt.Errorf("operation to edit is nil")
	}
	if op.RequestBody == nil {
		op.RequestBody = &oas3.RequestBodyRef{}
	}
	description = strings.TrimSpace(description)
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	if len(contentType) == 0 {
		return fmt.Errorf("content type [%s] is empty", contentType)
	}
	if len(op.RequestBody.Ref) > 0 {
		return fmt.Errorf("request body is reference for operationId [%s]", op.OperationID)
	}
	if op.RequestBody.Value == nil {
		op.RequestBody.Value = &oas3.RequestBody{}
	}
	op.RequestBody.Value.Description = description
	op.RequestBody.Value.Required = required
	if op.RequestBody.Value.Content == nil {
		op.RequestBody.Value.Content = oas3.NewContent()
	}
	op.RequestBody.Value.Content[contentType] = oas3.NewMediaType().WithSchemaRef(schemaRef)
	return nil
}
*/

func operationAddResponseBodySchemaRef(op *oas3.Operation, statusCode, description, contentType string, schemaRef *oas3.SchemaRef) error {
	if op == nil {
		return openapi3.ErrOperationNotSet
	}
	if schemaRef == nil {
		return fmt.Errorf("operation response to body to add is nil")
	}
	statusCode = strings.TrimSpace(statusCode)
	description = strings.TrimSpace(description)
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	if statusCode == "" || contentType == "" {
		return fmt.Errorf("status code [%s] or content type [%s] is empty", statusCode, contentType)
	}
	if op.Responses == nil {
		op.Responses = oas3.Responses{}
	}
	if op.Responses[statusCode] == nil {
		op.Responses[statusCode] = &oas3.ResponseRef{}
	}
	if len(op.Responses[statusCode].Ref) > 0 {
		return fmt.Errorf("response is a reference and not actual")
	}
	if op.Responses[statusCode].Value == nil {
		op.Responses[statusCode].Value = &oas3.Response{}
	}
	op.Responses[statusCode].Value.Description = &description
	if op.Responses[statusCode].Value.Content == nil {
		op.Responses[statusCode].Value.Content = oas3.NewContent()
	}
	op.Responses[statusCode].Value.Content[contentType] = oas3.NewMediaType().WithSchemaRef(schemaRef)
	return nil
}

/*
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
	resRef := op.Responses[status]
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
*/

func operationAddExternalDocs(op *oas3.Operation, docURL, docDescription string, preserveIfReqEmpty bool) error {
	if op == nil {
		return openapi3.ErrOperationNotSet
	}
	docURL = strings.TrimSpace(docURL)
	docDescription = strings.TrimSpace(docDescription)
	if len(docURL) > 0 || len(docDescription) > 0 {
		if preserveIfReqEmpty {
			if op.ExternalDocs == nil {
				op.ExternalDocs = &oas3.ExternalDocs{}
			}
			if len(docURL) > 0 {
				op.ExternalDocs.URL = docURL
			}
			if len(docDescription) > 0 {
				op.ExternalDocs.Description = docDescription
			}
		} else {
			op.ExternalDocs = &oas3.ExternalDocs{
				Description: docDescription,
				URL:         docURL}
		}
	}
	return nil
}

/*
type OperationEditSet struct {
	OperationEdits []OperationEdit
}

// SummariesMap returns a `map[string]string` where the keys are the operation's
// path and method, while the values are the sumamries.`
func (omSet *OperationEditSet) SummariesMap() map[string]string {
	mss := map[string]string{}
	for _, om := range omSet.OperationEdits {
		mss[om.PathMethod()] = om.Operation.Summary
	}
	return mss
}

func QueryOperationsByTags(spec *openapi3.Spec, tags []string) *OperationEditSet {
	tagsWantMatch := map[string]int{}
	for _, tag := range tags {
		tagsWantMatch[tag] = 1
	}
	opmSet := &OperationEditSet{OperationEdits: []OperationEdit{}}

	openapi3.VisitOperations(spec, func(path, method string, op *oas3.Operation) {
		if op == nil {
			return
		}
		for _, tagTry := range op.Tags {
			if _, ok := tagsWantMatch[tagTry]; ok {
				opmSet.OperationEdits = append(opmSet.OperationEdits,
					OperationEdit{
						OperationMore: openapi3.OperationMore{
							Path:      path,
							Method:    method,
							Operation: op}})
				return
			}
		}
	})

	return opmSet
}
*/
