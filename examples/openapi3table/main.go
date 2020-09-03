package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/grokify/swaggman/openapi3"
	"github.com/grokify/swaggman/openapi3/tohtml"
)

func main() {
	file := "spec_ringcentral_openapi3.yaml"

	spec, err := openapi3.ReadFile(file, true)
	if err != nil {
		log.Fatal(err)
	}

	pageParams := tohtml.PageParams{
		PageTitle:  spec.Info.Title,
		PageLink:   "https://developers.ringcentral.com",
		TableDomID: "apitable"}
	pageParams.AddSpec(spec)

	pageHTML := tohtml.SwaggmanUIPage(pageParams)

	filename := "api-regisry.html"
	err = ioutil.WriteFile(filename, []byte(pageHTML), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%s]\n", filename)

	fmt.Println("DONE")
}
