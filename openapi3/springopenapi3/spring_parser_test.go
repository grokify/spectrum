package springopenapi3

import (
	"testing"

	"github.com/grokify/spectrum/openapi3"
)

var parseLineTests = []struct {
	v                   string
	oasName             string
	oasType             string
	oasFormat           string
	oasDefault          interface{}
	explicitCustomTypes []string
}{
	{"private Boolean myPropBoolean;", "myPropBoolean", "boolean", "", nil, []string{}},
	{"private DateTime myPropDateTime;", "myPropDateTime", "string", "", nil, []string{}},
	{"private Integer myPropInteger = 1;", "myPropInteger", "integer", "", 1, []string{}},
	{"private Long myPropLong = 1;", "myPropLong", "integer", "int64", 1, []string{}},
	{"private String myPropString;", "myPropString", "string", "", nil, []string{}},
	{"private String myPropString = \"\";", "myPropString", "string", "", "", []string{}},
	{"private String myPropString = \"AVAILABLE\";", "myPropString", "string", "", "AVAILABLE", []string{}},
	{"private List<Integer> myPropStrings = new ArrayList<>();", "myPropStrings", "array", "", nil, []string{}},
	{"private List<String> myPropStrings = new ArrayList<>();", "myPropStrings", "array", "", nil, []string{}},
	{"private List<Integer> myPropStrings;", "myPropStrings", "array", "", nil, []string{}},
	{"private List<String> myPropStrings;", "myPropStrings", "array", "", nil, []string{}},
}

func TestParseLine(t *testing.T) {
	for _, tt := range parseLineTests {
		name, schemaRef, err := ParseSpringLineToSchemaRef(tt.v, tt.explicitCustomTypes)
		if err != nil {
			t.Errorf("fromspring.ParseSpringLineToSchema() [%v]", err)
		}
		schema := schemaRef.Value
		if tt.oasName != name || tt.oasType != openapi3.TypesRefString(schema.Type) || tt.oasFormat != schema.Format {
			t.Errorf(`fromspring.ParseSpringLineToSchema("%s") MISMATCH W[%v]G[%v] [%v][%v] [%v][%v]`, tt.v, tt.oasName, name, tt.oasType, schema.Type, tt.oasFormat, schema.Format)
		}
		//fmtutil.PrintJSON(schema)
	}
}

/*
const CampaignLeadSearchCriteriaSimple = `private List<Integer> leadIds = new ArrayList<>();
	private List<Integer> listIds = new ArrayList<>();
	private List<String> externIds = new ArrayList<>();
	private List<String> physicalStates;
	private List<String> agentDispositions;
	private List<String> leadPhoneNumbers = new ArrayList<>();
	private boolean orphanedLeadsOnly;
	private String callerId;
	private String leadPhoneNum;
	private List<Integer> campaignIds = new ArrayList<>();
	private String firstName;
	private String lastName;
	private String address1;
	private String address2;
	private String city;
	private String zip;
	private String emailAddress;
	private String auxData1;
	private String auxData2;
	private String auxData3;
	private String auxData4;
	private String auxData5;
	private Integer pendingAgentId;
	private Integer agentId;`
*/
