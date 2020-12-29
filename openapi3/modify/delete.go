package modify

import (
	"net/http"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/net/urlutil"
)

func SpecDeleteProperties(spec *oas3.Swagger, md SpecMetadata) {
	for _, opID := range md.OperationIDs {
		SpecDeleteOperations(spec,
			func(urlpath, method string, op *oas3.Operation) bool {
				if op != nil && op.OperationID == opID {
					return true
				}
				return false
			})
	}
	for _, epDel := range md.Endpoints {
		SpecDeleteOperations(spec,
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

func SpecDeleteOperations(spec *oas3.Swagger, delThis func(urlpath, method string, op *oas3.Operation) bool) {
	newPaths := oas3.Paths{}

	for urlpath, pathItem := range spec.Paths {
		newPathItem := oas3.PathItem{
			ExtensionProps: pathItem.ExtensionProps,
			Ref:            pathItem.Ref,
			Summary:        pathItem.Summary,
			Description:    pathItem.Description,
			Servers:        pathItem.Servers,
			Parameters:     pathItem.Parameters}
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
		if PathItemHasEndpoints(&newPathItem) {
			newPaths[urlpath] = &newPathItem
		}
	}
	spec.Paths = newPaths
}

func PathItemHasEndpoints(pathItem *oas3.PathItem) bool {
	if pathItem.Connect != nil || pathItem.Delete != nil ||
		pathItem.Get != nil || pathItem.Head != nil ||
		pathItem.Patch != nil || pathItem.Post != nil ||
		pathItem.Put != nil || pathItem.Trace != nil {
		return true
	}
	return false
}
