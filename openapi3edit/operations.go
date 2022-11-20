package openapi3edit

import (
	"fmt"
	"net/http"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
)

// SpecSetOperation sets an operation in a OpenAPI Specification.
func SpecSetOperation(spec *openapi3.Spec, path, method string, op oas3.Operation) error {
	if spec == nil {
		return fmt.Errorf("spec to add operation to is nil for path[%s] method [%s]", path, method)
	}
	pathItem, ok := spec.Paths[path]
	if !ok {
		pathItem = &oas3.PathItem{}
	}
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case http.MethodConnect:
		pathItem.Connect = &op
	case http.MethodDelete:
		pathItem.Delete = &op
	case http.MethodGet:
		pathItem.Get = &op
	case http.MethodHead:
		pathItem.Head = &op
	case http.MethodOptions:
		pathItem.Options = &op
	case http.MethodPatch:
		pathItem.Patch = &op
	case http.MethodPost:
		pathItem.Post = &op
	case http.MethodPut:
		pathItem.Put = &op
	case http.MethodTrace:
		pathItem.Trace = &op
	default:
		return fmt.Errorf("spec operation method to set not found path[%s] method[%s]", path, method)
	}
	spec.Paths[path] = pathItem
	return nil
}

func SpecOperationIDsFromSummaries(spec *openapi3.Spec, errorOnEmpty bool) error {
	empty := []string{}
	openapi3.VisitOperations(spec, func(path, method string, op *oas3.Operation) {
		op.Summary = strings.Join(strings.Split(op.Summary, " "), " ")
		op.OperationID = op.Summary
		if len(op.OperationID) == 0 {
			empty = append(empty, openapi3.PathMethod(path, method))
		}
	})
	if errorOnEmpty && len(empty) > 0 {
		return fmt.Errorf("no_opid: [%s]", strings.Join(empty, ", "))
	}
	return nil
}

// SpecOperationsOperationIDSummaryReplace sets the OperationID and Summary with a `map[string]string`
// where the keys are pathMethod values and the values are Summary strings.
// This currently converts a Summary into an OperationID by using the supplied `opIDFunc`.
func SpecOperationsOperationIDSummaryReplace(spec *openapi3.Spec, customMapPathMethodToSummary map[string]string, opIDFunc func(s string) string, forceOpID, forceSummary bool) {
	openapi3.VisitOperations(spec, func(path, method string, op *oas3.Operation) {
		op.OperationID = strings.TrimSpace(op.OperationID)
		op.Summary = strings.TrimSpace(op.Summary)
		pathMethod := openapi3.PathMethod(path, method)
		summaryTry, ok := customMapPathMethodToSummary[pathMethod]
		if !ok {
			return
		}
		opIDTry := summaryTry
		if opIDFunc != nil {
			opIDTry = opIDFunc(summaryTry)
		}
		if len(op.OperationID) == 0 || forceOpID {
			op.OperationID = opIDTry
		}
		if len(op.Summary) == 0 || forceSummary {
			op.Summary = summaryTry
		}
	})
}

func SpecAddCustomProperties(spec *openapi3.Spec, custom map[string]interface{}, addToOperations, addToSchemas bool) {
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

func SpecAddOperationMetas(spec *openapi3.Spec, metas map[string]openapi3.OperationMeta, overwrite bool) {
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
			err := operationAddExternalDocs(op, opMeta.DocsURL, opMeta.DocsDescription, true)
			if err != nil {
				return
			}
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

// SpecOperationsSecurityReplace rplaces the security requirement object of operations that meets its
// include and exclude filters. SecurityRequirement is specified by OpenAPI/Swagger standard version 3.
// See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#securityRequirementObject
func SpecOperationsSecurityReplace(spec *openapi3.Spec, pathMethodsInclude, pathMethodsExclude []string, securityRequirement map[string][]string) {
	pathMethodsExcludeMap := stringsutil.SliceToMap(stringsutil.SliceCondenseSpace(pathMethodsExclude, true, false))
	pathMethodsIncludeMap := stringsutil.SliceToMap(stringsutil.SliceCondenseSpace(pathMethodsInclude, true, false))

	openapi3.VisitOperations(spec, func(opPath, opMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		pathMethod := openapi3.PathMethod(opPath, opMethod)
		if _, ok := pathMethodsExcludeMap[pathMethod]; ok {
			return
		}
		if len(pathMethodsIncludeMap) > 0 { // only filter on explicit includes is more than one include.
			if _, ok := pathMethodsIncludeMap[pathMethod]; !ok {
				return
			}
		}
		op.Security = oas3.NewSecurityRequirements()
		op.Security.With(securityRequirement)
	})
}
