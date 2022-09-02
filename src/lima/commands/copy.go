package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
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
	zettels := make([]zettel_transacted.Zettel, len(hins))

	for i, h := range hins {
		var tz zettel_transacted.Zettel

		if tz, err = s.StoreObjekten().Read(h); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i] = tz
	}

	return
}
