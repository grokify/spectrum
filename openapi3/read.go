package openapi3

import (
	"fmt"
	"io/ioutil"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

// ReadFile does optional validation which is useful when
// merging incomplete spec files.
func ReadFile(file string, validate bool) (*oas3.Swagger, error) {
	if validate {
		return ReadFileLoader(file)
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ReadFile.ReadFile.Error.Filename [%v]", file))
	}
	swag := &oas3.Swagger{}
	err = swag.UnmarshalJSON(bytes)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("ReadFile.UnmarshalJSON.Error.Filename [%v]", file))
	}
	return swag, nil
}

func ReadFileLoader(file string) (*oas3.Swagger, error) {
	return oas3.NewSwaggerLoader().LoadSwaggerFromFile(file)
}

func ValidateSpec(specfile string) (bool, error) {
	bytes, err := ioutil.ReadFile(specfile)
	if err != nil {
		return false, errors.Wrap(err, "E_OPENAPI3_SPEC_READ_ERROR")
	}
	_, err = oas3.NewSwaggerLoader().LoadSwaggerFromData(bytes)
	if err != nil {
		return false, errors.Wrap(err, "E_OPENAPI3_SPEC_LOAD_VALIDATE_ERROR")
	}
	return true, nil
}
