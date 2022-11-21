package openapi3

import (
	"errors"
	"strings"

	"github.com/grokify/mogo/net/httputilmore"
)

// PathMethod returns a path-method string which can be used as a unique identifier for operations.
func PathMethod(opPath, opMethod string) string {
	opPath = strings.TrimSpace(opPath)
	opMethod = strings.ToUpper(strings.TrimSpace(opMethod))
	parts := []string{}
	if len(opPath) > 0 {
		parts = append(parts, opPath)
	}
	if len(opMethod) > 0 {
		parts = append(parts, opMethod)
	}
	return strings.Join(parts, " ")
}

type PathMethodSet struct {
	PathMethods map[string]int
}

func (pms *PathMethodSet) init() {
	if pms.PathMethods == nil {
		pms.PathMethods = map[string]int{}
	}
}

// Add adds pathmethod strings
func (pms *PathMethodSet) Add(pathmethods ...string) error {
	pms.init()
	for _, pm := range pathmethods {
		opPath, opMethod, err := ParsePathMethod(pm)
		if err != nil {
			return err
		}
		pathMethod := PathMethod(opPath, opMethod)
		pms.PathMethods[pathMethod]++
	}
	return nil
}

func (pms *PathMethodSet) Count() int {
	pms.init()
	return len(pms.PathMethods)
}

func (pms *PathMethodSet) Exists(opPath, opMethod string) bool {
	pms.init()
	pm := PathMethod(opPath, opMethod)
	_, ok := pms.PathMethods[pm]
	return ok
}

func (pms *PathMethodSet) StringExists(pathMethod string) bool {
	pms.init()
	_, ok := pms.PathMethods[pathMethod]
	return ok
}

var ErrPathMethodInvalid = errors.New("pathmethod string invalid")

func ParsePathMethod(pathmethod string) (string, string, error) {
	parts := strings.Split(pathmethod, " ")
	if len(parts) != 2 {
		return "", "", ErrPathMethodInvalid
	}
	method, err := httputilmore.ParseHTTPMethodString(parts[1])
	return strings.TrimSpace(parts[0]), method, err
}
