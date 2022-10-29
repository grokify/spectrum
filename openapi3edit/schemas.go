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

func SpecAddSchemaDir(spec *openapi3.Spec, dir string, fileRx *regexp.Regexp) error {
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
		spec.Components.Schemas[entryRoot] = &oas3.SchemaRef{Value: sch}
	}
	for _, sdir := range sdirs {
		err := SpecAddSchemaDir(spec, filepath.Join(dir, sdir.Name()), fileRx)
		if err != nil {
			return err
		}
	}
	return nil
}

// SpecModifySchemaRefs modifies `$ref` reference strings that match a supplied
// `*regexp.Regexp` and replaces that with a string. It was originally
// designed to convert `#schemas/` to `#components/schemas/`.
func SpecModifySchemaRefs(spec *openapi3.Spec, rx *regexp.Regexp, repl string) {
	if spec == nil || rx == nil {
		return
	}
	for _, paramRef := range spec.Components.Parameters {
		if paramRef == nil {
			continue
		}
		paramRef.Ref = rx.ReplaceAllString(paramRef.Ref, repl)
		if paramRef.Value != nil && paramRef.Value.Schema != nil {
			SchemaRefModifyRefs(paramRef.Value.Schema, rx, repl)
		}
	}

	for _, schRef := range spec.Components.Schemas {
		if schRef == nil {
			continue
		}
		if schRef.Ref == "" && schRef.Value != nil {
			// if `$ref` is populated as an Extension with mmessage type RawMessage,
			// manually convert to reference.
			if refValAny, ok := schRef.Value.Extensions["$ref"]; ok {
				switch reflectutil.TypeName(refValAny) {
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
		schRef.Ref = rx.ReplaceAllString(schRef.Ref, repl)
		SchemaRefModifyRefs(schRef, rx, repl)
	}
}

// SchemaRefModifyRefs modifies Schema reference schema pointers that match
// the supplied `*regexp.Regexp` with the replacement string. It was originally
// designed to convert `#schemas/` to `#components/schemas/`.
func SchemaRefModifyRefs(schRef *oas3.SchemaRef, rx *regexp.Regexp, repl string) {
	if schRef == nil || rx == nil {
		return
	}
	schRef.Ref = rx.ReplaceAllString(schRef.Ref, repl)
	if schRef.Value == nil {
		return
	}
	for _, propSchemaRef := range schRef.Value.Properties {
		SchemaRefModifyRefs(propSchemaRef, rx, repl)
	}
	if schRef.Value.AdditionalProperties != nil {
		SchemaRefModifyRefs(schRef.Value.AdditionalProperties, rx, repl)
	}
	if schRef.Value.Items != nil {
		SchemaRefModifyRefs(schRef.Value.Items, rx, repl)
	}
	if schRef.Value.Not != nil {
		SchemaRefModifyRefs(schRef.Value.Not, rx, repl)
	}
	for _, allOf := range schRef.Value.OneOf {
		SchemaRefModifyRefs(allOf, rx, repl)
	}
	for _, anyOf := range schRef.Value.OneOf {
		SchemaRefModifyRefs(anyOf, rx, repl)
	}
	for _, oneOf := range schRef.Value.OneOf {
		SchemaRefModifyRefs(oneOf, rx, repl)
	}
}

func SpecSchemaSetAdditionalPropertiesTrue(spec *openapi3.Spec, pointerBase string) []string {
	mods := []string{}
	for schName, schRef := range spec.Components.Schemas {
		if schRef == nil || schRef.Value == nil || schRef.Value.Type != openapi3.TypeObject {
			continue
		}
		if len(schRef.Value.Properties) == 0 && schRef.Value.AdditionalProperties == nil &&
			(schRef.Value.AdditionalPropertiesAllowed == nil || !*schRef.Value.AdditionalPropertiesAllowed) {
			additionalPropertiesAllowed := true
			schRef.Value.AdditionalPropertiesAllowed = &additionalPropertiesAllowed
			mods = append(mods, fmt.Sprintf("%s%s/%s", pointerBase, openapi3.PointerComponentsSchemas, schName))
		}

		for propName, propRef := range schRef.Value.Properties {
			if propRef == nil || propRef.Value == nil || propRef.Value.Type != openapi3.TypeObject {
				continue
			}
			if len(propRef.Value.Properties) == 0 && propRef.Value.AdditionalProperties == nil &&
				(propRef.Value.AdditionalPropertiesAllowed == nil || !*propRef.Value.AdditionalPropertiesAllowed) {
				additionalPropertiesAllowed := true
				propRef.Value.AdditionalPropertiesAllowed = &additionalPropertiesAllowed
				mods = append(mods, fmt.Sprintf("%s%s/%s/properties/%s", pointerBase, openapi3.PointerComponentsSchemas, schName, propName))
			}
		}
	}
	return mods
}
