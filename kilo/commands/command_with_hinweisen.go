package commands

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type CommandWithHinweisen interface {
	RunWithHinweisen(store_with_lock.Store, ...hinweis.Hinweis) error
}

type commandWithHinweisen struct {
	CommandWithHinweisen
}

func (c commandWithHinweisen) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	ids := make([]hinweis.Hinweis, len(args))

	for i, arg := range args {
		var h hinweis.Hinweis

		if h, err = hinweis.MakeBlindHinweis(arg); err != nil {
			err = errors.Error(err)
			return
		}

		ids[i] = h
	}

	c.RunWithHinweisen(store, ids...)

	return
}
