# Comparing Specs

Spectrum compares two specs across three independent dimensions:

- **Endpoints** — generic `path method` strings (path variables normalized)
- **Operation IDs**
- **Schema names**

Because endpoints are compared generically, an `operationId` rename does **not**
hide the fact that the same endpoint exists in both specs.

## Diff

`SpecsDiff` reports what is unique to each spec and what they share. _(Added in
[v1.21.0](../releases/v1.21.0.md).)_

```go
diff, err := openapi3.ReadSpecsDiff("v1.json", "v2.yaml", false)
if err != nil {
	log.Fatal(err)
}

fmt.Print(diff.String())
```

The result is a `SpecMetadataDiff`:

```go
diff.OnlyInSpec1 // in spec 1 but not spec 2 (e.g. removed / lost coverage)
diff.OnlyInSpec2 // in spec 2 but not spec 1 (e.g. added)
diff.Both        // present in both

// Each side is a SpecMetadata, so you can index by dimension:
diff.OnlyInSpec2.Endpoints    // endpoints added in spec 2
diff.OnlyInSpec1.SchemaNames  // schemas dropped from spec 2
diff.Both.OperationIDs        // shared operation IDs

diff.IsEmpty() // true when the specs are equivalent across all dimensions
```

If you already hold parsed specs, use `SpecsDiff(spec1, spec2)` directly, or
`SpecMetadata.Diff(md2)` to compare two metadata sets.

### From the command line

```bash
oas3diff v1.json v2.yaml
```

```
Endpoints: 9 in both, 1 only in spec 1, 60 only in spec 2
  - only in spec 1: /products/{}/releases/{} PUT
  + only in spec 2: /epics GET
  ...
OperationIDs: ...
SchemaNames: ...
```

## Intersection

When you only need what two specs have in common, use `SpecsIntersection`:

```go
idata := openapi3.SpecsIntersection(spec1, spec2)
idata.Intersection.Endpoints    // endpoints in both
idata.Intersection.SchemaNames  // schemas in both
```

`SpecsDiff` is the complement of this API and reuses the same `SpecMetadata`
structure, so the two compose cleanly.

## Common uses

- **Migration safety** — confirm a regenerated or re-authored spec didn't drop
  endpoints or schemas (`diff.OnlyInSpec1.Endpoints` should be empty).
- **Release scope** — measure what a newer spec adds (`diff.OnlyInSpec2`).
- **Drift detection** — diff a published spec against the one your code generates.
