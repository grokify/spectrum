package swagger2

import (
	"net/http"
	"strings"
)

func CopyEndpointsByTag(tag string, specOld, specNew Specification) (Specification, error) {
	var err error
	for url, path := range specOld.Paths {
		if path.Delete != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodDelete, *path.Delete, url, path, tag, specOld, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Get != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodGet, *path.Get, url, path, tag, specOld, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Head != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodHead, *path.Head, url, path, tag, specOld, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Options != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodOptions, *path.Options, url, path, tag, specOld, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Patch != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodPatch, *path.Patch, url, path, tag, specOld, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Post != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodPost, *path.Post, url, path, tag, specOld, specNew)
			if err != nil {
				return specNew, err
			}
		}
		if path.Put != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodPut, *path.Put, url, path, tag, specOld, specNew)
			if err != nil {
				return specNew, err
			}
		}
	}
	return specNew, nil
}

func copyOrIgnoreEndpoint(method string, endpoint Endpoint, url string, path Path, wantTag string, specOld, specNew Specification) (Specification, error) {
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
	pathNew, ok := specNew.Paths[url]
	if !ok {
		pathNew = Path{}
	}
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
