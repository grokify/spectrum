package openapi3

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/io/ioutilmore"
	"github.com/pkg/errors"
)

var jsonFileRx = regexp.MustCompile(`(?i)\.json\s*$`)

func MergeDirectory(dir string) (*oas3.Swagger, error) {
	return MergeDirectoryMore(dir, false, true)
}

func MergeDirectoryMore(dir string, validateEach, validateFinal bool) (*oas3.Swagger, error) {
	fileInfos, err := ioutilmore.DirEntriesRxSizeGt0(dir, ioutilmore.File, jsonFileRx)
	if err != nil {
		return nil, err
	}
	if len(fileInfos) == 0 {
		return nil, fmt.Errorf("No JSON files found in directory [%s]", dir)
	}
	filePaths := []string{}
	for _, fi := range fileInfos {
		filePaths = append(filePaths, filepath.Join(dir, fi.Name()))
	}
	return MergeFiles(filePaths, validateEach, validateFinal)
}

func MergeFiles(filepaths []string, validateEach, validateFinal bool) (*oas3.Swagger, error) {
	var specMaster *oas3.Swagger
	for i, fpath := range filepaths {
		thisSpec, err := ReadFile(fpath, validateEach)
		if err != nil {
			return specMaster, errors.Wrap(err, fmt.Sprintf("Filepath [%v]", fpath))
		}
		if i == 0 {
			specMaster = thisSpec
		} else {
			specMaster = Merge(specMaster, thisSpec)
		}
	}

	if validateFinal {
		bytes, err := specMaster.MarshalJSON()
		if err != nil {
			return specMaster, err
		}
		loader := oas3.NewSwaggerLoader()
		newSpec, err := loader.LoadSwaggerFromData(bytes)
		if err != nil {
			return newSpec, errors.Wrap(err, "Loader.LoadSwaggerFromData")
		}
		return newSpec, nil
	}
	return specMaster, nil
}

func Merge(specMaster, specExtra *oas3.Swagger) *oas3.Swagger {
	specMaster = MergeTags(specMaster, specExtra)
	specMaster = MergePaths(specMaster, specExtra)
	return MergeSchemas(specMaster, specExtra)
}

func MergeTags(specMaster, specExtra *oas3.Swagger) *oas3.Swagger {
	tagsMap := map[string]int{}
	for _, tag := range specMaster.Tags {
		tagsMap[tag.Name] = 1
	}
	for _, tag := range specExtra.Tags {
		tag.Name = strings.TrimSpace(tag.Name)
		if _, ok := tagsMap[tag.Name]; !ok {
			specMaster.Tags = append(specMaster.Tags, tag)
		}
	}
	return specMaster
}

func MergePaths(specMaster, specExtra *oas3.Swagger) *oas3.Swagger {
	for url, path := range specExtra.Paths {
		specMaster.Paths[url] = path
	}
	return specMaster
}

func MergeSchemas(specMaster, specExtra *oas3.Swagger) *oas3.Swagger {
	for schemaName, schema := range specExtra.Components.Schemas {
		specMaster.Components.Schemas[schemaName] = schema
	}
	return specMaster
}
