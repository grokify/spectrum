package openapi2

import (
	"encoding/json"
	"io/ioutil"
	"os"

	oas2 "github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"gopkg.in/yaml.v2"
)

func ConvertSwaggmanOAS2ToKinOAS2(smSpec *Specification) (*oas2.Swagger, error) {
	bytes, err := json.Marshal(smSpec)
	if err != nil {
		return nil, err
	}
	var kinSpec oas2.Swagger
	err = json.Unmarshal(bytes, &kinSpec)
	return &kinSpec, err
}

func ConvertOAS2FileToOAS3File(oas2file, oas3file string, perm os.FileMode, pretty bool) error {
	oas2, err := ReadOpenAPI2KinSpecFile(oas2file)
	if err != nil {
		return err
	}
	oas3, err := openapi2conv.ToV3Swagger(oas2)
	if err != nil {
		return err
	}
	if FilenameIsYAML(oas3file) {
		bytes, err := yaml.Marshal(oas3)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(oas3file, bytes, perm)
	}
	bytes, err := oas3.MarshalJSON()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(oas3file, bytes, perm)
}
