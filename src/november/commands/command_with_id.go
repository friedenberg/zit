package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
)

type CommandWithId interface {
	RunWithId(store store_with_lock.Store, ids ...id_set.Set) error
}

type commandWithId struct {
	CommandWithId
}

func (c commandWithId) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	ps := id_set.MakeProtoSet(
		&sha.Sha{},
		&hinweis.Hinweis{},
		&hinweis.HinweisWithIndex{},
		&ts.Time{},
	)

	ids := ps.MakeMany(args...)

	if err = c.RunWithId(store, ids...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
