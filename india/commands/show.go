package commands

import (
	"flag"
	"io"
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

			return commandWithZettels{c}
		},
	)
}

func (c Show) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	zettels := make([]_NamedZettel, len(args))

	for i, a := range args {
		var h _Hinweis

		if h, err = _MakeBlindHinweis(a); err != nil {
			err = _Error(err)
			return
		}

		var named _NamedZettel

		if named, err = zs.Read(h); err != nil {
			err = _Error(err)
			return
		}

		zettels[i] = named
	}

	switch c.Type {

	case _TypeAkte:
		return c.showAkten(u, zs, zettels)

	case _TypeZettel:
		return c.showZettels(u, zs, zettels)

	default:
		err = _Errorf("unsupported objekte type: %s", c.Type)
		return
	}

	return
}

func (c Show) showZettels(u _Umwelt, zs _Zettels, zettels []_NamedZettel) (err error) {
	f := _ZettelFormatsText{}

	ctx := _ZettelFormatContextWrite{
		Out:               u.Out,
		AkteReaderFactory: zs,
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

func (c Show) showAkten(u _Umwelt, zs _Zettels, zettels []_NamedZettel) (err error) {
	var ar io.ReadCloser

	for _, named := range zettels {
		if ar, err = zs.AkteReader(named.Zettel.Akte); err != nil {
			err = _Error(err)
			return
		}

		if ar == nil {
			err = _Errorf("akte reader is nil")
			return
		}

		defer _PanicIfError(ar.Close())

		if _, err = io.Copy(u.Out, ar); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
