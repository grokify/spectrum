package openapi3

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
)

func ContentToSchemaRefMap(content oas3.Content) map[string]string {
	mss := map[string]string{}
	// type oas3.Content map[string]*MediaType
	for ct, mediaType := range content {
		if mediaType.Schema != nil && len(mediaType.Schema.Ref) > 0 {
			mss[ct] = mediaType.Schema.Ref
		} else {
			mss[ct] = ""
		}
	}
	return mss
}
