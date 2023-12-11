package openapi3edit

import oas3 "github.com/getkin/kin-openapi/openapi3"

// PathsNullToEmpty converts a `path` property from `null` to
// an empty set `{}` to satisfy OpenAPI Generator which will
// fail on the following error "-attribute paths is not of type `object`"
func (se *SpecEdit) PathsNullToEmpty() {
	if se.SpecMore.Spec != nil && se.SpecMore.Spec.Paths == nil {
		se.SpecMore.Spec.Paths = oas3.NewPaths()
		// se.SpecMore.Spec.Paths = map[string]*oas3.PathItem{} // getkin v0.121.0 to v0.122.0
	}
}
