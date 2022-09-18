package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

var defaultProtoSet id_set.ProtoSet

func init() {
}

type CommandWithIds interface {
	RunWithIds(store *umwelt.Umwelt, ids id_set.Set) error
}

type commandWithIds struct {
	CommandWithIds
	id_set.ProtoSet
}

func (c commandWithIds) Run(store *umwelt.Umwelt, args ...string) (err error) {
	ps := c.ProtoSet

	if ps.Len() == 0 {
		ps = id_set.MakeProtoSet(
			id_set.ProtoId{
				MutableId: &sha.Sha{},
			},
			id_set.ProtoId{
				MutableId: &hinweis.Hinweis{},
				Expand: func(v string) (out string, err error) {
					var h hinweis.Hinweis
					h, err = store.StoreObjekten().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			id_set.ProtoId{
				MutableId: &ts.Time{},
			},
		)
	}

	ids := ps.Make(args...)

	if err = c.RunWithIds(store, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
