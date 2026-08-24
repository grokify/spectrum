package openapi3

import (
	"fmt"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

// SpecMetadataDiff holds a set comparison of two specs' metadata across the
// three dimensions tracked by SpecMetadata: endpoints (generic `path method`),
// operation IDs, and schema names. It is the complement of IntersectionData and
// is useful for spotting version drift, verifying a migration lost no coverage,
// or measuring what a newer spec adds.
type SpecMetadataDiff struct {
	OnlyInSpec1 SpecMetadata // present in spec 1 but not spec 2
	OnlyInSpec2 SpecMetadata // present in spec 2 but not spec 1
	Both        SpecMetadata // present in both
}

func NewSpecMetadataDiff() SpecMetadataDiff {
	return SpecMetadataDiff{
		OnlyInSpec1: NewSpecMetadata(),
		OnlyInSpec2: NewSpecMetadata(),
		Both:        NewSpecMetadata()}
}

// Sort sorts every dimension of each side of the diff.
func (d *SpecMetadataDiff) Sort() {
	d.OnlyInSpec1.Sort()
	d.OnlyInSpec2.Sort()
	d.Both.Sort()
}

// IsEmpty reports whether the two specs were identical across all dimensions,
// i.e. neither side had anything the other lacked.
func (d *SpecMetadataDiff) IsEmpty() bool {
	return d.OnlyInSpec1.IsEmpty() && d.OnlyInSpec2.IsEmpty()
}

// String renders a human-readable report of the differences, listing the
// items unique to each spec per dimension. Items common to both are summarized
// by count only.
func (d SpecMetadataDiff) String() string {
	d.Sort()
	var sb strings.Builder
	writeDim := func(name string, only1, only2, both []string) {
		fmt.Fprintf(&sb, "%s: %d in both, %d only in spec 1, %d only in spec 2\n",
			name, len(both), len(only1), len(only2))
		for _, v := range only1 {
			fmt.Fprintf(&sb, "  - only in spec 1: %s\n", v)
		}
		for _, v := range only2 {
			fmt.Fprintf(&sb, "  + only in spec 2: %s\n", v)
		}
	}
	writeDim("Endpoints", d.OnlyInSpec1.Endpoints, d.OnlyInSpec2.Endpoints, d.Both.Endpoints)
	writeDim("OperationIDs", d.OnlyInSpec1.OperationIDs, d.OnlyInSpec2.OperationIDs, d.Both.OperationIDs)
	writeDim("SchemaNames", d.OnlyInSpec1.SchemaNames, d.OnlyInSpec2.SchemaNames, d.Both.SchemaNames)
	return sb.String()
}

// Difference returns the metadata present in md but not in md2, evaluated
// independently for each dimension.
func (md *SpecMetadata) Difference(md2 SpecMetadata) SpecMetadata {
	return SpecMetadata{
		Endpoints:    stringsutil.SliceSubtract(md.Endpoints, md2.Endpoints),
		OperationIDs: stringsutil.SliceSubtract(md.OperationIDs, md2.OperationIDs),
		SchemaNames:  stringsutil.SliceSubtract(md.SchemaNames, md2.SchemaNames)}
}

// Diff returns a three-way comparison (only in md, only in md2, both) of the
// receiver against md2.
func (md *SpecMetadata) Diff(md2 SpecMetadata) SpecMetadataDiff {
	d := SpecMetadataDiff{
		OnlyInSpec1: md.Difference(md2),
		OnlyInSpec2: md2.Difference(*md),
		Both:        md.Intersection(md2)}
	d.Sort()
	return d
}

// SpecsDiff compares two specs across endpoints, operation IDs, and schema
// names, returning the items unique to each and those shared.
func SpecsDiff(spec1, spec2 *Spec) SpecMetadataDiff {
	sm1 := SpecMore{Spec: spec1}
	sm2 := SpecMore{Spec: spec2}
	md1 := sm1.Metadata()
	md2 := sm2.Metadata()
	return md1.Diff(md2)
}

// ReadSpecsDiff reads two OpenAPI 3 specs from file (JSON or YAML) and diffs
// them. It is a convenience wrapper over ReadFile and SpecsDiff.
func ReadSpecsDiff(file1, file2 string, validate bool) (SpecMetadataDiff, error) {
	spec1, err := ReadFile(file1, validate)
	if err != nil {
		return NewSpecMetadataDiff(), fmt.Errorf("reading spec 1 (%s): %w", file1, err)
	}
	spec2, err := ReadFile(file2, validate)
	if err != nil {
		return NewSpecMetadataDiff(), fmt.Errorf("reading spec 2 (%s): %w", file2, err)
	}
	return SpecsDiff(spec1, spec2), nil
}
