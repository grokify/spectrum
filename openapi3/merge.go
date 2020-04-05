package openapi3

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/data/table"
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
	sort.Strings(filepaths)
	var specMaster *oas3.Swagger
	for i, fpath := range filepaths {
		thisSpec, err := ReadFile(fpath, validateEach)
		if err != nil {
			return specMaster, errors.Wrap(err, fmt.Sprintf("Filepath [%v]", fpath))
		}
		if i == 0 {
			specMaster = thisSpec
		} else {
			specMaster, err = Merge(specMaster, thisSpec)
			if err != nil {
				return nil, err
			}
		}
	}

	if validateFinal {
		bytes, err := specMaster.MarshalJSON()
		if err != nil {
			return specMaster, err
		}
		newSpec, err := oas3.NewSwaggerLoader().LoadSwaggerFromData(bytes)
		if err != nil {
			return newSpec, errors.Wrap(err, "Loader.LoadSwaggerFromData")
		}
		return newSpec, nil
	}
	return specMaster, nil
}

func Merge(specMaster, specExtra *oas3.Swagger) (*oas3.Swagger, error) {
	specMaster = MergeTags(specMaster, specExtra)
	specMaster, err := MergePaths(specMaster, specExtra)
	if err != nil {
		return specMaster, err
	}
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

// MergeWithTables performs a spec merge and returns comparison
// tables. This is useful to combine with github.com/grokify/gocharts/data/table
// WriteXLSX() to write out comparison tables for debugging.
func MergeWithTables(spec1, spec2 *oas3.Swagger) (*oas3.Swagger, []*table.TableData, error) {
	tbls := []*table.TableData{}
	sm1 := SpecMore{Spec: spec1}
	sm2 := SpecMore{Spec: spec2}
	tbls = append(tbls, sm1.OperationsTable())
	tbls[0].Name = "Spec1"
	tbls = append(tbls, sm2.OperationsTable())
	tbls[1].Name = "Spec2"
	specf, err := Merge(spec1, spec2)
	if err != nil {
		return specf, tbls, err
	}
	smf := SpecMore{Spec: specf}
	tbls = append(tbls, smf.OperationsTable())
	tbls[2].Name = "SpecFinal"
	return specf, tbls, nil
}

func MergePaths(specMaster, specExtra *oas3.Swagger) (*oas3.Swagger, error) {
	for url, pathItem := range specExtra.Paths {
		if pathInfoMaster, ok := specMaster.Paths[url]; !ok || pathInfoMaster == nil {
			specMaster.Paths[url] = &oas3.PathItem{}
		}
		if pathItem.Connect != nil {
			if specMaster.Paths[url].Connect == nil {
				specMaster.Paths[url].Connect = pathItem.Connect
			} else if !reflect.DeepEqual(pathItem.Connect, specMaster.Paths[url].Connect) {
				return specMaster, fmt.Errorf("E_OPERATION_COLLISION_CONNECT [%v]", pathItem.Connect.OperationID)
			}
		}
		if pathItem.Delete != nil {
			if specMaster.Paths[url].Delete == nil {
				specMaster.Paths[url].Delete = pathItem.Delete
			} else if !reflect.DeepEqual(pathItem.Delete, specMaster.Paths[url].Delete) {
				return specMaster, fmt.Errorf("E_OPERATION_COLLISION_DELETE [%v]", pathItem.Delete.OperationID)
			}
		}
		if pathItem.Get != nil {
			if specMaster.Paths[url].Get == nil {
				specMaster.Paths[url].Get = pathItem.Get
			} else if !reflect.DeepEqual(pathItem.Get, specMaster.Paths[url].Get) {
				return specMaster, fmt.Errorf("E_OPERATION_COLLISION_GET [%v]", pathItem.Get.OperationID)
			}
		}
		if pathItem.Head != nil {
			if specMaster.Paths[url].Head == nil {
				specMaster.Paths[url].Head = pathItem.Head
			} else if !reflect.DeepEqual(pathItem.Head, specMaster.Paths[url].Head) {
				return specMaster, fmt.Errorf("E_OPERATION_COLLISION_HEAD [%v]", pathItem.Head.OperationID)
			}
		}
		if pathItem.Options != nil {
			if specMaster.Paths[url].Options == nil {
				specMaster.Paths[url].Options = pathItem.Options
			} else if !reflect.DeepEqual(pathItem.Options, specMaster.Paths[url].Options) {
				return specMaster, fmt.Errorf("E_OPERATION_COLLISION_OPTIONS [%v]", pathItem.Options.OperationID)
			}
		}
		if pathItem.Patch != nil {
			if specMaster.Paths[url].Patch == nil {
				specMaster.Paths[url].Patch = pathItem.Patch
			} else if !reflect.DeepEqual(pathItem.Patch, specMaster.Paths[url].Patch) {
				return specMaster, fmt.Errorf("E_OPERATION_COLLISION_PATCH [%v]", pathItem.Patch.OperationID)
			}
		}
		if pathItem.Post != nil {
			if specMaster.Paths[url].Post == nil {
				specMaster.Paths[url].Post = pathItem.Post
			} else if !reflect.DeepEqual(pathItem.Post, specMaster.Paths[url].Post) {
				return specMaster, fmt.Errorf("E_OPERATION_COLLISION_POST [%v]", pathItem.Post.OperationID)
			}
		}
		if pathItem.Put != nil {
			if specMaster.Paths[url].Put == nil {
				specMaster.Paths[url].Put = pathItem.Put
			} else if !reflect.DeepEqual(pathItem.Put, specMaster.Paths[url].Put) {
				return specMaster, fmt.Errorf("E_OPERATION_COLLISION_PUT [%v]", pathItem.Put.OperationID)
			}
		}
		if pathItem.Trace != nil {
			if specMaster.Paths[url].Trace == nil {
				specMaster.Paths[url].Trace = pathItem.Trace
			} else if !reflect.DeepEqual(pathItem.Trace, specMaster.Paths[url].Trace) {
				return specMaster, fmt.Errorf("E_OPERATION_COLLISION_TRACE [%v]", pathItem.Trace.OperationID)
			}
		}
	}
	return specMaster, nil
}

func MergeSchemas(specMaster, specExtra *oas3.Swagger) (*oas3.Swagger, error) {
	for schemaName, schemaExtra := range specExtra.Components.Schemas {
		if schemaExtra == nil {
			continue
		}
		if schemaMaster, ok := specMaster.Components.Schemas[schemaName]; ok {
			if schemaMaster != nil {
				if !reflect.DeepEqual(schemaMaster, schemaExtra) {
					return nil, fmt.Errorf("E_SCHEMA_COLLISION [%v]", schemaName)
				}
				continue
			}
		}
		specMaster.Components.Schemas[schemaName] = schemaExtra
	}
	return specMaster, nil
}

func WriteFileDirMerge(outfile, inputDir string, perm os.FileMode) error {
	spec, err := MergeDirectory(inputDir)
	if err != nil {
		return errors.Wrap(err, "E_OPENAPI3_MERGE_DIRECTORY_FAILED")
	}

	bytes, err := spec.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "E_SWAGGER2_JSON_ENCODING_FAILED")
	}

	err = ioutil.WriteFile(outfile, bytes, perm)
	if err != nil {
		return errors.Wrap(err, "E_SWAGGER2_WRITE_FAILED")
	}
	return nil
}
