package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type Show struct {
	Format string
	// Filter lua.Flag
}

func init() {
	registerCommandWithQuery(
		"show",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Show{}

			f.StringVar(&c.Format, "format", "log", "format")
			// f.Var(&c.Filter, "filter", "lua filter")

			return c
		},
	)
}

func (c Show) CompletionGattung() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Show) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
		// gattung.Etikett,
		// gattung.Typ,
		// gattung.Bestandsaufnahme,
		// gattung.Kasten,
	)
}

func (c Show) RunWithQuery(u *umwelt.Umwelt, ms *query.Group) (err error) {
	var f schnittstellen.FuncIter[*sku.Transacted]

	if f, err = u.MakeFormatFunc(c.Format, u.Out()); err != nil {
		err = errors.Wrap(err)
		return
	}

	f = iter.MakeSyncSerializer(f)

	if err = u.GetStore().QueryWithCwd(
		ms,
		func(sk *sku.Transacted) (err error) {
			// if c.Filter.String() != "" {
			// 	if !c.Filter.ContainsSku(sk) {
			// 		return
			// 	}
			// }

			if err = f(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
