package ontology

import (
	"strings"

	"github.com/grokify/mogo/text/stringcase"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	ActionCreate                            = "Create"
	ActionRead                              = "Get"
	ActionUpdate                            = "Update"
	ActionDelete                            = "Delete"
	ActionList                              = "List"
	ActionCreateDescription                 = ActionCreate
	ActionReadDescription                   = "Describe"
	ActionUpdateDescription                 = ActionUpdate
	ActionDeleteDescription                 = ActionDelete
	ActionListDescription                   = ActionList
	ResourceSchemaNameResponseSuffixDefault = "response"
	DefaultCase                             = stringcase.CamelCase
)

var DefaultCaseToFunc = stringcase.ToCamelCase

type Ontology struct {
	OperationIDCase          string
	SchemaNameCase           string
	SchemaNameReponseSuffix  string
	PathVarCase              string
	PathIDPrefix             string
	PathIDSuffix             string
	SpecFileCase             string
	SpecFilePrefix           string
	SpecFileSuffix           string
	SpecFileExt              string
	SpecFileResourceIsPlural bool
}

type TagOnology struct {
	Ontology
	ResourceNameSingular      string
	ResourceNameSingularTitle string
	ResourceNameSingularShort string
	ResourceNamePlural        string
	ResourceNamePluralTitle   string
	DeterminerSinglar         string
	DeterminerPlural          string
}

type TagOnologyData struct {
	Tag                        string
	CreateOperationID          string
	ReadOperationID            string
	UpdateOperationID          string
	DeleteOperationID          string
	ListOperationID            string
	CreateSummary              string
	ReadSummary                string
	UpdateSummary              string
	DeleteSummary              string
	ListSummary                string
	ResourcePathVar            string
	ResourceSchemaNameRequest  string
	ResourceSchemaNameResponse string
	SpecFilename               string
}

func (to *TagOnology) Data() TagOnologyData {
	return TagOnologyData{
		Tag:                        to.ResourceNamePluralTitle,
		CreateSummary:              to.CreateSummary(),
		ReadSummary:                to.ReadSummary(),
		UpdateSummary:              to.UpdateSummary(),
		DeleteSummary:              to.DeleteSummary(),
		ListSummary:                to.ListSummary(),
		CreateOperationID:          to.CreateOperationID(),
		ReadOperationID:            to.ReadOperationID(),
		UpdateOperationID:          to.UpdateOperationID(),
		DeleteOperationID:          to.DeleteOperationID(),
		ListOperationID:            to.ListOperationID(),
		ResourcePathVar:            to.ResourcePathVar(),
		ResourceSchemaNameRequest:  to.ResourceSchemaNameRequest(),
		ResourceSchemaNameResponse: to.ResourceSchemaNameResponse(),
		SpecFilename:               to.SpecFilename()}
}

func (tod *TagOnologyData) Equal(data TagOnologyData) ([]string, bool) {
	mismatchFields := []string{}
	if tod.Tag != data.Tag {
		mismatchFields = append(mismatchFields, "Tag")
	}
	if tod.CreateSummary != data.CreateSummary {
		mismatchFields = append(mismatchFields, "CreateSummary")
	}
	if tod.ReadSummary != data.ReadSummary {
		mismatchFields = append(mismatchFields, "ReadSummary")
	}
	if tod.UpdateSummary != data.UpdateSummary {
		mismatchFields = append(mismatchFields, "UpdateSummary")
	}
	if tod.DeleteSummary != data.DeleteSummary {
		mismatchFields = append(mismatchFields, "DeleteSummary")
	}
	if tod.ListSummary != data.ListSummary {
		mismatchFields = append(mismatchFields, "ListSummary")
	}
	if tod.CreateOperationID != data.CreateOperationID {
		mismatchFields = append(mismatchFields, "CreateOperationID")
	}
	if tod.ReadOperationID != data.ReadOperationID {
		mismatchFields = append(mismatchFields, "ReadOperationID")
	}
	if tod.UpdateOperationID != data.UpdateOperationID {
		mismatchFields = append(mismatchFields, "UpdateOperationID")
	}
	if tod.DeleteOperationID != data.DeleteOperationID {
		mismatchFields = append(mismatchFields, "DeleteOperationID")
	}
	if tod.ListOperationID != data.ListOperationID {
		mismatchFields = append(mismatchFields, "ListOperationID")
	}
	if tod.ResourcePathVar != data.ResourcePathVar {
		mismatchFields = append(mismatchFields, "ResourcePathVar")
	}
	if tod.ResourceSchemaNameRequest != data.ResourceSchemaNameRequest {
		mismatchFields = append(mismatchFields, "ResourceSchemaNameRequest")
	}
	if tod.ResourceSchemaNameResponse != data.ResourceSchemaNameResponse {
		mismatchFields = append(mismatchFields, "ResourceSchemaNameResponse")
	}
	if tod.SpecFilename != data.SpecFilename {
		mismatchFields = append(mismatchFields, "SpecFilename")
	}
	return mismatchFields, len(mismatchFields) == 0
}

