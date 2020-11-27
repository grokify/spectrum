# Identify Missing Descriptions

Descriptions are important to aid understanding of various objects in the OpenAPI spec.

Swaggman provides an ability to list operation parameters and schema properties with missing descriptions.

## Missing Operation Parameter Descriptions

```golang
specmore := openapi3.SpecMore{Spec: spec}

// OperationPropertiesWithoutDescriptions returns a
// map[string]map[string]int as a `gotilla/maputil.MapStringMapStringInt`
missing := OperationPropertiesWithoutDescriptions()

// OperationParametersWithoutDescriptionsWriteFile
// will write the operationIds and param names to a file
err := OperationParametersWithoutDescriptionsWriteFile(
    "missing-descs_op-params.txt")
```

## Missing Schema Property Descriptions

```golang
specmore := openapi3.SpecMore{Spec: spec}

// SchemaPropertiesWithoutDescriptions returns a
// map[string]map[string]int as a `gotilla/maputil.MapStringMapStringInt`
missing := SchemaPropertiesWithoutDescriptions()
// SchemaPropertiesWithoutDescriptionsWriteFile
// will write the schema names and property names to a file

err := SchemaPropertiesWithoutDescriptionsWriteFile(
    "missing-descs_schema-props.txt")
```