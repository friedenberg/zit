package commands

import (
	"flag"

	"github.com/friedenberg/zit/juliett/user_ops"
)

type Edit struct {
	IncludeAkte bool
}

func init() {
	registerCommand(
		"edit",
		func(f *flag.FlagSet) Command {
			c := &Edit{}

			f.BoolVar(&c.IncludeAkte, "include-akte", true, "check out and open the akte")

			return c
		},
	)
}

func (c Edit) Run(u _Umwelt, args ...string) (err error) {
	checkoutOp := user_ops.Checkout{
		Umwelt: u,
		Options: _ZettelsCheckinOptions{
			IncludeAkte: c.IncludeAkte,
			Format:      _ZettelFormatsText{},
		},
	}

	var checkoutResults user_ops.CheckoutResults

	if checkoutResults, err = checkoutOp.Run(args...); err != nil {
		err = _Error(err)
		return
	}

	if err = (user_ops.OpenFiles{}).Run(checkoutResults.FilesAkten...); err != nil {
		err = _Error(err)
		return
	}

	vimArgs := []string{
		"-c",
		"set ft=zit.zettel",
		"-c",
		"source ~/.vim/syntax/zit.zettel.vim",
	}

	if err = _OpenVimWithArgs(vimArgs, checkoutResults.FilesZettelen...); err != nil {
		err = _Error(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:  u,
		Options: checkoutOp.Options,
	}

	if _, err = checkinOp.Run(checkoutResults.FilesZettelen...); err != nil {
		err = _Error(err)
		return
	}

	return
}
