package openapi3

import (
	"errors"
	"slices"
	"strings"
	"testing"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

const audienceSpec = `
openapi: 3.1.0
info: {title: Orders, version: "1.0"}
servers: [{url: https://api.example.com}]
tags:
  - {name: Orders}
  - {name: Admin}
  - {name: Partners}
paths:
  /orders:
    get:
      operationId: list-orders
      tags: [Orders]
      x-audience: public
      security: [{apiKey: []}]
      parameters: [{$ref: "#/components/parameters/Limit"}]
      responses:
        "200":
          description: Orders
          content:
            application/json:
              schema: {$ref: "#/components/schemas/OrderList"}
        default: {$ref: "#/components/responses/Error"}
    post:
      operationId: create-order
      tags: [Orders]
      x-audience: [partner, public]
      requestBody: {$ref: "#/components/requestBodies/NewOrder"}
      responses:
        "201": {description: Created}
  /orders/bulk-sync:
    x-audience: partner
    post:
      operationId: bulk-sync
      tags: [Partners]
      responses:
        "202": {description: Accepted}
  /admin/reconcile:
    post:
      operationId: reconcile
      tags: [Admin]
      security: [{session: []}]
      responses:
        "200":
          description: Report
          content:
            application/json:
              schema: {$ref: "#/components/schemas/Report"}
  /admin/flags:
    get:
      operationId: flags
      tags: [Admin]
      x-audience: public
      x-internal: true
      responses:
        "200": {description: Flags}
  /health:
    get:
      operationId: health
      x-internal: false
      responses:
        "200": {description: OK}
components:
  parameters:
    Limit: {name: limit, in: query, schema: {type: integer}}
  requestBodies:
    NewOrder:
      content:
        application/json:
          schema: {$ref: "#/components/schemas/Order"}
  responses:
    Error:
      description: Error
      headers:
        X-Request-Id: {$ref: "#/components/headers/RequestId"}
      content:
        application/json:
          schema: {$ref: "#/components/schemas/Problem"}
  headers:
    RequestId: {schema: {type: string}}
  securitySchemes:
    apiKey: {type: http, scheme: bearer}
    session: {type: apiKey, in: cookie, name: sid}
  schemas:
    OrderList:
      type: object
      properties:
        items: {type: array, items: {$ref: "#/components/schemas/Order"}}
    Order:
      allOf:
        - $ref: "#/components/schemas/Base"
        - type: object
          properties:
            payment: {oneOf: [{$ref: "#/components/schemas/Card"}, {$ref: "#/components/schemas/Invoice"}]}
            meta: {type: object, additionalProperties: {$ref: "#/components/schemas/Meta"}}
    Base: {type: object, properties: {id: {type: string}}}
    Card: {type: object}
    Invoice: {type: object}
    Meta: {type: string}
    Problem: {type: object}
    Report: {type: object}
`

func readAudienceSpec(t *testing.T) *SpecMore {
	t.Helper()
	spec, err := Parse([]byte(audienceSpec))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return &SpecMore{Spec: spec}
}

func exportAudience(t *testing.T, sm *SpecMore, reader, def Audience, opts *FilterOptions) *Spec {
	t.Helper()
	keep, err := AudienceFilter(reader, def)
	if err != nil {
		t.Fatal(err)
	}
	out, err := sm.ExportByFilter(keep, opts)
	if err != nil {
		t.Fatalf("ExportByFilter(%s): %v", reader, err)
	}
	return out
}

func operationIDs(spec *Spec) []string {
	var ids []string
	for _, item := range spec.Paths.Map() {
		for _, op := range item.Operations() {
			ids = append(ids, op.OperationID)
		}
	}
	slices.Sort(ids)
	return ids
}

func componentNames[M ~map[string]V, V any](m M) []string {
	return sortedKeys(m)
}

func TestExportByFilterAudiences(t *testing.T) {
	sm := readAudienceSpec(t)
	tests := []struct {
		reader, def Audience
		want        []string
	}{
		{AudiencePublic, AudienceInternal, []string{"create-order", "list-orders"}},
		{AudiencePartner, AudienceInternal, []string{"bulk-sync", "create-order", "list-orders"}},
		{AudienceInternal, AudienceInternal, []string{"bulk-sync", "create-order", "flags", "health", "list-orders", "reconcile"}},
		// A spec that marks only what is internal: unmarked is public.
		{AudiencePublic, AudiencePublic, []string{"create-order", "health", "list-orders", "reconcile"}},
	}
	for _, tt := range tests {
		got := operationIDs(exportAudience(t, sm, tt.reader, tt.def, nil))
		if !slices.Equal(got, tt.want) {
			t.Errorf("reader %s, default %s: got %v, want %v", tt.reader, tt.def, got, tt.want)
		}
	}
}

func TestExportByFilterPrunesComponents(t *testing.T) {
	sm := readAudienceSpec(t)
	out := exportAudience(t, sm, AudiencePublic, AudienceInternal, nil)

	if got, want := componentNames(out.Components.Schemas), []string{"Base", "Card", "Invoice", "Meta", "Order", "OrderList", "Problem"}; !slices.Equal(got, want) {
		t.Errorf("schemas: got %v, want %v", got, want)
	}
	if got := componentNames(out.Components.Parameters); !slices.Equal(got, []string{"Limit"}) {
		t.Errorf("parameters: got %v", got)
	}
	if got := componentNames(out.Components.RequestBodies); !slices.Equal(got, []string{"NewOrder"}) {
		t.Errorf("request bodies: got %v", got)
	}
	if got := componentNames(out.Components.Responses); !slices.Equal(got, []string{"Error"}) {
		t.Errorf("responses: got %v", got)
	}
	if got := componentNames(out.Components.Headers); !slices.Equal(got, []string{"RequestId"}) {
		t.Errorf("headers: got %v", got)
	}
	if got := componentNames(out.Components.SecuritySchemes); !slices.Equal(got, []string{"apiKey"}) {
		t.Errorf("security schemes: got %v", got)
	}
	var tags []string
	for _, tg := range out.Tags {
		tags = append(tags, tg.Name)
	}
	if !slices.Equal(tags, []string{"Orders"}) {
		t.Errorf("tags: got %v", tags)
	}
	if out.Paths.Value("/admin/reconcile") != nil {
		t.Error("an internal path is in the public spec")
	}
}

func TestExportByFilterRemovesExtensionsAndKeepsSource(t *testing.T) {
	sm := readAudienceSpec(t)
	out := exportAudience(t, sm, AudiencePartner, AudienceInternal, &FilterOptions{RemoveExtensions: []string{XAudience, XInternal}})

	b, err := out.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), XAudience) {
		t.Errorf("x-audience left in the output: %s", b)
	}
	src, err := sm.Spec.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(src), XAudience) || !strings.Contains(string(src), "reconcile") {
		t.Error("the source spec was changed")
	}
}

