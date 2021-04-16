package openapi3

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/simplego/io/ioutilmore"
	"github.com/pkg/errors"
)

var jsonFileRx = regexp.MustCompile(`(?i)\.json\s*$`)

func MergeDirectory(dir string, mergeOpts *MergeOptions) (*oas3.Swagger, int, error) {
	var filePaths []string
	var err error
	if mergeOpts != nil && mergeOpts.FileRx != nil {
		filePaths, err = ioutilmore.DirEntriesRxSizeGt0Filepaths(dir, ioutilmore.File, mergeOpts.FileRx)
	} else {
		filePaths, err = ioutilmore.DirEntriesRxSizeGt0Filepaths(dir, ioutilmore.File, jsonFileRx)
	}

	if err != nil {
		return nil, len(filePaths), err
	}

	spec, err := MergeFiles(filePaths, mergeOpts)
	return spec, len(filePaths), err
}

func MergeFiles(filepaths []string, mergeOpts *MergeOptions) (*oas3.Swagger, error) {
	sort.Strings(filepaths)
	validateEach := false
	validateFinal := true
	if mergeOpts != nil {
		validateEach = mergeOpts.ValidateEach
		validateFinal = mergeOpts.ValidateFinal
	}
	var specMaster *oas3.Swagger
	for i, fpath := range filepaths {
		thisSpec, err := ReadFile(fpath, validateEach)
		if err != nil {
			return specMaster, errors.Wrap(err, fmt.Sprintf("ReadSpecError [%v] ValidateEach [%v]", fpath, validateEach))
		}
		if i == 0 {
			specMaster = thisSpec
		} else {
			specMaster, err = Merge(specMaster, thisSpec, fpath, mergeOpts)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Merging [%v]", fpath))
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
			return newSpec, errors.Wrap(err, "Loader.LoadSwaggerFromData (MergeFiles().ValidateFinal)")
		}
		return newSpec, nil
	}
	return specMaster, nil
}