func (to *TagOnology) SpecFilename() string {
	casefunc := stringcase.FuncToWantCaseOrDefault(to.Ontology.SpecFileCase, DefaultCaseToFunc)
	resourceName := to.ResourceNameSingular
	if to.Ontology.SpecFileResourceIsPlural {
		resourceName = to.ResourceNamePlural
	}
	filename := resourceName
	if len(to.Ontology.SpecFilePrefix) > 0 {
		filename = to.Ontology.SpecFilePrefix + " " + filename
	}
	if len(to.Ontology.SpecFileSuffix) > 0 {
		filename = filename + " " + to.Ontology.SpecFileSuffix
	}
	return casefunc(filename) + to.Ontology.SpecFileExt
}

func (to *TagOnology) ResourcePathVar() string {
	resourceNameSingle := strings.TrimSpace(to.ResourceNameSingularShort)
	if len(resourceNameSingle) == 0 {
		resourceNameSingle = to.ResourceNameSingular
	}
	casefunc := stringcase.FuncToWantCaseOrDefault(to.Ontology.PathVarCase, DefaultCaseToFunc)
	return casefunc(to.Ontology.PathIDPrefix + " " + resourceNameSingle + " " + to.Ontology.PathIDSuffix)
}

func (to *TagOnology) ResourceSchemaNameRequest() string {
	casefunc := stringcase.FuncToWantCaseOrDefault(to.Ontology.SchemaNameCase, DefaultCaseToFunc)
	return casefunc(to.ResourceNameSingular)
}

func (to *TagOnology) ResourceSchemaNameResponse() string {
	casefunc := stringcase.FuncToWantCaseOrDefault(to.Ontology.SchemaNameCase, DefaultCaseToFunc)
	return casefunc(to.ResourceNameSingular + " " + to.Ontology.SchemaNameReponseSuffix)
}

func (to *TagOnology) CreateSummary() string { return to.ActionSummary(ActionCreateDescription, false) }
func (to *TagOnology) ReadSummary() string   { return to.ActionSummary(ActionReadDescription, false) }
func (to *TagOnology) UpdateSummary() string { return to.ActionSummary(ActionUpdateDescription, false) }
func (to *TagOnology) DeleteSummary() string { return to.ActionSummary(ActionDeleteDescription, false) }
func (to *TagOnology) ListSummary() string   { return to.ActionSummary(ActionListDescription, true) }

func (to *TagOnology) ActionSummary(action string, plural bool) string {
	if plural {
		return stringsutil.CondenseSpace(strings.Join(
			[]string{action, to.DeterminerPlural, to.ResourceNamePlural}, " "))
	}
	return stringsutil.CondenseSpace(strings.Join(
		[]string{action, to.DeterminerSinglar, to.ResourceNameSingular}, " "))
}

func (to *TagOnology) CreateOperationID() string { return to.ActionOperationID(ActionCreate, false) }
func (to *TagOnology) ReadOperationID() string   { return to.ActionOperationID(ActionRead, false) }
func (to *TagOnology) UpdateOperationID() string { return to.ActionOperationID(ActionUpdate, false) }
func (to *TagOnology) DeleteOperationID() string { return to.ActionOperationID(ActionDelete, false) }
func (to *TagOnology) ListOperationID() string   { return to.ActionOperationID(ActionList, true) }

func (to *TagOnology) ActionOperationID(action string, plural bool) string {
	xcasefuncOpID := stringcase.FuncToWantCaseOrNoOp(to.Ontology.OperationIDCase)
	if plural {
		return xcasefuncOpID(strings.Join([]string{action, to.ResourceNamePlural}, " "))
	}
	return xcasefuncOpID(strings.Join([]string{action, to.ResourceNameSingular}, " "))
}
