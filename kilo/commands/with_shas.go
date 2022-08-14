package commands

import (
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/sha"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type WithShas interface {
	RunWithShas(store store_with_lock.Store, shas ...sha.Sha) error
}

type withShas struct {
	WithShas
}

func (c withShas) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	shas := make([]sha.Sha, len(args))

	for i, arg := range args {
		var sha sha.Sha

		if err = sha.Set(arg); err != nil {
			err = errors.Error(err)
			return
		}

		shas[i] = sha
	}

	if err = c.RunWithShas(store, shas...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
