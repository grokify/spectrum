package swagger2

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/grokify/gotilla/io/ioutilmore"
	"gopkg.in/yaml.v2"
)

func ConvertSwaggmanOAS2ToKinOAS2(smSpec *Specification) (*openapi2.Swagger, error) {
	bytes, err := json.Marshal(smSpec)
	if err != nil {
		return nil, err
	}
	var kinSpec openapi2.Swagger
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
	indent := ""
	if pretty {
		indent = "  "
	}
	return ioutilmore.WriteFileJSON(oas3file, oas3, perm, "", indent)
}
