package openapi3lint

type RuleCollections []RuleCollection

type RuleCollection interface {
	Name() string
	RuleNames() []string
	RuleExists(ruleName string) bool
	Rule(ruleName string) (Rule, error)
}
