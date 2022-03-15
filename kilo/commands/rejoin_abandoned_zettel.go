package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/india/store_with_lock"
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

func (c RejoinAbandonedZettel) RunWithShas(store store_with_lock.Store, shas ..._Sha) (err error) {
	for _, sha := range shas {
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

		_Outf("[%s %s] (adopted)\n", named.Hinweis, named.Sha)
	}

	return
}
