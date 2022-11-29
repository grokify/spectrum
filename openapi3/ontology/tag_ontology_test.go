package ontology

import (
	"fmt"
	"strings"
	"testing"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/text/stringcase"
)

var onologyTestStruct = Ontology{
	OperationIDCase:          stringcase.KebabCase,
	SchemaNameCase:           stringcase.SnakeCase,
	SchemaNameReponseSuffix:  "response",
	PathVarCase:              stringcase.SnakeCase,
	PathIDPrefix:             "",
	PathIDSuffix:             "id",
	SpecFileCase:             stringcase.SnakeCase,
	SpecFilePrefix:           "spec",
	SpecFileExt:              ".yaml",
	SpecFileResourceIsPlural: false,
}

var tagOnologyTests = []struct {
	tagOntology     TagOnology
	tagOntologyData TagOnologyData
}{
	{
		tagOntology: TagOnology{
			ResourceNameSingular:      "pet animal",
			ResourceNameSingularTitle: "Pet Animal",
			ResourceNameSingularShort: "pet",
			ResourceNamePlural:        "pet animals",
			ResourceNamePluralTitle:   "Pet Animals",
			DeterminerSinglar:         "a",
			DeterminerPlural:          "",
		},
		tagOntologyData: TagOnologyData{
			Tag:                        "Pet Animals",
			CreateOperationID:          "create-pet-animal",
			ReadOperationID:            "get-pet-animal",
			UpdateOperationID:          "update-pet-animal",
			DeleteOperationID:          "delete-pet-animal",
			ListOperationID:            "list-pet-animals",
			CreateSummary:              "Create a pet animal",
			ReadSummary:                "Describe a pet animal",
			UpdateSummary:              "Update a pet animal",
			DeleteSummary:              "Delete a pet animal",
			ListSummary:                "List pet animals",
			ResourcePathVar:            "pet_id",
			ResourceSchemaNameRequest:  "pet_animal",
			ResourceSchemaNameResponse: "pet_animal_response",
			SpecFilename:               "spec_pet_animal.yaml",
		},
	},
}

// TestParseJSONPointer ensures the `ParseJSONPointer` is working properly.
func TestTagOnologyTests(t *testing.T) {
	for _, tt := range tagOnologyTests {
		to := tt.tagOntology
		to.Ontology = onologyTestStruct
		gotData := to.Data()
		wantData := tt.tagOntologyData
		if gotData.CreateSummary != wantData.CreateSummary {
			fmt.Printf("MISMATCH_SUMMARIES:\nwant [%s]\nrecv [%s]\n", wantData.CreateSummary, gotData.CreateSummary)
		}
		mismatchFields, isEqual := gotData.Equal(wantData)
		if !isEqual {
			fmtutil.PrintJSON(gotData)
			fmtutil.PrintJSON(wantData)
			t.Errorf("ontology.TagOnologyData.Data(...) Mistmatch fields [%s] want [%v] got [%v]",
				strings.Join(mismatchFields, ","),
				wantData, gotData)
		}
	}
}
