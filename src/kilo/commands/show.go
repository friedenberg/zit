package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	zettel_stored "github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	"github.com/friedenberg/zit/src/india/store_with_lock"
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

func (c Show) RunWithId(store store_with_lock.Store, ids ...id.Id) (err error) {
	zettels := make([]zettel_stored.Transacted, len(ids))

	for i, a := range ids {
		var tz zettel_stored.Transacted

		if tz, err = store.StoreObjekten().Read(a); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = tz
	}

	switch c.Type {

	case zk_types.TypeAkte:
		return c.showAkten(store, zettels)

	case zk_types.TypeZettel:
		return c.showZettels(store, zettels)

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}
}

func (c Show) showZettels(store store_with_lock.Store, zettels []zettel_stored.Transacted) (err error) {
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

func (c Show) showAkten(store store_with_lock.Store, zettels []zettel_stored.Transacted) (err error) {
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
