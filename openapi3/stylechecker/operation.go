package stylechecker

import (
	"strconv"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/text/stringcase"
	"github.com/grokify/swaggman/openapi3"
)

func SpecCheckOperations(spec *oas3.Swagger, ruleset RuleSet) PolicyViolationsSets {
	vsets := NewPolicyViolationsSets()

	openapi3.VisitOperations(spec,
		func(path, method string, op *oas3.Operation) {
			if op == nil {
				return
			}
			//opLoc := "#/paths/" + jsonutil.PropertyNameEscape(path) + "/" + strings.ToLower(method)
			opLoc := PointerSubEscapeAll("#/paths/%s/%s", path, strings.ToLower(method))
			opId := strings.TrimSpace(op.OperationID)
			if len(opId) == 0 {
				if ruleset.HasRule(RuleOpIdNonEmpty) {
					vsets.AddSimple(RuleOpIdNonEmpty, opLoc, "")
				}
			} else {
				if ruleset.HasRule(RuleOpIdStyleCamelCase) &&
					!stringcase.IsCamelCase(opId) {
					vsets.AddSimple(
						RuleOpIdStyleCamelCase,
						opLoc+"/operationId",
						opId)
				}
			}
			// Check Parameters
			vsets.UpsertSets(
				ParametersCheck(
					op.Parameters,
					opLoc+"/parameters",
					ruleset))
			if ruleset.HasRule(RuleTagCaseFirstAlphaUpper) {
				for i, tag := range op.Tags {
					if !stringcase.IsFirstAlphaUpper(tag) {
						vsets.AddSimple(
							RuleTagCaseFirstAlphaUpper,
							opLoc+"/tags/"+strconv.Itoa(i),
							tag)
					}
				}
			}
		},
	)
	return vsets
}
