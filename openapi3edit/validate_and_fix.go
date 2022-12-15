package openapi3edit

import (
	"fmt"
	"reflect"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/net/httputilmore"
	"github.com/grokify/spectrum/openapi3"
)

/*
func OperationHasParameter(paramNameWant string, op *oas3.Operation) bool {
	paramNameWantLc := strings.ToLower(strings.TrimSpace(paramNameWant))
	for _, paramRef := range op.Parameters {
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
*/

// SortParameters sorts parameters according to an input name list.
// This used to resort parameters inline with path path parameters
// so they line up properly when rendered.
func SortParameters(sourceParams oas3.Parameters, sortOrder []string) oas3.Parameters {
	sortedParams := oas3.Parameters{}
	addedNames := map[string]int{}
	oldParamsMap := map[string]*oas3.ParameterRef{}
	for _, pRef := range sourceParams {
		if pRef.Value != nil {
			oldParamsMap[pRef.Value.Name] = pRef
		}
	}
	for _, sortedName := range sortOrder {
		if pRef, ok := oldParamsMap[sortedName]; ok {
			sortedParams = append(sortedParams, pRef)
			addedNames[sortedName] = 1
		}
	}
	for _, pRef := range sourceParams {
		if pRef.Value == nil {
			sortedParams = append(sortedParams, pRef)
		} else if _, ok := addedNames[pRef.Value.Name]; !ok {
			sortedParams = append(sortedParams, pRef)
		}
	}
	if len(sortedParams) != len(sourceParams) {
		panic("E_SORT_LEN_MISMATCH")
	}
	return sortedParams
}

// ValidateFixOperationPathParameters will add missing path parameters
// and re-sort parameters so required path parameters are on top and
// sorted by their position in the path.
func (se *SpecEdit) ValidateFixOperationPathParameters(fix bool) ([]*openapi3.OperationMeta, error) {
	errorOperations := []*openapi3.OperationMeta{}
	if se.SpecMore.Spec == nil {
		return errorOperations, openapi3.ErrSpecNotSet
	}
	openapi3.VisitOperations(
		se.SpecMore.Spec,
		func(path, method string, op *oas3.Operation) {
			if op == nil {
				return
			}
			varNamesPath := openapi3.ParsePathParametersParens(path)
			if len(varNamesPath) == 0 {
				return
			}
			varNamesMissing := []string{}
			fixed := false
			opm := openapi3.OperationMore{Operation: op}
			for _, varName := range varNamesPath {
				if !opm.HasParameter(varName) {
					if fix {
						newParamRef := &oas3.ParameterRef{
							Value: &oas3.Parameter{
								Name:     varName,
								In:       "path",
								Required: true,
								Schema: &oas3.SchemaRef{
									Value: &oas3.Schema{
										Type: "string",
									},
								},
							},
						}
						if op.Parameters == nil {
							params := make(oas3.Parameters, 0, 1)
							params = append(params, newParamRef)
							op.Parameters = params
						} else {
							op.Parameters = append(op.Parameters, newParamRef)
						}
						fixed = true
					} else {
						varNamesMissing = append(varNamesMissing, varName)
					}
				}
			}
			if fixed {
				op.Parameters = SortParameters(op.Parameters, varNamesPath)
			}
			if len(varNamesMissing) > 0 {
				om := openapi3.OperationToMeta(path, method, op, []string{})
				om.MetaNotes = append(om.MetaNotes,
					fmt.Sprintf("E_OP_MISSING_PATH_PARAMETER PARAM_NAMES[%s]",
						strings.Join(varNamesMissing, ",")))
				errorOperations = append(errorOperations, om)
			}
		},
	)

	if len(errorOperations) > 0 {
		return errorOperations, fmt.Errorf("E_NUM_VALIDATION_ERRORS [%v]", len(errorOperations))
	}
	return errorOperations, nil
}

