package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
)

type Copy struct {
	Edit bool
}

func init() {
	registerCommand(
		"cp",
		func(f *flag.FlagSet) Command {
			c := &Copy{}

			return commandWithLockedStore{commandWithHinweisen{c}}
		},
	)
}

func (c Copy) RunWithHinweisen(u _Umwelt, zs _Zettels, hins ..._Hinweis) (err error) {
	zettels := make([]_NamedZettel, len(hins))

	for i, h := range hins {
		var named _NamedZettel

		if named, err = zs.Read(h); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = named
	}

	return
}
