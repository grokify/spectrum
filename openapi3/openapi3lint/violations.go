package openapi3lint

import (
	"errors"
	"fmt"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/type/stringsutil"
)

type SpecCheckViolationsFunc func(spec *oas3.Swagger, rules Policy) PolicyViolationsSets

func SpecCheckViolations(spec *oas3.Swagger, rules Policy) (PolicyViolationsSets, error) {
	vsets := NewPolicyViolationsSets()

	vioFuncs := []SpecCheckViolationsFunc{
		SpecCheckOperations,
		SpecCheckPathItems,
		SpecCheckSchemas,
	}

	for _, vioFunc := range vioFuncs {
		err := vsets.UpsertSets(vioFunc(spec, rules))
		if err != nil {
			return vsets, err
		}
	}

	return vsets, nil
}

type PolicyViolationsSets struct {
	ByRule map[string]PolicyViolationsSet
}

func NewPolicyViolationsSets() PolicyViolationsSets {
	return PolicyViolationsSets{
		ByRule: map[string]PolicyViolationsSet{}}
}

func (sets *PolicyViolationsSets) AddSimple(ruleName, location, value string) {
	set, ok := sets.ByRule[ruleName]
	if !ok {
		set = NewPolicyViolationsSet(ruleName)
	}
	set.Violations = append(set.Violations,
		PolicyViolation{
			RuleName: ruleName,
			Location: location,
			Value:    value})
	sets.ByRule[ruleName] = set
}

func (sets *PolicyViolationsSets) UpsertSet(upsertSet PolicyViolationsSet) error {
	for _, vio := range upsertSet.Violations {
		ruleName := vio.RuleName
		if len(ruleName) == 0 {
			ruleName = upsertSet.RuleName
		}
		if len(ruleName) == 0 {
			return errors.New("violation & violationSet have no RuleName")
		}
		existingSet, ok := sets.ByRule[upsertSet.RuleName]
		if !ok {
			sets.ByRule[upsertSet.RuleName] = upsertSet
		} else {
			existingSet.Violations = append(
				existingSet.Violations, vio)
			sets.ByRule[ruleName] = existingSet
		}
	}
	return nil
}

func (sets *PolicyViolationsSets) UpsertSets(upsertSets PolicyViolationsSets) error {
	for ruleName, upsertSet := range upsertSets.ByRule {
		if ruleName != upsertSet.RuleName {
			return fmt.Errorf("set name mismatch sets[%s] set[%s]", ruleName, upsertSet.RuleName)
		}
		err := sets.UpsertSet(upsertSet)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sets *PolicyViolationsSets) LocationsByRule() ViolationLocationsByRuleSet {
	locs := map[string][]string{}
	for _, set := range sets.ByRule {
		for _, vio := range set.Violations {
			ruleName := vio.RuleName
			vioLocation := vio.Location
			vioValue := vio.Value
			if len(vioValue) > 0 {
				vioLocation += " [" + vioValue + "]"
			}
			if _, ok := locs[ruleName]; !ok {
				locs[ruleName] = []string{}
			}
			locs[ruleName] = append(locs[ruleName], vioLocation)
		}
	}
	vlrs := ViolationLocationsByRuleSet{
		ViolationLocationsByRule: locs}
	vlrs.Condense()
	return vlrs
}

type PolicyRule struct {
	Name         string
	StringFormat string
}

type PolicyViolationsSet struct {
	RuleName   string
	Violations []PolicyViolation
}

func NewPolicyViolationsSet(ruleName string) PolicyViolationsSet {
	return PolicyViolationsSet{
		RuleName:   ruleName,
		Violations: []PolicyViolation{}}
}

func (set *PolicyViolationsSet) Locations() PolicyViolationLocations {
	locations := PolicyViolationLocations{
		RuleName:  set.RuleName,
		Locations: []string{}}
	for _, v := range set.Violations {
		locations.Locations = append(locations.Locations, v.Location)
	}
	locations.Finalize()
	return locations
}

type PolicyViolationLocations struct {
	RuleName  string
	Locations []string
}

func (vl *PolicyViolationLocations) Finalize() {
	vl.Locations = stringsutil.SliceCondenseSpace(vl.Locations, true, true)
}

type PolicyViolation struct {
	RuleName  string
	RuleType  string
	Violation string
	Value     string
	Location  string
	Data      map[string]string
}

type ViolationLocationsByRuleSet struct {
	ViolationLocationsByRule map[string][]string
}

func (vlrs *ViolationLocationsByRuleSet) Condense() {
	for k, vals := range vlrs.ViolationLocationsByRule {
		vlrs.ViolationLocationsByRule[k] =
			stringsutil.SliceCondenseSpace(vals, true, true)
	}
}
