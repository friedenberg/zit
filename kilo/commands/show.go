package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/india/store_with_lock"
)

type Show struct {
	Type _Type
}

func init() {
	registerCommand(
		"show",
		func(f *flag.FlagSet) Command {
			c := &Show{
				Type: _TypeZettel,
			}

			f.Var(&c.Type, "type", "ObjekteType")

			return commandWithLockedStore{c}
		},
	)
}

func (c Show) RunWithLockedStore(store store_with_lock.Store, args ...string) (err error) {
	zettels := make([]_NamedZettel, len(args))

	for i, a := range args {
		var h _Hinweis

		if h, err = _MakeBlindHinweis(a); err != nil {
			err = _Error(err)
			return
		}

		var named _NamedZettel

		if named, err = store.Zettels().Read(h); err != nil {
			err = _Error(err)
			return
		}

		zettels[i] = named
	}

	switch c.Type {

	case _TypeAkte:
		return c.showAkten(store, zettels)

	case _TypeZettel:
		return c.showZettels(store, zettels)

	default:
		err = _Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}

func (c Show) showZettels(store store_with_lock.Store, zettels []_NamedZettel) (err error) {
	f := _ZettelFormatsText{}

	ctx := _ZettelFormatContextWrite{
		Out:               store.Out,
		AkteReaderFactory: store.Zettels(),
	}

	for _, named := range zettels {
		ctx.IncludeAkte = named.Zettel.AkteExt.String() == "md"

		ctx.Zettel = named.Zettel

		if _, err = f.WriteTo(ctx); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

func (c Show) showAkten(store store_with_lock.Store, zettels []_NamedZettel) (err error) {
	var ar io.ReadCloser

	for _, named := range zettels {
		if ar, err = store.Zettels().AkteReader(named.Zettel.Akte); err != nil {
			err = _Error(err)
			return
		}

		if ar == nil {
			err = _Errorf("akte reader is nil")
			return
		}

		defer _PanicIfError(ar.Close())

		if _, err = io.Copy(store.Out, ar); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
