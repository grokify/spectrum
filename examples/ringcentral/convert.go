package main

import (
	"fmt"

	"github.com/grokify/swaggman"
	"github.com/grokify/swaggman/postman2"
)

func main() {
	swagSpecFilepath := "ringcentral.swagger2.basic.json"
	pmanBaseFilepath := "ringcentral.postman2.base.json"
	pmanSpecFilepath := "ringcentral.postman2.basic.json"

	cfg := swaggman.Configuration{
		PostmanURLHostname: "{{RC_SERVER_HOSTNAME}}",
		PostmanHeaders: []postman2.Header{postman2.Header{
			Key:   "Authorization",
			Value: "Bearer {{my_access_token}}"}}}

	conv := swaggman.NewConverter(cfg)

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
