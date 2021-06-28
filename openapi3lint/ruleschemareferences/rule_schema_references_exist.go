package ruleschemareferences

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type RuleSchemaReferences struct {
	name string
}

func NewRule(ruleName string) (RuleSchemaReferences, error) {
	ruleNameCanonical := strings.ToLower(strings.TrimSpace(ruleName))
	rule := RuleSchemaReferences{
		name: ruleNameCanonical}
	if ruleNameCanonical != lintutil.RulenameSchemaHasReference &&
		ruleNameCanonical != lintutil.RulenameSchemaReferenceHasSchema {
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

func (rule RuleSchemaReferences) ProcessOperation(spec *openapi3.Spec, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}

func (rule RuleSchemaReferences) ProcessSpec(spec *openapi3.Spec, pointerBase string) []lintutil.PolicyViolation {
	violations := []lintutil.PolicyViolation{}

	sm := openapi3.SpecMore{Spec: spec}
	schemaNoRef, _, refNoSchema, err := sm.SchemaNamesStatus()
	if err != nil {
		return violations
	}
	if rule.name == lintutil.RulenameSchemaHasReference {
		for _, schemaName := range schemaNoRef {
			violations = append(violations, lintutil.PolicyViolation{
				RuleName: lintutil.RulenameSchemaHasReference,
				Location: openapi3.SchemaPointerExpand(pointerBase, schemaName)})
		}
	} else if rule.name == lintutil.RulenameSchemaReferenceHasSchema {
		for _, schemaName := range refNoSchema {
			violations = append(violations, lintutil.PolicyViolation{
				RuleName: lintutil.RulenameSchemaReferenceHasSchema,
				Location: openapi3.SchemaPointerExpand(pointerBase, schemaName)})
		}
	}
	return violations
}
