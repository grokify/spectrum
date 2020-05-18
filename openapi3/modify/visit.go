package modify

import (
	"net/http"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

/*
func VisitOperations(spec *oas3.Swagger, visitOp func(op *oas3.Operation)) {
	for _, pathItem := range spec.Paths {
		visitOp(pathItem.Connect)
		visitOp(pathItem.Delete)
		visitOp(pathItem.Get)
		visitOp(pathItem.Head)
		visitOp(pathItem.Options)
		visitOp(pathItem.Patch)
		visitOp(pathItem.Post)
		visitOp(pathItem.Put)
		visitOp(pathItem.Trace)
	}
}*/

func VisitOperations(spec *oas3.Swagger, visitOp func(path, method string, op *oas3.Operation)) {
	for path, pathItem := range spec.Paths {
		visitOp(path, http.MethodConnect, pathItem.Connect)
		visitOp(path, http.MethodDelete, pathItem.Delete)
		visitOp(path, http.MethodGet, pathItem.Get)
		visitOp(path, http.MethodHead, pathItem.Head)
		visitOp(path, http.MethodOptions, pathItem.Options)
		visitOp(path, http.MethodPatch, pathItem.Patch)
		visitOp(path, http.MethodPost, pathItem.Post)
		visitOp(path, http.MethodPut, pathItem.Put)
		visitOp(path, http.MethodTrace, pathItem.Trace)
	}
}
