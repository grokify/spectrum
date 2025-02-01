package openapi2

import (
	"encoding/json"
	"os"

	"github.com/getkin/kin-openapi/openapi2conv"
	yaml "gopkg.in/yaml.v3"
)

func ConvertSpectrumOAS2ToKinOAS2(smSpec *Specification) (*Spec, error) {
	bytes, err := json.Marshal(smSpec)
	if err != nil {
		return nil, err
	}
	var kinSpec Spec
	err = json.Unmarshal(bytes, &kinSpec)
	return &kinSpec, err
}

func ConvertOAS2FileToOAS3File(oas2file, oas3file string, perm os.FileMode, pretty bool) error {
	oas2, err := ReadOpenAPI2KinSpecFile(oas2file)
	if err != nil {
		return err
	}
	oas3, err := openapi2conv.ToV3(oas2)
	if err != nil {
		return err
	}
	if FilenameIsYAML(oas3file) {
		bytes, err := yaml.Marshal(oas3)
		if err != nil {
			return err
		}
		return os.WriteFile(oas3file, bytes, perm)
	}
	if bytes, err := oas3.MarshalJSON(); err != nil {
		return err
	} else {
		return os.WriteFile(oas3file, bytes, perm)
	}
}
