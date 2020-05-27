package modify

import (
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/swaggman/openapi3"
)

func SpecOperationsCount(spec *oas3.Swagger) uint {
	count := uint(0)
	openapi3.VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
		count++
	})
	return count
}

func SpecOperationIds(spec *oas3.Swagger) map[string]int {
	msi := map[string]int{}
	openapi3.VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
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

type OperationMoreSet struct {
	OperationMores []OperationMore
}

type OperationMore struct {
	UrlPath   string
	Method    string
	Operation *oas3.Operation
}

func QueryOperationsByTags(spec *oas3.Swagger, tags []string) *OperationMoreSet {
	tagsWantMatch := map[string]int{}
	for _, tag := range tags {
		tagsWantMatch[tag] = 1
	}
	opmSet := &OperationMoreSet{OperationMores: []OperationMore{}}
	// for path, pathInfo := range spec.Paths {
	openapi3.VisitOperations(spec, func(url, method string, op *oas3.Operation) {
		if op == nil {
			return
		}
		for _, tagTry := range op.Tags {
			if _, ok := tagsWantMatch[tagTry]; ok {
				opmSet.OperationMores = append(opmSet.OperationMores,
					OperationMore{
						UrlPath:   url,
						Method:    method,
						Operation: op})
				return
			}
		}
	})
	// }
	return opmSet
}
