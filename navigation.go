package swaggman

import (
	"strings"

	"github.com/grokify/gotilla/net/httputilmore"
)

type Navigation struct {
	Endpoints []Endpoint // Loadable from OpenAPI spec
	Tags      []Tag      // Additional metadata with nesting
}

func (nav *Navigation) Inflate() {
	if len(nav.Tags) == 0 {
		tagsMap := map[string]int{}
		for _, ep := range nav.Endpoints {
			ep.Tag = strings.TrimSpace(ep.Tag)
			if len(ep.Tag) > 0 {
				tagsMap[ep.Tag] = 1
			}
		}
		if len(tagsMap) > 0 {
			tagsSlice := []Tag{}
			for tag := range tagsMap {
				tagsSlice = append(tagsSlice,
					Tag{
						Name:      tag,
						SubTags:   []Tag{},
						Endpoints: []Endpoint{},
					},
				)
			}
		}
	} else {

	}
}

// Endpoint represents a simple endpoint
// lookup metadata. It supports only a simple
// tag.
type Endpoint struct {
	Method httputilmore.HTTPMethod
	Path   string
	Tag    string
}

// Tag represents a generic tag, including sub-tags.
type Tag struct {
	Name      string
	NamePath  string // TSV
	SubTags   []Tag
	Endpoints []Endpoint
}

func (tag *Tag) Name() string {
	if len(tag.TagPath) == 0 {
		return ""
	}
	tag.TagPath[len(tag.TagPath)-1] = strings.TrimSpace(
		tag.TagPath[len(tag.TagPath)-1],
	)
	return tag.TagPath[len(tag.TagPath)-1]
}
