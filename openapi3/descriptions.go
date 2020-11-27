package openapi3

import (
	"fmt"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/os/osutil"
	"github.com/grokify/gotilla/type/maputil"
)

// OperationPropertiesWithoutDescriptions returns a set of
// operation ids and parameters without descriptions. Descriptions
// for references aren't processed so they aren't analyzed and
// reported on. This returns a `MapStringMapStringInt` where the
// first key is the operation id and the second key is the
// parameter name.
func (sm *SpecMore) OperationPropertiesWithoutDescriptions() maputil.MapStringMapStringInt {
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
				missingDescs.Set(op.OperationID, paramRef.Value.Name, 1)
			}
		}
	})
	return missingDescs
}

// SchemaPropertiesWithoutDescriptions returns a set of
// schema names and properties without descriptions. Descriptions
// for references aren't processed so they aren't analyzed and
// reported on. This returns a `MapStringMapStringInt` where the
// first key is the component name and the second key is the
// property name.
func (sm *SpecMore) SchemaPropertiesWithoutDescriptions() maputil.MapStringMapStringInt {
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
				missingDescs.Set(schName, propName, 1)
			}
		}
	}
	return missingDescs
}

func (sm *SpecMore) OperationParametersWithoutDescriptionsWriteFile(filename string) error {
	missing := sm.OperationPropertiesWithoutDescriptions()
	arr := missing.Flatten("#/paths/...", "/", true, true)
	lines := []string{
		fmt.Sprintf("Missing Desc Operations [%d] Missing Desc Parameters [%d]", len(missing), len(arr)),
	}
	lines = append(lines, arr...)

	return osutil.CreateFileWithLines(filename, lines, "\n", true)
}

func (sm *SpecMore) SchemaPropertiesWithoutDescriptionsWriteFile(filename string) error {
	missing := sm.SchemaPropertiesWithoutDescriptions()
	arr := missing.Flatten("#/components/schemas", "/", true, true)
	lines := []string{
		fmt.Sprintf("Missing Desc Schemas [%d] Missing Desc Properties [%d]", len(missing), len(arr)),
	}
	lines = append(lines, arr...)

	return osutil.CreateFileWithLines(filename, lines, "\n", true)
}
