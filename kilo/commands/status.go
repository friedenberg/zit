package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/golf/checkout_store"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type Status struct {
}

func init() {
	registerCommand(
		"status",
		func(f *flag.FlagSet) Command {
			c := &Status{}

			return c
		},
	)
}

func (c Status) Run(u _Umwelt, args ...string) (err error) {
	if len(args) > 0 {
		stdprinter.Errf("args provided will be ignored")
	}

	var possible checkout_store.CwdFiles

	if possible, err = user_ops.NewGetPossibleZettels(u).Run(); err != nil {
		err = errors.Error(err)
		return
	}

	args = possible.Zettelen

	options := _ZettelsCheckinOptions{
		IncludeAkte: true,
		Format:      _ZettelFormatsText{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  u,
		Options: options,
	}

	if readResults, err = readOp.Run(args...); err != nil {
		err = errors.Error(err)
		return
	}

	for h, z := range readResults.Zettelen {
		if z.Internal.Zettel.Equals(z.External.Zettel) {
			continue
		}

		stdprinter.Outf("[%s] (different)\n", h)
	}

	return
}
