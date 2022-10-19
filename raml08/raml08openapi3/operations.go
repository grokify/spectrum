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
	"github.com/grokify/spectrum/openapi3edit"
)

// ReadFileOperations reads a RAML v0.8 file and returns a set of `openapi3edit.OperationMore` structs.
// The properties `path`, `method`, `summary`, `description` are populated. OpenAPI `summary` is populated
// by the `displayName` property. Currently, this reads a JSON formatted file into a map[string]interface.
func ReadFileOperations(filename string) (*openapi3edit.OperationMoreSet, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	msa := map[string]any{}
	err = json.Unmarshal(bytes, &msa)
	if err != nil {
		return nil, err
	}
	omSet := &openapi3edit.OperationMoreSet{
		OperationMores: []openapi3edit.OperationMore{}}
	err = msaPaths("", msa, omSet)
	return omSet, err
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
func msaPaths(basePath string, msa map[string]any, omSet *openapi3edit.OperationMoreSet) error {
	if omSet == nil {
		return ErrOperationMoreSetMissing
	} else if msa == nil {
		return nil
	}
	if omSet.OperationMores == nil {
		omSet.OperationMores = []openapi3edit.OperationMore{}
	}
	basePath = strings.TrimSpace(basePath)
	if len(basePath) > 0 {
		// only do if not at root.
		pathOms, err := operationMoresFromPathItem(basePath, msa)
		if err != nil {
			return err
		} else if len(pathOms) > 0 {
			omSet.OperationMores = append(omSet.OperationMores, pathOms...)
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
		err := msaPaths(childAbsPath, childMSA, omSet)
		if err != nil {
			return err
		}
	}
	return nil
}

func operationMoresFromPathItem(opPath string, opPathItem map[string]any) ([]openapi3edit.OperationMore, error) {
	oms := []openapi3edit.OperationMore{}
	for k, valAny := range opPathItem {
		// check if current `opPathItem`` property is an HTTP method, and add operations if so.
		_, err := httputilmore.ParseHTTPMethod(k)
		if err != nil { // err means not known HTTP Method
			continue
		}
		om := openapi3edit.OperationMore{
			Path:      opPath,
			Method:    strings.ToUpper(strings.TrimSpace(k)),
			Operation: &oas3.Operation{}}
		opMSA := valAny.(map[string]any)
		if descAny, ok := opMSA[RAMLKeyDescription]; ok {
			if descStr, ok := descAny.(string); ok {
				om.Operation.Description = descStr
			} else {
				return oms, ErrRAMLDescriptionNotString
			}
		}
		if dispNameAny, ok := opMSA[RAMLKeyDisplayName]; ok {
			if dispNameStr, ok := dispNameAny.(string); ok {
				om.Operation.Summary = dispNameStr
			} else {
				return oms, ErrRAMLDispNameNotString
			}
		}
		oms = append(oms, om)
	}
	return oms, nil
}
