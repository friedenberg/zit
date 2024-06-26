package user_ops

import (
	"fmt"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"github.com/google/shlex"
)

type EachAkte struct {
	*umwelt.Umwelt
	Utility string
}

func (c EachAkte) Run(
	zsc sku.CheckedOutLikeSet,
) (err error) {
	if zsc.Len() == 0 {
		return
	}

	var akten []string

	if err = zsc.Each(
		func(col sku.CheckedOutLike) (err error) {
			cofs := col.(*store_fs.CheckedOut)
			akten = append(akten, cofs.External.GetAkteFD().GetPath())

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	v := fmt.Sprintf("running utility: %q", c.Utility)

	if err = c.PrinterHeader()(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	var args []string

	if args, err = shlex.Split(c.Utility); err != nil {
		err = errors.Wrap(err)
		return
	}

	args = append(args, akten...)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = c.Out()
	cmd.Stdin = c.In()
	cmd.Stderr = c.Err()

	if err = cmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
