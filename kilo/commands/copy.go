package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/hotel/zettels"
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

func (c Copy) RunWithHinweisen(u *umwelt.Umwelt, zs zettels.Zettels, hins ...hinweis.Hinweis) (err error) {
	zettels := make([]stored_zettel.Named, len(hins))

	for i, h := range hins {
		var named stored_zettel.Named

		if named, err = zs.Read(h); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = named
	}

	return
}
