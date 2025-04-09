package openapi3edit

import (
	"fmt"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/mogo/type/strslices"
	"github.com/grokify/spectrum/openapi3"
)

func (se *SpecEdit) TagsModifyMore(opts *TagsModifyOpts) {
	if opts == nil {
		return
	}
	if se.SpecMore.Spec == nil {
		return
	}
	if len(opts.TagsMap) > 0 {
		se.TagsModify(opts.TagsMap)
	}
	if len(opts.TagURLsMap) > 0 {
		openapi3.VisitOperations(se.SpecMore.Spec, opts.ModifyTagsOperationFunc)
	}
}

// TagsModify renames tags using mapping of old tags to new tags.
func (se *SpecEdit) TagsModify(mapTagsOldToNew map[string]string) {
	spec := se.SpecMore.Spec
	if spec == nil {
		return
	}
	changeTags := map[string]string{}
	for old, new := range mapTagsOldToNew {
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
func (se *SpecEdit) TagsOrder(explicitSortedTagNames []string) error {
	if se.SpecMore.Spec == nil {
		return openapi3.ErrSpecNotSet
	}
	curTags := se.SpecMore.Spec.Tags

	opTagNames := se.SpecMore.TagsMap(&openapi3.TagsOpts{InclDefs: false, InclOps: true})
	for tagName := range opTagNames {
		curTags = append(curTags, &oas3.Tag{Name: tagName})
	}

	newTags, err := TagsOrder(curTags, explicitSortedTagNames)
	if err != nil {
		return err
	}
	se.SpecMore.Spec.Tags = newTags
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
func (se *SpecEdit) SpecTagsCondense() {
	if se.SpecMore.Spec == nil {
		return
	}
	//sm := openapi3.SpecMore{Spec: spec}
	opTags := se.SpecMore.TagsMap(&openapi3.TagsOpts{InclDefs: false, InclOps: true})
	newTags := oas3.Tags{}
	for _, curTag := range se.SpecMore.Spec.Tags {
		if _, ok := opTags[curTag.Name]; ok {
			newTags = append(newTags, curTag)
		}
	}
	se.SpecMore.Spec.Tags = newTags
}

// TagsModifyOpts is used with `SpecEdit.TagsModifyMore()`.
type TagsModifyOpts struct {
	// TagURLsMap is a map where the keys are the new tags to add in CSV format while the values
	// are a set of URL suffixes.
	TagURLsMap map[string][]string
	// TagsMap is a map where the keys are the old tag to modify and the value is the tag to use.
	TagsMap map[string]string
	// TagGroupsSet is a tag group set which can be added using Redocly's `x-tagGroups` property
	// as `spec.Extensions["x-tagGroups"] = opts.TagGroupsSet.TagGroups``
	// TagGroupsSet openapi3.TagGroupSet
}

// ModifyTagsOperationFunc satisfies the function signature used in `openapi3.VisitOperations`.`
func (tmo *TagsModifyOpts) ModifyTagsOperationFunc(path, method string, op *oas3.Operation) {
	if op == nil {
		return
	}
	for tagTry, urlSuffixes := range tmo.TagURLsMap {
		tags := strings.Split(tagTry, ",")
		tags = stringsutil.SliceCondenseSpace(tags, true, false)
		if strslices.IndexMore(
			urlSuffixes,
			path, true, true, stringsutil.MatchStringSuffix) > -1 {
			op.Tags = tags
		}
	}
}

// NewTagsSimple returns a `Tags` struct that can be assigned to `openapi3.Spec.Tags`.
func NewTagsSimple(tagNames []string) oas3.Tags {
	var tags oas3.Tags
	for _, tagName := range tagNames {
		tags = append(tags, &oas3.Tag{Name: tagName})
	}
	return tags
}
