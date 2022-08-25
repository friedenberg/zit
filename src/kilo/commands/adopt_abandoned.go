package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type AdoptAbandoned struct {
}

func init() {
	registerCommand(
		"adopt-abandoned",
		func(f *flag.FlagSet) Command {
			c := &AdoptAbandoned{}

			return commandWithLockedStore{c}
		},
	)
}

func (c AdoptAbandoned) Description() string {
	return "creates a new hinweis for a zettel that has somehow gotten detached"
}

func (c AdoptAbandoned) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	zettels := make([]stored_zettel.Named, len(args))

	for i, a := range args {
		var sha sha.Sha

		if err = sha.Set(a); err != nil {
			err = errors.Error(err)
			return
		}

		var stored stored_zettel.Transacted

		if stored, err = store.Zettels().Read(sha); err != nil {
			err = errors.Error(err)
			return
		}

		var tz stored_zettel.Transacted

		if tz, err = store.Zettels().Create(stored.Named.Stored.Zettel); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = tz.Named
		stdprinter.Outf("%s (adopted)\n", tz.Named)
	}

	return
}
