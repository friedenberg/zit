package commands

import (
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/sha"
	"github.com/friedenberg/zit/delta/hinweis"
	"github.com/friedenberg/zit/delta/id"
	"github.com/friedenberg/zit/india/store_with_lock"
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
			if id, err = hinweis.MakeBlindHinweis(arg); err != nil {
				err = errors.Error(err)
				return
			}
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
