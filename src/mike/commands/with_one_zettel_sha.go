package commands

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
	"github.com/friedenberg/zit/zettel_transacted"
)

type WithOneZettelSha interface {
	RunWithZettel(store store_with_lock.Store, zettel ...zettel_transacted.Transacted) error
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

	zettels := make([]zettel_transacted.Transacted, len(args))

	for i, arg := range args {
		var sha sha.Sha

		if err = sha.Set(arg); err != nil {
			err = errors.Error(err)
			return
		}

		if zettels[i], err = store.StoreObjekten().Read(sha); err != nil {
			err = errors.Error(err)
			return
		}
	}

	c.RunWithZettel(store, zettels...)

	return
}
