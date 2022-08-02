package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/node_type"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type Show struct {
	Type node_type.Type
}

func init() {
	registerCommand(
		"show",
		func(f *flag.FlagSet) Command {
			c := &Show{
				Type: node_type.TypeZettel,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{commandWithId{c}}
		},
	)
}

func (c Show) RunWithId(store store_with_lock.Store, ids ...id.Id) (err error) {
	zettels := make([]stored_zettel.Named, len(ids))

	for i, a := range ids {
		var named stored_zettel.Named

		if named, err = store.Zettels().Read(a); err != nil {
			err = errors.Error(err)
			return
		}

		zettels[i] = named
	}

	switch c.Type {

	case node_type.TypeAkte:
		return c.showAkten(store, zettels)

	case node_type.TypeZettel:
		return c.showZettels(store, zettels)

	default:
		err = errors.Errorf("unsupported objekte type: %s", c.Type)
		return
	}
}

func (c Show) showZettels(store store_with_lock.Store, zettels []stored_zettel.Named) (err error) {
	f := zettel_formats.Text{}

	ctx := zettel.FormatContextWrite{
		Out:               store.Out,
		AkteReaderFactory: store.Zettels(),
	}

	for _, named := range zettels {
		ctx.IncludeAkte = named.Zettel.AkteExt.String() == "md"

		ctx.Zettel = named.Zettel

		if _, err = f.WriteTo(ctx); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (c Show) showAkten(store store_with_lock.Store, zettels []stored_zettel.Named) (err error) {
	var ar io.ReadCloser

	for _, named := range zettels {
		if ar, err = store.Zettels().AkteReader(named.Zettel.Akte); err != nil {
			err = errors.Error(err)
			return
		}

		if ar == nil {
			err = errors.Errorf("akte reader is nil")
			return
		}

		defer stdprinter.PanicIfError(ar.Close())

		if _, err = io.Copy(store.Out, ar); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
