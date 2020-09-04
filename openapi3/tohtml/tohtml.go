package tohtml

import (
	"fmt"
	"html"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/gotilla/net/urlutil"
	"github.com/grokify/gotilla/type/stringsutil"
	"github.com/grokify/swaggman/openapi3"
)

func ToTable(spec *oas3.Swagger, baseURL string) *table.TableData {
	baseURL = strings.TrimSpace(baseURL)
	tbl := table.NewTableData()
	tbl.Name = spec.Info.Title

	tbl.Columns = []string{
		"Method", "Path", "OperationID", "Summary", "Tags", "API Group", "Throttling", "App Permission", "User Permissions",
	}

	openapi3.VisitOperations(spec, func(
		path, method string, op *oas3.Operation) {
		pathLink := path

		if len(baseURL) > 0 {
			tagSlug := ""
			if len(op.Tags) > 0 {
				op.Tags = stringsutil.SliceCondenseSpace(op.Tags, true, false)
				if len(op.Tags) > 0 {
					tagSlug = strings.Replace(op.Tags[0], " ", "-", -1)
				}
			}
			pathLink = fmt.Sprintf(
				"<a href=\"%s\">%s</a>",
				urlutil.JoinAbsolute(baseURL, tagSlug, op.OperationID),
				path)
		}

		row := []string{
			html.EscapeString(method),
			pathLink,
			html.EscapeString(op.OperationID),
			html.EscapeString(op.Summary),
			html.EscapeString(strings.Join(op.Tags, ", ")),
			html.EscapeString(openapi3.GetExtensionPropStringOrEmpty(op.ExtensionProps, "x-api-group")),
			html.EscapeString(openapi3.GetExtensionPropStringOrEmpty(op.ExtensionProps, "x-throttling-group")),
			html.EscapeString(openapi3.GetExtensionPropStringOrEmpty(op.ExtensionProps, "x-app-permission")),
			html.EscapeString(openapi3.GetExtensionPropStringOrEmpty(op.ExtensionProps, "x-user-permission")),
		}
		tbl.Records = append(tbl.Records, row)
	})
	return &tbl
}

func ToTableHTML(spec *oas3.Swagger, domID, baseURL string) string {
	return table.ToHTML(ToTable(spec, baseURL), domID, false)
}

func ToHTMLPage(spec *oas3.Swagger) string {
	return ""
}

func ToDocuments(spec *oas3.Swagger) []map[string]interface{} {
	return table.ToDocuments(ToTable(spec, ""))
}
