package commands

import (
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/delta/hinweis"
	"github.com/friedenberg/zit/india/store_with_lock"
	"github.com/friedenberg/zit/juliett/user_ops"
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