// OperationsRequestBodyMove moves `requestBody` `$ref` to the operation
// which appears to be supported by more tools.
func (se *SpecEdit) OperationsRequestBodyMove(move bool) ([]*openapi3.OperationMeta, error) {
	errorOperations := []*openapi3.OperationMeta{}
	if se.SpecMore.Spec == nil {
		return errorOperations, openapi3.ErrSpecNotSet
	}
	//specMore := openapi3.SpecMore{Spec: spec}
	openapi3.VisitOperations(
		se.SpecMore.Spec,
		func(path, method string, op *oas3.Operation) {
			if op == nil || op.RequestBody == nil {
				return
			}
			if len(op.RequestBody.Ref) > 0 {
				if move {
					requestBodyRef := se.SpecMore.ComponentRequestBody(op.RequestBody.Ref)
					if requestBodyRef != nil {
						op.RequestBody = requestBodyRef
					}
				} else {
					om := openapi3.OperationToMeta(path, method, op, []string{})
					om.MetaNotes = append(om.MetaNotes,
						fmt.Sprintf("E_REQUEST_BODY_DEFINITION REF[%s]", op.RequestBody.Ref))
					errorOperations = append(errorOperations, om)
				}
			}
		},
	)
	if len(errorOperations) > 0 {
		return errorOperations, fmt.Errorf("E_NUM_VALIDATION_ERRORS [%v]", len(errorOperations))
	}
	return errorOperations, nil
}

// ValidateFixOperationResponseTypes looks for `application/json` responses
// with response schema types that are not `array` or `object`. If the responses
// is a string or integer, it will reset the response mime type to `text/plain`.
func (se *SpecEdit) ValidateFixOperationResponseTypes(fix bool) ([]*openapi3.OperationMeta, error) {
	errorOperations := []*openapi3.OperationMeta{}
	if se.SpecMore.Spec == nil {
		return errorOperations, openapi3.ErrSpecNotSet
	}
	openapi3.VisitOperations(
		se.SpecMore.Spec,
		func(path, method string, op *oas3.Operation) {
			if op == nil {
				return
			}
			for _, resRef := range op.Responses {
				if resRef == nil || resRef.Value == nil {
					continue
				}
				response := resRef
				for mediaTypeOrig, mtRef := range response.Value.Content {
					mediaType := strings.ToLower(strings.TrimSpace(mediaTypeOrig))
					if mediaType == httputilmore.ContentTypeAppJSON {
						schemaRef := mtRef.Schema
						if len(schemaRef.Ref) == 0 {
							schema := schemaRef.Value
							schemaType := schema.Type
							if fix && (schemaType == openapi3.TypeString || schemaType == openapi3.TypeInteger) {
								delete(response.Value.Content, mediaTypeOrig)
								if mtRefTry, ok := response.Value.Content[httputilmore.ContentTypeTextPlain]; ok {
									if !reflect.DeepEqual(mtRef, mtRefTry) {
										om := openapi3.OperationToMeta(path, method, op, []string{})
										om.MetaNotes = append(om.MetaNotes,
											fmt.Sprintf("E_BAD_MIME_TYPE_AND_SCHEMA_COLLISION MT[%s] type[%s]", mediaType, schemaType))
										errorOperations = append(errorOperations, om)
									}
								} else {
									response.Value.Content[httputilmore.ContentTypeTextPlain] = mtRef
								}
							} else if schemaType != openapi3.TypeObject && schemaType != openapi3.TypeArray {
								om := openapi3.OperationToMeta(path, method, op, []string{})
								om.MetaNotes = append(om.MetaNotes,
									fmt.Sprintf("E_BAD_MIME_TYPE_AND_SCHEMA MT[%s] type[%s]", mediaType, schemaType))
								errorOperations = append(errorOperations, om)
							}
						}
					}
				}
			}
		},
	)

	if len(errorOperations) > 0 {
		return errorOperations, fmt.Errorf("E_NUM_VALIDATION_ERRORS [%v]", len(errorOperations))
	}
	return errorOperations, nil
}
