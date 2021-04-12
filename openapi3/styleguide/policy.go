package styleguide

import oas3 "github.com/getkin/kin-openapi/openapi3"

/*
type Policy struct {
	Rules map[string]Rule
}
*/

const (
	RuleStatusDisabled    = "disabled"
	RuleStatusError       = "error"
	RuleStatusHint        = "hint"
	RuleStatusInformation = "information"
	RuleStatusWarning     = "warning"
)

type Rule struct {
	Name   string
	Status string
	Func   func(spec *oas3.Swagger, ruleset Policy) PolicyViolationsSets
}
