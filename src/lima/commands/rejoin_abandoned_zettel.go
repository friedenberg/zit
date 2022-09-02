package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
)

type RejoinAbandonedZettel struct {
}

func init() {
	registerCommand(
		"rejoin-abandoned-zettel",
		func(f *flag.FlagSet) Command {
			c := &RejoinAbandonedZettel{}

			return commandWithLockedStore{withShas{c}}
		},
	)
}

func (c RejoinAbandonedZettel) RunWithShas(store store_with_lock.Store, shas ...sha.Sha) (err error) {
	for _, sha := range shas {
		var stored zettel_transacted.Zettel

		if stored, err = store.StoreObjekten().Read(sha); err != nil {
			err = errors.Wrap(err)
			return
		}

		var tz zettel_transacted.Zettel

		if tz, err = store.StoreObjekten().Create(stored.Named.Stored.Zettel); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.PrintOutf("%s (adopted)", tz.Named)
	}

	return
}
