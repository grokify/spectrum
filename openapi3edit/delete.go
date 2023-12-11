package openapi3edit

import (
	"net/http"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/spectrum/openapi3"
)

func (se *SpecEdit) DeleteProperties(md openapi3.SpecMetadata) {
	if se.SpecMore.Spec == nil {
		return
	}
	spec := se.SpecMore.Spec
	for _, opID := range md.OperationIDs {
		se.DeleteOperations(
			func(urlpath, method string, op *oas3.Operation) bool {
				if op != nil && op.OperationID == opID {
					return true
				}
				return false
			})
	}
	for _, epDel := range md.Endpoints {
		se.DeleteOperations(
			func(urlpath, method string, op *oas3.Operation) bool {
				if op == nil {
					return false
				}
				if epDel == urlutil.EndpointString(urlpath, method, false) ||
					epDel == urlutil.EndpointString(urlpath, method, true) {
					return true
				}
				return false
			})
	}
	for _, schemaNameDel := range md.SchemaNames {
		for schemaNameTry := range spec.Components.Schemas {
			if schemaNameDel == schemaNameTry {
				delete(spec.Components.Schemas, schemaNameTry)
			}
		}
	}
}

func (se *SpecEdit) DeleteOperations(delThis func(urlpath, method string, op *oas3.Operation) bool) {
	if se.SpecMore.Spec == nil {
		return
	}
	// newPaths := oas3.Paths{} // getkin v0.121.0 to v0.122.0
	newPaths := oas3.NewPaths()

	pathsMap := se.SpecMore.Spec.Paths.Map()
	for urlpath, pathItem := range pathsMap {
		// for urlpath, pathItem := range se.SpecMore.Spec.Paths { // getkin v0.121.0 to v0.122.0
		newPathItem := oas3.PathItem{
			// ExtensionProps: pathItem.ExtensionProps,
			Extensions:  pathItem.Extensions,
			Ref:         pathItem.Ref,
			Summary:     pathItem.Summary,
			Description: pathItem.Description,
			Servers:     pathItem.Servers,
			Parameters:  pathItem.Parameters}
		if pathItem.Connect != nil && !delThis(urlpath, http.MethodConnect, pathItem.Connect) {
			newPathItem.Connect = pathItem.Connect
		}
		if pathItem.Delete != nil && !delThis(urlpath, http.MethodDelete, pathItem.Delete) {
			newPathItem.Delete = pathItem.Delete
		}
		if pathItem.Get != nil && !delThis(urlpath, http.MethodGet, pathItem.Get) {
			newPathItem.Get = pathItem.Get
		}
		if pathItem.Head != nil && !delThis(urlpath, http.MethodHead, pathItem.Head) {
			newPathItem.Head = pathItem.Head
		}
		if pathItem.Options != nil && !delThis(urlpath, http.MethodOptions, pathItem.Options) {
			newPathItem.Options = pathItem.Options
		}
		if pathItem.Patch != nil && !delThis(urlpath, http.MethodPatch, pathItem.Patch) {
			newPathItem.Patch = pathItem.Patch
		}
		if pathItem.Post != nil && !delThis(urlpath, http.MethodPost, pathItem.Post) {
			newPathItem.Post = pathItem.Post
		}
		if pathItem.Put != nil && !delThis(urlpath, http.MethodPut, pathItem.Put) {
			newPathItem.Put = pathItem.Put
		}
		if pathItem.Trace != nil && !delThis(urlpath, http.MethodTrace, pathItem.Trace) {
			newPathItem.Trace = pathItem.Trace
		}
		if openapi3.PathItemHasEndpoints(&newPathItem) {
			newPaths.Set(urlpath, &newPathItem)
			// newPaths[urlpath] = &newPathItem // getkin v0.121.0 to v0.122.0
		}
	}
	se.SpecMore.Spec.Paths = newPaths
}
