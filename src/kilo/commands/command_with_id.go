package commands

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
)

type CommandWithId interface {
	RunWithId(store store_with_lock.Store, ids ...id.Id) error
}

type commandWithId struct {
	CommandWithId
}

func (c commandWithId) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	ids := make([]id.Id, len(args))

	for i, arg := range args {
		var id id.Id
		var sha sha.Sha

		if err = sha.Set(arg); err != nil {
			var hwi hinweis.HinweisWithIndex

			if err = hwi.Set(arg); err != nil {
				err = errors.Error(err)
				return
			}

			id = hwi
		} else {
			id = sha
		}

		ids[i] = id
	}

	if err = c.RunWithId(store, ids...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
