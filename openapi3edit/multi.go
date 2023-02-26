package openapi3edit

import (
	"fmt"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/spectrum/openapi3"
)

type SpecMoreModifyMultiOpts struct {
	OperationsDeleteFunc     func(opPath, opMethod string, op *oas3.Operation) bool
	OperationsRenameIDsFunc  func(string, string, *oas3.Operation)
	OperationsRemoveSecurity bool
	OperationsShowIDs        bool
	OperationsExec           bool
	Paths                    SpecPathsModifyOpts
	PathsShow                bool
	PathsExec                bool
	TagsOperationFunc        func(string, string, *oas3.Operation)
	Tags                     map[string]string
	TagsShow                 bool
	TagsExec                 bool
}

// SpecMoreModifyMulti is used to perform multiple updates on an OpenAPI 3 spec.
func SpecMoreModifyMulti(sm *openapi3.SpecMore, opts SpecMoreModifyMultiOpts) error {
	se := SpecEdit{SpecMore: *sm}
	if opts.OperationsShowIDs {
		// fmtutil.PrintJSON(SpecOperationIds(sm.Spec))
		oldIDs := sm.OperationIDsCounts()
		if opts.OperationsShowIDs {
			err := fmtutil.PrintJSON(oldIDs)
			if err != nil {
				return err
			}
		}
		for id, count := range oldIDs {
			if count != 1 {
				return fmt.Errorf("E_OPERATION_ID_BAD_COUNT ID[%s]COUNT[%d]", id, count)
			}
		}
	}
	if opts.OperationsExec {
		if opts.OperationsDeleteFunc != nil {
			se.DeleteOperations(opts.OperationsDeleteFunc)
		}
		if opts.OperationsRenameIDsFunc != nil {
			openapi3.VisitOperations(sm.Spec, opts.OperationsRenameIDsFunc)
		}
		// UpdateOperationIds(sm.Spec, opts.OperationIdsRename)
		newIDs := sm.OperationIDsCounts()
		if opts.OperationsShowIDs {
			err := fmtutil.PrintJSON(newIDs)
			if err != nil {
				return err
			}
		}
		for id, count := range newIDs {
			if count != 1 {
				return fmt.Errorf("E_OPERATION_ID_BAD_COUNT_AFTER_RENAME ID[%s]COUNT[%d]", id, count)
			}
		}

		if opts.OperationsRemoveSecurity {
			err := se.OperationsSecurityRemove([]string{})
			if err != nil {
				return err
			}
		}
	}
	// Update Paths
	if opts.PathsShow {
		err := fmtutil.PrintJSON(InspectPaths(sm.Spec))
		if err != nil {
			return err
		}
	}
	if opts.PathsExec {
		err := se.PathsModify(opts.Paths)
		if err != nil {
			return errorsutil.Wrap(err, "specModifyMulti")
		}
		if opts.PathsShow {
			err := fmtutil.PrintJSON(InspectPaths(sm.Spec))
			if err != nil {
				return err
			}
		}
	}

	// Update Tags
	if opts.TagsOperationFunc != nil || len(opts.Tags) > 0 {
		if opts.TagsShow {
			err := fmtutil.PrintJSON(sm.TagsMap(&openapi3.TagsOpts{InclDefs: true, InclOps: true}))
			if err != nil {
				return err
			}
		}
		if opts.TagsExec {
			if opts.TagsOperationFunc != nil {
				openapi3.VisitOperations(sm.Spec, opts.TagsOperationFunc)
			}
			se := SpecEdit{SpecMore: *sm}
			se.TagsModify(opts.Tags)
			if opts.TagsShow {
				err := fmtutil.PrintJSON(sm.TagsMap(&openapi3.TagsOpts{InclDefs: true, InclOps: true}))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
