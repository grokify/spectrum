package openapi3

import (
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

func SpecGetComponentRequestBody(spec *oas3.Swagger, componentPath string) *oas3.RequestBodyRef {
	componentPathParts := strings.Split(strings.TrimSpace(componentPath), "/")
	if len(componentPathParts) != 4 ||
		componentPathParts[0] != "#" ||
		componentPathParts[1] != "components" ||
		componentPathParts[2] != "requestBodies" ||
		len(componentPathParts[3]) == 0 {
		return nil
	}
	if reqBodyRef, ok := spec.Components.RequestBodies[componentPathParts[3]]; ok {
		return reqBodyRef
	}
	return nil
}
