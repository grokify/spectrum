package modify

import (
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
)

func SpecTags(spec *oas3.Swagger) map[string]int {
	tagsMap := map[string]int{}
	for _, tag := range spec.Tags {
		tagName := strings.TrimSpace(tag.Name)
		if len(tagName) > 0 {
			if _, ok := tagsMap[tagName]; !ok {
				tagsMap[tagName] = 0
			}
			tagsMap[tagName]++
		}
	}
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
