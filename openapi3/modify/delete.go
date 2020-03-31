package modify

import (
	"net/http"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/net/urlutil"
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
		/*for _, pathItem := range spec.Paths {
			PathItemDeleteOperationID(pathItem, opID)
		}*/
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
	for urlpath, pathItem := range spec.Paths {
		if delThis(urlpath, http.MethodConnect, pathItem.Connect) {
			pathItem.Connect = nil
		}
		if delThis(urlpath, http.MethodDelete, pathItem.Delete) {
			pathItem.Delete = nil
		}
		if delThis(urlpath, http.MethodGet, pathItem.Get) {
			pathItem.Get = nil
		}
		if delThis(urlpath, http.MethodHead, pathItem.Head) {
			pathItem.Head = nil
		}
		if delThis(urlpath, http.MethodPatch, pathItem.Patch) {
			pathItem.Patch = nil
		}
		if delThis(urlpath, http.MethodPost, pathItem.Post) {
			pathItem.Post = nil
		}
		if delThis(urlpath, http.MethodPut, pathItem.Put) {
			pathItem.Put = nil
		}
		if delThis(urlpath, http.MethodTrace, pathItem.Trace) {
			pathItem.Trace = nil
		}
	}
}

func PathItemDeleteOperationID(pathItem *oas3.PathItem, opID string) {
	if pathItem.Connect != nil && pathItem.Connect.OperationID == opID {
		pathItem.Connect = nil
	}
	if pathItem.Delete != nil && pathItem.Delete.OperationID == opID {
		pathItem.Delete = nil
	}
	if pathItem.Get != nil && pathItem.Get.OperationID == opID {
		pathItem.Get = nil
	}
	if pathItem.Head != nil && pathItem.Head.OperationID == opID {
		pathItem.Head = nil
	}
	if pathItem.Patch != nil && pathItem.Patch.OperationID == opID {
		pathItem.Patch = nil
	}
	if pathItem.Post != nil && pathItem.Post.OperationID == opID {
		pathItem.Post = nil
	}
	if pathItem.Put != nil && pathItem.Put.OperationID == opID {
		pathItem.Put = nil
	}
	if pathItem.Trace != nil && pathItem.Trace.OperationID == opID {
		pathItem.Trace = nil
	}
}
