package swagger2postman

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grokify/swagger2postman-go/postman2"
	"github.com/grokify/swagger2postman-go/swagger2"
)

type Configuration struct {
	PostmanURLHostname string            `json:"postmanURLHostname,omitempty"`
	PostmanHeaders     []postman2.Header `json:"postmanHeaders,omitempty"`
}

type Converter struct {
	Swagger swagger2.Specification
}

func Convert(swag swagger2.Specification) postman2.Collection {
	cfg := Configuration{}

	pman := postman2.Collection{
		Info: postman2.CollectionInfo{
			Name:        swag.Info.Title,
			Description: swag.Info.Description,
			Schema:      "https://schema.getpostman.com/json/collection/v2.0.0/collection.json"}}

	for url, path := range swag.Paths {
		if url != "/v1.0/account/{accountId}/extension" {
			continue
		}
		fmt.Println(url)

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

func Merge(cfg Configuration, pman postman2.Collection, swag swagger2.Specification) postman2.Collection {
	if len(pman.Info.Name) == 0 {
		pman.Info.Name = swag.Info.Title
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
			defaultValue := ""
			for _, parameter := range endpoint.Parameters {
				if parameter.Name == baseVariable {
					if len(strings.TrimSpace(parameter.Default)) > 0 {
						defaultValue = strings.TrimSpace(parameter.Default)
					}
					break
				}
			}
			postmanUrl.AddVariable(baseVariable, defaultValue)
		}
	}

	return postmanUrl
}
