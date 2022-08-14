package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/india/store_with_lock"
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
		var stored stored_zettel.Transacted

		if stored, err = store.Zettels().Read(sha); err != nil {
			err = errors.Error(err)
			return
		}

		var tz stored_zettel.Transacted

		if tz, err = store.Zettels().Create(stored.Zettel); err != nil {
			err = errors.Error(err)
			return
		}

		stdprinter.Outf("%s (adopted)\n", tz.Named)
	}

	return
}
