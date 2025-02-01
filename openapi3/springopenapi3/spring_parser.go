package springopenapi3

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/spectrum/openapi3"
)

const (
	/*TypeArray            = openapi3.TypeArray
	TypeBoolean          = openapi3.TypeBoolean
	TypeInteger          = openapi3.TypeInteger
	TypeObject           = openapi3.TypeObject
	TypeString           = openapi3.TypeString*/
	FormatStringDate     = openapi3.FormatDate
	FormatStringDateTime = openapi3.FormatDateTime
	FormatIntegerInt64   = openapi3.FormatInt64
)

var (
	rxSpringLine             = regexp.MustCompile(`^(?:private\s+)?(\S+)\s+(\S+)\s*;\s*$`)
	rxSpringLineBoolDef      = regexp.MustCompile(`^private\s+[Bb]oolean\s+(\S+)\s+=\s+(true|false);\s*$`)
	rxSpringLineIntOrLongDef = regexp.MustCompile(`^private\s+(Integer|Long)\s+(\S+)\s+=\s+(\d+);\s*$`)
	rxSpringLineStringDef    = regexp.MustCompile(`^private\s+String\s+(\S+)\s+=\s+"(.*)"\s*;\s*$`)
	rxSpringLineIntArrayDef  = regexp.MustCompile(`^private\s+List<([^<>]+)>\s+(\S+)\b`)
)

// ParseSpringPropertyLinesSliceToSchema takes a set of string slices
// and attempts to parse one property per set of lines.
func ParseSpringPropertyLinesSliceToSchema(groups [][]string) (map[string]*oas3.SchemaRef, error) {
	mss := map[string]*oas3.SchemaRef{}
	for _, group := range groups {
		name, prop, err := ParseSpringPropertyLinesToSchema(group)
		if err != nil {
			return mss, err
		} else if name == "" || prop == nil {
			continue
		}
		if prop != nil {
			mss[name] = oas3.NewSchemaRef("", prop)
		} else {
			mss[name] = oas3.NewSchemaRef(name, nil)
		}
	}
	return mss, nil
}

// ParseSpringPropertyLinesToSchema parses a set of lines looking for
// a property line. Only one property line is matched in this set.
func ParseSpringPropertyLinesToSchema(lines []string) (string, *oas3.Schema, error) {
	for _, line := range lines {
		name, prop, err := ParseSpringLineToSchema(line)
		if err != nil { // not every line is designed to match
			continue
		}
		return name, prop, nil
	}
	return "", nil, nil
}

func lineToBoolDef(line string) (string, oas3.Schema) {
	m1 := rxSpringLineBoolDef.FindAllStringSubmatch(line, -1)
	if len(m1) > 0 {
		propName := m1[0][1]
		boolDefaultVal := m1[0][2]
		sch := oas3.Schema{
			Type: openapi3.NewTypesRef(openapi3.TypeBoolean)}
		if boolDefaultVal == "true" {
			sch.Default = true
		} else {
			sch.Default = false
		}
		return propName, sch
	}
	return "", oas3.Schema{}
}

func lineToIntOrLongDef(line string) (string, oas3.Schema, error) {
	// error needs to return empty name.
	m1 := rxSpringLineIntOrLongDef.FindAllStringSubmatch(line, -1)
	if len(m1) > 0 {
		intOrLong := strings.ToLower(strings.TrimSpace(m1[0][1]))
		propName := m1[0][2]
		intDefaultVal := m1[0][3]
		defaultVal, err := strconv.Atoi(intDefaultVal)
		if err != nil {
			return "", oas3.Schema{}, err
		}
		sch := oas3.Schema{
			Type:    openapi3.NewTypesRef(openapi3.TypeInteger),
			Default: defaultVal}
		if intOrLong == "long" {
			sch.Format = FormatIntegerInt64
		}
		return propName, sch, nil
	}
	return "", oas3.Schema{}, nil
}

func lineToStringDef(line string) (string, oas3.Schema) {
	// error needs to return empty name.
	m1 := rxSpringLineStringDef.FindAllStringSubmatch(line, -1)
	if len(m1) > 0 {
		propName := m1[0][1]
		sch := oas3.Schema{
			Type:    openapi3.NewTypesRef(openapi3.TypeString),
			Default: strings.TrimSpace(m1[0][2])}
		return propName, sch
	}
	return "", oas3.Schema{}
}

