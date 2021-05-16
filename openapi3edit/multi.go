package openapi3edit

import (
	"fmt"

	oas3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/spectrum/openapi3"
	"github.com/pkg/errors"
)

type SpecMoreModifyMultiOpts struct {
	OperationsDeleteFunc     func(urlpath, method string, op *oas3.Operation) bool
	OperationsRenameIdsFunc  func(string, string, *oas3.Operation)
	OperationsRemoveSecurity bool
	OperationsShowIds        bool
	OperationsExec           bool
	Paths                    SpecPathsModifyOpts
	PathsShow                bool
	PathsExec                bool
	TagsOperationFunc        func(string, string, *oas3.Operation)
	Tags                     map[string]string
	TagsShow                 bool
	TagsExec                 bool
}

// SpecMoreModifyMulti is used to perform multiple updates on
// an OpenAPI 3 spec.
func SpecMoreModifyMulti(sm *openapi3.SpecMore, opts SpecMoreModifyMultiOpts) error {
	if opts.OperationsShowIds {
		fmtutil.PrintJSON(SpecOperationIds(sm.Spec))
		oldIds := SpecOperationIds(sm.Spec)
		if opts.OperationsShowIds {
			fmtutil.PrintJSON(oldIds)
		}
		for id, count := range oldIds {
			if count != 1 {
				return fmt.Errorf("E_OPERATION_ID_BAD_COUNT ID[%s]COUNT[%d]", id, count)
			}
		}
	}
	if opts.OperationsExec {
		if opts.OperationsDeleteFunc != nil {
			SpecDeleteOperations(sm.Spec, opts.OperationsDeleteFunc)
		}
		if opts.OperationsRenameIdsFunc != nil {
			openapi3.VisitOperations(sm.Spec, opts.OperationsRenameIdsFunc)
		}
		// UpdateOperationIds(sm.Spec, opts.OperationIdsRename)
		newIds := SpecOperationIds(sm.Spec)
		if opts.OperationsShowIds {
			fmtutil.PrintJSON(newIds)
		}
		for id, count := range newIds {
			if count != 1 {
				return fmt.Errorf("E_OPERATION_ID_BAD_COUNT_AFTER_RENAME ID[%s]COUNT[%d]", id, count)
			}
		}

		if opts.OperationsRemoveSecurity {
			RemoveOperationsSecurity(sm.Spec)
		}
	}
	// Update Paths
	if opts.PathsShow {
		fmtutil.PrintJSON(InspectPaths(sm.Spec))
	}
	if opts.PathsExec {
		err := SpecPathsModify(sm.Spec, opts.Paths)
		if err != nil {
			return errors.Wrap(err, "SpecModifyMulti")
		}
		if opts.PathsShow {
			fmtutil.PrintJSON(InspectPaths(sm.Spec))
		}
	}

	// Update Tags
	if opts.TagsOperationFunc != nil || len(opts.Tags) > 0 {
		if opts.TagsShow {
			fmtutil.PrintJSON(sm.TagsMap(true, true))
		}
		if opts.TagsExec {
			if opts.TagsOperationFunc != nil {
				openapi3.VisitOperations(sm.Spec, opts.TagsOperationFunc)
			}
			SpecTagsModify(sm.Spec, opts.Tags)
			if opts.TagsShow {
				fmtutil.PrintJSON(sm.TagsMap(true, true))
			}
		}
	}

	return nil
}
