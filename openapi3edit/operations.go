package openapi3edit

import (
	"fmt"
	"net/url"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/encoding/jsonpointer"
	"github.com/grokify/mogo/net/http/pathmethod"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
)

// SetOperation sets an operation in a OpenAPI Specification.
func (se *SpecEdit) SetOperation(path, method string, op oas3.Operation) error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	}
	spec := se.SpecMore.Spec
	if spec == nil {
		return fmt.Errorf("spec to add operation to is nil for path[%s] method [%s]", path, method)
	}
	/*
		pathItem, ok := spec.Paths[path]
		if !ok {
			pathItem = &oas3.PathItem{}
		}
	*/
	pathItem := spec.Paths.Find(path) // getkin v0.121.0 to v0.122.0
	if pathItem == nil {
		pathItem = &oas3.PathItem{}
	}
	method = strings.ToUpper(strings.TrimSpace(method))
	pathItem.SetOperation(method, &op)
	/*
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
	*/
	spec.Paths.Set(path, pathItem)
	return nil
}

func (se *SpecEdit) OperationIDsFromSummaries(errorOnEmpty bool) error {
	if se.SpecMore.Spec == nil {
		return nil
	}
	spec := se.SpecMore.Spec
	empty := []string{}
	openapi3.VisitOperations(spec, func(path, method string, op *oas3.Operation) {
		op.Summary = strings.Join(strings.Split(op.Summary, " "), " ")
		op.OperationID = op.Summary
		if len(op.OperationID) == 0 {
			empty = append(empty, pathmethod.PathMethod(path, method))
		}
	})
	if errorOnEmpty && len(empty) > 0 {
		return fmt.Errorf("no_opid: [%s]", strings.Join(empty, ", "))
	}
	return nil
}

// OperationsOperationIDSummaryReplace sets the OperationID and Summary with a `map[string]string`
// where the keys are pathMethod values and the values are Summary strings.
// This currently converts a Summary into an OperationID by using the supplied `opIDFunc`.
func (se *SpecEdit) OperationsOperationIDSummaryReplace(customMapPathMethodToSummary map[string]string, opIDFunc func(s string) string, forceOpID, forceSummary bool) {
	if se.SpecMore.Spec == nil {
		return
	}
	spec := se.SpecMore.Spec
	openapi3.VisitOperations(spec, func(path, method string, op *oas3.Operation) {
		op.OperationID = strings.TrimSpace(op.OperationID)
		op.Summary = strings.TrimSpace(op.Summary)
		pathMethod := pathmethod.PathMethod(path, method)
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

func (se *SpecEdit) AddCustomProperties(custom map[string]interface{}, addToOperations, addToSchemas bool) {
	if se.SpecMore.Spec == nil || len(custom) == 0 {
		return
	}
	spec := se.SpecMore.Spec
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

func (se *SpecEdit) AddOperationMetas(metas map[string]openapi3.OperationMeta, overwrite bool) {
	if se.SpecMore.Spec == nil || len(metas) == 0 {
		return
	}
	spec := se.SpecMore.Spec
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
			ope := NewOperationEdit(skipPath, skipMethod, op)
			err := ope.SetExternalDocs(opMeta.DocsURL, opMeta.DocsDescription, true)
			// err := operationAddExternalDocs(op, opMeta.DocsURL, opMeta.DocsDescription, true)
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
			/*
				if op.ExtensionProps.Extensions == nil {
					op.ExtensionProps.Extensions = map[string]interface{}{}
				}
				op.ExtensionProps.Extensions[openapi3.XThrottlingGroup] = opMeta.XThrottlingGroup
			*/
			if op.Extensions == nil {
				op.Extensions = map[string]any{}
			}
			op.Extensions[openapi3.XThrottlingGroup] = opMeta.XThrottlingGroup
		}
	})
}

// OperationsSecurityReplace rplaces the security requirement object of operations that meets its
// include and exclude filters. SecurityRequirement is specified by OpenAPI/Swagger standard version 3.
// See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#securityRequirementObject
func (se *SpecEdit) OperationsSecurityReplace(pathMethodsInclude, pathMethodsExclude []string, securityRequirement map[string][]string) {
	if se.SpecMore.Spec == nil {
		return
	}
	spec := se.SpecMore.Spec
	pathMethodsExcludeMap := stringsutil.SliceToMap(stringsutil.SliceCondenseSpace(pathMethodsExclude, true, false))
	pathMethodsIncludeMap := stringsutil.SliceToMap(stringsutil.SliceCondenseSpace(pathMethodsInclude, true, false))

	openapi3.VisitOperations(spec, func(opPath, opMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		pathMethod := pathmethod.PathMethod(opPath, opMethod)
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

const (
	ErrMsgSchemaKeyCollision    = "schemaKeySollision"
	ErrMsgOperationMissingID    = "operationMissingID"
	ErrMsgOperationRequestCTGT1 = "operationRequestContentTypeGT1"
)

func (se *SpecEdit) OperationsRequestBodyFlattenSchemas(schemaKeySuffix string) error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	}
	errs := url.Values{}
	openapi3.VisitOperations(se.SpecMore.Spec, func(opPath, opMethod string, op *oas3.Operation) {
		if op == nil || op.RequestBody == nil ||
			len(strings.TrimSpace(op.RequestBody.Ref)) > 0 || // have ref
			op.RequestBody.Value == nil ||
			len(op.RequestBody.Value.Content) == 0 { // have no content
			return
		}
		pm := pathmethod.PathMethod(opPath, opMethod)
		if len(op.RequestBody.Value.Content) > 1 {
			errs.Add(ErrMsgOperationRequestCTGT1, pm)
			return
		}
		if len(strings.TrimSpace(op.OperationID)) == 0 {
			errs.Add(ErrMsgOperationMissingID, pm)
			return
		}
		for _, mt := range op.RequestBody.Value.Content {
			if mt.Schema == nil ||
				len(strings.TrimSpace(mt.Schema.Ref)) > 0 ||
				mt.Schema.Value == nil {
				return
			}
			schKey := strings.TrimSpace(op.OperationID) + schemaKeySuffix
			if _, ok := se.SpecMore.Spec.Components.Schemas[schKey]; ok {
				errs.Add(ErrMsgSchemaKeyCollision, pm)
				return
			}
			schRef := mt.Schema
			se.SpecMore.Spec.Components.Schemas[schKey] = schRef
			schKeyPointer := jsonpointer.PointerSubEscapeAll(openapi3.PointerComponentsSchemasFormat, schKey)
			mt.Schema = oas3.NewSchemaRef(schKeyPointer, nil)
		}
	})

	if len(errs) > 0 {
		enc := errs.Encode()
		return fmt.Errorf("operation flattening faled: (%s)", enc)
	}
	return nil
}
