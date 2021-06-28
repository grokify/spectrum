package openapi3

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/grokify/simplego/type/stringsutil"
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

func (set *TagGroupSet) AddToSpec(spec *Spec) error {
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

/*
func (sm *SpecMore) TagsWithoutGroups() ([]string, []string, []string, error) {
	tgs, err := sm.TagGroups()
	if err != nil {
		return []string{}, []string{}, []string{}, err
	}
	allTags := []string{}

	topTags := stringsutil.SliceCondenseSpace(sm.Tags(true, false), true, true)
	allTags = append(allTags, topTags...)

	opsTags := stringsutil.SliceCondenseSpace(sm.Tags(false, true), true, true)
	allTags = append(allTags, opsTags...)

	allTags = stringsutil.SliceCondenseSpace(allTags, true, true)
	return allTags, topTags, opsTags, nil
}
*/

func TagsWithoutGroups(spec *Spec, tagGroupSet TagGroupSet) []string {
	missing := []string{}
	for _, tag := range spec.Tags {
		if !tagGroupSet.Exists(tag.Name) {
			missing = append(missing, tag.Name)
		}
	}
	return missing
}

// TagGroups parses a TagGroupSet from an OpenAPI3 spec.
func (sm *SpecMore) TagGroups() (TagGroupSet, error) {
	tgs := NewTagGroupSet()
	iface, ok := sm.Spec.ExtensionProps.Extensions[XTagGroupsPropertyName]
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
	delete(sm.Spec.ExtensionProps.Extensions, XTagGroupsPropertyName)
	sm.Spec.ExtensionProps.Extensions[XTagGroupsPropertyName] = tagGroups
	return tgs, nil
}
