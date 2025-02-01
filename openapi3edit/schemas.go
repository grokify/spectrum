package openapi3edit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/path/filepathutil"
	"github.com/grokify/mogo/reflect/reflectutil"
	"github.com/grokify/spectrum/openapi3"
)

func (se *SpecEdit) AddSchemaDir(dir string, fileRx *regexp.Regexp) error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	}
	entries, err := osutil.ReadDirMore(dir, nil, true, true, false)
	if err != nil {
		return err
	}
	sdirs := []os.DirEntry{}
	for _, entry := range entries {
		if entry.IsDir() {
			sdirs = append(sdirs, entry)
		} else if fileRx != nil && fileRx.MatchString(entry.Name()) {
			continue
		}

		sch, err := openapi3.ReadSchemaFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}
		delete(sch.Extensions, "$schema")
		entryRoot := filepathutil.TrimExt(entry.Name())
		se.SpecMore.Spec.Components.Schemas[entryRoot] = &oas3.SchemaRef{Value: sch}
	}
	for _, sdir := range sdirs {
		err := se.AddSchemaDir(filepath.Join(dir, sdir.Name()), fileRx)
		if err != nil {
			return err
		}
	}
	return nil
}

// SchemaRefsModifyRx modifies `$ref` reference strings that match a supplied
// `*regexp.Regexp` and replaces that with a string. It was originally
// designed to convert `#schemas/` to `#components/schemas/`.
func (se *SpecEdit) SchemaRefsModifyRx(rx *regexp.Regexp, repl string) {
	if se.SpecMore.Spec == nil || rx == nil {
		return
	}
	se.SchemaRefsModify(func(s string) string {
		return rx.ReplaceAllString(s, repl)
	})
}

// SchemaRefsModify modifys schema reference JSON pointers. The xf function
// must return the entire JSON pointer.
func (se *SpecEdit) SchemaRefsModify(xf func(string) string) {
	if se.SpecMore.Spec == nil || xf == nil {
		return
	}
	spec := se.SpecMore.Spec

	for _, paramRef := range spec.Components.Parameters {
		if paramRef == nil {
			continue
		}
		// paramRef.Ref = rx.ReplaceAllString(paramRef.Ref, repl)
		paramRef.Ref = xf(paramRef.Ref)
		if paramRef.Value != nil && paramRef.Value.Schema != nil {
			//SchemaRefModifyRefs(paramRef.Value.Schema, rx, repl)
			SchemaRefModifyRefs(paramRef.Value.Schema, xf)
		}
	}

	openapi3.VisitOperations(spec, func(opPath, opMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		// Operation Parameters
		for _, paramRef := range op.Parameters {
			paramRef.Ref = xf(paramRef.Ref)
			if paramRef.Value != nil && paramRef.Value.Schema != nil {
				SchemaRefModifyRefs(paramRef.Value.Schema, xf)
			}
		}
		// Operation Requests
		if op.RequestBody != nil {
			op.RequestBody.Ref = xf(op.RequestBody.Ref)
			if op.RequestBody.Value != nil {
				for _, mediaType := range op.RequestBody.Value.Content {
					SchemaRefModifyRefs(mediaType.Schema, xf)
				}
			}
		}
		// Operation Responses
		respsMap := op.Responses.Map()
		for _, respRef := range respsMap {
			// for _, respRef := range op.Responses { // getkin v0.121.0 to v0.122.0
			respRef.Ref = xf(respRef.Ref)
			if respRef.Value == nil {
				continue
			}
			for _, mediaType := range respRef.Value.Content {
				SchemaRefModifyRefs(mediaType.Schema, xf)
			}
		}
	})

	for _, schRef := range spec.Components.Schemas {
		if schRef == nil {
			continue
		}
		if schRef.Ref == "" && schRef.Value != nil {
			// if `$ref` is populated as an Extension with mmessage type RawMessage,
			// manually convert to reference.
			if refValAny, ok := schRef.Value.Extensions["$ref"]; ok {
				switch reflectutil.NameOf(refValAny, false) {
				case "RawMessage": // json.RawMessage
					refValJRM, ok := refValAny.(json.RawMessage)
					if !ok {
						panic("type `RawMessage` does not coerce to `json.RawMessage`")
					}
					refStr := ""
					err := json.Unmarshal(refValJRM, &refStr)
					if err != nil {
						panic("cannot unmarshal `json.RawMessage` to string")
					}
					schRef.Ref = refStr
					delete(schRef.Value.Extensions, "$ref")
				case "string":
					refStr, ok := refValAny.(string)
					if ok {
						schRef.Ref = strings.TrimSpace(refStr)
						delete(schRef.Value.Extensions, "$ref")
					}
				}
			}
		}
		// schRef.Ref = rx.ReplaceAllString(schRef.Ref, repl)
		// SchemaRefModifyRefs(schRef, rx, repl)
		schRef.Ref = xf(schRef.Ref)
		SchemaRefModifyRefs(schRef, xf)
	}
}

