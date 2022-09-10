package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type CommandWithId interface {
	RunWithId(store *umwelt.Umwelt, ids ...id_set.Set) error
}

type commandWithId struct {
	CommandWithId
}

func (c commandWithId) Run(store *umwelt.Umwelt, args ...string) (err error) {
	ps := id_set.MakeProtoSet(
		&sha.Sha{},
		&hinweis.Hinweis{},
		&hinweis.HinweisWithIndex{},
		&ts.Time{},
	)

	ids := ps.MakeMany(args...)

	if err = c.RunWithId(store, ids...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
