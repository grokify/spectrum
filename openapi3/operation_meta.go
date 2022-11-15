package openapi3

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

// OperationToMeta converts a path, method and operation to an `*OperationMeta`.
// The function returns `nil` if any of the items are empty.
func OperationToMeta(url, method string, op *oas3.Operation, inclTags []string) *OperationMeta {
	if url == "" || method == "" || op == nil {
		return nil
	}
	if len(inclTags) > 0 {
		inclTagsMap := map[string]int{}
		for _, inclTag := range inclTags {
			inclTagsMap[inclTag]++
		}
		haveMatch := false
		for _, opTag := range op.Tags {
			if _, ok := inclTagsMap[opTag]; ok {
				haveMatch = true
				break
			}
		}
		if !haveMatch {
			return nil
		}
	}
	return &OperationMeta{
		OperationID: strings.TrimSpace(op.OperationID),
		Summary:     strings.TrimSpace(op.Summary),
		Method:      strings.ToUpper(strings.TrimSpace(method)),
		Path:        strings.TrimSpace(url),
		Tags:        op.Tags,
		MetaNotes:   []string{}}
}

// OperationMeta is used to hold additional information
// for a spec operation.
type OperationMeta struct {
	OperationID      string
	DocsDescription  string
	DocsURL          string
	Method           string
	Path             string
	SecurityScopes   []string
	Summary          string
	Tags             []string
	MetaNotes        []string
	XThrottlingGroup string
}

func (om *OperationMeta) TrimSpace() {
	om.OperationID = strings.TrimSpace(om.OperationID)
	om.DocsURL = strings.TrimSpace(om.DocsURL)
	om.DocsDescription = strings.TrimSpace(om.DocsDescription)
	om.SecurityScopes = stringsutil.SliceCondenseSpace(om.SecurityScopes, true, false)
	om.Tags = stringsutil.SliceCondenseSpace(om.Tags, true, false)
	om.XThrottlingGroup = strings.TrimSpace(om.XThrottlingGroup)
}

func OperationSetRequestBodySchemaRef(op *oas3.Operation, mediaType string, schemaRef *oas3.SchemaRef) {
	if op.RequestBody == nil {
		op.RequestBody = &oas3.RequestBodyRef{}
	}
	if op.RequestBody.Value == nil {
		op.RequestBody.Value = &oas3.RequestBody{
			Content: oas3.NewContent()}
	}
	op.RequestBody.Value.Content[mediaType] = oas3.NewMediaType().WithSchemaRef(schemaRef)
}

/*
	op.RequestBody = &oas3.RequestBodyRef{
		Value: &oas3.RequestBody{
			Content: map[string]*oas3.MediaType{
				"application/json": &oas3.MediaType{
					Schema: &oas3.SchemaRef{
						Ref: ref,
					},
					// Example: //
				},
			},
		},
	}
*/

func OperationSetResponseBodySchemaRef(op *oas3.Operation, status, description, mediaType string, schemaRef *oas3.SchemaRef) error {
	description = strings.TrimSpace(description)
	if len(description) == 0 {
		return fmt.Errorf("no response description for operationId [%s]", op.OperationID)
	}
	if op.Responses == nil {
		op.Responses = oas3.Responses{}
	}
	status = strings.TrimSpace(status)
	mediaType = strings.ToLower(strings.TrimSpace(mediaType))
	if _, ok := op.Responses[status]; !ok || op.Responses[status] == nil {
		op.Responses[status] = &oas3.ResponseRef{}
	}
	resRef := op.Responses[status]
	if resRef.Value == nil {
		resRef.Value = &oas3.Response{}
	}
	resRef.Value.Description = &description
	if resRef.Value.Content == nil {
		resRef.Value.Content = oas3.NewContent()
	}
	resRef.Value.Content[mediaType] = oas3.NewMediaType().WithSchemaRef(schemaRef)
	return nil
}

// OperationSecurityScopes retrieves a flat list of security scopes for
// an operation.
func OperationSecurityScopes(op *oas3.Operation, fullyQualified bool) []string {
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

// PathMethod returns a path-method string which can be used as a unique identifier for operations.
func PathMethod(opPath, opMethod string) string {
	opPath = strings.TrimSpace(opPath)
	opMethod = strings.ToUpper(strings.TrimSpace(opMethod))
	parts := []string{}
	if len(opPath) > 0 {
		parts = append(parts, opPath)
	}
	if len(opMethod) > 0 {
		parts = append(parts, opMethod)
	}
	return strings.Join(parts, " ")
}
