package topostman2

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/swaggman/openapi3"
	"github.com/grokify/swaggman/postman2"
)

func CreateTagsAndTagGroups(pman postman2.Collection, spec *oas3.Swagger) (postman2.Collection, error) {
	oas3specMore := openapi3.SpecMore{Spec: spec}
	tagGroupSet, err := oas3specMore.TagGroups()
	// tagGroupSet, err := openapi3.SpecTagGroups(spec)
	if err != nil {
		return pman, err
	}
	if len(tagGroupSet.TagGroups) > 0 {
		return addFoldersFromTagGroups(pman, tagGroupSet, spec.Tags)
	}
	return addFoldersFromTags(pman, spec.Tags), nil
}

func addFoldersFromTagGroups(pman postman2.Collection, tgSet openapi3.TagGroupSet, tags oas3.Tags) (postman2.Collection, error) {
	tagsMore := openapi3.TagsMore{Tags: tags}
	for _, tg := range tgSet.TagGroups {
		tg.Name = strings.TrimSpace(tg.Name)
		if len(tg.Name) == 0 && len(tg.Tags) > 0 {
			return pman, fmt.Errorf("E_TAG_GROUP_EMPTY_NAME TAGS [%s]", strings.Join(tg.Tags, ","))
		}
		topFolder := pman.GetOrNewFolder(tg.Name)
		if topFolder.Item == nil {
			topFolder.Item = []*postman2.Item{}
		}
		for _, tagName := range tg.Tags {
			tagName = strings.TrimSpace(tagName)
			if len(tagName) == 0 {
				continue
			}
			subFolder := &postman2.Item{Name: tagName}
			tag := tagsMore.Get(tagName)
			if tag != nil {
				subFolder.Description = strings.TrimSpace(tag.Description)
			}
			topFolder.UpsertSubItem(subFolder)
		}
		pman.SetFolder(topFolder)
	}
	return pman, nil
}

func addFoldersFromTags(pman postman2.Collection, tags oas3.Tags) postman2.Collection {
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		folder := &postman2.Item{
			Name:        strings.TrimSpace(tag.Name),
			Description: strings.TrimSpace(tag.Description)}
		if len(folder.Name) == 0 {
			continue
		}
		pman.SetFolder(folder)
	}
	return pman
}
