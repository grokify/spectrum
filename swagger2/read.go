package swagger2

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/ghodss/yaml"
)

func ReadOpenAPI2SpecFile(filename string) (*Specification, error) {
	var spec Specification
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return &spec, err
	}
	rx := regexp.MustCompile(`.ya?ml$`)
	if rx.MatchString(strings.ToLower(strings.TrimSpace(filename))) {
		err = yaml.Unmarshal(bytes, &spec)
	} else {
		err = json.Unmarshal(bytes, &spec)
	}
	return &spec, err
}

func ReadOpenAPI2KinSpecFile(filename string) (*openapi2.Swagger, error) {
	var swag openapi2.Swagger
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return &swag, err
	}
	rx := regexp.MustCompile(`.ya?ml$`)
	if rx.MatchString(strings.ToLower(strings.TrimSpace(filename))) {
		err = yaml.Unmarshal(bytes, &swag)
	} else {
		err = json.Unmarshal(bytes, &swag)
	}
	return &swag, err
}
