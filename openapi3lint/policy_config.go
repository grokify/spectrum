package openapi3lint

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/stringsutil"
)

type PolicyConfig struct {
	Name                 string                `json:"name"`
	Version              string                `json:"version"`
	LastUpdated          time.Time             `json:"lastUpdated,omitempty"`
	IncludeStandardRules bool                  `json:"includeStandardRules"`
	Rules                map[string]RuleConfig `json:"rules,omitempty"`
	NonStandardRules     []string              `json:"nonStandardRules,omitempty"`
	xRuleCollections     RuleCollections       `json:"-"`
}

func NewPolicyConfigFile(filename string) (PolicyConfig, error) {
	pol := PolicyConfig{}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return pol, err
	}
	return pol, json.Unmarshal(bytes, &pol)
}

const (
	RuleTypeAll        = "all"
	RuleTypeStandard   = "standard"
	RuleTypeXDefined   = "xdefined"
	RuleTypeXUndefined = "xundefined"
)

func (polCfg *PolicyConfig) RuleNames() map[string][]string {
	ruleNamesMap := map[string][]string{
		RuleTypeAll:        {},
		RuleTypeStandard:   {},
		RuleTypeXDefined:   {},
		RuleTypeXUndefined: {}}
	stdRules := NewRuleCollectionStandard()
	xRuleNames := map[string]int{} // defined = 1, undefined 0
	for ruleName := range polCfg.Rules {
		ruleNamesMap[RuleTypeAll] = append(ruleNamesMap[RuleTypeAll], ruleName)
		if polCfg.IncludeStandardRules &&
			stdRules.RuleExists(ruleName) {
			ruleNamesMap[RuleTypeStandard] =
				append(ruleNamesMap[RuleTypeStandard], ruleName)
			continue
		}
		if len(polCfg.xRuleCollections) == 0 {
			xRuleNames[ruleName] = 0
		} else {
			for _, ruleCollection := range polCfg.xRuleCollections {
				if ruleCollection.RuleExists(ruleName) {
					xRuleNames[ruleName] = 1
				} else {
					if _, ok := xRuleNames[ruleName]; !ok {
						xRuleNames[ruleName] = 0
					}
				}
			}
		}
	}
	for ruleName, ruleVal := range xRuleNames {
		if ruleVal >= 1 {
			ruleNamesMap[RuleTypeXDefined] =
				append(ruleNamesMap[RuleTypeXDefined], ruleName)
		} else {
			ruleNamesMap[RuleTypeXUndefined] =
				append(ruleNamesMap[RuleTypeXUndefined], ruleName)
		}
	}
	for ruleType, ruleNames := range ruleNamesMap {
		ruleNamesMap[ruleType] =
			stringsutil.SliceCondenseSpace(ruleNames, true, true)
	}
	return ruleNamesMap
}

func (polCfg *PolicyConfig) AddRuleCollection(collection RuleCollection) {
	polCfg.xRuleCollections = append(polCfg.xRuleCollections, collection)
}

type RuleConfig struct {
	Severity string `json:"severity"`
}

func (polCfg *PolicyConfig) Policy() (Policy, error) {
	pol := NewPolicy()
	stdRules := NewRuleCollectionStandard()
	ruleCollectionsMap := map[string][]string{}

	for ruleName, ruleCfg := range polCfg.Rules {
		if polCfg.IncludeStandardRules {
			if stdRules.RuleExists(ruleName) {
				if _, ok := ruleCollectionsMap[ruleName]; !ok {
					ruleCollectionsMap[ruleName] = []string{}
				}
				ruleCollectionsMap[ruleName] = append(ruleCollectionsMap[ruleName], stdRules.Name())
				rule, err := stdRules.Rule(ruleName)
				if err != nil {
					return pol, errorsutil.Wrap(err, "standard error not found. PolicyConfig.Policy()")
				}
				if err = pol.AddRule(rule, ruleCfg.Severity, true); err != nil {
					return pol, errorsutil.Wrap(err, fmt.Sprintf("Policy.AddRule() [%s]", ruleName))
				}
				/*
					if err := pol.addRuleWithPriorError(stdRules.Rule(ruleName, ruleCfg.Severity)); err != nil {
						return pol, errorsutil.Wrap(err, fmt.Sprintf("pol.addRuleWithPriorError [%s]", ruleName))
					}*/
			}
		}
		for _, collection := range polCfg.xRuleCollections {
			if collection.RuleExists(ruleName) {
				if _, ok := ruleCollectionsMap[ruleName]; !ok {
					ruleCollectionsMap[ruleName] = []string{}
				}
				ruleCollectionsMap[ruleName] = append(ruleCollectionsMap[ruleName], collection.Name())
				rule, err := collection.Rule(ruleName)
				if err != nil {
					return pol, errorsutil.Wrap(err, "collection rule exists but not found. PolicyConfig.Policy()")
				}
				if err = pol.AddRule(rule, ruleCfg.Severity, true); err != nil {
					return pol, errorsutil.Wrap(err, fmt.Sprintf("Policy.AddRule() [%s]", ruleName))
				}
				/*if err := pol.addRuleWithPriorError(collection.Rule(ruleName, ruleCfg.Severity)); err != nil {
					return pol, errorsutil.Wrap(err, fmt.Sprintf("pol.addRuleWithPriorError [%s]", ruleName))
				}*/
			}
		}
	}

	collisions := map[string][]string{}
	for ruleName, collections := range ruleCollectionsMap {
		if len(collections) > 1 {
			collisions[ruleName] = collections
		}
	}

	if len(collisions) > 0 {
		bytes, err := json.Marshal(collisions)
		if err != nil {
			return pol, errorsutil.Wrap(err, fmt.Sprintf("json.Marshal [%s]", string(bytes)))
		}
		return pol, fmt.Errorf("rule collisions: %s", string(bytes))
	}

	return pol, nil
}
