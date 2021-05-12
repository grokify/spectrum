package openapi3

import (
	"testing"
)

var schemaPointerExpandTests = []struct {
	prefix     string
	schemaName string
	want       string
}{
	{"", "FooBar", "#/components/schemas/FooBar"},
	{"spec.json", "FooBar", "spec.json#/components/schemas/FooBar"},
	{"spec.json", "#/components/schemas/FooBar", "spec.json#/components/schemas/FooBar"},
}

// TestDMYHM2ParseTime ensures timeutil.DateDMYHM2 is parsed to GMT timezone.
func TestSchemaPointerExpand(t *testing.T) {
	for _, tt := range schemaPointerExpandTests {
		got := SchemaPointerExpand(tt.prefix, tt.schemaName)
		if got != tt.want {
			t.Errorf("openapi3.SchemaPointerExpand(\"%s\",\"%s\") Mismatch: want [%v], got [%v]",
				tt.prefix, tt.schemaName, tt.want, got)
		}
	}
}
