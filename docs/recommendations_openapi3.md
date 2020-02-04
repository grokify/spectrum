# OpenAPI 3 Spec

This document is a recommended set of minimal OpenAPI 3 spec properties so that a spec can be used by various ecosystem tools such as API References, API Explorers and Client SDK generators.

## Operation

| Property | Requirement | Notes |
|----------|-------------|-------|
| `description` | `MUST` | Used by API References |
| `operationId` | `MUST` | May be used to auto-generate client SDK method names |
| `responses` | `MUST` | Minimally must have 2xx successful response. Other responses including errors are desirable |
| `summary` | `MUST` | Used in API References, such as Swagger UI and ReadMe.io |
| `tags` | `MUST` | There should have 1 and only 1 tag. Tags are used to organize endpoints in API References and Client SDKs. More than 1 tag may not be supported will in some software |

### Operation Parameter

| Property | Requirement | Notes |
|----------|-------------|-------|
| `in` | `MUST` | This describes where the parameter appears. |
| `name` | `MUST` | This is the name of the parameter. |
| `required` | `MUST` | This is the name of the parameter. |
| `description` | `SHOULD` | This describes where the parameter appears. |
| `schema.type` or `schema.$ref` | `MUST` | Type property must be present. Schema paramters are typically not objects which would be defined by a `$ref`, though JSON bodes are. |
| `schema.format` | `SHOULD` | Format property should be present. For `integer` type, if using `long`, explicitly set `format` to `int64`. For Date Time properties, only set `format` to `date-time` or `date` if the fields correspond to IETF RFC-3339. If date/time formats do not correspond to RFC-3339, leave `format` empty and add format information in the `description` property |

## Schema

| Property | Requirement | Notes |
|----------|-------------|-------|
| `required` | `MUST` | When required fields are present, they should be indicated |

### Schema Property

| Property | Requirement | Notes |
|----------|-------------|-------|
| `description` | `SHOULD` | Description should be included when available. For fields with ambiguous `format` information such as non-RFC-3339 date/times, the format should be included in the description. |
| `type` or `$ref` | `MUST` | Type property must be present. |
| `format` | `SHOULD` | Format property should be present. For `integer` type, if using `long`, explicitly set `format` to `int64`. For Date Time properties, only set `format` to `date-time` or `date` if the fields correspond to IETF RFC-3339. If date/time formats do not correspond to RFC-3339, leave `format` empty and add format information in the `description` property |