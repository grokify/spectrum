package openapi3

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/table/tabulator"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/http/pathmethod"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/type/stringsutil"
	"golang.org/x/exp/slices"
	"sigs.k8s.io/yaml"
)

var ErrSpecNotSet = errors.New("spec not set")

type Spec = oas3.T

type Operation = oas3.Operation

type SpecMore struct {
	Spec *Spec
}

func ReadSpecMore(path string, validate bool) (*SpecMore, error) {
	spec, err := ReadFile(path, validate)
	if err != nil {
		return nil, err
	}
	return &SpecMore{Spec: spec}, nil
}

func (sm *SpecMore) Clone() (*Spec, error) {
	if sm.Spec == nil {
		return nil, nil
	}
	bytes, err := sm.Spec.MarshalJSON()
	if err != nil {
		return nil, err
	}
	loader := oas3.NewLoader()
	return loader.LoadFromData(bytes)
}

func (sm *SpecMore) SchemasCount() int {
	if sm.Spec == nil {
		return -1
	} else if sm.Spec.Components.Schemas == nil {
		return 0
	}
	return len(sm.Spec.Components.Schemas)
}

func (sm *SpecMore) OperationsTable(columns *tabulator.ColumnSet, filterFunc func(path, method string, op *oas3.Operation) bool, addlColFuncs *OperationMoreStringFuncMap) (*table.Table, error) {
	return operationsTable(sm.Spec, columns, filterFunc, addlColFuncs)
}

func operationsTable(spec *Spec, columns *tabulator.ColumnSet, filterFunc func(path, method string, op *oas3.Operation) bool, addlColFuncs *OperationMoreStringFuncMap) (*table.Table, error) {
	if columns == nil {
		columns = OpTableColumnsDefault(false)
	}
	title := ""
	if spec.Info != nil {
		title = spec.Info.Title
	}
	tbl := table.NewTable(title)
	tbl.Columns = columns.DisplayTexts()

	// specMore := SpecMore{Spec: spec}
	// tgs, err := specMore.TagGroups()
	// if err != nil {
	// 	return nil, err
	// }

	VisitOperations(spec, func(path, method string, op *oas3.Operation) {
		if filterFunc != nil && !filterFunc(path, method, op) {
			return
		}
		row := []string{}

		for _, text := range columns.Columns {
			switch text.Slug {
			case "tags":
				row = append(row, strings.Join(op.Tags, ", "))
			case "method":
				row = append(row, method)
			case "path":
				row = append(row, path)
			case "operationId":
				row = append(row, op.OperationID)
			case "summary":
				row = append(row, op.Summary)
			case "description":
				row = append(row, op.Description)
			// case XTagGroups:
			//	row = append(row, strings.Join(
			//		tgs.GetTagGroupNamesForTagNames(op.Tags...), ", "))
			case "securityScopes":
				om := OperationMore{Operation: op}
				row = append(row, strings.Join(om.SecurityScopes(false), ", "))
			case XThrottlingGroup:
				// row = append(row, GetExtensionPropStringOrEmpty(op.ExtensionProps, XThrottlingGroup))
				row = append(row, GetExtensionPropStringOrEmpty(op.Extensions, XThrottlingGroup))
			case "docsURL":
				if op.ExternalDocs != nil {
					row = append(row, op.ExternalDocs.URL)
				}
			default:
				if addlColFuncs != nil {
					colFunc := addlColFuncs.Func(text.Slug)
					// for XTagGroups send OperationMoreStringFuncMap[XTagGroups] = tgs.
					if colFunc != nil {
						row = append(row, colFunc(&OperationMore{
							Path:      path,
							Method:    method,
							Operation: op}))
						continue
					}
				}
				// row = append(row, GetExtensionPropStringOrEmpty(op.ExtensionProps, text.Slug))
				row = append(row, GetExtensionPropStringOrEmpty(op.Extensions, text.Slug))
			}
		}

		tbl.Rows = append(tbl.Rows, row)
	})
	return &tbl, nil
}

