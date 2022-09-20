package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type CommandWithIds interface {
	RunWithIds(store *umwelt.Umwelt, ids id_set.Set) error
}

type commandWithIds struct {
	CommandWithIds
	id_set.ProtoSet
}

func (c commandWithIds) getIdProtoSet(u *umwelt.Umwelt) (is id_set.ProtoSet) {
	is = c.ProtoSet

	if is.Len() == 0 {
		is = id_set.MakeProtoSet(
			id_set.ProtoId{
				MutableId: &sha.Sha{},
			},
			id_set.ProtoId{
				MutableId: &hinweis.Hinweis{},
				Expand: func(v string) (out string, err error) {
					var h hinweis.Hinweis
					h, err = u.StoreObjekten().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			id_set.ProtoId{
				MutableId: &ts.Time{},
			},
		)
	}

	return
}

func (c commandWithIds) Run(u *umwelt.Umwelt, args ...string) (err error) {
	ps := c.getIdProtoSet(u)

	ids := ps.Make(args...)

	if err = c.RunWithIds(u, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
