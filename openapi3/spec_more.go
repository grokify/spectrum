package openapi3

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/encoding/jsonutil"
)

type SpecMore struct {
	Spec *oas3.Swagger
}

func ReadSpecMore(path string, validate bool) (*SpecMore, error) {
	spec, err := ReadFile(path, validate)
	if err != nil {
		return nil, err
	}
	return &SpecMore{Spec: spec}, nil
}

/*
func (s *SpecMore) Paths() urlpath.SpecPaths {
	return urlpath.InspectPaths(s.Spec)
}
*/
func (s *SpecMore) OperationsTable() *table.TableData {
	tbl := table.NewTableData()
	tbl.Columns = []string{"operationId", "path", "url", "tags"}
	ops := s.OperationMetas()
	for _, op := range ops {
		tbl.Records = append(tbl.Records, []string{
			op.OperationID,
			op.Path,
			op.Method,
			strings.Join(op.Tags, ",")})
	}
	return &tbl
}

func (s *SpecMore) OperationMetas() []OperationMeta {
	ometas := []OperationMeta{}
	if s.Spec == nil {
		return ometas
	}
	for url, path := range s.Spec.Paths {
		if path.Connect != nil {
			ometas = append(ometas, opToMeta(url, http.MethodConnect, path.Connect))
		}
		if path.Delete != nil {
			ometas = append(ometas, opToMeta(url, http.MethodDelete, path.Delete))
		}
		if path.Get != nil {
			ometas = append(ometas, opToMeta(url, http.MethodGet, path.Get))
		}
		if path.Head != nil {
			ometas = append(ometas, opToMeta(url, http.MethodHead, path.Head))
		}
		if path.Options != nil {
			ometas = append(ometas, opToMeta(url, http.MethodOptions, path.Options))
		}
		if path.Patch != nil {
			ometas = append(ometas, opToMeta(url, http.MethodPatch, path.Patch))
		}
		if path.Post != nil {
			ometas = append(ometas, opToMeta(url, http.MethodPost, path.Post))
		}
		if path.Put != nil {
			ometas = append(ometas, opToMeta(url, http.MethodPut, path.Put))
		}
		if path.Trace != nil {
			ometas = append(ometas, opToMeta(url, http.MethodTrace, path.Trace))
		}
	}

	return ometas
}

func opToMeta(url, method string, op *openapi3.Operation) OperationMeta {
	return OperationMeta{
		OperationID: strings.TrimSpace(op.OperationID),
		Method:      strings.ToUpper(strings.TrimSpace(method)),
		Path:        strings.TrimSpace(url),
		Tags:        op.Tags}
}

type OperationMeta struct {
	OperationID string
	Method      string
	Path        string
	Tags        []string
}

func (s *SpecMore) WriteFileJSON(filename string, perm os.FileMode, pretty bool) error {
	jsonData, err := s.Spec.MarshalJSON()
	if err != nil {
		return err
	}
	if pretty {
		jsonData = jsonutil.PrettyPrint(jsonData, "", "  ")
	}
	return ioutil.WriteFile(filename, jsonData, perm)
}
