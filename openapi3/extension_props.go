package openapi3

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

// GetExtensionPropAsString converts extension prop value from `json.RawMessage` to `string`.
func GetExtensionPropAsStringOrEmpty(xprops oas3.ExtensionProps, key string) string {
	str, err := GetExtensionPropAsString(xprops, key)
	if err != nil {
		return ""
	}
	return str
}

// GetExtensionPropAsString converts extension prop value from `json.RawMessage` to `string`.
func GetExtensionPropAsString(xprops oas3.ExtensionProps, key string) (string, error) {
	iface, ok := xprops.Extensions[key]
	if !ok {
		return "", fmt.Errorf("extension prop key [%s] not found", key)
	}
	// Important to use %s instead of %v.
	return strings.Trim(fmt.Sprintf("%s", iface), "\""), nil
}