func TestExportByFilterDeterministic(t *testing.T) {
	sm := readAudienceSpec(t)
	var first []byte
	for i := range 5 {
		b, err := exportAudience(t, sm, AudiencePublic, AudienceInternal, nil).MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			first = b
		} else if string(b) != string(first) {
			t.Fatalf("run %d differs:\n%s\n%s", i, first, b)
		}
	}
}

func TestExportByFilterMissingRef(t *testing.T) {
	spec, err := Parse([]byte(`
openapi: 3.1.0
info: {title: T, version: "1"}
paths:
  /a:
    get:
      x-audience: public
      responses:
        "200":
          description: A
          content:
            application/json:
              schema: {$ref: "#/components/schemas/Missing"}
`))
	if err != nil {
		t.Fatal(err)
	}
	keep, err := AudienceFilter(AudiencePublic, AudienceInternal)
	if err != nil {
		t.Fatal(err)
	}
	_, err = (&SpecMore{Spec: spec}).ExportByFilter(keep, nil)
	if err == nil {
		t.Fatal("want an error for a reference to a missing schema")
	}
}

func TestAudienceFilterErrors(t *testing.T) {
	if _, err := AudienceFilter("everyone", AudienceInternal); !errors.Is(err, ErrUnknownAudience) {
		t.Errorf("unknown reader: got %v", err)
	}
	keep, err := AudienceFilter(AudiencePublic, AudienceInternal)
	if err != nil {
		t.Fatal(err)
	}
	for _, ext := range []map[string]any{
		{XAudience: "external-public"},
		{XAudience: []any{"public", 3}},
		{XAudience: []any{}},
		{XInternal: "yes"},
	} {
		if _, err := keep("/x", "get", nil, &oas3.Operation{Extensions: ext}); err == nil {
			t.Errorf("%v: want an error", ext)
		}
	}
}

func TestParseAudience(t *testing.T) {
	a, err := ParseAudience(" Public ")
	if err != nil || a != AudiencePublic {
		t.Errorf("got %q, %v", a, err)
	}
	if !AudiencePublic.VisibleTo(AudienceInternal) || AudiencePartner.VisibleTo(AudiencePublic) {
		t.Error("VisibleTo order is wrong")
	}
}
