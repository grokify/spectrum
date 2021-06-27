package openapi3html

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/spectrum/openapi3"
)

type PageParams struct {
	PageTitle     string
	PageLink      string
	TableDomID    string
	Spec          *oas3.Swagger
	ColumnSet     *table.ColumnSet
	OpsFilterFunc func(path, method string, op *oas3.Operation) bool
	TableJSON     []byte
}

func (pp *PageParams) PageLinkHTML() string {
	pp.PageLink = strings.TrimSpace(pp.PageLink)
	if len(pp.PageLink) == 0 {
		return html.EscapeString(pp.PageTitle)
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>", pp.PageLink,
		html.EscapeString(pp.PageTitle))
}

func (pp *PageParams) AddSpec(spec *oas3.Swagger) error {
	sm := openapi3.SpecMore{Spec: spec}
	tbl, err := sm.OperationsTable(pp.ColumnSet, pp.OpsFilterFunc)
	if err != nil {
		return err
	}
	return pp.AddOperationsTable(tbl)
}

func (pp *PageParams) AddOperationsTable(tbl *table.Table) error {
	docs := tbl.ToDocuments()
	jdocs, err := json.Marshal(docs)
	if err != nil {
		return err
	}
	pp.TableJSON = jdocs
	return nil
}

func (pp *PageParams) TableJSONBytesOrEmpty() []byte {
	if len(pp.TableJSON) > 0 {
		return pp.TableJSON
	}
	if pp.Spec != nil {
		pp.AddSpec(pp.Spec)
		return pp.TableJSON
	}
	return []byte("[]")
}

func (pp *PageParams) TabulatorColumnsJSONBytesOrEmpty() []byte {
	if pp.ColumnSet == nil || len(pp.ColumnSet.Columns) == 0 {
		colSet := openapi3.OpTableColumnsDefault(false)
		tcols := table.BuildColumnsTabulator(colSet.Columns)
		return tcols.MustColumnsJSON()
	}
	tcols := table.BuildColumnsTabulator(pp.ColumnSet.Columns)
	return tcols.MustColumnsJSON()
}

func (pp *PageParams) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	WriteSpectrumUIPage(f, *pp)
	return nil
}

/*
func DefaultColumns() text.Texts {
	return text.Texts{
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
}
*/
