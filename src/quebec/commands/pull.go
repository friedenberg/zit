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
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/quebec/remote_pull"
)

type Pull struct {
	GattungSet gattungen.MutableSet
	All        bool
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
			f.BoolVar(&c.All, "all", false, "pull all Objekten")

			return c
		},
	)
}

func (c Pull) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet()

	if c.GattungSet.Contains(gattung.Zettel) {
		is.AddMany(
			id_set.ProtoId{
				Setter: &sha.Sha{},
			},
			id_set.ProtoId{
				Setter: &hinweis.Hinweis{},
				Expand: func(v string) (out string, err error) {
					var h hinweis.Hinweis
					h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			id_set.ProtoId{
				Setter: &kennung.Etikett{},
				Expand: func(v string) (out string, err error) {
					var e kennung.Etikett
					e, err = u.StoreObjekten().Abbr().ExpandEtikettString(v)
					out = e.String()
					return
				},
			},
			id_set.ProtoId{
				Setter: &kennung.Typ{},
			},
			id_set.ProtoId{
				Setter: &ts.Time{},
			},
		)
	}

	if c.GattungSet.Contains(gattung.Typ) {
		is.AddMany(
			id_set.ProtoId{
				Setter: &kennung.Typ{},
			},
		)
	}

	if c.GattungSet.Contains(gattung.Transaktion) {
		is.AddMany(
			id_set.ProtoId{
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

		if c.All {
			errors.Log().Print("-all is set but arguments passed in. Ignore -all.")
		}
	} else if !c.All {
		err = errors.Normalf("Refusing to pull all unless -all is set.")
		return
	} else {
		args = []string{}
	}

	ps := c.ProtoIdSet(u)

	var ids id_set.Set

	if ids, err = ps.Make(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	filter := id_set.Filter{
		AllowEmpty: c.All,
		Set:        ids,
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	var client remote_pull.Client

	if client, err = remote_pull.MakeClient(u, from); err != nil {
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
