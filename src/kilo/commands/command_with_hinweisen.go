package commands

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/kilo/user_ops"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
)

type CommandWithHinweisen interface {
	RunWithHinweisen(store_with_lock.Store, ...hinweis.Hinweis) error
}

type commandWithHinweisen struct {
	CommandWithHinweisen
}

func (c commandWithHinweisen) RunWithLockedStore(
	store store_with_lock.Store,
	args ...string,
) (err error) {
	var hins []hinweis.Hinweis

	if hins, err = (user_ops.GetHinweisenFromArgs{}).RunMany(args...); err != nil {
		err = errors.Error(err)
		return
	}

	c.RunWithHinweisen(store, hins...)

	return
}
