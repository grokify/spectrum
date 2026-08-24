# Spectrum

Spectrum is a multi-purpose OpenAPI specification toolkit for Go. It wraps
[`getkin/kin-openapi`](https://github.com/getkin/kin-openapi) with convenience
methods for reading, inspecting, comparing, merging, linting, editing, and
converting OpenAPI specs — as an importable library and as a set of CLIs.

## What it does

### OpenAPI 3

- **Read & parse** JSON or YAML specs into a rich `SpecMore` wrapper
- **Inspect** operations, schemas, parameters, and spec metadata
- **Compare** two specs — intersection and diff across endpoints, operation IDs, and schema names
- **Merge** multiple specs into one
- **Lint** against configurable rule policies
- **Fix** and programmatically edit specs
- **Validate** specs
- **Convert** to Postman Collection v2

### OpenAPI 2 (Swagger)

- **Convert** OAS2 → OAS3
- **Convert** OAS2 → Postman Collection v2
- **Merge** multiple OAS2 specs

## Getting started

- [Getting Started](getting-started.md) — install the library and CLIs, read and inspect your first spec

## User guide

- [Inspecting specs](guide/inspect.md) — operations, schemas, and metadata via `SpecMore`
- [Comparing specs](guide/compare.md) — intersection and diff
- [Merging specs](guide/merge.md) — combine multiple specs
- [Converting specs](guide/convert.md) — OAS2 → OAS3, OpenAPI → Postman, CSV export
- [Fixing specs](openapi3_fix.md) — programmatic fixes
- [Recommendations](recommendations_openapi3.md) — authoring guidance

## Linting

- [Overview](openapi3lint.md)
- [Custom rules](openapi3lint/custom_rules.md)

## CLI reference

- [Command-line tools](cli.md) — `oas2meta`, `oas3diff`, `oas3validate`, `oas3lint`, `openapi2to3`, `openapi2postman`, `openapi2csv`

## Releases

- [v1.21.0](releases/v1.21.0.md)
