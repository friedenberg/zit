package query

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
)

type ErrBlobMissing struct {
	ObjectId
	env_dir.ErrBlobMissing
}

// TODO add recovery text
func (e ErrBlobMissing) Error() string {
	return fmt.Sprintf(
		"Blob for %q with sha %q does not exist locally.",
		e.ObjectId,
		e.GetShaLike(),
	)
}

func (e ErrBlobMissing) Is(target error) bool {
	_, ok := target.(ErrBlobMissing)
	return ok
}

func IsErrBlobMissing(err error) bool {
	return errors.Is(err, ErrBlobMissing{})
}
