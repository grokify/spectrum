package modify

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
)

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
}
