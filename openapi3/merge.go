package openapi3

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/io/ioutilmore"
)

var jsonFileRx = regexp.MustCompile(`(?i)\.json\s*$`)

func MergeDirectory(dir string) (*oas3.Swagger, error) {
	fileInfos, err := ioutilmore.DirEntriesRxSizeGt0(dir, ioutilmore.File, jsonFileRx)
	if err != nil {
		return nil, err
	}
	if len(fileInfos) == 0 {
		return nil, fmt.Errorf("No JSON files found in directory [%s]", dir)
	}
	loader := oas3.NewSwaggerLoader()
	var specMaster *oas3.Swagger
	for i, fi := range fileInfos {
		thisSpecFilepath := filepath.Join(dir, fi.Name())
		thisSpec, err := loader.LoadSwaggerFromFile(thisSpecFilepath)
		if err != nil {
			return specMaster, err
		}
		if i == 0 {
			specMaster = thisSpec
		} else {
			specMaster = Merge(specMaster, thisSpec)
		}
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
	for _, tag := range specMaster.Components.Tags {
		tagsMap[tag.Name] = 1
	}
	for _, tag := range specExtra.Components.Tags {
		tag.Name = strings.TrimSpace(tag.Name)
		if _, ok := tagsMap[tag.Name]; !ok {
			specMaster.Components.Tags = append(specMaster.Components.Tags, tag)
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
