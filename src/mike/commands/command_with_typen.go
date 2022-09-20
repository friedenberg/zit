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
	ps := id_set.MakeProtoIdList(
		id_set.ProtoId{
			MutableId: &typ.Typ{},
		},
	)

	var ids id_set.Set

	if ids, err = ps.Make(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithTypen(store, ids.Typen()...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
