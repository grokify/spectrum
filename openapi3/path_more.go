package openapi3

import (
	"errors"
	"fmt"
	"reflect"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/net/http/httputilmore"
)

type PathItemMore struct {
	PathItem *oas3.PathItem
}

func (pm *PathItemMore) AddPathItemOperations(add *oas3.PathItem, overwriteOpration bool) error {
	if add == nil {
		return nil
	} else if pm.PathItem == nil {
		return errors.New("path item is not set")
	}
	methods := httputilmore.Methods()
	for _, method := range methods {
		opAdd := add.GetOperation(method)
		if opAdd == nil {
			continue
		} else if overwriteOpration {
			pm.PathItem.SetOperation(method, opAdd)
		} else {
			opSrc := pm.PathItem.GetOperation(method)
			if opSrc == nil {
				pm.PathItem.SetOperation(method, opAdd)
			} else if !reflect.DeepEqual(opAdd, opSrc) {
				return fmt.Errorf("operation collision on op id (%s)", opSrc.OperationID)
			}
		}
	}
	return nil
}
