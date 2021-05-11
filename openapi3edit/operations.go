package openapi3edit

import (
	"fmt"
	"net/http"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/swaggman/openapi3"
)

func OperationAddExternalDocs(op *oas3.Operation, url, description string, preserveIfReqEmpty bool) {
	if op == nil {
		return
	}
	url = strings.TrimSpace(url)
	description = strings.TrimSpace(description)
	if len(url) > 0 || len(description) > 0 {
		if preserveIfReqEmpty {
			if op.ExternalDocs == nil {
				op.ExternalDocs = &oas3.ExternalDocs{}
			}
			if len(url) > 0 {
				op.ExternalDocs.URL = url
			}
			if len(description) > 0 {
				op.ExternalDocs.Description = description
			}
		} else {
			op.ExternalDocs = &oas3.ExternalDocs{
				Description: description,
				URL:         url}
		}
	}
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
			OperationAddExternalDocs(op, opMeta.DocsURL, opMeta.DocsDescription, true)
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
