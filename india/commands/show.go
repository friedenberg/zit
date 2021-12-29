package commands

import (
	"flag"
)

type Show struct {
}

func init() {
	registerCommand(
		"show",
		func(f *flag.FlagSet) Command {
			c := &Show{}

			return commandWithZettels{c}
		},
	)
}

func (c Show) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	f := _ZettelFormatsText{}

	ctx := _ZettelFormatContextWrite{
		Out:               u.Out,
		AkteReaderFactory: zs,
	}

	for _, a := range args {
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

		ctx.IncludeAkte = named.Zettel.AkteExt.String() == "md"

		// 		if !ctx.IncludeAkte {
		// 			v := fmt.Sprintf(
		// 				"%s.%s",
		// 				named.Zettel.Akte.String(),
		// 				named.Zettel.AkteExt.String(),
		// 			)

		// 			if err = named.Zettel.AkteExt.Set(v); err != nil {
		// 				err = _Error(err)
		// 				return
		// 			}
		// 		}

		ctx.Zettel = named.Zettel

		if _, err = f.WriteTo(ctx); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