func OpTableColumnsDefault(inclDocsURL bool) *tabulator.ColumnSet {
	cols := []tabulator.Column{
		{
			Display: "Tags",
			Slug:    "tags",
			Width:   150},
		{
			Display: "Method",
			Slug:    "method",
			Width:   70},
		{
			Display: "Path",
			Slug:    "path",
			Width:   800},
		{
			Display: "OperationID",
			Slug:    "operationId",
			Width:   150},
		{
			Display: "Summary",
			Slug:    "summary",
			Width:   150},
		{
			Display: "SecurityScopes",
			Slug:    "securityScopes",
			Width:   150},
		{
			Display: "XThrottlingGroup",
			Slug:    XThrottlingGroup,
			Width:   150},
		{
			Display: "DocsURL",
			Slug:    "docsURL",
			Width:   150},
	}
	/*
		if inclDocsURL {
			cols = append(cols, tabulator.Column{
				Display: "DocsURL",
				Slug:    "docsURL",
				Width:   150})
		}
	*/
	return &tabulator.ColumnSet{Columns: cols}
}

func OpTableColumnsRingCentral() *tabulator.ColumnSet {
	columns := OpTableColumnsDefault(false)
	rcCols := []tabulator.Column{
		{
			Display: "API Group",
			Slug:    "x-api-group",
			Width:   150},
		{
			Display: "Throttling",
			Slug:    "x-throttling-group",
			Width:   150},
		{
			Display: "App Permission",
			Slug:    "x-app-permission",
			Width:   150},
		{
			Display: "User Permissions",
			Slug:    "x-user-permission",
			Width:   150},
	}
	columns.Columns = append(columns.Columns, rcCols...)
	return columns
	//return &table.ColumnSet{Columns: columns}
}

func (sm *SpecMore) Operations(inclTags []string) *OperationMores {
	if sm.Spec == nil {
		return nil
	}
	// func QueryOperationsByTags(spec *openapi3.Spec, tags []string) *OperationEditSet {
	tagsWantMatch := map[string]int{}
	for _, tag := range inclTags {
		tagsWantMatch[tag] = 1
	}
	// opmSet := &OperationMoreSet{OperationMores: []OperationMore{}}

	oms := &OperationMores{}

	VisitOperations(sm.Spec, func(path, method string, op *oas3.Operation) {
		if len(tagsWantMatch) == 0 {
			*oms = append(*oms,
				OperationMore{
					Path:      path,
					Method:    method,
					Operation: op})
			return
		}
		for _, opTagTry := range op.Tags {
			if _, ok := tagsWantMatch[opTagTry]; ok {
				*oms = append(*oms,
					OperationMore{
						Path:      path,
						Method:    method,
						Operation: op})
				return
			}
		}
	})

	return oms
}

func (sm *SpecMore) OperationMetasMap(inclTags []string) map[string]OperationMeta {
	oms := sm.OperationMetas(inclTags)
	omsMap := map[string]OperationMeta{}
	for _, om := range oms {
		omsMap[om.PathMethod()] = om
	}
	return omsMap
}

func (sm *SpecMore) OperationMetas(inclTags []string) []OperationMeta {
	if sm.Spec == nil {
		return []OperationMeta{}
	}
	oms := []*OperationMeta{}
	pathsMap := sm.Spec.Paths.Map()
	for url, path := range pathsMap {
		// for url, path := range sm.Spec.Paths {
		oms = append(oms,
			OperationToMeta(url, http.MethodConnect, path.Connect, inclTags),
			OperationToMeta(url, http.MethodDelete, path.Delete, inclTags),
			OperationToMeta(url, http.MethodGet, path.Get, inclTags),
			OperationToMeta(url, http.MethodHead, path.Head, inclTags),
			OperationToMeta(url, http.MethodOptions, path.Options, inclTags),
			OperationToMeta(url, http.MethodPatch, path.Patch, inclTags),
			OperationToMeta(url, http.MethodPost, path.Post, inclTags),
			OperationToMeta(url, http.MethodPut, path.Put, inclTags),
			OperationToMeta(url, http.MethodTrace, path.Trace, inclTags))
	}

	oms2 := []OperationMeta{}
	for _, om := range oms {
		if om != nil {
			oms2 = append(oms2, *om)
		}
	}
	return oms2
}

func (sm *SpecMore) OperationsCount() int {
	if sm.Spec == nil {
		return -1
	}
	count := 0
	VisitOperations(sm.Spec, func(skipPath, skipMethod string, op *oas3.Operation) {
		count++
	})
	return count
}

// OperationCountsByTag returns a histogram for operations by tag.
func (sm *SpecMore) OperationCountsByTag() *histogram.Histogram {
	hist := histogram.NewHistogram("Operation Counts by Tag")
	hist.Bins = sm.TagsMap(&TagsOpts{InclOps: true})
	hist.Inflate()
	return hist
}

