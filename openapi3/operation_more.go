package openapi3

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/net/http/pathmethod"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

var ErrOperationNotSet = errors.New("operation not set")

const (
	LocationParameter = "parameter"
	LocationRequest   = "request"
	LocationResponse  = "response"
)

type OperationMore struct {
	Path      string
	Method    string
	Operation *oas3.Operation
}

type OperationMoreStringFunc func(opm *OperationMore) string

type OperationMoreStringFuncMap map[string]OperationMoreStringFunc

func (opmmap *OperationMoreStringFuncMap) Func(key string) OperationMoreStringFunc {
	opmmapIndexable := map[string]OperationMoreStringFunc(*opmmap)
	wantFunc, ok := opmmapIndexable[key]
	if !ok {
		return nil
	}
	return wantFunc
}

func (om *OperationMore) HasParameter(paramNameWant string) bool {
	paramNameWantLc := strings.ToLower(strings.TrimSpace(paramNameWant))
	for _, paramRef := range om.Operation.Parameters {
		if paramRef.Value == nil {
			continue
		}
		param := paramRef.Value
		param.Name = strings.TrimSpace(param.Name)
		paramNameTryLc := strings.ToLower(param.Name)
		if paramNameWantLc == paramNameTryLc {
			return true
		}
	}
	return false
}

func (om *OperationMore) Meta() *OperationMeta {
	return OperationToMeta(om.Path, om.Method, om.Operation, []string{})
}

func (om *OperationMore) PathMethod() string {
	return pathmethod.PathMethod(om.Path, om.Method)
}

// RequestMediaTypes returns a sorted slice of request media types.
func (om *OperationMore) RequestMediaTypes(spec *Spec) ([]string, error) {
	if om.Operation == nil {
		return []string{}, ErrOperationNotSet
	}
	op := om.Operation
	if op.RequestBody == nil {
		return []string{}, nil
	}
	if len(strings.TrimSpace(op.RequestBody.Ref)) == 0 {
		if op.RequestBody.Value != nil {
			return maputil.StringKeys(op.RequestBody.Value.Content, nil), nil
		}
		return []string{}, nil
	}
	if spec == nil {
		return []string{}, errors.New("spec needed to trace json pointer")
	}
	sm := SpecMore{Spec: spec}
	return sm.RequestBodyRefMediaTypes(op.RequestBody)
}

func (sm *SpecMore) RequestBodyRefMediaTypes(ref *oas3.RequestBodyRef) ([]string, error) {
	if sm.Spec == nil {
		return []string{}, ErrSpecNotSet
	}
	if ref == nil {
		return []string{}, nil
	}
	if len(strings.TrimSpace(ref.Ref)) > 0 {
		reqRef, err := sm.RequestBodyRef(ref.Ref)
		if err != nil {
			return []string{}, err
		}
		return sm.RequestBodyRefMediaTypes(reqRef)
	}
	if ref.Value != nil {
		return maputil.StringKeys(ref.Value.Content, nil), nil
	}
	return []string{}, nil
}

var rxJSONPointerComponentsRequestBodies = regexp.MustCompile(`^(.*?)#/components/requestBodies/(.+)$`)

func (sm *SpecMore) RequestBodyRef(componentKeyOrPointer string) (*oas3.RequestBodyRef, error) {
	if reqRef, ok := sm.Spec.Components.RequestBodies[componentKeyOrPointer]; ok {
		return reqRef, nil
	}
	if strings.Contains(componentKeyOrPointer, PointerComponentsRequestBodies) {
		m := rxJSONPointerComponentsRequestBodies.FindStringSubmatch(componentKeyOrPointer)
		if len(m) == 0 {
			return nil, errors.New("json pointer does not match request bodies")
		}
		reqBodyKey := m[2]
		if reqRef, ok := sm.Spec.Components.RequestBodies[reqBodyKey]; ok {
			return reqRef, nil
		}
	}
	return nil, errors.New("request body component not found")
}

