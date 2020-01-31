package csv

import (
	"net/http"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/swaggman/swagger2"
)

func TableFromSpecFiles(files []string, includeFilename bool) (*table.TableData, error) {
	tbl := table.NewTableData()
	tblp := &tbl
	tbl.Columns = []string{}
	if includeFilename {
		tbl.Columns = append(tbl.Columns, "Filename")
	}
	tbl.Columns = append(tbl.Columns, []string{"Path", "Method", "Summary", "Description"}...)
	for _, file := range files {
		spec, err := swagger2.ReadOpenAPI2KinSpecFile(file)
		if err != nil {
			return tblp, err
		}
		prefix := []string{}
		if includeFilename {
			_, filename := filepath.Split(file)
			prefix = []string{filename}
		}
		tblp = TableAddOpenAPI2Spec(tblp, spec, prefix)
	}
	return tblp, nil
}

func TableAddOpenAPI2Spec(tbl *table.TableData, spec *openapi2.Swagger, prefix []string) *table.TableData {
	for url, path := range spec.Paths {
		tbl = TableAddOpenAPI2Path(tbl, path, append(prefix, url))
	}
	return tbl
}

// prefix can be `filename`,`path`
func TableAddOpenAPI2Path(tbl *table.TableData, path *openapi2.PathItem, prefix []string) *table.TableData {
	if path.Delete != nil {
		tbl.Records = append(tbl.Records, pathOpenApi2ToRow(prefix, path.Delete, http.MethodDelete))
	}
	if path.Get != nil {
		tbl.Records = append(tbl.Records, pathOpenApi2ToRow(prefix, path.Get, http.MethodGet))
	}
	if path.Head != nil {
		tbl.Records = append(tbl.Records, pathOpenApi2ToRow(prefix, path.Head, http.MethodHead))
	}
	if path.Options != nil {
		tbl.Records = append(tbl.Records, pathOpenApi2ToRow(prefix, path.Options, http.MethodOptions))
	}
	if path.Patch != nil {
		tbl.Records = append(tbl.Records, pathOpenApi2ToRow(prefix, path.Patch, http.MethodDelete))
	}
	if path.Post != nil {
		tbl.Records = append(tbl.Records, pathOpenApi2ToRow(prefix, path.Post, http.MethodDelete))
	}
	if path.Put != nil {
		tbl.Records = append(tbl.Records, pathOpenApi2ToRow(prefix, path.Put, http.MethodDelete))
	}
	return tbl
}

func pathOpenApi2ToRow(prefix []string, op *openapi2.Operation, method string) []string {
	row := prefix
	row = append(row, method, op.Summary, op.Description)
	return row
}
