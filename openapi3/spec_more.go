package openapi3

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/gotilla/type/stringsutil"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/encoding/jsonutil"
	"github.com/grokify/gotilla/text"
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

func (s *SpecMore) SchemaCount() int {
	if s.Spec == nil {
		return -1
	} else if s.Spec.Components.Schemas == nil {
		return 0
	}
	return len(s.Spec.Components.Schemas)
}

func (s *SpecMore) OperationsTable(columns *text.TextSet) (*table.Table, error) {
	return operationsTable(s.Spec, columns)
}

func operationsTable(spec *oas3.Swagger, columns *text.TextSet) (*table.Table, error) {
	if columns == nil {
		columns = &text.TextSet{Texts: DefaultColumns()}
	}
	tbl := table.NewTable()
	tbl.Name = spec.Info.Title
	tbl.Columns = columns.DisplayTexts()

	tgs, err := SpecTagGroups(spec)
	if err != nil {
		return nil, err
	}

	VisitOperations(spec, func(path, method string, op *oas3.Operation) {
		row := []string{}

		for _, text := range columns.Texts {
			switch text.Slug {
			case "method":
				row = append(row, method)
			case "path":
				row = append(row, path)
			case "operationId":
				row = append(row, op.OperationID)
			case "summary":
				row = append(row, op.Summary)
			case "tags":
				row = append(row, strings.Join(op.Tags, ", "))
			case "x-tag-groups":
				row = append(row, strings.Join(
					tgs.GetTagGroupNamesForTagNames(op.Tags...), ", "))
			default:
				row = append(row, GetExtensionPropStringOrEmpty(op.ExtensionProps, text.Slug))
			}
		}

		tbl.Records = append(tbl.Records, row)
	})
	return &tbl, nil
}

func DefaultColumns() []text.Text {
	texts := []text.Text{
		{
			Display: "Method",
			Slug:    "method"},
		{
			Display: "Path",
			Slug:    "path"},
		{
			Display: "OperationID",
			Slug:    "operationId"},
		{
			Display: "Summary",
			Slug:    "summary"},
		{
			Display: "Tags",
			Slug:    "tags"},
	}
	return texts
}

/*
func (s *SpecMore) OperationsTableOld() (*table.TableData, error) {
	tbl := table.NewTableData()
	tbl.Name = "Operations"
	tgs, err := SpecTagGroups(s.Spec)
	if err != nil {
		return nil, err
	}
	addTagGroups := false
	if len(tgs.TagGroups) > 0 {
		addTagGroups = true
		tbl.Columns = []string{"OperationId", "Summary", "Path", "Method", "Tag Groups", "Tags"}
	} else {
		tbl.Columns = []string{"OperationId", "Summary", "Path", "Method", "Tags"}
	}
	ops := s.OperationMetas()
	for _, op := range ops {
		if addTagGroups {
			tagGroupNames := tgs.GetTagGroupNamesForTagNames(op.Tags...)
			tbl.Records = append(tbl.Records, []string{
				op.OperationID,
				op.Summary,
				op.Path,
				op.Method,
				strings.Join(tagGroupNames, ","),
				strings.Join(stringsutil.SliceCondenseSpace(op.Tags, true, true), ",")})
		} else {
			tbl.Records = append(tbl.Records, []string{
				op.OperationID,
				op.Summary,
				op.Path,
				op.Method,
				strings.Join(op.Tags, ",")})
		}
	}
	return &tbl, nil
}
*/

func (s *SpecMore) OperationMetas() []OperationMeta {
	ometas := []OperationMeta{}
	if s.Spec == nil {
		return ometas
	}
	for url, path := range s.Spec.Paths {
		if path.Connect != nil {
			ometas = append(ometas, OperationToMeta(url, http.MethodConnect, path.Connect))
		}
		if path.Delete != nil {
			ometas = append(ometas, OperationToMeta(url, http.MethodDelete, path.Delete))
		}
		if path.Get != nil {
			ometas = append(ometas, OperationToMeta(url, http.MethodGet, path.Get))
		}
		if path.Head != nil {
			ometas = append(ometas, OperationToMeta(url, http.MethodHead, path.Head))
		}
		if path.Options != nil {
			ometas = append(ometas, OperationToMeta(url, http.MethodOptions, path.Options))
		}
		if path.Patch != nil {
			ometas = append(ometas, OperationToMeta(url, http.MethodPatch, path.Patch))
		}
		if path.Post != nil {
			ometas = append(ometas, OperationToMeta(url, http.MethodPost, path.Post))
		}
		if path.Put != nil {
			ometas = append(ometas, OperationToMeta(url, http.MethodPut, path.Put))
		}
		if path.Trace != nil {
			ometas = append(ometas, OperationToMeta(url, http.MethodTrace, path.Trace))
		}
	}

	return ometas
}

func (s *SpecMore) OperationsCount() uint {
	return uint(len(s.OperationMetas()))
}

func (sm *SpecMore) SchemaNames() []string {
	schemaNames := []string{}
	for schemaName := range sm.Spec.Components.Schemas {
		schemaNames = append(schemaNames, schemaName)
	}
	return stringsutil.SliceCondenseSpace(schemaNames, true, true)
}

func (sm *SpecMore) SchemaNameExists(schemaName string, includeNil bool) bool {
	for schemaNameTry, schemaRef := range sm.Spec.Components.Schemas {
		if schemaNameTry == schemaName {
			if includeNil {
				return true
			} else if schemaRef == nil {
				return false
			}
			schemaRef.Ref = strings.TrimSpace(schemaRef.Ref)
			if len(schemaRef.Ref) > 0 {
				return true
			}
			if schemaRef.Value == nil {
				return false
			} else {
				return true
			}
		}
	}
	return false
}

func (s *SpecMore) WriteFileJSON(filename string, perm os.FileMode, prefix, indent string) error {
	jsonData, err := s.Spec.MarshalJSON()
	if err != nil {
		return err
	}
	pretty := false
	if len(prefix) > 0 || len(indent) > 0 {
		pretty = true
	}
	if pretty {
		jsonData = jsonutil.PrettyPrint(jsonData, "", "  ")
	}
	return ioutil.WriteFile(filename, jsonData, perm)
}

func (sm *SpecMore) WriteFileXLSX(filename string) error {
	tbl, err := sm.OperationsTable(nil)
	if err != nil {
		return err
	}
	return table.WriteXLSX(filename, tbl)
}

type TagsMore struct {
	Tags oas3.Tags
}

func (tg *TagsMore) Get(tagName string) *oas3.Tag {
	for _, tag := range tg.Tags {
		if tagName == tag.Name {
			return tag
		}
	}
	return nil
}
