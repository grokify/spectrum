package raml08openapi3

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/net/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/spectrum/openapi3"
)

// ReadFileOperations reads a RAML v0.8 file and returns a set of `openapi3edit.OperationMore` structs.
// The properties `path`, `method`, `summary`, `description` are populated. OpenAPI `summary` is populated
// by the `displayName` property. Currently, this reads a JSON formatted file into a map[string]interface.
// This is useful after converting a RAML v0.8 spec using https://github.com/daviemakz/oas-raml-converter-cli.
func ReadFileOperations(filename string) (*openapi3.OperationMores, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	msa := map[string]any{}
	err = json.Unmarshal(bytes, &msa)
	if err != nil {
		return nil, err
	}
	// omSet := &openapi3.OperationMoreSet{
	//	OperationMores: []openapi3.OperationMore{}}
	oms := &openapi3.OperationMores{}
	err = msaPaths("", msa, oms)
	return oms, err
}

const (
	RAMLKeyDescription = "description"
	RAMLKeyDisplayName = "displayName"
)

var (
	ErrOperationMoreSetMissing  = errors.New("required parameter operationSetMore is empty")
	ErrRAMLDescriptionNotString = errors.New("format for RAML description is not string")
	ErrRAMLDispNameNotString    = errors.New("format for RAML description is not string")
)

// msaPaths is a recursive function. Use "" for the basePath for RAML root. The structure of RAML
// appears to be that properties for metadata and sub-paths are co-mingled when sub-paths starting
// with a slash `/`. This method walks the operation segements and collects the following information
// using the `openapi3edit.OperationMore` struct: `path`, `method`, `summary`, `description`. OpenAPI
// `summary` is populated by the `displayName` property.
func msaPaths(basePath string, msa map[string]any, oms *openapi3.OperationMores) error {
	if len(msa) == 0 {
		return nil
	} else if oms == nil {
		return ErrOperationMoreSetMissing
	} else if oms == nil {
		oms = &openapi3.OperationMores{}
	}
	basePath = strings.TrimSpace(basePath)
	if len(basePath) > 0 {
		// only do if not at root.
		pathOms, err := operationMoresFromPathItem(basePath, msa)
		if err != nil {
			return err
		} else if len(pathOms) > 0 {
			*oms = append(*oms, pathOms...)
		}
	}

	for k, msaAny := range msa {
		k = strings.TrimSpace(k)
		if len(k) == 0 || strings.Index(k, "/") != 0 {
			continue
		}
		childMSA, ok := msaAny.(map[string]any)
		if !ok {
			return fmt.Errorf("value is not map[string]any for key [%s]", k)
		}
		childAbsPath := urlutil.JoinAbsolute(basePath, k)
		err := msaPaths(childAbsPath, childMSA, oms)
		if err != nil {
			return err
		}
	}
	return nil
}

func operationMoresFromPathItem(opPath string, opPathItem map[string]any) ([]openapi3.OperationMore, error) {
	oms := []openapi3.OperationMore{}
	for k, valAny := range opPathItem {
		// check if current `opPathItem`` property is an HTTP method, and add operations if so.
		methodCanonical, err := httputilmore.ParseHTTPMethod(k)
		if err != nil { // err means not known HTTP Method
			continue
		}
		om := openapi3.OperationMore{
			Path:      opPath,
			Method:    string(methodCanonical),
			Operation: &oas3.Operation{}}
		opMSA := valAny.(map[string]any)
		if descAny, ok := opMSA[RAMLKeyDescription]; ok {
			if descStr, ok := descAny.(string); ok {
				om.Operation.Description = strings.TrimSpace(descStr)
			} else {
				return oms, ErrRAMLDescriptionNotString
			}
		}
		if dispNameAny, ok := opMSA[RAMLKeyDisplayName]; ok {
			if dispNameStr, ok := dispNameAny.(string); ok {
				om.Operation.Summary = strings.TrimSpace(dispNameStr)
			} else {
				return oms, ErrRAMLDispNameNotString
			}
		}
		oms = append(oms, om)
	}
	return oms, nil
}
