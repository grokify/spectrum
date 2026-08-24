# Merging Specs

Spectrum merges multiple OpenAPI 3 specs into a single spec — useful for
assembling one API document from per-service or per-domain fragments.

## Merge files

```go
merged, err := openapi3.MergeFiles(
	[]string{"users.yaml", "orders.yaml", "billing.yaml"},
	nil, // *MergeOptions, or nil for defaults
)
if err != nil {
	log.Fatal(err)
}
```

## Merge a directory

Merge every spec in a directory:

```go
merged, count, err := openapi3.MergeDirectory("specs/", nil)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("merged %d specs\n", count)
```

## Merge two specs

`Merge` combines a master spec with an extra spec. The `specExtraNote` is used to
annotate provenance of the merged-in content:

```go
merged, err := openapi3.Merge(master, extra, "from orders service", nil)
```

## Options

Pass `*MergeOptions` to control conflict handling. A skip-on-collision policy is
available via `NewMergeOptionsSkip()`:

```go
opts := openapi3.NewMergeOptionsSkip()
merged, err := openapi3.MergeFiles(files, opts)
```

Finer-grained helpers exist when you only need to merge a subset:
`MergePaths`, `MergeParameters`, `MergeResponses`, `MergeSchemas`, and
`MergeTags`.
