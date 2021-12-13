package openapi3

import (
	"strings"

	"github.com/grokify/mogo/net/urlutil"
)

func BuildApiUrlOAS(specServerURL, overrideServerURL, specPath string) string {
	overrideServerURL = strings.TrimSpace(overrideServerURL)
	specServerURL = strings.TrimSpace(specServerURL)
	specPath = strings.TrimSpace(specPath)
	serverURL := specServerURL
	if len(overrideServerURL) > 0 {
		serverURL = overrideServerURL
	}
	fullUrl := strings.Join([]string{serverURL, specPath}, "/")
	return urlutil.CondenseUri(fullUrl)
}
