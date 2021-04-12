package openapi3lint

import oas3 "github.com/getkin/kin-openapi/openapi3"

/*
type Policy struct {
	Rules map[string]Rule
}
*/

const (
	SeverityDisabled    = "disabled"
	SeverityError       = "error"
	SeverityHint        = "hint"
	SeverityInformation = "information"
	SeverityWarning     = "warning"
)

type Rule struct {
	Name     string
	Severity string
	Func     func(spec *oas3.Swagger, ruleset Policy) PolicyViolationsSets
}
