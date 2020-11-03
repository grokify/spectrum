package openapi3

import (
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
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
