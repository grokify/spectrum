package openapi2

import (
	"net/http"
	"strings"

	"github.com/qdm12/reprint"
)

func CopyEndpointsByTag(tag string, specOld, specNew Specification) (Specification, error) {
	var err error
	for url, path := range specOld.Paths {
		if path.Delete != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodDelete, *path.Delete, url, path, tag, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Get != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodGet, *path.Get, url, path, tag, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Head != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodHead, *path.Head, url, path, tag, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Options != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodOptions, *path.Options, url, path, tag, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Patch != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodPatch, *path.Patch, url, path, tag, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Post != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodPost, *path.Post, url, path, tag, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Put != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodPut, *path.Put, url, path, tag, specNew)
			if err != nil {
				return specNew, err
			}
		}
	}
	return specNew, nil
}

func copyOrIgnoreEndpoint(method string, endpoint Endpoint, url string, path Path, wantTag string, specNew Specification) (Specification, error) {
	// TODO: copy referenced objects, e.g. schema objects, which may need `specOld`.
	wantTag = strings.TrimSpace(wantTag)
	if len(wantTag) != 0 {
		match := false
		for _, tryTag := range endpoint.Tags {
			if strings.TrimSpace(tryTag) == wantTag {
				match = true
			}
		}
		if !match {
			return specNew, nil
		}
	}
	if _, ok := specNew.Paths[url]; ok {
		return specNew, nil
	}

	pathNew := reprint.This(path).(Path) // ref: https://stackoverflow.com/a/77412997/1908967
	err := pathNew.SetEndpoint(method, endpoint)
	if err != nil {
		return specNew, err
	}
	if specNew.Paths == nil {
		specNew.Paths = map[string]Path{}
	}
	specNew.Paths[url] = pathNew
	return specNew, nil
}
