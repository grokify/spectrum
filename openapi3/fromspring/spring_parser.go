package fromspring

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/gotilla/type/stringsutil"
)

const (
	TypeBoolean          string = "boolean"
	TypeInteger          string = "integer"
	TypeString           string = "string"
	FormatStringDate     string = "date"
	FormatStringDateTime string = "date-time"
	FormatIntegerInt32   string = "int32"
	FormatIntegerInt64   string = "int64"
)

var (
	rxSpringLine             *regexp.Regexp = regexp.MustCompile(`^private\s+(\S+)\s+(\S+)\s*;\s*$`)
	rxSpringLineBoolDef      *regexp.Regexp = regexp.MustCompile(`^private\s+Boolean\s+(\S+)\s+=\s+(true|false);\s*$`)
	rxSpringLineIntDef       *regexp.Regexp = regexp.MustCompile(`^private\s+Integer\s+(\S+)\s+=\s+(\d+);\s*$`)
	rxSpringLineIntOrLongDef *regexp.Regexp = regexp.MustCompile(`^private\s+(Integer|Long)\s+(\S+)\s+=\s+(\d+);\s*$`)
	rxSpringLineStringDef    *regexp.Regexp = regexp.MustCompile(`^private\s+String\s+(\S+)\s+=\s+\"(.*)\"\s*;\s*$`)
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
		return name, &prop, nil
	}
	return "", nil, nil
}

func lineToBoolDef(line string) (string, oas3.Schema, error) {
	m1 := rxSpringLineBoolDef.FindAllStringSubmatch(line, -1)
	if len(m1) > 0 {
		propName := m1[0][1]
		boolDefaultVal := m1[0][2]
		sch := oas3.Schema{
			Type: TypeBoolean}
		if boolDefaultVal == "true" {
			sch.Default = true
		} else {
			sch.Default = false
		}
		return propName, sch, nil
	}
	return "", oas3.Schema{}, nil
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
			Type:    TypeInteger,
			Default: defaultVal}
		if intOrLong == "long" {
			sch.Format = FormatIntegerInt64
		}
		return propName, sch, nil
	}
	return "", oas3.Schema{}, nil
}

func lineToStringDef(line string) (string, oas3.Schema, error) {
	// error needs to return empty name.
	m1 := rxSpringLineStringDef.FindAllStringSubmatch(line, -1)
	if len(m1) > 0 {
		propName := m1[0][1]
		sch := oas3.Schema{
			Type:    TypeString,
			Default: strings.TrimSpace(m1[0][2])}
		return propName, sch, nil
	}
	return "", oas3.Schema{}, nil
}

// ParseSpringLineToSchema parses a Spring Java code line and
// attempts to extract a property name, type, format and default
// value.
func ParseSpringLineToSchema(line string) (string, oas3.Schema, error) {
	sch := oas3.Schema{}
	line = strings.Trim(line, " \t")

	name, sch, err := lineToStringDef(line)
	if err == nil && len(name) > 0 {
		return name, sch, nil
	}
	name, sch, err = lineToBoolDef(line)
	if err == nil && len(name) > 0 {
		return name, sch, nil
	}
	name, sch, err = lineToIntOrLongDef(line)
	if err != nil {
		return "", oas3.Schema{}, err
	} else if len(name) > 0 {
		return name, sch, nil
	}

	m := rxSpringLine.FindAllStringSubmatch(line, -1)
	if len(m) == 0 {
		return "", sch, errors.New("E_SPRING_TO_OAS_SCHEMA_NO_MATCH")
	} else if len(m) != 1 && len(m[1]) != 3 {
		return "", sch, fmt.Errorf("E_SPRING_TO_OAS_SCHEMA_NO_MATCH [%v]", m)
	}
	m2a := m[0]
	propName := m2a[2]
	javaTypeLc := strings.ToLower(strings.TrimSpace(m2a[1]))
	switch javaTypeLc {
	case "boolean":
		sch.Type = TypeBoolean
	case "date":
		sch.Type = TypeString
		sch.Format = FormatStringDate
	case "datetime":
		sch.Type = TypeString
		sch.Format = FormatStringDateTime
	case "integer":
		sch.Type = TypeInteger
	case "long":
		sch.Type = TypeInteger
		sch.Format = FormatIntegerInt64
	case "string":
		sch.Type = TypeString
	default:
		panic(javaTypeLc)
	}
	return propName, sch, nil
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
			curLines = []string{}
		}
	}
	columns = stringsutil.Slice2FilterLinesHaveIndex(columns, "@Column", 0)
	return columns
}
