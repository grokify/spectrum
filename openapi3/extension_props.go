package openapi3

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

// GetExtensionPropAsString converts extension prop value from `json.RawMessage` to `string`.
func GetExtensionPropAsStringOrEmpty(xprops oas3.ExtensionProps, key string) string {
	valIface, ok := xprops.Extensions[key]
	if !ok {
		return ""
	}
	return strings.Trim(fmt.Sprintf("%s", valIface), "\"")
}

// GetExtensionPropAsString converts extension prop value from `json.RawMessage` to `string`.
func GetExtensionPropAsString(xprops oas3.ExtensionProps, key string) (string, error) {
	valIface, ok := xprops.Extensions[key]
	if !ok {
		return "", fmt.Errorf("extension prop key [%s] not found", key)
	}
	return strings.Trim(fmt.Sprintf("%s", valIface), "\""), nil
}
