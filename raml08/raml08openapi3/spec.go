package raml08openapi3

import (
	"errors"
	"fmt"

	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3edit"
)

var ErrOperationMissing = errors.New("operation is not set")

// NewSpec returns an OpenAPI 3 spec holding the operations read from a RAML v0.8 file by
// `ReadFileOperations`, with the given API title and version. It returns an error when an
// operation is not set or when two operations share a path and method.
func NewSpec(oms openapi3.OperationMores, apiTitle, apiVersion string) (*openapi3.Spec, error) {
	spec := openapi3.NewSpec("", apiTitle, apiVersion)
	for _, om := range oms {
		if om.Operation == nil {
			return nil, fmt.Errorf("%w: path [%s] method [%s]", ErrOperationMissing, om.Path, om.Method)
		}
		ope := openapi3edit.OperationEdit{OperationMore: om}
		added, err := ope.AddToSpec(spec, false)
		if err != nil {
			return nil, err
		} else if !added {
			return nil, fmt.Errorf("duplicate operation: path [%s] method [%s]", om.Path, om.Method)
		}
	}
	return spec, nil
}
