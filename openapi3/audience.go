package openapi3

import (
	"errors"
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

const (
	// XAudience names who may see an operation or path: internal, partner
	// or public, as a string or a list of strings.
	XAudience = "x-audience"
	// XInternal marks an operation or path as internal when true, as
	// Redocly, Scalar and Stoplight read it.
	XInternal = "x-internal"
)

// Audience is who may see part of an API. Audiences are ordered from
// narrowest to widest, and each includes the narrower ones: what is public
// is also shown to partners and internal readers.
type Audience string

const (
	AudienceInternal Audience = "internal"
	AudiencePartner  Audience = "partner"
	AudiencePublic   Audience = "public"
)

var audienceRanks = map[Audience]int{
	AudienceInternal: 1,
	AudiencePartner:  2,
	AudiencePublic:   3,
}

// ErrUnknownAudience is returned for an audience other than internal,
// partner or public.
var ErrUnknownAudience = errors.New("unknown audience")

// ParseAudience reads an audience name, ignoring case and surrounding space.
func ParseAudience(s string) (Audience, error) {
	a := Audience(strings.ToLower(strings.TrimSpace(s)))
	if _, ok := audienceRanks[a]; !ok {
		return "", fmt.Errorf("%w: %q", ErrUnknownAudience, s)
	}
	return a, nil
}

// VisibleTo reports whether something for audience a may be shown to
// reader: a is reader's audience or a wider one.
func (a Audience) VisibleTo(reader Audience) bool {
	ra, ok := audienceRanks[a]
	return ok && ra >= audienceRanks[reader]
}

// ExtensionAudience reads an audience from extensions: internal when
// x-internal is true, which wins over x-audience so that conflicting marks
// never widen who sees something; else x-audience, as a string or a list
// (the widest listed audience wins). ok is false when neither says
// anything.
func ExtensionAudience(ext map[string]any) (a Audience, ok bool, err error) {
	if v, has := ext[XInternal]; has {
		internal, isBool := v.(bool)
		if !isBool {
			return "", false, fmt.Errorf("%s: want a boolean, got %T", XInternal, v)
		}
		if internal {
			return AudienceInternal, true, nil
		}
	}
	if v, has := ext[XAudience]; has {
		a, err := audienceValue(v)
		if err != nil {
			return "", false, fmt.Errorf("%s: %w", XAudience, err)
		}
		return a, true, nil
	}
	return "", false, nil
}

func audienceValue(v any) (Audience, error) {
	switch t := v.(type) {
	case string:
		return ParseAudience(t)
	case []string:
		return widestAudience(t)
	case []any:
		names := make([]string, 0, len(t))
		for _, e := range t {
			s, ok := e.(string)
			if !ok {
				return "", fmt.Errorf("want a string or a list of strings, got a list with %T", e)
			}
			names = append(names, s)
		}
		return widestAudience(names)
	default:
		return "", fmt.Errorf("want a string or a list of strings, got %T", v)
	}
}

func widestAudience(names []string) (Audience, error) {
	if len(names) == 0 {
		return "", errors.New("empty list")
	}
	var widest Audience
	for _, n := range names {
		a, err := ParseAudience(n)
		if err != nil {
			return "", err
		}
		if widest == "" || audienceRanks[a] > audienceRanks[widest] {
			widest = a
		}
	}
	return widest, nil
}

// AudienceFilter keeps the operations reader may see. An operation's
// audience is its own x-audience or x-internal, else its path item's, else
// def. Pass AudienceInternal as def so an operation nobody marked is never
// published by mistake, or AudiencePublic for a spec that marks only what
// is internal with x-internal.
func AudienceFilter(reader, def Audience) (OperationFilterFunc, error) {
	if _, ok := audienceRanks[reader]; !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownAudience, reader)
	}
	if _, ok := audienceRanks[def]; !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownAudience, def)
	}
	return func(path, method string, pathItem *oas3.PathItem, op *oas3.Operation) (bool, error) {
		a, ok, err := ExtensionAudience(op.Extensions)
		if err != nil {
			return false, fmt.Errorf("%s %s: %w", strings.ToUpper(method), path, err)
		}
		if !ok && pathItem != nil {
			if a, ok, err = ExtensionAudience(pathItem.Extensions); err != nil {
				return false, fmt.Errorf("%s: %w", path, err)
			}
		}
		if !ok {
			a = def
		}
		return a.VisibleTo(reader), nil
	}, nil
}
