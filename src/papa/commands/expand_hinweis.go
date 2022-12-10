package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type ExpandHinweis struct {
}

func init() {
	registerCommand(
		"expand-hinweis",
		func(f *flag.FlagSet) Command {
			c := &ExpandHinweis{}

			return commandWithIds{CommandWithIds: c}
		},
	)
}

func (c ExpandHinweis) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
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

func (c ExpandHinweis) RunWithIds(s *umwelt.Umwelt, ids id_set.Set) (err error) {
	hins := ids.Hinweisen()

	for _, h := range hins.Elements() {
		errors.Out().Print(h)
	}

	return
}
