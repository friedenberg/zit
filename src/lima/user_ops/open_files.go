package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
)

type OpenFiles struct {
}

func (c OpenFiles) Run(args ...string) (err error) {
	if len(args) == 0 {
		return
	}

	if err = open_file_guard.OpenFiles(args...); err != nil {
		err = errors.Errorf("%q: %s", args, err)
		return
	}

	return
}