func (sm *SpecMore) OperationIDs() []string {
	ids := []string{}
	VisitOperations(sm.Spec, func(path, method string, op *oas3.Operation) {
		if op == nil {
			return
		}
		ids = append(ids, op.OperationID)
	})
	return stringsutil.SliceCondenseSpace(ids, false, true)
}

func (sm *SpecMore) OperationIDsCounts() map[string]int {
	msi := map[string]int{}
	VisitOperations(sm.Spec, func(skipPath, skipMethod string, op *oas3.Operation) {
		msi[op.OperationID]++
	})
	return msi
}

// OperationIDsLocations returns a `map[string][]string` where the keys are
// operationIDs and the values are pathMethods for use in analyzing if there are
// duplicate operationIDs.
func (sm *SpecMore) OperationIDsLocations() map[string][]string {
	vals := map[string][]string{}
	VisitOperations(sm.Spec, func(opPath, opMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		pathMethod := pathmethod.PathMethod(opPath, opMethod)
		op.OperationID = strings.TrimSpace(op.OperationID)
		vals[op.OperationID] = append(vals[op.OperationID], pathMethod)
	})
	return vals
}

func (sm *SpecMore) OperationByID(wantOperationID string) (path, method string, op *oas3.Operation, err error) {
	wantOperationID = strings.TrimSpace(wantOperationID)
	VisitOperations(sm.Spec, func(thisPath, thisMethod string, thisOp *oas3.Operation) {
		if thisOp == nil {
			return
		}
		if wantOperationID == strings.TrimSpace(thisOp.OperationID) {
			path = thisPath
			method = thisMethod
			op = thisOp
		}
	})
	if len(strings.TrimSpace(method)) == 0 {
		err = fmt.Errorf("operation_not_found [%s]", wantOperationID)
	}
	return path, method, op, err
}

var (
	ErrPathNotFound      = errors.New("path not found")
	ErrOperationNotFound = errors.New("operation not found")
)

func (sm *SpecMore) OperationByPathMethod(path, method string) (*oas3.Operation, error) {
	method = strings.ToUpper(strings.TrimSpace(method))
	_, err := httputilmore.ParseHTTPMethod(method)
	if err != nil {
		return nil, err
	}

	pathItem := sm.Spec.Paths.Find(path)
	if pathItem == nil {
		return nil, nil
	}
	// pathItem, ok := sm.Spec.Paths[path]
	// if !ok {
	// 	return nil, nil
	// }

	return pathItem.GetOperation(method), nil
}

func (sm *SpecMore) SetOperation(path, method string, op *oas3.Operation) {
	path = strings.TrimSpace(path)
	if strings.Index(path, "/") != 0 {
		path = "/" + path
	}
	if sm.Spec.Paths == nil {
		// sm.Spec.Paths = map[string]*oas3.PathItem{}
		sm.Spec.Paths = oas3.NewPaths()
	}
	pathItem := sm.Spec.Paths.Find(path)
	if pathItem == nil {
		pathItem = &oas3.PathItem{}
	}
	/*
		// code here is for getkin v0.121.0, broken in v0.122.0
		pathItem, ok := sm.Spec.Paths[path]
		if !ok {
			pathItem = &oas3.PathItem{}
		}
	*/
	method = strings.ToUpper(strings.TrimSpace(method))
	switch method {
	case http.MethodConnect:
		pathItem.Connect = op
	case http.MethodDelete:
		pathItem.Delete = op
	case http.MethodGet:
		pathItem.Get = op
	case http.MethodHead:
		pathItem.Head = op
	case http.MethodOptions:
		pathItem.Options = op
	case http.MethodPatch:
		pathItem.Patch = op
	case http.MethodPost:
		pathItem.Post = op
	case http.MethodPut:
		pathItem.Put = op
	case http.MethodTrace:
		pathItem.Trace = op
	}
	// sm.Spec.Paths[path] = pathItem // code here is for getkin v0.121.0, broken in v0.122.0
	sm.Spec.Paths.Set(path, pathItem)
}

// Ontology returns a populated `Ontology` struct for the spec. If no spec
// is supplied, a zero value is returned.
func (sm *SpecMore) Ontology() Ontology {
	return Ontology{
		Operations:  sm.OperationMetasMap([]string{}),
		Parameters:  sm.Spec.Components.Parameters,
		SchemaNames: sm.SchemaNames()}
}

