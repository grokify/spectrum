# Identify Missing Descriptions

Descriptions are important to aid understanding of various objects in the OpenAPI spec.

Swaggman provides an ability to list operation parameters and schema properties with missing descriptions.

## Missing Operation Parameter Descriptions

```golang
specmore := openapi3.SpecMore{Spec: spec}

// OperationPropertiesDescriptionStatus returns a
// map[string]map[string]int as a `simplego/maputil.MapStringMapStringInt`
// where `1` indicates with desc and `0` indicates without desc.
status := specmore.OperationPropertiesDescriptionStatus()

// OperationParametersWithoutDescriptionsWriteFile
// will write the operationIds and param names to a file
err := specmore.OperationParametersWithoutDescriptionsWriteFile(
    "missing-descs_op-params.txt")
```

## Missing Schema Property Descriptions

```golang
specmore := openapi3.SpecMore{Spec: spec}

// SchemaPropertiesDescriptionStatus returns a
// map[string]map[string]int as a `simplego/maputil.MapStringMapStringInt`
// where `1` indicates with desc and `0` indicates without desc.
status := specmore.SchemaPropertiesDescriptionStatus()

// SchemaPropertiesWithoutDescriptionsWriteFile
// will write the schema names and property names to a file
err := specmore.SchemaPropertiesWithoutDescriptionsWriteFile(
    "missing-descs_schema-props.txt")
```