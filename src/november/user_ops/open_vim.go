package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
)

type OpenVim struct {
	Options []string
}

type OpenVimResults struct {
}

func (c OpenVim) Run(args ...string) (results OpenVimResults, err error) {
	if len(args) == 0 {
		err = errors.Errorf("no files to edit in vim")
		return
	}

	vimArgs := make([]string, 0, len(c.Options)*2)

	for _, o := range c.Options {
		vimArgs = append(vimArgs, "-c", o)
	}

	if err = files.OpenVimWithArgs(vimArgs, args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
