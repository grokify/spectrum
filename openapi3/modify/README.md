# Swaggman OpenAPI 3 Inspect, Modify & Compare

Swaggman `modify` is a library to assist in inspecting and modifying OpenAPI specs.

OpenAPI specifications can be large and have many endpoints which can make it difficult to manage. Additionally, some services may consist of many specs created by different people, teams and software, so some ability to make various specs consistent is desirable, especially when the individual specs need to be merged into a master spec.

Key Features include:

* Inspect: Various functions to examine aspects of a OpenAPI 3 spec including OperationIDs, paths, endpoint, schemas, tags, etc.
* Modify: Ability to modify various properties programmatically.
* Intersection: Ability to compare two specs and show the overlap.

## Usage

Steps for clean merging of multiple specs.

1. Examine all specs for consistent operationIds and tags
1. Ensure that all specs can be merged to common server base URL and paths
1. Optionally delete endpoint security from each spec and add it to merged spec
1. Check / deletee overlapping operationIDs, endpoints (method+path) and schema components
1. Validate resulting spec

### Inspect & Modify

Use `SpecMoreModifyMulti` and `SpecMoreModifyMultiOpts` to handle 
to inspect and modify mulitple files. 

### Compare

Use `modify.SpecsIntersection()`

```
// spec1 and spec2 are *github.com/getkin/kin-openapi/openapi3.Swagger
intersectionData := modify.SpecsIntersection(spec1, spec2)
intersectionData.Sort()
```

### Delete

After running intersection, you can use the resulting data to delete those items from a spec using `SpecDeleteProperties`. Be sure to validate afterwards.

This is useful when merging to specs with an overlap. To check for cleanliness of merging, you can:

1. run an intersection
1. delete the intersection from one of the sepcs and ensures it still validates
1. merge the specs


## Examples

### Add Bearer Token Auth

```go
modify.SecuritySchemeAddBearertoken(
    spec, "", "",
    []string{},
    []string{
        "Authentication",
    },
)
```