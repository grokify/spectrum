package openapi3

import "fmt"

type Format struct {
	StandardType   string
	XFormat        string
	StandardFormat string
}

func Formats() map[string]Format {
	data := map[string]Format{
		"uint16": {
			XFormat:        "uint16",
			StandardType:   "string",
			StandardFormat: "int32",
		},
		"uint32": {
			XFormat:        "uint32",
			StandardType:   "string",
			StandardFormat: "int64",
		},
	}
	for k, v := range data {
		if k != v.XFormat {
			panic(fmt.Sprintf("mismatch [%s] [%s]", k, v.XFormat))
		}
	}
	return data
}