func Merge(specMaster, specExtra *oas3.Swagger, specExtraNote string, mergeOpts *MergeOptions) (*oas3.Swagger, error) {
	specMaster = MergeTags(specMaster, specExtra)
	specMaster, err := MergeParameters(specMaster, specExtra, specExtraNote, mergeOpts)
	if err != nil {
		return specMaster, err
	}
	specMaster, err = MergeSchemas(specMaster, specExtra, specExtraNote, mergeOpts)
	if err != nil {
		return specMaster, err
	}
	specMaster, err = MergePaths(specMaster, specExtra)
	if err != nil {
		return specMaster, err
	}
	specMaster, err = MergeResponses(specMaster, specExtra, specExtraNote, mergeOpts)
	if err != nil {
		return specMaster, err
	}
	return MergeRequestBodies(specMaster, specExtra, specExtraNote)
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
func MergeWithTables(spec1, spec2 *oas3.Swagger, specExtraNote string, mergeOpts *MergeOptions) (*oas3.Swagger, []*table.Table, error) {
	tbls := []*table.Table{}
	sm1 := SpecMore{Spec: spec1}
	sm2 := SpecMore{Spec: spec2}
	tbls1, err := sm1.OperationsTable(mergeOpts.TableColumns, mergeOpts.TableOpFilterFunc)
	if err != nil {
		return nil, nil, err
	}
	tbls = append(tbls, tbls1)
	tbls[0].Name = "Spec1"
	tbls2, err := sm2.OperationsTable(mergeOpts.TableColumns, mergeOpts.TableOpFilterFunc)
	if err != nil {
		return nil, nil, err
	}
	tbls = append(tbls, tbls2)
	tbls[1].Name = "Spec2"
	specf, err := Merge(spec1, spec2, specExtraNote, mergeOpts)
	if err != nil {
		return specf, tbls, err
	}
	smf := SpecMore{Spec: specf}
	tblsf, err := smf.OperationsTable(mergeOpts.TableColumns, mergeOpts.TableOpFilterFunc)
	if err != nil {
		return nil, nil, err
	}
	tbls = append(tbls, tblsf)

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

func MergeParameters(specMaster, specExtra *oas3.Swagger, specExtraNote string, mergeOpts *MergeOptions) (*oas3.Swagger, error) {
	if specMaster.Components.Parameters == nil {
		specMaster.Components.Parameters = map[string]*oas3.ParameterRef{}
	}
	for pName, pExtra := range specExtra.Components.Parameters {
		if pExtra == nil {
			continue
		} else if pMaster, ok := specMaster.Components.Parameters[pName]; ok {
			if pMaster == nil {
				specMaster.Components.Parameters[pName] = pExtra
			} else {
				if mergeOpts == nil {
					mergeOpts = &MergeOptions{}
				}
				if mergeOpts.CollisionCheckResult == CollisionCheckSkip {
					continue
				} else if reflect.DeepEqual(pExtra, pMaster) {
					continue
				} else if mergeOpts.CollisionCheckResult == CollisionCheckOverwrite {
					specExtra.Components.Parameters[pName] = pExtra
				} else {
					return nil, fmt.Errorf("E_SCHEMA_COLLISION [%v] EXTRA_COMPONENTS_PARAMETER [%s]", specExtraNote, pName)
				}
			}
		} else {
			specMaster.Components.Parameters[pName] = pExtra
		}
	}
	return specMaster, nil
}

func MergeResponses(specMaster, specExtra *oas3.Swagger, specExtraNote string, mergeOpts *MergeOptions) (*oas3.Swagger, error) {
	if specMaster.Components.Responses == nil {
		specMaster.Components.Responses = map[string]*oas3.ResponseRef{}
	}
	for rName, rExtra := range specExtra.Components.Responses {
		if rExtra == nil {
			continue
		} else if rMaster, ok := specMaster.Components.Responses[rName]; ok {
			if rMaster == nil {
				specMaster.Components.Responses[rName] = rExtra
			} else {
				if mergeOpts == nil {
					mergeOpts = &MergeOptions{}
				}
				if mergeOpts.CollisionCheckResult == CollisionCheckSkip {
					continue
				} else if reflect.DeepEqual(rExtra, rMaster) {
					continue
				} else {
					return nil, fmt.Errorf("E_SCHEMA_COLLISION [%v] EXTRA_COMPONENTS_RESPONSE [%s]", specExtraNote, rName)
				}
			}
		} else {
			specMaster.Components.Responses[rName] = rExtra
		}
	}
	return specMaster, nil
}

func MergeSchemas(specMaster, specExtra *oas3.Swagger, specExtraNote string, mergeOpts *MergeOptions) (*oas3.Swagger, error) {
	for schemaName, schemaExtra := range specExtra.Components.Schemas {
		if schemaExtra == nil {
			continue
		} else if schemaMaster, ok := specMaster.Components.Schemas[schemaName]; ok {
			if schemaMaster == nil {
				specMaster.Components.Schemas[schemaName] = schemaExtra
			} else {
				if mergeOpts == nil {
					mergeOpts = &MergeOptions{}
				}
				checkCollisionResult := mergeOpts.CheckSchemaCollision(schemaName, schemaMaster, schemaExtra, specExtraNote)
				if checkCollisionResult != CollisionCheckSame &&
					mergeOpts.CollisionCheckResult != CollisionCheckSkip {
					if mergeOpts.CollisionCheckResult == CollisionCheckOverwrite {
						delete(specMaster.Components.Schemas, schemaName)
						specMaster.Components.Schemas[schemaName] = schemaExtra
					} else if mergeOpts.CollisionCheckResult == CollisionCheckError {
						return nil, fmt.Errorf("E_SCHEMA_COLLISION [%v] EXTRA_SPEC [%s]", schemaName, specExtraNote)
					}
				}
				/*
					if !reflect.DeepEqual(schemaMaster, schemaExtra) {
						return nil, fmt.Errorf("E_SCHEMA_COLLISION [%v] EXTRA_SPEC [%s]", schemaName, specExtraNote)
					}*/
				continue
			}
		} else {
			specMaster.Components.Schemas[schemaName] = schemaExtra
		}
	}
	return specMaster, nil
}

func MergeRequestBodies(specMaster, specExtra *oas3.Swagger, specExtraNote string) (*oas3.Swagger, error) {
	for rbName, rbExtra := range specExtra.Components.RequestBodies {
		if rbExtra == nil {
			continue
		} else if rbMaster, ok := specMaster.Components.RequestBodies[rbName]; ok {
			if rbMaster == nil {
				if specMaster.Components.RequestBodies == nil {
					specMaster.Components.RequestBodies = map[string]*oas3.RequestBodyRef{}
				}
				specMaster.Components.RequestBodies[rbName] = rbExtra
			} else if !reflect.DeepEqual(rbMaster, rbExtra) {
				return nil, fmt.Errorf("E_SCHEMA_COLLISION [%v] EXTRA_SPEC [%s]", rbName, specExtraNote)
			}
		} else {
			if specMaster.Components.RequestBodies == nil {
				specMaster.Components.RequestBodies = map[string]*oas3.RequestBodyRef{}
			}
			specMaster.Components.RequestBodies[rbName] = rbExtra
		}
	}
	return specMaster, nil
}

func WriteFileDirMerge(outfile, inputDir string, perm os.FileMode, mergeOpts *MergeOptions) (int, error) {
	spec, num, err := MergeDirectory(inputDir, mergeOpts)
	if err != nil {
		return num, errors.Wrap(err, "E_OPENAPI3_MERGE_DIRECTORY_FAILED")
	}

	bytes, err := spec.MarshalJSON()
	if err != nil {
		return num, errors.Wrap(err, "E_SWAGGER2_JSON_ENCODING_FAILED")
	}

	err = ioutil.WriteFile(outfile, bytes, perm)
	if err != nil {
		return num, errors.Wrap(err, "E_SWAGGER2_WRITE_FAILED")
	}
	return num, nil
}
