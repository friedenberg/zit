package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/konfig"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/golf/typ"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Show struct {
	gattung.Gattung
	Format string
}

func init() {
	registerCommand(
		"show",
		func(f *flag.FlagSet) Command {
			c := &Show{
				Gattung: gattung.Zettel,
			}

			f.Var(&c.Gattung, "gattung", "Gattung")
			f.StringVar(&c.Format, "format", "text", "format")

			cwi := commandWithIds{
				CommandWithIds: c,
			}

			return CommandV2{
				Command:        cwi,
				WithCompletion: cwi,
			}
		},
	)
}

func (c Show) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
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

func (c Show) RunWithIds(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	switch c.Gattung {

	case gattung.Akte:
		return c.showAkten(store, ids)

	case gattung.Zettel:
		fv := zettel.MakeFormatValue(store.Out(), store.Konfig())

		if err = fv.Set(c.Format); err != nil {
			err = errors.Normal(err)
			return
		}

		return c.showZettels(store, ids, fv)

	case gattung.Transaktion:
		return c.showTransaktions(store, ids)

	case gattung.Typ:
		var ev typ.FormatterValue

		if err = ev.Set(c.Format); err != nil {
			err = errors.Normal(err)
			return
		}

		return c.showTypen(
			store,
			ids,
			ev.FuncFormatter(
				store.Out(),
				store.StoreObjekten(),
			),
		)

	case gattung.Konfig:
		var ev konfig.FormatterValue

		if err = ev.Set(c.Format); err != nil {
			err = errors.Normal(err)
			return
		}

		return c.showKonfig(
			store,
			ev.FuncFormatter(
				store.Out(),
				store.StoreObjekten(),
			),
		)

	default:
		err = errors.Errorf("unsupported Gattung: %s", c.Gattung)
		return
	}
}

func (c Show) showZettels(
	store *umwelt.Umwelt,
	ids id_set.Set,
	fv *zettel.FormatValue,
) (err error) {
	w := collections.MakeChain(
		zettel.WriterIds{
			Filter: id_set.Filter{
				Set: ids,
			},
		}.WriteZettelTransacted,
		store.StoreWorkingDirectory().ZettelTransactedWriter(
			zettel.MakeWriterZettel(
				zettel.MakeSerializedFormatWriter(
					fv.Format,
					store.Out(),
					store.StoreObjekten(),
					store.Konfig(),
				),
			),
		),
	)

	if err = store.StoreObjekten().Zettel().ReadAllSchwanzenTransacted(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Show) showAkten(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	zettels := make([]zettel.Transacted, ids.Len())

	for i, is := range ids.AnyShasOrHinweisen() {
		var tz zettel.Transacted

		if tz, err = store.StoreObjekten().Zettel().ReadOne(is); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i] = tz
	}

	var ar io.ReadCloser

	for _, named := range zettels {
		if ar, err = store.StoreObjekten().AkteReader(named.Objekte.Akte); err != nil {
			err = errors.Wrap(err)
			return
		}

		if ar == nil {
			err = errors.Errorf("akte reader is nil")
			return
		}

		defer errors.Deferred(&err, ar.Close)

		if _, err = io.Copy(store.Out(), ar); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Show) showTransaktions(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	for _, is := range ids.Timestamps() {
		var t *transaktion.Transaktion

		if t, err = store.StoreObjekten().ReadTransaktion(is); err != nil {
			errors.PrintErrf("error: %s", err)
			continue
		}

		errors.Out().Printf("transaktion: %#v", t)

		t.Each(
			func(o *sku.Sku) (err error) {
				errors.Out().Print(o)
				return
			},
		)
	}

	return
}

func (c Show) showTypen(
	u *umwelt.Umwelt,
	ids id_set.Set,
	f collections.WriterFunc[*typ.Transacted],
) (err error) {
	f1 := collections.MakeSyncSerializer(f)

	typen := typ.MakeMutableSet(ids.Typen()...)
	if err = typen.EachPtr(
		collections.MakeChain(
			func(t *kennung.Typ) (err error) {
				//TODO-P2 move to store_objekten
				ty := u.Konfig().GetTyp(t.String())

				if ty == nil {
					return
				}

				if err = f1(ty); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Show) showKonfig(
	u *umwelt.Umwelt,
	f collections.WriterFunc[*konfig.Transacted],
) (err error) {
	f1 := collections.MakeSyncSerializer(f)

	var k *konfig.Transacted

	if k, err = u.StoreObjekten().Konfig().Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f1(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
