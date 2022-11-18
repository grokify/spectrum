package openapi3

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

const (
	XTagGroups       = "x-tag-groups"
	XThrottlingGroup = "x-throttling-group"
)

type ExtensionPropsParent interface{}

// ExtensionPropStringOrEmpty converts extension prop value from `json.RawMessage` to `string`.
func (om *OperationMore) ExtensionPropStringOrEmpty(key string) string {
	if om.Operation == nil {
		return ""
	}
	str, err := GetExtensionPropString(om.Operation.ExtensionProps, key)
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
