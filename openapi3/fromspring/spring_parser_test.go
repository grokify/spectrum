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
	{"private Boolean userManagedByRC;", "userManagedByRC", "boolean", "", nil},
	{"private DateTime creationDate;", "creationDate", "string", "date-time", nil},
	{"private Integer softphoneId = 1;", "softphoneId", "integer", "int64", 1},
	{"private String initLoginBaseState;", "initLoginBaseState", "string", "", nil},
	{"private String initLoginBaseState = \"AVAILABLE\";", "initLoginBaseState", "string", "", "AVAILABLE"},
	{"private String manualOutboundDefaultCallerId = \"\";", "manualOutboundDefaultCallerId", "string", "", nil},
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
