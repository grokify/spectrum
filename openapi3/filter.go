package openapi3

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

// OperationFilterFunc reports whether to keep an operation. pathItem is the
// operation's path item, for filters that read path-level fields.
type OperationFilterFunc func(path, method string, pathItem *oas3.PathItem, op *oas3.Operation) (bool, error)

// FilterOptions configures ExportByFilter.
type FilterOptions struct {
	// RemoveExtensions are extensions removed from the document, its
	// path items and kept operations, such as x-audience and x-internal
	// once they have done their job.
	RemoveExtensions []string
}

// ErrRefNotFound is returned when a kept operation refers to a component
// the spec doesn't have.
var ErrRefNotFound = errors.New("reference not found")

// ExportByFilter returns a new spec with the operations keep accepts, in
// paths and webhooks, and only the components, security schemes and tags
// they use, following references through every part of a component, such
// as allOf, oneOf and additionalProperties. Path items left with no
// operations are dropped. The source spec is not changed, and the result
// is checked to have no unresolved references.
func (sm *SpecMore) ExportByFilter(keep OperationFilterFunc, opts *FilterOptions) (*Spec, error) {
	if sm.Spec == nil {
		return nil, ErrSpecNotSet
	}
	if keep == nil {
		return nil, errors.New("filter func is nil")
	}
	if opts == nil {
		opts = &FilterOptions{}
	}
	src, err := sm.Clone()
	if err != nil {
		return nil, fmt.Errorf("clone spec: %w", err)
	}

	dest := &Spec{
		OpenAPI:           src.OpenAPI,
		Info:              src.Info,
		Servers:           src.Servers,
		Security:          src.Security,
		ExternalDocs:      src.ExternalDocs,
		JSONSchemaDialect: src.JSONSchemaDialect,
		Extensions:        src.Extensions,
		Paths:             oas3.NewPaths(),
		Components:        &oas3.Components{},
	}
	if src.Components != nil {
		dest.Components.Extensions = src.Components.Extensions
	}
	var kept []*oas3.Operation
	if src.Paths != nil {
		for _, path := range sortedKeys(src.Paths.Map()) {
			item, ops, err := filterPathItem(path, src.Paths.Value(path), keep)
			if err != nil {
				return nil, err
			}
			if item != nil {
				dest.Paths.Set(path, item)
				kept = append(kept, ops...)
			}
		}
	}
	for _, name := range sortedKeys(src.Webhooks) {
		item, ops, err := filterPathItem(name, src.Webhooks[name], keep)
		if err != nil {
			return nil, err
		}
		if item != nil {
			if dest.Webhooks == nil {
				dest.Webhooks = map[string]*oas3.PathItem{}
			}
			dest.Webhooks[name] = item
			kept = append(kept, ops...)
		}
	}

	if err := copyComponents(src, dest); err != nil {
		return nil, err
	}
	if err := copySecuritySchemes(src, dest, kept); err != nil {
		return nil, err
	}
	dest.Tags = usedTags(src.Tags, kept)
	removeExtensions(dest, opts.RemoveExtensions)

	// Reload the result, which resolves every reference and fails on any
	// that points nowhere.
	b, err := dest.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshal filtered spec: %w", err)
	}
	out, err := oas3.NewLoader().LoadFromData(b)
	if err != nil {
		return nil, fmt.Errorf("load filtered spec: %w", err)
	}
	return out, nil
}

// filterPathItem returns a copy of item with the operations keep accepts,
// or nil when it accepts none.
func filterPathItem(path string, item *oas3.PathItem, keep OperationFilterFunc) (*oas3.PathItem, []*oas3.Operation, error) {
	if item == nil {
		return nil, nil, nil
	}
	ops := item.Operations()
	out := &oas3.PathItem{
		Extensions:  item.Extensions,
		Summary:     item.Summary,
		Description: item.Description,
		Servers:     item.Servers,
		Parameters:  item.Parameters,
	}
	var kept []*oas3.Operation
	for _, method := range sortedKeys(ops) {
		op := ops[method]
		ok, err := keep(path, method, item, op)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			out.SetOperation(method, op)
			kept = append(kept, op)
		}
	}
	if len(kept) == 0 {
		return nil, nil, nil
	}
	return out, kept, nil
}

// copyComponents copies into dest every component its paths and webhooks
// refer to, and every component those refer to in turn.
func copyComponents(src, dest *Spec) error {
	b, err := json.Marshal(struct {
		Paths    *oas3.Paths               `json:"paths,omitempty"`
		Webhooks map[string]*oas3.PathItem `json:"webhooks,omitempty"`
	}{dest.Paths, dest.Webhooks})
	if err != nil {
		return fmt.Errorf("marshal paths: %w", err)
	}
	queue, err := localRefs(b)
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	for len(queue) > 0 {
		ref := queue[0]
		queue = queue[1:]
		kind, name, ok := componentRef(ref)
		if !ok {
			return fmt.Errorf("%w: %s is not a component", ErrRefNotFound, ref)
		}
		key := kind + "/" + name
		if seen[key] {
			continue
		}
		seen[key] = true
		v, err := copyComponent(src.Components, dest.Components, kind, name)
		if err != nil {
			return fmt.Errorf("%s: %w", ref, err)
		}
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("marshal %s: %w", ref, err)
		}
		more, err := localRefs(b)
		if err != nil {
			return err
		}
		queue = append(queue, more...)
	}
	return nil
}

