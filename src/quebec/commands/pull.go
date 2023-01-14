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
	All        bool
	RewriteTai bool
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
			f.BoolVar(&c.RewriteTai, "rewrite-tai", false, "generate new Taimstamps for pulled Objektes")

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

	defer errors.DeferredCloser(&err, client)

	p := collections.MakePool[zettel.Transacted]()

	inflator := objekte.MakeTransactedInflator[
		zettel.Objekte,
		*zettel.Objekte,
		hinweis.Hinweis,
		*hinweis.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	](
		client.ObjekteReader,
		func(sh sha.Sha) (rc sha.ReadCloser, err error) {
			errors.Todo(errors.P2, "move to own constructor")
			var or io.ReadCloser

			if or, err = client.AkteReader(sh); err != nil {
				err = errors.Wrap(err)
				return
			}

			errors.Log().Printf("got reader for sha: %s", sh)

			var ow sha.WriteCloser

			if ow, err = u.StoreObjekten().AkteWriter(); err != nil {
				err = errors.Wrap(err)
				return
			}

			rc = sha.MakeReadCloserTee(or, ow)

			return
		},
		&zettel.FormatObjekte{
			IgnoreTypErrors: true,
		},
		// objekte.MakeParserStorerWithCustomFormat[zettel.Objekte, *zettel.Objekte](
		// 	u.StoreObjekten(),
		// 	&zettel.FormatObjekte{
		// 		IgnoreTypErrors: true,
		// 	},
		// ),
		objekte.MakeNopAkteParser[zettel.Objekte, *zettel.Objekte](),
		p,
	)

	if err = client.SkusFromFilter(
		filter,
		func(sk sku.Sku2) (err error) {
			if sk.Gattung != gattung.Zettel {
				return
			}

			//TODO-P1 check for akte sha
			//TODO-P1 write akte
			if u.Standort().HasObjekte(sk.Gattung, sk.ObjekteSha) {
				errors.Log().Printf("already have objekte: %s", sk.ObjekteSha)
				return
			}

			errors.Log().Printf("need objekte: %s", sk.ObjekteSha)

			var t *zettel.Transacted

			if t, err = inflator.InflateFromSku2(sk); err != nil {
				err = errors.Wrapf(err, "Sku: %s", sk)
				return
			}

			f := &zettel.FormatObjekte{
				IgnoreTypErrors: true,
			}

			var ow sha.WriteCloser

			if ow, err = u.StoreObjekten().ObjekteWriter(t.GetGattung()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, ow)

			if _, err = f.Format(ow, &t.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

			t.Sku.ObjekteSha = ow.Sha()
			t.Sku.AkteSha = t.AkteSha()

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
