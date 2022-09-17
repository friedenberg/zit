package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type WithOneZettelSha interface {
	RunWithZettel(store *umwelt.Umwelt, zettel ...zettel_transacted.Zettel) error
}

type withOneZettelSha struct {
	WithOneZettelSha
	Count int
}

func (c withOneZettelSha) Run(store *umwelt.Umwelt, args ...string) (err error) {
	if len(args) != c.Count {
		err = errors.Errorf("exactly %d argument expected, but got %d", c.Count, len(args))
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
