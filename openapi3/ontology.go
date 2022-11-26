package openapi3

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
	OperationIDs   map[string][]string `json:"operationIDs"`
	SchemaNames    []string            `json:"schemaNames"`
	ParameterNames map[string][]string `json:"parameterNames"`
}

func NewOntology() Ontology {
	return Ontology{
		OperationIDs:   map[string][]string{},
		SchemaNames:    []string{},
		ParameterNames: map[string][]string{}}
}
