package main

import (
	"fmt"

	"github.com/grokify/swagger2postman-go"
	"github.com/grokify/swagger2postman-go/postman2"
)

func main() {
	swagSpecFilepath := "ringcentral.swagger2.basic.json"
	pmanBaseFilepath := "ringcentral.postman2.base.json"
	pmanSpecFilepath := "ringcentral.postman2.basic.json"

	cfg := swagger2postman.Configuration{
		PostmanURLHostname: "{{RC_SERVER_HOSTNAME}}",
		PostmanHeaders: []postman2.Header{postman2.Header{
			Key:   "Authorization",
			Value: "Bearer {{my_access_token}}"}}}

	conv := swagger2postman.NewConverter(cfg)

	merge := true
	var err error

	if merge {
		err = conv.MergeConvert(swagSpecFilepath, pmanBaseFilepath, pmanSpecFilepath)
	} else {
		err = conv.Convert(swagSpecFilepath, pmanSpecFilepath)
	}

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Wrote %v\n", pmanSpecFilepath)
	}

	fmt.Println("DONE")
}