// copyComponent copies one component from src to dest and returns it.
func copyComponent(src, dest *oas3.Components, kind, name string) (any, error) {
	if src == nil {
		return nil, ErrRefNotFound
	}
	switch kind {
	case "schemas":
		return copyEntry(src.Schemas, &dest.Schemas, name)
	case "parameters":
		return copyEntry(src.Parameters, &dest.Parameters, name)
	case "headers":
		return copyEntry(src.Headers, &dest.Headers, name)
	case "requestBodies":
		return copyEntry(src.RequestBodies, &dest.RequestBodies, name)
	case "responses":
		return copyEntry(src.Responses, &dest.Responses, name)
	case "securitySchemes":
		return copyEntry(src.SecuritySchemes, &dest.SecuritySchemes, name)
	case "examples":
		return copyEntry(src.Examples, &dest.Examples, name)
	case "links":
		return copyEntry(src.Links, &dest.Links, name)
	case "callbacks":
		return copyEntry(src.Callbacks, &dest.Callbacks, name)
	default:
		return nil, fmt.Errorf("%w: unknown component kind %q", ErrRefNotFound, kind)
	}
}

func copyEntry[M ~map[string]V, V any](src M, dest *M, name string) (any, error) {
	v, ok := src[name]
	if !ok {
		return nil, ErrRefNotFound
	}
	if *dest == nil {
		*dest = M{}
	}
	(*dest)[name] = v
	return v, nil
}

// copySecuritySchemes copies the security schemes the document and the
// kept operations require.
func copySecuritySchemes(src, dest *Spec, ops []*oas3.Operation) error {
	reqs := slices.Clone(dest.Security)
	for _, op := range ops {
		if op.Security != nil {
			reqs = append(reqs, *op.Security...)
		}
	}
	for _, req := range reqs {
		for _, name := range sortedKeys(req) {
			if _, err := copyComponent(src.Components, dest.Components, "securitySchemes", name); err != nil {
				return fmt.Errorf("security scheme %q: %w", name, err)
			}
		}
	}
	return nil
}

// usedTags returns the tags the operations use, in the source's order,
// followed by any used tag the source doesn't declare.
func usedTags(tags oas3.Tags, ops []*oas3.Operation) oas3.Tags {
	used := map[string]bool{}
	for _, op := range ops {
		for _, t := range op.Tags {
			used[t] = true
		}
	}
	var out oas3.Tags
	for _, t := range tags {
		if t != nil && used[t.Name] {
			out = append(out, t)
		}
	}
	return out
}

func removeExtensions(spec *Spec, names []string) {
	if len(names) == 0 {
		return
	}
	strip := func(ext map[string]any) {
		for _, n := range names {
			delete(ext, n)
		}
	}
	strip(spec.Extensions)
	items := []*oas3.PathItem{}
	if spec.Paths != nil {
		for _, item := range spec.Paths.Map() {
			items = append(items, item)
		}
	}
	for _, item := range spec.Webhooks {
		items = append(items, item)
	}
	for _, item := range items {
		strip(item.Extensions)
		for _, op := range item.Operations() {
			strip(op.Extensions)
		}
	}
}

// localRefs returns the "#/..." references in a JSON document.
func localRefs(b []byte) ([]string, error) {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, fmt.Errorf("read references: %w", err)
	}
	var refs []string
	var walk func(any)
	walk = func(v any) {
		switch t := v.(type) {
		case map[string]any:
			if r, ok := t["$ref"].(string); ok && strings.HasPrefix(r, "#/") {
				refs = append(refs, r)
			}
			for _, k := range sortedKeys(t) {
				walk(t[k])
			}
		case []any:
			for _, e := range t {
				walk(e)
			}
		}
	}
	walk(v)
	return refs, nil
}

// componentRef splits "#/components/{kind}/{name}[/...]" into its kind and
// unescaped name.
func componentRef(ref string) (kind, name string, ok bool) {
	parts := strings.Split(strings.TrimPrefix(ref, "#/"), "/")
	if len(parts) < 3 || parts[0] != "components" {
		return "", "", false
	}
	unescape := strings.NewReplacer("~1", "/", "~0", "~")
	return parts[1], unescape.Replace(parts[2]), true
}

func sortedKeys[M ~map[string]V, V any](m M) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
