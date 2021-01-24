package topostman2

import (
	"github.com/grokify/swaggman/postman2"
)

//const DefaultContentTypePreferences string = `multipart/form-data,application/json,application/x-www-form-urlencoded,application/xml,text/plain`
//var defaultContentTypePreferencesSlice = strings.Split(DefaultContentTypePreferences, ",")

// Configuration is a Swaggman configuration that holds information on how
// to create the Postman 2.0 collection including overriding Swagger 2.0
// spec values.
type Configuration struct {
	// PostmanServerURLBasePath supports setting the base path as an environment variable
	// such as {{MY_API_BASE_URL}}
	PostmanServerURLBasePath string            `json:"postmanServerUrlApiBasePath,omitempty"`
	PostmanServerURL         string            `json:"postmanServerUrl,omitempty"`
	PostmanURLHostname       string            `json:"postmanURLHostname,omitempty"`
	PostmanHeaders           []postman2.Header `json:"postmanHeaders,omitempty"`
	UseXTagGroups            bool              `json:"useXTagGroups,omitempty"`
	RequestBodyFunc          func(urlPath string) string
}
