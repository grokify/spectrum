package openapi3edit

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/stringsutil/transform"
	"github.com/grokify/spectrum/openapi3"
)

func (se *SpecEdit) paramPathNamesModifyComponents(xf func(s string) string) error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	}
	spec := se.SpecMore.Spec
	namesStart := maputil.MapStringSlice(se.SpecMore.ParamPathNamesPaths())
	namesStartSlice := maputil.StringKeys(namesStart, nil)
	// check for unique outnames
	_, _, err := transform.TransformMap(xf, namesStartSlice)
	if err != nil {
		return err
	}
	for _, paramRef := range spec.Components.Parameters {
		if paramRef == nil ||
			paramRef.Value == nil ||
			strings.ToLower(strings.TrimSpace(paramRef.Value.In)) != openapi3.InPath {
			continue
		}
		paramRef.Value.Name = xf(paramRef.Value.Name)
	}
	namesComplete := maputil.MapStringSlice(se.SpecMore.ParamPathNamesPaths())
	if len(namesStart) != len(namesComplete) {
		return fmt.Errorf("conversion mismatch: start name count [%d] end name count [%d]", len(namesStart), len(namesComplete))
	}
	return nil
}

func (se *SpecEdit) paramPathNamesModifyOperations(xf func(s string) string) (map[string]string, error) {
	if se.SpecMore.Spec == nil {
		return map[string]string{}, openapi3.ErrSpecNotSet
	}
	if xf == nil {
		return map[string]string{}, nil
	}
	namesStart := url.Values(se.SpecMore.ParamPathNamesOperations())
	namesStartSlice := maputil.StringKeys(namesStart, nil)
	// check for unique outnames
	xfMap, _, err := transform.TransformMap(xf, namesStartSlice)
	if err != nil {
		return xfMap, errorsutil.Wrap(err, "transform.TransformMap")
	}
	openapi3.VisitOperations(se.SpecMore.Spec, func(opPath, opMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		for _, paramRef := range op.Parameters {
			if paramRef.Value == nil {
				continue
			}
			if strings.ToLower(strings.TrimSpace(paramRef.Value.In)) != openapi3.InPath {
				continue
			}
			paramRef.Value.Name = xf(paramRef.Value.Name)
		}
	})
	namesComplete := url.Values(se.SpecMore.ParamPathNamesOperations())
	if len(namesStart) != len(namesComplete) {
		return xfMap, fmt.Errorf("conversion mismatch: start name count [%d] end name count [%d]", len(namesStart), len(namesComplete))
	}
	return xfMap, nil
}

func (se *SpecEdit) paramPathNamesModifyPaths(xf func(s string) string) error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	}
	spec := se.SpecMore.Spec
	pathsCountStart := len(se.SpecMore.Spec.Paths)
	namesStart := url.Values(se.SpecMore.ParamPathNamesPaths())
	namesStartSlice := maputil.StringKeys(namesStart, nil)
	// check for unique outnames
	_, _, err := transform.TransformMap(xf, namesStartSlice)
	if err != nil {
		return err
	}
	pathBefores := maputil.StringKeys(spec.Paths, nil)
	pathsMap := map[string]string{}
	for _, pathBefore := range pathBefores {
		pathItem, ok := spec.Paths[pathBefore]
		if !ok {
			panic("path not found")
		}
		pathAfter := PathTemplateParamMod(pathBefore, xf)
		if pathAfter != pathBefore {
			spec.Paths[pathAfter] = pathItem
			pathsMap[pathBefore] = pathAfter
			delete(spec.Paths, pathBefore)
		}
	}

	if !maputil.UniqueValues(pathsMap) {
		return errors.New("path strcase collisions")
	}
	pathsCountComplete := len(spec.Paths)
	if pathsCountStart != pathsCountComplete {
		return fmt.Errorf("conversion mismatch: start path count [%d] end path count [%d]", pathsCountStart, pathsCountComplete)
	}
	namesComplete := url.Values(se.SpecMore.ParamPathNamesPaths())
	if len(namesStart) != len(namesComplete) {
		return fmt.Errorf("conversion mismatch: start name count [%d] end name count [%d]", len(namesStart), len(namesComplete))
	}
	return nil
}

// ParamPathNamesModify should result in a spec that validates and performs post-modification validation.
func (se *SpecEdit) ParamPathNamesModify(xf func(string) string) (map[string]string, error) {
	// Operations must come before components.
	xfMap, err := se.paramPathNamesModifyOperations(xf)
	if err != nil {
		return xfMap, errorsutil.Wrap(err, "SpecEdit.ParamPathNamesModifyOperations")
	}
	err = se.paramPathNamesModifyComponents(xf)
	if err != nil {
		return xfMap, errorsutil.Wrap(err, "SpecEdit.ParamPathNamesModifyComponents")
	}
	err = se.paramPathNamesModifyPaths(xf)
	if err != nil {
		return xfMap, errorsutil.Wrap(err, "SpecEdit.ParamPathNamesModifyPaths")
	}
	return xfMap, se.SpecMore.Validate()
}

// PathTemplateParamMod takes a URL path with templated parameters like `{pet_id}`.`
func PathTemplateParamMod(p string, xf func(string) string) string {
	m := openapi3.RxPathParam.FindAllStringSubmatch(p, -1)
	ids0 := map[string]int{}
	for _, mi := range m {
		niOut := xf(mi[1])
		rx := regexp.MustCompile(regexp.QuoteMeta(mi[0]))
		p = rx.ReplaceAllString(p, "{"+niOut+"}")
		ids0[mi[0]]++
	}
	mverify := openapi3.RxPathParam.FindAllStringSubmatch(p, -1)
	ids1 := map[string]int{}
	for _, mi := range mverify {
		ids1[mi[0]]++
	}
	if len(ids0) != len(ids1) {
		panic("mismatch")
	}
	return p
}
