package modify

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
)

// RemoveOperationsSecurity removes the security property
// for all operations. It is useful when building a spec
// to get individual specs to validate before setting the
// correct security property.
func RemoveOperationsSecurity(spec *oas3.Swagger) {
	VisitOperations(spec, func(op *oas3.Operation) {
		if op == nil {
			return
		}
		op.Security = &oas3.SecurityRequirements{}
	})
}
