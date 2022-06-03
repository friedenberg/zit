package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type ReunifyFamily struct {
}

func init() {
	registerCommand(
		"reunify-family",
		func(f *flag.FlagSet) Command {
			c := &ReunifyFamily{}

			return commandWithLockedStore{withOneZettelSha{WithOneZettelSha: c, Count: 2}}
		},
	)
}

func (c ReunifyFamily) RunWithZettel(store store_with_lock.Store, zettel ..._NamedZettel) (err error) {
	mutter := zettel[0]
	kinder := zettel[1]

	mutter.Kinder = kinder.Sha
	kinder.Mutter = mutter.Sha

	if err = store.Zettels().UpdateNoKinder(mutter); err != nil {
		err = errors.Error(err)
		return
	}

	if err = store.Zettels().UpdateNoKinder(kinder); err != nil {
		err = errors.Error(err)
		return
	}

	_Outf("%#v\n", mutter)
	_Outf("%#v\n", kinder)

	return
}
