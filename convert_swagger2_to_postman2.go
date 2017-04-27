package swaggman

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/grokify/swaggman/postman2"
	"github.com/grokify/swaggman/postman2/simple"
	"github.com/grokify/swaggman/swagger2"
)

type Configuration struct {
	PostmanURLHostname string            `json:"postmanURLHostname,omitempty"`
	PostmanHeaders     []postman2.Header `json:"postmanHeaders,omitempty"`
}

type Converter struct {
	Configuration Configuration
	Swagger       swagger2.Specification
}

func NewConverter(cfg Configuration) Converter {
	return Converter{Configuration: cfg}
}

func (conv *Converter) MergeConvert(swaggerFilepath string, pmanBaseFilepath string, pmanSpecFilepath string) error {
	swag, err := swagger2.ReadSwagger2Spec(swaggerFilepath)
	if err != nil {
		return err
	}

	pman, err := simple.ReadCanonicalCollection(pmanBaseFilepath)
	if err != nil {
		return err
	}

	pm := Merge(conv.Configuration, pman, swag)

	bytes, err := json.MarshalIndent(pm, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pmanSpecFilepath, bytes, 0644)
}

func (conv *Converter) Convert(swaggerFilepath string, pmanSpecFilepath string) error {
	swag, err := swagger2.ReadSwagger2Spec(swaggerFilepath)
	if err != nil {
		return err
	}
	pm := Convert(conv.Configuration, swag)

	bytes, err := json.MarshalIndent(pm, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pmanSpecFilepath, bytes, 0644)
}

func Convert(cfg Configuration, swag swagger2.Specification) postman2.Collection {
	return Merge(cfg, postman2.Collection{}, swag)
}

func Merge(cfg Configuration, pman postman2.Collection, swag swagger2.Specification) postman2.Collection {
	if len(pman.Info.Name) == 0 {
		pman.Info.Name = strings.TrimSpace(swag.Info.Title)
	}
	if len(pman.Info.Description) == 0 {
		pman.Info.Description = strings.TrimSpace(swag.Info.Description)
	}
	if len(pman.Info.Schema) == 0 {
		pman.Info.Schema = "https://schema.getpostman.com/json/collection/v2.0.0/collection.json"
	}

	for url, path := range swag.Paths {
		if url != "/v1.0/account/{accountId}/extension" {
			//continue
		}

		if len(path.Get.Tags) > 0 {
			if len(strings.TrimSpace(path.Get.Tags[0])) > 0 {
				pmItem := Swagger2PathToPostman2ApiItem(cfg, swag, url, "GET", path.Get)
				pmFolderName := strings.TrimSpace(path.Get.Tags[0])
				pmFolder := pman.GetOrNewFolder(pmFolderName)
				pmFolder.Item = append(pmFolder.Item, pmItem)
				pman.SetFolder(pmFolder)
			}
		}
		if len(path.Post.Tags) > 0 {
			if len(strings.TrimSpace(path.Post.Tags[0])) > 0 {
				pmItem := Swagger2PathToPostman2ApiItem(cfg, swag, url, "POST", path.Post)
				pmFolderName := strings.TrimSpace(path.Post.Tags[0])
				pmFolder := pman.GetOrNewFolder(pmFolderName)
				pmFolder.Item = append(pmFolder.Item, pmItem)
				pman.SetFolder(pmFolder)
			}
		}
		if len(path.Put.Tags) > 0 {
			if len(strings.TrimSpace(path.Put.Tags[0])) > 0 {
				pmItem := Swagger2PathToPostman2ApiItem(cfg, swag, url, "PUT", path.Put)
				pmFolderName := strings.TrimSpace(path.Put.Tags[0])
				pmFolder := pman.GetOrNewFolder(pmFolderName)
				pmFolder.Item = append(pmFolder.Item, pmItem)
				pman.SetFolder(pmFolder)
			}
		}
		if len(path.Delete.Tags) > 0 {
			if len(strings.TrimSpace(path.Delete.Tags[0])) > 0 {
				pmItem := Swagger2PathToPostman2ApiItem(cfg, swag, url, "DELETE", path.Delete)
				pmFolderName := strings.TrimSpace(path.Delete.Tags[0])
				pmFolder := pman.GetOrNewFolder(pmFolderName)
				pmFolder.Item = append(pmFolder.Item, pmItem)
				pman.SetFolder(pmFolder)
			}
		}
	}

	return pman
}

