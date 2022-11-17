package openapi3

import (
	"encoding/json"
	"net/url"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/type/maputil"
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
	schemaRefs := url.Values{}
	if om.Operation == nil {
		return schemaRefs
	}
	op := om.Operation
	for _, paramRef := range op.Parameters {
		if paramRef == nil {
			continue
		}
		if len(paramRef.Ref) > 0 {
			schemaRefs.Add(paramRef.Ref, LocationParameter)
		}
		if paramRef.Value == nil {
			continue
		}
		if len(paramRef.Value.Schema.Ref) > 0 {
			schemaRefs.Add(paramRef.Value.Schema.Ref, LocationParameter)
		}
		if paramRef.Value.Schema.Value != nil && paramRef.Value.Schema.Value.Items != nil {
			if len(paramRef.Value.Schema.Value.Items.Ref) > 0 {
				schemaRefs.Add(paramRef.Value.Schema.Value.Items.Ref, LocationParameter)
			}
		}
	}
	if op.RequestBody != nil {
		if len(op.RequestBody.Ref) > 0 {
			schemaRefs.Add(op.RequestBody.Ref, LocationRequest)
		}
		if op.RequestBody.Value != nil {
			for _, mediaType := range op.RequestBody.Value.Content {
				if mediaType.Schema == nil {
					continue
				}
				if len(strings.TrimSpace(mediaType.Schema.Ref)) > 0 {
					schemaRefs.Add(mediaType.Schema.Ref, LocationRequest)
				}
			}
		}
	}
	for _, respRef := range op.Responses {
		if respRef == nil {
			continue
		}
		if len(respRef.Ref) > 0 {
			schemaRefs.Add(respRef.Ref, LocationResponse)
		}
		if respRef.Value == nil {
			continue
		}
		for _, mediaType := range respRef.Value.Content {
			if mediaType.Schema == nil {
				continue
			}
			if len(strings.TrimSpace(mediaType.Schema.Ref)) > 0 {
				schemaRefs.Add(mediaType.Schema.Ref, LocationResponse)
			}
		}
	}
	return maputil.MapStringSliceCondenseSpace(schemaRefs, true, true)
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
