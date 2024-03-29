package postman2

import (
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
)

const DefaultMediaTypePreferences string = `multipart/form-data,application/json,application/x-www-form-urlencoded,application/xml,text/plain`

func DefaultMediaTypePreferencesSlice() []string {
	return strings.Split(DefaultMediaTypePreferences, ",")
}

func AppendPostmanHeaderValueLower(headers []Header, headerName string, options, preferenceOrder []string) ([]Header, string) {
	headerName = strings.TrimSpace(headerName)
	headerValue := stringsutil.SliceChooseOnePreferredLowerTrimSpace(options, preferenceOrder)
	if len(headerName) > 0 && len(headerValue) > 0 {
		headers = append(headers, Header{
			Key:   headerName,
			Value: headerValue})
	}
	return headers, headerValue
}

func AddOperationReqResMediaTypeHeaders(headers []Header, operation *oas3.Operation, spec *openapi3.Spec, reqPreferences []string, resPreferences []string) ([]Header, string, string, error) {
	om := openapi3.OperationMore{Operation: operation}
	reqMediaTypes, err := om.RequestMediaTypes(spec)
	if err != nil {
		return []Header{}, "", "", err
	}
	headers, reqMediaType := AppendPostmanHeaderValueLower(
		headers,
		httputilmore.HeaderContentType,
		reqMediaTypes,
		reqPreferences,
	)
	headers, resMediaType := AppendPostmanHeaderValueLower(
		headers,
		httputilmore.HeaderAccept,
		om.ResponseMediaTypes(),
		resPreferences,
	)
	return headers, reqMediaType, resMediaType, nil
}
