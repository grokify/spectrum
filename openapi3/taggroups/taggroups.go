package taggroups

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

const XTagGroupsPropertyName = "x-tag-groups"

type TagGroupSet struct {
	TagGroups []TagGroup
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
