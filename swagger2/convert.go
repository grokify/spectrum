package swagger2

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi2"
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
