package openapi3lint

import "path/filepath"

type PolicyViolationsSetsByFile struct {
	Sets map[string]PolicyViolationsSets
}

func NewPolicyViolationsSetsByFile() PolicyViolationsSetsByFile {
	return PolicyViolationsSetsByFile{
		Sets: map[string]PolicyViolationsSets{}}
}

func (byFile *PolicyViolationsSetsByFile) LocationsByRule(filenameOnly, skipEmpty bool) map[string]ViolationLocationsByRuleSet {
	res := map[string]ViolationLocationsByRuleSet{}
	for filename, vset := range byFile.Sets {
		if filenameOnly {
			_, file := filepath.Split(filename)
			filename = file
		}
		locs := vset.LocationsByRule()
		if skipEmpty && locs.Count() == 0 {
			continue
		}
		res[filename] = locs
	}
	return res
}
