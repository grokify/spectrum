package swagger2

import (
	"encoding/json"
)

type Specification struct {
	Host     string          `json:"host,omitempty"`
	Info     Info            `json:"info,omitempty"`
	BasePath string          `json:"basePath,omitempty"`
	Schemes  []string        `json:"schemes,omitempty"`
	Paths    map[string]Path `json:"paths,omitempty"`
}

func NewSpecificationFromBytes(data []byte) (Specification, error) {
	spec := Specification{}
	err := json.Unmarshal(data, &spec)
	return spec, err
}

type Info struct {
	Description    string `json:"description,omitempty"`
	Version        string `json:"version,omitempty"`
	Title          string `json:"title,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty"`
}

type Path struct {
	Get Endpoint `json:"get,omitempty"`
}

type Endpoint struct {
	Tags        []string    `json:"tags,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	OperationId string      `json:"operationId,omitempmty"`
	Description string      `json:"description,omitempty"`
	Consumes    []string    `json:"consumes,omitempty"`
	Produces    []string    `json:"produces,omitempty"`
	Parameters  []Parameter `json:"parameters"`
}

type Parameter struct {
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	In          string `json:"in,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Default     string `json:"default,omitempty"`
}
