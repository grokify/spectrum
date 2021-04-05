package stylechecker

/*

 https://stackoverflow.com/questions/1128305/regex-for-pascalcased-words-aka-camelcased-with-leading-uppercase-letter
 https://gist.github.com/manjeettahkur/ff114ef92d8ffee1b797091ff77ea89f

*/
/*
const (
	CaseCamel  = "camelCase"
	CaseKebab  = "kebab-case"
	CasePascal = "PascalCase"
	CaseSnake  = "snake_case"
)

type CaseValidator struct{}

func ValidateCase(caseType, s string) (bool, error) {
	switch caseType {
	case CaseCamel:
		{
			return IsCamelCase(s), nil
		}
	case CasePascal:
		{
			return IsPascalCase(s), nil
		}
	}
	return false, fmt.Errorf("unkown string case type [%s]", caseType)
}

var (
	rxCamelCase  = regexp.MustCompile(`^[a-z][0-9A-Za-z]*$`)
	rxPascalCase = regexp.MustCompile(`^[A-Z][0-9A-Za-z]*`)
	rxIdSuffix   = regexp.MustCompile(`(I[dD])$`)
)

func IsCamelCase(input string) bool {
	if !rxCamelCase.MatchString(input) {
		return false
	}
	m := rxIdSuffix.FindStringSubmatch(input)
	if len(m) == 2 {
		if m[1] != "Id" {
			return false
		}
	}
	return true
}

func IsPascalCase(input string) bool {
	return rxPascalCase.MatchString(input)
}

var rxFirstAlphaUpper = regexp.MustCompile(`^[A-Z]`)

func IsFirstAlphaUpper(s string) bool {
	if rxFirstAlphaUpper.MatchString(s) {
		return true
	}
	return false
}
*/
