package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type OpenVim struct {
	Options []string
}

type OpenVimResults struct {
}

func (c OpenVim) Run(
	u *umwelt.Umwelt,
	args ...string,
) (results OpenVimResults, err error) {
	vimArgs := make([]string, 0, len(c.Options)*2)

	for _, o := range c.Options {
		vimArgs = append(vimArgs, "-c", o)
	}

	v := "vim started"

	if err = u.PrinterHeader()(&v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = files.OpenVimWithArgs(vimArgs, args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	v = "vim exited"

	if err = u.PrinterHeader()(&v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
