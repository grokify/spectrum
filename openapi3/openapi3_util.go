package openapi3

import (
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/net/urlutil"
)

/*
func MediaTypesToSlice(typesMap map[string]*oas3.MediaType) []string {
	slice := []string{}
	for thisType := range typesMap {
		slice = append(slides, thisType)
	}
	return slice
}
*/
// OperationRequestMediaTypes returns a sorted slice of
// request media types.
func OperationRequestMediaTypes(op *oas3.Operation) []string {
	mediaTypes := []string{}
	if op.RequestBody != nil {
		if op.RequestBody.Value != nil {
			for mediaType := range op.RequestBody.Value.Content {
				mediaType = strings.TrimSpace(mediaType)
				if len(mediaType) > 0 {
					mediaTypes = append(mediaTypes, mediaType)
				}
			}
		}
	}
	sort.Strings(mediaTypes)
	return mediaTypes
}

// OperationResponseMediaTypes returns a sorted slice of
// response media types.
func OperationResponseMediaTypes(op *oas3.Operation) []string {
	mediaTypes := []string{}
	for _, resp := range op.Responses {
		for mediaType := range resp.Value.Content {
			mediaType = strings.TrimSpace(mediaType)
			if len(mediaType) > 0 {
				mediaTypes = append(mediaTypes, mediaType)
			}
		}
	}
	sort.Strings(mediaTypes)
	return mediaTypes
}

// ServerURL returns the OAS3 Spec URL for the index
// specified.
func ServerURL(spec *oas3.Swagger, index int) string {
	if index+1 > len(spec.Servers) {
		return ""
	}
	server := spec.Servers[index]
	return strings.TrimSpace(server.URL)
}

// BasePath extracts the base path from a OAS URL
// which can include variables.
func BasePath(spec *oas3.Swagger) (string, error) {
	serverURL := ServerURL(spec, 0)
	if len(serverURL) == 0 {
		return "", nil
	}
	serverURLParsed, err := urlutil.ParseURLTemplate(serverURL)
	if err != nil {
		return "", err
	}
	return serverURLParsed.Path, nil
}
