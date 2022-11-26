package openapi3lint

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/log/severity"
	"github.com/grokify/mogo/path/filepathutil"
	"github.com/grokify/mogo/text/stringcase"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type Policy struct {
	//rules       map[string]Rule
	policyRules map[string]PolicyRule
}

func NewPolicy() Policy {
	return Policy{
		//rules:       map[string]Rule{},
		policyRules: map[string]PolicyRule{}}
}

func (pol *Policy) AddRule(rule Rule, sev string, errorOnCollision bool) error {
	ruleName := rule.Name()
	if len(strings.TrimSpace(ruleName)) == 0 {
		return errors.New("rule has no name Policy.AddRule()")
	}
	if !stringcase.IsKebabCase(ruleName) {
		return fmt.Errorf("rule to add name must be in in kebab-case format [%s]", ruleName)
	}
	if errorOnCollision {
		if _, ok := pol.policyRules[ruleName]; ok {
			return fmt.Errorf("duplicate rule [%s] Policy.AddRule()", ruleName)
		}
	}
	canonicalSeverity := severity.SeverityError
	if len(strings.TrimSpace(sev)) > 0 {
		canonicalSeverityTry, err := severity.Parse(sev)
		if err != nil {
			return fmt.Errorf("severity not found [%s] Policy.AddRule()", sev)
		}
		canonicalSeverity = canonicalSeverityTry
	}
	pol.policyRules[ruleName] = PolicyRule{
		Rule:     rule,
		Severity: canonicalSeverity}
	return nil
}

/*
func (pol *Policy) addRuleWithPriorError(rule Rule, sev string, err error) error {
	if err != nil {
		return err
	}
	return pol.AddRule(rule, sev, true)
}
*/
/*
func (pol *Policy) AddRule(rule Rule, errorOnCollision bool) error {
	if len(rule.Name()) == 0 {
		return errors.New("rule to add must have non-empty name")
	}
	if !stringcase.IsKebabCase(rule.Name()) {
		return fmt.Errorf("rule to add name must be in in kebab-case format [%s]", rule.Name())
	}
	if _, ok := pol.rules[rule.Name()]; ok {
		if errorOnCollision {
			return fmt.Errorf("add rule collision for [%s]", rule.Name())
		}
	}
	pol.rules[rule.Name()] = rule
	return nil
}
*/

func (pol *Policy) RuleNames() []string {
	ruleNames := []string{}
	for rn := range pol.policyRules {
		ruleNames = append(ruleNames, rn)
	}
	sort.Strings(ruleNames)
	return ruleNames
}

func (pol *Policy) ValidateSpec(spec *openapi3.Spec, pointerBase, filterSeverity string) (*lintutil.PolicyViolationsSets, error) {
	vsets := lintutil.NewPolicyViolationsSets()

	unknownScopes := []string{}
	for _, policyRule := range pol.policyRules {
		_, err := lintutil.ParseScope(policyRule.Rule.Scope())
		if err != nil {
			unknownScopes = append(unknownScopes, policyRule.Rule.Scope())
		}
	}
	if len(unknownScopes) > 0 {
		return nil, fmt.Errorf("bad policy: rules have unknown scopes [%s]",
			strings.Join(unknownScopes, ","))
	}

	vsetsOps, err := pol.processRulesOperation(spec, pointerBase, filterSeverity)
	if err != nil {
		return vsets, err
	}
	err = vsets.UpsertSets(vsetsOps)
	if err != nil {
		return vsets, err
	}

	vsetsSpec, err := pol.processRulesSpecification(spec, pointerBase, filterSeverity)
	if err != nil {
		return vsets, err
	}
	err = vsets.UpsertSets(vsetsSpec)
	if err != nil {
		return vsets, err
	}

	return vsets, nil
}

