package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/swaggman/openapi3"
	"github.com/grokify/swaggman/openapi3/openapi3html"
)

func main() {
	file := "spec_ringcentral_openapi3.yaml"

	spec, err := openapi3.ReadFile(file, true)
	if err != nil {
		log.Fatal(err)
	}

	pageParams := openapi3html.PageParams{
		PageTitle:  spec.Info.Title,
		PageLink:   "https://developers.ringcentral.com",
		TableDomID: "apitable"}
	pageParams.AddSpec(spec, ColumnTexts())

	pageHTML := openapi3html.SwaggmanUIPage(pageParams)

	filename := "api-regisry.html"
	err = ioutil.WriteFile(filename, []byte(pageHTML), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%s]\n", filename)

	fmt.Println("DONE")
}

func ColumnTexts() *table.ColumnSet {
	columns := []table.Column{
		{
			Display: "Method",
			Slug:    "method"},
		{
			Display: "Path",
			Slug:    "path"},
		{
			Display: "OperationID",
			Slug:    "operationId"},
		{
			Display: "Summary",
			Slug:    "summary"},
		{
			Display: "Tags",
			Slug:    "tags"},
		{
			Display: "API Group",
			Slug:    "x-api-group"},
		{
			Display: "Throttling",
			Slug:    "x-throttling-group"},
		{
			Display: "App Permission",
			Slug:    "x-app-permission"},
		{
			Display: "User Permissions",
			Slug:    "x-user-permission"},
		{
			Display: "Docs Level",
			Slug:    "x-docs-level"},
	}
	return &table.ColumnSet{Columns: columns}
}
