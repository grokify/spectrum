package openapi3

import (
	"testing"
)

var pathVarsToGenericTests = []struct {
	v    string
	want string
}{
	{"/user/{userId}", "/user/{}"},
	{"/user/{userId}/email/{emailId}", "/user/{}/email/{}"},
}

// TestDMYHM2ParseTime ensures timeutil.DateDMYHM2 is parsed to GMT timezone.
func TestPathVarsToGeneric(t *testing.T) {
	for _, tt := range pathVarsToGenericTests {
		got := PathVarsToGeneric(tt.v)
		if got != tt.want {
			t.Errorf("modify.PathVarsToGeneric(\"%v\") Mismatch: want [%v], got [%v]", tt.v, tt.want, got)
		}
	}
}
