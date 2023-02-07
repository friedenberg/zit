package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
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

func (c ExpandHinweis) ProtoIdSet(u *umwelt.Umwelt) (is kennung.ProtoIdSet) {
	is = kennung.MakeProtoIdSet(
		kennung.ProtoId{
			Setter: &kennung.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h kennung.Hinweis
				h, err = u.StoreObjekten().GetAbbrStore().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
	)

	return
}

func (c ExpandHinweis) RunWithIds(s *umwelt.Umwelt, ids kennung.Set) (err error) {
	hins := ids.Hinweisen.Copy()

	for _, h := range hins.Elements() {
		errors.Out().Print(h)
	}

	return
}
