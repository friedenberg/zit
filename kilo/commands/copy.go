package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
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

func (c Copy) RunWithHinweisen(s store_with_lock.Store, hins ...hinweis.Hinweis) (err error) {
	zettels := make([]stored_zettel.Transacted, len(hins))

	for i, h := range hins {
		var tz stored_zettel.Transacted

		if tz, err = s.Zettels().Read(h); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = tz
	}

	return
}
