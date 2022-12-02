package taxonomy

import (
	"errors"
	"fmt"
	"os"

	mapslice "github.com/ake-persson/mapslice-json"
	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/spectrum/openapi3"
	"sigs.k8s.io/yaml"
)

const XTaxonomy = "x-taxonomy"

type Category struct {
	Key         string     `json:"-"`
	Slug        string     `json:"slug,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Tags        []oas3.Tag `json:"tags,omitempty" yaml:"tags,omitempty"`
}

func (cat *Category) Taxonomy(catsFilepath string) (Taxonomy, error) {
	if len(cat.Key) == 0 {
		return Taxonomy{}, errors.New("categgory `Key` not set")
	}
	return Taxonomy{
		Category: CategoryRef{Ref: fmt.Sprintf("%s#/%s", catsFilepath, cat.Key)},
		Slug:     cat.Slug}, nil
}

type Categories []Category

var ErrCategoryNotFound = errors.New("category not found")

// SpecAddXTaxonomy expects `catTitle` to be the same as the OAS3 `tag` name.
func (cats *Categories) SpecAddXTaxonomy(spec *openapi3.Spec, catTitle, catsFilepath string) error {
	if spec == nil {
		return openapi3.ErrSpecNotSet
	}
	cat, err := cats.Category(catTitle)
	if err != nil {
		return errorsutil.Wrapf(ErrCategoryNotFound, "category (%s)", catTitle)
	}
	tax, err := cat.Taxonomy(catsFilepath)
	if err != nil {
		return err
	}
	tax.AddToSpec(spec)
	return nil
}

func (cats *Categories) Category(title string) (Category, error) {
	for _, cat := range *cats {
		if cat.Title == title {
			return cat, nil
		}
	}
	return Category{}, errorsutil.Wrapf(ErrCategoryNotFound, "category (%s)", title)
}

func (cats *Categories) MapSlice() mapslice.MapSlice {
	ms := mapslice.MapSlice{}
	for _, cat := range *cats {
		ms = append(ms, mapslice.MapItem{Key: cat.Key, Value: cat})
	}
	return ms
}

func (cats *Categories) MarshalJSON(prefix, indent string) ([]byte, error) {
	return jsonutil.MarshalSimple(cats.MapSlice(), prefix, indent)
}

func (cats *Categories) MarshalYAML() ([]byte, error) {
	if jbytes, err := cats.MarshalJSON("", ""); err != nil {
		return []byte{}, err
	} else {
		return yaml.JSONToYAML(jbytes)
	}
}

func (cats *Categories) WriteFileJSON(filename string, perm os.FileMode, prefix, indent string) error {
	if jsonData, err := cats.MarshalJSON(prefix, indent); err != nil {
		return err
	} else {
		return os.WriteFile(filename, jsonData, perm)
	}
}

func (cats *Categories) WriteFileYAML(filename string, perm os.FileMode) error {
	if ybytes, err := cats.MarshalYAML(); err != nil {
		return err
	} else {
		return os.WriteFile(filename, ybytes, perm)
	}
}
