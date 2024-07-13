package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Show struct {
	Format string
}

func init() {
	registerCommandWithQuery(
		"show",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Show{}

			f.StringVar(&c.Format, "format", "log", "format")

			return c
		},
	)
}

func (c Show) CompletionGattung() ids.Genre {
	return ids.MakeGenre(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
		gattung.Bestandsaufnahme,
		gattung.Kasten,
	)
}

func (c Show) DefaultGattungen() ids.Genre {
	return ids.MakeGenre(
		gattung.Zettel,
	)
}

func (c Show) RunWithQuery(
	u *umwelt.Umwelt,
	eqwk *query.Group,
) (err error) {
	var f interfaces.FuncIter[*sku.Transacted]

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
