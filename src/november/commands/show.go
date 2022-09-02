package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/echo/transaktion"
	"github.com/friedenberg/zit/src/echo/zettel"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/store_objekten"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
)

type Show struct {
	Type zk_types.Type
}

func init() {
	registerCommand(
		"show",
		func(f *flag.FlagSet) Command {
			c := &Show{
				Type: zk_types.TypeZettel,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{commandWithId{c}}
		},
	)
}

func (c Show) RunWithId(store store_with_lock.Store, ids ...id_set.Set) (err error) {
	switch c.Type {

	case zk_types.TypeAkte:
		return c.showAkten(store, ids)

	case zk_types.TypeZettel:
		return c.showZettels(store, ids)

	case zk_types.TypeTransaktion:
		return c.showTransaktions(store, ids)

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}
}

func (c Show) showZettels(store store_with_lock.Store, ids []id_set.Set) (err error) {
	zettels := make([]zettel_transacted.Zettel, len(ids))

	for i, is := range ids {
		var idd id.Id
		ok := false

		if idd, ok = is.AnyShaOrHinweis(); !ok {
			errors.PrintErrf("unsupported id: %s", is)
			err = nil
			continue
		}

		var tz zettel_transacted.Zettel

		if tz, err = store.StoreObjekten().Read(idd); err != nil {
			if errors.Is(err, store_objekten.ErrNotFound{}) {
				err = errors.Normal(err)
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		zettels[i] = tz
	}

	f := zettel.Text{}

	ctx := zettel.FormatContextWrite{
		Out:               store.Out,
		AkteReaderFactory: store.StoreObjekten(),
	}

	for _, named := range zettels {
		if typKonfig, ok := store.Konfig.Typen[named.Named.Stored.Zettel.Typ.String()]; ok {
			ctx.IncludeAkte = typKonfig.InlineAkte
		} else {
			ctx.IncludeAkte = named.Named.Stored.Zettel.Typ.String() == "md"
		}

		ctx.Zettel = named.Named.Stored.Zettel

		if _, err = f.WriteTo(ctx); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Show) showAkten(store store_with_lock.Store, ids []id_set.Set) (err error) {
	zettels := make([]zettel_transacted.Zettel, len(ids))

	for i, is := range ids {
		var idd id.Id
		ok := false

		if idd, ok = is.AnyShaOrHinweis(); !ok {
			errors.PrintErrf("unsupported id: %s", is)
			err = nil
			continue
		}

		var tz zettel_transacted.Zettel

		if tz, err = store.StoreObjekten().Read(idd); err != nil {
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

		if _, err = io.Copy(store.Out, ar); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Show) showTransaktions(store store_with_lock.Store, ids []id_set.Set) (err error) {
	for _, is := range ids {
		var idd id.Id
		ok := false

		if idd, ok = is.Any(&ts.Time{}); !ok {
			errors.PrintErrf("unsupported id: %s", is)
			err = nil
			continue
		}

		tid := idd.(ts.Time)

		var t transaktion.Transaktion

		if t, err = store.StoreObjekten().ReadTransaktion(tid); err != nil {
			errors.PrintErrf("%s", err)
			continue
		}

		errors.PrintOutf("%#v", t)
	}

	return
}
