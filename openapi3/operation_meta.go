package openapi3

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func OperationToMeta(url, method string, op *openapi3.Operation) OperationMeta {
	return OperationMeta{
		OperationID: strings.TrimSpace(op.OperationID),
		Summary:     strings.TrimSpace(op.Summary),
		Method:      strings.ToUpper(strings.TrimSpace(method)),
		Path:        strings.TrimSpace(url),
		Tags:        op.Tags,
		MetaNotes:   []string{}}
}

type OperationMeta struct {
	OperationID string
	Summary     string
	Method      string
	Path        string
	Tags        []string
	MetaNotes   []string
}
