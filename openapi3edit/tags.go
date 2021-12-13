package openapi3edit

import (
	"fmt"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
)

func SpecTagsModify(spec *openapi3.Spec, changeTagsRaw map[string]string) {
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

	openapi3.VisitOperations(spec, func(skipPath, skipMethod string, op *oas3.Operation) {
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
func SpecTagsOrder(spec *openapi3.Spec, explicitSortedTagNames []string) error {
	curTags := spec.Tags

	sm := openapi3.SpecMore{Spec: spec}
	opTagNames := sm.TagsMap(false, true)
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
	tagMap := map[string]*oas3.Tag{} // curTags
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
			// } else {
			// skip
			// return newTags, fmt.Errorf("spectrum/openapi3/smodify.TagsOrder.Explicit.E_EXPLICIT_TAG_NAME_NOT_FOUND [%v]", tagName)
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
			return newTags, fmt.Errorf("spectrum/openapi3/modify.TagsOrder.Remaining.sE_EXPLICIT_TAG_NAME_NOT_FOUND [%v]", tagName)
		}
	}
	return newTags, nil
}

// SpecTagsCondense removes unused tags from the top
// level specification by comparing with tags used
// in operations.
func SpecTagsCondense(spec *openapi3.Spec) {
	sm := openapi3.SpecMore{Spec: spec}
	opTags := sm.TagsMap(false, true)
	newTags := oas3.Tags{}
	for _, curTag := range spec.Tags {
		if _, ok := opTags[curTag.Name]; ok {
			newTags = append(newTags, curTag)
		}
	}
	spec.Tags = newTags
}

type UpdateTagsOpts struct {
	TagURLsMap   map[string][]string
	TagsMap      map[string]string
	TagGroupsSet openapi3.TagGroupSet
}

func (uto *UpdateTagsOpts) ModifyTagsOperationFunc(path, method string, op *oas3.Operation) {
	if op == nil {
		return
	}
	for tagTry, urlSuffixes := range uto.TagURLsMap {
		tags := strings.Split(tagTry, ",")
		tags = stringsutil.SliceCondenseSpace(tags, true, false)
		if stringsutil.SliceIndexMore(
			urlSuffixes,
			path, true, true, stringsutil.MatchHasSuffix) > -1 {
			op.Tags = tags
		}
	}
}
