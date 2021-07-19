package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/grokify/gocharts/data/table/tabulator"
	"github.com/grokify/spectrum/openapi3"
	"github.com/grokify/spectrum/openapi3/openapi3html"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	OpenAPISpec string `short:"o" long:"openapispec" description:"Input OpenAPI Spec File" required:"true"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	spec, err := openapi3.ReadFile(opts.OpenAPISpec, true)
	if err != nil {
		log.Fatal(err)
	}

	pageParams := openapi3html.PageParams{
		PageTitle:  spec.Info.Title,
		PageLink:   "https://developers.ringcentral.com",
		TableDomID: "apitable",
		ColumnSet:  ColumnTexts()}
	pageParams.AddSpec(spec)

	pageHTML := openapi3html.SpectrumUIPage(pageParams)

	filename := "api-regisry.html"
	err = ioutil.WriteFile(filename, []byte(pageHTML), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%s]\n", filename)

	fmt.Println("DONE")
}

func ColumnTexts() *tabulator.ColumnSet {
	columns := []tabulator.Column{
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
	return &tabulator.ColumnSet{Columns: columns}
}
