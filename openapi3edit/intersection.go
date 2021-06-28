package openapi3edit

import (
	"sort"

	"github.com/grokify/simplego/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
)

type SpecMetadata struct {
	Endpoints    []string
	OperationIDs []string
	SchemaNames  []string
}

func NewSpecMetadata(spec *openapi3.Spec) SpecMetadata {
	md := SpecMetadata{
		OperationIDs: []string{},
		Endpoints:    []string{}}
	if spec != nil {
		mapOpIDs := SpecOperationIds(spec)
		for key := range mapOpIDs {
			md.OperationIDs = append(md.OperationIDs, key)
		}
		md.Endpoints = SpecEndpoints(spec, true)
		sm := openapi3.SpecMore{Spec: spec}
		md.SchemaNames = sm.SchemaNames()
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

func SpecsIntersection(spec1, spec2 *openapi3.Spec) IntersectionData {
	idata := IntersectionData{
		Spec1: NewSpecMetadata(spec1),
		Spec2: NewSpecMetadata(spec2)}
	idata.Intersection = idata.Spec1.Intersection(idata.Spec2)
	return idata
}
