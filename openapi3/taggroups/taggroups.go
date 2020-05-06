package taggroups

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

type TagGroupSet struct {
	TagGroups []TagGroup
}

func (set *TagGroupSet) Exists(tagName string) bool {
	for _, tg := range set.TagGroups {
		for _, tagName := range tg.Tags {
			if tagName == tg.Name {
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
	/*	if spec.Extensions == nil {

		}*/
	spec.ExtensionProps.Extensions["x-tag-groups"] = set.TagGroups
	return nil
}

type TagGroup struct {
	Name    string   `json:"name"`
	Popular bool     `json:"popular"`
	Tags    []string `json:"tags"`
}

func TagsWithoutGroups(spec *oas3.Swagger, tagGroupSet TagGroupSet) []string {
	missing := []string{}
	for i, tag := range spec.Tags {
		tagName := strings.TrimSpace(tag.Name)
		if tagName != tag.Name {
			spec.Tags[i].Name = tagName
		}
		if !tagGroupSet.Exists(tagName) {
			missing = append(missing, tagName)
		}

	}
	return missing
}
