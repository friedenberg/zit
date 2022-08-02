package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/hotel/zettels"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type checkoutArgType int

const (
	checkoutArgTypeNormal = checkoutArgType(iota)
	checkoutArgTypeAll
	checkoutArgTypeEtiketten
)

type argOptions struct {
	All       bool
	Etiketten bool
}

func (o argOptions) validateArgs(args ...string) (t checkoutArgType, err error) {
	if o.All {
		if o.Etiketten {
			err = errors.Errorf("cannot have -all and -etiketten set")
			return
		}

		if len(args) > 0 {
			err = errors.Errorf("cannot have args when -all is set")
			return
		}

		t = checkoutArgTypeAll
	} else if o.Etiketten {
		t = checkoutArgTypeEtiketten
	}

	return
}

type Checkout struct {
	argOptions
	IncludeAkte bool
	Force       bool
}

func init() {
	registerCommand(
		"checkout",
		func(f *flag.FlagSet) Command {
			c := &Checkout{}

			f.BoolVar(&c.argOptions.All, "all", false, "include all zettels in the current directory")
			f.BoolVar(&c.argOptions.Etiketten, "etiketten", false, "treat the arguments as Etiketten instead of Hinweisen")
			f.BoolVar(&c.IncludeAkte, "include-akte", false, "check out akte as well")
			f.BoolVar(&c.Force, "force", false, "force update checked out zettels, even if they will overwrite existing checkouts")

			return c
		},
	)
}

func (c Checkout) Run(u *umwelt.Umwelt, args ...string) (err error) {
	var t checkoutArgType

	if t, err = c.argOptions.validateArgs(args...); err != nil {
		err = errors.Error(err)
		return
	}

	switch t {
	case checkoutArgTypeAll:
		getHinweisenOp := user_ops.GetAllHinweisen{
			Umwelt: u,
		}

		var getHinweisenResults user_ops.GetAllHinweisenResults

		if getHinweisenResults, err = getHinweisenOp.Run(); err != nil {
			err = errors.Error(err)
			return
		}

		args = getHinweisenResults.HinweisenStrings

	case checkoutArgTypeEtiketten:
		queryOp := user_ops.GetZettelsFromQuery{
			Umwelt: u,
		}

		var zettelResults user_ops.ZettelResults

		var slice etikett.Slice

		if slice, err = etikett.NewSliceFromStrings(args...); err != nil {
			err = errors.Error(err)
			return
		}

		query := stored_zettel.FilterEtikettSet{
			Set: slice.ToSet(),
		}

		if zettelResults, err = queryOp.Run(query); err != nil {
			err = errors.Error(err)
			return
		}

		args = zettelResults.HinweisStrings()

	case checkoutArgTypeNormal:
	default:
		break
	}

	checkinOptions := zettels.CheckinOptions{
		IgnoreMissingHinweis: true,
		AddMdExtension:       true,
		IncludeAkte:          c.IncludeAkte,
		Format:               zettel_formats.Text{},
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
			stdprinter.Outf("[%s %s] (already checked out)\n", cz.Internal.Hinweis, cz.Internal.Sha)
			continue
		}

		if c.Force {
			toCheckOut = append(toCheckOut, h.String())
		} else {
			stdprinter.Errf("[%s] (external has changes)\n", h)
			continue
		}
	}

	options := zettels.CheckinOptions{
		IncludeAkte: c.IncludeAkte,
		Format:      zettel_formats.Text{},
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