func (om *OperationMore) RequestBodySchemaRef() []string {
	schemaRefs := []string{}
	if om.Operation == nil {
		return schemaRefs
	}
	if om.Operation.RequestBody == nil {
		return schemaRefs
	}
	if len(om.Operation.RequestBody.Ref) > 0 {
		schemaRefs = append(schemaRefs, om.Operation.RequestBody.Ref)
	}
	if om.Operation.RequestBody.Value != nil {
		ctMap := ContentToSchemaRefMap(om.Operation.RequestBody.Value.Content)
		keys := maputil.Values(ctMap)
		schemaRefs = append(schemaRefs, keys...)
	}
	return schemaRefs
}

// ResponseMediaTypes returns a sorted slice of response media types. Media type values are
// deduped against multiple response statuses.
func (om *OperationMore) ResponseMediaTypes() []string {
	if om.Operation == nil {
		return []string{}
	}
	mediaTypes := []string{}
	respsMap := om.Operation.Responses.Map() // added for getkin v0.121.0 to v0.122.0 breaking change
	for _, respRef := range respsMap {
		// for _, respRef := range om.Operation.Responses {
		if respRef.Value != nil {
			mt := maputil.StringKeys(respRef.Value.Content, nil)
			mediaTypes = append(mediaTypes, mt...)
		}
	}
	mediaTypes = slicesutil.Dedupe(mediaTypes)
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
	respsMap := op.Responses.Map()
	for _, respRef := range respsMap {
		// for _, respRef := range op.Responses { // added for getkin v0.121.0 to v0.122.0 breaking change
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
	mss := maputil.MapStringSlice(schemaRefs)
	return mss.CondenseSpace(true, true)
	// return maputil.MapStringSliceCondenseSpace(schemaRefs, true, true)
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

type OperationMores []OperationMore

/*
type OperationMoreSet struct {
	OperationMores []OperationMore
}
*/

// SummariesMap returns a `map[string]string` where the keys are the operation's
// path and method, while the values are the summaries.
func (oms *OperationMores) SummariesMap() map[string]string {
	mss := map[string]string{}
	for _, om := range *oms {
		mss[om.PathMethod()] = om.Operation.Summary
	}
	return mss
}

func OperationMoresForPath(url string, pathItem *oas3.PathItem) []OperationMore {
	pathOps := []OperationMore{}
	if pathItem == nil {
		return pathOps
	}
	if pathItem.Connect != nil {
		pathOps = append(pathOps, OperationMore{Path: url,
			Operation: pathItem.Connect, Method: http.MethodConnect})
	}
	if pathItem.Delete != nil {
		pathOps = append(pathOps, OperationMore{Path: url,
			Operation: pathItem.Delete, Method: http.MethodDelete})
	}
	if pathItem.Get != nil {
		pathOps = append(pathOps, OperationMore{Path: url,
			Operation: pathItem.Get, Method: http.MethodGet})
	}
	if pathItem.Head != nil {
		pathOps = append(pathOps, OperationMore{Path: url,
			Operation: pathItem.Head, Method: http.MethodHead})
	}
	if pathItem.Options != nil {
		pathOps = append(pathOps, OperationMore{Path: url,
			Operation: pathItem.Options, Method: http.MethodOptions})
	}
	if pathItem.Patch != nil {
		pathOps = append(pathOps, OperationMore{Path: url,
			Operation: pathItem.Patch, Method: http.MethodPatch})
	}
	if pathItem.Post != nil {
		pathOps = append(pathOps, OperationMore{Path: url,
			Operation: pathItem.Post, Method: http.MethodPost})
	}
	if pathItem.Put != nil {
		pathOps = append(pathOps, OperationMore{Path: url,
			Operation: pathItem.Put, Method: http.MethodPut})
	}
	if pathItem.Trace != nil {
		pathOps = append(pathOps, OperationMore{Path: url,
			Operation: pathItem.Trace, Method: http.MethodTrace})
	}
	return pathOps
}
