package openapi3lint1

import (
	"strings"

	"github.com/grokify/simplego/type/stringsutil"
)

type Policy struct {
	rulesMap map[string]Rule
}

func NewPolicySimple(rules []string) Policy {
	pol := Policy{rulesMap: map[string]Rule{}}
	rules = stringsutil.SliceCondenseSpace(rules, true, true)
	for i, rule := range rules {
		rules[i] = strings.ToLower(rule)
		pol.rulesMap[rule] = Rule{
			Name:     rule,
			Severity: SeverityError}
	}
	return pol
}

func (pol *Policy) Validate() error {
	for ruleName, rule := range pol.rulesMap {
		if ruleName != rule.Name {
			if len(rule.Name) == 0 {
				rule.Name = ruleName
				pol.rulesMap[ruleName] = rule
			}
		}
	}
	return nil
}

func (pol *Policy) HasRule(rule string) bool {
	rule = strings.ToLower(strings.TrimSpace(rule))
	if _, ok := pol.rulesMap[rule]; ok {
		return true
	}
	return false
}

func (pol *Policy) HasPathItemRules() bool {
	return pol.HasRulePrefix(PrefixPathParam)
}

func (pol *Policy) HasSchemaEnumStyleRules() bool {
	return pol.HasRulePrefix(PrefixSchemaPropertyEnum)
}

func (pol *Policy) HasRulePrefix(prefix string) bool {
	for rule := range pol.rulesMap {
		if strings.Index(rule, prefix) == 0 {
			return true
		}
	}
	return false
}

func (pol *Policy) RulesWithPrefix(prefix string) []string {
	rules := []string{}
	for rule := range pol.rulesMap {
		if strings.Index(rule, prefix) == 0 {
			rules = append(rules, rule)
		}
	}
	return rules
}