func lineToArrayDef(line string, explicitCustomTypes []string) (string, *oas3.SchemaRef) {
	// error needs to return empty name.
	// private List<Integer> leadIds = new ArrayList<>();
	m1 := rxSpringLineIntArrayDef.FindAllStringSubmatch(line, -1)
	//fmt.Println("lineToArrayDef")
	//fmtutil.PrintJSON(m1)
	if len(m1) > 0 {
		javaType := strings.TrimSpace(m1[0][1])
		javaTypeLc := strings.ToLower(javaType)
		propName := m1[0][2]
		switch javaTypeLc {
		case openapi3.TypeInteger:
			sch := oas3.Schema{
				Type: openapi3.NewTypesRef(openapi3.TypeArray),
				Items: oas3.NewSchemaRef("",
					&oas3.Schema{
						Type: openapi3.NewTypesRef(openapi3.TypeInteger)})}
			sr := oas3.NewSchemaRef("", &sch)
			return propName, sr
		case openapi3.TypeString:
			sch := oas3.Schema{
				Type: openapi3.NewTypesRef(openapi3.TypeArray),
				Items: oas3.NewSchemaRef("",
					&oas3.Schema{
						Type: openapi3.NewTypesRef(openapi3.TypeString)})}
			sr := oas3.NewSchemaRef("", &sch)
			return propName, sr
		default:
			for _, exType := range explicitCustomTypes {
				exType = strings.TrimSpace(exType)
				if exType == javaType {
					sr := oas3.NewSchemaRef(schemaPath(exType), nil)
					return propName, sr
				}
			}
		}
	}
	return "", nil
}

// ParseSpringLinesToMapStringSchemaRefs parses a Spring Java code line and
// attempts to extract a property name, type, format and default value.
func ParseSpringLinesToMapStringSchemaRefs(lines, explicitCustomTypes []string) (map[string]*oas3.SchemaRef, error) {
	mss := map[string]*oas3.SchemaRef{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		name, prop, err := ParseSpringLineToSchemaRef(line, explicitCustomTypes)
		if err != nil {
			return mss, err
		} else if name == "" || prop == nil {
			continue
		} else {
			mss[name] = prop
		}
	}
	return mss, nil
}

// ParseSpringLineToSchemaRef parses a Spring Java code line and
// attempts to extract a property name, type, format and default
// value.
func ParseSpringLineToSchemaRef(line string, explicitCustomTypes []string) (string, *oas3.SchemaRef, error) {
	var sch oas3.Schema
	line = strings.Trim(line, " \t")

	name, sch := lineToStringDef(line)
	if len(name) > 0 {
		return name, oas3.NewSchemaRef("", &sch), nil
	}
	name, sch = lineToBoolDef(line)
	if len(name) > 0 {
		return name, oas3.NewSchemaRef("", &sch), nil
	}
	name, schRef := lineToArrayDef(line, explicitCustomTypes)
	if len(name) > 0 {
		return name, schRef, nil
	}
	name, sch, err := lineToIntOrLongDef(line)
	if err != nil {
		return "", nil, err
	} else if len(name) > 0 {
		return name, oas3.NewSchemaRef("", &sch), nil
	}

	m := rxSpringLine.FindAllStringSubmatch(line, -1)
	if len(m) == 0 {
		return "", nil, fmt.Errorf("spring-to-oas schema-no-match [%v]", line)
	} else if len(m) != 1 && len(m[1]) != 3 {
		return "", nil, fmt.Errorf("spring-to-oas schema-no-match [%v]", m)
	}
	m2a := m[0]
	propName := m2a[2]
	javaType := strings.TrimSpace(m2a[1])
	javaTypeLc := strings.ToLower(javaType)
	schemaRef := &oas3.SchemaRef{}
	switch javaTypeLc {
	case openapi3.TypeBoolean:
		schemaRef = oas3.NewSchemaRef("", &oas3.Schema{Type: openapi3.NewTypesRef(openapi3.TypeBoolean)})
	case "date":
		schemaRef = oas3.NewSchemaRef("", &oas3.Schema{
			Type: openapi3.NewTypesRef(openapi3.TypeString), Format: FormatStringDate})
	case "datetime":
		schemaRef = oas3.NewSchemaRef("", &oas3.Schema{
			Type:        openapi3.NewTypesRef(openapi3.TypeString),
			Description: "Date-time in Java format. Example: `2019-01-01T01:01:01.000+0000`. Note this is not compatible with RFC-3339 which is used by OpenAPI 3.0 Spec because it doesn't have a `:` between hours and minutes.",
		})
	case openapi3.TypeInteger:
		schemaRef = oas3.NewSchemaRef("", &oas3.Schema{Type: openapi3.NewTypesRef(openapi3.TypeInteger)})
	case "long":
		schemaRef = oas3.NewSchemaRef("", &oas3.Schema{
			Type: openapi3.NewTypesRef(openapi3.TypeInteger), Format: FormatIntegerInt64})
	case openapi3.TypeString:
		schemaRef = oas3.NewSchemaRef("", &oas3.Schema{Type: openapi3.NewTypesRef(openapi3.TypeString)})
	default:
		found := false
		for _, exType := range explicitCustomTypes {
			if javaType == exType {
				schemaRef = oas3.NewSchemaRef(schemaPath(javaType), nil)
				found = true
			}
		}
		if !found {
			panic(fmt.Sprintf("TYPE [%v] LINE [%v]", javaTypeLc, line))
		}
	}
	return propName, schemaRef, nil
}

