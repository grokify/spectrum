# Converting Specs

## OpenAPI 2 (Swagger) → OpenAPI 3

Convert an OAS2 file to an OAS3 file:

```go
import "github.com/grokify/spectrum/openapi2"

err := openapi2.ConvertOAS2FileToOAS3File(
	"swagger.json", // input OAS2
	"openapi.json", // output OAS3
	0644,           // output file permission
	true,           // pretty-print
)
```

From the command line:

```bash
openapi2to3 -i swagger.json -o openapi.json -p
```

## OpenAPI → Postman Collection v2

Convert a Swagger/OAS2 spec to a Postman Collection, optionally merging onto a
base Postman file:

```bash
openapi2postman -s swagger.json -p collection.postman.json
openapi2postman -s swagger.json -b base.postman.json -p collection.postman.json
```

Programmatically, use the converter in `openapi2postman2`:

```go
import "github.com/grokify/spectrum/openapi2/openapi2postman2"

conv := openapi2postman2.NewConverter(openapi2postman2.Configuration{})
err := conv.MergeConvert(swaggerFile, basePostmanFile, outPostmanFile)
```

## Export operations to CSV

Scan a directory of specs and emit an operations table as CSV:

```bash
openapi2csv -d ./specs -r '\.json$' -o operations.csv
```

- `-d` — source directory
- `-r` — regexp matching spec filenames
- `-o` — output CSV path
