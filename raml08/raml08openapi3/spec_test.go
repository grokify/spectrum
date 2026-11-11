package raml08openapi3

import (
	"errors"
	"slices"
	"testing"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/spectrum/openapi3"
)

func TestNewSpec(t *testing.T) {
	oms, err := ReadFileOperations(petstoreFile)
	if err != nil {
		t.Fatal(err)
	}
	spec, err := NewSpec(*oms, "Pet Store", "v1")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Info.Title != "Pet Store" || spec.Info.Version != "v1" {
		t.Errorf("info = %q %q, want %q %q", spec.Info.Title, spec.Info.Version, "Pet Store", "v1")
	}
	sm := openapi3.SpecMore{Spec: spec}
	if got := sm.OperationsCount(); got != 5 {
		t.Errorf("got %d operations, want 5", got)
	}
	op, err := sm.OperationByPathMethod("/pets/{petId}", "GET")
	if err != nil {
		t.Fatal(err)
	}
	if op == nil {
		t.Fatal("operation GET /pets/{petId} not found")
	}
	if op.Summary != "Get pet" || !slices.Equal(op.Tags, []string{"Pet"}) {
		t.Errorf("operation = %q %v, want %q %v", op.Summary, op.Tags, "Get pet", []string{"Pet"})
	}
}

func TestNewSpecErrors(t *testing.T) {
	get := func() openapi3.OperationMore {
		return openapi3.OperationMore{Path: "/pets", Method: "GET", Operation: &oas3.Operation{Summary: "List pets"}}
	}
	t.Run("duplicate operation", func(t *testing.T) {
		if _, err := NewSpec(openapi3.OperationMores{get(), get()}, "Pet Store", "v1"); err == nil {
			t.Error("error = nil, want an error for a duplicate path and method")
		}
	})
	t.Run("operation not set", func(t *testing.T) {
		om := get()
		om.Operation = nil
		if _, err := NewSpec(openapi3.OperationMores{om}, "Pet Store", "v1"); !errors.Is(err, ErrOperationMissing) {
			t.Errorf("error = %v, want %v", err, ErrOperationMissing)
		}
	})
}
