package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
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
		var stored zettel_transacted.Transacted

		if stored, err = store.StoreObjekten().Read(sha); err != nil {
			err = errors.Error(err)
			return
		}

		var tz zettel_transacted.Transacted

		if tz, err = store.StoreObjekten().Create(stored.Named.Stored.Zettel); err != nil {
			err = errors.Error(err)
			return
		}

		stdprinter.Outf("%s (adopted)\n", tz.Named)
	}

	return
}
