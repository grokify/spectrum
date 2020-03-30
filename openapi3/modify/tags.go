package modify

import (
	"fmt"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

// SpecTags returns a set of tags present in the current
// spec.
func SpecTags(spec *oas3.Swagger, inclTop, inclOp bool) map[string]int {
	tagsMap := map[string]int{}
	if inclTop {
		for _, tag := range spec.Tags {
			tagName := strings.TrimSpace(tag.Name)
			if len(tagName) > 0 {
				if _, ok := tagsMap[tagName]; !ok {
					tagsMap[tagName] = 0
				}
				tagsMap[tagName]++
			}
		}
	}
	if inclOp {
		VisitOperations(spec, func(op *oas3.Operation) {
			if op == nil {
				return
			}
			for _, tagName := range op.Tags {
				tagName = strings.TrimSpace(tagName)
				if len(tagName) > 0 {
					if _, ok := tagsMap[tagName]; !ok {
						tagsMap[tagName] = 0
					}
					tagsMap[tagName]++
				}
			}
		})
	}
	return tagsMap
}

func SpecTagsModify(spec *oas3.Swagger, changeTagsRaw map[string]string) {
	changeTags := map[string]string{}
	for old, new := range changeTagsRaw {
		changeTags[strings.TrimSpace(old)] = strings.TrimSpace(new)
	}

TAG:
	for i, tag := range spec.Tags {
		tag.Name = strings.TrimSpace(tag.Name)
		if len(tag.Name) > 0 {
			for tOld, tNew := range changeTags {
				if tag.Name == tOld {
					tag.Name = tNew
					spec.Tags[i] = tag
					continue TAG
				}
			}
		}
	}

	VisitOperations(spec, func(op *oas3.Operation) {
		if op == nil {
			return
		}
	TAG:
		for i, tagName := range op.Tags {
			tagName = strings.TrimSpace(tagName)
			if len(tagName) > 0 {
				for tOld, tNew := range changeTags {
					if tagName == tOld {
						op.Tags[i] = tNew
						continue TAG
					}
				}
			}
		}
	})
}

// SpecTagsOrder sorts a specs tags based on an input set
// and explitcit sort order. The remaining tags are sorted
// alphabetically.
func SpecTagsOrder(spec *oas3.Swagger, explicitSortedTagNames []string) error {
	curTags := spec.Tags

	opTagNames := SpecTags(spec, false, true)
	for tagName := range opTagNames {
		curTags = append(curTags, &oas3.Tag{Name: tagName})
	}

	newTags, err := TagsOrder(curTags, explicitSortedTagNames)
	if err != nil {
		return err
	}
	spec.Tags = newTags

	return nil
}

// TagsOrder creates a list of ordered tags based on an input set
// and explitcit sort order. The remaining tags are sorted
// alphabetically.
func TagsOrder(curTags oas3.Tags, explicitSortedTagNames []string) (oas3.Tags, error) {
	tagMap := map[string]*oas3.Tag{}
	newTags := oas3.Tags{}
	for _, tag := range curTags {
		tag.Name = strings.TrimSpace(tag.Name)
		tagMap[tag.Name] = tag
	}
	for _, tagName := range explicitSortedTagNames {
		tagName = strings.TrimSpace(tagName)
		if tag, ok := tagMap[tagName]; ok {
			newTags = append(newTags, tag)
			delete(tagMap, tagName)
		} else {
			return newTags, fmt.Errorf("E_EXPLICIT_TAG_NAME_NOT_FOUND [%v]", tagName)
		}
	}
	remainingSorted := []string{}
	for tagName := range tagMap {
		remainingSorted = append(remainingSorted, tagName)
	}
	sort.Strings(remainingSorted)
	for _, tagName := range remainingSorted {
		if tag, ok := tagMap[tagName]; ok {
			newTags = append(newTags, tag)
		} else {
			return newTags, fmt.Errorf("E_EXPLICIT_TAG_NAME_NOT_FOUND [%v]", tagName)
		}
	}

	return newTags, nil
}
