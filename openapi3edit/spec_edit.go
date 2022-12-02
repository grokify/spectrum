package openapi3edit

import "github.com/grokify/spectrum/openapi3"

type SpecEdit struct {
	SpecMore openapi3.SpecMore
}

func (se *SpecEdit) ExtensionSet(key string, val any) {
	se.SpecMore.Spec.ExtensionProps.Extensions[key] = val
}

func (se *SpecEdit) SpecSet(spec *openapi3.Spec) {
	se.SpecMore = openapi3.SpecMore{Spec: spec}
}
