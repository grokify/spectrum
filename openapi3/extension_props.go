package openapi3

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

type ExtensionPropsParent interface{}

// GetOperationExtensionPropStringOrEmpty converts extension prop value from `json.RawMessage` to `string`.
func GetOperationExtensionPropStringOrEmpty(op oas3.Operation, key string) string {
	str, err := GetExtensionPropString(op.ExtensionProps, key)
	if err != nil {
		return ""
	}
	return str
}

// GetExtensionPropStringOrEmpty converts extension prop value from `json.RawMessage` to `string`.
func GetExtensionPropStringOrEmpty(xprops oas3.ExtensionProps, key string) string {
	str, err := GetExtensionPropString(xprops, key)
	if err != nil {
		return ""
	}
	return str
}

// GetExtensionPropString converts extension prop value from `json.RawMessage` to `string`.
func GetExtensionPropString(xprops oas3.ExtensionProps, key string) (string, error) {
	iface, ok := xprops.Extensions[key]
	if !ok {
		return "", fmt.Errorf("extension prop key [%s] not found", key)
	}
	// Important to use %s instead of %v.
	return strings.Trim(fmt.Sprintf("%s", iface), "\""), nil
}
