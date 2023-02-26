package openapi3edit

import "github.com/grokify/spectrum/openapi3"

type SpecEdit struct {
	SpecMore openapi3.SpecMore
}

func NewSpecEdit(spec *openapi3.Spec) SpecEdit {
	return SpecEdit{
		SpecMore: openapi3.SpecMore{Spec: spec}}
}

func (se *SpecEdit) ExtensionSet(key string, val any) {
	// se.SpecMore.Spec.ExtensionProps.Extensions[key] = val
	se.SpecMore.Spec.Extensions[key] = val
}

func (se *SpecEdit) SpecSet(spec *openapi3.Spec) {
	se.SpecMore = openapi3.SpecMore{Spec: spec}
}
