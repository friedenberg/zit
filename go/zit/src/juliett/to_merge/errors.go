package to_merge

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
)

func MakeErrMergeConflict(sk *store_fs.FDPair) (err *ErrMergeConflict) {
	err = &ErrMergeConflict{}

	if sk != nil {
		err.ResetWith(sk)
	}

	return
}

type ErrMergeConflict struct {
	store_fs.FDPair
}

func (e *ErrMergeConflict) Is(target error) bool {
	_, ok := target.(*ErrMergeConflict)
	return ok
}

func (e *ErrMergeConflict) Error() string {
	return fmt.Sprintf("merge conflict for fds: %v", &e.FDPair)
}
