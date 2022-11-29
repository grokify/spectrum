package openapi3

import (
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

// OperationMeta is used to hold additional information for a spec operation.
type OperationMeta struct {
	OperationID          string
	DocsDescription      string `json:"docsDescription,omitempty"`
	DocsURL              string `json:"docsURL,omitempty"`
	Method               string
	Path                 string
	SecurityScopes       []string `json:"securityScopes,omitempty"`
	Summary              string
	Tags                 []string
	MetaNotes            []string `json:"metaNotes,omitempty"`
	XThrottlingGroup     string   `json:"x-throttlingGroup,omitempty"`
	RequestBodySchemaRef string   `json:"requestBodySchemaRef"`
}

func (om *OperationMeta) TrimSpace() {
	om.OperationID = strings.TrimSpace(om.OperationID)
	om.DocsURL = strings.TrimSpace(om.DocsURL)
	om.DocsDescription = strings.TrimSpace(om.DocsDescription)
	om.SecurityScopes = stringsutil.SliceCondenseSpace(om.SecurityScopes, true, false)
	om.Tags = stringsutil.SliceCondenseSpace(om.Tags, true, false)
	om.XThrottlingGroup = strings.TrimSpace(om.XThrottlingGroup)
}

func (om *OperationMeta) PathMethod() string {
	return PathMethod(om.Path, om.Method)
}

/*
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
*/

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

/*
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
*/
