package modify

import (
	"errors"
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/type/stringsutil"
)

const SecuritySchemeDefaultNameApiKeyAuth = "ApiKeyAuth"

// AddAPIKey adds an API Key definition to the spec.
// https://swagger.io/docs/specification/authentication/api-keys/
func SecuritySchemeAddDefinitionApiKey(spec *oas3.Swagger, schemeName, location, name string) error {
	schemeName = strings.TrimSpace(schemeName)
	location = strings.TrimSpace(location)
	name = strings.TrimSpace(name)
	if len(schemeName) == 0 {
		schemeName = SecuritySchemeDefaultNameApiKeyAuth
	}
	if len(location) == 0 {
		return errors.New("API Key Security Scheme Location cannot be empty. Must be one of: [\"header\", \"query\", \"cookie\"]")
	} else if location != "header" && location != "query" && location != "cookie" {
		return fmt.Errorf("API Key Security Scheme Invalid Location [%s], must be one of: [\"header\", \"query\", \"cookie\"]", location)
	}
	if len(name) == 0 {
		return errors.New("API Key Security Scheme name cannot be empty.")
	}
	if spec.Components.SecuritySchemes == nil {
		spec.Components.SecuritySchemes = map[string]*oas3.SecuritySchemeRef{}
	}
	spec.Components.SecuritySchemes[schemeName] = &oas3.SecuritySchemeRef{
		Value: &oas3.SecurityScheme{
			Type: "apiKey",
			In:   location,
			Name: name,
		},
	}
	return nil
}

func SecuritySchemeAddDefinitionOperations(spec *oas3.Swagger, tags []string, keyName string) {
	keyName = strings.TrimSpace(keyName)
	if len(keyName) == 0 {
		keyName = SecuritySchemeDefaultNameApiKeyAuth
	}
	tagsMap := map[string]int{}
	for _, tagName := range tags {
		tagName = strings.TrimSpace(tagName)
		tagsMap[tagName] = 1
	}
	VisitOperations(spec, func(op *oas3.Operation) {
		if op == nil || !SliceIntersectionExists(tagsMap, op.Tags) {
			return
		}
		if op.Security == nil {
			op.Security = &oas3.SecurityRequirements{}
		}
		secreq := oas3.SecurityRequirement{}
		secreq[keyName] = []string{}
		*op.Security = append(*op.Security, secreq)
	})
}

func SliceIntersection(haystack map[string]int, needles []string, unique bool) []string {
	if unique {
		needles = stringsutil.SliceCondenseSpace(needles, true, false)
	}
	matches := []string{}
	for _, needle := range needles {
		if _, ok := haystack[needle]; ok {
			matches = append(matches, needle)
		}
	}
	return matches
}

func SliceIntersectionExists(haystack map[string]int, needles []string) bool {
	for _, needle := range needles {
		if _, ok := haystack[needle]; ok {
			return true
		}
	}
	return false
}

// RemoveOperationsSecurity removes the security property
// for all operations. It is useful when building a spec
// to get individual specs to validate before setting the
// correct security property.
func RemoveOperationsSecurity(spec *oas3.Swagger) {
	VisitOperations(spec, func(op *oas3.Operation) {
		if op == nil {
			return
		}
		op.Security = &oas3.SecurityRequirements{}
	})
}
