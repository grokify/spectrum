package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/grokify/swagger2postman-go"
	"github.com/grokify/swagger2postman-go/postman2"
	"github.com/grokify/swagger2postman-go/swagger2"

	"github.com/grokify/gotilla/fmt/fmtutil"
)

func getSwagger2Spec(filepath string) (swagger2.Specification, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return swagger2.Specification{}, err
	}
	return swagger2.NewSpecificationFromBytes(bytes)
}

func getPostman2BaseSpec(filepath string) (postman2.Collection, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return postman2.Collection{}, err
	}
	return postman2.NewCollectionFromBytes(bytes)
}

func main() {
	swagSpecFilepath := "ringcentral.swagger2.basic.json"
	pmanBaseFilepath := "ringcentral.postman2.base.json"
	pmanSpecFilepath := "ringcentral.postman2.basic.json"

	cfg := swagger2postman.Configuration{
		PostmanURLHostname: "{{RC_SERVER_HOSTNAME}}",
		PostmanHeaders: []postman2.Header{postman2.Header{
			Key:   "Authorization",
			Value: "Bearer {{myAccessToken}}"}}}

	swag, err := getSwagger2Spec(swagSpecFilepath)
	if err != nil {
		panic(err)
	}

	pman, err := getPostman2BaseSpec(pmanBaseFilepath)
	if err != nil {
		panic(err)
	}
	pman.InflateRawURLs()

	pm := swagger2postman.Merge(cfg, pman, swag)
	if 1 == 0 {
		fmtutil.PrintJSON(pm)
	}
	bytes, err := json.MarshalIndent(pm, "", "  ")
	if err == nil {
		ioutil.WriteFile(pmanSpecFilepath, bytes, 0644)
	}

	fmt.Println("DONE")
}
