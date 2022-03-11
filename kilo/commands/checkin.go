package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/juliett/user_ops"
)

type Checkin struct {
	Delete     bool
	IgnoreAkte bool
	All        bool
}

func init() {
	registerCommand(
		"checkin",
		func(f *flag.FlagSet) Command {
			c := &Checkin{}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.BoolVar(&c.IgnoreAkte, "ignore-akte", false, "do not change the akte")
			f.BoolVar(&c.All, "all", false, "")

			return commandWithZettels{c}
		},
	)
}

func (c Checkin) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	if c.All {
		if len(args) > 0 {
			_Errf("Ignoring args because -all is set\n")
		}

		var cwd string

		if cwd, err = os.Getwd(); err != nil {
			err = _Error(err)
			return
		}

		if args, err = zs.GetPossibleZettels(cwd); err != nil {
			err = _Error(err)
			return
		}
	}

	checkinOp := user_ops.Checkin{
		Umwelt: u,
		Store:  zs,
		Options: _ZettelsCheckinOptions{
			IncludeAkte: !c.IgnoreAkte,
			Format:      _ZettelFormatsText{},
		},
	}

	var results user_ops.CheckinResults

	if results, err = checkinOp.Run(args...); err != nil {
		err = _Error(err)
		return
	}

	if c.Delete {
		deleteOp := user_ops.DeleteCheckout{}

		if err = deleteOp.Run(results.Zettelen); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
