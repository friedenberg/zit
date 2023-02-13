package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
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
				MutableSet:   &c.GattungSet,
				SetterPolicy: collections.SetterPolicyReset,
			}

			f.Var(gsvs, "gattung", "Gattung")

			return c
		},
	)
}

func (c Pull) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
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

	ids := u.MakeMetaIdSet()

	if err = ids.SetMany(args...); err != nil {
		err = errors.Wrap(err)
		return
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

	if err = client.PullSkus(ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
