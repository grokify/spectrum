package taxonomy

import (
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3edit"
)

type Taxonomy struct {
	Category CategoryRef `json:"category"`
	Slug     string      `json:"slug"`
}

func (tax *Taxonomy) AddToSpec(spec *openapi3.Spec) {
	if spec != nil {
		se := openapi3edit.SpecEdit{}
		se.SpecSet(spec)
		se.ExtensionSet(XTaxonomy, tax)
	}
}

type CategoryRef struct {
	Ref string `json:"$ref"`
}
