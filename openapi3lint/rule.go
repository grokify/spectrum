package openapi3lint

import (
	"fmt"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/log/severity"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/grokify/swaggman/openapi3lint/lintutil"
)

type Rule interface {
	Name() string
	Scope() string
	Severity() string
	ProcessSpec(spec *oas3.Swagger, pointerBase string) *lintutil.PolicyViolationsSets
	ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation
}

type RulesMap map[string]Rule

func ValidateRules(rules map[string]Rule) error {
	unknownSeverities := []string{}
	for ruleName, rule := range rules {
		_, err := severity.Parse(rule.Severity())
		if err != nil {
			unknownSeverities = append(unknownSeverities, ruleName)
		}
	}
	if len(unknownSeverities) > 0 {
		unknownSeverities = stringsutil.Dedupe(unknownSeverities)
		sort.Strings(unknownSeverities)
		return fmt.Errorf(
			"rules with unknown severities [%s] valid [%s]",
			strings.Join(unknownSeverities, ","),
			strings.Join(severity.Severities(), ","))
	}
	return nil
}
