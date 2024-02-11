package to_merge

import (
	"fmt"

	"code.linenisgreat.com/zit/src/hotel/sku"
)

func MakeErrMergeConflict(sk *sku.ExternalFDs) (err *ErrMergeConflict) {
	err = &ErrMergeConflict{}

	if sk != nil {
		err.ResetWith(sk)
	}

	return
}

type ErrMergeConflict struct {
	sku.ExternalFDs
}

func (e *ErrMergeConflict) Is(target error) bool {
	_, ok := target.(*ErrMergeConflict)
	return ok
}

func (e *ErrMergeConflict) Error() string {
	return fmt.Sprintf("merge conflict for fds: %s", &e.ExternalFDs)
}
