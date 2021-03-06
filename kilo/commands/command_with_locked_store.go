package commands

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type CommandWithLockedStore interface {
	RunWithLockedStore(store_with_lock.Store, ...string) error
}

type commandWithLockedStore struct {
	CommandWithLockedStore
}

func (c commandWithLockedStore) Run(u _Umwelt, args ...string) (err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	if err = c.RunWithLockedStore(store, args...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
