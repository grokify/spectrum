package modify

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

func SpecOperationsCount(spec *oas3.Swagger) uint {
	count := uint(0)
	VisitOperations(spec, func(op *oas3.Operation) {
		if op == nil {
			return
		}
		count++
	})
	return count
}

func SpecOperationIds(spec *oas3.Swagger) map[string]int {
	msi := map[string]int{}
	VisitOperations(spec, func(op *oas3.Operation) {
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

func UpdateOperationIds(spec *oas3.Swagger, renameOpId func(orig string) string) {
	VisitOperations(spec, func(op *oas3.Operation) {
		if op == nil {
			return
		}
		op.OperationID = strings.TrimSpace(renameOpId(op.OperationID))
	})
	opIds := map[string]int{}
	VisitOperations(spec, func(op *oas3.Operation) {
		if op == nil {
			return
		}
		opIds[strings.TrimSpace(op.OperationID)] = 1
	})
	for opId, count := range opIds {
		if count > 1 {
			panic(fmt.Sprintf("OPID ID[%s] COUNT[%d]", opId, count))
		}
	}
}

func SpecAddCustomProperties(spec *oas3.Swagger, custom map[string]interface{}, addToOperations, addToSchemas bool) {
	if len(custom) == 0 {
		return
	}
	if addToOperations {
		VisitOperations(spec, func(op *oas3.Operation) {
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
