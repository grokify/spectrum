# OpenAPI3Lint - Custom Rules

Ther eare two parttst of using custom rules:

1. Custom Rule
1. Custom Rule Collection

Simple rules can be used directly with the simple Rule Collection, `RuleCollectionSimple`, however, more complex rules can be built using custom Rule Collections. Complex rules allow multiple rule names to be used with a single rule definition, with the Rule Collection handling the instantiation.

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

#### Functions

* `Name()` should return the name of a rule in kebab case.
* `Scope()` should return the type of object / property operated on. This affects the processing function provided. As of now, `operation` and `specfication` are supported.
* `Severity()` should return a syslog like severity level supported by `github.com/grokify/simplego/log/severity`. This should be updated for the `Policy` used.
* `ProcessSpec(spec *oas3.Swagger, pointerBase string)` is a function to process a rule at the top specfication level. `pointerBase` is used to provide JSON Pointer info before the `#`. This is executed when `Scope()` is set to `specification`.
* `ProcessOperation(spec *oas3.Swagger, op *oas3.Operation, opPointer, path, method string)` is executed when `Scope()` is set to `operation`.

## Rule Collection

```go
type RuleCollection interface {
	Name() string
	RuleNames() []string
	RuleExists(ruleName string) bool
	Rule(ruleName string) (Rule, error)
}
```
