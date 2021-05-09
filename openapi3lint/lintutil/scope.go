package lintutil

import (
	"fmt"
	"strings"
)

const (
	ScopeOperation     = "operation"
	ScopeSpecification = "specification"
)

var mapStringScope = map[string]string{
	"operation":     ScopeOperation,
	"oper":          ScopeOperation,
	"op":            ScopeOperation,
	"specification": ScopeSpecification,
	"spec":          ScopeSpecification,
	"":              ScopeSpecification,
}

func ParseScope(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if scope, ok := mapStringScope[s]; ok {
		return scope, nil
	}
	return "", fmt.Errorf("unknown scope [%s]", s)
}

func ScopeMatch(wantScope, tryScope string) bool {
	wantScopeCanonical, err := ParseScope(wantScope)
	if err != nil {
		return false
	}
	tryScopeCanonical, err := ParseScope(tryScope)
	if err != nil {
		return false
	}
	if wantScopeCanonical != tryScopeCanonical {
		return false
	}
	return true
}
