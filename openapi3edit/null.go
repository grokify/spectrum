package openapi3edit

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/spectrum/openapi3"
)

func NullToEmpty(spec *openapi3.Spec) {
	NullToEmptyPaths(spec)
}

// NullToEmptyPaths converts a `path` property from `null` to
// an empty set `{}` to satisfy OpenAPI Generator which will
// fail on the following error "-attribute paths is not of type `object`"
func NullToEmptyPaths(spec *openapi3.Spec) {
	if spec.Paths == nil {
		spec.Paths = map[string]*oas3.PathItem{}
	}
}
