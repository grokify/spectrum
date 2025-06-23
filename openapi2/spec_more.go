package openapi2

import (
	"net/http"
	"sort"

	"github.com/grokify/spectrum"
)

type SpecMore struct {
	Spec *Spec
}

func (m *SpecMore) Meta() *spectrum.SpecMeta {
	meta := spectrum.NewSpecMeta()
	if m.Spec == nil {
		return meta
	}
	meta.Names = m.Names()
	meta.Inflate()
	return meta
}

func (more *SpecMore) Names() spectrum.SpecMetaNames {
	out := spectrum.SpecMetaNames{
		Endpoints: more.Endpoints(),
		Models:    more.ModelNames(),
		Paths:     more.PathNames(),
	}
	return out
}

func (more *SpecMore) Endpoints() []string {
	var out []string
	if more.Spec == nil {
		return out
	}
	for k, pathItem := range more.Spec.Paths {
		if pathItem.Delete != nil {
			out = append(out, spectrum.EndpointString(http.MethodDelete, k))
		}
		if pathItem.Get != nil {
			out = append(out, spectrum.EndpointString(http.MethodGet, k))
		}
		if pathItem.Head != nil {
			out = append(out, spectrum.EndpointString(http.MethodHead, k))
		}
		if pathItem.Options != nil {
			out = append(out, spectrum.EndpointString(http.MethodOptions, k))
		}
		if pathItem.Patch != nil {
			out = append(out, spectrum.EndpointString(http.MethodPatch, k))
		}
		if pathItem.Post != nil {
			out = append(out, spectrum.EndpointString(http.MethodPost, k))
		}
		if pathItem.Put != nil {
			out = append(out, spectrum.EndpointString(http.MethodPut, k))
		}
	}
	sort.Strings(out)
	return out
}

func (more *SpecMore) ModelNames() []string {
	var out []string
	for k := range more.Spec.Definitions {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func (more *SpecMore) PathNames() []string {
	var out []string
	for k := range more.Spec.Paths {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
