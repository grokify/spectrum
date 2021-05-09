package openapi3lint1

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/encoding/jsonutil"
)

func SpecCheckPathItems(spec *oas3.Swagger, rules Policy) PolicyViolationsSets {
	vsets := NewPolicyViolationsSets()
	if !rules.HasPathItemRules() {
		return vsets
	}

	for path, pathItemRef := range spec.Paths {
		if pathItemRef == nil {
			continue
		}
		jsPointer := "#/paths/" + jsonutil.PropertyNameEscape(path)
		err := vsets.UpsertSets(
			ParametersCheck(pathItemRef.Parameters, jsPointer, rules))
		if err != nil {
			vsets.AddSimple(RuleInternalError, LocationPaths, err.Error())
		}
	}

	return vsets
}
