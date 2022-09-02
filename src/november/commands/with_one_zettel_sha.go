package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
)

type WithOneZettelSha interface {
	RunWithZettel(store store_with_lock.Store, zettel ...zettel_transacted.Zettel) error
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

	zettels := make([]zettel_transacted.Zettel, len(args))

	for i, arg := range args {
		var sha sha.Sha

		if err = sha.Set(arg); err != nil {
			err = errors.Wrap(err)
			return
		}

		if zettels[i], err = store.StoreObjekten().Read(sha); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	c.RunWithZettel(store, zettels...)

	return
}
