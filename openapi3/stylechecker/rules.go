package stylechecker

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/simplego/text/stringcase"
	"github.com/grokify/simplego/type/maputil"
	"github.com/grokify/simplego/type/stringsutil"
)

const (
	// Optional
	RuleOpDescNonEmpty                = "operation.description.non_empty"
	RuleOpIdNonEmpty                  = "operation.operationid.non_empty"
	RuleOpIdStyleCamelCase            = "operation.operationid.style.camelcase"
	RuleOpTagsCount1                  = "operation.tags.count.1"
	RulePathParamNameNonEmpty         = "path.param.name.non-empty"
	RulePathParamStyleCamelCase       = "path.param.style.camelcase"
	RuleSchemaPropEnumStylePascalCase = "schema.property.enum.style.pascalcase"
	RuleSchemaPropDescNonEmpty        = "schema.property.description.non-empty"
	RuleSchemaObjectPropsExist        = "schema.object.properties.exists"
	RuleTagCaseFirstAlphaUpper        = "tag.case.first.alpha.upper"

	// Mandatory
	RuleOpParameterNameNonEmpty = "operation.parameter.name.non-empty"

	// Prefixes
	PrefixPathParam          = "path.param."
	PrefixSchemaPropertyEnum = "schema.property.enum.style."
)

func RuleToCaseStyle(s string) string {
	infoMap := map[string]string{
		RuleSchemaPropEnumStylePascalCase: stringcase.CasePascal}
	if caseStyle, ok := infoMap[s]; ok {
		return caseStyle
	}
	return ""
}

type RuleSet struct {
	rulesMap map[string]int
}

func NewRuleSet(rules []string) RuleSet {
	rules = stringsutil.SliceCondenseSpace(rules, true, true)
	for i, rule := range rules {
		rules[i] = strings.ToLower(rule)
	}
	msi := maputil.NewMapStringIntSlice(rules)
	return RuleSet{rulesMap: msi}
}

func (set *RuleSet) HasRule(rule string) bool {
	rule = strings.ToLower(strings.TrimSpace(rule))
	if _, ok := set.rulesMap[rule]; ok {
		return true
	}
	return false
}

func (set *RuleSet) HasPathItemRules() bool {
	return set.HasRulePrefix(PrefixPathParam)
}

func (set *RuleSet) HasSchemaEnumStyleRules() bool {
	return set.HasRulePrefix(PrefixSchemaPropertyEnum)
}

func (set *RuleSet) HasRulePrefix(prefix string) bool {
	for rule := range set.rulesMap {
		if strings.Index(rule, prefix) == 0 {
			return true
		}
	}
	return false
}

func (set *RuleSet) RulesWithPrefix(prefix string) []string {
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
