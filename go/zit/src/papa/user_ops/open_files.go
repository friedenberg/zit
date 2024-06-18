package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type OpenFiles struct{}

func (c OpenFiles) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		return
	}

	if err = files.OpenFiles(args...); err != nil {
		err = errors.Wrapf(err, "%q", args)
		return
	}

	v := "opening files"

	if err = u.PrinterHeader()(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
