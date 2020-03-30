package modify

import (
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
		pm := PathMeta{
			Current: url}
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
