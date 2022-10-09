package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/transaktion"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/mike/umwelt"
)

type Show struct {
	gattung.Gattung
}

func init() {
	registerCommand(
		"show",
		func(f *flag.FlagSet) Command {
			c := &Show{
				Gattung: gattung.Zettel,
			}

			f.Var(&c.Gattung, "gattung", "Gattung")

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
					h, err = u.StoreObjekten().ExpandHinweisString(v)
					out = h.String()
					return
				},
			},
			id_set.ProtoId{
				MutableId: &etikett.Etikett{},
				Expand: func(v string) (out string, err error) {
					var e etikett.Etikett
					e, err = u.StoreObjekten().ExpandEtikettString(v)
					out = e.String()
					return
				},
			},
			id_set.ProtoId{
				MutableId: &typ.Typ{},
			},
			id_set.ProtoId{
				MutableId: &ts.Time{},
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
		return c.showZettels(store, ids)

	case gattung.Transaktion:
		return c.showTransaktions(store, ids)

	default:
		err = errors.Errorf("unsupported Gattung: %s", c.Gattung)
		return
	}
}

func (c Show) showZettels(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	zts := zettel_transacted.MakeSetUnique(0)
	w := zettel_transacted.MakeWriterMulti(
		nil,
		zettel_transacted.WriterZettelNamed{
			Writer: zettel_named.WriterFilter{
				NamedFilter: zettel_named.FilterIdSet{
					Set: ids,
				},
			},
		},
		zts,
	)

	if err = store.StoreObjekten().ReadAllSchwanzen(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	f := zettel.Text{}

	ctx := zettel.FormatContextWrite{
		Out:               store.Out(),
		AkteReaderFactory: store.StoreObjekten(),
	}

	zts.Each(
		func(zt zettel_transacted.Zettel) (err error) {
			if typKonfig, ok := store.Konfig().Typen[zt.Named.Stored.Zettel.Typ.String()]; ok {
				ctx.IncludeAkte = typKonfig.InlineAkte
			} else {
				ctx.IncludeAkte = zt.Named.Stored.Zettel.Typ.String() == "md"
			}

			ctx.Zettel = zt.Named.Stored.Zettel

			if _, err = f.WriteTo(ctx); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)

	return
}

func (c Show) showAkten(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	zettels := make([]zettel_transacted.Zettel, ids.Len())

	for i, is := range ids.AnyShasOrHinweisen() {
		var tz zettel_transacted.Zettel

		if tz, err = store.StoreObjekten().ReadOne(is); err != nil {
			err = errors.Wrap(err)
			return
		}

		zettels[i] = tz
	}

	var ar io.ReadCloser

	for _, named := range zettels {
		if ar, err = store.StoreObjekten().AkteReader(named.Named.Stored.Zettel.Akte); err != nil {
			err = errors.Wrap(err)
			return
		}

		if ar == nil {
			err = errors.Errorf("akte reader is nil")
			return
		}

		defer errors.PanicIfError(ar.Close)

		if _, err = io.Copy(store.Out(), ar); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Show) showTransaktions(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	for _, is := range ids.Timestamps() {
		var t transaktion.Transaktion

		if t, err = store.StoreObjekten().ReadTransaktion(is); err != nil {
			errors.PrintErrf("error: %s", err)
			continue
		}

		errors.PrintOutf("transaktion: %#v", t)

		for _, o := range t.Objekten {
			errors.Out().Print(o)
		}
	}

	return
}
