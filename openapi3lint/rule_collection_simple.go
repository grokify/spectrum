package openapi3lint

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type RuleCollectionSimple struct {
	name  string
	rules map[string]Rule
}

func (simple RuleCollectionSimple) AddRule(rule Rule) error {
	ruleName := rule.Name()
	if len(strings.TrimSpace(ruleName)) == 0 {
		return errors.New("rule has no name")
	}
	if simple.rules == nil {
		simple.rules = map[string]Rule{}
	}
	simple.rules[ruleName] = rule
	return nil
}

func (simple RuleCollectionSimple) Name() string {
	if len(strings.TrimSpace(simple.name)) > 0 {
		return simple.name
	}
	return strings.Join(simple.RuleNames(), ",")
}

func (simple RuleCollectionSimple) RuleNames() []string {
	names := []string{}
	for ruleName := range simple.rules {
		names = append(names, ruleName)
	}
	sort.Strings(names)
	return names
}

func (simple RuleCollectionSimple) RuleExists(ruleName string) bool {
	if _, ok := simple.rules[ruleName]; ok {
		return true
	}
	return false
}

func (simple RuleCollectionSimple) Rule(ruleName string) (Rule, error) {
	if rule, ok := simple.rules[ruleName]; ok {
		return rule, nil
	}
	return EmptyRule{}, fmt.Errorf("rule not found [%s]", ruleName)
}
