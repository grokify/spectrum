package lintutil

const (
	RulenameDatatypeIntFormatStandardExist = "datatype-int-format-standard-exist"
	RuleOpDescExist                        = "operation-description-exist"
	RuleOpIDExist                          = "operation-operationid-exist"

	RulenameOpIDStyleCamelCase  = "operation-operationid-style-camelcase"
	RulenameOpIDStyleKebabCase  = "operation-operationid-style-kebabcase"
	RulenameOpIDStylePascalCase = "operation-operationid-style-pascalcase"
	RulenameOpIDStyleSnakeCase  = "operation-operationid-style-snakecase"

	RulenameSchemaNameStylePascalCase = "schema-name-style-pascalcase"
	RulenameSchemaHasReference        = "schema-has-reference"
	RulenameSchemaReferenceHasSchema  = "schema-reference-has-schema"

	RulenameOpSummaryExist               = "operation-summary-exist"
	RulenameOpSummaryStyleFirstUpperCase = "operation-summary-style-first-uppercase"

	RuleOpTagsCountOneOnly = "operation-tags-count-one"
	RulePathParamNameExist = "path-param-name-exist"

	RulenamePathParamStyleCamelCase  = "path-param-style-camelcase"
	RulenamePathParamStyleKebabCase  = "path-param-style-kebabcase"
	RulenamePathParamStylePascalCase = "path-param-style-pascalcase"
	RulenamePathParamStyleSnakeCase  = "path-param-style-snakecase"

	RulenameSchemaPropEnumStyleCamelCase  = "schema-property-enum-style-camelcase"
	RulenameSchemaPropEnumStyleKebabCase  = "schema-property-enum-style-kebabcase"
	RulenameSchemaPropEnumStylePascalCase = "schema-property-enum-style-pascalcase"
	RulenameSchemaPropEnumStyleSnakeCase  = "schema-property-enum-style-snakecase"

	RuleSchemaPropDescExist = "property-description-exist"

	RulenameSchemaObjectPropsExist = "schema-object-properties-exist"

	RulenameTagStyleFirstUpperCase = "tag-style-first-uppercase"
)
