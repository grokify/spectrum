# Filtering Specs by Audience

One spec can describe an API for several audiences: what everyone may call,
what selected partners may call, and what only the API's own team uses.
Spectrum writes the spec for one audience from it, keeping only the
operations that audience may see and only the components, security schemes
and tags those operations use. _(Added in [v1.22.0](../releases/v1.22.0.md).)_

Filtering removes operations from the file itself. Marking them hidden for a
documentation tool is not enough on its own, since anyone can download the
spec the tool reads.

## Marking operations

Mark an operation, or a whole path, with `x-audience`:

```yaml
paths:
  /v1/orders:
    get:
      x-audience: public
  /v1/orders/bulk-sync:
    post:
      x-audience: partner
  /v1/admin/reconcile:
    post:
      x-audience: internal
```

| Audience | Who sees it |
|----------|-------------|
| `internal` | The API's own team |
| `partner` | Selected customers, and the API's own team |
| `public` | Everyone |

Each audience includes the narrower ones, so a public operation is in the
partner and internal specs too. `x-audience` may also be a list, such as
`[partner, public]`, in which case the widest audience listed applies.

`x-internal: true`, as Redocly, Scalar and Stoplight read it, marks an
operation as internal. When an operation has both marks, `x-internal: true`
wins, so conflicting marks never widen who sees something.

An operation's own mark comes first, then its path item's, then a default.
Use `internal` as the default so an operation nobody marked is never
published by mistake. For a spec that marks only what is internal with
`x-internal`, use `public` as the default.

## Exporting

```go
keep, err := openapi3.AudienceFilter(openapi3.AudiencePublic, openapi3.AudienceInternal)
if err != nil {
	log.Fatal(err)
}

sm := openapi3.SpecMore{Spec: spec}
public, err := sm.ExportByFilter(keep, &openapi3.FilterOptions{
	// Drop the marks once they have done their job.
	RemoveExtensions: []string{openapi3.XAudience, openapi3.XInternal},
})
if err != nil {
	log.Fatal(err)
}
```

`ExportByFilter`:

- keeps the operations the filter accepts, in paths and webhooks, and drops
  path items left with none
- copies every component the kept operations reach, following references
  through `allOf`, `oneOf`, `anyOf`, `items`, `additionalProperties`, shared
  responses, request bodies, parameters, headers, examples, links and
  callbacks
- keeps only the security schemes the document and the kept operations
  require, and only the tags they use, in the source's order
- leaves the source spec unchanged
- reloads the result and fails on any reference that points nowhere

The output is the same every time for the same input, so it can be checked
in and diffed.

`ExportByFilter` takes any `OperationFilterFunc`, so the same pruning works
for filters by tag, path prefix or another extension.

## From the command line

```bash
oas3filter -i openapi.json -o openapi-public.json -a public
oas3filter -i openapi.json -o openapi-partner.yaml -a partner
```

| Flag | Meaning |
|------|---------|
| `-i`, `--input` | Input spec, JSON or YAML |
| `-o`, `--output` | Output spec; YAML when it ends in `.yaml` or `.yml`, else JSON |
| `-a`, `--audience` | `internal`, `partner` or `public` |
| `--default` | Audience of unmarked operations. Default `internal` |
| `--keep-extensions` | Keep `x-audience` and `x-internal` in the output |
