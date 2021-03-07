package openapi3

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/grokify/simplego/encoding/jsonutil"
	"github.com/pkg/errors"
)

var rxYamlExtension = regexp.MustCompile(`(?i)\.ya?ml\s*$`)

// ReadFile does optional validation which is useful when
// merging incomplete spec files.
func ReadFile(oas3file string, validate bool) (*oas3.Swagger, error) {
	if validate {
		return ReadAndValidateFile(oas3file)
	}
	bytes, err := ioutil.ReadFile(oas3file)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ReadFile.ReadFile.Error.Filename [%v]", oas3file))
	}
	if rxYamlExtension.MatchString(oas3file) {
		bytes, err = yaml.YAMLToJSON(bytes)
		if err != nil {
			return nil, err
		}
	}
	spec := &oas3.Swagger{}
	err = spec.UnmarshalJSON(bytes)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ReadFile.UnmarshalJSON.Error.Filename [%s]", oas3file))
	}
	return spec, nil
}

// Parse will parse a byte array to an `*oas3.Swagger` struct.
// It will use JSON first. If unsuccessful, it will attempt to
// parse it as YAML.
func Parse(oas3Bytes []byte) (*oas3.Swagger, error) {
	spec := &oas3.Swagger{}
	err := spec.UnmarshalJSON(oas3Bytes)
	if err != nil {
		bytes, err2 := yaml.YAMLToJSON(oas3Bytes)
		if err2 != nil {
			return spec, err
		}
		spec = &oas3.Swagger{}
		err3 := spec.UnmarshalJSON(bytes)
		return spec, err3
	}
	return spec, err
}

func ReadAndValidateFile(oas3file string) (*oas3.Swagger, error) {
	bytes, err := ioutil.ReadFile(oas3file)
	if err != nil {
		return nil, errors.Wrap(err, "E_READ_FILE_ERROR")
	}
	spec, err := oas3.NewSwaggerLoader().LoadSwaggerFromData(bytes)
	if err != nil {
		return spec, errors.Wrap(err, fmt.Sprintf("E_OPENAPI3_SPEC_LOAD_VALIDATE_ERROR [%s]", oas3file))
	}
	_, err = ValidateMore(spec)
	return spec, err
}

type ValidationStatus struct {
	Status  bool
	Message string
	Context string
	OpenAPI string
}

/*
status: false
message: |-
  expected Object {
    title: 'Medium API',
    description: 'Articles that matter on social publishing platform'
  } to have key version
  	missing keys: version
context: "#/info"
openapi: 3.0.0
*/

func ValidateMore(spec *oas3.Swagger) (ValidationStatus, error) {
	vs := ValidationStatus{OpenAPI: "3.0.0"}
	version := strings.TrimSpace(spec.Info.Version)
	if len(version) == 0 {
		jdata, err := jsonutil.MarshalSimple(spec.Info, "", "  ")
		if err != nil {
			return vs, err
		}
		vs := ValidationStatus{
			Context: "#/info",
			Message: fmt.Sprintf("expect Object %s to have key version\nmissing keys:version", string(jdata)),
			OpenAPI: "3.0.0"}
		return vs, fmt.Errorf("E_OPENAPI3_MISSING_KEY [%s]", "info/version")
	}
	vs.Status = true
	return vs, nil
}

func Copy(spec *oas3.Swagger) (*oas3.Swagger, error) {
	bytes, err := spec.MarshalJSON()
	if err != nil {
		return nil, err
	}
	loader := oas3.NewSwaggerLoader()
	return loader.LoadSwaggerFromData(bytes)
}
