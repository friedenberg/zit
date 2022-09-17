package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
)

type OpenFiles struct {
}

func (c OpenFiles) Run(args ...string) (err error) {
	if len(args) == 0 {
		return
	}

	if err = files.OpenFiles(args...); err != nil {
		err = errors.Wrapf(err, "%q", args)
		return
	}

	return
}
