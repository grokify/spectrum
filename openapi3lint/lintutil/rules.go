package lintutil

import "sort"

const (
	RulenameDatatypeIntFormatIsInt32OrInt64 = "datatype-int-format-int32-int64"
	RuleOpDescExist                         = "operation-description-exist"
	RuleOpIdExist                           = "operation-operationid-exist"

	RulenameOpIdStyleCamelCase  = "operation-operationid-style-camelcase"
	RulenameOpIdStyleKebabCase  = "operation-operationid-style-kebabcase"
	RulenameOpIdStylePascalCase = "operation-operationid-style-pascalcase"
	RulenameOpIdStyleSnakeCase  = "operation-operationid-style-snakecase"

	RulenameSchemaNameStylePascalCase    = "schema-name-style-pascalcase"
	RulenameSchemaWithoutReference       = "schema-without-reference"
	RulenameSchemaReferenceWithoutSchema = "schema-reference-without-schema"

	RuleOpSummaryExist                = "operation-summary-exist"
	RuleOpSummaryCaseFirstCapitalized = "operation-summary-first-letter-capitalized"
	RuleOpTagsCountOneOnly            = "operation-tags-count-one"
	RulePathParamNameExist            = "path-param-name-exist"
	RulePathParamStyleCamelCase       = "path-param-style-camelcase"
	RuleSchemaPropEnumStylePascalCase = "property-enum-style-pascalcase"
	RuleSchemaPropDescExist           = "property-description-exist"
	RuleSchemaObjectPropsExist        = "schema-object-properties-exist"
	RuleTagCaseFirstCapitalized       = "tag-case-first-capitalized"
)

func StandardRules() []string {
	rules := []string{
		RulenameDatatypeIntFormatIsInt32OrInt64,
		RulenameOpIdStyleCamelCase,
		RulenameOpIdStyleKebabCase,
		RulenameOpIdStylePascalCase,
		RulenameOpIdStyleSnakeCase}
	sort.Strings(rules)
	return rules
}
