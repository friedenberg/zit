package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/transaktion"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type Show struct {
	Type gattung.Gattung
}

func init() {
	registerCommand(
		"show",
		func(f *flag.FlagSet) Command {
			c := &Show{
				Type: gattung.Zettel,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c Show) RunWithIds(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	switch c.Type {

	case gattung.Akte:
		return c.showAkten(store, ids)

	case gattung.Zettel:
		return c.showZettels(store, ids)

	case gattung.Transaktion:
		return c.showTransaktions(store, ids)

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}
}

func (c Show) showZettels(store *umwelt.Umwelt, ids id_set.Set) (err error) {
	var zts zettel_transacted.Set

	if zts, err = store.StoreObjekten().ReadMany(ids); err != nil {
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
			errors.PrintErrf("%s", err)
			continue
		}

		errors.PrintOutf("%#v", t)
	}

	return
}
