package openapi2

import (
	"net/http"
	"sort"

	"github.com/grokify/spectrum"
)

type SpecMore struct {
	Spec *Spec
}

func (sm *SpecMore) Meta() *spectrum.SpecMeta {
	meta := spectrum.NewSpecMeta()
	if sm.Spec == nil {
		return meta
	}
	meta.Names = sm.Names()
	meta.Inflate()
	return meta
}

func (sm *SpecMore) Names() spectrum.SpecMetaNames {
	out := spectrum.SpecMetaNames{
		Endpoints: sm.Endpoints(),
		Models:    sm.ModelNames(),
		Paths:     sm.PathNames(),
	}
	return out
}

func (sm *SpecMore) Endpoints() []string {
	var out []string
	if sm.Spec == nil {
		return out
	}
	for k, pathItem := range sm.Spec.Paths {
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

func (sm *SpecMore) ModelNames() []string {
	var out []string
	for k := range sm.Spec.Definitions {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func (sm *SpecMore) PathNames() []string {
	var out []string
	for k := range sm.Spec.Paths {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
