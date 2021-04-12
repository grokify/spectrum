package openapi3lint

import (
	"fmt"
	"regexp"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/grokify/simplego/text/stringcase"
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

const (
	SeverityDisabled    = "disabled"
	SeverityError       = "error"
	SeverityHint        = "hint"
	SeverityInformation = "information"
	SeverityWarning     = "warning"
)

func RuleToCaseStyle(s string) string {
	infoMap := map[string]string{
		RuleSchemaPropEnumStylePascalCase: stringcase.CasePascal}
	if caseStyle, ok := infoMap[s]; ok {
		return caseStyle
	}
	return ""
}

type Rule struct {
	Name     string
	Severity string
	Func     func(spec *oas3.Swagger, ruleset Policy) PolicyViolationsSets
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
