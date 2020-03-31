# Swaggman OpenAPI 3 Inspect, Modify & Compare

Swaggman `modify` is a library to assist in inspecting and modifying OpenAPI specs.

OpenAPI specifications can be large and have many endpoints which can make it difficult to manage. Additionally, some services may consist of many specs created by different people, teams and software, so some ability to make various specs consistent is desirable, especially when the individual specs need to be merged into a master spec.

Key Features include:

* Inspect: Various functions to examine aspects of a OpenAPI 3 spec including OperationIDs, paths, endpoint, schemas, tags, etc.
* Modify: Ability to modify various properties programmatically.
* Intersection: Ability to compare two specs and show the overlap.

## Usage

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
