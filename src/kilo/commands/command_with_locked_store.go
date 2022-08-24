package commands

import (
	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type CommandWithLockedStore interface {
	RunWithLockedStore(store_with_lock.Store, ...string) error
}

type commandWithLockedStore struct {
	CommandWithLockedStore
}

func (c commandWithLockedStore) Run(u *umwelt.Umwelt, args ...string) (err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print()

	defer errors.PanicIfError(store.Flush)

	logz.Print()

	if err = c.RunWithLockedStore(store, args...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
