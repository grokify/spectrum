package swagger2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/sjson"
)

// Specification represents a Swagger 2.0 specification.
type Specification struct {
	Swagger     string                `json:"swagger,omitempty"`
	Host        string                `json:"host,omitempty"`
	Info        *Info                 `json:"info,omitempty"`
	BasePath    string                `json:"basePath,omitempty"`
	Schemes     []string              `json:"schemes,omitempty"`
	Tags        []Tag                 `json:"tags,omitempty"`
	Produces    []string              `json:"produces,omitempty"`
	Consumes    []string              `json:"consumes,omitempty"`
	Paths       map[string]Path       `json:"paths,omitempty"`
	Definitions map[string]Definition `json:"definitions,omitempty"`
}

// NewSpecificationFromBytes returns a Swagger Specification from a byte array.
func NewSpecificationFromBytes(data []byte) (Specification, error) {
	spec := Specification{}
	err := json.Unmarshal(data, &spec)
	if err != nil {
		return spec, err
	}

	if len(spec.Schemes) == 0 {
		spec.Schemes = []string{"http"}
	}

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

// Tag represents a Swagger 2.0 spec tag object.
type Tag struct {
	Name         string        `json:"name,omitempty"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// ExternalDocs represents a Swagger 2.0 spec tag object. The
// URL property is required.
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// Path represents a Swagger 2.0 spec path object.
type Path struct {
	Get    *Endpoint `json:"get,omitempty"`
	Patch  *Endpoint `json:"patch,omitempty"`
	Post   *Endpoint `json:"post,omitempty"`
	Put    *Endpoint `json:"put,omitempty"`
	Delete *Endpoint `json:"delete,omitempty"`
	Ref    string    `json:"$ref,omitempty"`
}

// HasMethodWithTag returns true if the swagger path supports the specified method
func (p *Path) HasMethodWithTag(method string) bool {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case http.MethodGet:
		if p.Get != nil && len(p.Get.Tags) > 0 && len(strings.TrimSpace(p.Get.Tags[0])) > 0 {
			return true
		}
	case http.MethodPatch:
		if p.Patch != nil && len(p.Patch.Tags) > 0 && len(strings.TrimSpace(p.Patch.Tags[0])) > 0 {
			return true
		}
	case http.MethodPost:
		if p.Post != nil && len(p.Post.Tags) > 0 && len(strings.TrimSpace(p.Post.Tags[0])) > 0 {
			return true
		}
	case http.MethodPut:
		if p.Put != nil && len(p.Put.Tags) > 0 && len(strings.TrimSpace(p.Put.Tags[0])) > 0 {
			return true
		}
	case http.MethodDelete:
		if p.Delete != nil && len(p.Delete.Tags) > 0 && len(strings.TrimSpace(p.Delete.Tags[0])) > 0 {
			return true
		}
	}
	return false
}

// SetEndpoint sets the method endpoint
func (p *Path) SetEndpoint(method string, endpoint Endpoint) error {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case http.MethodGet:
		p.Get = &endpoint
	case http.MethodPost:
		p.Post = &endpoint
	case http.MethodPut:
		p.Put = &endpoint
	case http.MethodPatch:
		p.Patch = &endpoint
	case http.MethodDelete:
		p.Delete = &endpoint
	default:
		return fmt.Errorf("method [%v] not supported", method)
	}
	return nil
}

// Endpoint represents a Swagger 2.0 spec endpoint object.
type Endpoint struct {
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
	Description string              `json:"description,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"`
	Produces    []string            `json:"produces,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

// IsEmpty returns true if all endpoint fields are zero value
func (ep *Endpoint) IsEmpty() bool {
	if len(ep.Tags) > 0 ||
		len(ep.Summary) > 0 ||
		len(ep.OperationID) > 0 ||
		len(ep.Description) > 0 ||
		len(ep.Consumes) > 0 ||
		len(ep.Produces) > 0 ||
		len(ep.Parameters) > 0 ||
		len(ep.Responses) > 0 {
		return false
	}
	return true
}

// Response is a swagger response
type Response struct {
	Description string            `json:"description,omitempty"`
	Schema      *Schema           `json:"schema,omitempty"`
	Headers     map[string]Header `json:"headers,omitempty"`
	Examples    map[string]string `json:"examples,omitempty"`
}

// Schema is a swagger schema
type Schema struct {
	Ref string `json:"$ref,omitempty"`
}

// Header is a swagger header
type Header struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

// Definition is a swagger type definition
type Definition struct {
	Type       string              `json:"type,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

// Property is a swagger property
type Property struct {
	Description string      `json:"description,omitempty"`
	Format      string      `json:"format,omitempty"`
	Items       *Items      `json:"items,omitempty"`
	Type        string      `json:"type,omitempty"`
	Ref         string      `json:"$ref,omitempty"`
	Example     interface{} `json:"example,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
}

// Items are swagger type items
type Items struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}

// IsEmpty returns true if all item field are zero value
func (items *Items) IsEmpty() bool {
	if len(strings.TrimSpace(items.Type)) == 0 && len(strings.TrimSpace(items.Ref)) == 0 {
		return true
	}
	return false
}

// Parameter represents a Swagger 2.0 spec parameter object.
type Parameter struct {
	Name             string            `json:"name,omitempty"`
	Type             string            `json:"type,omitempty"`
	In               string            `json:"in,omitempty"`
	Description      string            `json:"description,omitempty"`
	Schema           *Schema           `json:"schema,omitempty"`
	Required         bool              `json:"required,omitempty"`
	CollectionFormat string            `json:"collectionFormat,omitempty"`
	Items            *Items            `json:"items,omitempty"`
	Default          interface{}       `json:"default,omitempty"`
	XExamples        map[string]string `json:"x-examples,omitempty"`
}

// GetJSONBodyFromDefinition returns a POST body example from the named swagger definition
func GetJSONBodyFromDefinition(name string, defs map[string]Definition) (string, error) {

	def := defs[name]
	var body string
	var err error

	if def.Properties == nil {
		switch def.Type {
		case "string":
			return `"string"`, nil
		}
	}

	for pName, prop := range def.Properties {

		if prop.Ref != "" {
			defName := prop.Ref[strings.LastIndex(prop.Ref, "/")+1:]
			b, err := GetJSONBodyFromDefinition(defName, defs)
			if err != nil {
				return "", err
			}
			body, err = sjson.SetRaw(body, pName, b)
			if err != nil {
				return "", err
			}
		} else if prop.Example != nil {
			body, err = sjson.Set(body, pName, prop.Example)
			if err != nil {
				return "", err
			}
		} else {

			switch prop.Type {
			case "string":
				if prop.Enum != nil {
					body, err = sjson.Set(body, pName, prop.Enum[0])
				} else {
					body, err = sjson.Set(body, pName, getExampleString(prop.Format))
				}
				if err != nil {
					return "", err
				}
			case "boolean":
				body, err = sjson.Set(body, pName, true)
				if err != nil {
					return "", err
				}
			case "number":
				if prop.Format == "double" {
					body, err = sjson.Set(body, pName, 3.14159)
				} else {
					body, err = sjson.Set(body, pName, 1202)
				}
				if err != nil {
					return "", err
				}
			case "integer":
				body, err = sjson.Set(body, pName, 1201)
				if err != nil {
					return "", err
				}
			case "object":
				body, err = sjson.SetRaw(body, pName, "{}")
				if err != nil {
					return "", err
				}
			case "array":
				if prop.Items != nil && !prop.Items.IsEmpty() {
					if prop.Items.Type != "" {
						switch prop.Items.Type {
						case "string":
							body, err = sjson.SetRaw(body, pName, `["aString", "aString"]`)
						default:
							body, err = sjson.SetRaw(body, pName, "[12345, 67890]")
						}
						if err != nil {
							return "", err
						}
					} else if prop.Items.Ref != "" {
						defName := prop.Items.Ref[strings.LastIndex(prop.Items.Ref, "/")+1:]
						b, err := GetJSONBodyFromDefinition(defName, defs)
						if err != nil {
							return "", err
						}
						body, err = sjson.SetRaw(body, pName, "["+b+"]")
						if err != nil {
							return "", err
						}
					}
				}

			default:
				body, err = sjson.Set(body, pName, "schema type not provided")
				if err != nil {
					return "", err
				}
			}
		}
	}

	return body, nil
}

func getExampleString(format string) string {
	switch format {
	case "date":
		now := time.Now()
		return now.Format("2006-01-02")
	case "date-time":
		now := time.Now()
		return now.Format(time.RFC3339)
	case "uri":
		return "http://www.ietf.org/rfc/rfc2396.txt"
	case "email":
		return "user@example.com"
	default:
		return "aString"
	}
}
