# Getting Started

## Install

### Library

```bash
go get github.com/grokify/spectrum@latest
```

The most-used package is `openapi3`:

```go
import "github.com/grokify/spectrum/openapi3"
```

### CLIs

Install any command directly:

```bash
go install github.com/grokify/spectrum/cmd/oas3diff@latest
go install github.com/grokify/spectrum/cmd/oas3validate@latest
go install github.com/grokify/spectrum/cmd/oas3lint@latest
go install github.com/grokify/spectrum/cmd/openapi2to3@latest
```

See the [CLI reference](cli.md) for the full list.

## Read a spec

`ReadFile` loads a spec from JSON or YAML (the format is auto-detected) and
optionally validates it. It returns a `*Spec`, which is an alias for
kin-openapi's `openapi3.T`.

```go
spec, err := openapi3.ReadFile("openapi.yaml", false) // pass true to validate
if err != nil {
	log.Fatal(err)
}
```

To parse bytes you already hold in memory, use `Parse`:

```go
spec, err := openapi3.Parse(specBytes)
```

## Inspect a spec

Wrap a spec in `SpecMore` for the convenience API, or read straight into one
with `ReadSpecMore`:

```go
sm, err := openapi3.ReadSpecMore("openapi.yaml", false)
if err != nil {
	log.Fatal(err)
}

fmt.Println("operations:", sm.OperationsCount())
fmt.Println("schemas:   ", sm.SchemasCount())
fmt.Println("operationIDs:", sm.OperationIDs())
```

See [Inspecting specs](guide/inspect.md) for the full inspection API.

## Next steps

- [Inspecting specs](guide/inspect.md)
- [Comparing specs](guide/compare.md)
- [Merging specs](guide/merge.md)
- [Converting specs](guide/convert.md)
