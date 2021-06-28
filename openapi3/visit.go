package openapi3

import (
	"fmt"
	"net/http"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/encoding/jsonutil"
)

const (
	jPtrParamFormat          = "#/components/parameters/%s"
	jPtrSchemasRoot          = "#/components/schemas/"
	jPtrSchemaPropertyFormat = "#/components/schemas/%s/properties/%s"
)

func VisitTypesFormats(spec *oas3.Swagger, visitTypeFormat func(jsonPointerRoot, oasType, oasFormat string)) {
	for schemaName, schemaRef := range spec.Components.Schemas {
		if schemaRef.Value == nil {
			continue
		}
		for propName, propRef := range schemaRef.Value.Properties {
			if propRef.Value == nil {
				continue
			}
			visitTypeFormat(
				fmt.Sprintf(jPtrSchemaPropertyFormat, schemaName, propName),
				propRef.Value.Type,
				propRef.Value.Format)
		}
	}
	for paramName, paramRef := range spec.Components.Parameters {
		if paramRef.Value == nil ||
			paramRef.Value.Schema == nil ||
			paramRef.Value.Schema.Value == nil {
			continue
		}
		visitTypeFormat(
			fmt.Sprintf(jPtrParamFormat, paramName),
			paramRef.Value.Schema.Value.Type,
			paramRef.Value.Schema.Value.Format)
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
					paramRef.Value.Schema.Value == nil {
					continue
				}
				visitTypeFormat(
					jsonutil.PointerSubEscapeAll(
						"#/paths/%s/%s/parameters/%d/schema", path, strings.ToLower(method), i),
					paramRef.Value.Schema.Value.Type,
					paramRef.Value.Schema.Value.Format)
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

func VisitOperations(spec *oas3.Swagger, visitOp func(path, method string, op *oas3.Operation)) {
	for path, pathItem := range spec.Paths {
		VisitOperationsPathItem(path, pathItem, visitOp)
		/*
			if pathItem.Connect != nil {
				visitOp(path, http.MethodConnect, pathItem.Connect)
			}
			if pathItem.Delete != nil {
				visitOp(path, http.MethodDelete, pathItem.Delete)
			}
			if pathItem.Get != nil {
				visitOp(path, http.MethodGet, pathItem.Get)
			}
			if pathItem.Head != nil {
				visitOp(path, http.MethodHead, pathItem.Head)
			}
			if pathItem.Options != nil {
				visitOp(path, http.MethodOptions, pathItem.Options)
			}
			if pathItem.Patch != nil {
				visitOp(path, http.MethodPatch, pathItem.Patch)
			}
			if pathItem.Post != nil {
				visitOp(path, http.MethodPost, pathItem.Post)
			}
			if pathItem.Put != nil {
				visitOp(path, http.MethodPut, pathItem.Put)
			}
			if pathItem.Trace != nil {
				visitOp(path, http.MethodTrace, pathItem.Trace)
			}
		*/
	}
}
