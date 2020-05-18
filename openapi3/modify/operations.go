package modify

import (
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

func SpecOperationsCount(spec *oas3.Swagger) uint {
	count := uint(0)
	VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		count++
	})
	return count
}

func SpecOperationIds(spec *oas3.Swagger) map[string]int {
	msi := map[string]int{}
	VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		op.OperationID = strings.TrimSpace(op.OperationID)
		if _, ok := msi[op.OperationID]; !ok {
			msi[op.OperationID] = 0
		}
		msi[op.OperationID]++
	})
	return msi
}

func SpecAddCustomProperties(spec *oas3.Swagger, custom map[string]interface{}, addToOperations, addToSchemas bool) {
	if len(custom) == 0 {
		return
	}
	if addToOperations {
		VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
			if op == nil {
				return
			}
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
