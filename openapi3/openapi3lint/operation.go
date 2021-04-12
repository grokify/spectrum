package openapi3lint

import (
	"strconv"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/text/stringcase"
	"github.com/grokify/swaggman/openapi3"
)

func SpecCheckOperations(spec *oas3.Swagger, ruleset Policy) PolicyViolationsSets {
	vsets := NewPolicyViolationsSets()

	openapi3.VisitOperations(spec,
		func(path, method string, op *oas3.Operation) {
			if op == nil {
				return
			}
			opLoc := PointerSubEscapeAll(
				"#/paths/%s/%s", path, strings.ToLower(method))
			opId := strings.TrimSpace(op.OperationID)
			if len(opId) == 0 {
				if ruleset.HasRule(RuleOpIdNotEmpty) {
					vsets.AddSimple(RuleOpIdNotEmpty, opLoc, "")
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
			requiredSummaryEmpty := false
			if ruleset.HasRule(RuleOpSummaryNotEmpty) {
				summaryCondensed := strings.TrimSpace(op.Summary)
				if len(summaryCondensed) == 0 {
					vsets.AddSimple(
						RuleOpSummaryNotEmpty,
						opLoc+"/summary",
						"")
					requiredSummaryEmpty = true
				}
			}
			if ruleset.HasRule(RuleOpSummaryCaseFirstCapitalized) &&
				!requiredSummaryEmpty {
				if !stringcase.IsFirstAlphaUpper(op.Summary) {
					vsets.AddSimple(
						RuleOpSummaryCaseFirstCapitalized,
						opLoc+"/summary",
						op.Summary)
				}
			}
			// Check Parameters
			err := vsets.UpsertSets(
				ParametersCheck(
					op.Parameters,
					opLoc+"/parameters",
					ruleset))
			if err != nil {
				vsets.AddSimple(RuleInternalError, opLoc, err.Error())
			}
			// Check Tags
			if ruleset.HasRule(RuleTagCaseFirstCapitalized) {
				for i, tag := range op.Tags {
					if !stringcase.IsFirstAlphaUpper(tag) {
						vsets.AddSimple(
							RuleTagCaseFirstCapitalized,
							opLoc+"/tags/"+strconv.Itoa(i),
							tag)
					}
				}
			}
		},
	)
	return vsets
}
