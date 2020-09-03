package tohtml

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/data/table"
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

func (pp *PageParams) AddSpec(spec *oas3.Swagger) error {
	return pp.AddTable(ToTable(spec, ""))
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
