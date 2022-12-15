package openapi3

import (
	"sort"

	"github.com/grokify/mogo/type/stringsutil"
)

type SpecMetadata struct {
	Endpoints    []string
	OperationIDs []string
	SchemaNames  []string
}

func NewSpecMetadata() SpecMetadata {
	return SpecMetadata{
		Endpoints:    []string{},
		OperationIDs: []string{},
		SchemaNames:  []string{}}
}

func (sm *SpecMore) Metadata() SpecMetadata {
	md := NewSpecMetadata()
	if sm.Spec == nil {
		return md
	}
	mapOpIDs := sm.OperationIDsCounts()
	for key := range mapOpIDs {
		md.OperationIDs = append(md.OperationIDs, key)
	}
	md.Endpoints = sm.PathMethods(true)
	md.SchemaNames = sm.SchemaNames()
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
		Spec1:        NewSpecMetadata(),
		Spec2:        NewSpecMetadata(),
		Intersection: NewSpecMetadata()}
}

func SpecsIntersection(spec1, spec2 *Spec) IntersectionData {
	sm1 := SpecMore{Spec: spec1}
	sm2 := SpecMore{Spec: spec2}
	idata := IntersectionData{
		Spec1: sm1.Metadata(),
		Spec2: sm2.Metadata()}
	idata.Intersection = idata.Spec1.Intersection(idata.Spec2)
	return idata
}
