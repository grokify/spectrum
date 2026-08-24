# CLI Reference

Spectrum ships several command-line tools under `cmd/`. Install any of them with:

```bash
go install github.com/grokify/spectrum/cmd/<name>@latest
```

Most tools use flags (via `go-flags`); pass `--help` to any command for the
authoritative list.

## oas3diff

Diff two OpenAPI 3 specs (JSON or YAML) across endpoints, operation IDs, and
schema names. _(Added in v1.21.0.)_

```bash
oas3diff <spec1> <spec2>
```

Positional arguments; see [Comparing specs](guide/compare.md).

## oas2meta

Print metadata (endpoints, operation IDs, schema names) for one or more specs as
JSON.

```bash
oas2meta <spec1> [<spec2> ...]
```

## oas3validate

Validate an OAS3 spec, optionally writing an operation/status-code report to XLSX.

```bash
oas3validate -s openapi.yaml
oas3validate -s openapi.yaml -x report.xlsx
```

- `-s`, `--specfile` — input OAS3 spec (required)
- `-x`, `--xlsxwrite` — write an XLSX report
- `--writeopstatus` — write an operation status-code report

## oas3lint

Lint an OAS3 spec (or directory) against a rule policy file.

```bash
oas3lint -p policy.json -i openapi.yaml
oas3lint -p policy.json -i ./specs -s error
```

- `-p`, `--policyfile` — lint policy file (required)
- `-i`, `--inputspec` — input spec file or directory
- `-s`, `--severity` — severity level filter

See [Linting](openapi3lint.md) and [Custom rules](openapi3lint/custom_rules.md).

## openapi2to3

Convert an OAS2 (Swagger) file to OAS3.

```bash
openapi2to3 -i swagger.json -o openapi.json -p
```

- `-i`, `--input` — input OAS2 file (required)
- `-o`, `--output` — output OAS3 file (required)
- `-p`, `--pretty` — pretty-print output

## openapi2postman

Convert a Swagger/OAS2 spec to a Postman Collection v2, optionally merging onto a
base Postman file.

```bash
openapi2postman -s swagger.json -p collection.postman.json
```

- `-s`, `--swagger` — input Swagger file (required)
- `-p`, `--postman` — output Postman file (required)
- `-b`, `--base` — base Postman file to merge onto

## openapi2csv

Scan a directory of specs and export an operations table to CSV.

```bash
openapi2csv -d ./specs -r '\.json$' -o operations.csv
```

- `-d`, `--dir` — source directory (required)
- `-r`, `--regexp` — filename match regexp (required)
- `-o`, `--output` — output CSV file (required)
