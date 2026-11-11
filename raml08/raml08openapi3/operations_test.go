package raml08openapi3

import (
	"errors"
	"slices"
	"testing"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/spectrum/openapi3"
)

const petstoreFile = "testdata/petstore.raml08.json"

func operationsByPathMethod(t *testing.T, oms openapi3.OperationMores) map[string]*oas3.Operation {
	t.Helper()
	ops := map[string]*oas3.Operation{}
	for _, om := range oms {
		key := om.Method + " " + om.Path
		if _, ok := ops[key]; ok {
			t.Fatalf("duplicate operation [%s]", key)
		}
		ops[key] = om.Operation
	}
	return ops
}

func TestReadFileOperations(t *testing.T) {
	oms, err := ReadFileOperations(petstoreFile)
	if err != nil {
		t.Fatal(err)
	}
	ops := operationsByPathMethod(t, *oms)
	tests := []struct {
		key         string
		summary     string
		description string
		tags        []string
	}{
		{"GET /pets", "List pets", "Returns all pets.", []string{"Pets"}},
		{"POST /pets", "Create pet", "", []string{"Pets"}},
		{"GET /pets/{petId}", "Get pet", "Returns one pet.", []string{"Pet"}},
		{"DELETE /pets/{petId}", "", "", []string{"Pet"}},
		{"GET /stores", "List stores", "", nil},
	}
	if len(ops) != len(tests) {
		t.Errorf("got %d operations, want %d", len(ops), len(tests))
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			op, ok := ops[tt.key]
			if !ok {
				t.Fatalf("operation [%s] not found", tt.key)
			}
			if op.Summary != tt.summary {
				t.Errorf("summary = %q, want %q", op.Summary, tt.summary)
			}
			if op.Description != tt.description {
				t.Errorf("description = %q, want %q", op.Description, tt.description)
			}
			if !slices.Equal(op.Tags, tt.tags) {
				t.Errorf("tags = %v, want %v", op.Tags, tt.tags)
			}
		})
	}
}

func TestOperationTagsAreNotShared(t *testing.T) {
	oms, err := operationMoresFromPathItem("/pets", map[string]any{
		"displayName": "Pets",
		"get":         map[string]any{},
		"post":        map[string]any{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(oms) != 2 {
		t.Fatalf("got %d operations, want 2", len(oms))
	}
	oms[0].Operation.Tags[0] = "Changed"
	if got := oms[1].Operation.Tags[0]; got != "Pets" {
		t.Errorf("second operation's tag = %q after editing the first, want %q", got, "Pets")
	}
}

func TestOperationMoresFromPathItemErrors(t *testing.T) {
	tests := []struct {
		name     string
		pathItem map[string]any
		want     error
	}{
		{"resource displayName not string", map[string]any{"displayName": 7, "get": map[string]any{}}, ErrRAMLDispNameNotString},
		{"method not object", map[string]any{"get": "List pets"}, ErrRAMLMethodNotObject},
		{"method displayName not string", map[string]any{"get": map[string]any{"displayName": true}}, ErrRAMLDispNameNotString},
		{"method description not string", map[string]any{"get": map[string]any{"description": 1}}, ErrRAMLDescriptionNotString},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := operationMoresFromPathItem("/pets", tt.pathItem)
			if !errors.Is(err, tt.want) {
				t.Errorf("error = %v, want %v", err, tt.want)
			}
		})
	}
}
