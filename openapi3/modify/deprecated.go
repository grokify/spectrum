package modify

import (
	"fmt"
	"reflect"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/swaggman/openapi3"
)

/*
func CopyOperation(src oas3.Operation) oas3.Operation {
	dst := oas3.Operation{
		ExtensionProps: src.ExtensionProps,
		Tags:           src.Tags,
		Summary:        src.Summary,
		Description:    src.Description,
		OperationID:    src.OperationID,
		Parameters:     src.Parameters,
		RequestBody:    src.RequestBody,
		Responses:      src.Responses,
		Callbacks:      src.Callbacks,
		Deprecated:     src.Deprecated,
		Security:       src.Security,
		Servers:        src.Servers,
		ExternalDocs:   src.ExternalDocs,
	}
	return dst
}
*/

func SpecSchemasSetDeprecated(spec *oas3.Swagger, newDeprecated bool) {
	for _, schemaRef := range spec.Components.Schemas {
		if len(schemaRef.Ref) == 0 && schemaRef.Value != nil {
			schemaRef.Value.Deprecated = newDeprecated
		}
	}
}

func SpecOperationsSetDeprecated(spec *oas3.Swagger, newDeprecated bool) {
	openapi3.VisitOperations(
		spec,
		func(path, method string, op *oas3.Operation) {
			/*if op == nil {
				return
			}*/

			// op = &*op

			op.Summary = "SUMSUMSUMARY"
			op.Deprecated = true
			op.Tags = append(op.Tags, "NEW_TAG")

			newOp := oas3.Operation{
				ExtensionProps: op.ExtensionProps,
				Tags:           op.Tags,
				Summary:        op.Summary,
				Description:    op.Description,
				OperationID:    op.OperationID,
				Parameters:     op.Parameters,
				RequestBody:    op.RequestBody,
				Responses:      op.Responses,
				Callbacks:      op.Callbacks,
				Deprecated:     op.Deprecated,
				Security:       op.Security,
				Servers:        op.Servers,
				ExternalDocs:   op.ExternalDocs,
			}

			//newOp := oas3.Operation{}
			// reflect.Copy([]oas3.Operation{*op}, []oas3.Operation{newOp})
			SpecSetOperation(spec, path, method, newOp)
			if 1 == 0 {
				ps := reflect.ValueOf(op)
				// struct
				s := ps.Elem()
				if s.Kind() == reflect.Struct {
					// exported field
					f := s.FieldByName("Deprecated")
					if f.IsValid() {
						// A Value can be changed only if it is
						// addressable and was not obtained by
						// the use of unexported struct fields.
						if f.CanSet() {
							// change value of N
							if f.Kind() == reflect.Bool {
								f.SetBool(false)
								fmt.Println("SETTING_REEFLECT")
								//x := int64(7)
								/*if !f.OverflowInt(x) {
									f.Set(x)
								}*/
							}
						}
					}
				}
			}

			if 1 == 0 && op.OperationID == "listGlipGroups" {
				op.Deprecated = false
				fmt.Printf("DEP  [%v]\n", op.Deprecated)

				fmtutil.PrintJSON(op)
				panic("ABC_VISIT")
				fmt.Printf("DEP  [%v]\n", op.Deprecated)
				if 1 == 0 {
					panic("YYY")
					fmt.Printf("DEP  [%v]\n", op.Deprecated)
					op.Deprecated = false
					fmt.Printf("DEP  [%v]\n", op.Deprecated)
					fmtutil.PrintJSON(op)
					fmt.Printf("SSUM [%s]\n", op.Summary)
					op.Summary = op.Summary + " " + "ABC "
					fmt.Printf("SSUM [%s]\n", op.Summary)
				}
				//panic("ZZZ_LISTGLIPGROUPS")
			}

		},
	)
}

/*

func SpecOperationsRemoveDeprecated(spec *oas3.Swagger) {
	openapi3.VisitOperations(
		spec,
		func(path, method string, op *oas3.Operation) {
			if op == nil {
				return
			}
			op.Deprecated = false
			// SpecSetOperation(spec, path, method, &*op)
			// op = &*op
			if 1 == 1 && op.OperationID == "listGlipGroups" {
				op.Deprecated = false
				fmt.Printf("DEP  [%v]\n", op.Deprecated)

				fmtutil.PrintJSON(*op)
				fmt.Printf("DEP  [%v]\n", op.Deprecated)
				if 1 == 0 {
					panic("YYY")
					fmt.Printf("DEP  [%v]\n", op.Deprecated)
					op.Deprecated = true
					fmt.Printf("DEP  [%v]\n", op.Deprecated)
					fmtutil.PrintJSON(op)
					fmt.Printf("SSUM [%s]\n", op.Summary)
					op.Summary = op.Summary + " " + "ABC "
					fmt.Printf("SSUM [%s]\n", op.Summary)
				}
				panic("ZZZ_LISTGLIPGROUPS")
			}
		},
	)
}
*/
