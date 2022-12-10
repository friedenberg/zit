package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/papa/umwelt"
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
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
	)

	return
}

func (c Copy) RunWithIds(s *umwelt.Umwelt, ids id_set.Set) (err error) {
	hins := ids.Hinweisen()

	zettels := make([]zettel.Transacted, hins.Len())

	for i, h := range hins.Elements() {
		var tz zettel.Transacted

		if tz, err = s.StoreObjekten().Zettel().ReadHinweisSchwanzen(h); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i] = tz
	}

	return
}
