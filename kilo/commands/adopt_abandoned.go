package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
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
	zettels := make([]_NamedZettel, len(args))

	for i, a := range args {
		var sha _Sha

		if err = sha.Set(a); err != nil {
			err = errors.Error(err)
			return
		}

		var stored _StoredZettel

		if stored, err = store.Zettels().ReadZettel(sha); err != nil {
			err = errors.Error(err)
			return
		}

		var named _NamedZettel

		if named, err = store.Zettels().Create(stored.Zettel); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = named
		_Outf("[%s %s] (adopted)\n", named.Hinweis, named.Sha)
	}

	return
}
