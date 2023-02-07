package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_transfers"
)

type Pull struct {
	GattungSet gattungen.MutableSet
}

func init() {
	registerCommand(
		"pull",
		func(f *flag.FlagSet) Command {
			c := &Pull{
				GattungSet: gattungen.MakeMutableSet(gattung.Zettel),
			}

			gsvs := collections.MutableValueSet2[gattung.Gattung, *gattung.Gattung]{
				MutableSetLike: &c.GattungSet,
				SetterPolicy:   collections.SetterPolicyReset,
			}

			f.Var(gsvs, "gattung", "Gattung")

			return c
		},
	)
}

func (c Pull) ProtoIdSet(u *umwelt.Umwelt) (is kennung.ProtoIdSet) {
	is = kennung.MakeProtoIdSet()

	if c.GattungSet.Contains(gattung.Zettel) {
		is.AddMany(
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
	}

	if c.GattungSet.Contains(gattung.Typ) {
		is.AddMany(
			kennung.ProtoId{
				Setter: &kennung.Typ{},
			},
		)
	}

	if c.GattungSet.Contains(gattung.Transaktion) {
		is.AddMany(
			kennung.ProtoId{
				Setter: &ts.Time{},
			},
		)
	}

	return
}

func (c Pull) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Normalf("must specify kasten to pull from")
		return
	}

	from := args[0]

	if len(args) > 1 {
		args = args[1:]
	} else {
		err = errors.Normalf("Nothing to pull")
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

	var client remote_transfers.PullClient

	if client, err = remote_transfers.MakePullClient(u, from); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, client)

	if err = client.PullSkus(filter, c.GattungSet.Copy()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
