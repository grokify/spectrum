package extensions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grokify/spectrum/openapi3lint"
	"github.com/grokify/spectrum/openapi3lint/extensions/ruletaggroupexist"
)

type RuleCollectionExtensions struct {
	name      string
	ruleNames map[string]int
}

func NewRuleCollectionExtensions() RuleCollectionExtensions {
	rules := RuleCollectionExtensions{
		name:      "Spectrum OpenAPI 3 Lint Extensions Rule Collection",
		ruleNames: map[string]int{}}
	names := rules.RuleNames()
	for _, name := range names {
		rules.ruleNames[name] = 1
	}
	return rules
}

func (std RuleCollectionExtensions) Name() string {
	return std.name
}

func (std RuleCollectionExtensions) RuleExists(ruleName string) bool {
	if _, ok := std.ruleNames[ruleName]; ok {
		return true
	}
	return false
}

func (std RuleCollectionExtensions) RuleNames() []string {
	rulenames := []string{
		ruletaggroupexist.RuleName,
	}
	sort.Strings(rulenames)
	return rulenames
}

func (std RuleCollectionExtensions) Rule(name, severity string) (openapi3lint.Rule, error) {
	sev := severity
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case ruletaggroupexist.RuleName:
		return ruletaggroupexist.NewRule(sev), nil
	}

	return openapi3lint.EmptyRule{}, fmt.Errorf("NewExtensionRule: rule [%s] not found", name)
}
