package modify

import (
	"fmt"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/swaggman/openapi3"
	"github.com/pkg/errors"
)

type SpecMoreModifyMultiOpts struct {
	OperationIdsRename           func(string) string
	OperationIdsRenameExec       bool
	OperationIdsShow             bool
	Paths                        SpecPathsModifyOpts
	PathsShow                    bool
	PathsExec                    bool
	Tags                         map[string]string
	TagsShow                     bool
	TagsExec                     bool
	OperationsRemoveSecurityExec bool
}

// SpecMoreModifyMulti is used to perform multiple updates on
// an OpenAPI 3 spec.
func SpecMoreModifyMulti(sm *openapi3.SpecMore, opts SpecMoreModifyMultiOpts) error {
	if opts.OperationIdsShow {
		fmtutil.PrintJSON(SpecOperationIds(sm.Spec))
		oldIds := SpecOperationIds(sm.Spec)
		if opts.OperationIdsShow {
			fmtutil.PrintJSON(oldIds)
		}
		for id, count := range oldIds {
			if count != 1 {
				return fmt.Errorf("E_OPERATION_ID_BAD_COUNT ID[%s]COUNT[%d]", id, count)
			}
		}
	}
	if opts.OperationIdsRenameExec {
		UpdateOperationIds(sm.Spec, opts.OperationIdsRename)
		newIds := SpecOperationIds(sm.Spec)
		if opts.OperationIdsShow {
			fmtutil.PrintJSON(newIds)
		}
		for id, count := range newIds {
			if count != 1 {
				return fmt.Errorf("E_OPERATION_ID_BAD_COUNT_AFTER_RENAME ID[%s]COUNT[%d]", id, count)
			}
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
	if len(opts.Tags) > 0 {
		if opts.TagsShow {
			fmtutil.PrintJSON(SpecTags(sm.Spec))
		}
		if opts.TagsExec {
			SpecTagsModify(sm.Spec, opts.Tags)
			if opts.TagsShow {
				fmtutil.PrintJSON(SpecTags(sm.Spec))
			}
		}
	}

	if opts.OperationsRemoveSecurityExec {
		RemoveOperationsSecurity(sm.Spec)
	}

	return nil
}
