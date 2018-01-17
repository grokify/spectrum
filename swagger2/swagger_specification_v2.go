package swagger2

import (
	"encoding/json"
	"io/ioutil"
)

// Specification represents a Swagger 2.0 specification.
type Specification struct {
	Swagger                        string                         `json:"swagger,omitempty"`
	Host                           string                         `json:"host,omitempty"`
	Info                           *Info                          `json:"info,omitempty"`
	BasePath                       string                         `json:"basePath,omitempty"`
	Schemes                        []string                       `json:"schemes,omitempty"`
	Paths                          map[string]Path                `json:"paths,omitempty"`
	Definitions                    map[string]Definition          `json:"definitions,omitempty"`
	XAmazonApigatewayDocumentation XAmazonApigatewayDocumentation `json:"x-amazon-apigateway-documentation,omitempty"`
}

// NewSpecificationFromBytes returns a Swagger Specification from a byte array.
func NewSpecificationFromBytes(data []byte) (Specification, error) {
	spec := Specification{}
	err := json.Unmarshal(data, &spec)
	return spec, err
}

// ReadSwagger2Spec returns a Swagger Specification from a filepath.
func ReadSwagger2Spec(filepath string) (Specification, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Specification{}, err
	}
	return NewSpecificationFromBytes(bytes)
}

// Info represents a Swagger 2.0 spec info object.
type Info struct {
	Description    string `json:"description,omitempty"`
	Version        string `json:"version,omitempty"`
	Title          string `json:"title,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty"`
}

// Path represents a Swagger 2.0 spec path object.
type Path struct {
	Get    *Endpoint `json:"get,omitempty"`
	Post   *Endpoint `json:"post,omitempty"`
	Put    *Endpoint `json:"put,omitempty"`
	Delete *Endpoint `json:"delete,omitempty"`
	Ref    string    `json:"$ref,omitempty"`
}

// Endpoint represents a Swagger 2.0 spec endpoint object.
type Endpoint struct {
	Tags                         []string                     `json:"tags,omitempty"`
	Summary                      string                       `json:"summary,omitempty"`
	OperationID                  string                       `json:"operationId,omitempty"`
	Description                  string                       `json:"description,omitempty"`
	Consumes                     []string                     `json:"consumes,omitempty"`
	Produces                     []string                     `json:"produces,omitempty"`
	Parameters                   []Parameter                  `json:"parameters,omitempty"`
	Responses                    map[string]Response          `json:"responses,omitempty"`
	XAmazonApigatewayIntegration XAmazonApigatewayIntegration `json:"x-amazon-apigateway-integration,omitempty"`
}

type Response struct {
	Description string            `json:"description,omitempty"`
	Schema      Schema            `json:"schema,omitempty"`
	Headers     map[string]Header `json:"headers,omitempty"`
}

type Schema struct {
	Ref string `json:"$ref,omitempty"`
}

type Header struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type Definition struct {
	Type       string              `json:"type,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type Property struct {
	Description string `json:"description,omitempty"`
	Format      string `json:"format,omitempty"`
	Items       Items  `json:"items,omitempty"`
	Type        string `json:"type,omitempty"`
	Ref         string `json:"$ref,omitempty"`
}

type Items struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}

type XAmazonApigatewayIntegration struct {
	Responses           map[string]XAmazonApigatewayIntegrationResponse `json:"responses,omitempty"`
	PassthroughBehavior string                                          `json:"passthroughBehavior,omitempty"`
	RequestTemplates    map[string]string                               `json:"requestTemplates,omitempty"`
	Type                string                                          `json:"type,omitempty"`
}

type XAmazonApigatewayIntegrationResponse struct {
	StatusCode         string            `json:"statusCode,omitempty"`
	ResponseParameters map[string]string `json:"responseParameters,omitempty"`
	ResponseTemplates  map[string]string `json:"responseTemplates,omitempty"`
}

type XAmazonApigatewayDocumentation struct {
	Version            string              `json:"version,omitempty"`
	CreatedDate        string              `json:"createdDate,omitempty"`
	DocumentationParts []DocumentationPart `json:"documentationParts,omitempty"`
}

type DocumentationPart struct {
	Location   XAmazonApigatewayDocumentationPartLocation   `json:"location,omitempty"`
	Properties XAmazonApigatewayDocumentationPartProperties `json:"properties,omitempty"`
}

type XAmazonApigatewayDocumentationPartLocation struct {
	Type       string `json:"type,omitempty"`
	Method     string `json:"method,omitempty"`
	Path       string `json:"path,omitempty"`
	StatusCode string `json:"statusCode,omitempty"`
	Name       string `json:"name,omitempty"`
}

type XAmazonApigatewayDocumentationPartProperties struct {
	Tags        []string                                `json:"tags,omitempty"`
	Summary     string                                  `json:"summary,omitempty"`
	Description string                                  `json:"description,omitempty"`
	Info        *XAmazonApigatewayDocumentationPartInfo `json:"info,omitempty"`
}

type XAmazonApigatewayDocumentationPartInfo struct {
	Description string `json:"description,omitempty"`
}

// Parameter represents a Swagger 2.0 spec parameter object.
type Parameter struct {
	Name        string      `json:"name,omitempty"`
	Type        string      `json:"type,omitempty"`
	In          string      `json:"in,omitempty"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}
