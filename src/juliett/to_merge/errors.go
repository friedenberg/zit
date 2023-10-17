package to_merge

import (
	"fmt"

	"github.com/friedenberg/zit/src/hotel/sku"
)

type ErrMergeConflict struct {
	sku.ExternalFDs
}

func (e ErrMergeConflict) Is(target error) bool {
	_, ok := target.(ErrMergeConflict)
	return ok
}

func (e ErrMergeConflict) Error() string {
	return fmt.Sprintf("merge conflict for fds: %s", e.ExternalFDs)
}
