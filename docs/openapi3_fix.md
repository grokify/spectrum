# OpenAPI 3 Specs Fixer

In addition to enabling introspection and programmatic modification of
OpenAPI 3 specifications, swaggman can fix/update some issues with
specs.

## Path Parameters - Examination and Resolution

Path parameters are required to be defined. The following will add operation
path parameters if they are missing. It will also move path parameters to
the of the parameter list and maintain the order in which they appeaer in
the URL path. Otheer parameters will maintain their original order.

```go
ops, err = modify.ValidateFixOperationPathParameters(specFinal, true)
if err != nil {
    fmtutil.PrintJSON(ops)
    log.Fatal(err)
}
```

## Response Type - Examination and Resolution

Sometimes a spec can be misdefined to use a `application/json` response MIME
type when the schema returned doesn't ssupport JSON, e.g. with the type is
`string` or `integer`. The following will examine and optionally update the
type to `text/plain` to resolve then issue when the response is mis-classified
as `application/json` or some other incompatible MIME type.

```go
ops, err := modify.ValidateFixOperationResponseTypes(specFinal, true)
if err != nil {
    fmtutil.PrintJSON(ops)
    log.Fatal(err)
}
```

## Move Request Bodies

Some OpenAPI 3 spec defintions can use request body references like the following
which may not be supported by all tools.

```json
{
    "requestBody": {
        "$ref": "#/components/requestBodies/MyObject"
    }
}
```

Some tools are better able to handle a `requestBody` definition
as follows:

```json
{
    "requestBody": {
        "content": {
            "application/json": {
                "schema": {
                    "$ref": "#/components/schemas/MyObject"
                }
            }
        }
    }
}
```

The following will move the request body definition so that `content`
and MIME types are specified directly in the operation definition.

```go
ops, err = modify.MoveRequestBodies(specFinal, true)
if err != nil {
    fmtutil.PrintJSON(ops)
    log.Fatal(err)
}
```