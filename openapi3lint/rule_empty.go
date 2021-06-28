package openapi3lint

import (
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3lint/lintutil"
)

type EmptyRule struct{}

func (rule EmptyRule) Name() string  { return "" }
func (rule EmptyRule) Scope() string { return "" }
func (rule EmptyRule) ProcessSpec(spec *openapi3.Spec, pointerBase string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}
func (rule EmptyRule) ProcessOperation(spec *openapi3.Spec, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation {
	return []lintutil.PolicyViolation{}
}
