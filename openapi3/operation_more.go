package openapi3

import (
	"encoding/json"
	"net/url"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	LocationParameter = "parameter"
	LocationRequest   = "request"
	LocationResponse  = "response"
)

type OperationMore struct {
	Operation *oas3.Operation
}

/*
func MediaTypesToSlice(typesMap map[string]*oas3.MediaType) []string {
	slice := []string{}
	for thisType := range typesMap {
		slice = append(slides, thisType)
	}
	return slice
}
*/
// RequestMediaTypes returns a sorted slice of request media types.
func (om *OperationMore) RequestMediaTypes() []string {
	op := om.Operation
	mediaTypes := []string{}
	if op == nil {
		return mediaTypes
	}
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

// ResponseMediaTypes returns a sorted slice of response media types.
func (om *OperationMore) ResponseMediaTypes() []string {
	op := om.Operation
	mediaTypes := []string{}
	if op == nil {
		return mediaTypes
	}
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

// JSONPointers returns a `map[string][]string` where the keys
// are JSON pointers and the value slice is a slice of locations.
func (om *OperationMore) JSONPointers() map[string][]string {
	op := om.Operation
	schemaRefs := url.Values{}
	if op == nil {
		return schemaRefs
	}
	if op == nil {
		return schemaRefs
	}
	for _, paramRef := range op.Parameters {
		if paramRef == nil {
			continue
		}
		if len(paramRef.Ref) > 0 {
			schemaRefs.Add(paramRef.Ref, LocationParameter)
		}
	}
	if op.RequestBody != nil && len(op.RequestBody.Ref) > 0 {
		schemaRefs.Add(op.RequestBody.Ref, LocationParameter)
	}
	for _, respRef := range op.Responses {
		if respRef == nil {
			continue
		}
		if len(respRef.Ref) > 0 {
			schemaRefs.Add(respRef.Ref, LocationResponse)
		}
	}
	return schemaRefs
}

// SecurityScopes retrieves a flat list of security scopes for an operation.
func (om *OperationMore) SecurityScopes(fullyQualified bool) []string {
	op := om.Operation
	if op == nil {
		return []string{}
	}
	securityScopes := []string{}
	if op == nil || op.Security == nil {
		return securityScopes
	}
	seqReqRaw := SecurityRequirementsToRaw(*op.Security)
	for _, secReq := range seqReqRaw {
		for secSchemeName, scopes := range secReq {
			if fullyQualified {
				secSchemeNameTrimmed := strings.TrimSpace(secSchemeName)
				for _, scope := range scopes {
					scope = strings.TrimSpace(scope)
					if len(scope) > 0 {
						securityScopes = append(securityScopes,
							secSchemeNameTrimmed+"."+scope)
					}
				}
			} else {
				securityScopes = append(securityScopes, scopes...)
			}
		}
	}
	return stringsutil.SliceCondenseSpace(securityScopes, true, false)
}

// SecurityRequirementsToRaw returns a raw SecurityRequirements slice
// to be used for iterating over elements.
func SecurityRequirementsToRaw(secReqs oas3.SecurityRequirements) []map[string][]string {
	bytes, err := json.Marshal(secReqs)
	if err != nil {
		panic(err)
	}
	raw := []map[string][]string{}
	err = json.Unmarshal(bytes, &raw)
	if err != nil {
		panic(err)
	}
	return raw
}