func (pol *Policy) processRulesSpecification(spec *openapi3.Spec, pointerBase, filterSeverity string) (*lintutil.PolicyViolationsSets, error) {
	if spec == nil {
		return nil, errors.New("cannot process nil spec")
	}
	vsets := lintutil.NewPolicyViolationsSets()

	for _, policyRule := range pol.policyRules {
		if !lintutil.ScopeMatch(lintutil.ScopeSpecification, policyRule.Rule.Scope()) {
			continue
		}
		inclRule, err := severity.SeverityInclude(filterSeverity, policyRule.Severity)
		if err != nil {
			return vsets, err
		}
		// fmt.Printf("FILTER_SEV [%v] ITEM_SEV [%v] INCL [%v]\n", filterSeverity, rule.Severity(), inclRule)
		if inclRule {
			//fmt.Printf("PROC RULE name[%s] scope[%s] sev[%s]\n", rule.Name(), rule.Scope(), rule.Severity())
			vsets.AddViolations(policyRule.Rule.ProcessSpec(spec, pointerBase))
		}
	}
	return vsets, nil
}

func (pol *Policy) processRulesOperation(spec *openapi3.Spec, pointerBase, filterSeverity string) (*lintutil.PolicyViolationsSets, error) {
	vsets := lintutil.NewPolicyViolationsSets()

	severityErrorRules := []string{}
	unknownSeverities := []string{}

	openapi3.VisitOperations(spec,
		func(path, method string, op *oas3.Operation) {
			if op == nil {
				return
			}
			opPointer := jsonutil.PointerSubEscapeAll(
				"%s#/paths/%s/%s", pointerBase, path, strings.ToLower(method))
			for _, policyRule := range pol.policyRules {
				if !lintutil.ScopeMatch(lintutil.ScopeOperation, policyRule.Rule.Scope()) {
					continue
				}
				//fmt.Printf("HERE [%s] RULE [%s] Scope [%s]\n", path, rule.Name(), rule.Scope())
				inclRule, err := severity.SeverityInclude(filterSeverity, policyRule.Severity)
				//fmt.Printf("INCL_RULE? [%v] RULE [%s]\n", inclRule, rule.Name())
				if err != nil {
					severityErrorRules = append(severityErrorRules, policyRule.Rule.Name())
					unknownSeverities = append(unknownSeverities, policyRule.Severity)
				} else if inclRule {
					vsets.AddViolations(policyRule.Rule.ProcessOperation(spec, op, opPointer, path, method))
				}
			}
		},
	)

	if len(severityErrorRules) > 0 || len(unknownSeverities) > 0 {
		severityErrorRules = stringsutil.Dedupe(severityErrorRules)
		sort.Strings(severityErrorRules)
		return vsets, fmt.Errorf(
			"rules with unknown severities rules[%s] severities[%s] valid[%s]",
			strings.Join(unknownSeverities, ","),
			strings.Join(severityErrorRules, ","),
			strings.Join(severity.Severities(), ","))
	}

	return vsets, nil
}

var ErrNoSpecFiles = errors.New("no spec files supplied")

// ValidateSpecFiles executes the policy against a set of one or more spec files.
// `sev` is the severity as specified by `github.com/grokify/mogo/log/severity`.
// A benefit of using this over `ValidateSpec()` when validating multiple files
// is that this will automatically inject the filename as a JSON pointer base.`
func (pol *Policy) ValidateSpecFiles(filterSeverity string, specfiles []string) (*lintutil.PolicyViolationsSets, error) {
	if len(specfiles) == 0 {
		return nil, ErrNoSpecFiles
	}
	severityLevel, err := severity.Parse(filterSeverity)
	if err != nil {
		return nil, err
	}

	vsets := lintutil.NewPolicyViolationsSets()
	for _, file := range specfiles {
		spec, err := openapi3.ReadFile(file, false)
		if err != nil {
			return nil, err
		}
		vsetsRule, err := pol.ValidateSpec(spec, filepathutil.FilepathLeaf(file), severityLevel)
		if err != nil {
			return nil, err
		}
		err = vsets.UpsertSets(vsetsRule)
		if err != nil {
			return nil, err
		}
	}

	return vsets, nil
}
