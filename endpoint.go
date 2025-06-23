package spectrum

import "strings"

func EndpointString(method, path string) string {
	method = strings.ToUpper(strings.TrimSpace(method))
	path = strings.TrimSpace(path)
	return strings.Join([]string{method, path}, " ")
}
