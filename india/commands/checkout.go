package commands

import (
	"flag"
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

			return commandWithZettels{c}
		},
	)
}

func (c Checkout) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	u.Lock.Lock()
	defer _PanicIfError(u.Lock.Unlock())

	if len(args) == 0 {
		if c.All {
			var hins []_Hinweis

			if _, hins, err = zs.Hinweisen().All(); err != nil {
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

	if checkedOut, err = zs.ReadCheckedOut(checkinOptions, args...); err != nil {
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

	if _, err = zs.Checkout(options, args...); err != nil {
		err = _Error(err)
		return
	}

	return
}
