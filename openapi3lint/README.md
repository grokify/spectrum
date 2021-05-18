# OpenAPI 3 Lint

This is a Go-based OpenAPI 3 spec linter.

## Why Spectrum Linter?

There are a few linters available.

The reasons this exists are:

1. written in Go so its easy to use/modify for Go devs
2. policy violations are grouped by rule, vs. line number for easier mitigation
3. policy violations are identified by JSON Schema pointere vs. line number for easier identification when using merged files

## Standard Rules

The following rules are built into Spectral.

1. `datatype-int-format-int32-int64`: ensures data types of `integer` have a standard format of `int32` or `int64`.
1. `operation-operationid-style-camelcase`: ensures operationIds use camelCase
1. `operation-operationid-style-kebabcase`: ensures operationIds use kebab-case
1. `operation-operationid-style-pascalcase`: ensures operationIds use PascalCase
1. `operation-operationid-style-snakecase`: ensures operationIds use snake_case
1. `schema-has-reference`: ensures schemas have references
1. `schema-reference-has-schema`: ensures schma JSON pointers reference existing schemas

## Other Linters

There are other linters available. To date, Spectrum Linter hasn't beeen inspired by any of them, though there is a desire and effort to align on rule names and potentially rule definitions to achieve similar behavior.

1. Mermade OAS-Kit - https://github.com/mermade/oas-kit
    1. https://mermade.github.io/oas-kit/default-rules.html
    1. https://mermade.github.io/oas-kit/linter-rules.html
1. Spectral - https://github.com/stoplightio/spectral
    1. Inspired by Speccy
1. Speccy - https://github.com/wework/speccy
    1. Inspired by OAS-Kit