package openapi3edit

import (
	"sort"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/type/stringsutil"
)

type SpecMetadata struct {
	Endpoints    []string
	OperationIDs []string
	SchemaNames  []string
}

func NewSpecMetadata(spec *oas3.Swagger) SpecMetadata {
	md := SpecMetadata{
		OperationIDs: []string{},
		Endpoints:    []string{}}
	if spec != nil {
		mapOpIDs := SpecOperationIds(spec)
		for key := range mapOpIDs {
			md.OperationIDs = append(md.OperationIDs, key)
		}
		md.Endpoints = SpecEndpoints(spec, true)
		md.SchemaNames = SpecSchemaNames(spec)
	}
	return md
}

func (md *SpecMetadata) Intersection(md2 SpecMetadata) SpecMetadata {
	idata := SpecMetadata{
		Endpoints:    stringsutil.SliceIntersection(md.Endpoints, md2.Endpoints),
		OperationIDs: stringsutil.SliceIntersection(md.OperationIDs, md2.OperationIDs),
		SchemaNames:  stringsutil.SliceIntersection(md.SchemaNames, md2.SchemaNames)}
	return idata
}

func (md *SpecMetadata) IsEmpty() bool {
	if len(md.Endpoints) == 0 &&
		len(md.OperationIDs) == 0 &&
		len(md.SchemaNames) == 0 {
		return true
	}
	return false
}

func (md *SpecMetadata) Sort() {
	sort.Strings(md.Endpoints)
	sort.Strings(md.OperationIDs)
	sort.Strings(md.SchemaNames)
}

type IntersectionData struct {
	Spec1        SpecMetadata
	Spec2        SpecMetadata
	Intersection SpecMetadata
}

func (idata *IntersectionData) Sort() {
	idata.Spec1.Sort()
	idata.Spec2.Sort()
	idata.Intersection.Sort()
}

func NewIntersectionData() IntersectionData {
	return IntersectionData{
		Spec1:        NewSpecMetadata(nil),
		Spec2:        NewSpecMetadata(nil),
		Intersection: NewSpecMetadata(nil)}
}

func SpecsIntersection(spec1, spec2 *oas3.Swagger) IntersectionData {
	idata := IntersectionData{
		Spec1: NewSpecMetadata(spec1),
		Spec2: NewSpecMetadata(spec2)}
	idata.Intersection = idata.Spec1.Intersection(idata.Spec2)
	return idata
}
