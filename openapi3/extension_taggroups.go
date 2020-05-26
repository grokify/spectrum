package openapi3

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/type/stringsutil"
)

const XTagGroupsPropertyName = "x-tag-groups"

type TagGroupSet struct {
	TagGroups []TagGroup
}

func NewTagGroupSet() TagGroupSet {
	return TagGroupSet{TagGroups: []TagGroup{}}
}

func (set *TagGroupSet) Exists(tagName string) bool {
	for _, tg := range set.TagGroups {
		for _, tgTagName := range tg.Tags {
			if tagName == tgTagName {
				return true
			}
		}
	}
	return false
}

func (set *TagGroupSet) GetTagGroupNamesForTagNames(wantTagNames ...string) []string {
	tagGroupNames := []string{}
	for _, tg := range set.TagGroups {
		for _, tgTagName := range tg.Tags {
			for _, wantTagName := range wantTagNames {
				if wantTagName == tgTagName {
					tagGroupNames = append(tagGroupNames, tg.Name)
				}
			}
		}
	}
	return stringsutil.SliceCondenseSpace(tagGroupNames, true, true)
}

func (set *TagGroupSet) AddToSpec(spec *oas3.Swagger) error {
	if len(set.TagGroups) == 0 {
		return nil
	}
	missing := TagsWithoutGroups(spec, *set)
	if len(missing) > 0 {
		return fmt.Errorf("E_TAGS_WITHOUT_GROUPS [%s]", strings.Join(missing, ","))
	}
	spec.ExtensionProps.Extensions[XTagGroupsPropertyName] = set.TagGroups
	return nil
}

type TagGroup struct {
	Name    string   `json:"name"`
	Popular bool     `json:"popular"`
	Tags    []string `json:"tags"`
}

func TagsWithoutGroups(spec *oas3.Swagger, tagGroupSet TagGroupSet) []string {
	missing := []string{}
	for _, tag := range spec.Tags {
		if !tagGroupSet.Exists(tag.Name) {
			missing = append(missing, tag.Name)
		}
	}
	return missing
}

// SpecTagGroups parses a TagGroupSet from an OpenAPI3 spec.
func SpecTagGroups(spec *oas3.Swagger) (TagGroupSet, error) {
	tgs := NewTagGroupSet()
	iface, ok := spec.ExtensionProps.Extensions[XTagGroupsPropertyName]
	if !ok {
		return tgs, nil
	}

	tagGroups := []TagGroup{}
	if reflect.TypeOf(iface) == reflect.TypeOf(tagGroups) {
		tgs.TagGroups = iface.([]TagGroup)
		return tgs, nil
	}

	// message is stored as `json.RawMessage` when the data
	// is read in from JSON, vs. set via code.
	rawMessage := iface.(json.RawMessage)
	err := json.Unmarshal(rawMessage, &tagGroups)
	if err != nil {
		return tgs, err
	}
	tgs.TagGroups = tagGroups
	delete(spec.ExtensionProps.Extensions, XTagGroupsPropertyName)
	spec.ExtensionProps.Extensions[XTagGroupsPropertyName] = tagGroups
	return tgs, nil
}
