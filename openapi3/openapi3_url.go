package openapi3

import (
	"strings"

	"github.com/grokify/mogo/net/urlutil"
)

func BuildAPIURLOAS(specServerURL, overrideServerURL, specPath string) string {
	overrideServerURL = strings.TrimSpace(overrideServerURL)
	specServerURL = strings.TrimSpace(specServerURL)
	specPath = strings.TrimSpace(specPath)
	serverURL := specServerURL
	if len(overrideServerURL) > 0 {
		serverURL = overrideServerURL
	}
	return urlutil.CondenseURI(
		strings.Join([]string{serverURL, specPath}, "/"))
}
