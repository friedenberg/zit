package repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

func (local *Local) CancelWithError(err error) {
	local.Context.Cancel(errors.WrapN(1, err))
}
