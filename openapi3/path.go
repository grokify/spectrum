package openapi3

import (
	"net/http"
	"regexp"
	"sort"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

var RxPathParam = regexp.MustCompile(`\{([^\{\}]+)\}`)

func PathParams(p string) []string {
	params := []string{}
	m := RxPathParam.FindAllStringSubmatch(p, -1)
	for _, mi := range m {
		params = append(params, mi[1])
	}
	return params
}

func PathItemHasEndpoints(pathItem *oas3.PathItem) bool {
	if pathItem.Connect != nil || pathItem.Delete != nil ||
		pathItem.Get != nil || pathItem.Head != nil ||
		pathItem.Options != nil ||
		pathItem.Patch != nil || pathItem.Post != nil ||
		pathItem.Put != nil || pathItem.Trace != nil {
		return true
	}
	return false
}

func PathMethods(pathItem *oas3.PathItem) []string {
	methods := []string{}
	if pathItem.Connect != nil {
		methods = append(methods, http.MethodConnect)
	}
	if pathItem.Delete != nil {
		methods = append(methods, http.MethodDelete)
	}
	if pathItem.Get != nil {
		methods = append(methods, http.MethodGet)
	}
	if pathItem.Head != nil {
		methods = append(methods, http.MethodHead)
	}
	if pathItem.Options != nil {
		methods = append(methods, http.MethodOptions)
	}
	if pathItem.Patch != nil {
		methods = append(methods, http.MethodPatch)
	}
	if pathItem.Post != nil {
		methods = append(methods, http.MethodPost)
	}
	if pathItem.Put != nil {
		methods = append(methods, http.MethodPut)
	}
	if pathItem.Trace != nil {
		methods = append(methods, http.MethodTrace)
	}
	return methods
}

var rxPathVarToGeneric = regexp.MustCompile(`{[^}{]*}`)

func PathVarsToGeneric(input string) string {
	return rxPathVarToGeneric.ReplaceAllString(input, "{}")
}

func PathMatchGeneric(path1, path2 string) bool {
	return PathVarsToGeneric(path1) == PathVarsToGeneric(path2)
}

func (sm *SpecMore) PathMethods(generic bool) []string {
	endpoints := []string{}
	if sm.Spec == nil {
		return endpoints
	}
	for url, pathItem := range sm.Spec.Paths {
		if generic {
			url = PathVarsToGeneric(url)
		}
		pathMethods := PathMethods(pathItem)
		for _, pathMethod := range pathMethods {
			endpoints = append(endpoints, url+" "+pathMethod)
		}
	}
	sort.Strings(endpoints)
	return endpoints
}
