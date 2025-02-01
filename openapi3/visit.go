package openapi3

import (
	"fmt"
	"net/http"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/encoding/jsonpointer"
)

const (
	jPtrParamFormat = "#/components/parameters/%s"
	// jPtrSchemasRoot          = "#/components/schemas/"
	jPtrSchemaPropertyFormat = "#/components/schemas/%s/properties/%s"
)

func VisitTypesFormats(spec *Spec, visitTypeFormat func(jsonPointerRoot, oasType, oasFormat string)) {
	for schemaName, schemaRef := range spec.Components.Schemas {
		if schemaRef.Value == nil {
			continue
		}
		for propName, propRef := range schemaRef.Value.Properties {
			if propRef.Value == nil || propRef.Value.Type == nil {
				continue
			}
			for _, t := range *propRef.Value.Type {
				visitTypeFormat(
					fmt.Sprintf(jPtrSchemaPropertyFormat, schemaName, propName),
					t,
					propRef.Value.Format)
			}
		}
	}
	for paramName, paramRef := range spec.Components.Parameters {
		if paramRef.Value == nil ||
			paramRef.Value.Schema == nil ||
			paramRef.Value.Schema.Value == nil ||
			paramRef.Value.Schema.Value.Type == nil {
			continue
		}
		for _, t := range *paramRef.Value.Schema.Value.Type {
			visitTypeFormat(
				fmt.Sprintf(jPtrParamFormat, paramName),
				t,
				paramRef.Value.Schema.Value.Format)
		}
	}
	VisitOperations(
		spec,
		func(path, method string, op *oas3.Operation) {
			if op == nil {
				return
			}
			for i, paramRef := range op.Parameters {
				if paramRef.Value == nil ||
					paramRef.Value.Schema == nil ||
					paramRef.Value.Schema.Value == nil ||
					paramRef.Value.Schema.Value.Type == nil {
					continue
				}
				for _, t := range *paramRef.Value.Schema.Value.Type {
					visitTypeFormat(
						jsonpointer.PointerSubEscapeAll(
							"#/paths/%s/%s/parameters/%d/schema", path, strings.ToLower(method), i),
						t,
						paramRef.Value.Schema.Value.Format)
				}
			}
		},
	)
}

func VisitOperationsPathItem(path string, pathItem *oas3.PathItem, visitOp func(path, method string, op *oas3.Operation)) {
	pathURL := path
	if pathItem == nil {
		return
	}
	if pathItem.Connect != nil {
		visitOp(pathURL, http.MethodConnect, pathItem.Connect)
	}
	if pathItem.Delete != nil {
		visitOp(pathURL, http.MethodDelete, pathItem.Delete)
	}
	if pathItem.Get != nil {
		visitOp(pathURL, http.MethodGet, pathItem.Get)
	}
	if pathItem.Head != nil {
		visitOp(pathURL, http.MethodHead, pathItem.Head)
	}
	if pathItem.Options != nil {
		visitOp(pathURL, http.MethodOptions, pathItem.Options)
	}
	if pathItem.Patch != nil {
		visitOp(pathURL, http.MethodPatch, pathItem.Patch)
	}
	if pathItem.Post != nil {
		visitOp(pathURL, http.MethodPost, pathItem.Post)
	}
	if pathItem.Put != nil {
		visitOp(pathURL, http.MethodPut, pathItem.Put)
	}
	if pathItem.Trace != nil {
		visitOp(pathURL, http.MethodTrace, pathItem.Trace)
	}
}

func VisitOperations(spec *Spec, visitOp func(path, method string, op *oas3.Operation)) {
	pathsMap := spec.Paths.Map()
	for path, pathItem := range pathsMap {
		// for path, pathItem := range spec.Paths { // getkin v0.121.0 to v0.122.0
		VisitOperationsPathItem(path, pathItem, visitOp)
	}
}
