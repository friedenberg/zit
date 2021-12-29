package commands

import (
	"flag"
)

type Checkout struct {
	All         bool
	IncludeAkte bool
}

func init() {
	registerCommand(
		"checkout",
		func(f *flag.FlagSet) Command {
			c := &Checkout{}

			f.BoolVar(&c.All, "all", false, "include all zettels in the current directory")
			f.BoolVar(&c.IncludeAkte, "include-akte", false, "check out akte as well")

			return commandWithZettels{c}
		},
	)
}

func (c Checkout) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
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
