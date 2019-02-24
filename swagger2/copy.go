package swagger2

import (
	"fmt"
	"net/http"
	"strings"
)

func CopyEndpointsByTag(tag string, specOld, specNew Specification) (Specification, error) {
	var err error
	for url, path := range specOld.Paths {
		if path.Get != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodGet, *path.Get, url, path, tag, specOld, specNew)
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
		if path.Delete != nil {
			specNew, err = copyOrIgnoreEndpoint(http.MethodDelete, *path.Delete, url, path, tag, specOld, specNew)
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

// EndpointCount returns a count of the endpoints for a specification.
func EndpointCount(spec Specification) int {
	endpoints := map[string]int{}
	for url, path := range spec.Paths {
		url = strings.TrimSpace(url)
		if path.Get != nil && !path.Get.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodGet, url)] = 1
		}
		if path.Patch != nil && !path.Patch.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodPatch, url)] = 1
		}
		if path.Post != nil && !path.Post.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodPost, url)] = 1
		}
		if path.Put != nil && !path.Put.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodPut, url)] = 1
		}
		if path.Delete != nil && !path.Delete.IsEmpty() {
			endpoints[fmt.Sprintf("%s %s", http.MethodDelete, url)] = 1
		}
	}
	return len(endpoints)
}
