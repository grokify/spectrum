package openapi3postman2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/spectrum/ext/taggroups"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/postman2"
	"github.com/grokify/spectrum/postman2/simple"
)

// Converter is the struct that manages the conversion.
type Converter struct {
	Configuration Configuration
	OpenAPISpec   *openapi3.Spec
}

// NewConverter instantiates a new converter.
func NewConverter(cfg Configuration) Converter {
	return Converter{Configuration: cfg}
}

// MergeConvert builds a Postman 2.0 spec using a base Postman 2.0 collection
// and a OpenAPI 3.0 spec.
func (conv *Converter) MergeConvert(openapiFilepath string, pmanBaseFilepath string, pmanSpecFilepath string) error {
	oas3spec, err := openapi3.ReadFile(openapiFilepath, true)
	if err != nil {
		return errorsutil.Wrap(err,
			fmt.Sprintf(
				"cannot read OpenAPI 3 spec [%s] openapi3postman2.Converter.MergeConvert << openapi3.ReadFile",
				openapiFilepath))
	}

	pmanBaseFilepath = strings.TrimSpace(pmanBaseFilepath)
	if len(pmanBaseFilepath) > 0 {
		pman, err := simple.ReadCanonicalCollection(pmanBaseFilepath)
		if err != nil {
			err = errorsutil.Wrap(err,
				fmt.Sprintf(
					"cannot read Postman Collection [%s] openapi3postman2.Converter.MergeConvert << simple.ReadCanonicalCollection",
					pmanBaseFilepath))
			return err
		}

		pm, err := Merge(conv.Configuration, pman, oas3spec)
		if err != nil {
			return err
		}

		bytes, err := json.MarshalIndent(pm, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile(pmanSpecFilepath, bytes, 0600)
	}
	return conv.ConvertFile(openapiFilepath, pmanSpecFilepath)
}

// ConvertFile builds a Postman 2.0 spec using an OpenAPI 3.0 spec.
func (conv *Converter) ConvertFile(openapiFilepath string, pmanSpecFilepath string) error {
	oas3Loader := oas3.NewLoader()
	oas3spec, err := oas3Loader.LoadFromFile(openapiFilepath)
	if err != nil {
		return err
	}
	pm, err := ConvertSpec(conv.Configuration, oas3spec)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(pm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pmanSpecFilepath, bytes, 0600)
}

// ConvertSpec creates a Postman 2.0 collection from a configuration and Swagger 2.0 spec
func ConvertSpec(cfg Configuration, oas3spec *openapi3.Spec) (postman2.Collection, error) {
	return Merge(cfg, postman2.Collection{}, oas3spec)
}

// Merge creates a Postman 2.0 collection from a configuration, base Postman
// 2.0 collection and Swagger 2.0 spec
func Merge(cfg Configuration, pman postman2.Collection, oas3spec *openapi3.Spec) (postman2.Collection, error) {
	if len(pman.Info.Name) == 0 {
		pman.Info.Name = strings.TrimSpace(oas3spec.Info.Title)
	}
	if len(pman.Info.Description) == 0 {
		pman.Info.Description = strings.TrimSpace(oas3spec.Info.Description)
	}
	if len(pman.Info.Schema) == 0 {
		pman.Info.Schema = "https://schema.getpostman.com/json/collection/v2.0.0/collection.json"
	}

	pman, err := CreateTagsAndTagGroups(pman, oas3spec)
	if err != nil {
		return pman, err
	}

	// tagGroupSet, err := openapi3.SpecTagGroups(oas3spec)
	// oas3specMore := openapi3.SpecMore{Spec: oas3spec}
	// tagGroupSet, err := oas3specMore.TagGroups()
	tagGroupSet, err := taggroups.SpecTagGroups(oas3spec)
	if err != nil {
		return pman, err
	}

	urls := []string{}
	for url := range oas3spec.Paths {
		urls = append(urls, url)
	}
	sort.Strings(urls)

	for _, url := range urls {
		path := oas3spec.Paths[url] // *PathItem

		if path.Delete != nil {
			pitem, err := Openapi3OperationToPostman2APIItem(cfg, oas3spec, url, http.MethodDelete, path.Delete)
			if err != nil {
				return pman, err
			}
			pman = postmanAddItemToFolders(pman, pitem, path.Delete.Tags, tagGroupSet)
		}
		if path.Get != nil {
			pitem, err := Openapi3OperationToPostman2APIItem(cfg, oas3spec, url, http.MethodGet, path.Get)
			if err != nil {
				return pman, err
			}
			pman = postmanAddItemToFolders(pman, pitem, path.Get.Tags, tagGroupSet)
		}
		if path.Patch != nil {
			pitem, err := Openapi3OperationToPostman2APIItem(cfg, oas3spec, url, http.MethodPatch, path.Patch)
			if err != nil {
				return pman, err
			}
			pman = postmanAddItemToFolders(pman, pitem, path.Patch.Tags, tagGroupSet)
		}
		if path.Post != nil {
			pitem, err := Openapi3OperationToPostman2APIItem(cfg, oas3spec, url, http.MethodPost, path.Post)
			if err != nil {
				return pman, err
			}
			pman = postmanAddItemToFolders(pman, pitem, path.Post.Tags, tagGroupSet)
		}
		if path.Put != nil {
			pitem, err := Openapi3OperationToPostman2APIItem(cfg, oas3spec, url, http.MethodPut, path.Put)
			if err != nil {
				return pman, err
			}
			pman = postmanAddItemToFolders(pman, pitem, path.Put.Tags, tagGroupSet)
		}
	}
	return pman, nil
}

func postmanAddItemToFolders(pman postman2.Collection, pmItem *postman2.Item, tagNames []string, tagGroupSet taggroups.TagGroupSet) postman2.Collection {
	for _, tagName := range tagNames {
		tagGroupNames := tagGroupSet.GetTagGroupNamesForTagNames(tagName)
		if len(tagGroupNames) == 0 {
			pmFolder := pman.GetOrNewFolder(tagName)
			pmFolder.Item = append(pmFolder.Item, pmItem)
			pman.SetFolder(pmFolder)
		} else {
			for _, tagGroupName := range tagGroupNames {
				pmFolder := pman.GetOrNewFolder(tagGroupName)
				if pmFolder.Item == nil {
					pmFolder.Item = []*postman2.Item{}
				}
				// Tags
				modded := false
				for i, pmfSubItem := range pmFolder.Item {
					if pmfSubItem.Name == tagName {
						if pmfSubItem.Item == nil {
							if pmfSubItem.Item == nil {
								pmfSubItem.Item = []*postman2.Item{}
							}
						}
						pmfSubItem.Item = append(pmfSubItem.Item, pmItem)
						pmFolder.Item[i] = pmfSubItem
						modded = true
					}
				}
				if modded {
					pman.SetFolder(pmFolder)
				}
			}
		}
	}

	return pman
}

/*
func postmanAddItemToFolder(pman postman2.Collection, pmItem *postman2.Item, pmFolderName string) postman2.Collection {
	pmFolder := pman.GetOrNewFolder(pmFolderName)
	pmFolder.Item = append(pmFolder.Item, pmItem)
	pman.SetFolder(pmFolder)
	return pman
}
*/

func Openapi3OperationToPostman2APIItem(cfg Configuration, oas3spec *openapi3.Spec, oasUrl string, method string, operation *oas3.Operation) (*postman2.Item, error) {
	pmUrl := BuildPostmanURL(cfg, oas3spec, oasUrl, operation)
	item := &postman2.Item{
		Name: operation.Summary,
		Request: &postman2.Request{
			Method: strings.ToUpper(method),
			URL:    &pmUrl,
		},
	}

	if len(strings.TrimSpace(operation.Description)) > 0 {
		item.Request.Description = strings.TrimSpace(operation.Description)
	}

	headers := cfg.PostmanHeaders

	headers, _, _, err := postman2.AddOperationReqResMediaTypeHeaders(
		headers, operation, oas3spec,
		postman2.DefaultMediaTypePreferencesSlice(),
		postman2.DefaultMediaTypePreferencesSlice(),
	)
	if err != nil {
		return nil, err
	}

	item.Request.Header = headers

	params := ParamsOpenAPI3ToPostman(operation.Parameters)
	if len(params.Query) > 0 {
		item.Request.URL.Query = params.Query
	}
	if len(params.Variable) > 0 {
		item.Request.URL.Variable = params.Variable
	}

	if cfg.RequestBodyFunc != nil {
		bodyString := strings.TrimSpace(cfg.RequestBodyFunc(oasUrl))
		if len(bodyString) > 0 {
			item.Request.Body = &postman2.RequestBody{
				Mode: "raw",
				Raw:  bodyString}
		}
	}

	return item, nil
}

func BuildPostmanURL(cfg Configuration, spec *openapi3.Spec, specPath string, operation *oas3.Operation) postman2.URL {
	specMore := openapi3.SpecMore{Spec: spec}
	specServerURL := specMore.ServerURL(0)
	partsOverrideURL := []string{}
	cfg.PostmanServerURL = strings.TrimSpace(cfg.PostmanServerURL)
	if len(cfg.PostmanServerURL) > 0 {
		partsOverrideURL = append(partsOverrideURL, cfg.PostmanServerURL)
	}
	cfg.PostmanServerURLBasePath = strings.TrimSpace(cfg.PostmanServerURLBasePath)
	if len(cfg.PostmanServerURLBasePath) > 0 {
		partsOverrideURL = append(partsOverrideURL, cfg.PostmanServerURLBasePath)
	}
	overrideServerURL := urlutil.JoinAbsolute(partsOverrideURL...)

	specURLString := openapi3.BuildApiURLOAS(specServerURL, overrideServerURL, specPath)
	pmanURLString := postman2.APIURLOasToPostman(specURLString)
	pmanURL := postman2.NewURL(pmanURLString)
	pmanURL = PostmanURLAddDefaultsOAS3(pmanURL, operation)
	return pmanURL
}

var postmanUrlDefaultsRx *regexp.Regexp = regexp.MustCompile(`^\s*(:(.+))\s*$`)

func PostmanURLAddDefaultsOAS3(pmanURL postman2.URL, operation *oas3.Operation) postman2.URL {
	for _, part := range pmanURL.Path {
		match := postmanUrlDefaultsRx.FindAllStringSubmatch(part, -1)
		if len(match) > 0 {
			baseVariable := match[0][2]
			var defaultValue interface{}
			for _, parameterRef := range operation.Parameters {
				if parameterRef == nil || parameterRef.Value == nil {
					continue
				}
				if parameterRef.Value.Name != baseVariable {
					continue
				}
				schemaRef := parameterRef.Value.Schema
				if schemaRef == nil || schemaRef.Value == nil {
					continue
				}
				if schemaRef.Value.Default != nil {
					defaultValue = schemaRef.Value.Default
					pmanURL.AddVariable(baseVariable, defaultValue)
				}
				/*
					if parameter.Name == baseVariable {
						defaultValue = parameter.Default
						break
					}*/
			}
			//pmanURL.AddVariable(baseVariable, defaultValue)
		}
	}
	return pmanURL
}

// ParamsOpenAPI3ToPostman returns a slices of Postman parameters
// for a slice of OpenAPI 3 parameters.
func ParamsOpenAPI3ToPostman(oparams []*oas3.ParameterRef) postman2.URLParameters {
	pparams := postman2.NewURLParameters()
	for _, oparamRef := range oparams {
		if oparamRef == nil || oparamRef.Value == nil {
			continue
		}
		oparam := oparamRef.Value
		if oparam.In == oas3.ParameterInQuery {
			pparams.Query = append(pparams.Query,
				postman2.URLQuery{
					Key:         oparam.Name,
					Value:       schemaToString(oparam.Schema),
					Description: oparam.Description,
					Disabled:    true,
				},
			)
		}
	}
	return pparams
}

func schemaToString(schemaRef *oas3.SchemaRef) string {
	if schemaRef == nil || schemaRef.Value == nil {
		return ""
	}
	schema := schemaRef.Value
	parts := []string{}
	schema.Type = strings.TrimSpace(schema.Type)
	schema.Format = strings.TrimSpace(schema.Format)
	if len(schema.Type) > 0 {
		parts = append(parts, schema.Type)
	}
	if len(schema.Format) > 0 {
		parts = append(parts, schema.Format)
	}
	if strings.ToLower(schema.Type) == "array" {
		if schema.Items != nil && schema.Items.Value != nil {
			parts = append(parts, schema.Items.Value.Type)
		}
	}
	if len(parts) > 0 {
		if len(parts) == 2 && parts[0] == "array" && parts[1] == "string" {
			parts = append(parts, "csv")
		}
		return "<" + strings.Join(parts, ".") + ">"
	} else {
		return ""
	}
}
