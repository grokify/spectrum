# OpenAPI3Lint

`openapi3lint` is a linting tool for OpenAPI 3 specifications using Go.

Features include:

* Standard and custom rules
* Violation locations by JSON pointer, including file/URL roots
* Violation locations grouped by rule

It is designed to scale to a large set of APIs split across many files.

## Standard Rules

### Command Line Application

Standard rules can be executed through the `cmd/oas3lint` CLI program. It takes two parameters:

* `-i` for the OAS3 specifictio file or diectory. If a directory, it will ead in all JSON/YAML/YML extension files.
* `-p` for the linter Policy config file.
* `-s` is optional aand used to select the severity level used. If none is selected, `error` is used.

### Policy File Format

The Policy file uses the following syntax:

```json
{
    "rules":{
        "datatype-int-format-int32-int64": {
            "severity": "off"
        },
        "operation-operationid-style-camelcase":{
            "severity": "error"
        }
    }
}
```

### Severity Levels

`openapi3lint` uses Syslog-like severity levels defined in `github.com/grokify/simplego/log/severity`, including:

```go
const (
	SeverityDisabled      = "disabled"
	SeverityEmergency     = "emerg"
	SeverityAlert         = "alert"
	SeverityCritical      = "crit"
	SeverityError         = "err"
	SeverityWarning       = "warning"
	SeverityNotice        = "notice"
	SeverityInformational = "info"
	SeverityDebug         = "debug"
)
```

The following aliases can be used:

```go
var mapStringSeverity = map[string]string{
	"disabled":      SeverityDisabled,
	"disable":       SeverityDisabled,
	"off":           SeverityDisabled,
	"emergency":     SeverityEmergency,
	"emerg":         SeverityEmergency,
	"panic":         SeverityEmergency, // deprecated by syslog
	"exception":     SeverityEmergency, // used by PostgreSQL
	"alert":         SeverityAlert,
	"critical":      SeverityCritical,
	"crit":          SeverityCritical,
	"error":         SeverityError,
	"err":           SeverityError,
	"warning":       SeverityWarning,
	"warn":          SeverityWarning,
	"notice":        SeverityNotice,
	"informational": SeverityInformational,
	"info":          SeverityInformational,
	"debug":         SeverityDebug,
	"hint":          SeverityDebug, // used by Spectral
}
```

### Standard Rules List

The following standard rules are built into the `openapi3lint` More are coming soon and this is under active development, including refactoring.

* `datatype-int-format-int32-int64`: reports if `type: integer` doesn't have a standard `format` set to `int32` or `int64`
* `operation-operationid-style-camelcase`: reports if `operationId` is not camel case
* `operation-operationid-style-kebabcase`: reports if `operationId` is not kebab case
* `operation-operationid-style-pascalcase`: reports if `operationId` is not pascal case
* `operation-operationid-style-snakecase`: reports if `operationId` is not snake case

## Custom Rules

Custom rules are created using the `Rule` interface. After implementing aa custom rule, load it into a `Policy` to execute.

Use `Policy.AddRule(rule Rule, errorOnCollision bool)` to add a rule.

### Rule Interface

A rule has the following interface:

```go
type Rule interface {
	Name() string
	Scope() string
	Severity() string
	ProcessSpec(spec *oas3.Swagger, pointerBase string) *lintutil.PolicyViolationsSets
	ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string) []lintutil.PolicyViolation
}
```

Functions:

* `Name()` should return the name of a rule in kebab case.
* `Scope()` should return the type of object / property operated on. This affects the processing function provided. As of now, `operration` and `specfication` are supported.
* `Severity()` should return a syslog like severity level supported by `github.com/grokify/simplego/log/severity`. This should be updated for the `Policy` used.
* `ProcessSpec(spec *oas3.Swagger, pointerBase string)` is a function to process a rule at the top specfication level. `pointerBase` is used to provide JSON Pointer info before the `#`. This is executed when `Scope()` is set to `specification`.
* `ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string)` is executed when `Scope()` is set to `operation`.