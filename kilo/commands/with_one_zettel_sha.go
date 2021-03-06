package commands

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type WithOneZettelSha interface {
	RunWithZettel(store store_with_lock.Store, zettel ..._NamedZettel) error
}

type withOneZettelSha struct {
	WithOneZettelSha
	Count int
}

func (c withOneZettelSha) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	if len(args) != c.Count {
		err = errors.Errorf("exactly %d argument expected, but got %d\n", c.Count, len(args))
		return
	}

	zettels := make([]_NamedZettel, len(args))

	for i, arg := range args {
		var sha _Sha

		if err = sha.Set(arg); err != nil {
			err = errors.Error(err)
			return
		}

		if zettels[i], err = store.Zettels().Read(sha); err != nil {
			err = errors.Error(err)
			return
		}
	}

	c.RunWithZettel(store, zettels...)

	return
}
