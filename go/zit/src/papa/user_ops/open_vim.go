package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type OpenVim struct {
	Options []string
}

func (c OpenVim) Run(
	u *env.Env,
	args ...string,
) (err error) {
	op := store_fs.Open{
		Options: c.Options,
	}

	if err = op.Run(u.PrinterHeader(), args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
