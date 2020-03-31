package modify

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/net/urlutil"
)

type SpecPaths struct {
	Servers openapi3.Servers
	Paths   []PathMeta
}

func InspectPaths(spec *openapi3.Swagger) SpecPaths {
	specPaths := SpecPaths{
		Servers: spec.Servers,
		Paths:   []PathMeta{}}
	for url := range spec.Paths {
		pm := PathMeta{Current: url}
		specPaths.Paths = append(specPaths.Paths, pm)
	}
	return specPaths
}

type PathMeta struct {
	Current string
	New     string
}

type SpecPathsModifyOpts struct {
	ServerPathExec          bool
	ServerPathNew           string
	OpPathRenameNewBase     string
	OpPathRenameNewBaseExec bool
	OpPathRenameFunc        func(string) string
	OpPathRenameFuncExec    bool
}

func SpecPathsModify(spec *oas3.Swagger, opts SpecPathsModifyOpts) error {
	if opts.ServerPathExec {
		opts.ServerPathNew = strings.TrimSpace(opts.ServerPathNew)
		for i, svr := range spec.Servers {
			newServerURL, err := urlutil.ModifyPath(svr.URL, opts.ServerPathNew)
			if err != nil {
				return err
			}
			svr.URL = newServerURL
			spec.Servers[i] = svr
		}
	}
	if opts.OpPathRenameFuncExec {
		oldPathURLs := map[string]int{}
		for oldPathURL := range spec.Paths {
			oldPathURLs[oldPathURL] = 1
		}
		for oldPathURL := range oldPathURLs {
			pathItem := spec.Paths[oldPathURL]
			newPathURL := strings.TrimSpace(opts.OpPathRenameFunc(oldPathURL))
			if len(newPathURL) > 0 && newPathURL != oldPathURL {
				spec.Paths[newPathURL] = pathItem
				delete(spec.Paths, oldPathURL)
			}
		}
	} else if opts.OpPathRenameNewBaseExec {
		opts.OpPathRenameNewBase = strings.TrimSpace(opts.OpPathRenameNewBase)
		if strings.Index(opts.OpPathRenameNewBase, "/") != 0 {
			// path needs to start with "/", even if not a root path.
			opts.OpPathRenameNewBase = "/" + opts.OpPathRenameNewBase
		}
		if len(opts.OpPathRenameNewBase) > 0 {
			oldPathURLs := map[string]int{}
			for oldPathURL := range spec.Paths {
				oldPathURLs[oldPathURL] = 1
			}
			for oldPathURL := range oldPathURLs {
				pathItem := spec.Paths[oldPathURL]
				newPathURL := urlutil.Join(opts.OpPathRenameNewBase, oldPathURL)
				spec.Paths[newPathURL] = pathItem
				delete(spec.Paths, oldPathURL)
			}
		}
	}
	return nil
}

func SpecEndpoints(spec *oas3.Swagger, generic bool) []string {
	endpoints := []string{}
	for url, pathItem := range spec.Paths {
		if generic {
			url = PathVarsToGeneric(url)
		}
		pathMethods := PathMethods(pathItem)
		for _, pathMethod := range pathMethods {
			endpoints = append(endpoints, url+" "+pathMethod)
		}
	}
	return endpoints
}

type Endpoint struct {
	Path      string
	Method    string
	Operation *oas3.Operation
}

func (ep *Endpoint) String() string {
	ep.Path = strings.TrimSpace(ep.Path)
	ep.Method = strings.TrimSpace(ep.Method)
	return ep.Path + " " + ep.Method
}

func PathEndpoints(url string, pathItem *oas3.PathItem) []Endpoint {
	pathOps := []Endpoint{}
	if pathItem.Connect != nil {
		pathOps = append(pathOps, Endpoint{Path: url,
			Operation: pathItem.Connect, Method: http.MethodConnect})
	}
	if pathItem.Delete != nil {
		pathOps = append(pathOps, Endpoint{Path: url,
			Operation: pathItem.Delete, Method: http.MethodDelete})
	}
	if pathItem.Get != nil {
		pathOps = append(pathOps, Endpoint{Path: url,
			Operation: pathItem.Get, Method: http.MethodGet})
	}
	if pathItem.Head != nil {
		pathOps = append(pathOps, Endpoint{Path: url,
			Operation: pathItem.Head, Method: http.MethodHead})
	}
	if pathItem.Options != nil {
		pathOps = append(pathOps, Endpoint{Path: url,
			Operation: pathItem.Options, Method: http.MethodOptions})
	}
	if pathItem.Patch != nil {
		pathOps = append(pathOps, Endpoint{Path: url,
			Operation: pathItem.Patch, Method: http.MethodPatch})
	}
	if pathItem.Post != nil {
		pathOps = append(pathOps, Endpoint{Path: url,
			Operation: pathItem.Post, Method: http.MethodPost})
	}
	if pathItem.Put != nil {
		pathOps = append(pathOps, Endpoint{Path: url,
			Operation: pathItem.Put, Method: http.MethodPut})
	}
	if pathItem.Trace != nil {
		pathOps = append(pathOps, Endpoint{Path: url,
			Operation: pathItem.Trace, Method: http.MethodTrace})
	}
	return pathOps
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
	gen1 := PathVarsToGeneric(path1)
	gen2 := PathVarsToGeneric(path2)
	if gen1 != gen2 {
		return false
	}
	return true
}
