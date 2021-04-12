package openapi3lint

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

func (set *Policy) HasRule(rule string) bool {
	rule = strings.ToLower(strings.TrimSpace(rule))
	if _, ok := set.rulesMap[rule]; ok {
		return true
	}
	return false
}

func (set *Policy) HasPathItemRules() bool {
	return set.HasRulePrefix(PrefixPathParam)
}

func (set *Policy) HasSchemaEnumStyleRules() bool {
	return set.HasRulePrefix(PrefixSchemaPropertyEnum)
}

func (set *Policy) HasRulePrefix(prefix string) bool {
	for rule := range set.rulesMap {
		if strings.Index(rule, prefix) == 0 {
			return true
		}
	}
	return false
}

func (set *Policy) RulesWithPrefix(prefix string) []string {
	rules := []string{}
	for rule := range set.rulesMap {
		if strings.Index(rule, prefix) == 0 {
			rules = append(rules, rule)
		}
	}
	return rules
}