func schemaPath(object string) string {
	return "#/components/schemas/" + object
}

// ParseSpringLineToSchema parses a Spring Java code line and
// attempts to extract a property name, type, format and default
// value. DEPRECATED
func ParseSpringLineToSchema(line string) (string, *oas3.Schema, error) {
	var sch oas3.Schema
	line = strings.Trim(line, " \t")

	name, sch := lineToStringDef(line)
	if len(name) > 0 {
		return name, &sch, nil
	}
	name, sch = lineToBoolDef(line)
	if len(name) > 0 {
		return name, &sch, nil
	} /*
		name, sch, err = lineToArrayDef(line)
		if err == nil && len(name) > 0 {
			return name, &sch, nil
		}*/
	name, sch, err := lineToIntOrLongDef(line)
	if err != nil {
		return "", nil, err
	} else if len(name) > 0 {
		return name, &sch, nil
	}

	m := rxSpringLine.FindAllStringSubmatch(line, -1)
	if len(m) == 0 {
		return "", &sch, fmt.Errorf("E_SPRING_TO_OAS_SCHEMA_NO_MATCH [%v]", line)
	} else if len(m) != 1 && len(m[1]) != 3 {
		return "", &sch, fmt.Errorf("E_SPRING_TO_OAS_SCHEMA_NO_MATCH [%v]", m)
	}
	m2a := m[0]
	propName := m2a[2]
	javaTypeLc := strings.ToLower(strings.TrimSpace(m2a[1]))
	switch javaTypeLc {
	case openapi3.TypeBoolean:
		sch.Type = openapi3.NewTypesRef(openapi3.TypeBoolean)
	case openapi3.FormatDate:
		sch.Type = openapi3.NewTypesRef(openapi3.TypeString)
		sch.Format = FormatStringDate
	case "datetime":
		sch.Type = openapi3.NewTypesRef(openapi3.TypeString)
		sch.Format = FormatStringDateTime
	case openapi3.TypeInteger:
		sch.Type = openapi3.NewTypesRef(openapi3.TypeInteger)
	case "long":
		sch.Type = openapi3.NewTypesRef(openapi3.TypeInteger)
		sch.Format = FormatIntegerInt64
	case openapi3.TypeString:
		sch.Type = openapi3.NewTypesRef(openapi3.TypeString)
	default:
		panic(fmt.Sprintf("TYPE [%v] LINE [%v]", javaTypeLc, line))
	}
	return propName, &sch, nil
}

// ParseSpringCodeColumnsRaw takes a set of Java code lines
// and groups them into lines per property. Not all Java
// code may be formatted in a way to take advantage of this.
func ParseSpringCodeColumnsRaw(input []string) [][]string {
	columns := [][]string{}
	curLines := []string{}
	for _, line := range input {
		if len(line) == 0 {
			if len(curLines) == 0 {
				continue
			} else if strings.Index(curLines[0], "@") == 0 {
				columns = append(columns, curLines)
				curLines = []string{}
			} else {
				curLines = []string{}
			}
			continue
		}
		curLines = append(curLines, line)
	}

	if len(curLines) > 0 {
		if strings.Index(curLines[0], "@") == 0 {
			columns = append(columns, curLines)
		}
	}
	columns = stringsutil.Slice2FilterLinesHaveIndex(columns, "@Column", 0)
	return columns
}
