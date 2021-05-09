package openapi3lint1

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/swaggman/openapi3"
)

func SpecCheckDataTypesFormats(spec *oas3.Swagger, ruleset Policy) PolicyViolationsSets {
	vsets := NewPolicyViolationsSets()
	if ruleset.HasRule(RuleDatatypeIntFormatIsInt32OrInt64) {
		openapi3.VisitTypesFormats(
			spec,
			func(jsonPointerRoot, oasType, oasFormat string) {
				if oasType == openapi3.TypeInteger &&
					oasFormat != openapi3.FormatInt32 &&
					oasFormat != openapi3.FormatInt64 {
					vsets.AddSimple(
						RuleDatatypeIntFormatIsInt32OrInt64,
						urlutil.JoinAbsolute(jsonPointerRoot, "format"),
						oasFormat)
				}
			},
		)
	}
	return vsets
}
