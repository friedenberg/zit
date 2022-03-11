package commands

import (
	"flag"

	"github.com/friedenberg/zit/india/store_with_lock"
)

type Checkout struct {
	All         bool
	IncludeAkte bool
	Force       bool
}

func init() {
	registerCommand(
		"checkout",
		func(f *flag.FlagSet) Command {
			c := &Checkout{}

			f.BoolVar(&c.All, "all", false, "include all zettels in the current directory")
			f.BoolVar(&c.IncludeAkte, "include-akte", false, "check out akte as well")
			f.BoolVar(&c.Force, "force", false, "force update checked out zettels, even if they will overwrite existing checkouts")

			return commandWithLockedStore{c}
		},
	)
}

func (c Checkout) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	if len(args) == 0 {
		if c.All {
			var hins []_Hinweis

			if _, hins, err = store.Hinweisen().All(); err != nil {
				err = _Error(err)
				return
			}

			for _, h := range hins {
				args = append(args, h.String())
			}
		} else {
			_Errf("nothing to checkout\n")
			return
		}
	}

	var checkedOut map[_Hinweis]_ZettelCheckedOut

	checkinOptions := _ZettelsCheckinOptions{
		IgnoreMissingHinweis: true,
		AddMdExtension:       true,
		IncludeAkte:          c.IncludeAkte,
		Format:               _ZettelFormatsText{},
	}

	if checkedOut, err = store.Zettels().ReadCheckedOut(checkinOptions, args...); err != nil {
		err = _Error(err)
		return
	}

	toCheckOut := make([]string, 0, len(args))

	for h, cz := range checkedOut {
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

	//TODO use user_op
	if _, err = store.Zettels().Checkout(options, args...); err != nil {
		err = _Error(err)
		return
	}

	return
}
