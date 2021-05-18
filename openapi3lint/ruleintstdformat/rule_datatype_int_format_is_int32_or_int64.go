package ruleintstdformat

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleDatatypeIntFormatIsInt32OrInt64 struct {
	name     string
	severity string
}

func NewRuleDatatypeIntFormatIsInt32OrInt64(sev string) RuleDatatypeIntFormatIsInt32OrInt64 {
	return RuleDatatypeIntFormatIsInt32OrInt64{
		name:     lintutil.RulenameDatatypeIntFormatIsInt32OrInt64,
		severity: sev}
}

func (rule RuleDatatypeIntFormatIsInt32OrInt64) Name() string {
	return rule.name
}

func (rule RuleDatatypeIntFormatIsInt32OrInt64) Severity() string {
	return rule.severity
}

func (rule RuleDatatypeIntFormatIsInt32OrInt64) Scope() string {
	return lintutil.ScopeSpecification
}

func (rule RuleDatatypeIntFormatIsInt32OrInt64) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}

func (rule RuleDatatypeIntFormatIsInt32OrInt64) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	vios := []lintutil.PolicyViolation{}
	openapi3.VisitTypesFormats(
		spec,
		func(jsonPointerRoot, oasType, oasFormat string) {
			if oasType == openapi3.TypeInteger &&
				oasFormat != openapi3.FormatInt32 &&
				oasFormat != openapi3.FormatInt64 {
				vios = append(vios, lintutil.PolicyViolation{
					RuleName: rule.Name(),
					Location: urlutil.JoinAbsolute(pointerBase+jsonPointerRoot, "format"),
					Value:    oasFormat})
			}
		},
	)
	return vios
}
