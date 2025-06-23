package openapi2

import (
	"github.com/grokify/spectrum"
)

func NewSpecMetaFilepath(filepath string) (*spectrum.SpecMeta, error) {
	if s, err := ReadOpenAPI2KinSpecFile(filepath); err != nil {
		return nil, err
	} else {
		sm := SpecMore{Spec: s}
		return sm.Meta(), nil
	}
}

func NewSpecMetaSetFilepaths(filepaths []string) (*spectrum.SpecMetaSet, error) {
	set := spectrum.NewSpecMetaSet()
	for _, fp := range filepaths {
		if m, err := NewSpecMetaFilepath(fp); err != nil {
			return nil, err
		} else {
			set.Data[fp] = *m
		}
	}
	return set, nil
}
