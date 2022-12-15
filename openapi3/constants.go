package openapi3

const (
	OASVersionLatest  = "3.1.0"
	OASVersionDefault = "3.0.3"
	apiVersionDefault = "0.0.1"

	TypeArray      = "array"
	TypeBoolean    = "boolean"
	TypeInteger    = "integer"
	TypeObject     = "object"
	TypeString     = "string"
	FormatDate     = "date"
	FormatDateTime = "date-time"
	FormatInt32    = "int32"
	FormatInt64    = "int64"

	PropertyOperationID = "operationId"
	PropertySummary     = "summary"
	PropertyTags        = "tags"

	InCookie = "cookie"
	InHeader = "header"
	InPath   = "path"
	InQuery  = "query"

	PointerComponentsSchemas       = "#/components/schemas"
	PointerComponentsSchemasFormat = `#/components/schemas/%s`
)
