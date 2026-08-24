package openapi3

import (
	"reflect"
	"testing"
)

const diffSpec1JSON = `{
  "openapi":"3.0.0","info":{"title":"S1","version":"1.0.0"},
  "paths":{
    "/a":{"get":{"operationId":"getA","responses":{"200":{"description":"ok"}}}},
    "/b":{"get":{"operationId":"getB","responses":{"200":{"description":"ok"}}}}
  },
  "components":{"schemas":{
    "Foo":{"type":"object"},
    "Bar":{"type":"object"}
  }}
}`

// diffSpec2JSON keeps endpoint GET /a but renames its operationId (getA -> listA),
// drops GET /b, adds GET /c, and swaps schema Bar for Baz.
const diffSpec2JSON = `{
  "openapi":"3.0.0","info":{"title":"S2","version":"1.0.0"},
  "paths":{
    "/a":{"get":{"operationId":"listA","responses":{"200":{"description":"ok"}}}},
    "/c":{"get":{"operationId":"getC","responses":{"200":{"description":"ok"}}}}
  },
  "components":{"schemas":{
    "Foo":{"type":"object"},
    "Baz":{"type":"object"}
  }}
}`

func TestSpecsDiff(t *testing.T) {
	spec1, err := Parse([]byte(diffSpec1JSON))
	if err != nil {
		t.Fatalf("Parse(spec1) error: %v", err)
	}
	spec2, err := Parse([]byte(diffSpec2JSON))
	if err != nil {
		t.Fatalf("Parse(spec2) error: %v", err)
	}

	d := SpecsDiff(spec1, spec2)

	checks := []struct {
		name string
		got  []string
		want []string
	}{
		// GET /a is shared; the differing operationId does not affect the endpoint set.
		{"Endpoints OnlyInSpec1", d.OnlyInSpec1.Endpoints, []string{"/b GET"}},
		{"Endpoints OnlyInSpec2", d.OnlyInSpec2.Endpoints, []string{"/c GET"}},
		{"Endpoints Both", d.Both.Endpoints, []string{"/a GET"}},
		// getA was renamed to listA, so no operationId is shared.
		{"OperationIDs OnlyInSpec1", d.OnlyInSpec1.OperationIDs, []string{"getA", "getB"}},
		{"OperationIDs OnlyInSpec2", d.OnlyInSpec2.OperationIDs, []string{"getC", "listA"}},
		{"OperationIDs Both", d.Both.OperationIDs, []string{}},
		{"SchemaNames OnlyInSpec1", d.OnlyInSpec1.SchemaNames, []string{"Bar"}},
		{"SchemaNames OnlyInSpec2", d.OnlyInSpec2.SchemaNames, []string{"Baz"}},
		{"SchemaNames Both", d.Both.SchemaNames, []string{"Foo"}},
	}
	for _, c := range checks {
		if !reflect.DeepEqual(c.got, c.want) {
			t.Errorf("%s: want %v, got %v", c.name, c.want, c.got)
		}
	}

	if d.IsEmpty() {
		t.Error("SpecMetadataDiff.IsEmpty() = true, want false for differing specs")
	}
}

func TestSpecsDiffIdentical(t *testing.T) {
	spec, err := Parse([]byte(diffSpec1JSON))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	d := SpecsDiff(spec, spec)
	if !d.IsEmpty() {
		t.Errorf("SpecsDiff(spec, spec).IsEmpty() = false, want true\n%s", d.String())
	}
}
