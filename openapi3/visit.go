package openapi3

import (
	"net/http"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

func VisitOperations(spec *oas3.Swagger, visitOp func(path, method string, op *oas3.Operation)) {
	for path, pathItem := range spec.Paths {
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
	}
}
