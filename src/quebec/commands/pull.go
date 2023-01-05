package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/quebec/remote_pull"
)

type Pull struct {
	gattung.Gattung
	All bool
}

func init() {
	registerCommand(
		"pull",
		func(f *flag.FlagSet) Command {
			c := &Pull{
				Gattung: gattung.Zettel,
			}

			f.Var(&c.Gattung, "gattung", "Gattung")
			f.BoolVar(&c.All, "all", false, "pull all Objekten")

			return c
		},
	)
}

func (c Pull) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	switch c.Gattung {

	default:
		is = id_set.MakeProtoIdSet(
			id_set.ProtoId{
				MutableId: &sha.Sha{},
			},
			id_set.ProtoId{
				MutableId: &hinweis.Hinweis{},
				Expand: func(v string) (out string, err error) {
					var h hinweis.Hinweis
					h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			id_set.ProtoId{
				MutableId: &kennung.Etikett{},
				Expand: func(v string) (out string, err error) {
					var e kennung.Etikett
					e, err = u.StoreObjekten().Abbr().ExpandEtikettString(v)
					out = e.String()
					return
				},
			},
			id_set.ProtoId{
				MutableId: &kennung.Typ{},
			},
			id_set.ProtoId{
				MutableId: &ts.Time{},
			},
		)

	case gattung.Typ:
		is = id_set.MakeProtoIdSet(
			id_set.ProtoId{
				MutableId: &kennung.Typ{},
			},
		)

	case gattung.Transaktion:
		is = id_set.MakeProtoIdSet(
			id_set.ProtoId{
				MutableId: &ts.Time{},
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

	defer errors.Deferred(&err, client.Close)

	p := collections.MakePool[zettel.Transacted]()

	inflator := objekte.MakeTransactedInflator[
		zettel.Objekte,
		*zettel.Objekte,
		hinweis.Hinweis,
		*hinweis.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	](
		func(sk sku.SkuLike) (rc sha.ReadCloser, err error) {
			var or io.ReadCloser

			if or, err = client.ObjekteReaderForSku(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			errors.Log().Printf("got reader for sku: %s", sk)

			var ow sha.WriteCloser

			if ow, err = u.StoreObjekten().WriteCloserObjektenGattung(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			rc = sha.MakeReadCloserTee(or, ow)

			return
		},
		client.AkteReader,
		&zettel.FormatObjekte{
			IgnoreTypErrors: true,
		},
		nil,
		// objekte.MakeNopAkteParser[zettel.Objekte, *zettel.Objekte](),
		p,
	)

	if err = client.SkusFromFilter(
		filter,
		func(sk sku.SkuLike) (err error) {
			if sk.GetGattung() != gattung.Zettel {
				return
			}

			if u.StoreObjekten().Zettel().HasObjekte(sk.GetObjekteSha()) {
				errors.Log().Printf("already have objekte: %s", sk.GetObjekteSha())
				return
			}

			errors.Log().Printf("need objekte: %s", sk.GetObjekteSha())

			//TODO-P1 check for akte sha
			//TODO-P1 write akte

			var t *zettel.Transacted

			if t, err = inflator.Inflate(
				ts.Now(),
				sk,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = u.StoreObjekten().Zettel().Inherit(t); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
