package openapi3edit

import (
	"fmt"
	"net/http"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/swaggman/openapi3"
)

type OperationEditor struct {
	Operation *oas3.Operation
}

func (oedit *OperationEditor) AddExternalDocs(docURL, docDescription string, preserveIfReqEmpty bool) {
	operationAddExternalDocs(oedit.Operation, docURL, docDescription, preserveIfReqEmpty)
}

func (oedit *OperationEditor) AddRequestBodySchemaRef(description string, required bool, contentType string, schemaRef *oas3.SchemaRef) error {
	return operationAddRequestBodySchemaRef(oedit.Operation, description, required, contentType, schemaRef)
}

func (oedit *OperationEditor) AddResponseBodySchemaRef(statusCode, description, contentType string, schemaRef *oas3.SchemaRef) error {
	return operationAddResponseBodySchemaRef(oedit.Operation, statusCode, description, contentType, schemaRef)
}

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

func operationAddResponseBodySchemaRef(op *oas3.Operation, statusCode, description, contentType string, schemaRef *oas3.SchemaRef) error {
	if op == nil {
		return fmt.Errorf("operation to edit is nil")
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
		op.Responses[statusCode].Value = &oas3.Response{
			Description: &description,
		}
	}
	if op.Responses[statusCode].Value.Content == nil {
		op.Responses[statusCode].Value.Content = oas3.NewContent()
	}
	op.Responses[statusCode].Value.Content[contentType] = oas3.NewMediaType().WithSchemaRef(schemaRef)
	return nil
}

func operationAddExternalDocs(op *oas3.Operation, docURL, docDescription string, preserveIfReqEmpty bool) error {
	if op == nil {
		return fmt.Errorf("operation to edit is nil")
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

func SpecOperationsCount(spec *oas3.Swagger) uint {
	count := uint(0)
	openapi3.VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
		count++
	})
	return count
}

func SpecSetOperation(spec *oas3.Swagger, path, method string, op oas3.Operation) {
	pathItem, ok := spec.Paths[path]
	if !ok {
		pathItem = &oas3.PathItem{}
	}
	method = strings.ToUpper(strings.TrimSpace(method))
	switch method {
	case http.MethodGet:
		pathItem.Get = &op
	case http.MethodPost:
		pathItem.Post = &op
	case http.MethodPut:
		pathItem.Put = &op
	case http.MethodPatch:
		pathItem.Patch = &op
	}

}

func SpecOperationIds(spec *oas3.Swagger) map[string]int {
	msi := map[string]int{}
	openapi3.VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
		op.OperationID = strings.TrimSpace(op.OperationID)
		if _, ok := msi[op.OperationID]; !ok {
			msi[op.OperationID] = 0
		}
		msi[op.OperationID]++
	})
	return msi
}

func SpecOperationIdsFromSummaries(spec *oas3.Swagger, errorOnEmpty bool) error {
	empty := []string{}
	openapi3.VisitOperations(spec, func(path, method string, op *oas3.Operation) {
		op.Summary = strings.Join(strings.Split(op.Summary, " "), " ")
		op.OperationID = op.Summary
		if len(op.OperationID) == 0 {
			empty = append(empty, path+" "+method)
		}
	})
	if errorOnEmpty && len(empty) > 0 {
		return fmt.Errorf("no_opid: [%s]", strings.Join(empty, ", "))
	}
	return nil
}

func SpecAddCustomProperties(spec *oas3.Swagger, custom map[string]interface{}, addToOperations, addToSchemas bool) {
	if spec == nil || len(custom) == 0 {
		return
	}
	if addToOperations {
		openapi3.VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
			for key, val := range custom {
				op.Extensions[key] = val
			}
		})
	}
	if addToSchemas {
		for _, schema := range spec.Components.Schemas {
			if schema.Value != nil {
				for key, val := range custom {
					schema.Value.Extensions[key] = val
				}
			}
		}
	}
}

func SpecAddOperationMetas(spec *oas3.Swagger, metas map[string]openapi3.OperationMeta, overwrite bool) {
	if spec == nil || len(metas) == 0 {
		return
	}
	openapi3.VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		opMeta, ok := metas[op.OperationID]
		if !ok {
			return
		}
		opMeta.TrimSpace()
		writeDocs := false
		writeScopes := false
		writeThrottling := false
		if overwrite {
			writeDocs = true
			writeScopes = true
			writeThrottling = true
		}
		if writeDocs {
			operationAddExternalDocs(op, opMeta.DocsURL, opMeta.DocsDescription, true)
		}
		if writeScopes {
			if len(opMeta.SecurityScopes) > 0 {
				op.Security = &oas3.SecurityRequirements{
					map[string][]string{"oauth": opMeta.SecurityScopes},
				}
			} else {
				op.Security = nil
			}
		}
		if writeThrottling {
			if op.ExtensionProps.Extensions == nil {
				op.ExtensionProps.Extensions = map[string]interface{}{}
			}
			op.ExtensionProps.Extensions[openapi3.XThrottlingGroup] = opMeta.XThrottlingGroup
		}
	})
}

type OperationMoreSet struct {
	OperationMores []OperationMore
}

type OperationMore struct {
	UrlPath   string
	Method    string
	Operation *oas3.Operation
}

func QueryOperationsByTags(spec *oas3.Swagger, tags []string) *OperationMoreSet {
	tagsWantMatch := map[string]int{}
	for _, tag := range tags {
		tagsWantMatch[tag] = 1
	}
	opmSet := &OperationMoreSet{OperationMores: []OperationMore{}}
	// for path, pathInfo := range spec.Paths {
	openapi3.VisitOperations(spec, func(url, method string, op *oas3.Operation) {
		if op == nil {
			return
		}
		for _, tagTry := range op.Tags {
			if _, ok := tagsWantMatch[tagTry]; ok {
				opmSet.OperationMores = append(opmSet.OperationMores,
					OperationMore{
						UrlPath:   url,
						Method:    method,
						Operation: op})
				return
			}
		}
	})
	// }
	return opmSet
}