// ParameterNames returns a `map[string][]string` where they key is the
// key in `#/components/parameters` and the values are both references
// and names. There should only be either a reference or a name but this
// structure allows capture of both.
func (sm *SpecMore) ParameterNames() map[string][]string {
	mss := map[string][]string{}
	if sm.Spec == nil {
		return mss
	}
	for paramKey, paramRef := range sm.Spec.Components.Parameters {
		if _, ok := mss[paramKey]; !ok {
			mss[paramKey] = []string{}
		}
		if len(paramRef.Ref) > 0 {
			mss[paramKey] = append(mss[paramKey], paramRef.Ref)
		}
		if paramRef.Value == nil {
			continue
		}
		if len(paramRef.Value.Name) > 0 {
			mss[paramKey] = append(mss[paramKey], paramRef.Value.Name)
		}
	}
	return mss
}

func (sm *SpecMore) SchemaNames() []string {
	schemaNames := []string{}
	for schemaName := range sm.Spec.Components.Schemas {
		schemaNames = append(schemaNames, schemaName)
	}
	return stringsutil.SliceCondenseSpace(schemaNames, true, true)
}

var rxSchemas = regexp.MustCompile(`"([^"]*#/components/schemas/([^"]+))"`)

func (sm *SpecMore) SchemaPointers(dedupe bool) ([]string, []string, error) {
	bytes, err := sm.MarshalJSON("", "")
	if err != nil {
		return []string{}, []string{}, err
	}
	pointers := []string{}
	names := []string{}
	m := rxSchemas.FindAllStringSubmatch(string(bytes), -1)
	for _, mx := range m {
		if len(mx) == 3 {
			pointers = append(pointers, mx[1])
			names = append(names, mx[2])
		}
	}
	return stringsutil.SliceCondenseSpace(pointers, dedupe, true),
		stringsutil.SliceCondenseSpace(names, dedupe, true),
		nil
}

func (sm *SpecMore) SchemaNamesStatus() (schemaNoReference, both, referenceNoSchema []string, err error) {
	haveNames := sm.SchemaNames()
	_, havePointers, err := sm.SchemaPointers(true)
	if err != nil {
		return
	}
	schemaNoReference, both, referenceNoSchema = stringsutil.SlicesCompare(haveNames, havePointers)
	return
}

func (sm *SpecMore) SchemaNameExists(schemaName string, includeNil bool) bool {
	for schemaNameTry, schemaRef := range sm.Spec.Components.Schemas {
		if schemaNameTry == schemaName {
			if includeNil {
				return true
			} else if schemaRef == nil {
				return false
			}
			schemaRef.Ref = strings.TrimSpace(schemaRef.Ref)
			if len(schemaRef.Ref) > 0 {
				return true
			}
			if schemaRef.Value == nil {
				return false
			} else {
				return true
			}
		}
	}
	return false
}

// SchemaRef returns a top level `SchemaRef` under `Components` based on
// map name or JSON pointer. It returns `nil` if the `schemaName` is not
// found.
func (sm *SpecMore) SchemaRef(schemaName string) *oas3.SchemaRef {
	if sm.Spec == nil {
		return nil
	}
	if strings.Contains(schemaName, PointerComponentsSchemas) {
		ptr, err := ParseJSONPointer(schemaName)
		if err != nil {
			return nil
		}
		schNameTry, ok := ptr.IsTopSchema()
		if !ok {
			return nil
		}
		schemaName = schNameTry
	}

	if schRef, ok := sm.Spec.Components.Schemas[schemaName]; ok {
		return schRef
	}
	return nil
}

func (sm *SpecMore) SchemaRefSet(schemaName string, schemaRef *oas3.SchemaRef) error {
	schemaName = strings.TrimSpace(schemaName)
	if schemaRef != nil {
		if sm.Spec.Components.Schemas == nil {
			sm.Spec.Components.Schemas = map[string]*oas3.SchemaRef{}
		}
		if schemaRef.Value != nil {
			if 1 == 0 && len(schemaRef.Value.Description) == 0 {
				return fmt.Errorf("no description for schema component [%s]", schemaName)
			}
		}
	}
	sm.Spec.Components.Schemas[schemaName] = schemaRef
	return nil
}

