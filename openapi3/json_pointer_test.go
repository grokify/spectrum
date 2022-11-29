package openapi3

import (
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

var jsonPointerTests = []struct {
	jsonPointer string
	document    string
	path        []string
	isTopParam  bool
	isTopSchema bool
}{
	{"mydoc.yaml#/components/schemas/FooBar", "mydoc.yaml", []string{"components", "schemas", "FooBar"}, false, true},
	// {"#/components/schemas/FooBar", "", []string{"components", "schemas", "FooBar"}, false, true},
}

// TestJSONPointers ensures the `ParseJSONPointer` is working properly.
func TestJSONPointers(t *testing.T) {
	for _, tt := range jsonPointerTests {
		ptr, err := ParseJSONPointer(tt.jsonPointer)
		if err != nil {
			t.Errorf("openapi3.ParseJSONPointer(\"%s\") Error [%s]",
				tt.jsonPointer, err.Error())
		}
		if ptr.Document != tt.document {
			t.Errorf("JSONPointer.Document Mismatch: want [%v], got [%v]",
				tt.document, ptr.Document)
		}
		if !slices.Equal(ptr.Path, tt.path) {
			t.Errorf("JSONPointer.Path Mismatch: want [%v], got [%v]",
				strings.Join(tt.path, ", "), strings.Join(ptr.Path, ", "))
		}
		_, gotIsTopParam := ptr.IsTopParameter()
		if gotIsTopParam != tt.isTopParam {
			t.Errorf("JSONPointer.IsTopParameter() Mismatch: want [%v], got [%v]",
				tt.isTopParam, gotIsTopParam)
		}
		_, gotIsTopSchema := ptr.IsTopSchema()
		if gotIsTopSchema != tt.isTopSchema {
			t.Errorf("JSONPointer.IsTopSchema() Mismatch: want [%v], got [%v]",
				tt.isTopSchema, gotIsTopSchema)
		}
	}
}
