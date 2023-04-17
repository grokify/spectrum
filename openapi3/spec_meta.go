package openapi3

import (
	"regexp"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/os/osutil"
)

// NewSpec returns a new OpenAPI 3 spec that will validate.
// Specifically, it includes an OAS version, sets `info` to
// be an empty object instead of null and sets apiVersion.
func NewSpec(oasVersion, apiTitle, apiVersion string) *Spec {
	oasVersion = strings.TrimSpace(oasVersion)
	if len(oasVersion) == 0 {
		oasVersion = OASVersionLatest
	}
	apiVersion = strings.TrimSpace(apiVersion)
	if len(apiVersion) == 0 {
		apiVersion = apiVersionDefault
	}
	return &Spec{
		OpenAPI: oasVersion,
		Info: &oas3.Info{
			Title:   strings.TrimSpace(apiTitle),
			Version: apiVersion}}
}

type SpecMetas struct {
	Metas []SpecMeta
}

func (metas *SpecMetas) Filepaths(validOnly bool) []string {
	files := []string{}
	for _, meta := range metas.Metas {
		if validOnly && !meta.IsValid {
			continue
		}
		meta.Filepath = strings.TrimSpace(meta.Filepath)
		if len(meta.Filepath) > 0 {
			files = append(files, meta.Filepath)
		}
	}
	return files
}

type SpecMeta struct {
	Filepath        string
	Version         int
	IsValid         bool
	ValidationError string
}

func ReadSpecMetasDir(dir string, rx *regexp.Regexp) (SpecMetas, error) {
	metas := SpecMetas{Metas: []SpecMeta{}}
	entries, err := osutil.ReadDirMore(dir, rx, false, true, false)

	if err != nil {
		return metas, err
	}

	return ReadSpecMetasFiles(entries.Names(dir))
}

func ReadSpecMetasFiles(files []string) (SpecMetas, error) {
	metas := SpecMetas{Metas: []SpecMeta{}}
	for _, f := range files {
		_, err := ReadFile(f, true)
		meta := SpecMeta{
			Filepath: f,
			Version:  3}
		if err != nil {
			meta.ValidationError = err.Error()
		} else {
			meta.IsValid = true
		}
		metas.Metas = append(metas.Metas, meta)
	}
	return metas, nil
}

func (metas *SpecMetas) Merge(validatesOnly bool, mergeOpts *MergeOptions) (SpecMore, error) {
	return MergeSpecMetas(metas, validatesOnly, mergeOpts)
}

func MergeSpecMetas(metas *SpecMetas, validatesOnly bool, mergeOpts *MergeOptions) (SpecMore, error) {
	specMore := SpecMore{}
	filepaths := metas.Filepaths(validatesOnly)
	spec, err := MergeFiles(filepaths, mergeOpts)
	if err != nil {
		return specMore, err
	}
	specMore.Spec = spec
	return specMore, nil
}
