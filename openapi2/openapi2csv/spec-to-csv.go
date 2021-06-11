package openapi2csv

import (
	"net/http"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/grokify/gocharts/data/table"
	oas2 "github.com/grokify/spectrum/openapi2"
)

func TableFromSpecFiles(files []string, includeFilename bool) (*table.Table, error) {
	tbl := table.NewTable()
	tblp := &tbl
	tbl.Columns = []string{}
	if includeFilename {
		tbl.Columns = append(tbl.Columns, "Filename")
	}
	tbl.Columns = append(tbl.Columns, []string{"Path", "Method", "OperationID", "Summary", "Description"}...)
	for _, file := range files {
		spec, err := oas2.ReadOpenAPI2KinSpecFile(file)
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

func TableAddOpenAPI2Spec(tbl *table.Table, spec *openapi2.Swagger, prefix []string) *table.Table {
	for url, path := range spec.Paths {
		tbl = TableAddOpenAPI2Path(tbl, path, append(prefix, url))
	}
	return tbl
}

// prefix can be `filename`,`path`
func TableAddOpenAPI2Path(tbl *table.Table, path *openapi2.PathItem, prefix []string) *table.Table {
	if path.Delete != nil {
		tbl.Rows = append(tbl.Rows, pathOpenApi2ToRow(prefix, path.Delete, http.MethodDelete))
	}
	if path.Get != nil {
		tbl.Rows = append(tbl.Rows, pathOpenApi2ToRow(prefix, path.Get, http.MethodGet))
	}
	if path.Head != nil {
		tbl.Rows = append(tbl.Rows, pathOpenApi2ToRow(prefix, path.Head, http.MethodHead))
	}
	if path.Options != nil {
		tbl.Rows = append(tbl.Rows, pathOpenApi2ToRow(prefix, path.Options, http.MethodOptions))
	}
	if path.Patch != nil {
		tbl.Rows = append(tbl.Rows, pathOpenApi2ToRow(prefix, path.Patch, http.MethodPatch))
	}
	if path.Post != nil {
		tbl.Rows = append(tbl.Rows, pathOpenApi2ToRow(prefix, path.Post, http.MethodPost))
	}
	if path.Put != nil {
		tbl.Rows = append(tbl.Rows, pathOpenApi2ToRow(prefix, path.Put, http.MethodPut))
	}
	return tbl
}

func pathOpenApi2ToRow(prefix []string, op *openapi2.Operation, method string) []string {
	row := prefix
	row = append(row, method, op.OperationID, op.Summary, op.Description)
	return row
}