// ServerURL returns the OAS3 Spec URL for the index
// specified.
func (sm *SpecMore) ServerURL(index uint) string {
	if int(index)+1 > len(sm.Spec.Servers) {
		return ""
	}
	server := sm.Spec.Servers[index]
	return strings.TrimSpace(server.URL)
}

// ServerURLBasePath extracts the base path from a OAS URL which can include variables.
func (sm *SpecMore) ServerURLBasePath(index uint) (string, error) {
	serverURL := sm.ServerURL(index)
	if len(serverURL) == 0 {
		return "", nil
	}
	serverURLParsed, err := urlutil.ParseURLTemplate(serverURL)
	if err != nil {
		return "", err
	}
	return serverURLParsed.Path, nil
}

// OperationsDescriptionInfo returns information on operations with and without descriptions.
func (sm *SpecMore) OperationsDescriptionInfo() map[string][]string {
	data := map[string][]string{
		"opWithDesc":      {},
		"opWoutDesc":      {},
		"opWithDescCount": {},
		"opWoutDescCount": {},
	}
	VisitOperations(sm.Spec, func(opPath, opMethod string, op *oas3.Operation) {
		if op == nil {
			return
		}
		pathMethod := pathmethod.PathMethod(opPath, opMethod)
		op.Description = strings.TrimSpace(op.Description)
		if len(op.Description) == 0 {
			data["opWoutDesc"] = append(data["opWoutDesc"], pathMethod)
		} else {
			data["opWithDesc"] = append(data["opWithDesc"], pathMethod)
		}
	})
	data["opWithDescCount"] = append(data["opWithDescCount"], strconv.Itoa(len(data["opWithDesc"])))
	data["opWoutDescCount"] = append(data["opWoutDescCount"], strconv.Itoa(len(data["opWoutDesc"])))
	return data
}

func (sm *SpecMore) SpecTagStats() SpecTagStats {
	stats := SpecTagStats{
		TagStats:      SpecTagCounts{},
		TagsAll:       sm.Tags(&TagsOpts{InclDefs: true, InclOps: true}),
		TagsDefs:      sm.Tags(&TagsOpts{InclDefs: true, InclOps: false}),
		TagsOps:       sm.Tags(&TagsOpts{InclDefs: false, InclOps: true}),
		TagCountsAll:  sm.TagsMap(&TagsOpts{InclDefs: true, InclOps: true}),
		TagCountsDefs: sm.TagsMap(&TagsOpts{InclDefs: true, InclOps: false}),
		TagCountsOps:  sm.TagsMap(&TagsOpts{InclDefs: false, InclOps: true}),
	}
	VisitOperations(sm.Spec, func(skipPath, skipMethod string, op *oas3.Operation) {
		op.Tags = stringsutil.SliceCondenseSpace(op.Tags, true, true)
		stats.TagStats.OpsTotal++
		if len(op.Tags) > 0 {
			stats.TagStats.OpsWithTags++
		} else {
			stats.TagStats.OpsWithoutTags++
		}
	})
	return stats
}

// TagsOpts represents additional settings for tag read functions.
type TagsOpts struct {
	InclDefs       bool
	InclOps        bool
	OpsTagsJoin    bool // default is false.
	OpsTagsJoinSep string
}

func TagsOptsDefault() *TagsOpts {
	return &TagsOpts{
		InclDefs: true,
		InclOps:  true,
	}
}

func (sm *SpecMore) Tags(opts *TagsOpts) []string {
	tags := []string{}
	tagsMap := sm.TagsMap(opts)
	for tag := range tagsMap {
		tags = append(tags, tag)
	}
	return stringsutil.SliceCondenseSpace(tags, true, true)
}

// TagsValidate checks to see if the tag names in the Spec tags property and operations match.
func (sm *SpecMore) TagsValidate() bool {
	return slices.Equal(
		sm.Tags(&TagsOpts{InclDefs: true}),
		sm.Tags(&TagsOpts{InclOps: true}))
}

