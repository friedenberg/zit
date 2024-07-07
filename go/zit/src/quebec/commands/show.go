package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Show struct {
	Format string
}

func init() {
	registerCommandWithExternalQuery(
		"show",
		func(f *flag.FlagSet) CommandWithExternalQuery {
			c := &Show{}

			f.StringVar(&c.Format, "format", "log", "format")

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
	)
}

func (c Show) RunWithExternalQuery(
	u *umwelt.Umwelt,
	eqwk sku.ExternalQuery,
) (err error) {
	var f schnittstellen.FuncIter[*sku.Transacted]

	if f, err = u.MakeFormatFunc(c.Format, u.Out()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryWithKasten(
		eqwk,
		iter.MakeSyncSerializer(f),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
