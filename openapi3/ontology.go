package openapi3

import oas3 "github.com/getkin/kin-openapi/openapi3"

// OntologySet is a set of ontologies which can be used to understand
// ontologies by a key, such as filename or by tag.
type OntologySet struct {
	Ontologies map[string]Ontology `json:"ontologies"`
}

// Ontology returns the naming structure of an OpenAPI Spec. It is useful
// for understanding the naming conventions of an existing OpenAPI Spec.
// For example, the relationship of operationIDs to paths and the relationship
// of parameter name component keys to parameter names.
type Ontology struct {
	Operations  map[string]OperationMeta `json:"operationIDs"`
	SchemaNames []string                 `json:"schemaNames"`
	Parameters  oas3.ParametersMap       `json:"parameters"`
}

func NewOntology() Ontology {
	return Ontology{
		Operations:  map[string]OperationMeta{},
		SchemaNames: []string{},
		Parameters:  oas3.ParametersMap{}}
}
