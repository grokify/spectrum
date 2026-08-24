# Inspecting Specs

`SpecMore` wraps a spec (`*openapi3.T`) with convenience methods for reading its
structure without hand-walking the object graph.

```go
sm, err := openapi3.ReadSpecMore("openapi.yaml", false)
if err != nil {
	log.Fatal(err)
}
```

You can also wrap a spec you already have:

```go
sm := openapi3.SpecMore{Spec: spec}
```

## Counts

```go
sm.OperationsCount() // number of operations across all paths
sm.SchemasCount()    // number of component schemas
```

## Operations

```go
sm.OperationIDs()                 // []string of every operationId
sm.OperationIDsCounts()           // map[operationId]count (spot duplicates)
sm.OperationCountsByTag()         // histogram of operations per tag

// Look up a single operation by ID or by path+method:
path, method, op, err := sm.OperationByID("getFeature")
op, err := sm.OperationByPathMethod("/features/{id}", http.MethodGet)
```

`PathMethods(true)` returns every endpoint as a generic `path method` string
(path variables normalized), which is the basis for spec comparison:

```go
for _, ep := range sm.PathMethods(true) {
	fmt.Println(ep) // e.g. "/features/{} GET"
}
```

## Schemas

```go
sm.SchemaNames()                  // []string of component schema names
sm.SchemaRef("Feature")           // *openapi3.SchemaRef for one schema

// Reconcile referenced vs defined schemas:
noRef, both, refNoSchema, err := sm.SchemaNamesStatus()
// noRef        - defined but never referenced
// refNoSchema  - referenced but not defined (dangling $ref)
```

## Metadata

`Metadata()` bundles the three comparable dimensions into one struct:

```go
md := sm.Metadata()
md.Endpoints    // []string generic "path method"
md.OperationIDs // []string
md.SchemaNames  // []string
```

This is the same `SpecMetadata` used by [spec comparison](compare.md).

## Tabular export

Emit operations as a table for CSV/XLSX reporting:

```go
if err := sm.WriteFileCSV("operations.csv", nil); err != nil {
	log.Fatal(err)
}
```
