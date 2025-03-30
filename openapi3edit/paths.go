package openapi3edit

import (
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/spectrum/openapi3"
)

type Paths struct {
	*oas3.Paths
}

func (p Paths) PathKeys() []string {
	var paths []string
	pathsMap := p.Map()
	for k := range pathsMap {
		paths = append(paths, k)
	}
	sort.Strings(paths)
	return paths
}

type SpecPaths struct {
	Servers oas3.Servers
	Paths   []PathMeta
}

func InspectPaths(spec *openapi3.Spec) SpecPaths {
	specPaths := SpecPaths{
		Servers: spec.Servers,
		Paths:   []PathMeta{}}
	pathsMap := spec.Paths.Map()
	for url := range pathsMap {
		// for url := range spec.Paths {
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

func (se *SpecEdit) PathsModify(opts SpecPathsModifyOpts) error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	}
	spec := se.SpecMore.Spec
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
		pathsMap := spec.Paths.Map()
		for oldPathURL := range pathsMap {
			// for oldPathURL := range spec.Paths { // getkin v0.121.0 to v0.122.0
			oldPathURLs[oldPathURL] = 1
		}
		for oldPathURL := range oldPathURLs {
			// pathItem := spec.Paths[oldPathURL] // getkin v0.121.0 to v0.122.0
			pathItem := spec.Paths.Find(oldPathURL)
			newPathURL := strings.TrimSpace(opts.OpPathRenameFunc(oldPathURL))
			if len(newPathURL) > 0 && newPathURL != oldPathURL {
				// spec.Paths[newPathURL] = pathItem // getkin v0.121.0 to v0.122.0
				// delete(spec.Paths, oldPathURL) // getkin v0.121.0 to v0.122.0
				spec.Paths.Set(newPathURL, pathItem)
				spec.Paths.Set(oldPathURL, nil) // getkin v0.121.0 to v0.122.0
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
			oldPathsMap := spec.Paths.Map()
			for oldPathURL := range oldPathsMap {
				// for oldPathURL := range spec.Paths {
				oldPathURLs[oldPathURL] = 1
			}
			for oldPathURL := range oldPathURLs {
				// pathItem := spec.Paths[oldPathURL] // getkin v0.121.0 to v0.122.0
				pathItem := spec.Paths.Find(oldPathURL)
				newPathURL := urlutil.Join(opts.OpPathRenameNewBase, oldPathURL)
				// spec.Paths[newPathURL] = pathItem // getkin v0.121.0 to v0.122.0
				// delete(spec.Paths, oldPathURL) // getkin v0.121.0 to v0.122.0
				spec.Paths.Set(newPathURL, pathItem)
				spec.Paths.Set(oldPathURL, nil)
			}
		}
	}
	return nil
}

/*
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
*/
