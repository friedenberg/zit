package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/id_set"
	"github.com/friedenberg/zit/src/echo/transaktion"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	"github.com/friedenberg/zit/src/india/store_objekten"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
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
	zettels := make([]zettel_transacted.Transacted, len(ids))

	for i, is := range ids {
		var idd id.Id
		ok := false

		if idd, ok = is.AnyShaOrHinweis(); !ok {
			stdprinter.Errf("unsupported id: %s", is)
			err = nil
			continue
		}

		var tz zettel_transacted.Transacted

		if tz, err = store.StoreObjekten().Read(idd); err != nil {
			if errors.Is(err, store_objekten.ErrNotFound{}) {
				err = errors.Normal(err)
			} else {
				err = errors.Error(err)
			}

			return
		}

		zettels[i] = tz
	}

	f := zettel_formats.Text{}

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
			err = errors.Error(err)
			return
		}
	}

	return
}

func (c Show) showAkten(store store_with_lock.Store, ids []id_set.Set) (err error) {
	zettels := make([]zettel_transacted.Transacted, len(ids))

	for i, is := range ids {
		var idd id.Id
		ok := false

		if idd, ok = is.AnyShaOrHinweis(); !ok {
			stdprinter.Errf("unsupported id: %s", is)
			err = nil
			continue
		}

		var tz zettel_transacted.Transacted

		if tz, err = store.StoreObjekten().Read(idd); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = tz
	}

	var ar io.ReadCloser

	for _, named := range zettels {
		if ar, err = store.StoreObjekten().AkteReader(named.Named.Stored.Zettel.Akte); err != nil {
			err = errors.Error(err)
			return
		}

		if ar == nil {
			err = errors.Errorf("akte reader is nil")
			return
		}

		defer stdprinter.PanicIfError(ar.Close)

		if _, err = io.Copy(store.Out, ar); err != nil {
			err = errors.Error(err)
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
			stdprinter.Errf("unsupported id: %s", is)
			err = nil
			continue
		}

		tid := idd.(ts.Time)

		var t transaktion.Transaktion

		if t, err = store.StoreObjekten().ReadTransaktion(tid); err != nil {
			stdprinter.Errf("%s\n", err)
			continue
		}

		stdprinter.Outf("%#v\n", t)
	}

	return
}
