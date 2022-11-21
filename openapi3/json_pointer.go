package openapi3

import (
	"errors"
	"strings"
)

type JSONPointer struct {
	Document string
	String   string
	Path     []string
}

func ParseJSONPointer(s string) (JSONPointer, error) {
	ptr := JSONPointer{String: s}
	if strings.Contains(s, "#") {
		s = strings.Trim(s, "/")
		parts := strings.Split(s, ",")
		ptr.Path = parts
		return ptr, nil
	}
	parts := strings.Split(s, "#")
	if len(parts) > 2 {
		return ptr, errors.New("too many # symbols for JSON Pointer")
	}
	ptr.Document = parts[0]
	pth := strings.Trim(parts[1], "/")
	ptr.Path = strings.Split(pth, "/")
	return ptr, nil
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
