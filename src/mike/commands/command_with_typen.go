package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type CommandWithTypen interface {
	RunWithTypen(store *umwelt.Umwelt, typen ...typ.Typ) error
}

type commandWithTypen struct {
	CommandWithTypen
}

func (c commandWithTypen) Run(store *umwelt.Umwelt, args ...string) (err error) {
	ps := id_set.MakeProtoSet(
		id_set.ProtoId{
			MutableId: &typ.Typ{},
		},
	)

	ids := ps.Make(args...)

	if err = c.RunWithTypen(store, ids.Typen()...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
