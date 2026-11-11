// oas3filter writes the part of an OpenAPI 3 spec (JSON or YAML) that one
// audience may see, as marked by x-audience (internal, partner or public)
// or x-internal: true, with only the components, security schemes and tags
// the kept operations use.
//
// Usage:
//
//	oas3filter -i openapi.json -o openapi-public.json -a public
//
// Operations with no mark are internal unless --default says otherwise, so
// an unmarked operation is never published by mistake. The output is JSON,
// or YAML when the output file ends in .yaml or .yml, and is the same for
// the same input.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/grokify/spectrum/openapi3"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Input          string `short:"i" long:"input" description:"Input OpenAPI 3 spec file" required:"true"`
	Output         string `short:"o" long:"output" description:"Output spec file (.json, .yaml or .yml)" required:"true"`
	Audience       string `short:"a" long:"audience" description:"Audience to write: internal, partner or public" required:"true"`
	Default        string `long:"default" description:"Audience of operations with no x-audience or x-internal" default:"internal"`
	KeepExtensions bool   `long:"keep-extensions" description:"Keep x-audience and x-internal in the output"`
}

var rxYAML = regexp.MustCompile(`(?i)\.ya?ml$`)

func main() {
	var opts Options
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}
	if err := run(opts); err != nil {
		slog.Error("oas3filter", "error", err)
		os.Exit(2)
	}
}

func run(opts Options) error {
	reader, err := openapi3.ParseAudience(opts.Audience)
	if err != nil {
		return err
	}
	def, err := openapi3.ParseAudience(opts.Default)
	if err != nil {
		return fmt.Errorf("--default: %w", err)
	}
	spec, err := openapi3.ReadFile(opts.Input, true)
	if err != nil {
		return err
	}
	keep, err := openapi3.AudienceFilter(reader, def)
	if err != nil {
		return err
	}
	fopts := &openapi3.FilterOptions{}
	if !opts.KeepExtensions {
		fopts.RemoveExtensions = []string{openapi3.XAudience, openapi3.XInternal}
	}
	sm := openapi3.SpecMore{Spec: spec}
	out, err := sm.ExportByFilter(keep, fopts)
	if err != nil {
		return err
	}
	om := openapi3.SpecMore{Spec: out}
	if rxYAML.MatchString(opts.Output) {
		err = om.WriteFileYAML(opts.Output, 0o600)
	} else {
		err = om.WriteFileJSON(opts.Output, 0o600, "", "  ")
	}
	if err != nil {
		return err
	}
	slog.Info("wrote filtered spec",
		"audience", reader,
		"output", opts.Output,
		"operations", om.OperationsCount(),
		"of", sm.OperationsCount(),
		"schemas", om.SchemasCount(),
		"of_schemas", sm.SchemasCount())
	return nil
}
