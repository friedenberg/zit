package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
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
	zettels := make([]zettel_stored.Transacted, len(hins))

	for i, h := range hins {
		var tz zettel_stored.Transacted

		if tz, err = s.StoreObjekten().Read(h); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = tz
	}

	return
}
