package openapi3

import (
	"regexp"
	"strings"

	"github.com/grokify/gotilla/io/ioutilmore"
)

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
	_, files, err := ioutilmore.ReadDirRx(dir, rx, true)

	if err != nil {
		return metas, err
	}

	return ReadSpecMetasFiles(files)
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
