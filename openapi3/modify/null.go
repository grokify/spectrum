package modify

import oas3 "github.com/getkin/kin-openapi/openapi3"

func NullToEmpty(spec *oas3.Swagger) {
	NullToEmptyPaths(spec)
}

// NullToEmptyPaths converts a `path` property from `null` to
// an empty set `{}` to satisfy OpenAPI Generator which will
// fail on the following error "-attribute paths is not of type `object`"
func NullToEmptyPaths(spec *oas3.Swagger) {
	if spec.Paths == nil {
		spec.Paths = map[string]*oas3.PathItem{}
	}
}
