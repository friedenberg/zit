package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
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

		var stored stored_zettel.Stored

		if stored, err = store.Zettels().ReadZettel(sha); err != nil {
			err = errors.Error(err)
			return
		}

		var named stored_zettel.Named

		if named, err = store.Zettels().Create(stored.Zettel); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = named
		stdprinter.Outf("[%s %s] (adopted)\n", named.Hinweis, named.Sha)
	}

	return
}
