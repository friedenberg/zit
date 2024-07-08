package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type OpenVim struct {
	Options []string
}

func (c OpenVim) Run(
	u *umwelt.Umwelt,
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
