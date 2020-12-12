package openapi3

import (
	"reflect"
	"regexp"
)

type CollisionCheckResult int

const (
	CollisionCheckSame CollisionCheckResult = iota
	CollisionCheckOverwrite
	CollisionCheckError
	CollisionCheckSkip
)

type MergeOptions struct {
	FileRx               *regexp.Regexp
	SchemaFunc           func(schemaName string, sch1, sch2 interface{}, hint2 string) CollisionCheckResult
	CollisionCheckResult CollisionCheckResult
	ValidateEach         bool
	ValidateFinal        bool
}

func NewMergeOptionsSkip() *MergeOptions {
	return &MergeOptions{
		SchemaFunc: SchemaCheckCollisionSkip}
}

func (mo *MergeOptions) CheckSchemaCollision(schemaName string, sch1, sch2 interface{}, hint2 string) CollisionCheckResult {
	if mo.CollisionCheckResult == CollisionCheckSkip {
		mo.SchemaFunc = SchemaCheckCollisionSkip
	} else if mo.SchemaFunc == nil {
		mo.SchemaFunc = SchemaCheckCollisionDefault
	}
	return mo.SchemaFunc(schemaName, sch1, sch2, hint2)
}

func SchemaCheckCollisionDefault(schemaName string, item1, item2 interface{}, item2Note string) CollisionCheckResult {
	if reflect.DeepEqual(item1, item2) {
		return CollisionCheckSame
	}
	return CollisionCheckError
}

func SchemaCheckCollisionSkip(schemaName string, item1, item2 interface{}, item2Note string) CollisionCheckResult {
	if reflect.DeepEqual(item1, item2) {
		return CollisionCheckSame
	}
	return CollisionCheckSkip
}
