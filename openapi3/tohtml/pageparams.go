package tohtml

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/text"
	"github.com/grokify/swaggman/openapi3"
)

type PageParams struct {
	PageTitle  string
	PageLink   string
	TableDomID string
	TableJSON  []byte
}

func (pp *PageParams) PageLinkHTML() string {
	pp.PageLink = strings.TrimSpace(pp.PageLink)
	if len(pp.PageLink) == 0 {
		return html.EscapeString(pp.PageTitle)
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>", pp.PageLink,
		html.EscapeString(pp.PageTitle))
}

func (pp *PageParams) AddSpec(spec *oas3.Swagger, columns *text.TextSet) error {
	sm := openapi3.SpecMore{Spec: spec}
	tbl, err := sm.OperationsTable(columns)
	if err != nil {
		return err
	}
	return pp.AddTable(tbl)
}

func (pp *PageParams) AddTable(tbl *table.TableData) error {
	docs := table.ToDocuments(tbl)
	jdocs, err := json.Marshal(docs)
	if err != nil {
		return err
	}
	pp.TableJSON = jdocs
	return nil
}
