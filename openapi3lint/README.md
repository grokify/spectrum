# OpenAPI 3 Lint

This is a Go-based OpenAPI 3 spec linter.

Spectrum is designed to support a multi-file, multi-user, async editing process where linting reports need to be resilinent to mutiple changes to the specs occurring between the time the validation is run and resolved.

## Why Spectrum Linter?

There are a few linters available.

The reasons this exists are:

1. written in Go so its easy to use/modify for Go devs
2. policy violations are grouped by rule, vs. line number for easier mitigation
3. policy violations are identified by JSON Schema pointere vs. line number for easier identification when using merged files

## Standard Rules

The following rules are built into Spectrum.

1. `datatype-int-format-standard-exist`: ensures data types of `integer` have a standard format (`int32` or `int64`)
1. `operation-operationid-style-camelcase`: ensures operationIds use camelCase
1. `operation-operationid-style-kebabcase`: ensures operationIds use kebab-case
1. `operation-operationid-style-pascalcase`: ensures operationIds use PascalCase
1. `operation-operationid-style-snakecase`: ensures operationIds use snake_case
1. `operation-summary-exist` ensures a summary exists.
1. `operation-summary-style-first-uppercase`: ensures summary starts with capitalized first character
1. `path-param-style-camelcase`: path parms are camel case
1. `path-param-style-kebabcase`: path parms are kebab case
1. `path-param-style-pascalcase`: path parms are Pascal case
1. `path-param-style-snakecase`: path parms are snake case
1. `schema-has-reference`: ensures schemas have references
1. `schema-object-properties-exist`: schema of type `object` have `properties` or `additionalProperties` defined
1. `schema-property-enum-style-camelcase`: schema property enums are camel case
1. `schema-property-enum-style-kebabcase`: schema property enums are kebab case
1. `schema-property-enum-style-pascalcase`: schema property enums are Pascal case
1. `schema-property-enum-style-snakecase`: schema property enums are snake case
1. `schema-reference-has-schema`: ensures schma JSON pointers reference existing schemas
1. `tag-style-first-uppercase`: Tag names have capitalized first character

## Other Linters

There are other linters available. To date, Spectrum Linter hasn't beeen inspired by any of them, though there is a desire and effort to align on rule names and potentially rule definitions to achieve similar behavior.

1. Mermade OAS-Kit - https://github.com/mermade/oas-kit
    1. https://mermade.github.io/oas-kit/default-rules.html
    1. https://mermade.github.io/oas-kit/linter-rules.html
1. Spectral - https://github.com/stoplightio/spectral
    1. Inspired by Speccy
1. Speccy - https://github.com/wework/speccy
    1. Inspired by OAS-Kit