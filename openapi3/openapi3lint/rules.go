package openapi3lint

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/simplego/text/stringcase"
	"github.com/grokify/simplego/type/stringsutil"
)

const (
	// Optional
	RuleOpDescNotEmpty                = "operation.description.not_empty"
	RuleOpIdNotEmpty                  = "operation.operationid.not_empty"
	RuleOpIdStyleCamelCase            = "operation.operationid.style.camelcase"
	RuleOpSummaryNotEmpty             = "operation.summary.not_empty"
	RuleOpSummaryCaseFirstCapitalized = "operation.summary.case.first.capitalized"
	RuleOpTagsCountOneOnly            = "operation.tags.count.one_only"
	RulePathParamNameNonEmpty         = "path.param.name.not_empty"
	RulePathParamStyleCamelCase       = "path.param.style.camelcase"
	RuleSchemaPropEnumStylePascalCase = "schema.property.enum.style.pascalcase"
	RuleSchemaPropDescNotEmpty        = "schema.property.description.not_empty"
	RuleSchemaObjectPropsExist        = "schema.object.properties.exists"
	RuleTagCaseFirstCapitalized       = "tag.case.first.capitalized"

	// Mandatory
	RuleOpParameterNameNonEmpty = "operation.parameter.name.non-empty"

	// Prefixes
	PrefixPathParam          = "path.param."
	PrefixSchemaPropertyEnum = "schema.property.enum.style."

	RuleInternalError = "internal.error"

	LocationPaths   = "#/paths"
	LocationSchemas = "#/components/schemas"
)

func RuleToCaseStyle(s string) string {
	infoMap := map[string]string{
		RuleSchemaPropEnumStylePascalCase: stringcase.CasePascal}
	if caseStyle, ok := infoMap[s]; ok {
		return caseStyle
	}
	return ""
}

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

var rxSlashMore = regexp.MustCompile(`/+`)

func PointerCondense(s string) string {
	return rxSlashMore.ReplaceAllString(s, "/")
}

func PointerSubEscapeAll(format string, vars ...interface{}) string {
	if len(vars) == 0 {
		return format
	}
	for i, v := range vars {
		if vString, ok := v.(string); ok {
			vars[i] = jsonutil.PropertyNameEscape(vString)
		}
	}
	return PointerCondense(fmt.Sprintf(format, vars...))
}
