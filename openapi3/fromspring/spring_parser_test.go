package fromspring

import (
	"testing"
)

var parseLineTests = []struct {
	v          string
	oasName    string
	oasType    string
	oasFormat  string
	oasDefault interface{}
}{
	{"private Boolean myPropBoolean;", "myPropBoolean", "boolean", "", nil},
	{"private DateTime myPropDateTime;", "myPropDateTime", "string", "date-time", nil},
	{"private Integer myPropInteger = 1;", "myPropInteger", "integer", "int64", 1},
	{"private Long myPropLong = 1;", "myPropLong", "integer", "int64", 1},
	{"private String myPropString;", "myPropString", "string", "", nil},
	{"private String myPropString = \"\";", "myPropString", "string", "", ""},
	{"private String myPropString = \"AVAILABLE\";", "myPropString", "string", "", "AVAILABLE"},
}

func TestParseLine(t *testing.T) {
	for _, tt := range parseLineTests {
		name, schema, err := ParseSpringLineToSchema(tt.v)
		if err != nil {
			t.Errorf("fromspring.ParseSpringLineToSchema() [%v]", err)
		}
		if tt.oasName != name || tt.oasType != schema.Type || tt.oasFormat != schema.Format {
			t.Errorf(`fromspring.ParseSpringLineToSchema("%s") MISMATCH W[%v]G[%v] [%v][%v] [%v][%v]`, tt.v, tt.oasName, name, tt.oasType, schema.Type, tt.oasFormat, schema.Format)
		}
		//fmtutil.PrintJSON(schema)
	}
}
