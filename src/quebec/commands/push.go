package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_push"
)

type Push struct {
	gattung.Gattung
}

func init() {
	registerCommand(
		"push",
		func(f *flag.FlagSet) Command {
			c := &Push{
				Gattung: gattung.Zettel,
			}

			f.Var(&c.Gattung, "gattung", "Gattung")

			return c
		},
	)
}

func (c Push) ProtoIdSet(u *umwelt.Umwelt) (is kennung.ProtoIdSet) {
	switch c.Gattung {

	default:
		is = kennung.MakeProtoIdSet(
			kennung.ProtoId{
				Setter: &sha.Sha{},
			},
			kennung.ProtoId{
				Setter: &kennung.Hinweis{},
				Expand: func(v string) (out string, err error) {
					var h kennung.Hinweis
					h, err = u.StoreObjekten().GetAbbrStore().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			kennung.ProtoId{
				Setter: &kennung.Etikett{},
				Expand: func(v string) (out string, err error) {
					var e kennung.Etikett
					e, err = u.StoreObjekten().GetAbbrStore().ExpandEtikettString(v)
					out = e.String()
					return
				},
			},
			kennung.ProtoId{
				Setter: &kennung.Typ{},
			},
			kennung.ProtoId{
				Setter: &ts.Time{},
			},
		)

	case gattung.Typ:
		is = kennung.MakeProtoIdSet(
			kennung.ProtoId{
				Setter: &kennung.Typ{},
			},
		)

	case gattung.Transaktion:
		is = kennung.MakeProtoIdSet(
			kennung.ProtoId{
				Setter: &ts.Time{},
			},
		)
	}

	return
}

func (c Push) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Normalf("must specify kasten to push from")
		return
	}

	from := args[0]

	if len(args) > 1 {
		args = args[1:]
	} else {
		err = errors.Normalf("Nothing to push.")
		return
	}

	ps := c.ProtoIdSet(u)

	var ids kennung.Set

	if ids, err = ps.Make(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	filter := kennung.Filter{
		Set: ids,
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	var client remote_push.Client

	if client, err = remote_push.MakeClient(u, from); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = client.SendNeededSkus(filter); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, client)

	return
}
