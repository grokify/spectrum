package ruletaggroupexist

import (
	"strconv"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

const (
	RuleName = "x-tag-has-group"
)

type RuleTagHasGroup struct {
	name     string
	severity string
}

func NewRule(sev string) RuleTagHasGroup {
	return RuleTagHasGroup{
		name:     RuleName,
		severity: sev}
}

func (rule RuleTagHasGroup) Name() string {
	return RuleName
}

func (rule RuleTagHasGroup) Scope() string {
	return lintutil.ScopeOperation
}

func (rule RuleTagHasGroup) Severity() string {
	return rule.severity
}

func (rule RuleTagHasGroup) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	if spec == nil || op == nil || len(op.Tags) == 0 {
		return nil
	}
	sm := openapi3.SpecMore{Spec: spec}
	tagGroups, err := sm.TagGroups()
	if err != nil {
		vio := lintutil.PolicyViolation{
			RuleName: rule.Name(),
			Location: "#/x-tag-groups",
			Value:    err.Error()}
		return []lintutil.PolicyViolation{vio}
	}

	vios := []lintutil.PolicyViolation{}
	for i, tagName := range op.Tags {
		if !tagGroups.Exists(tagName) {
			vios = append(vios, lintutil.PolicyViolation{
				RuleName: rule.Name(),
				Location: urlutil.JoinAbsolute(opPointer, openapi3.PropertyTags, strconv.Itoa(i)),
				Value:    tagName})
		}
	}
	return vios
}

func (rule RuleTagHasGroup) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}
