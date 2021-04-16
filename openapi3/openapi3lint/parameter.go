package openapi3lint

import (
	"strconv"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/text/stringcase"
)

func ParametersCheck(params oas3.Parameters, jsPointerParameters string, rules Policy) PolicyViolationsSets {
	vsets := NewPolicyViolationsSets()

	for i, paramRef := range params {
		if paramRef.Value != nil {
			paramName := paramRef.Value.Name
			if strings.ToLower(paramRef.Value.In) == "path" {
				jsPointerPath := PointerCondense(jsPointerParameters +
					"/parameters/" + strconv.Itoa(i) + "/name")
				if len(paramName) == 0 {
					vsets.AddSimple(
						RulePathParamNameExist,
						jsPointerPath, paramName)
				} else {
					if rules.HasRule(RulePathParamStyleCamelCase) &&
						!stringcase.IsCamelCase(paramName) {
						vsets.AddSimple(
							RulePathParamStyleCamelCase,
							jsPointerPath, paramName)
					}
				}
			}
		}
	}
	return vsets
}
