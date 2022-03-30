package commands

import (
	"flag"

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
		_Errf("args provided will be ignored")
	}

	getPossibleOp := user_ops.GetPossibleZettels{
		Umwelt: u,
	}

	var getPossibleResults user_ops.GetPossibleZettelsResults

	if getPossibleResults, err = getPossibleOp.Run(); err != nil {
		err = _Error(err)
		return
	}

	args = getPossibleResults.Hinweisen

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
		err = _Error(err)
		return
	}

	for h, z := range readResults.Zettelen {
		if z.Internal.Zettel.Equals(z.External.Zettel) {
			continue
		}

		_Outf("[%s] (different)\n", h)
	}

	return
}
