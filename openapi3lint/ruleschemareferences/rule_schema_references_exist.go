package ruleschemareferences

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleSchemaReferences struct {
	name                              string
	severity                          string
	checkSchemaWithoutReference       bool
	checkSchemaReferenceWithoutSchema bool
}

func NewRuleSchemaReferences(sev, ruleName string) (RuleSchemaReferences, error) {
	ruleNameCanonical := strings.ToLower(strings.TrimSpace(ruleName))
	rule := RuleSchemaReferences{
		name:     ruleNameCanonical,
		severity: sev}
	if ruleNameCanonical != lintutil.RulenameSchemaWithoutReference &&
		ruleNameCanonical != lintutil.RulenameSchemaReferenceWithoutSchema {
		return rule, fmt.Errorf("rule [%s] not supported", ruleName)
	}
	return rule, nil
}

func (rule RuleSchemaReferences) Name() string {
	return rule.name
}

func (rule RuleSchemaReferences) Scope() string {
	return lintutil.ScopeSpecification
}

func (rule RuleSchemaReferences) Severity() string {
	return rule.severity
}

func (rule RuleSchemaReferences) ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}

func (rule RuleSchemaReferences) ProcessSpec(spec *oas3.Swagger, pointerBase string) []lintutil.PolicyViolation {
	violations := []lintutil.PolicyViolation{}

	sm := openapi3.SpecMore{Spec: spec}
	schemaNoRef, _, refNoSchema, err := sm.SchemaNamesStatus()
	if err != nil {
		return violations
	}
	if rule.name == lintutil.RulenameSchemaWithoutReference {
		for _, schemaName := range schemaNoRef {
			violations = append(violations, lintutil.PolicyViolation{
				RuleName: lintutil.RulenameSchemaWithoutReference,
				Location: openapi3.SchemaPointerExpand(pointerBase, schemaName)})
		}
	} else if rule.name == lintutil.RulenameSchemaReferenceWithoutSchema {
		for _, schemaName := range refNoSchema {
			violations = append(violations, lintutil.PolicyViolation{
				RuleName: lintutil.RulenameSchemaReferenceWithoutSchema,
				Location: openapi3.SchemaPointerExpand(pointerBase, schemaName)})
		}
	}

	return violations
}
