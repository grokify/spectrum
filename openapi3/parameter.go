package openapi3

import (
	"net/url"
	"regexp"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/encoding/jsonpointer"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/stringsutil"
)

// ParameterNames covers path template parameters, operation parameter names/referecnes, compoents keys/names
// Parameter path names only. This is useful for viewing and modifying parameter template naems.
type ParameterNames struct {
	Components maputil.MapStringSlice
	Operations maputil.MapStringSlice
	Paths      maputil.MapStringSlice
}

func (pn *ParameterNames) Names() []string {
	names := map[string]int{}
	for name := range pn.Components {
		names[name]++
	}
	for name := range pn.Operations {
		names[name]++
	}
	for name := range pn.Paths {
		names[name]++
	}
	return maputil.StringKeys(names, nil)
}

func NewParameterNames() ParameterNames {
	return ParameterNames{
		Components: maputil.MapStringSlice{},
		Operations: maputil.MapStringSlice{},
		Paths:      maputil.MapStringSlice{},
	}
}

// ParameterPathNames returns a set of parameter names. Parameter names exist in (1) path URLs,
// (2) operation paramters and (3) spec component parameters.
func (sm *SpecMore) ParameterPathNames() ParameterNames {
	if sm.Spec == nil {
		return NewParameterNames()
	}
	return ParameterNames{
		Components: sm.ParamPathNamesComponents(),
		Operations: sm.ParamPathNamesOperations(),
		Paths:      sm.ParamPathNamesPaths()}
}

func (sm *SpecMore) ParamPathNamesComponents() map[string][]string {
	names := url.Values{}
	if sm.Spec == nil {
		return names
	}
	for paramKey, paramRef := range sm.Spec.Components.Parameters {
		if paramRef.Value == nil {
			continue
		}
		if strings.ToLower(strings.TrimSpace(paramRef.Value.In)) != InPath {
			continue
		}
		if len(paramRef.Value.Name) > 0 {
			jpath := jsonpointer.PointerSubEscapeAll(`#/components/parameters/%s/name`, paramKey)
			names.Add(paramRef.Value.Name, jpath)
		}
	}
	return names
}

func (sm *SpecMore) ParamPathNamesOperations() map[string][]string {
	names := url.Values{}
	if sm.Spec == nil {
		return names
	}
	VisitOperations(sm.Spec, func(opPath, opMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		for i, paramRef := range op.Parameters {
			if paramRef.Value == nil {
				continue
			}
			if strings.ToLower(strings.TrimSpace(paramRef.Value.In)) != InPath {
				continue
			}
			jpath := jsonpointer.PointerSubEscapeAll(`#/paths/%s/%s/parameters/%d/name`, opPath, opMethod, i)
			names.Add(paramRef.Value.Name, jpath)
		}
	})
	return names
}

func (sm *SpecMore) ParamPathNamesPaths() map[string][]string {
	names := url.Values{}
	if sm.Spec == nil {
		return names
	}
	for pathURL := range sm.Spec.Paths {
		m := PathParams(pathURL)
		if len(m) > 0 {
			jpath := jsonpointer.PointerSubEscapeAll(`#/paths/%s`, pathURL)
			for _, paramName := range m {
				names.Add(paramName, jpath)
			}
		}
	}
	return names
}

var rxParens = regexp.MustCompile(`{([^}{}]+)}`)

func ParsePathParametersParens(urlPath string) []string {
	paramNames := []string{}
	m := rxParens.FindAllStringSubmatch(urlPath, -1)
	if len(m) == 0 {
		return paramNames
	}
	for _, n := range m {
		if len(n) == 2 {
			varName := strings.TrimSpace(n[1])
			paramNames = append(paramNames, varName)
		}
	}
	if len(paramNames) > 0 {
		paramNames = stringsutil.SliceCondenseSpace(paramNames, true, false)
	}
	return paramNames
}
