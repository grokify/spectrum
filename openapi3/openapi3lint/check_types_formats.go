package openapi3lint

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
		/*
			openapi3.VisitOperations(
				spec,
				func(path, method string, op *oas3.Operation) {
					if op == nil {
						return
					}
					for i, paramRef := range op.Parameters {
						if paramRef.Value == nil ||
							paramRef.Value.Schema == nil ||
							paramRef.Value.Schema.Value == nil {
							continue
						}
						if paramRef.Value.Schema.Value.Type == openapi3.TypeInteger &&
							paramRef.Value.Schema.Value.Format != openapi3.FormatInt32 &&
							paramRef.Value.Schema.Value.Format != openapi3.FormatInt64 {

							opLoc := PointerSubEscapeAll(
								"#/paths/%s/%s/parameters/%d/schema/format", path, strings.ToLower(method), i)

							vsets.AddSimple(
								RuleDatatypeIntFormatIsInt32OrInt64,
								opLoc,
								paramRef.Value.Schema.Value.Format)
						}
					}
				},
			)*/
	}

	return vsets
}