// TagsMap returns a set of operations with tags present in the current spec.
func (sm *SpecMore) TagsMap(opts *TagsOpts) map[string]int {
	tagsMap := map[string]int{}
	if opts == nil {
		opts = TagsOptsDefault()
	}
	if opts.InclDefs {
		for _, tag := range sm.Spec.Tags {
			tagName := strings.TrimSpace(tag.Name)
			if len(tagName) > 0 {
				if _, ok := tagsMap[tagName]; !ok {
					tagsMap[tagName] = 0 // don't increment unless present in operations
				}
			}
		}
	}
	if opts.InclOps {
		VisitOperations(sm.Spec, func(skipPath, skipMethod string, op *oas3.Operation) {
			if op == nil {
				return
			}
			if opts.OpsTagsJoin {
				tags := strings.Join(op.Tags, opts.OpsTagsJoinSep)
				tagsMap[tags]++
			} else {
				for _, tagName := range op.Tags {
					tagName = strings.TrimSpace(tagName)
					if len(tagName) > 0 {
						if _, ok := tagsMap[tagName]; !ok {
							tagsMap[tagName] = 0
						}
						tagsMap[tagName]++
					}
				}
			}
		})
	}
	return tagsMap
}

type SpecStats struct {
	OperationsCount int
	SchemasCount    int
}

// SpecTagStats represents tags data for tag definitions defined at the root level of a spec
// and tags referenced in operations.
type SpecTagStats struct {
	TagStats      SpecTagCounts
	TagsAll       []string
	TagsDefs      []string
	TagsOps       []string
	TagCountsAll  map[string]int
	TagCountsDefs map[string]int
	TagCountsOps  map[string]int
}

type SpecTagCounts struct {
	OpsWithTags    int
	OpsWithoutTags int
	OpsTotal       int
}

func (sm *SpecMore) Stats() SpecStats {
	return SpecStats{
		OperationsCount: sm.OperationsCount(),
		SchemasCount:    sm.SchemasCount(),
	}
}

func (sm *SpecMore) MarshalJSON(prefix, indent string) ([]byte, error) {
	bytes, err := sm.Spec.MarshalJSON()
	if err != nil {
		return bytes, err
	}

	if len(prefix) > 0 || len(indent) > 0 {
		return jsonutil.IndentBytes(bytes, prefix, indent)
	}

	return bytes, nil
}

func (sm *SpecMore) MarshalYAML() ([]byte, error) {
	if jbytes, err := sm.MarshalJSON("", ""); err != nil {
		return []byte{}, err
	} else {
		return yaml.JSONToYAML(jbytes)
	}
}

func (sm *SpecMore) PrintJSON(prefix, indent string) error {
	if bytes, err := sm.MarshalJSON(prefix, indent); err != nil {
		return err
	} else {
		_, err = fmt.Println(string(bytes))
		return err
	}
}

func (sm *SpecMore) WriteFileCSV(filename string, addlColFuncs *OperationMoreStringFuncMap) error {
	if tbl, err := sm.OperationsTable(nil, nil, addlColFuncs); err != nil {
		return err
	} else {
		return tbl.WriteCSV(filename)
	}
}

func (sm *SpecMore) WriteFileJSON(filename string, perm os.FileMode, prefix, indent string) error {
	if jsonData, err := sm.MarshalJSON(prefix, indent); err != nil {
		return err
	} else {
		return os.WriteFile(filename, jsonData, perm)
	}
}

// WriteFileXLSX writes the spec in XLSX Open XML format. If supplied, the `filterFunc` must return `true` for
// an operation to be included in the file. if non-nil, `addlColFuncs` is used to add additional columns using
// `OperationMoreStringFuncMap` should return a map where the keys are column names and the values are
// functions that return a string to populate the cell.
func (sm *SpecMore) WriteFileXLSX(filename string, columns *tabulator.ColumnSet, filterFunc func(path, method string, op *oas3.Operation) bool, addlColFuncs *OperationMoreStringFuncMap) error {
	if columns == nil {
		columns = OpTableColumnsDefault(true)
	}
	if tbl, err := sm.OperationsTable(columns, filterFunc, addlColFuncs); err != nil {
		return err
	} else {
		tbl.FormatAutoLink = true
		return table.WriteXLSX(filename, []*table.Table{tbl})
	}
}

func (sm *SpecMore) WriteFileYAML(filename string, perm os.FileMode) error {
	if ybytes, err := sm.MarshalYAML(); err != nil {
		return err
	} else {
		return os.WriteFile(filename, ybytes, perm)
	}
}

type TagsMore struct {
	Tags oas3.Tags
}

func (tg *TagsMore) Get(tagName string) *oas3.Tag {
	for _, tag := range tg.Tags {
		if tagName == tag.Name {
			return tag
		}
	}
	return nil
}
