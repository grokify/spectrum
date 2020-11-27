package openapi3

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/os/osutil"
	"github.com/grokify/gotilla/type/maputil"
)

// OperationPropertiesDescriptionStatus returns a set of
// operationIds and parameters with description status where `1`
// indicates a description and `0` indicates no descriptions.
// Descriptions for references aren't processed so they aren't
// analyzed and reported on. This returns a `MapStringMapStringInt`
// where the first key is the operationIds and the second key is the
// parameter name.
func (sm *SpecMore) OperationPropertiesDescriptionStatus() maputil.MapStringMapStringInt {
	missingDescs := maputil.MapStringMapStringInt{}
	VisitOperations(sm.Spec, func(path, method string, op *oas3.Operation) {
		if op == nil {
			return
		}
		for _, paramRef := range op.Parameters {
			if paramRef == nil {
				continue
			}
			// Is a reference
			if len(strings.TrimSpace(paramRef.Ref)) > 0 {
				continue
			}
			// Is not a reference but has no value.
			if paramRef.Value == nil {
				continue
			}
			descTry := strings.TrimSpace(paramRef.Value.Description)
			if len(descTry) == 0 {
				missingDescs.Set(op.OperationID, paramRef.Value.Name, 0)
			} else {
				missingDescs.Set(op.OperationID, paramRef.Value.Name, 1)
			}
		}
	})
	return missingDescs
}

// SchemaPropertiesDescriptionStatus returns a set of
// schema names and properties with description status where `1`
// indicates a description and `0` indicates no descriptions.
// Descriptions for references aren't processed so they aren't
// analyzed and reported on. This returns a `MapStringMapStringInt`
// where the first key is the component name and the second key is the
// property name.
func (sm *SpecMore) SchemaPropertiesDescriptionStatus() maputil.MapStringMapStringInt {
	missingDescs := maputil.MapStringMapStringInt{}
	for schName, schRef := range sm.Spec.Components.Schemas {
		if len(schRef.Ref) > 0 || schRef.Value == nil {
			continue
		}
		for propName, propRef := range schRef.Value.Properties {
			if propRef == nil ||
				len(propRef.Ref) > 0 ||
				propRef.Value == nil {
				continue
			}
			desc := strings.TrimSpace(propRef.Value.Description)
			if len(desc) == 0 {
				missingDescs.Set(schName, propName, 0)
			} else {
				missingDescs.Set(schName, propName, 1)
			}
		}
	}
	return missingDescs
}

func (sm *SpecMore) OperationParametersWithoutDescriptionsWriteFile(filename string) error {
	missing := sm.OperationPropertiesDescriptionStatus()
	arr := missing.Flatten("#/paths/...", "/",
		maputil.MapStringMapStringIntFuncExactMatch(1),
		true, true)
	withCount1, withCount2 := missing.CountsWithVal(1, " ~~~ ")
	woutCount1, woutCount2 := missing.CountsWithVal(0, " ~~~ ")
	allCount1, allCount2 := missing.Counts(" ~~~ ")
	lines := []string{
		fmt.Sprintf("Operations Missing/Have/All [%d/%d/%d] Params Missing/Have/All [%d/%d/%d]",
			woutCount1, withCount1, allCount1,
			woutCount2, withCount2, allCount2),
	}
	lines = append(lines, arr...)

	return osutil.CreateFileWithLines(filename, lines, "\n", true)
}

func (sm *SpecMore) SchemaPropertiesWithoutDescriptionsWriteFile(filename string) error {
	missing := sm.SchemaPropertiesDescriptionStatus()
	arr := missing.Flatten("#/components/schemas", "/",
		maputil.MapStringMapStringIntFuncExactMatch(1),
		true, true)
	withCount1, withCount2 := missing.CountsWithVal(1, " ~~~ ")
	woutCount1, woutCount2 := missing.CountsWithVal(0, " ~~~ ")
	allCount1, allCount2 := missing.Counts(" ~~~ ")
	lines := []string{
		fmt.Sprintf("Schemas Missing/Have/All [%d/%d/%d] Props Missing/Have/All [%d/%d/%d]",
			woutCount1, withCount1, allCount1,
			woutCount2, withCount2, allCount2),
	}
	lines = append(lines, arr...)

	return osutil.CreateFileWithLines(filename, lines, "\n", true)
}
