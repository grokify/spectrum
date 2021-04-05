package stylechecker

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/encoding/jsonutil"
)

func SpecCheckPathItems(spec *oas3.Swagger, rules RuleSet) PolicyViolationsSets {
	vsets := NewPolicyViolationsSets()
	if !rules.HasPathItemRules() {
		return vsets
	}

	for path, pathItemRef := range spec.Paths {
		if pathItemRef == nil {
			continue
		}
		jsPointer := "#/paths/" + jsonutil.PropertyNameEscape(path)
		vsets.UpsertSets(
			ParametersCheck(pathItemRef.Parameters, jsPointer, rules))
	}

	return vsets
}
