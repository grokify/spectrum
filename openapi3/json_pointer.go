package openapi3

import (
	"github.com/grokify/mogo/encoding/jsonpointer"
)

type JSONPointer jsonpointer.JSONPointer

func ParseJSONPointer(s string) (JSONPointer, error) {
	ptr, err := jsonpointer.ParseJSONPointer(s)
	return JSONPointer(ptr), err
}

const (
	PathComponents = "components"
	PathParameters = "parameters"
	PathPath       = "path"
	PathSchemas    = "schemas"
)

func (p *JSONPointer) IsTopParameter() (string, bool) {
	if len(p.Path) == 3 && p.Path[0] == PathComponents && p.Path[1] == PathParameters {
		return p.Path[2], true
	}
	return "", false
}

func (p *JSONPointer) IsTopSchema() (string, bool) {
	if len(p.Path) == 3 && p.Path[0] == PathComponents && p.Path[1] == PathSchemas {
		return p.Path[2], true
	}
	return "", false
}