func Swagger2PathToPostman2ApiItem(cfg Configuration, swag swagger2.Specification, url string, method string, endpoint swagger2.Endpoint) postman2.ApiItem {
	item := postman2.ApiItem{}

	item.Name = endpoint.Summary

	item.Request = postman2.Request{
		Method: strings.ToUpper(method)}

	item.Request.URL = BuildPostmanURL(cfg, swag, url, endpoint)

	headers := []postman2.Header{}

	if len(endpoint.Consumes) > 0 {
		if len(strings.TrimSpace(endpoint.Consumes[0])) > 0 {
			headers = append(headers, postman2.Header{
				Key:   "Content-Type",
				Value: strings.TrimSpace(endpoint.Consumes[0])})
		}
	}
	if len(endpoint.Produces) > 0 {
		if len(strings.TrimSpace(endpoint.Produces[0])) > 0 {
			headers = append(headers, postman2.Header{
				Key:   "Accept",
				Value: strings.TrimSpace(endpoint.Consumes[0])})
		}
	}
	for _, header := range cfg.PostmanHeaders {
		headers = append(headers, header)
	}

	item.Request.Header = headers

	return item
}

func BuildPostmanURL(cfg Configuration, swag swagger2.Specification, swaggerUrl string, endpoint swagger2.Endpoint) postman2.URL {
	URLParts := []string{}

	// Set URL path parts
	if len(strings.TrimSpace(cfg.PostmanURLHostname)) > 0 {
		URLParts = append(URLParts, strings.TrimSpace(cfg.PostmanURLHostname))
	} else if len(strings.TrimSpace(swag.Host)) > 0 {
		URLParts = append(URLParts, strings.TrimSpace(swag.Host))
	}

	if len(strings.TrimSpace(swag.BasePath)) > 0 {
		URLParts = append(URLParts, strings.TrimSpace(swag.BasePath))
	}
	if len(strings.TrimSpace(swaggerUrl)) > 0 {
		URLParts = append(URLParts, strings.TrimSpace(swaggerUrl))
	}

	// Create URL
	rawPostmanUrl := strings.TrimSpace(strings.Join(URLParts, "/"))
	rx1 := regexp.MustCompile(`/+`)
	rawPostmanUrl = rx1.ReplaceAllString(rawPostmanUrl, "/")
	rx2 := regexp.MustCompile(`^/+`)
	rawPostmanUrl = rx2.ReplaceAllString(rawPostmanUrl, "")

	// Add URL Scheme
	if len(swag.Schemes) > 0 {
		for _, scheme := range swag.Schemes {
			if len(strings.TrimSpace(scheme)) > 0 {
				rawPostmanUrl = strings.Join([]string{scheme, rawPostmanUrl}, "://")
				break
			}
		}
	}

	rx3 := regexp.MustCompile(`(^|[^\{])\{([^\/\{\}]+)\}([^\}]|$)`)
	rawPostmanUrl = rx3.ReplaceAllString(rawPostmanUrl, "$1:$2$3")

	postmanUrl := postman2.NewURL(rawPostmanUrl)

	// Set Default URL Path Parameters
	rx4 := regexp.MustCompile(`^\s*(:(.+))\s*$`)

	for _, part := range postmanUrl.Path {
		rs4 := rx4.FindAllStringSubmatch(part, -1)
		if len(rs4) > 0 {
			baseVariable := rs4[0][2]
			//defaultValue := ""
			var defaultValue interface{}
			for _, parameter := range endpoint.Parameters {
				if parameter.Name == baseVariable {
					defaultValue = parameter.Default
					/*
						if len(strings.TrimSpace(parameter.Default)) > 0 {
							defaultValue = strings.TrimSpace(parameter.Default)
						}
					*/
					break
				}
			}
			postmanUrl.AddVariable(baseVariable, defaultValue)
		}
	}

	return postmanUrl
}
