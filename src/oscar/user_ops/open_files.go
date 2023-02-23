package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type OpenFiles struct{}

func (c OpenFiles) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		return
	}

  errors.Err().Printf("files %#v", args)

	if err = files.OpenFiles(args...); err != nil {
		err = errors.Wrapf(err, "%q", args)
		return
	}

	v := "opening files"

	if err = u.PrinterHeader()(&v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
