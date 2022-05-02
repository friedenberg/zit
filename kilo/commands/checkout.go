package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type Checkout struct {
	All         bool
	IncludeAkte bool
	Force       bool
	Etiketten   bool
}

func init() {
	registerCommand(
		"checkout",
		func(f *flag.FlagSet) Command {
			c := &Checkout{}

			f.BoolVar(&c.All, "all", false, "include all zettels in the current directory")
			f.BoolVar(&c.IncludeAkte, "include-akte", false, "check out akte as well")
			f.BoolVar(&c.Force, "force", false, "force update checked out zettels, even if they will overwrite existing checkouts")
			f.BoolVar(&c.Etiketten, "etiketten", false, "treat the arguments as Etiketten instead of Hinweisen")

			return c
		},
	)
}

func (c Checkout) ValidateArgs(args ...string) (err error) {
	if c.All && c.Etiketten {
		err = errors.Errorf("cannot have -all and -etiketten set")
		return
	}

	if c.All && len(args) > 0 {
		err = errors.Errorf("cannot have args when -all is set")
		return
	}

	if !c.All && len(args) == 0 {
		err = errors.Errorf("nothing to check out")
		return
	}

	return
}

func (c Checkout) Run(u _Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		if c.All {
			getHinweisenOp := user_ops.GetAllHinweisen{
				Umwelt: u,
			}

			var getHinweisenResults user_ops.GetAllHinweisenResults

			if getHinweisenResults, err = getHinweisenOp.Run(); err != nil {
				err = errors.Error(err)
				return
			}

			args = getHinweisenResults.HinweisenStrings
		}
	}

	checkinOptions := _ZettelsCheckinOptions{
		IgnoreMissingHinweis: true,
		AddMdExtension:       true,
		IncludeAkte:          c.IncludeAkte,
		Format:               _ZettelFormatsText{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  u,
		Options: checkinOptions,
	}

	if readResults, err = readOp.Run(args...); err != nil {
		err = errors.Error(err)
		return
	}

	toCheckOut := make([]string, 0, len(args))

	for h, cz := range readResults.Zettelen {
		if cz.External.Path == "" {
			toCheckOut = append(toCheckOut, h.String())
			continue
		}

		if cz.Internal.Zettel.Equals(cz.External.Zettel) {
			_Outf("[%s %s] (already checked out)\n", cz.Internal.Hinweis, cz.Internal.Sha)
			continue
		}

		if c.Force {
			toCheckOut = append(toCheckOut, h.String())
		} else {
			_Errf("[%s] (external has changes)\n", h)
			continue
		}
	}

	options := _ZettelsCheckinOptions{
		IncludeAkte: c.IncludeAkte,
		Format:      _ZettelFormatsText{},
	}

	checkoutOp := user_ops.Checkout{
		Umwelt:  u,
		Options: options,
	}

	if _, err = checkoutOp.Run(args...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
