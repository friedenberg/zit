package commands

import (
	"flag"
)

type Copy struct {
	Edit bool
}

func init() {
	registerCommand(
		"cp",
		func(f *flag.FlagSet) Command {
			c := &Copy{}

			return commandWithZettels{commandWithHinweisen{c}}
		},
	)
}

func (c Copy) RunWithHinweisen(u _Umwelt, zs _Zettels, hins ..._Hinweis) (err error) {
	zettels := make([]_NamedZettel, len(hins))

	for i, h := range hins {
		var named _NamedZettel

		if named, err = zs.Read(h); err != nil {
			err = _Error(err)
			return
		}

		zettels[i] = named
	}

	return
}