// SchemaRefModifyRefsRx modifies Schema reference schema pointers that match
// the supplied `*regexp.Regexp` with the replacement string. It was originally
// designed to convert `#schemas/` to `#components/schemas/`.
func SchemaRefModifyRefsRx(schRef *oas3.SchemaRef, rx *regexp.Regexp, repl string) {
	SchemaRefModifyRefs(schRef, func(s string) string {
		return rx.ReplaceAllString(s, repl)
	})
}

func SchemaRefModifyRefs(schRef *oas3.SchemaRef, xf func(string) string) {
	if schRef == nil || xf == nil {
		return
	}
	schRef.Ref = xf(schRef.Ref)
	if schRef.Value == nil {
		return
	}
	for _, propSchemaRef := range schRef.Value.Properties {
		SchemaRefModifyRefs(propSchemaRef, xf)
	}
	if schRef.Value.AdditionalProperties.Schema != nil {
		SchemaRefModifyRefs(schRef.Value.AdditionalProperties.Schema, xf)
	}
	if schRef.Value.Items != nil {
		SchemaRefModifyRefs(schRef.Value.Items, xf)
	}
	if schRef.Value.Not != nil {
		SchemaRefModifyRefs(schRef.Value.Not, xf)
	}
	for _, allOf := range schRef.Value.OneOf {
		SchemaRefModifyRefs(allOf, xf)
	}
	for _, anyOf := range schRef.Value.OneOf {
		SchemaRefModifyRefs(anyOf, xf)
	}
	for _, oneOf := range schRef.Value.OneOf {
		SchemaRefModifyRefs(oneOf, xf)
	}
}

func (se *SpecEdit) SchemaSetAdditionalPropertiesTrue(pointerBase string) []string {
	mods := []string{}
	if se.SpecMore.Spec == nil {
		return mods
	}
	for schName, schRef := range se.SpecMore.Spec.Components.Schemas {
		if schRef == nil || schRef.Value == nil || !openapi3.TypesRefIs(schRef.Value.Type, openapi3.TypeObject) {
			continue
		}
		/*
			if len(schRef.Value.Properties) == 0 &&
				schRef.Value.AdditionalProperties == nil &&
				(schRef.Value.AdditionalPropertiesAllowed == nil || !*schRef.Value.AdditionalPropertiesAllowed) {
				additionalPropertiesAllowed := true
				schRef.Value.AdditionalPropertiesAllowed = &additionalPropertiesAllowed
				mods = append(mods, fmt.Sprintf("%s%s/%s", pointerBase, openapi3.PointerComponentsSchemas, schName))
			}
		*/
		if !openapi3.AdditionalPropertiesAllowed(schRef.Value.AdditionalProperties) {
			additionalPropertiesAllowed := true
			schRef.Value.AdditionalProperties.Has = &additionalPropertiesAllowed
			mods = append(mods, fmt.Sprintf("%s%s/%s", pointerBase, openapi3.PointerComponentsSchemas, schName))
		}
		for propName, propRef := range schRef.Value.Properties {
			if propRef == nil || propRef.Value == nil || !openapi3.TypesRefIs(propRef.Value.Type, openapi3.TypeObject) {
				continue
			}
			/*
				if len(propRef.Value.Properties) == 0 &&
					propRef.Value.AdditionalProperties == nil &&
					(propRef.Value.AdditionalPropertiesAllowed == nil || !*propRef.Value.AdditionalPropertiesAllowed) {
					additionalPropertiesAllowed := true
					propRef.Value.AdditionalPropertiesAllowed = &additionalPropertiesAllowed
					mods = append(mods, fmt.Sprintf("%s%s/%s/properties/%s", pointerBase, openapi3.PointerComponentsSchemas, schName, propName))
				}
			*/
			if !openapi3.AdditionalPropertiesAllowed(propRef.Value.AdditionalProperties) {
				additionalPropertiesAllowed := true
				propRef.Value.AdditionalProperties.Has = &additionalPropertiesAllowed
				mods = append(mods, fmt.Sprintf("%s%s/%s/properties/%s", pointerBase, openapi3.PointerComponentsSchemas, schName, propName))
			}
		}
	}
	return mods
}
