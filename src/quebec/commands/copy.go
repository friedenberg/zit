package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Copy struct {
	Edit bool
}

func init() {
	registerCommand(
		"cp",
		func(f *flag.FlagSet) Command {
			c := &Copy{}

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c Copy) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			Setter: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().GetAbbrStore().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
	)

	return
}

func (c Copy) RunWithIds(s *umwelt.Umwelt, ids id_set.Set) (err error) {
	hins := ids.Hinweisen.Copy()

	zettels := make([]*zettel.Transacted, hins.Len())

	for i, h := range hins.Elements() {
		var tz *zettel.Transacted

		if tz, err = s.StoreObjekten().Zettel().ReadOne(h); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i] = tz
	}

	return
}
